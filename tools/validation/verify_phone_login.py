"""只用本地隔离账号和无效手机号验证验证码登录；会清理名片和用户。"""
import argparse,json,re,secrets
from pathlib import Path
from contract_compare import decode,differences,request
from verify_auth_crud import sql


def main():
    p=argparse.ArgumentParser(description=__doc__);p.add_argument('--sessions',type=Path,required=True);p.add_argument('--output',type=Path,required=True);p.add_argument('--go-port',default='8798');args=p.parse_args();ss=json.loads(args.sessions.read_text())['ptadmin'];marker='audit1118-phone';phone='00000001118';results=[]
    assert sql("SELECT count(*) FROM chain_auth_user WHERE account='"+marker+"' OR phone='"+phone+"' OR remark5='"+marker+"'")=='0'
    owner=''
    try:
        body={'account':marker,'phone':phone,'password':secrets.token_urlsafe(24),'username':'本地验证码验收','level':1,'accountType':'3','channel':'2','isAgreeCall':False,'handleResult':'1','remark5':marker}
        status,_,blob=request(args.go_port,ss[args.go_port]['token'],{'method':'POST','path':'/saas-auth/chainAuthUser/add','body':body});assert status==200 and decode(blob)['code']=='200'
        owner=sql("SELECT id FROM chain_auth_user WHERE remark5='"+marker+"'");assert owner.isdigit()
        sql("INSERT INTO bus_business_card(user_id,phone,is_default,remark5) VALUES ("+owner+",'"+phone+"',1,'"+marker+"')")
        runs={}
        for port in ('8765',args.go_port):
            def call(path,body,token):
                status,_,raw=request(port,token,{'method':'POST','path':'/saas-auth'+path,'body':body});return status,decode(raw)
            status,login=call('/chainAuthChat/loginByAccount',{'account':marker},ss[port]['token']);assert login['code']=='200' and isinstance(login['data'],dict)
            token=login['data']['token'];cases=[]
            for path in ('/chainAuthChat/loginPhone','/chainAuthChat/loginChatPhone'):
                status,value=call('/chainAuthChat/code/phone',{'account':phone},token);code=value['data'];assert re.fullmatch('[0-9]{6}',code),'本机必须为不发送真实短信的验收配置'
                status,value=call(path,{'phone':phone,'code':code,'openid':marker,'pushClientid':marker},token)
                # 1:1 复刻 JAR 编译缺陷:两侧恒 500/code999,零落库。
                assert status==500 and value.get('code')=='999',f"JAR 缺陷路由两侧必须一致: {port} {path} HTTP{status} code={value.get('code')}"
                cases.append([status,value])
                status,value=call(path,{'phone':phone,'code':'invalid'},token)
                assert status==500 and value.get('code')=='999',f"JAR 缺陷路由两侧必须一致(invalid): {port} {path}"
                cases.append([status,value])
            sql("UPDATE chain_auth_user SET openid='"+marker+"' WHERE id="+owner)
            duplicate=marker+'-official'
            sql("INSERT INTO chain_auth_user(official_openid,remark5) VALUES ('"+duplicate+"','"+marker+"')")
            status,value=call('/wechat/updByOpenid',{'openid':marker,'officialOpenid':duplicate},token);assert value['code']=='200'
            assert sql("SELECT count(*) FROM chain_auth_user WHERE remark5='"+marker+"' AND phone IS NULL AND openid IS NULL")=='0'
            assert sql('SELECT official_openid FROM chain_auth_user WHERE id='+owner)==duplicate
            sql('UPDATE chain_auth_user SET official_openid=NULL WHERE id='+owner)
            cases.append([status,value])
            before_rule=sql('SELECT id,mo_get_sum_score FROM bus_score_rule ORDER BY id')
            sql("INSERT INTO bus_score_log(user_id,score_get_type,remark5) VALUES ("+owner+",1,'"+marker+"')")
            status,value=call('/chainAuthChat/updByOpenId',{'id':int(owner),'phone':phone,'openid':marker,'nickname':'本地资料更新验收'},token)
            assert value['code']=='200' and sql('SELECT nickname FROM chain_auth_user WHERE id='+owner)=='本地资料更新验收'
            assert sql('SELECT id,mo_get_sum_score FROM bus_score_rule ORDER BY id')==before_rule,'已领过注册积分不得再次扣全局积分池'
            sql("DELETE FROM bus_score_log WHERE user_id="+owner+" AND remark5='"+marker+"'")
            sql('UPDATE chain_auth_user SET nickname=NULL WHERE id='+owner)
            cases.append([status,value]);runs[port]=cases
        for i,(left,right) in enumerate(zip(runs['8765'],runs[args.go_port])):
            diff=differences(left,right)
            known=i<4 and left[0]==500 and left[1].get('code')=='999' and right[0]==500 and right[1].get('code')=='999'
            go_ok=(right[0]==500 and right[1].get('code')=='999') if i<4 else (right[0]==200 and right[1].get('code')=='200')
            results.append({'case':i,'differences':diff,'jar_compilation_defect_replicated':known,'go_verified':go_ok})
            print('phone/official fixture',i,'PASS (1:1 JAR 编译缺陷)' if known else ('PASS' if not diff else 'FAIL '+','.join(diff)),flush=True)

    finally:
        if owner:
            sql("DELETE FROM bus_score_log WHERE user_id="+owner+" AND remark5='"+marker+"'")
        if owner:sql('DELETE FROM bus_business_card WHERE user_id='+owner+" AND remark5='"+marker+"'")
        sql("DELETE FROM chain_auth_user WHERE remark5='"+marker+"'")
    args.output.write_text(json.dumps(results,ensure_ascii=False,indent=2)+'\n');return int(any(x['differences'] or not x['go_verified'] for x in results))

if __name__=='__main__':raise SystemExit(main())
