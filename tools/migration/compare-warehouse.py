#!/usr/bin/env python3
# 六期批1 读对比:Java 8765 vs Go 8799,同库同接口全量 JSON 对比
import hashlib
from pathlib import Path
import json, urllib.request, urllib.error, sys, os, urllib.parse

JAVA, GO = "http://localhost:8765", "http://localhost:8799"

def login(base):
    data = urllib.parse.urlencode({"username":os.environ["MIGRATION_USER"], "password":os.environ["MIGRATION_PASSWORD"], "code":"0000"}).encode()
    req = urllib.request.Request(base + "/saas-auth/login", data=data,
        headers={"Content-Type": "application/x-www-form-urlencoded"})
    return json.load(urllib.request.urlopen(req, timeout=10))["data"]["token"]

JT, GT = login(JAVA), login(GO)


def post(base, tok, path, body, raw=False):
    data = body if raw else json.dumps(body, ensure_ascii=False).encode()
    req = urllib.request.Request(base + path, data=data,
        headers={"Authorization": "Bearer " + tok, "Content-Type": "application/json"})
    try:
        with urllib.request.urlopen(req, timeout=25) as r:
            return r.status, r.read().decode("utf-8", "replace")
    except urllib.error.HTTPError as e:
        return e.code, e.read().decode("utf-8", "replace")
    except Exception as e:
        return -1, str(e)

def norm(s):
    try:
        obj = json.loads(s, parse_float=lambda n: {"$number":n}, parse_int=lambda n: {"$number":n})
        return json.dumps(obj, sort_keys=True, ensure_ascii=False)
    except Exception:
        return "RAW:" + s

CASES = json.load(open(sys.argv[1] if len(sys.argv) > 1 else Path(__file__).with_name("cases") / "phase6-cases-b3.json"))
from concurrent.futures import ThreadPoolExecutor
ok = bad = 0
fails = []
lock = __import__("threading").Lock()

def run(case):
    global ok, bad
    path, body = case
    sj, rj = post(JAVA, JT, path, body)
    sg, rg = post(GO, GT, path, body)
    same = sj == sg and norm(rj) == norm(rg)
    with lock:
        if same:
            ok += 1
        else:
            bad += 1
            fails.append((path, body, sj, sg, rj, rg))

with ThreadPoolExecutor(max_workers=1) as ex:
    list(ex.map(run, CASES))
print(f"\n===== {ok} 一致 / {bad} 不一致 =====")
for f in fails:
    print(f"\n[X] {f[0]} body={json.dumps(f[1], ensure_ascii=False)}\n  java({f[2]}): {len(f[4])} chars sha256={hashlib.sha256(f[4].encode()).hexdigest()}\n  go({f[3]}):   {len(f[5])} chars sha256={hashlib.sha256(f[5].encode()).hexdigest()}")

sys.exit(1 if bad else 0)
