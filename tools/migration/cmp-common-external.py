import urllib.request,urllib.error,json,os

def req(base,path,body=None,token=None):
 h={}
 if token:h['Authorization']='Bearer '+token
 if body is not None:h['Content-Type']='application/json'
 r=urllib.request.Request(base+path,None if body is None else json.dumps(body).encode(),h)
 try:
  with urllib.request.urlopen(r,timeout=15) as r:return r.status,r.read().decode()
 except urllib.error.HTTPError as e:return e.code,e.read().decode()
def login(base):
 data=urllib.parse.urlencode(dict(username=os.getenv('CMP_USER','hyqcbjhd'),password=os.environ['CMP_PASSWORD'],code='0000')).encode()
 r=urllib.request.Request(base+'/saas-auth/login',data,{'Content-Type':'application/x-www-form-urlencoded'})
 return json.load(urllib.request.urlopen(r))['data']['token']
def norm(b):
 try:return json.dumps(json.loads(b,parse_float=str,parse_int=str),sort_keys=True,ensure_ascii=False)
 except:return b
cases=[('GET','/upload/preShareFIle',None),('POST','/upload/preFIle',{'file_id':99977799}),('POST','/upload/getShareToken',{}),('POST','/upload/getByRelId',{'id':99977799}),('POST','/upload/getError',{}),('POST','/TenOcrController/getOcrResult',{}),('POST','/TenOcrController/getOcrResult',{'id':99977799,'type':1}),('GET','/upload/imageShow?fileName=__phase6_missing__.jpg',None),('GET','/upload/imageShow?fileName=__phase6_missing__.jpg&credit=missing',None),('GET','/upload/downById?fileName=__phase6_missing__.jpg',None),('GET','/upload/zjjgDownById?fileName=__phase6_missing__.jpg',None),('GET','/upload/previewStream?fileName=__phase6_missing__.pdf',None),('GET','/upload/zjjgPreviewStream?fileName=__phase6_missing__.pdf',None),('GET','/upload/getVideo?fileName=__phase6_missing__.mp4',None),('GET','/upload/play?fileName=__phase6_missing__.mp4',None),('GET','/download',None),('GET','/api/v1/download/getFile',None),('GET','/api/v1/download/limit',None)]
bases=['http://localhost:8765','http://localhost:8799'];ts=[login(b) for b in bases];ok=0
for method,path,body in cases:
 rs=[req(b,'/saas-common'+path,body,t) for b,t in zip(bases,ts)]
 same=rs[0][0]==rs[1][0] and norm(rs[0][1])==norm(rs[1][1]);ok+=same
 print('PASS' if same else 'FAIL',path,rs if not same else rs[0][0])
print(f'{ok}/{len(cases)}');raise SystemExit(ok!=len(cases))
