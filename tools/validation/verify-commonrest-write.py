import os,json,urllib.request,urllib.error,subprocess,sys
from pathlib import Path
# Use fixed fixture IDs only; SQL setup is local and always cleaned in finally.
import json,urllib.request,urllib.error,os,sys
from pathlib import Path
# Secrets are provided by the local acceptance environment, never stored in this script.
def login(port):
 data=urllib.parse.urlencode({'username':os.environ.get('CMP_USER','hyqcbjhd'),'password':os.environ['CMP_PASSWORD'],'code':'0000'}).encode()
 return json.load(urllib.request.urlopen(urllib.request.Request(f'http://localhost:{port}/saas-auth/login',data=data)))['data']['token']
from urllib import parse
jt,gt=login(8765),login(8799)

def sql(q):
 p=subprocess.run(['mysql','-uroot','-h127.0.0.1','chain','-N','-B','-e',q],capture_output=True,text=True)
 if p.returncode: raise RuntimeError(p.stderr)
 return p.stdout.strip()
def post(port,tok,p,b):
 req=urllib.request.Request(f'http://localhost:{port}/saas-common'+p,data=json.dumps(b).encode(),headers={'Authorization':'Bearer '+tok,'Content-Type':'application/json'})
 try:r=urllib.request.urlopen(req,timeout=20)
 except urllib.error.HTTPError as e:r=e
 return r.status,r.read().decode()
def norm(x):return json.dumps(json.loads(x,parse_float=str,parse_int=str),sort_keys=True,ensure_ascii=False)
names=['ChainMsg','ChainAreaInfo','ChainCommonAddress','ChainCommonImg','ChainCommonLog','ChainCommonTemplate','ChainCommonTempmx','ChainScreenConfig','ChainRequestLog','ChainCommonMqlog','ChainBasicBussiness','ChainTlzfOrder','ChainTlzfLog','ChainTlzfContractLog']
snake=lambda n: ''.join('_'+c.lower() if c.isupper() and i else c.lower() for i,c in enumerate(n))
passed=0
for n in names:
 table=snake(n);path='/'+n[0].lower()+n[1:];fid=99977770
 if sql(f"SELECT count(*) FROM information_schema.tables WHERE table_schema='chain' AND table_name='{table}'")=='0':
  a,x=post(8765,jt,path+'/updateById',{'id':fid,'remark1':'phase6-common-updated'});b,y=post(8799,gt,path+'/updateById',{'id':fid,'remark1':'phase6-common-updated'})
  assert a==b and norm(x)==norm(y),(n,a,b,x,y)
  passed+=1;print('MISSING_TABLE_MATCH',n,a,flush=True);continue
 try:
  result=[]
  for port,tok in [(8765,jt),(8799,gt)]:
   sql(f"DELETE FROM {table} WHERE id={fid}")
   cols='id,remark1,create_time';vals=f"{fid},'phase6-common-fixture','2020-01-02 03:04:05'"
   if n=='ChainAreaInfo': cols+=',code';vals+=',99977770'
   sql(f"INSERT INTO {table}({cols})VALUES({vals})")
   a,x=post(port,tok,path+'/updateById',{'id':fid,'remark1':'phase6-common-updated','updateTime':'2020-01-02T04:05:06' if n not in ['ChainCommonMqlog','ChainTlzfOrder','ChainTlzfLog','ChainTlzfContractLog'] else '2020-01-02 04:05:06'})
   db=sql(f"SELECT remark1,update_time FROM {table} WHERE id={fid}")
   b,y=post(port,tok,path+'/deleteById',{'id':fid})
   count=sql(f"SELECT count(*) FROM {table} WHERE id={fid}")
   result.append((a,norm(x),db,b,norm(y),count))
  if result[0]!=result[1]:print('FAIL',n,result,flush=True)
  else:passed+=1; print('PASS',n,result[0][0],result[0][-1],flush=True)
 finally:sql(f"DELETE FROM {table} WHERE id={fid}")
print('RESULT',passed,'/',len(names));sys.exit(passed!=len(names))
