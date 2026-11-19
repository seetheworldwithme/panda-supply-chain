"""本机文件 fixture 验证图片回显，不读取既有私有附件。"""
import argparse,json,tempfile
from pathlib import Path
from contract_compare import request,compare
from verify_auth_crud import sql


def main():
    p=argparse.ArgumentParser(description=__doc__);p.add_argument('--sessions',type=Path,required=True);p.add_argument('--output',type=Path,required=True);p.add_argument('--go-port',default='8798');args=p.parse_args();ss=json.loads(args.sessions.read_text())['ptadmin'];marker='audit1118-image';results=[]
    assert sql("SELECT count(*) FROM chain_common_file WHERE remark5='"+marker+"'")=='0'
    assert sql('SELECT count(*) FROM chain_common_file WHERE relation_id=99977799')=='0'
    with tempfile.TemporaryDirectory(prefix='panda-image-fixture-') as temp:
        Path(temp,'fixture.bin').write_bytes(b'local-fixture-1118')
        try:
            for typ in ('3','4'):
                sql("INSERT INTO chain_common_file(file_path,file_name,file_type,relation_id,remark5) VALUES ('"+temp+"','fixture.bin','"+typ+"',99977799,'"+marker+"')")
            for i,body in enumerate([{'relation_id':99977799},{'id':99977799},{'relation_id':99977799,'file_type':'4'}]):
                case={'method':'POST','path':'/saas-common/upload/getImgByRelId','body':body}
                runs=[request(port,ss[port]['token'],case) for port in ('8765',args.go_port)];diff=compare(*runs)
                results.append({'case':i,'differences':diff});print('image fixture',i,'PASS' if not diff else 'FAIL '+','.join(diff),flush=True)
        finally:sql("DELETE FROM chain_common_file WHERE remark5='"+marker+"'")
    args.output.write_text(json.dumps(results,ensure_ascii=False,indent=2)+'\n');return int(any(x['differences'] for x in results))

if __name__=='__main__':raise SystemExit(main())
