"""Read Spring route annotations directly from original Boot JAR class files.

Does not load Java classes, execute handlers, start services, or read JAR configs.
Usage: python3 tools/validation/jar_routes.py --source-root PATH --output PATH
"""
import argparse
import hashlib
import io
import itertools
import json
import struct
import zipfile
from pathlib import Path


class ClassFile:
    def __init__(self, blob):
        self.f = io.BytesIO(blob)
        if self.f.read(4) != b"\xca\xfe\xba\xbe":
            raise ValueError("not a Java class")
        self.f.read(4)
        count = self.u2()
        self.cp = [None] * count
        i = 1
        while i < count:
            tag = self.u1()
            if tag == 1:
                self.cp[i] = self.f.read(self.u2()).decode("utf-8", "replace")
            elif tag in (3, 4):
                self.cp[i] = struct.unpack(">i" if tag==3 else ">f",self.f.read(4))[0]
            elif tag in (5, 6):
                self.f.read(8)
                i += 1
            elif tag in (7, 8, 16, 19, 20):
                self.cp[i] = self.u2()
            elif tag in (9, 10, 11, 12, 17, 18):
                self.f.read(4)
            elif tag == 15:
                self.f.read(3)
            else:
                raise ValueError(f"unknown constant tag {tag}")
            i += 1
        self.f.read(6)
        self.f.read(self.u2() * 2)
        self.members()  # fields
        self.methods = self.members()
        self.annotations = self.attributes()

    def u1(self):
        return struct.unpack(">B", self.f.read(1))[0]

    def u2(self):
        return struct.unpack(">H", self.f.read(2))[0]

    def value(self):
        tag = chr(self.u1())
        if tag in "BCDFIJSZs":
            return self.cp[self.u2()]
        if tag == "e":
            return [self.cp[self.u2()], self.cp[self.u2()]]
        if tag == "c":
            return self.cp[self.u2()]
        if tag == "@":
            return self.annotation()
        if tag == "[":
            return [self.value() for _ in range(self.u2())]
        raise ValueError(f"unknown annotation value tag {tag}")

    def annotation(self):
        name = self.cp[self.u2()]
        pairs = {}
        for _ in range(self.u2()):
            key = self.cp[self.u2()]
            pairs[key] = self.value()
        return name, pairs

    def attributes(self):
        annotations = []
        for _ in range(self.u2()):
            name = self.cp[self.u2()]
            length = struct.unpack(">I", self.f.read(4))[0]
            end = self.f.tell() + length
            if name == "RuntimeVisibleAnnotations":
                annotations += [self.annotation() for _ in range(self.u2())]
            elif name == "RuntimeVisibleParameterAnnotations":
                for index in range(self.u1()):
                    items=[self.annotation() for _ in range(self.u2())]
                    annotations.append(("__parameter__",{"index":index,"annotations":items}))
            self.f.seek(end)
        return annotations

    def members(self):
        members = []
        for _ in range(self.u2()):
            self.u2()  # access
            name = self.cp[self.u2()]
            descriptor = self.cp[self.u2()]
            members.append((name, descriptor, self.attributes()))
        return members


PREFIX = "Lorg/springframework/web/bind/annotation/"


def mappings(annotations):
    for name, values in annotations:
        kind = name.removeprefix(PREFIX).removesuffix(";")
        if not name.startswith(PREFIX) or kind not in (
            "RequestMapping", "GetMapping", "PostMapping", "PutMapping",
            "DeleteMapping", "PatchMapping",
        ):
            continue
        paths = values.get("value", values.get("path", [""]))
        methods = [x[1] for x in values.get("method", [])]
        if not methods and kind != "RequestMapping":
            methods = [kind.removesuffix("Mapping").upper()]
        yield paths, methods


def class_routes(blob, service):
    cf = ClassFile(blob)
    names = {name for name, _ in cf.annotations}
    if not names.intersection({PREFIX + "RestController;", "Lorg/springframework/stereotype/Controller;"}):
        return []
    bases = list(mappings(cf.annotations)) or [([""], [])]
    routes = []
    for name, descriptor, annotations in cf.methods:
        parameters=__import__('re').findall(r'\[*(?:L[^;]+;|[BCDFIJSZ])',descriptor.split(')',1)[0][1:])
        body=None
        for annotation,details in annotations:
            if annotation != '__parameter__':continue
            for an,values in details['annotations']:
                if an==PREFIX+'RequestBody;':
                    typ=parameters[details['index']]
                    kind='object'
                    if typ=='Ljava/util/List;' or typ.startswith('['):kind='array'
                    if typ in ('Ljava/lang/String;','Ljava/lang/Object;'):kind='raw'
                    body={'kind':kind,'required':bool(values.get('required',True))}
        for (base_paths, base_methods), (paths, methods) in itertools.product(bases, mappings(annotations)):
            for base, path, method in itertools.product(base_paths, paths, methods or base_methods or ["ANY"]):
                full = "/" + "/".join(x.strip("/") for x in (service, base, path) if x.strip("/"))
                routes.append({"method": method, "path": full, "java_method": name, "descriptor": descriptor,"body":body})
    return routes


def inventory(root):
    routes, jars = [], []
    for service in ("auth", "common", "contract", "logistics", "capital", "bill", "warehouse"):
        module = "saas-trade-business" if service == "warehouse" else "saas-supply-chain"
        target = root / module / ("saas-" + service) / "target"
        paths = [p for p in target.glob("*.jar") if p.name != "original.jar"]
        if len(paths) != 1:
            raise ValueError(f"expected one executable JAR in {target}")
        path = paths[0]
        jars.append({"path": str(path.relative_to(root)), "sha256": hashlib.file_digest(path.open("rb"), "sha256").hexdigest()})
        with zipfile.ZipFile(path) as jar:
            for member in sorted(jar.namelist()):
                if not member.startswith("BOOT-INF/classes/") or not member.endswith(".class"):
                    continue
                # Controllers are annotated classes; filenames need not end in Controller.
                blob = jar.read(member)
                if b"RestController" not in blob and b"stereotype/Controller" not in blob:
                    continue
                for route in class_routes(blob, "saas-" + service):
                    routes.append({**route, "class": member, "jar": str(path.relative_to(root))})
    return {"jars": jars, "routes": routes}


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--source-root", required=True, type=Path)
    parser.add_argument("--output", required=True, type=Path)
    parser.add_argument("--go-body-contracts",type=Path)
    args = parser.parse_args()
    result = inventory(args.source_root)
    args.output.parent.mkdir(parents=True, exist_ok=True)
    args.output.write_text(json.dumps(result, ensure_ascii=False, indent=2) + "\n")
    if args.go_body_contracts:
        lines=['// Code generated by tools/validation/jar_routes.py from original JAR annotations. DO NOT EDIT.','package platform','var javaBodyContracts=map[string]bodyContract{']
        for r in sorted(result['routes'],key=lambda r:(r['method'],r['path'])):
            if not r['body']:continue
            key=r['method']+' '+r['path']
            lines.append(json.dumps(key)+':{kind:'+json.dumps(r['body']['kind'])+',required:'+str(r['body']['required']).lower()+'},')
        lines.append('}')
        args.go_body_contracts.write_text('\n'.join(lines)+'\n')
    print(f"{len(result['jars'])} JARs, {len(result['routes'])} method/path mappings")
