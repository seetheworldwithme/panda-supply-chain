"""Local zjjg fixtures for multi-sheet import; original chain rows stay untouched.

A temporary local zjjg account lets an existing enterprise JWT exercise its
normal seven-character enterprise code. The account is never used for login.
Only marker-owned contracts, suppliers, logs and this fixture account are removed.
"""
import argparse
from datetime import datetime, timedelta
import io
import json
import re
import subprocess
import zipfile
from pathlib import Path
from contract_compare import decode, differences, Number
from verify_contract_import import excel, upload


def sql(query):
    return subprocess.check_output(['mysql','--protocol=SOCKET','--socket=/tmp/mysql.sock','-uroot','-N','-B','--raw','zjjg','-e',query],text=True).strip()


def multi_excel(sheets):
    out=io.BytesIO()
    with zipfile.ZipFile(out,'w') as z:
        ct='<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types"><Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/><Default Extension="xml" ContentType="application/xml"/><Override PartName="/xl/workbook.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.sheet.main+xml"/>'
        wb='<workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"><sheets>'
        rel='<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">'
        for i,(fields,values) in enumerate(sheets,1):
            blob=excel(fields,values or {},omit_missing=True)
            with zipfile.ZipFile(io.BytesIO(blob)) as template:
                content=template.read('xl/worksheets/sheet1.xml')
                if values is None:
                    # Header-only sheets contain no business records.
                    content=re.sub(rb'<row r="2">.*?</row>',b'',content)
                z.writestr(f'xl/worksheets/sheet{i}.xml',content)
            ct+=f'<Override PartName="/xl/worksheets/sheet{i}.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.worksheet+xml"/>'
            wb+=f'<sheet name="Sheet{i}" sheetId="{i}" r:id="rId{i}"/>'
            rel+=f'<Relationship Id="rId{i}" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/worksheet" Target="worksheets/sheet{i}.xml"/>'
        z.writestr('[Content_Types].xml',ct+'</Types>')
        z.writestr('_rels/.rels','<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="xl/workbook.xml"/></Relationships>')
        z.writestr('xl/workbook.xml',wb+'</sheets></workbook>')
        z.writestr('xl/_rels/workbook.xml.rels',rel+'</Relationships>')
    return out.getvalue()


def norm(value,row_id):
    if isinstance(value,dict):
        out={}
        for k,v in value.items():
            if k=='id' and isinstance(v,Number) and v==row_id:out[k]=Number('0')
            elif k=='supplierId' and isinstance(v,Number):out[k]=Number('0')
            elif k=='imBatchNumber' and isinstance(v,str) and re.fullmatch(r'[0-9a-f-]{36}',v):out[k]='{batch}'
            else:out[k]=norm(v,row_id)
        return out
    if isinstance(value,list):return [norm(v,row_id) for v in value]
    return value


def main():
    p=argparse.ArgumentParser(description=__doc__)
    p.add_argument('--sessions',type=Path,required=True)
    p.add_argument('--go-port',default='8798')
    p.add_argument('--output',type=Path,required=True)
    p.add_argument('--populated',action='store_true',help='每张业务 sheet 填入隔离记录并核对落库')
    args=p.parse_args()
    sessions=json.loads(args.sessions.read_text())['hyqcbjhd']
    source=(Path(__file__).resolve().parents[2]/'backend/internal/commonrest/auth_meta_gen.go').read_text()
    def fields(name):
        block=re.search(r'"'+name+r'":\s*\{(.*?)(?=\n\t"Chain|\n\})',source,re.S)[1]
        return re.findall(r'\{"([^"]*)", "([^"]*)", "([^"]*)", (true|false), (true|false), (true|false)(?:, (?:true|false), (?:true|false))?\}',block)
    names=['ChainContract','ChainLogistics','ChainLogisticsWeighingList','ChainCapitalSettlement','ChainCapitalDetailed','ChainBill','ChainBill','ChainBillRate']
    meta={n:fields(n) for n in names}
    assert all(meta.values())
    user_marker='audit1118-jwt-fixture'
    assert sql("SELECT count(*) FROM chain_auth_user WHERE account='hyqcbjhd' OR remark5='"+user_marker+"'")=='0','fixture account collision'
    results=[]
    tables=['chain_capital_settlement','chain_capital_detailed','chain_bill','chain_bill_rate','chain_logistics','chain_logistics_weighing_list']
    def cleanup(where,sup_where,marker):
        batches=sql('SELECT im_batch_number FROM chain_contract WHERE '+where).splitlines()
        for table in tables:
            sql("DELETE FROM "+table+" WHERE remark5='"+marker+"'")
        for batch in batches:
            if re.fullmatch(r'[0-9a-f-]{36}',batch):sql("DELETE FROM chain_common_log WHERE operation LIKE '%"+batch+"%'")
        sql('DELETE FROM chain_contract WHERE '+where);sql('DELETE FROM chain_contract_supplier WHERE '+sup_where)
    try:
        sql("INSERT INTO chain_auth_user(account,username,level,account_type,handle_result,password,remark5) VALUES ('hyqcbjhd','本地兼容验收',1,3,'1','fixture-not-a-login', '"+user_marker+"')")
        paths=['chainContract/importOnce','chainContract/importAlone','chainCgContract/importAlone','chainXsContract/importAlone','chainContract/importStep']
        for ix,path in enumerate(paths):
            code='AUDIT1118-MULTI'+str(ix);supplier='本地多表验收供应商'+str(ix);buyer='本地多表验收客户'+str(ix)
            where="contract_code='"+code+"'";sup_where="enterprise_name IN ('"+supplier+"','"+buyer+"')"
            assert sql('SELECT count(*) FROM chain_contract WHERE '+where)=='0'
            assert sql('SELECT count(*) FROM chain_contract_supplier WHERE '+sup_where)=='0'
            values={}
            for _,col,kind,_,dictionary,required in meta['ChainContract']:
                if required=='true' or dictionary=='true':values[col]={'i':'1','d':'10.00','D':'2020-01-02'}.get(kind,'本地验收')
            values.update(contract_code=code,contract_type='1',sales_unit=supplier,buyer_unit=buyer,month='2020-01',create_time='2020-01-02 03:04:05',update_time='2020-01-02 03:04:05')
            chosen=names if not path.endswith('/importStep') else ['ChainContract','ChainLogistics','ChainCapitalSettlement','ChainCapitalDetailed','ChainBill','ChainBill','ChainBillRate']
            marker='audit1118-multi-'+str(ix)
            for table in tables:assert sql("SELECT count(*) FROM "+table+" WHERE remark5='"+marker+"'")=='0','fixture marker collision'
            sheets=[]
            for i,n in enumerate(chosen):
                row=values if i==0 else None
                if args.populated and i>0:
                    row={col:{'i':'1','d':'10.00','D':'2020-01-02','t':'2020-01-02 03:04:05'}.get(kind,'本地验收') for _,col,kind,_,dictionary,required in meta[n] if required=='true' or dictionary=='true'}
                    row.update(contract_code=code,bill_number=code+'-'+str(i),bill_code='AUDIT1118-'+str(i),month='2099-'+str(ix+1).zfill(2),remark5=marker,create_time='2020-01-02 03:04:05',update_time='2020-01-02 03:04:05',type='1',detype='1',bill_type='1')
                sheets.append((meta[n],row))
            blob=multi_excel(sheets)
            runs=[];dbs=[];codes=[];statuses=[]
            try:
                for port in ('9002',args.go_port):
                    started=datetime.now();before_log=int(sql('SELECT IFNULL(max(id),0) FROM chain_common_log'))
                    status,_,raw=upload(port,sessions['8765' if port=='9002' else port]['token'],'/saas-contract/'+path,blob,'zjjg')
                    ended=datetime.now()
                    new_logs=sql("SELECT id FROM chain_common_log WHERE id>"+str(before_log)+" AND user_id='7625'").splitlines()
                    result=decode(raw)
                    if args.populated and isinstance(result.get('data'),dict) and result['data'].get('errInfos'): print(port,'隔离数据校验错误',result['data']['errInfos'],flush=True)
                    row_id=sql('SELECT id FROM chain_contract WHERE '+where)
                    runs.append(norm(result,row_id));codes.append(result.get('code'));statuses.append(status)
                    dbs.append((sql('SELECT contract_code,contract_type,company_id,create_user_id,update_user_id,sum_price,contract_num,cust_type,system_code,create_time,update_time FROM chain_contract WHERE '+where),sql('SELECT s.enterprise_name,s.partner_type,s.cust_type FROM chain_contract c JOIN chain_contract_supplier s ON s.id=c.supplier_id WHERE '+where)))
                    supplier_id=sql('SELECT supplier_id FROM chain_contract WHERE '+where)
                    if args.populated:
                        snapshot=[]
                        for table in tables:
                            cols=sql('SHOW COLUMNS FROM '+table).splitlines();cols=[x.split('\t')[0] for x in cols]
                            rows=sql("SELECT "+','.join('`'+c+'`' for c in cols)+" FROM "+table+" WHERE remark5='"+marker+"' ORDER BY id").splitlines()
                            data=[]
                            for line in rows:
                                record=dict(zip(cols,line.split('\t')));record.pop('id',None);record.pop('im_batch_number',None)
                                if record.get('contract_id')==row_id:record['contract_id']='{contract-id}'
                                if table=='chain_capital_settlement' and record.get('supplier_id') not in (None,'NULL'):
                                    assert record['supplier_id']==supplier_id,'结算单供应商外键必须指向导入合同的供应商'
                                    record['supplier_id']='{supplier-id}'
                                if table=='chain_logistics' and record.get('create_time') not in (None,'NULL'):
                                    stamp=datetime.fromisoformat(record['create_time'])
                                    assert started-timedelta(seconds=2)<=stamp<=ended+timedelta(seconds=2),'自动填充时间必须落在实际请求窗口内'
                                    record['create_time']='{server-time-within-request}'
                                data.append(record)
                            snapshot.append([table,data])
                        dbs[-1]=(dbs[-1],snapshot)
                    cleanup(where,sup_where,marker)
                    if new_logs:
                        assert all(x.isdigit() for x in new_logs)
                        sql("DELETE FROM chain_common_log WHERE user_id='7625' AND id IN ("+','.join(new_logs)+")")
                diff=differences(runs[0],runs[1])
                if statuses[0]!=statuses[1]:diff.append('HTTP:status')
                if dbs[0]!=dbs[1]:
                    diff.extend(differences(json.loads(json.dumps(dbs[0])),json.loads(json.dumps(dbs[1])),'$database'))
                    if args.populated:
                        for (table,left),(_,right) in zip(dbs[0][1],dbs[1][1]):
                            if len(left)==len(right):
                                for lrow,rrow in zip(left,right):
                                    for k in lrow.keys()|rrow.keys():
                                        if lrow.get(k)!=rrow.get(k): print('fixture',table,k,repr(lrow.get(k)),repr(rrow.get(k)),flush=True)
                    else: print('fixture contract differs',flush=True)
                if not dbs[0] or not dbs[1] or codes!=[Number('200'),Number('200')]:diff.append('$business:import-failed')
                if args.populated:
                    for run in runs:
                        if isinstance(run.get('data'),dict) and run['data'].get('errs'):diff.append('$business:sheet-validation-failed')
                    for snap in dbs:
                        for table,data in snap[1]:
                            if not data and not (path.endswith('/importStep') and table in ('chain_logistics','chain_logistics_weighing_list')):diff.append('$database:'+table+':missing-row')
                results.append({'populated':args.populated,'endpoint':path,'differences':diff,'java_status':statuses[0],'go_status':statuses[1],'cleaned':True})
                args.output.write_text(json.dumps(results,ensure_ascii=False,indent=2)+'\n')
                print(path,'PASS' if not diff else 'FAIL '+','.join(diff),flush=True)
            finally:
                cleanup(where,sup_where,marker)
    finally:
        sql("DELETE FROM chain_auth_user WHERE account='hyqcbjhd' AND remark5='"+user_marker+"'")
    args.output.write_text(json.dumps(results,ensure_ascii=False,indent=2)+'\n')
    return int(any(x['differences'] for x in results))


if __name__=='__main__':raise SystemExit(main())
