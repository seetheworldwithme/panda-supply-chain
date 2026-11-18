import json,urllib.request,urllib.error,urllib.parse,subprocess,io,zipfile,xml.etree.ElementTree as ET
from decimal import Decimal
import os

def token(port):
 filename=os.environ.get('JAVA_TOKEN_FILE' if port==8765 else 'GO_TOKEN_FILE')
 if filename:return open(filename).read().strip()
 data=urllib.parse.urlencode({'username':os.environ['MIGRATION_USERNAME'],'password':os.environ['MIGRATION_PASSWORD'],'code':'0000'}).encode()
 request=urllib.request.Request('http://localhost:'+str(port)+'/saas-auth/login',data=data,headers={'Content-Type':'application/x-www-form-urlencoded'})
 return json.load(urllib.request.urlopen(request,timeout=10))['data']['token']
TOK={port:token(port) for port in [8765,8799]}

def req(port,path,body={},method='POST',ctype='application/json'):
 data=body if isinstance(body,bytes) else json.dumps(body).encode()
 r=urllib.request.Request('http://localhost:'+str(port)+'/saas-warehouse'+path,data=None if method=='GET' else data,method=method,headers={'Authorization':'Bearer '+TOK[port],'Content-Type':ctype})
 try:
  with urllib.request.urlopen(r,timeout=30) as x:return x.status,x.read()
 except urllib.error.HTTPError as x:return x.code,x.read()

def sql(q):return subprocess.check_output(['mysql','-uroot','-N','-B','chain','-e',q],text=True).strip()
def norm(b):return json.loads(b,parse_float=Decimal)
def cmp(path,body={},method='POST'):
 a,b=req(8765,path,body,method),req(8799,path,body,method)
 assert a[0]==b[0] and norm(a[1])==norm(b[1]),(path,body,a[0],b[0],a[1][:120],b[1][:120])

def main():
 N=0
 for ix,domain in enumerate(['Coal','Coking','Railway','Washing']):
  base='/xmmc'+domain+'Info'
  for body in [{},{'page':{'current':1,'size':3}},{'city':'太原','page':{'current':1,'size':3}},{'factoryTaxPriceMin':1,'page':{'current':1,'size':3}},{'searchStr':'ignored','page':{'current':1,'size':3}}]:
   cmp(base+'/listByPage',body);N+=1
  a,b=req(8765,base+'/generateAndDownloadExcel',method='GET'),req(8799,base+'/generateAndDownloadExcel',method='GET')
  def head(blob):
   z=zipfile.ZipFile(io.BytesIO(blob));ns={'m':'http://schemas.openxmlformats.org/spreadsheetml/2006/main'}
   shared=[]
   if 'xl/sharedStrings.xml' in z.namelist():shared=[''.join(e.itertext()) for e in ET.fromstring(z.read('xl/sharedStrings.xml'))]
   sheet=ET.fromstring(z.read('xl/worksheets/sheet1.xml'))
   out=[]
   for c in sheet.find('m:sheetData/m:row',ns):
    if c.attrib.get('t')=='s':out.append(shared[int(c.find('m:v',ns).text)])
    elif c.attrib.get('t')=='inlineStr':out.append(''.join(c.find('m:is',ns).itertext()))
    else:out.append(c.find('m:v',ns).text)
   return out
  assert a[0]==200 and b[0]==200 and head(a[1])==head(b[1]),(base,'header mismatch',head(a[1]),head(b[1]));N+=1
  cmp(base+'/importExcel',{});N+=1
  cmp(base+'/updateById',{'id':99977799});N+=1
  table='xmmc_'+domain.lower()+'_info';id=99977770+ix
  try:
   results=[]
   for port in [8765,8799]:
    sql(f'DELETE FROM {table} WHERE id={id}; INSERT INTO {table}(id,city,area,update_time) VALUES ({id},"六期验收","测试县","2025-01-01 00:00:00"); DELETE FROM bus_common_log WHERE pro_id={id}')
    status,body=req(port,base+'/updateById',{'id':id,'city':'六期更新'})
    assert status==200 and norm(body)['code']==200,(base,port,body)
    assert sql(f'SELECT city FROM {table} WHERE id={id}')=='六期更新'
    assert sql(f'SELECT COUNT(*) FROM bus_common_log WHERE pro_id={id} AND operation_type=2')=='1'
    results.append((norm(body),sql(f'SELECT user_id,user_name,user_phone,pro_id,operation_type FROM bus_common_log WHERE pro_id={id}')))
   assert results[0]==results[1],(base,results)
   N+=1
  finally:
   sql(f'DELETE FROM {table} WHERE id={id}; DELETE FROM bus_common_log WHERE pro_id={id}')
 cmp('/wechat/transfer');N+=1
 print('PASS',N,'information/peripheral assertions; fixture rows and logs cleaned')

if __name__=="__main__":main()
