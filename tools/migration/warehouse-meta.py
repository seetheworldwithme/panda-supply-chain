# 六期 Go 元数据生成器:解析 67 实体 + DB information_schema → internal/warehouse/meta_gen.go
import re, os, glob, json, subprocess, collections
from pathlib import Path

BASE = "/Volumes/coding/HBuilderxProject/panda/saas-hyqc/saas-trade-business/saas-warehouse/src/main/java/com/saas/warehouse"
CORE = "/Volumes/coding/HBuilderxProject/panda/saas-hyqc/saas-trade-business/saas-warehouse-core/src/main/java/com/saas/warehouse"
OUT = str(Path(__file__).resolve().parents[2] / "backend/internal/warehouse/meta_gen.go")

def read(p):
    return open(p, encoding="utf-8", errors="replace").read()

def camel_to_underline(s):
    out = []
    for i, c in enumerate(s):
        if c.isupper():
            if i > 0:
                out.append('_')
            out.append(c.lower())
        else:
            out.append(c)
    return ''.join(out)

# ── mapper listByPage/listNoPage SQL 提取 ──
MAPDIR = BASE + "/business/mapper"

def mapper_sqls(name):
    p = os.path.join(MAPDIR, name + "Mapper.java")
    if not os.path.exists(p):
        return {}, {}
    src = read(p)
    out = {}
    for m in re.finditer(r'@Select\("((?:[^"\\]|\\.)*)"\)\s*(?:@\w+(?:\([^)]*\))?\s*)*[\w<>.,\s\[\]?]+?\s+(\w+)\(', src):
        out[m.group(2)] = m.group(1)
    return out

# ── mapper listByPage/listNoPage SQL 提取 ──
MAPDIR = BASE + "/business/mapper"

def mapper_sqls(name):
    p = os.path.join(MAPDIR, name + "Mapper.java")
    if not os.path.exists(p):
        return {}
    src = read(p)
    out = {}
    for m in re.finditer(r'@Select\("((?:[^"\\]|\\.)*)"\)\s*(?:@\w+(?:\([^)]*\))?\s*)*[\w<>.,\s\[\]?]+?\s+(\w+)\(', src):
        out[m.group(2)] = m.group(1)
    return out

# ── DB 列与类型 ──
def db_schema():
    q = "SELECT table_name, column_name, data_type FROM information_schema.columns WHERE table_schema='chain' ORDER BY table_name, ordinal_position"
    r = subprocess.run(["mysql", "-h127.0.0.1", "-uroot", "chain", "-N", "-B", "-e", q], capture_output=True, text=True)
    cols = collections.defaultdict(dict)  # table → col → datatype
    for line in r.stdout.splitlines():
        t, c, d = line.split('\t')
        cols[t][c] = d
    return cols

DB = db_schema()

KIND = {"String": "s", "Integer": "i", "Long": "l", "BigDecimal": "d",
        "Boolean": "b", "boolean": "b", "LocalDate": "D", "LocalDateTime": "t"}

def parse_entity(path):
    src = read(path)
    name = os.path.basename(path)[:-5]
    fields = []
    anns = []
    for line in src.split("\n"):
        t = line.strip()
        if t.startswith("@"):
            anns.append(t)
            continue
        m = re.match(r'(?:private|protected|public)\s+([\w.<>\[\], ]+?)\s+(\w+)\s*;', t)
        if m:
            ftype, fname = m.group(1).strip(), m.group(2)
            block = "\n".join(anns)
            excel = re.search(r'@ExcelColumn\(\s*value\s*=\s*"([^"]*)"', block)
            jf = '@JsonFormat' in block
            fm = re.search(r'FieldFill\.(\w+)', block)
            fill = fm.group(1) if fm else None
            notbank = '@NotBankColumn' in block
            syscol = '@SystemCodeColumn' in block
            exist_false = '@TableField' in block and ('exist = false' in block or 'exist=false' in block)
            if ftype in KIND:
                kind = KIND[ftype]
                if kind == "t" and not jf:
                    kind = "T"
            elif ftype.startswith("Page<"):
                kind = "page"
            elif ftype.startswith("List<"):
                kind = "list"
            else:
                kind = "x"
            fields.append({
                "name": fname, "type": ftype, "kind": kind, "excel": excel, "jf": jf,
                "fill": fill, "notbank": notbank, "sys": syscol,
                "carrier": exist_false or ftype not in KIND,
            })
            anns = []  # annotations belong only to the immediately following field
        elif t and not t.startswith("@") and not t.startswith("*") and not t.startswith("/*") and not t.startswith("//"):
            if fields or "{" in t:
                anns = []
    return name, fields

def golist(s):
    """SQL 字符串 → Go 反引号串(内嵌双引号场景安全)"""
    if '`' not in s and '\r' not in s and '\n' not in s:
        return '`' + s + '`'
    return json.dumps(s)

def emit():
    ents = []
    mismatches = []
    EXTRA_ENTS = ["/Volumes/coding/HBuilderxProject/panda/saas-hyqc/saas-hyqc/hyqc-core/src/main/java/com/saas/cores/model/ChainAuthUser.java",
                  "/Volumes/coding/HBuilderxProject/panda/saas-hyqc/saas-hyqc/hyqc-core/src/main/java/com/saas/cores/model/BusMemberBenefits.java",
                  "/Volumes/coding/HBuilderxProject/panda/saas-hyqc/saas-hyqc/hyqc-core/src/main/java/com/saas/cores/model/BusVideos.java"]
    for path in sorted(glob.glob(BASE + "/business/entity/*.java")) + sorted(glob.glob(CORE + "/core/model/*.java")) + [e for e in EXTRA_ENTS if os.path.exists(e)]:
        name, fields = parse_entity(path)
        table = camel_to_underline(name)
        dbcols = DB.get(table)
        msq = mapper_sqls(name)
        ents.append((name, table, fields, dbcols, msq))

    lines = []
    lines.append('// meta_gen.go 由 tools/migration/warehouse-meta.py 从 Java 实体+chain 库 information_schema 生成,勿手改。')
    lines.append('// 每实体:table/实体全键(keys)/getSql 候选(q)/Excel 列(excel)/DECIMAL(dates 等)响应整形清单。')
    lines.append('package warehouse')
    lines.append('')
    lines.append('type sqlf struct { Name, Kind string } // Kind: i/D 原值绑定,其余 %v% 包裹')
    lines.append('')
    lines.append('var ents = map[string]*ent{')
    for name, table, fields, dbcols, msq in ents:
        keys = [f["name"] for f in fields]
        noncarrier = [f for f in fields if not f["carrier"]]
        dbmap = {}
        if dbcols:
            colsset = set(dbcols)
            for f in noncarrier:
                cn = camel_to_underline(f["name"])
                dbmap[f["name"]] = cn if cn in colsset else cn
        # getSql 候选:有 @ExcelColumn 注解(含空头)的非载体字段
        q = [f for f in noncarrier if f["excel"]]
        # excel 导入导出列:@ExcelColumn 且 value 非空
        excel = [f for f in noncarrier if f["excel"] and f["excel"].group(1) != ""]
        decimals = [f["name"] for f in noncarrier if f["type"] == "BigDecimal"]
        dates = [f["name"] for f in noncarrier if f["type"] == "LocalDate"]
        iso = [f["name"] for f in noncarrier if f["type"] == "LocalDateTime" and not f["jf"]]
        bools = [f["name"] for f in noncarrier if f["type"] in ("Boolean", "boolean")]
        syscol = bool(dbcols and "system_code" in dbcols)
        # 对齐检查
        if dbcols is None:
            mismatches.append(f"{name}: 表 {table} 不存在于 chain 库!")
        else:
            colsset = set(dbcols)
            for f in noncarrier:
                cn = camel_to_underline(f["name"])
                if cn not in colsset:
                    mismatches.append(f"{name}.{f['name']} → 列 {cn} 不在表 {table}")
            ecols = {camel_to_underline(f['name']) for f in noncarrier}
            extra = [c for c in dbcols if c not in ecols]
            if extra:
                mismatches.append(f"{name}: 表 {table} 有实体外列 {extra} (SELECT * 会被 Java 丢弃)")
        varname = 'ent' + name
        lines.append(f'\t"{name}": {varname},')
    lines.append('}')
    lines.append('')
    for name, table, fields, dbcols, msq in ents:
        varname = 'ent' + name
        keys = [f["name"] for f in fields]
        noncarrier = [f for f in fields if not f["carrier"]]
        q = [f for f in noncarrier if f["excel"]]
        excel = [f for f in noncarrier if f["excel"] and f["excel"].group(1) != ""]
        decimals = [f["name"] for f in noncarrier if f["type"] == "BigDecimal"]
        dates = [f["name"] for f in noncarrier if f["type"] == "LocalDate"]
        iso = [f["name"] for f in noncarrier if f["type"] == "LocalDateTime" and not f["jf"]]
        bools = [f["name"] for f in noncarrier if f["type"] in ("Boolean", "boolean")]
        syscol = bool(dbcols and "system_code" in dbcols)
        listsql = msq.get("listByPage", "select * from " + table + " ${searchSql} and CONCAT(id) like #{maps.searchStr} order by create_time desc")
        listnosql = msq.get("listNoPage", listsql)
        lines.append(f'var {varname} = &ent{{')
        lines.append(f'\ttable: {json.dumps(table)},')
        lines.append('\tkeys: []string{' + ', '.join(json.dumps(k, ensure_ascii=False) for k in keys) + '},')
        lines.append('\tdbcols: []string{' + ', '.join(json.dumps(camel_to_underline(f["name"]), ensure_ascii=False) for f in noncarrier) + '},')
        lines.append('\tq: []sqlf{' + ', '.join('{%s, %s}' % (json.dumps(f["name"], ensure_ascii=False), json.dumps('i' if f["kind"] == "i" else ('D' if f["kind"] == "D" else 'w'))) for f in q) + '},')
        if excel:
            lines.append('\texcel: []impField{')
            for f in excel:
                kind = KIND[f["type"]]
                notbank = 'true' if f["notbank"] else 'false'
                sysf = 'true' if f["sys"] else 'false'
                lines.append('\t\t{%s, %s, %s, %s, false, %s},' % (
                    json.dumps(f["excel"].group(1), ensure_ascii=False),
                    json.dumps(camel_to_underline(f["name"])), json.dumps(kind), sysf, notbank))
            lines.append('\t},')
        else:
            lines.append('\texcel: nil,')
        for lst, label in ((decimals, 'decimals'), (dates, 'dates'), (iso, 'isoTimes'), (bools, 'bools')):
            if lst:
                lines.append('\t%s: []string{%s},' % (label, ', '.join(json.dumps(x) for x in lst)))
        lines.append('\tlistSQL: ' + golist(listsql) + ',')
        lines.append('\tlistNoSQL: ' + golist(listnosql) + ',')
        lines.append('\tsearchLit: %s,' % ('true' if '${maps.searchStr}' in listsql else 'false'))
        # 全字段绑定类型(Jackson 校验用):声明序键值对
        lines.append('\ttypes: []jf{' + ', '.join('{%s, %s}' % (json.dumps(f["name"], ensure_ascii=False), json.dumps(f["kind"])) for f in fields) + '},')
        lines.append('\tsyscol: %s,' % ('true' if syscol else 'false'))
        lines.append('}')
        lines.append('')
    with open(OUT, 'w', encoding='utf-8') as f:
        f.write('\n'.join(lines))
    print(f"生成 {OUT}: {len(ents)} 实体")
    print("对齐检查:")
    for m in mismatches:
        print(" ", m)

emit()
