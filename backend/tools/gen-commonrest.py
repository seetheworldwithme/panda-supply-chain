# Generate commonrest entity metadata from the read-only Java source and local chain schema.
import re, os, glob, json, subprocess, collections
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
            anns = []
        elif t and not t.startswith("@") and not t.startswith("*") and not t.startswith("/*") and not t.startswith("//"):
            if fields or "{" in t:
                anns = []
    return name, fields

def golist(s):
    """SQL 字符串 → Go 反引号串(内嵌双引号场景安全)"""
    if '`' not in s and '\r' not in s and '\n' not in s:
        return '`' + s + '`'
    return json.dumps(s)


from pathlib import Path
ROOT=Path(__file__).resolve().parents[1]
SRC=Path(os.environ.get("PANDA_JAVA_SOURCE", "/Volumes/coding/HBuilderxProject/panda/saas-hyqc/saas-supply-chain"))
BASE=SRC/"saas-common/src/main/java/com/saas/common"
# reuse exact source parser, preserving fields and annotation metadata.
parse=parse_entity;snake=camel_to_underline
names=['BusMqConsumeLog','ChainMsg','ChainAreaInfo','ChainCommonAddress','ChainCommonImg','ChainCommonLog','ChainCommonTemplate','ChainCommonTempmx','ChainScreenConfig','ChainRequestLog','ChainCommonMqlog','ChainBasicBussiness','ChainTlzfOrder','ChainTlzfLog','ChainTlzfContractLog','ChainBasicEnterprise','ChainGlobalConfiguration','ChainEqbConfig','ChainEqbFiles','ChainEqbFlow','ChainEqbSingers','ChainCommonDict','ChainTlzfConfig','ChainTlzfSaleEtc','ChainCommonFile']
lines=['// Code generated from read-only Java source by backend/tools/gen-commonrest.py. DO NOT EDIT.','package commonrest','type sqlf struct {Name,Kind string}','var ents=map[string]*ent{']
entries=[]
for n in names:
 files=list(BASE.rglob(n+'.java')) or list((SRC/'saas-core').rglob(n+'.java'));f=files[0]
 _,fields=parse(str(f));non=[x for x in fields if not x['carrier']];mf=list(BASE.rglob(('ChainBasicEep' if n=='ChainBasicEnterprise' else n)+'Mapper.java')) or list((SRC/'saas-core').rglob(n+'Mapper.java')); mp=mf[0].read_text() if mf else ''
 sqls={m[2]:m[1] for m in re.finditer(r'@Select\("((?:[^"\\]|\\.)*)"\)\s*(?:@\w+(?:\([^)]*\))?\s*)*[\w<>.,\s\[\]?]+?\s+(\w+)\(',mp)}
 j=lambda x:json.dumps(x,ensure_ascii=False)
 arr=lambda a:'[]string{'+','.join(map(j,a))+'}'
 table=snake(n);fill=any(x['fill'] for x in non)
 e=['table:'+j(table),'fill:'+str(fill).lower(),'echo:'+str(n not in ['BusMqConsumeLog','ChainCommonMqlog','ChainRequestLog','ChainGlobalConfiguration','ChainEqbFiles','ChainTlzfOrder','ChainTlzfLog','ChainTlzfContractLog','ChainTlzfConfig','ChainTlzfSaleEtc']).lower(),'keys:'+arr([x['name'] for x in fields]),'dbcols:'+arr([snake(x['name']) for x in non]),'q:[]sqlf{'+','.join('{'+j(x['name'])+','+j(x['kind'] if x['kind'] in ['i','D'] else 'w')+'}' for x in non if x['excel'])+'}','types:[]jf{'+','.join('{'+j(x['name'])+','+j(x['kind'])+'}' for x in fields)+'}', 'syscol:'+str('system_code' in DB.get(table,{})).lower()]
 for typ,key in [('BigDecimal','decimals'),('LocalDate','dates'),('Boolean','bools')]:e += [key+':'+arr([x['name'] for x in non if x['type']==typ])]
 e += ['isoTimes:'+arr([x['name'] for x in non if x['type']=='LocalDateTime' and not x['jf']])]
 e += ['excel:[]impField{'+','.join('{'+','.join([j(x['excel'][1]),j(snake(x['name'])),j(KIND.get(x['type'],x['kind'])),str(x['sys']).lower(),'false',str(x['notbank']).lower()])+'}' for x in non if x['excel'] and x['excel'][1])+'}']
 for meth,key in [('listByPage','listSQL'),('listNoPage','listNoSQL')]:e+=[key+':'+j(sqls.get(meth,sqls.get('listByPage','')))]
 entries.append('\nvar ent'+n+'=&ent{'+',\n'.join(e)+'}\n')
 lines += [j(n)+':ent'+n+',']
lines+=['}']+entries
(ROOT/'internal/commonrest/meta_gen.go').write_text('\n'.join(lines))
