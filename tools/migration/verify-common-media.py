"""Local Go media upload DB assertions. CMP_PASSWORD must be provided; fixtures cleaned."""
import urllib.request,urllib.parse,json,os,subprocess,io,zipfile,base64,pathlib,shutil
ROOT=pathlib.Path(__file__).resolve().parents[2];BASE='http://localhost:8799'
def sql(q):return subprocess.check_output(['mysql','-h127.0.0.1','-uroot','chain','-NB','-e',q],text=True).strip()
def login():
 data=urllib.parse.urlencode(dict(username='hyqcbjhd',password=os.environ['CMP_PASSWORD'],code='0000')).encode();return json.load(urllib.request.urlopen(urllib.request.Request(BASE+'/saas-auth/login',data,{'Content-Type':'application/x-www-form-urlencoded'})))['data']['token']
def upload(name,b,t):
 boundary='phase6media';body=f'--{boundary}\r\nContent-Disposition: form-data; name="file"; filename="{name}"\r\nContent-Type: application/octet-stream\r\n\r\n'.encode()+b+f'\r\n--{boundary}--\r\n'.encode();r=urllib.request.Request(BASE+'/saas-common/upload/uploadFiles',body,{'Authorization':'Bearer '+t,'Content-Type':'multipart/form-data; boundary='+boundary});return json.load(urllib.request.urlopen(r,timeout=90))
def diskpath(p):
 p=pathlib.Path(p);return p if p.is_absolute() else ROOT/'backend'/p

def xlsx():
 b=io.BytesIO()
 with zipfile.ZipFile(b,'w') as z:
  z.writestr('[Content_Types].xml','''<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types"><Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/><Default Extension="xml" ContentType="application/xml"/><Override PartName="/xl/workbook.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.sheet.main+xml"/><Override PartName="/xl/worksheets/sheet1.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.worksheet+xml"/></Types>''')
  z.writestr('_rels/.rels','''<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="xl/workbook.xml"/></Relationships>''')
  z.writestr('xl/workbook.xml','''<workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"><sheets><sheet name="fixture" sheetId="1" r:id="rId1"/></sheets></workbook>''')
  z.writestr('xl/_rels/workbook.xml.rels','''<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/worksheet" Target="worksheets/sheet1.xml"/></Relationships>''')
  z.writestr('xl/worksheets/sheet1.xml','''<worksheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><sheetData><row r="1"><c r="A1" t="inlineStr"><is><t>phase6 media fixture</t></is></c></row></sheetData></worksheet>''')
 return b.getvalue()

ids=[];files=[];dirs=[];tok=login()
try:
 png=base64.b64decode('iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8/x8AAusB9Y9Zl1sAAAAASUVORK5CYII=')
 r=upload('phase6-99977782.png',png,tok);assert r['code']==200 and len(r['data'])==1,r;f=r['data'][0];rid=int(f['id']);cid=int(f['compressedId']);ids.extend([rid,cid]);row=sql(f'SELECT compressed_id FROM chain_common_file WHERE id={rid}');assert int(row)==cid and cid>0;row=sql(f'SELECT file_path,file_name,is_compressed_file,effective,source_id FROM chain_common_file WHERE id={cid}').split('\t');assert row[2:]==['1','1',str(rid)],row;cp=diskpath(row[0])/row[1];files.append(cp);assert cp.read_bytes()==png;op=diskpath(row[0])/f['fileName'];files.append(op);assert op.read_bytes()==png;print('PASS image compressed row + source FK + exact small-image bytes')
 r=upload('phase6-99977783.xlsx',xlsx(),tok);assert r['code']==200 and len(r['data'])==1,r;f=r['data'][0];rid=int(f['id']);ids.append(rid);row=sql(f'SELECT file_path,file_name FROM chain_common_file WHERE id={rid}').split('\t');files.append(diskpath(row[0])/row[1]);rows=sql(f'SELECT id,file_path,file_name FROM chain_common_img WHERE file_id={rid}').splitlines();assert len(rows)==1,rows
 for row in rows:
  iid,p,n=row.split('\t');p=diskpath(p);dirs.append(p);assert (p/n).read_bytes().startswith(b'\x89PNG');sql(f'DELETE FROM chain_common_img WHERE id={iid}')
 print('PASS spreadsheet source row + one-sheet PNG preview row + PNG bytes')
 # Current local Java Doc2PdfUtil throws before adding response: file saved, no preview, [] result.
 before=set(sql("SELECT id FROM chain_common_file WHERE file_real_name='phase6-99977784.doc'").split())
 r=upload('phase6-99977784.doc',b'phase6 doc fixture',tok);assert r['code']==200 and r['data']==[],r;after=set(sql("SELECT id FROM chain_common_file WHERE file_real_name='phase6-99977784.doc'").split());new=after-before;assert len(new)==1,new
 rid=int(new.pop());ids.append(rid);row=sql(f'SELECT file_path,file_name FROM chain_common_file WHERE id={rid}').split('\t');files.append(diskpath(row[0])/row[1]);assert sql(f'SELECT count(*) FROM chain_common_img WHERE file_id={rid}')=='0';print('PASS doc converter Java stub parity: saved file + [] response + no preview')
finally:
 if ids:sql('DELETE FROM chain_common_img WHERE file_id IN ('+','.join(map(str,ids))+')');sql('DELETE FROM chain_common_file WHERE id IN ('+','.join(map(str,ids))+')')
 for p in files:p.unlink(missing_ok=True)
 for p in dirs:shutil.rmtree(p,ignore_errors=True)
print('3/3 media write fixture groups passed; DB/files cleaned')
