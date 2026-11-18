import os
#!/usr/bin/env python3
# cmp-logistics: 四期物流 62 端点读对比(Java 8765 vs Go 8799 同库)
# JSON 数字按字面文本保留("0.00"≠"0"),对象键排序后比较,数组保序
import json, subprocess, sys, urllib.request, urllib.error

JAVA, GO = "http://127.0.0.1:8765", "http://127.0.0.1:8799"

def login(base):
    data = __import__("urllib.parse",fromlist=["urlencode"]).urlencode({"username":os.environ["MIGRATION_USER"],"password":os.environ["MIGRATION_PASSWORD"],"code":"0000"}).encode()
    req = urllib.request.Request(base + "/saas-auth/login", data=data,
        headers={"Content-Type": "application/x-www-form-urlencoded"})
    with urllib.request.urlopen(req, timeout=15) as r:
        return json.load(r)["data"]["token"]

JT = login(JAVA); GT = login(GO)
print(f"tokens java={len(JT)} go={len(GT)}")

def post(base, tok, path, body, extra=None):
    req = urllib.request.Request(base + "/saas-logistics/" + path,
        data=json.dumps(body, ensure_ascii=False).encode() if body is not None else b"")
    h = {"Authorization": "Bearer " + tok, "Content-Type": "application/json"}
    if extra: h.update(extra)
    for k, v in h.items(): req.add_header(k, v)
    try:
        with urllib.request.urlopen(req, timeout=30) as r:
            return r.status, r.read().decode("utf-8", "replace")
    except urllib.error.HTTPError as e:
        return e.code, e.read().decode("utf-8", "replace")
    except Exception as e:
        return -1, str(e)

def norm(s, sort_recs=False):
    try:
        obj = json.loads(s, parse_float=str, parse_int=str)
        if sort_recs and isinstance(obj, dict) and isinstance(obj.get("data"), dict) \
           and isinstance(obj["data"].get("records"), list):
            obj["data"]["records"] = sorted(obj["data"]["records"], key=lambda r: str(r.get("id")))
            obj["data"].pop("pages", None)
        return json.dumps(obj, sort_keys=True, ensure_ascii=False)
    except Exception:
        return "RAW:" + s.strip()

CASES = []
SORT_REC = set()
def C(name, path, body=None, extra=None, sort_recs=False):
    CASES.append((name, path, body, extra))
    if sort_recs:
        SORT_REC.add(name)

P = "page"
pg = {"current": 1, "size": 3}

# ── /chainLogistics ──
C("lo.list.all", "chainLogistics/list", {})
C("lo.list.54", "chainLogistics/list", {"company_id": 54})
C("lo.list.badkey", "chainLogistics/list", {"no_such_col": 1})
C("lo.listByPage.search", "chainLogistics/listByPage", {"searchStr": "煤", "page": pg})
C("lo.listByPage.plain", "chainLogistics/listByPage", {"page": pg})
C("lo.listByPage.loading", "chainLogistics/listByPage", {"loadingAddress": "山西", "page": pg})
C("lo.listByPage.unloading", "chainLogistics/listByPage", {"unloadingPoint": "河北", "page": pg})
C("lo.listByPage.billNo", "chainLogistics/listByPage", {"billNumber": "YD2023", "page": pg})
C("lo.listByPage.transTime", "chainLogistics/listByPage", {"transTime": "2023-01-18T00:41:37", "page": pg})
C("lo.listByPage.month", "chainLogistics/listByPage", {"startTime": "2023-01-05", "endTime": "2023-03-01", "page": pg})
C("lo.listByPage.platform", "chainLogistics/listByPage", {"carriagePlatform": 4, "page": pg})
C("lo.listByPage.notBind", "chainLogistics/listByPage", {"notBindOrder": True, "orderId": 211, "page": pg})
C("lo.listByPage.orderId", "chainLogistics/listByPage", {"orderId": 211, "page": pg})
C("lo.listByPage.inject", "chainLogistics/listByPage", {"searchStr": "'", "page": pg})
C("lo.listByPage.blankstr", "chainLogistics/listByPage", {"searchStr": "YD20230118004137", "billNumber": "", "startTime": "", "endTime": "", "loadingAddress": "", "unloadingPoint": "", "carriagePlatform": "", "orderId": "", "notBindOrder": "", "transTime": "", "page": pg})  # 全字段空串绑定(Java ""→null 不拼条件),小结果集避开并列序
C("lo.listByPage.badtime", "chainLogistics/listByPage", {"transTime": "2023-01-18 00:41:37", "page": pg})
C("lo.getById.4493", "chainLogistics/getById", {"id": 4493})
C("lo.getById.missing", "chainLogistics/getById", {"id": 999999999})
C("lo.getById.noid", "chainLogistics/getById", {})
C("lo.deleteById.missing", "chainLogistics/deleteById", {"id": 999999999})
C("lo.delBatch.empty", "chainLogistics/delBatch", {"ids": []})
C("lo.delBatch.null", "chainLogistics/delBatch", {"ids": None})
C("lo.countBuses.all", "chainLogistics/countBuses", {})
C("lo.countBuses.405", "chainLogistics/countBuses", {"companyId": "405"})
C("lo.countBuses.xzqh", "chainLogistics/countBuses", {"xzqhCode": "130000"})
C("lo.countBuses.zero", "chainLogistics/countBuses", {"xzqhCode": "0"})
C("lo.countMonBuses.empty", "chainLogistics/countMonBuses", {})
C("lo.countMonBuses.h1", "chainLogistics/countMonBuses", {"startTime": "2023-01-01", "endTime": "2023-06-30"})
C("lo.countMonBuses.co", "chainLogistics/countMonBuses", {"startTime": "2022-11-01", "endTime": "2023-02-31", "companyId": "405"})
C("lo.countMonBuses.baddate", "chainLogistics/countMonBuses", {"startTime": "abc", "endTime": "2023-06-30"})
C("lo.updTransLocation.missingbody", "chainLogistics/updTransLocation", {})

# ── /chainLogisticsBill ──
C("bill.list.all", "chainLogisticsBill/list", {})
C("bill.list.405", "chainLogisticsBill/list", {"company_id": 405})
C("bill.list.badkey", "chainLogisticsBill/list", {"no_such_col": 1})
C("bill.listByPage.search", "chainLogisticsBill/listByPage", {"searchStr": "海", "page": pg})
C("bill.listByPage.plain", "chainLogisticsBill/listByPage", {"page": pg})
C("bill.listByPage.null", "chainLogisticsBill/listByPage", {"searchStr": None, "page": pg})
C("bill.listByPage.missing", "chainLogisticsBill/listByPage", {"page": pg})
C("bill.listByPage.inject", "chainLogisticsBill/listByPage", {"searchStr": "'", "page": pg})
C("bill.getById.10", "chainLogisticsBill/getById", {"id": 10})
C("bill.getById.missing", "chainLogisticsBill/getById", {"id": 999999999})
C("bill.getById.noid", "chainLogisticsBill/getById", {})
C("bill.getSystemCode", "chainLogisticsBill/getSystemCode", {})
C("bill.selectQuat.empty", "chainLogisticsBill/selectQuat", {})
C("bill.selectQuat.t4", "chainLogisticsBill/selectQuat", {"timeType": 4})
C("bill.selectQuat.t1", "chainLogisticsBill/selectQuat", {"timeType": 1})
C("bill.selectQuat.str", "chainLogisticsBill/selectQuat", {"timeType": "4"})
C("bill.getContracts.211", "chainLogisticsBill/getContracts", {"id": 211})
C("bill.getContracts.noid", "chainLogisticsBill/getContracts", {})
C("bill.getSumPrice", "chainLogisticsBill/getSumPrice", {})
C("bill.countBuses.all", "chainLogisticsBill/countBuses", {})
C("bill.countBuses.405", "chainLogisticsBill/countBuses", {"companyId": "405"})
C("bill.countMonBuses.empty", "chainLogisticsBill/countMonBuses", {})
C("bill.countMonBuses.range", "chainLogisticsBill/countMonBuses", {"startTime": "2023-01-01", "endTime": "2023-06-30"})
C("bill.countLogist.t4", "chainLogisticsBill/countLogist", {"timeType": 4})
C("bill.countLogist.m3", "chainLogisticsBill/countLogist", {"startTime": "2023-01-01", "endTime": "2023-03-01"})
C("bill.countLogist.d1", "chainLogisticsBill/countLogist", {"startTime": "2023-01-01", "endTime": "2023-01-31"})
C("bill.countLogist.rev", "chainLogisticsBill/countLogist", {"startTime": "2023-03-01", "endTime": "2023-01-01"})
C("bill.countLogist.missing", "chainLogisticsBill/countLogist", {})
C("bill.countLogist.same", "chainLogisticsBill/countLogist", {"startTime": "2023-01-05", "endTime": "2023-01-09"})
C("bill.selOrders.13662", "chainLogisticsBill/selOrders", {"contractId": 13662, "page": pg})
C("bill.selOrders.none", "chainLogisticsBill/selOrders", {"contractId": 999999999, "page": pg})
C("bill.selOrders.noid", "chainLogisticsBill/selOrders", {"page": pg})
C("bill.selAllOrders.13662", "chainLogisticsBill/selAllOrders", {"contractId": 13662})
C("bill.selAllOrders.none", "chainLogisticsBill/selAllOrders", {"contractId": 999999999})
C("bill.selAllOrders.noid", "chainLogisticsBill/selAllOrders", {})

# ── /chainLogisticsContract ──
C("lc.list.all", "chainLogisticsContract/list", {})
C("lc.listByPage.plain", "chainLogisticsContract/listByPage", {"page": pg})
C("lc.listByPage.search", "chainLogisticsContract/listByPage", {"searchStr": "1", "page": pg})
C("lc.listByPage.missing", "chainLogisticsContract/listByPage", {"page": pg})
C("lc.getById.real", "chainLogisticsContract/getById", {"id": 1})
C("lc.getById.missing", "chainLogisticsContract/getById", {"id": 999999999})
C("lc.deleteById.missing", "chainLogisticsContract/deleteById", {"id": 999999999})
C("lc.delBatch.empty", "chainLogisticsContract/delBatch", {"ids": []})
C("lc.delBatch.null", "chainLogisticsContract/delBatch", {"ids": None})

# ── /chainLogisticsThirdWaybill ──
C("tw.list.all", "chainLogisticsThirdWaybill/list", {})
C("tw.listByPage.plain", "chainLogisticsThirdWaybill/listByPage", {"page": pg})
C("tw.listByPage.search", "chainLogisticsThirdWaybill/listByPage", {"searchStr": "YD", "page": pg})
C("tw.listByPage.clean1", "chainLogisticsThirdWaybill/listByPage", {"cleanFlag": True, "page": pg})
C("tw.listByPage.clean0", "chainLogisticsThirdWaybill/listByPage", {"cleanFlag": False, "page": pg})
C("tw.listByPage.billNo", "chainLogisticsThirdWaybill/listByPage", {"billNumber": "YD2020", "page": pg})
C("tw.getById.3491", "chainLogisticsThirdWaybill/getById", {"id": 3491})
C("tw.getById.missing", "chainLogisticsThirdWaybill/getById", {"id": 999999999})
C("tw.addBatch.dup", "chainLogisticsThirdWaybill/addBatch", [{"billNumber": "A"}, {"billNumber": "A"}],
  {"BatchNumber": "CMPBATCH1", "BatchCount": "2"})
C("tw.addBatch.countmismatch", "chainLogisticsThirdWaybill/addBatch", [{"billNumber": "CMPX1"}],
  {"BatchNumber": "CMPBATCH2", "BatchCount": "5"})
C("tw.addBatch.nocount", "chainLogisticsThirdWaybill/addBatch", [{"billNumber": "CMPX2"}],
  {"BatchNumber": "CMPBATCH3"})
C("tw.updateBatch.dup", "chainLogisticsThirdWaybill/updateBatch", [{"billNumber": "A"}, {"billNumber": "A"}],
  {"BatchCount": "2"})
C("tw.updateBatch.nocount", "chainLogisticsThirdWaybill/updateBatch", [{"billNumber": "A"}])
C("tw.deleteById.missing", "chainLogisticsThirdWaybill/deleteById", {"id": 999999999})
C("tw.delBatch.empty", "chainLogisticsThirdWaybill/delBatch", {"ids": []})

# ── /chainLogisticsWeighingList ──
C("wl.list.all", "chainLogisticsWeighingList/list", {})
C("wl.listByPage.plain", "chainLogisticsWeighingList/listByPage", {"page": pg})
C("wl.listByPage.search", "chainLogisticsWeighingList/listByPage", {"searchStr": "10", "page": pg})
C("wl.listByPage.order", "chainLogisticsWeighingList/listByPage", {"orderId": 10, "page": pg})
C("wl.listByPage.notBind", "chainLogisticsWeighingList/listByPage", {"notBindOrder": True, "orderId": 10, "page": pg})
C("wl.getById.10", "chainLogisticsWeighingList/getById", {"id": 10})
C("wl.getById.missing", "chainLogisticsWeighingList/getById", {"id": 999999999})
C("wl.deleteById.missing", "chainLogisticsWeighingList/deleteById", {"id": 999999999})
C("wl.delBatch.empty", "chainLogisticsWeighingList/delBatch", {"ids": []})

ok = fail = 0
fails = []
for name, path, body, extra in CASES:
    js, jraw = post(JAVA, JT, path, body, extra)
    gs, graw = post(GO, GT, path, body, extra)
    sr = name in SORT_REC
    jn, gn = norm(jraw, sr), norm(graw, sr)
    if js == gs and jn == gn:
        ok += 1
        print(f"PASS {name}")
    else:
        fail += 1
        print(f"FAIL {name} java={js} go={gs}")
        if jn != gn:
            print(f"  J: {jn[:600]}")
            print(f"  G: {gn[:600]}")
        fails.append(name)
print(f"\n==== cmp-logistics: {ok} PASS / {fail} FAIL / {len(CASES)} total")
if fails:
    print("fails:", ",".join(fails))
sys.exit(1 if fail else 0)
