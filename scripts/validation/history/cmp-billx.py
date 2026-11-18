import os
#!/usr/bin/env python3
# cmp-billx: 五期批2 票据 5 控制器读对比(Java 8765 vs Go 8799 同库)
import json, sys, urllib.request, urllib.error

JAVA, GO = "http://127.0.0.1:8765", "http://127.0.0.1:8799"

def login(base):
    data = __import__("urllib.parse",fromlist=["urlencode"]).urlencode({"username":os.environ["MIGRATION_USER"],"password":os.environ["MIGRATION_PASSWORD"],"code":"0000"}).encode()
    req = urllib.request.Request(base + "/saas-auth/login", data=data,
        headers={"Content-Type": "application/x-www-form-urlencoded"})
    with urllib.request.urlopen(req, timeout=15) as r:
        return json.load(r)["data"]["token"]

JT = login(JAVA); GT = login(GO)

def post(base, tok, path, body):
    req = urllib.request.Request(base + path,
        data=json.dumps(body, ensure_ascii=False).encode() if body is not None else b"{}",
        headers={"Authorization": "Bearer " + tok, "Content-Type": "application/json"})
    try:
        with urllib.request.urlopen(req, timeout=60) as r:
            return r.status, r.read().decode("utf-8", "replace")
    except urllib.error.HTTPError as e:
        return e.code, e.read().decode("utf-8", "replace")
    except Exception as e:
        return -1, str(e)

def norm(s):
    try:
        obj = json.loads(s, parse_float=str, parse_int=str)
        return json.dumps(obj, sort_keys=True, ensure_ascii=False)
    except Exception:
        return "RAW:" + s.strip()

CASES = []
def C(name, path, body=None):
    CASES.append((name, path, body))

pg = {"current": 1, "size": 3}

# ── /chainBill 主 ──
C("bl.list.all", "/saas-bill/chainBill/list", {})
C("bl.list.type", "/saas-bill/chainBill/list", {"type": 1})
C("bl.list.badkey", "/saas-bill/chainBill/list", {"no_such_col": 1})
C("bl.lbp.plain", "/saas-bill/chainBill/listByPage", {"page": pg})
C("bl.lbp.search", "/saas-bill/chainBill/listByPage", {"searchStr": "煤", "page": pg})
C("bl.lbp.type", "/saas-bill/chainBill/listByPage", {"type": 1, "page": pg})
C("bl.lbp.companyName", "/saas-bill/chainBill/listByPage", {"companyName": "山西", "page": pg})
C("bl.lbp.month", "/saas-bill/chainBill/listByPage", {"month": "2024-01", "page": pg})
C("bl.lbp.openDate", "/saas-bill/chainBill/listByPage", {"openDate": "2024-01-10", "page": pg})
C("bl.lbp.openDateBad", "/saas-bill/chainBill/listByPage", {"openDate": "2024/01/10", "page": pg})
C("bl.lbp.range", "/saas-bill/chainBill/listByPage", {"startTime": "2024-01-15", "endTime": "2024-03-20", "page": pg})
C("bl.lbp.relFlag", "/saas-bill/chainBill/listByPage", {"relationFlag": True, "page": pg})
C("bl.lbp.relFlagFalse", "/saas-bill/chainBill/listByPage", {"relationFlag": False, "page": pg})
C("bl.lbp.zeroVal", "/saas-bill/chainBill/listByPage", {"type": "0", "page": pg})  # getSql "0" 值跳过
C("bl.lbp.inType", "/saas-bill/chainBill/listByPage", {"inType": "电", "page": pg})
C("bl.lbp.billCode", "/saas-bill/chainBill/listByPage", {"billCode": "140", "page": pg})
C("bl.lbp.blanks", "/saas-bill/chainBill/listByPage", {"companyName": "", "inType": "", "openDate": "", "month": "", "startTime": "", "endTime": "", "relationFlag": "", "searchStr": "", "type": "", "page": pg})
C("bl.lbp.bigpage", "/saas-bill/chainBill/listByPage", {"page": {"current": 2, "size": 5}})
C("bl.getById.1", "/saas-bill/chainBill/getById", {"id": 1})
C("bl.getById.miss", "/saas-bill/chainBill/getById", {"id": 999999999})
C("bl.delBatch.empty", "/saas-bill/chainBill/delBatch", {"ids": []})
C("bl.delBatch.null", "/saas-bill/chainBill/delBatch", {"ids": None})
C("bl.delBatch.noids", "/saas-bill/chainBill/delBatch", {})
C("bl.getSments.1", "/saas-bill/chainBill/getSments", {"id": 1})
C("bl.getSments.miss", "/saas-bill/chainBill/getSments", {"id": 999999999})
C("bl.sumPrice.plain", "/saas-bill/chainBill/getSumPrice", {})
C("bl.sumPrice.type", "/saas-bill/chainBill/getSumPrice", {"type": 2})
C("bl.sumPrice.relFlag", "/saas-bill/chainBill/getSumPrice", {"relationFlag": True})
C("bl.sumPrice.contract", "/saas-bill/chainBill/getSumPrice", {"contractId": "529"})
C("bl.sysCode.1", "/saas-bill/chainBill/getSystemCode", {"codeType": "1"})
C("bl.sysCode.2", "/saas-bill/chainBill/getSystemCode", {"codeType": "2"})
C("bl.selectBusCount.plain", "/saas-bill/chainBill/selectBusCount", {})
C("bl.selectBusCount.range", "/saas-bill/chainBill/selectBusCount", {"startTime": "2024-01-01", "endTime": "2026-12-31"})
C("bl.countBill.1", "/saas-bill/chainBill/countBill", {"timeType": 1})
C("bl.countBill.3", "/saas-bill/chainBill/countBill", {"timeType": 3})
C("bl.countBill.4", "/saas-bill/chainBill/countBill", {"timeType": 4})
C("bl.countBill.null", "/saas-bill/chainBill/countBill", {})
C("bl.countTaxRate.4", "/saas-bill/chainBill/countTaxRate", {"startTime": "2026-01-01", "endTime": "2026-06-01"})
C("bl.countTaxRate.month", "/saas-bill/chainBill/countTaxRate", {"startTime": "2024-01-01", "endTime": "2024-06-01"})
C("bl.countTaxRate.nulls", "/saas-bill/chainBill/countTaxRate", {})
C("bl.selZeroTax", "/saas-bill/chainBill/selZeroTax", {})
C("bl.selJxzb", "/saas-bill/chainBill/selJxzb", {})
C("bl.searchRelBill.plain", "/saas-bill/chainBill/searchRelBill", {"page": pg})
C("bl.searchRelBill.type", "/saas-bill/chainBill/searchRelBill", {"type": 1, "page": pg})
C("bl.searchRelBill.search", "/saas-bill/chainBill/searchRelBill", {"searchStr": "煤", "page": pg})
C("bl.countBuses.plain", "/saas-bill/chainBill/countBuses", {})
C("bl.countBuses.405", "/saas-bill/chainBill/countBuses", {"companyId": "405"})
C("bl.countBuses.xzqh", "/saas-bill/chainBill/countBuses", {"xzqhCode": "130000"})
C("bl.countBuses.t1", "/saas-bill/chainBill/countBuses", {"timeType": "1"})
C("bl.countMonBuses.range", "/saas-bill/chainBill/countMonBuses", {"startTime": "2024-01-01", "endTime": "2024-06-01"})
C("bl.countMonBuses.nulls", "/saas-bill/chainBill/countMonBuses", {})
C("bl.countPtTax.plain", "/saas-bill/chainBill/countPtTax", {})
C("bl.countPtTax.405", "/saas-bill/chainBill/countPtTax", {"companyId": "405"})
C("bl.countYearBuses", "/saas-bill/chainBill/countYearBuses", {})

# ── Cg/Xs 变体 ──
for base in ["chainCgBill", "chainXsBill"]:
    C(f"{base}.list.all", f"/saas-bill/{base}/list", {})
    C(f"{base}.lbp.type", f"/saas-bill/{base}/listByPage", {"type": 2, "page": pg})
    C(f"{base}.getSumPrice", f"/saas-bill/{base}/getSumPrice", {"type": 2})
    C(f"{base}.selectBusCount", f"/saas-bill/{base}/selectBusCount", {})
    C(f"{base}.countTaxRate.4", f"/saas-bill/{base}/countTaxRate", {"startTime": "2026-01-01", "endTime": "2026-03-01"})
    C(f"{base}.countTaxRate.nulls", f"/saas-bill/{base}/countTaxRate", {})
    C(f"{base}.selZeroTax", f"/saas-bill/{base}/selZeroTax", {})
    C(f"{base}.selJxzb", f"/saas-bill/{base}/selJxzb", {})
    C(f"{base}.countBuses", f"/saas-bill/{base}/countBuses", {})
    C(f"{base}.countMonBuses.range", f"/saas-bill/{base}/countMonBuses", {"startTime": "2024-01-01", "endTime": "2024-03-01"})
    C(f"{base}.countMonBuses.nulls", f"/saas-bill/{base}/countMonBuses", {})
    C(f"{base}.countPtTax", f"/saas-bill/{base}/countPtTax", {})
    C(f"{base}.sysCode", f"/saas-bill/{base}/getSystemCode", {"codeType": "2"})
    C(f"{base}.countBill.null", f"/saas-bill/{base}/countBill", {})

# ── /chainBillRate ──
C("rt.list.all", "/saas-bill/chainBillRate/list", {})
C("rt.lbp.plain", "/saas-bill/chainBillRate/listByPage", {"page": pg})
C("rt.lbp.search", "/saas-bill/chainBillRate/listByPage", {"searchStr": "2", "page": pg})
C("rt.lbp.month", "/saas-bill/chainBillRate/listByPage", {"month": "2024-01", "page": pg})
C("rt.lbp.zeroSkip", "/saas-bill/chainBillRate/listByPage", {"month": "0", "page": pg})
C("rt.getById.1", "/saas-bill/chainBillRate/getById", {"id": 1})
C("rt.delBatch.empty", "/saas-bill/chainBillRate/delBatch", {"ids": []})
C("rt.delBatch.null", "/saas-bill/chainBillRate/delBatch", {"ids": None})
C("rt.getSumPrice.plain", "/saas-bill/chainBillRate/getSumPrice", {})
C("rt.getSumPrice.405", "/saas-bill/chainBillRate/getSumPrice", {"companyId": "405"})
C("rt.getSumPrice.range", "/saas-bill/chainBillRate/getSumPrice", {"startTime": "2024-01-01", "endTime": "2026-12-31"})
C("rt.countMonBuses.range", "/saas-bill/chainBillRate/countMonBuses", {"startTime": "2024-01-01", "endTime": "2024-06-01"})
C("rt.countMonBuses.nulls", "/saas-bill/chainBillRate/countMonBuses", {})
C("rt.countMonBuses.405", "/saas-bill/chainBillRate/countMonBuses", {"companyId": "405", "startTime": "2024-01-01", "endTime": "2024-03-01"})

# ── /chainBillSettlement ──
C("bs.list.all", "/saas-bill/chainBillSettlement/list", {})
C("bs.lbp.plain", "/saas-bill/chainBillSettlement/listByPage", {"page": pg})
C("bs.lbp.search", "/saas-bill/chainBillSettlement/listByPage", {"searchStr": "1", "page": pg})
C("bs.lbp.nullstr", "/saas-bill/chainBillSettlement/listByPage", {"page": pg})
C("bs.getById.1", "/saas-bill/chainBillSettlement/getById", {"id": 1})
C("bs.delBatch.empty", "/saas-bill/chainBillSettlement/delBatch", {"ids": []})

fails = 0
for name, path, body in CASES:
    sj, rj = post(JAVA, JT, path, body)
    sg, rg = post(GO, GT, path, body)
    ok = sj == sg and norm(rj) == norm(rg)
    if not ok:
        fails += 1
        print(f"FAIL {name} [{sj}/{sg}]")
        nj, ng = norm(rj), norm(rg)
        if nj != ng:
            print(f"  J: {nj[:320]}")
            print(f"  G: {ng[:320]}")
    else:
        print(f"ok   {name}")
print(f"\n{len(CASES)-fails}/{len(CASES)} passed")
sys.exit(1 if fails else 0)
