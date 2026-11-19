"""本机隔离账号登录回归；口令随机生成，仅保留响应差异路径。"""
import argparse,base64,json,secrets
from pathlib import Path
from contract_compare import decode,differences,request,Number
from verify_auth_crud import sql


def normalize(value):
    if isinstance(value,dict):
        result={k:normalize(v) for k,v in value.items()}
        if 'token' in result:
            claims=json.loads(base64.urlsafe_b64decode(result['token'].split('.')[1]+'==='))
            assert claims['sub']=='audit1118-login' and claims['userId']==row_id
            result['token']='{validated-token}'
        return result
    if isinstance(value,list):return [normalize(x) for x in value]
    return value


def main():
    global row_id
    p=argparse.ArgumentParser(description=__doc__);p.add_argument('--sessions',type=Path,required=True);p.add_argument('--output',type=Path,required=True);p.add_argument('--go-port',default='8798');args=p.parse_args()
    session=json.loads(args.sessions.read_text())['ptadmin'];marker='audit1118-login';password=secrets.token_urlsafe(24);results=[]
    assert sql("SELECT count(*) FROM chain_auth_user WHERE account='"+marker+"' OR remark5='"+marker+"'")=='0'
    assert sql("SELECT count(*) FROM chain_auth_chat_authorization WHERE app_identification='"+marker+"' OR remark5='"+marker+"'")=='0'
    try:
        case={'method':'POST','path':'/saas-auth/chainAuthUser/add','body':{'account':marker,'username':'本地登录验收','password':password,'phone':'19900001118','handleResult':'1','level':1,'accountType':'3','isAgreeCall':False,'applys':'2','remark5':marker}}
        status,_,blob=request(args.go_port,session[args.go_port]['token'],case);assert status==200 and decode(blob)['code']=='200'
        row_id=sql("SELECT id FROM chain_auth_user WHERE remark5='"+marker+"'");assert row_id.isdigit()
        if sql('SELECT count(*) FROM chain_basic_enterprise WHERE id=405')=='1':
            sql("UPDATE chain_auth_user SET company_id=405 WHERE remark5='"+marker+"'")
        sql("INSERT INTO chain_auth_chat_authorization(app_identification,apply_name,secret,remark5,create_time,update_time) VALUES ('"+marker+"','本地测试','本地虚拟配置','"+marker+"','2020-01-02 03:04:05','2020-01-02 03:04:05')")
        cases=[('/chainAuthChat/loginByAccount',{'account':marker},True),('/chainAuthChat/buslogin',{'phone':'19900001118','level':1},True),('/chainAuthUser/loginCarSystem',{'account':marker,'password':password},True),('/chainAuthUser/loginCarSystem',{'account':marker,'password':'invalid'},False),('/chainAuthChat/getAppInfo',{'appIdentification':'audit1118-no-app'},False),('/chainAuthChat/getAppInfo',{'appIdentification':marker},True)]
        for path,body,nonempty in cases:
            case={'method':'POST','path':'/saas-auth'+path,'body':body};runs=[];statuses=[]
            for port in ('8765',args.go_port):
                status,_,blob=request(port,session[port]['token'],case);value=decode(blob)
                if nonempty:assert value['code']=='200' and isinstance(value.get('data'),dict),'必须真正登录成功'
                runs.append(normalize(value));statuses.append(status)
            diff=differences(*runs)
            if statuses[0]!=statuses[1]:diff.append('HTTP:status')
            results.append({'path':path,'success_required':nonempty,'differences':diff})
            print(path,'PASS' if not diff else 'FAIL '+','.join(diff),flush=True)
    finally:
        sql("DELETE FROM chain_auth_user WHERE account='"+marker+"' AND remark5='"+marker+"'")
        sql("DELETE FROM chain_auth_chat_authorization WHERE app_identification='"+marker+"' AND remark5='"+marker+"'")
    args.output.write_text(json.dumps(results,ensure_ascii=False,indent=2)+'\n');return int(any(x['differences'] for x in results))

if __name__=='__main__':raise SystemExit(main())
