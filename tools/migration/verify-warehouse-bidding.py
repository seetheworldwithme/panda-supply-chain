"""Isolated local HTTP vendor fixture; requires a freshly built GO_AUDIT_BINARY."""
from pathlib import Path
exec(Path(__file__).with_name('verify-warehouse-information.py').read_text().split('def main():')[0])
import http.server,threading,time
marker='migration-warehouse-bidding-99977782'
state={'fail':False}
class Vendor(http.server.BaseHTTPRequestHandler):
 def do_POST(self):
  assert json.loads(self.rfile.read(int(self.headers['Content-Length'])))=={'page':1,'pageSize':50,'status':2,'requireType':0,'tradeSources':'quote'}
  if state['fail']:self.send_response(503);self.end_headers();return
  self.send_response(200);self.send_header('Content-Type','application/json');self.end_headers()
  self.wfile.write(json.dumps({'data':[{'coalCode':marker,'preorderBeginTime':'2099-01-01 00:00:00','preorderEndTime':'2099-01-02 00:00:00','price':123.45,'buyNum':1}]}).encode())
 def log_message(self,*args):pass
vendor=http.server.HTTPServer(('127.0.0.1',18998),Vendor)
threading.Thread(target=vendor.serve_forever,daemon=True).start()
env=dict(os.environ,PORT='18898',BIDDING_SOURCE_URL='http://127.0.0.1:18998/vendor')
cwd=Path(__file__).resolve().parents[2]/'backend'
proc=subprocess.Popen([os.environ['GO_AUDIT_BINARY']],cwd=cwd,env=env,stdout=subprocess.DEVNULL,stderr=subprocess.DEVNULL)
def call():
 request=urllib.request.Request('http://127.0.0.1:18898/saas-warehouse/busBiddingInfo/getBiddingInfo',data=b'{}',headers={'Authorization':'Bearer '+TOK[8799],'Content-Type':'application/json'})
 try:
  r=urllib.request.urlopen(request,timeout=10);return r.status,json.load(r)
 except urllib.error.HTTPError as r:return r.code,json.load(r)
try:
 assert sql('SELECT COUNT(*) FROM bus_bidding_info WHERE id=99977781')=='0'
 assert sql(f'SELECT COUNT(*) FROM bus_bidding_info WHERE coal_code="{marker}"')=='0'
 sql('INSERT INTO bus_bidding_info(id,preorder_begin_time,coal_code) VALUES(99977781,"1900-01-01 00:00:00","migration-bidding-seed")')
 for _ in range(100):
  try:
   urllib.request.urlopen('http://127.0.0.1:18898/healthz',timeout=0.2);break
  except Exception:time.sleep(0.1)
 else:raise AssertionError('Go audit service did not start')
 for _ in range(2):
  status,j=call();assert status==200 and j['code']==200,(status,j)
  assert sql(f'SELECT COUNT(*) FROM bus_bidding_info WHERE coal_code="{marker}"')=='1'
 assert sql(f'SELECT price,buy_num FROM bus_bidding_info WHERE coal_code="{marker}"')=='123.45\t1.00'
 state['fail']=True
 status,j=call();assert status==500 and j['code']==999,(status,j)
 print('PASS bidding vendor payload, HTTP fetch, DB insert, strict cutoff idempotence, vendor error')
finally:
 sql(f'DELETE FROM bus_bidding_info WHERE coal_code="{marker}"; DELETE FROM bus_bidding_info WHERE id=99977781 AND coal_code="migration-bidding-seed"')
 proc.terminate();proc.wait(timeout=10);vendor.shutdown()
