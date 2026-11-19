"""Strict, local-only, read response comparison without storing business payloads.

Session file is private JSON: {account: {port: {token: ...}}}. Never prints tokens,
profile data, passwords or response bodies. Numeric strings differ from numbers;
null differs from absent keys; array order and decimal spelling are preserved.
"""
import argparse
import ast
import hashlib
import io
import json
import os
import sys
import urllib.error
import urllib.request
import zipfile
import xml.etree.ElementTree as ET
from pathlib import Path


class Number(str):
    pass


def decode(blob):
    return json.loads(blob, parse_int=Number, parse_float=Number)


def differences(a, b, path="$"):
    if type(a) is not type(b):
        return [path + ":type"]
    if isinstance(a, dict):
        out = [path + "." + k + ":missing" for k in sorted(a.keys() ^ b.keys())]
        for k in sorted(a.keys() & b.keys()):
            out.extend(differences(a[k], b[k], path + "." + k))
            if len(out) >= 30:
                break
        return out[:30]
    if isinstance(a, list):
        if len(a) != len(b):
            return [path + ":length"]
        out = []
        for i, (x, y) in enumerate(zip(a, b)):
            out.extend(differences(x, y, f"{path}[{i}]"))
            if len(out) >= 30:
                break
        return out[:30]
    return [] if a == b else [path + ":value"]


def spreadsheet(blob):
    ns = {"m": "http://schemas.openxmlformats.org/spreadsheetml/2006/main"}
    with zipfile.ZipFile(io.BytesIO(blob)) as z:
        shared = []
        if "xl/sharedStrings.xml" in z.namelist():
            shared = ["".join(x.itertext()) for x in ET.fromstring(z.read("xl/sharedStrings.xml"))]
        result = []
        for name in sorted(x for x in z.namelist() if x.startswith("xl/worksheets/sheet") and x.endswith(".xml")):
            rows = []
            for row in ET.fromstring(z.read(name)).findall("m:sheetData/m:row", ns):
                cells = []
                for c in row:
                    value = c.find("m:v", ns)
                    value = value.text if value is not None else ""
                    if c.get("t") == "s":
                        value = shared[int(value)]
                    elif c.get("t") == "inlineStr":
                        value = "".join(c.find("m:is", ns).itertext())
                    cells.append([c.get("r"), value])
                rows.append(cells)
            result.append(rows)
        return result


READ_NAMES = {"list", "getById", "listByPage", "exportExcel"}


def matrix(root):
    """Explicit read allowlist; do not infer safety from HTTP method or a get prefix."""
    routes = json.loads((root / "docs/validation/java-go-20261118/java-jar-routes.json").read_text())["routes"]
    cases = []
    for r in routes:
        path = r["path"]
        end = path.rsplit("/", 1)[-1]
        # Legacy user detail lazily persists a nickname; test it with isolated fixtures.
        if path == "/saas-auth/chainAuthUser/getById":
            continue
        if end not in READ_NAMES or r["method"] not in ("POST", "GET"):
            continue
        bodies = {"list": [{"id": 99977799}, {"id": 1}], "getById": [{"id": 99977799}, {"id": 1}, {}],
                  "listByPage": [{}, {"searchStr": "", "page": {"current": 1, "size": 3}},
                                 {"searchStr": "1", "page": {"current": 1, "size": 3}}],
                  "exportExcel": [{}]}[end]
        cases.extend({"method": r["method"], "path": path, "body": b} for b in bodies)
    return cases


def request(port, token, case, applytype=""):
    headers = {"Content-Type": "application/json"}
    if token:
        headers["Authorization"] = "Bearer " + token
    if applytype:
        headers["applytype"] = applytype
    req = urllib.request.Request(f"http://127.0.0.1:{port}" + case["path"],
                                 data=json.dumps(case["body"]).encode() if case["method"] != "GET" else None,
                                 method=case["method"], headers=headers)
    try:
        response = urllib.request.urlopen(req, timeout=30)
    except urllib.error.HTTPError as e:
        response = e
    with response:
        return response.status, dict(response.headers), response.read()


def compare(a, b):
    out = []
    if a[0] != b[0]:
        out.append("HTTP:status")
    try:
        ja, jb = decode(a[2]), decode(b[2])
    except (ValueError, UnicodeDecodeError):
        if a[2].startswith(b"PK") and b[2].startswith(b"PK"):
            out.extend(differences(spreadsheet(a[2]), spreadsheet(b[2]), "$workbook"))
        elif a[2] != b[2]:
            out.append("$raw:bytes")
    else:
        out.extend(differences(ja, jb))
    # Ignore charset spelling/order, but do not hide JSON vs binary mismatches.
    media = lambda r: next((v.split(";")[0].strip().lower() for k, v in r[1].items() if k.lower() == "content-type"), "")
    if media(a) != media(b):
        out.append("HTTP:content-type")
    return out


def main():
    p = argparse.ArgumentParser(description=__doc__)
    p.add_argument("--sessions", type=Path, required=True)
    p.add_argument("--output", type=Path, required=True)
    p.add_argument("--account", default="ptadmin")
    p.add_argument("--service", default="")
    p.add_argument("--case-file", type=Path)
    p.add_argument("--applytype", default="")
    p.add_argument("--anonymous", action="store_true")
    p.add_argument("--include-export", action="store_true")
    p.add_argument("--go-port", default="8799")
    args = p.parse_args()
    root = Path(__file__).resolve().parents[2]
    sessions = json.loads(args.sessions.read_text())[args.account]
    cases = json.loads(args.case_file.read_text()) if args.case_file else matrix(root)
    cases = [c for c in cases if args.service in c["path"] and (args.include_export or not c["path"].endswith("/exportExcel"))]
    result = []
    for i, c in enumerate(cases):
        if c["path"].rsplit("/", 1)[-1] not in READ_NAMES:
            raise ValueError("case outside reviewed read allowlist: " + c["path"])
        entry = {"case": c, "account": "anonymous" if args.anonymous else args.account, "applytype": args.applytype}
        try:
            a, b = [request(port, "" if args.anonymous else sessions[port]["token"], c, args.applytype) for port in ("8765", args.go_port)]
            diff = compare(a, b)
            entry.update(java_status=a[0], go_status=b[0], differences=diff,
                         java_sha256=hashlib.sha256(a[2]).hexdigest(), go_sha256=hashlib.sha256(b[2]).hexdigest())
            for name, response in (("java", a), ("go", b)):
                try:
                    obj = decode(response[2])
                    entry[name + "_code"] = str(obj.get("code")) if isinstance(obj, dict) else None
                except (ValueError, UnicodeDecodeError):
                    pass
        except Exception as e:
            entry["differences"] = ["transport-or-decode:" + type(e).__name__]
        result.append(entry)
        if (i + 1) % 100 == 0:
            print(f"{i+1}/{len(cases)}, mismatches={sum(bool(x['differences']) for x in result)}", flush=True)
    args.output.parent.mkdir(parents=True, exist_ok=True)
    args.output.write_text(json.dumps(result, ensure_ascii=False, indent=2) + "\n")
    print(f"RESULT {len(result) - sum(bool(x['differences']) for x in result)}/{len(result)} equal", flush=True)
    return int(any(x["differences"] for x in result))


if __name__ == "__main__":
    sys.exit(main())
