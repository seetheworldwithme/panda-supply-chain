import json,urllib.request,urllib.error,os,sys
from pathlib import Path
# Secrets are provided by the local acceptance environment, never stored in this script.
def login(port):
 data=urllib.parse.urlencode({'username':os.environ.get('CMP_USER','hyqcbjhd'),'password':os.environ['CMP_PASSWORD'],'code':'0000'}).encode()
 return json.load(urllib.request.urlopen(urllib.request.Request(f'http://localhost:{port}/saas-auth/login',data=data)))['data']['token']
from urllib import parse
jt,gt=login(8765),login(8799)
names=['busMqConsumeLog','chainMsg','chainAreaInfo','chainCommonAddress','chainCommonImg','chainCommonLog','chainCommonTemplate','chainCommonTempmx','chainScreenConfig','chainRequestLog','chainCommonMqlog','chainBasicBussiness','chainTlzfOrder','chainTlzfLog','chainTlzfContractLog']
cases=[]
for name in names:
 for end,body in [('list',{'id':99977799}),('getById',{'id':1}),('getById',{'id':99977799}),('listByPage',{}),('listByPage',{'searchStr':'','page':{'current':1,'size':3}}),('listByPage',{'searchStr':'1','page':{'current':1,'size':2}}),('listByPage',{'searchStr':'','page':{'current':1,'size':3},'startTime':'2020-01-01','endTime':'2030-01-01'}),('listByPage',{'searchStr':'','id':1,'page':{'current':1,'size':3}})]:
  cases.append((f'/{name}/{end}',body))
for p,b in [('/chainAreaInfo/selByPet',{'parentCode':'0'}),('/chainBasicMainPage/selProvinces',{}),('/chainBasicMainPage/selectCount',{}),('/chainBasicMainPage/ptlist',{'searchStr':''}),('/chainBasicMainPage/countPtLogistArea',{}),('/chainTlzfOrder/selSfEtc',{}),('/chainTlzfOrder/selSfEtc',{'saleCompanyId':405}),('/chainTlzfOrder/selPtAmt',{}),('/chainTlzfOrder/selRelTraceFee',{'contractId':99977799})]:cases.append((p,b))
def post(port,tok,p,b):
 req=urllib.request.Request(f'http://localhost:{port}/saas-common'+p,data=json.dumps(b).encode(),headers={'Authorization':'Bearer '+tok,'Content-Type':'application/json'})
 try:r=urllib.request.urlopen(req,timeout=20)
 except urllib.error.HTTPError as e:r=e
 return r.status,r.read().decode()
def norm(x):
 try:return json.dumps(json.loads(x,parse_float=str,parse_int=str),sort_keys=True,ensure_ascii=False)
 except:return x
fail=[]
for p,b in cases:
 a,x=post(8765,jt,p,b);c,y=post(8799,gt,p,b)
 if a!=c or norm(x)!=norm(y):fail.append((p,b,a,c,x,y));print('FAIL',p,b,a,c,x[:180],y[:180],flush=True)
print('RESULT',len(cases)-len(fail),'/',len(cases),flush=True)
Path(os.environ.get('CMP_FAILURES','/tmp/commonrest-failures.json')).write_text(json.dumps(fail,ensure_ascii=False,indent=2))
sys.exit(bool(fail))
