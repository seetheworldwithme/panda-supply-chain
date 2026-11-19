"""Fixture-only contract Excel import comparison against the local Java JAR.

Creates an XLSX in memory using generated Java annotation metadata. Validates
the response, contract row and supplier relation, then cleans only its markers.
"""
import argparse
import html
import io
import json
import re
import subprocess
import urllib.request
import urllib.error
import zipfile
from pathlib import Path
from contract_compare import compare
SCHEMA='chain'


def sql(query):
    return subprocess.check_output(["mysql", "--protocol=SOCKET", "--socket=/tmp/mysql.sock", "-uroot", "-N", "-B", "--raw", SCHEMA, "-e", query], text=True).strip()


def excel(fields, values, omit_missing=False):
    grid = [[f[0] for f in fields], [values.get(f[1], "") for f in fields]]
    sheet = ['<worksheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><sheetData>']
    for i, row in enumerate(grid, 1):
        sheet.append(f'<row r="{i}">')
        for j, value in enumerate(row):
            if omit_missing and i==2 and fields[j][1] not in values:continue
            n, col = j + 1, ""
            while n:
                n, k = divmod(n - 1, 26)
                col = chr(65 + k) + col
            sheet.append(f'<c r="{col}{i}" t="inlineStr"><is><t>{html.escape(str(value))}</t></is></c>')
        sheet.append('</row>')
    sheet.append('</sheetData></worksheet>')
    blob = io.BytesIO()
    with zipfile.ZipFile(blob, "w") as z:
        z.writestr('[Content_Types].xml', '<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types"><Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/><Default Extension="xml" ContentType="application/xml"/><Override PartName="/xl/workbook.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.sheet.main+xml"/><Override PartName="/xl/worksheets/sheet1.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.worksheet+xml"/></Types>')
        z.writestr('_rels/.rels', '<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="xl/workbook.xml"/></Relationships>')
        z.writestr('xl/workbook.xml', '<workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"><sheets><sheet name="Sheet1" sheetId="1" r:id="rId1"/></sheets></workbook>')
        z.writestr('xl/_rels/workbook.xml.rels', '<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/worksheet" Target="worksheets/sheet1.xml"/></Relationships>')
        z.writestr('xl/worksheets/sheet1.xml', ''.join(sheet))
    return blob.getvalue()


def upload(port, token, path, blob, applytype=''):
    if port == '9002': path=path.removeprefix('/saas-contract')
    boundary = "panda-local-contract-audit"
    body = (f'--{boundary}\r\nContent-Disposition: form-data; name="file"; filename="fixture.xlsx"\r\nContent-Type: application/vnd.openxmlformats-officedocument.spreadsheetml.sheet\r\n\r\n').encode() + blob + f'\r\n--{boundary}--\r\n'.encode()
    headers={'Authorization':'Bearer '+token,'Content-Type':'multipart/form-data; boundary='+boundary}
    if applytype:headers['applytype']=applytype
    req = urllib.request.Request(f'http://127.0.0.1:{port}'+path, data=body, headers=headers)
    try:
        response = urllib.request.urlopen(req, timeout=30)
    except urllib.error.HTTPError as e:
        response = e
    with response:
        return response.status, dict(response.headers), response.read()


def main():
    p = argparse.ArgumentParser(description=__doc__)
    p.add_argument('--sessions', type=Path, required=True)
    p.add_argument('--go-port', default='8798')
    p.add_argument('--java-port', default='8765', choices=['8765','9002'])
    p.add_argument('--schema',default='chain',choices=['chain','zjjg'])
    p.add_argument('--output', type=Path, required=True)
    args = p.parse_args()
    global SCHEMA
    SCHEMA=args.schema
    sessions = json.loads(args.sessions.read_text())['ptadmin' if args.schema=='zjjg' else 'hyqcbjhd']
    source = (Path(__file__).resolve().parents[2] / 'backend/internal/commonrest/auth_meta_gen.go').read_text()
    block = re.search(r'"ChainContract":\s*\{(.*?)(?=\n\t"Chain|\n\})',source,re.S)[1]
    fields = re.findall(r'\{"([^"]*)", "([^"]*)", "([^"]*)", (true|false), (true|false), (true|false)(?:, (?:true|false), (?:true|false))?\}', block)
    assert fields, 'contract annotation metadata missing'
    results = []
    for index, endpoint in enumerate(('chainContract/importExcelCg', 'chainContract/importExcelXs', 'chainCgContract/importExcelCg', 'chainCgContract/importExcelXs', 'chainXsContract/importExcelCg', 'chainXsContract/importExcelXs')):
        code = 'AUDIT1118-C' + str(index)
        supplier = '本地兼容验收供应商' + str(index)
        buyer = '本地兼容验收客户' + str(index)
        where = "contract_code='" + code + "'"
        suppliers = "enterprise_name IN ('" + supplier + "','" + buyer + "')"
        assert sql('SELECT count(*) FROM chain_contract WHERE '+where) == '0'
        assert sql('SELECT count(*) FROM chain_contract_supplier WHERE '+suppliers) == '0'
        values = {}
        for _, col, kind, _, dictionary, required in fields:
            if required == 'true' or dictionary == 'true':
                values[col] = {'i': '1', 'd': '10.00', 'D': '2020-01-02'}.get(kind, '本地验收')
        values.update(contract_code=code, contract_type='2' if endpoint.endswith('Xs') else '1', sales_unit=supplier, buyer_unit=buyer, month='2020-01', create_time='2020-01-02 03:04:05', update_time='2020-01-02 03:04:05')
        responses, rows = [], []
        try:
            for port in (args.java_port, args.go_port):
                responses.append(upload(port, sessions['8765' if port==args.java_port else port]['token'], '/saas-contract/'+endpoint, excel(fields, values),'zjjg' if args.schema=='zjjg' else ''))
                row = sql('SELECT contract_code,contract_type,contract_num,unit_price,sum_price,cust_type,create_user_id,update_user_id,system_code,create_time,update_time FROM chain_contract WHERE '+where)
                relation = sql('SELECT s.enterprise_name,s.partner_type,s.cust_type FROM chain_contract c JOIN chain_contract_supplier s ON s.id=c.supplier_id WHERE '+where)
                rows.append((row, relation))
                sql('DELETE FROM chain_contract WHERE '+where)
                sql('DELETE FROM chain_contract_supplier WHERE '+suppliers)
            diff = compare(*responses)
            if rows[0] != rows[1]:
                diff.append('$database:contract-or-supplier')
                columns=['contract_code','contract_type','contract_num','unit_price','sum_price','cust_type','create_user_id','update_user_id','system_code','create_time','update_time']
                print('different contract DB columns:',[k for k,a,b in zip(columns,rows[0][0].split('\t'),rows[1][0].split('\t')) if a!=b],flush=True)
            # Both errors do not count as successful business implementation.
            for response in responses:
                if response[0] != 200 or json.loads(response[2]).get('code') != 200:
                    diff.append('$business:import-failed')
                    payload=json.loads(response[2]); data=payload.get('data')
                    if isinstance(data,list) and all(isinstance(v,str) for v in data):
                        print('fixture validation errors:',data[:8],flush=True)
            results.append({'endpoint': endpoint, 'differences': diff, 'java_status':responses[0][0], 'go_status':responses[1][0], 'java_code':json.loads(responses[0][2]).get('code'), 'go_code':json.loads(responses[1][2]).get('code'), 'java_fixture':bool(rows[0][0]),'go_fixture':bool(rows[1][0]), 'fixture_created': bool(rows[0][0]) and bool(rows[1][0]), 'cleaned': True})
            print(endpoint, 'PASS' if not diff else 'FAIL '+','.join(diff), flush=True)
        finally:
            sql('DELETE FROM chain_contract WHERE '+where)
            sql('DELETE FROM chain_contract_supplier WHERE '+suppliers)
    args.output.write_text(json.dumps(results, ensure_ascii=False, indent=2)+'\n')
    return int(any(x['differences'] for x in results))


if __name__ == '__main__':
    raise SystemExit(main())
