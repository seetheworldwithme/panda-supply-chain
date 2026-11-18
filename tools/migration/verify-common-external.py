import urllib.request,urllib.error,urllib.parse,json,os,subprocess,tempfile,pathlib,shutil

def sql(q):
 return subprocess.check_output(['mysql','-h127.0.0.1','-uroot','chain','-NB','-e',q],text=True).strip()
def login(base):
 data=urllib.parse.urlencode(dict(username='hyqcbjhd',password=os.environ['CMP_PASSWORD'],code='0000')).encode()
 return json.load(urllib.request.urlopen(urllib.request.Request(base+'/saas-auth/login',data,{'Content-Type':'application/x-www-form-urlencoded'})))['data']['token']
def get(base,path,t):
 with urllib.request.urlopen(urllib.request.Request(base+'/saas-common'+path,headers={'Authorization':'Bearer '+t}),timeout=15) as r:return r.status,r.read()
def post(base,path,body,t):
 with urllib.request.urlopen(urllib.request.Request(base+'/saas-common'+path,json.dumps(body).encode(),{'Authorization':'Bearer '+t,'Content-Type':'application/json'}),timeout=15) as r:return r.status,json.load(r)
def upload(base,path,t,fields={}):
 boundary='phase6fixtureboundary';body=b''
 for k,v in fields.items():body+=f'--{boundary}\r\nContent-Disposition: form-data; name="{k}"\r\n\r\n{v}\r\n'.encode()
 body+=f'--{boundary}\r\nContent-Disposition: form-data; name="file"; filename="phase6-fixture.txt"\r\nContent-Type: text/plain\r\n\r\n'.encode()+b'phase6-fixture\r\n'+f'--{boundary}--\r\n'.encode()
 with urllib.request.urlopen(urllib.request.Request(base+'/saas-common'+path,body,{'Authorization':'Bearer '+t,'Content-Type':'multipart/form-data; boundary='+boundary}),timeout=20) as r:return json.load(r)

D=tempfile.mkdtemp(prefix='phase6-ext-');p=pathlib.Path(D);(p/'phase6.mp4').write_bytes(b'0123456789');ID=99977781
sql(f"DELETE FROM chain_common_file WHERE id={ID}");sql(f"INSERT INTO chain_common_file(id,file_real_name,file_name,file_path,file_type) VALUES({ID},'fixture.mp4','phase6.mp4','{D}','5')")
created=[]
try:
 for base in ['http://localhost:8765','http://localhost:8799']:
  tok=login(base)
  st,b=get(base,f'/upload/imageShow?fileName={ID}',tok);assert st==200 and b==b'0123456789';print('PASS image stream fixture',base)
  st,b=get(base,f'/upload/getVideo?fileName={ID}',tok);assert st==200 and b==b'0123456789';print('PASS video bytes fixture',base)
  st,b=get(base,f'/upload/downById?fileName={ID}',tok);assert len(b)==1024 and b[:10]==b'0123456789';print('PASS download padded bytes fixture',base)
  _,r=post(base,'/upload/getByRelId',{'id':ID},tok);assert r['data']==[{'id':ID,'name':'fixture.mp4','fileName':'phase6.mp4','fileType':'5'}],r;print('PASS rel fixture',base)
  if base.endswith('8799'):
   _,r=post(base,'/upload/getShareToken',{'id':ID},tok);assert len(r['data'])==36;print('PASS share token in-memory fixture',base)
  else:
   try: post(base,'/upload/getShareToken',{'id':ID},tok)
   except urllib.error.HTTPError as e: assert e.code==500;print('KNOWN Redis unavailable: Java share token HTTP500')
  r=upload(base,'/upload/uploadFiles',tok)
  if base.endswith('8765') and r['data']==[]: print('KNOWN Java configured upload directory unavailable: returns []')
  else:
   assert r['code']==200 and len(r['data'])==1,r;f=r['data'][0];rid=int(f['id']);created.append(rid);row=sql(f'SELECT file_path,file_name,file_type FROM chain_common_file WHERE id={rid}').split('\t');assert row[2]=='5',row;fp=pathlib.Path(row[0])/row[1];fp=fp if fp.is_absolute() else pathlib.Path(__file__).resolve().parents[2]/'backend'/fp;assert fp.read_bytes()==b'phase6-fixture';fp.unlink();print('PASS uploadFiles DB + bytes',base)
 # Go chunk assembler, one part: returned id and effective DB flag, bytes then cleanup
 base='http://localhost:8799';tok=login(base);name='phase6-99977781.txt';r=upload(base,'/upload/bigUpload',tok,{'fileName':name,'chunkNumber':1,'totalChunks':1});assert r['code']==200 and r['data']['fileName']==name,r;rid=int(r['data']['id']);created.append(rid);row=sql(f'SELECT file_path,file_name,effective FROM chain_common_file WHERE id={rid}').split('\t');assert row[2]=='1';fp=pathlib.Path(row[0])/row[1];fp=fp if fp.is_absolute() else pathlib.Path(__file__).resolve().parents[2]/'backend'/fp;assert fp.read_bytes()==b'phase6-fixture';fp.unlink();print('PASS Go bigUpload DB + bytes')
finally:
 sql('DELETE FROM chain_common_file WHERE id IN ('+','.join(map(str,[ID]+created))+')');shutil.rmtree(D)
print('Positive local fixtures passed; two Java environment differences explicitly excluded; fixture rows/files removed')
