import os
#!/usr/bin/env python3
# cmp-capital: 五期批1 资金 Detailed/Settlement(主+Cg/Xs)读对比(Java 8765 vs Go 8799 同库)
# JSON 数字按字面文本保留,对象键排序后比较,数组保序
import json, sys, urllib.request, urllib.error

JAVA, GO = "http://127.0.0.1:8765", "http://127.0.0.1:8799"

def login(base):
    data = __import__("urllib.parse",fromlist=["urlencode"]).urlencode({"username":os.environ["MIGRATION_USER"],"password":os.environ["MIGRATION_PASSWORD"],"code":"0000"}).encode()
    req = urllib.request.Request(base + "/saas-auth/login", data=data,
        headers={"Content-Type": "application/x-www-form-urlencoded"})
    with urllib.request.urlopen(req, timeout=15) as r:
        return json.load(r)["data"]["token"]

JT = login(JAVA); GT = login(GO)
print(f"tokens java={len(JT)} go={len(GT)}")

def post(base, tok, path, body, files=None):
    if files:
        b = "----X" + str(abs(hash(path)) % 10**8)
        payload = (f"--{b}\r\nContent-Disposition: form-data; name=\"file\"; filename=\"t.xlsx\"\r\n"
                   "Content-Type: application/octet-stream\r\n\r\n").encode() + files + f"\r\n--{b}--\r\n".encode()
        req = urllib.request.Request(base + path, data=payload,
            headers={"Authorization": "Bearer " + tok, "Content-Type": "multipart/form-data; boundary=" + b})
    else:
        req = urllib.request.Request(base + path,
            data=json.dumps(body, ensure_ascii=False).encode() if body is not None else b"{}",
            headers={"Authorization": "Bearer " + tok, "Content-Type": "application/json"})
    try:
        with urllib.request.urlopen(req, timeout=30) as r:
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

# ── /chainCapitalDetailed 主 ──
C("dt.list.all", "/saas-capital/chainCapitalDetailed/list", {})
C("dt.list.405", "/saas-capital/chainCapitalDetailed/list", {"company_id": 405})
C("dt.list.detype", "/saas-capital/chainCapitalDetailed/list", {"detype": "1"})
C("dt.list.badkey", "/saas-capital/chainCapitalDetailed/list", {"no_such_col": 1})
C("dt.lbp.plain", "/saas-capital/chainCapitalDetailed/listByPage", {"page": pg})
C("dt.lbp.search", "/saas-capital/chainCapitalDetailed/listByPage", {"searchStr": "银行", "page": pg})
C("dt.lbp.detype", "/saas-capital/chainCapitalDetailed/listByPage", {"detype": "2", "page": pg})
C("dt.lbp.khx", "/saas-capital/chainCapitalDetailed/listByPage", {"khxFlag": True, "page": pg})
C("dt.lbp.hxCode", "/saas-capital/chainCapitalDetailed/listByPage", {"hxCode": "HXJL", "page": pg})
C("dt.lbp.opName", "/saas-capital/chainCapitalDetailed/listByPage", {"oppositeName": "山西", "page": pg})
C("dt.lbp.transTime", "/saas-capital/chainCapitalDetailed/listByPage", {"transTime": "2023-06-15", "page": pg})
C("dt.lbp.transTimeBad", "/saas-capital/chainCapitalDetailed/listByPage", {"transTime": "2023/06/15", "page": pg})
C("dt.lbp.range", "/saas-capital/chainCapitalDetailed/listByPage", {"startTime": "2023-01-01", "endTime": "2023-12-31", "page": pg})
C("dt.lbp.blanks", "/saas-capital/chainCapitalDetailed/listByPage", {"detype": "", "hxCode": "", "oppositeName": "", "transTime": "", "oppositeAccount": "", "bankName": "", "startTime": "", "endTime": "", "searchStr": "", "page": pg})
C("dt.lbp.nulllike", "/saas-capital/chainCapitalDetailed/listByPage", {"page": pg, "searchStr": None})
C("dt.lbp.inject", "/saas-capital/chainCapitalDetailed/listByPage", {"searchStr": "'", "page": pg})
C("dt.lbp.bigpage", "/saas-capital/chainCapitalDetailed/listByPage", {"page": {"current": 2, "size": 5}})
C("dt.getById.1", "/saas-capital/chainCapitalDetailed/getById", {"id": 1})
C("dt.getById.missing", "/saas-capital/chainCapitalDetailed/getById", {"id": 999999999})
C("dt.getById.noid", "/saas-capital/chainCapitalDetailed/getById", {})
C("dt.delBatch.empty", "/saas-capital/chainCapitalDetailed/delBatch", {"ids": []})
C("dt.delBatch.null", "/saas-capital/chainCapitalDetailed/delBatch", {"ids": None})
C("dt.delBatch.noids", "/saas-capital/chainCapitalDetailed/delBatch", {})
C("dt.getSments.1", "/saas-capital/chainCapitalDetailed/getSments", {"id": 1})
C("dt.getSments.missing", "/saas-capital/chainCapitalDetailed/getSments", {"id": 999999999})
C("dt.sumPrice.plain", "/saas-capital/chainCapitalDetailed/getSumPrice", {})
C("dt.sumPrice.detype", "/saas-capital/chainCapitalDetailed/getSumPrice", {"detype": "1"})
C("dt.sumPrice.year", "/saas-capital/chainCapitalDetailed/getSumPrice", {"timeType": 4})
C("dt.sumPrice.contract", "/saas-capital/chainCapitalDetailed/getSumPrice", {"contractId": "529"})
C("dt.ptSum.plain", "/saas-capital/chainCapitalDetailed/getPtSumPrice", {})
C("dt.ptSum.405", "/saas-capital/chainCapitalDetailed/getPtSumPrice", {"companyId": "405"})
C("dt.ptSum.xzqh", "/saas-capital/chainCapitalDetailed/getPtSumPrice", {"xzqhCode": "130000"})
C("dt.ptSum.zero", "/saas-capital/chainCapitalDetailed/getPtSumPrice", {"xzqhCode": "0"})
C("dt.banks", "/saas-capital/chainCapitalDetailed/getBanks", {})
C("dt.getDetails.plain", "/saas-capital/chainCapitalDetailed/getDetails", {"page": pg})
C("dt.getDetails.contract", "/saas-capital/chainCapitalDetailed/getDetails", {"contractId": "529", "page": pg})
C("dt.getDetails.search", "/saas-capital/chainCapitalDetailed/getDetails", {"searchStr": "银行", "page": pg})
C("dt.sysCode.1", "/saas-capital/chainCapitalDetailed/getSystemCode", {"codeType": "1"})
C("dt.sysCode.2", "/saas-capital/chainCapitalDetailed/getSystemCode", {"codeType": "2"})
C("dt.sysCode.none", "/saas-capital/chainCapitalDetailed/getSystemCode", {})
C("dt.cancelBal.empty", "/saas-capital/chainCapitalDetailed/cancelBalancing", [])
C("dt.cancelBal.obj", "/saas-capital/chainCapitalDetailed/cancelBalancing", {})

# ── Cg/Xs Detailed 变体(抽测) ──
for base in ["chainCapitalCgDetailed", "chainCapitalXsDetailed"]:
    C(f"{base}.list.all", f"/saas-capital/{base}/list", {})
    C(f"{base}.lbp.detype", f"/saas-capital/{base}/listByPage", {"detype": "1", "page": pg})
    C(f"{base}.countDetail.4", f"/saas-capital/{base}/countDetail", {"timeType": 4})
    C(f"{base}.countDetail.range", f"/saas-capital/{base}/countDetail", {"startTime": "2026-01-01", "endTime": "2026-03-01"})
    C(f"{base}.countDetail.nulldates", f"/saas-capital/{base}/countDetail", {})
    C(f"{base}.countDetail.month", f"/saas-capital/{base}/countDetail", {"startTime": "2026-01-01", "endTime": "2026-04-01"})
    C(f"{base}.getSumPrice", f"/saas-capital/{base}/getSumPrice", {})
    C(f"{base}.sysCode", f"/saas-capital/{base}/getSystemCode", {"codeType": "1"})

# ── /chainCapitalSettlement 主 ──
C("st.list.all", "/saas-capital/chainCapitalSettlement/list", {})
C("st.list.type", "/saas-capital/chainCapitalSettlement/list", {"type": 1})
C("st.list.badkey", "/saas-capital/chainCapitalSettlement/list", {"no_such_col": 1})
C("st.lbp.plain", "/saas-capital/chainCapitalSettlement/listByPage", {"page": pg})
C("st.lbp.search", "/saas-capital/chainCapitalSettlement/listByPage", {"searchStr": "结算", "page": pg})
C("st.lbp.type", "/saas-capital/chainCapitalSettlement/listByPage", {"type": 2, "page": pg})
C("st.lbp.khx", "/saas-capital/chainCapitalSettlement/listByPage", {"khxFlag": True, "page": pg})
C("st.lbp.supplier", "/saas-capital/chainCapitalSettlement/listByPage", {"supplierName": "山西", "page": pg})
C("st.lbp.month", "/saas-capital/chainCapitalSettlement/listByPage", {"month": "2026-01", "page": pg})
C("st.lbp.range", "/saas-capital/chainCapitalSettlement/listByPage", {"startTime": "2026-01-15", "endTime": "2026-03-20", "page": pg})
C("st.lbp.relFlag", "/saas-capital/chainCapitalSettlement/listByPage", {"relBillFlag": True, "page": pg})
C("st.lbp.blanks", "/saas-capital/chainCapitalSettlement/listByPage", {"type": "", "hxCode": "", "supplierId": "", "supplierName": "", "contractId": "", "setStartTime": "", "setEndTime": "", "relBillFlag": "", "relationId": "", "month": "", "startTime": "", "endTime": "", "searchStr": "", "page": pg})
C("st.lbp.badDate", "/saas-capital/chainCapitalSettlement/listByPage", {"setStartTime": "2026/01/01", "page": pg})
C("st.lbp.badFlag", "/saas-capital/chainCapitalSettlement/listByPage", {"relBillFlag": "abc", "page": pg})
C("st.getById.1", "/saas-capital/chainCapitalSettlement/getById", {"id": 1})
C("st.getById.missing", "/saas-capital/chainCapitalSettlement/getById", {"id": 999999999})
C("st.delBatch.empty", "/saas-capital/chainCapitalSettlement/delBatch", {"ids": []})
C("st.getDetails.1", "/saas-capital/chainCapitalSettlement/getDetails", {"id": 1})
C("st.getDetails.missing", "/saas-capital/chainCapitalSettlement/getDetails", {"id": 999999999})
C("st.getBills.1", "/saas-capital/chainCapitalSettlement/getBills", {"id": 1})
C("st.getBills.missing", "/saas-capital/chainCapitalSettlement/getBills", {"id": 999999999})
C("st.sumPrice.plain", "/saas-capital/chainCapitalSettlement/getSumPrice", {})
C("st.sumPrice.type", "/saas-capital/chainCapitalSettlement/getSumPrice", {"type": 1})
C("st.sumPrice.relFlag", "/saas-capital/chainCapitalSettlement/getSumPrice", {"relBillFlag": True})
C("st.ptSum.plain", "/saas-capital/chainCapitalSettlement/getPtSumPrice", {})
C("st.ptSum.405", "/saas-capital/chainCapitalSettlement/getPtSumPrice", {"companyId": "405"})
C("st.ptSum.xzqh", "/saas-capital/chainCapitalSettlement/getPtSumPrice", {"xzqhCode": "130000"})
C("st.sysCode.1", "/saas-capital/chainCapitalSettlement/getSystemCode", {"codeType": "1"})
C("st.sysCode.2", "/saas-capital/chainCapitalSettlement/getSystemCode", {"codeType": "2"})
C("st.selectQuat.1", "/saas-capital/chainCapitalSettlement/selectQuantity", {"timeType": 1})
C("st.selectQuat.2", "/saas-capital/chainCapitalSettlement/selectQuantity", {"timeType": 2})
C("st.selectQuat.3", "/saas-capital/chainCapitalSettlement/selectQuantity", {"timeType": 3})
C("st.selectQuat.4", "/saas-capital/chainCapitalSettlement/selectQuantity", {"timeType": 4})
C("st.selectQuat.null", "/saas-capital/chainCapitalSettlement/selectQuantity", {})
C("st.countSett.4", "/saas-capital/chainCapitalSettlement/countSett", {"timeType": 4})
C("st.countSett.range", "/saas-capital/chainCapitalSettlement/countSett", {"startTime": "2026-01-01", "endTime": "2026-03-01"})
C("st.countSett.nulldates", "/saas-capital/chainCapitalSettlement/countSett", {})

# ── Cg/Xs Settlement 变体 ──
for base in ["chainCapitalCgSettlement", "chainCapitalXsSettlement"]:
    C(f"{base}.list.all", f"/saas-capital/{base}/list", {})
    C(f"{base}.lbp.type", f"/saas-capital/{base}/listByPage", {"type": 1, "page": pg})
    C(f"{base}.selectQuat.4", f"/saas-capital/{base}/selectQuantity", {"timeType": 4})
    C(f"{base}.selectQuat.null", f"/saas-capital/{base}/selectQuantity", {})
    C(f"{base}.countSett.4", f"/saas-capital/{base}/countSett", {"timeType": 4})
    C(f"{base}.countSett.range", f"/saas-capital/{base}/countSett", {"startTime": "2026-01-01", "endTime": "2026-04-01"})
    C(f"{base}.countSett.nulldates", f"/saas-capital/{base}/countSett", {})
    C(f"{base}.getSumPrice", f"/saas-capital/{base}/getSumPrice", {})
    C(f"{base}.sysCode", f"/saas-capital/{base}/getSystemCode", {"codeType": "2"})

# ── 主路径 dashboard 注册版对照(countSett/selectQuantity/countDetail 走一期实现) ──
C("dash.countSett.4", "/saas-capital/chainCapitalSettlement/countSett", {"timeType": 4})
C("dash.selectQuat.1", "/saas-capital/chainCapitalSettlement/selectQuantity", {"timeType": 1})
C("dash.countDetail.4", "/saas-capital/chainCapitalDetailed/countDetail", {"timeType": 4})
C("dash.countDetail.range", "/saas-capital/chainCapitalDetailed/countDetail", {"startTime": "2026-01-01", "endTime": "2026-03-01"})

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
            print(f"  J: {nj[:400]}")
            print(f"  G: {ng[:400]}")
    else:
        print(f"ok   {name}")
print(f"\n{len(CASES)-fails}/{len(CASES)} passed")
sys.exit(1 if fails else 0)
