import os
#!/usr/bin/env python3
# cmp-capital-small: 五期批1.5 小控制器 7 个 + payment 读对比(Java 8765 vs Go 8799 同库)
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
        with urllib.request.urlopen(req, timeout=30) as r:
            return r.status, r.read().decode("utf-8", "replace")
    except urllib.error.HTTPError as e:
        return e.code, e.read().decode("utf-8", "replace")
    except Exception as e:
        return -1, str(e)

def norm(s):
    try:
        obj = json.loads(s, parse_float=str, parse_int=str)
        if isinstance(obj, dict) and isinstance(obj.get("data"), list) and obj["data"] and \
           isinstance(obj["data"][0], dict) and "id" in obj["data"][0]:
            obj["data"] = sorted(obj["data"], key=lambda r: str(r.get("id")))
        return json.dumps(obj, sort_keys=True, ensure_ascii=False)
    except Exception:
        return "RAW:" + s.strip()

CASES = []
def C(name, path, body=None):
    CASES.append((name, path, body))

pg = {"current": 1, "size": 3}

# ── /chainCapitalAccounts ──
C("acc.list.all", "/saas-capital/chainCapitalAccounts/list", {})
C("acc.list.owner", "/saas-capital/chainCapitalAccounts/list", {"owner_type": 1})
C("acc.lbp.owner", "/saas-capital/chainCapitalAccounts/listByPage", {"ownerType": 1, "page": pg})
C("acc.lbp.noowner", "/saas-capital/chainCapitalAccounts/listByPage", {"page": pg})  # Java 语法错 999
C("acc.lbp.search", "/saas-capital/chainCapitalAccounts/listByPage", {"ownerType": 2, "searchStr": "银行", "page": pg})
C("acc.getById.1", "/saas-capital/chainCapitalAccounts/getById", {"id": 1})
C("acc.getById.miss", "/saas-capital/chainCapitalAccounts/getById", {"id": 999999999})
C("acc.delBatch.empty", "/saas-capital/chainCapitalAccounts/delBatch", {"ids": []})
C("acc.delBatch.null", "/saas-capital/chainCapitalAccounts/delBatch", {"ids": None})
C("acc.delBatch.noids", "/saas-capital/chainCapitalAccounts/delBatch", {})

# ── /chainCapitalBkcall ──
C("bk.list.all", "/saas-capital/chainCapitalBkcall/list", {})
C("bk.lbp.plain", "/saas-capital/chainCapitalBkcall/listByPage", {"page": pg})
C("bk.lbp.search", "/saas-capital/chainCapitalBkcall/listByPage", {"searchStr": "PAY", "page": pg})
C("bk.lbp.status0", "/saas-capital/chainCapitalBkcall/listByPage", {"status": "0", "page": pg})  # "0" 值跳过
C("bk.lbp.status", "/saas-capital/chainCapitalBkcall/listByPage", {"status": "FIN", "page": pg})
C("bk.lbp.orderType", "/saas-capital/chainCapitalBkcall/listByPage", {"orderType": "PAY", "page": pg})
C("bk.getById.1", "/saas-capital/chainCapitalBkcall/getById", {"id": 1})
C("bk.getById.msgtext", "/saas-capital/chainCapitalBkcall/getById", {"id": 2})
C("bk.delBatch.empty", "/saas-capital/chainCapitalBkcall/delBatch", {"ids": []})

# ── /chainCapitalTransinfo ──
C("tr.list.all", "/saas-capital/chainCapitalTransinfo/list", {})
C("tr.lbp.plain", "/saas-capital/chainCapitalTransinfo/listByPage", {"page": pg})
C("tr.lbp.search", "/saas-capital/chainCapitalTransinfo/listByPage", {"searchStr": "山西", "page": pg})
C("tr.lbp.t1", "/saas-capital/chainCapitalTransinfo/listByPage", {"transType": "1", "page": pg})
C("tr.lbp.t2", "/saas-capital/chainCapitalTransinfo/listByPage", {"transType": "2", "page": pg})
C("tr.lbp.t0", "/saas-capital/chainCapitalTransinfo/listByPage", {"transType": "0", "page": pg})
C("tr.lbp.range", "/saas-capital/chainCapitalTransinfo/listByPage", {"startTime": "2022-01-01", "endTime": "2026-12-31", "page": pg})
C("tr.lbp.busStatus", "/saas-capital/chainCapitalTransinfo/listByPage", {"busStatus": "FIN", "page": pg})
C("tr.getById.1", "/saas-capital/chainCapitalTransinfo/getById", {"id": 1})
C("tr.delBatch.empty", "/saas-capital/chainCapitalTransinfo/delBatch", {"ids": []})

# ── /chainCapitalCount ──
C("cnt.list.all", "/saas-capital/chainCapitalCount/list", {})
C("cnt.lbp.plain", "/saas-capital/chainCapitalCount/listByPage", {"page": pg})
C("cnt.lbp.search", "/saas-capital/chainCapitalCount/listByPage", {"searchStr": "HXJL", "page": pg})
C("cnt.lbp.createTime", "/saas-capital/chainCapitalCount/listByPage", {"createTime": "2021-04-25 00:00:00", "page": pg})
C("cnt.lbp.createTimeBad", "/saas-capital/chainCapitalCount/listByPage", {"createTime": "2021/04/25", "page": pg})
C("cnt.getById.1", "/saas-capital/chainCapitalCount/getById", {"id": 1})

# ── /chainCapitalCountMx(listByPage 恒 999) ──
C("mx.list.all", "/saas-capital/chainCapitalCountMx/list", {})
C("mx.lbp.plain", "/saas-capital/chainCapitalCountMx/listByPage", {"page": pg})
C("mx.lbp.search", "/saas-capital/chainCapitalCountMx/listByPage", {"searchStr": "HXJL", "page": pg})
C("mx.getById.1", "/saas-capital/chainCapitalCountMx/getById", {"id": 1})

# ── /chainCapitalDetailSettlement ──
C("ds.list.all", "/saas-capital/chainCapitalDetailSettlement/list", {})
C("ds.lbp.plain", "/saas-capital/chainCapitalDetailSettlement/listByPage", {"page": pg})
C("ds.lbp.search", "/saas-capital/chainCapitalDetailSettlement/listByPage", {"searchStr": "1", "page": pg})
C("ds.lbp.nullstr", "/saas-capital/chainCapitalDetailSettlement/listByPage", {"page": pg})
C("ds.lbp.inject", "/saas-capital/chainCapitalDetailSettlement/listByPage", {"searchStr": "' or '1'='1", "page": pg})
C("ds.getById.1", "/saas-capital/chainCapitalDetailSettlement/getById", {"id": 1})

# ── /chainBillSettRec ──
C("sr.list.all", "/saas-capital/chainBillSettRec/list", {})
C("sr.lbp.plain", "/saas-capital/chainBillSettRec/listByPage", {"page": pg})
C("sr.lbp.relDate", "/saas-capital/chainBillSettRec/listByPage", {"relDate": "2026-11-17", "page": pg})
C("sr.lbp.relDateBad", "/saas-capital/chainBillSettRec/listByPage", {"relDate": "2026/11/17", "page": pg})
C("sr.getById.1", "/saas-capital/chainBillSettRec/getById", {"id": 1})
C("sr.delBatch.empty", "/saas-capital/chainBillSettRec/delBatch", {"ids": []})
C("sr.relieve.empty", "/saas-capital/chainBillSettRec/relieve", {"ids": []})
C("sr.relieve.null", "/saas-capital/chainBillSettRec/relieve", {})

# ── /chainPayment(本地失败语义) ──
C("pay.payment.nopayee", "/saas-capital/chainPayment/payment", {})
C("pay.payment.ghost", "/saas-capital/chainPayment/payment", {"payeeEtc": "不存在企业XYZ"})
C("pay.payment.real", "/saas-capital/chainPayment/payment", {"payeeEtc": "山西国盛恒泰能源科技有限公司", "orderAmt": 1, "contractId": "1", "contractCode": "C", "orderSubtype": "DRT"})
C("pay.confirmPay.nocode", "/saas-capital/chainPayment/confirmPay", {"id": "1"})
C("pay.confirmPay.code", "/saas-capital/chainPayment/confirmPay", {"id": "1", "code": "123456"})
C("pay.cancelPay.nocode", "/saas-capital/chainPayment/cancelPay", {"refundList": []})
C("pay.cancelPay.code", "/saas-capital/chainPayment/cancelPay", {"code": "123456", "refundList": []})
C("pay.refund.nocode", "/saas-capital/chainPayment/refund", {"refundList": []})
C("pay.refund.code", "/saas-capital/chainPayment/refund", {"code": "123456", "refundList": []})
C("pay.orderDetails", "/saas-capital/chainPayment/orderDetails", {})
C("pay.orderDetails.full", "/saas-capital/chainPayment/orderDetails", {"mchNbr": "C1", "mchOrdNbr": "O1", "mchTransNbr": "T1"})
C("pay.billInfo", "/saas-capital/chainPayment/billInfo", {"orderDate": "20260101", "pageNum": 1, "pageSize": 5})
C("pay.callback", "/saas-capital/chainPayment/transMsgCallBack", {"msgNbr": "1", "status": "WCF", "orderType": "PAY"})

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
            print(f"  J: {nj[:350]}")
            print(f"  G: {ng[:350]}")
    else:
        print(f"ok   {name}")
print(f"\n{len(CASES)-fails}/{len(CASES)} passed")
sys.exit(1 if fails else 0)
