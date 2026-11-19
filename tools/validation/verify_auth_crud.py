"""Local fixture-only Java/Go CRUD contract assertions. Never uses remote MySQL.

Requires --sessions (private token JSON), --go-port and --output. Each fixture is
guarded by a unique remark5 marker; finally deletes only rows with that marker.
Generated IDs are normalized after checking the returned ID against the DB.
"""
import argparse
import json
import subprocess
from pathlib import Path
from contract_compare import Number, decode, differences, request

DOMAINS = ["AuthenLog", "Authority", "AuthorityRole", "AuthorityUser", "MenuRole", "MenuUser", "Secret", "UseRole"]


def sql(query):
    return subprocess.check_output(["mysql", "--protocol=SOCKET", "--socket=/tmp/mysql.sock", "-uroot", "-N", "-B", "--raw", "chain", "-e", query], text=True).strip()


def normalize(value, row_id, key=""):
    if isinstance(value, Number):
        return Number("0") if key == "id" and value == row_id else value
    if isinstance(value, str):
        return "删除成功，id为：{fixture-id}" if value == "删除成功，id为：" + row_id else value
    if isinstance(value, dict):
        return {k: normalize(v, row_id, k) for k, v in value.items()}
    if isinstance(value, list):
        return [normalize(v, row_id) for v in value]
    return value


def main():
    p = argparse.ArgumentParser(description=__doc__)
    p.add_argument("--sessions", type=Path, required=True)
    p.add_argument("--go-port", default="8798")
    p.add_argument("--output", type=Path, required=True)
    args = p.parse_args()
    sessions = json.loads(args.sessions.read_text())["ptadmin"]
    results = []
    for i, domain in enumerate(DOMAINS):
        name = "chainAuth" + domain
        table = "chain_auth_" + "".join(("_" + c.lower()) if c.isupper() else c for c in domain).lstrip("_")
        marker = "audit1118-crud-" + str(i)
        predicate = "remark5='" + marker + "'"
        assert sql("SELECT count(*) FROM " + table + " WHERE " + predicate) == "0", "fixture marker collision"
        snapshots = []
        try:
            for port in ("8765", args.go_port):
                update_time = "2020-01-02 03:04:05" if domain in ("AuthenLog", "MenuUser") else "2020-01-02T03:04:05"
                body = {"remark5": marker, "createTime": "2020-01-02 03:04:05", "updateTime": update_time}
                case = {"path": "/saas-auth/" + name + "/add", "method": "POST", "body": body}
                status, _, raw = request(port, sessions[port]["token"], case)
                value = decode(raw)
                assert status == 200 and value["code"] == "200", f"{name} {port} add failed: HTTP{status}, code={value.get('code')}"
                row_id = sql("SELECT id FROM " + table + " WHERE " + predicate)
                assert row_id and "\n" not in row_id, "must create exactly one fixture"
                if isinstance(value.get("data"), dict):
                    assert value["data"]["id"] == row_id, "echoed ID must reference inserted row"
                run = [normalize(value, row_id)]
                run.append(sql("SELECT * FROM " + table + " WHERE " + predicate).split("\t")[1:])
                case["path"] = "/saas-auth/" + name + "/getById"
                case["body"] = {"id": int(row_id)}
                status, _, raw = request(port, sessions[port]["token"], case)
                assert status == 200
                run.append(normalize(decode(raw), row_id))
                case["path"] = "/saas-auth/" + name + "/updateById"
                case["body"] = {"id": int(row_id), "remark1": "fixture updated", "updateTime": update_time}
                status, _, raw = request(port, sessions[port]["token"], case)
                value = decode(raw)
                assert status == 200 and value["code"] == "200", "fixture update failed"
                run.append(normalize(value, row_id))
                assert sql("SELECT remark1 FROM " + table + " WHERE " + predicate) == "fixture updated"
                run.append(sql("SELECT * FROM " + table + " WHERE " + predicate).split("\t")[1:])
                case["path"] = "/saas-auth/" + name + "/deleteById"
                case["body"] = {"id": int(row_id)}
                status, _, raw = request(port, sessions[port]["token"], case)
                assert status == 200
                run.append(normalize(decode(raw), row_id))
                assert sql("SELECT count(*) FROM " + table + " WHERE " + predicate) == "0", "fixture delete failed"
                snapshots.append(run)
            diff = differences(snapshots[0], snapshots[1])
            results.append({"controller": name, "differences": diff, "db_assertions": "add/update/delete", "cleaned": True})
            print(name, "PASS" if not diff else "FAIL " + ",".join(diff), flush=True)
        finally:
            sql("DELETE FROM " + table + " WHERE " + predicate)
    marker="audit1118-batch-rollback"
    predicate="remark5='"+marker+"'"
    assert sql("SELECT count(*) FROM chain_auth_authen_log WHERE "+predicate)=="0"
    rollback=[]
    try:
        for port in ('8765',args.go_port):
            sql("INSERT INTO chain_auth_authen_log(remark1,remark5) VALUES ('before','"+marker+"'),('before','"+marker+"')")
            ids=[int(x) for x in sql("SELECT id FROM chain_auth_authen_log WHERE "+predicate+" ORDER BY id").splitlines()]
            status,_,raw=request(port,sessions[port]['token'],{'method':'POST','path':'/saas-auth/chainAuthAuthenLog/updBatch','body':[{'id':ids[0],'remark1':'after'},{'id':ids[1],'remark1':'x'*5000}]})
            unchanged=sql("SELECT count(*) FROM chain_auth_authen_log WHERE "+predicate+" AND remark1='before'")=="2"
            rollback.append(status==500 and unchanged)
            sql("DELETE FROM chain_auth_authen_log WHERE "+predicate)
        diff=[] if all(rollback) else ['$database:batch-rollback']
        results.append({'controller':'chainAuthAuthenLog/updBatch','differences':diff,'java_rollback':rollback[0],'go_rollback':rollback[1],'cleaned':True})
        print('chainAuthAuthenLog batch SQL failure rollback', 'PASS' if not diff else 'FAIL',flush=True)
    finally:
        sql("DELETE FROM chain_auth_authen_log WHERE "+predicate)
    args.output.write_text(json.dumps(results, ensure_ascii=False, indent=2) + "\n")
    return int(any(x["differences"] for x in results))


if __name__ == "__main__":
    raise SystemExit(main())
