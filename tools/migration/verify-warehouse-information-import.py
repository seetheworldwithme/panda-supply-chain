from pathlib import Path
exec(Path(__file__).with_name('verify-warehouse-information.py').read_text().split('def main():')[0])
from xml.sax.saxutils import escape
from tempfile import TemporaryDirectory

def legacy_workbook(data):
 with TemporaryDirectory() as work:
  folder = Path(work)
  source = folder / 'fixture.xlsx'
  source.write_bytes(data)
  subprocess.run([
   os.environ.get('PANDA_SOFFICE_BIN', 'soffice'),
   '-env:UserInstallation=' + (folder / 'profile').as_uri(),
   '--headless', '--convert-to', 'xls:MS Excel 97', '--outdir', work, str(source),
  ], check=True, capture_output=True, timeout=30)
  return (folder / 'fixture.xls').read_bytes()


def workbook(blob,values):
 z=zipfile.ZipFile(io.BytesIO(blob));ns={'m':'http://schemas.openxmlformats.org/spreadsheetml/2006/main'}
 shared=[]
 if 'xl/sharedStrings.xml' in z.namelist():shared=[''.join(e.itertext()) for e in ET.fromstring(z.read('xl/sharedStrings.xml'))]
 root=ET.fromstring(z.read('xl/worksheets/sheet1.xml'));r=root.find('m:sheetData/m:row',ns)
 cells=[]
 for c in r:
  v=c.find('m:v',ns)
  title=shared[int(v.text)] if c.attrib.get('t')=='s' else ''.join(c.itertext())
  if title in values:
   ref=''.join(x for x in c.attrib['r'] if x.isalpha())+'2'
   cells.append(f'<c r="{ref}" t="inlineStr"><is><t>{escape(str(values[title]))}</t></is></c>')
 data=z.read('xl/worksheets/sheet1.xml').decode();i=data.index('</sheetData>');data=data[:i]+'<row r="2">'+''.join(cells)+'</row>'+data[i:]
 out=io.BytesIO()
 with zipfile.ZipFile(out,'w') as dst:
  for name in z.namelist():dst.writestr(name,data if name=='xl/worksheets/sheet1.xml' else z.read(name))
 return out.getvalue()
for ix,domain in enumerate(['Coal','Coking','Railway','Washing']):
 base='/xmmc'+domain+'Info';table='xmmc_'+domain.lower()+'_info';id=99977774+ix
 assert sql(f'SELECT COUNT(*) FROM {table} WHERE id={id}')=='0'
 template=req(8765,base+'/generateAndDownloadExcel',method='GET')[1]
 file=workbook(template,{'自增id':id,'*地区':'太原市','*县城':'古交市','地区':'太原市','县城':'古交市','县市':'古交市','市区':'太原市','县区':'古交市','所属市':'太原市','所属县市':'古交市'})
 ext = 'xlsx'
 if os.getenv('MIGRATION_XLS') == '1':
  file = legacy_workbook(file)
  ext = 'xls'
 bound='warehouse-migration-fixture'
 body=(f'--{bound}\r\nContent-Disposition: form-data; name="file"; filename="fixture.{ext}"\r\nContent-Type: application/octet-stream\r\n\r\n').encode()+file+f'\r\n--{bound}--\r\n'.encode()
 snapshots=[]
 try:
  for port in [8765,8799]:
   status,b=req(port,base+'/importExcel',body,ctype='multipart/form-data; boundary='+bound)
   j=norm(b)
   print(domain,port,status,j['code'])
   assert j['code']==200,(domain,port,j)
   id=int(j['data']['dataList'][0]['id'])
   assert 99977700<=id<99977800
   assert sql(f'SELECT COUNT(*) FROM {table} WHERE id={id}')=='1'
   assert sql(f'SELECT province_code,province,city_code,area_code FROM {table} WHERE id={id}')=='14\t山西省\t1401\t140181'
   for row in j['data']['dataList']:
    for key in ['id','importBatchNumber','createTime','updateTime']:row.pop(key,None)
   snapshots.append(j)
   sql(f'DELETE FROM {table} WHERE id={id}')
  assert snapshots[0]==snapshots[1],(domain,'import response mismatch',snapshots)
 finally:sql(f'DELETE FROM {table} WHERE id={id}')

print('PASS four information imports; 8 writes verified and cleaned')
