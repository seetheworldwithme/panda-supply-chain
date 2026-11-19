"""本机鉴权、畸形请求和随机图像的协议断言；不保存响应正文。"""
import argparse,io,json,struct,urllib.request,urllib.error
from pathlib import Path
from contract_compare import compare,decode,Number


def fetch(port,token,path,body,method='POST'):
    headers={'Content-Type':'application/json'}
    if token:headers['Authorization']='Bearer '+token
    req=urllib.request.Request('http://127.0.0.1:'+port+path,data=body,headers=headers,method=method)
    try:r=urllib.request.urlopen(req,timeout=30)
    except urllib.error.HTTPError as e:r=e
    with r:return r.status,dict(r.headers),r.read()


def jpeg_size(data):
    assert data[:2]==b'\xff\xd8'
    i=2
    while i<len(data):
        assert data[i]==255
        while data[i]==255:i+=1
        marker=data[i];i+=1
        if marker in (216,217):continue
        length=int.from_bytes(data[i:i+2],'big')
        if marker in (192,193,194):return tuple(reversed(struct.unpack('>HH',data[i+3:i+7])))
        i+=length
    raise ValueError('缺少 JPEG 尺寸')


def main():
    p=argparse.ArgumentParser(description=__doc__);p.add_argument('--sessions',type=Path,required=True);p.add_argument('--output',type=Path,required=True);p.add_argument('--go-port',default='8798');args=p.parse_args();session=json.loads(args.sessions.read_text())['ptadmin'];results=[]
    def record(name,diff):
        results.append({'case':name,'differences':diff});print(name,'PASS' if not diff else 'FAIL '+','.join(diff),flush=True)
    for body in (b'',b'null',b'[]',b'{',b'"string"'):
        runs=[fetch(port,session[port]['token'],'/saas-common/chainAreaInfo/getById',body) for port in ('8765',args.go_port)]
        diff=compare(*runs)
        if any(x[0]!=500 or decode(x[2]).get('code')!=Number('999') for x in runs):diff.append('expected-Java-binding-failure')
        record('invalid-body:'+body.decode(),diff)
    runs=[fetch(port,session[port]['token'],'/saas-common/chainAreaInfo/getById',b'{}{}') for port in ('8765',args.go_port)]
    record('Java-first-JSON-value',compare(*runs))
    path='/saas-common/chainAreaInfo/list';body=b'{"id":99977799}'
    runs=[fetch(port,session['8765']['token'],path,body) for port in ('8765',args.go_port)]
    diff=compare(*runs)
    if any(decode(x[2]).get('code')!=Number('200') for x in runs):diff.append('legacy-token-rejected')
    record('original-Java-token-on-both',diff)
    runs=[fetch(port,'',path,b'{') for port in ('8765',args.go_port)]
    diff=compare(*runs)
    if any(decode(x[2]).get('code')!=Number('403') for x in runs):diff.append('anonymous-not-denied')
    record('authentication-before-body-binding',diff)
    for query,size in [('',(67,23)),('?width=90&height=30',(90,30))]:
        runs=[fetch(port,'','/saas-auth/code/imageGet'+query,None,'GET') for port in ('8765',args.go_port)]
        diff=[]
        for label,run in zip(('java','go'),runs):
            headers={k.lower():v for k,v in run[1].items()}
            if run[0]!=200 or headers.get('content-type')!='image/jpeg' or jpeg_size(run[2])!=size or headers.get('cache-control')!='no-store, no-cache':diff.append(label+':image-contract')
        record('imageGet'+query,diff)
    args.output.write_text(json.dumps(results,ensure_ascii=False,indent=2)+'\n');return int(any(x['differences'] for x in results))

if __name__=='__main__':raise SystemExit(main())
