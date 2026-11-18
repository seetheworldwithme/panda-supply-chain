import json,urllib.request,urllib.error,subprocess,os,urllib.parse

def sql(q):
 r=subprocess.run(['mysql','-h127.0.0.1','-uroot','chain','-N','-B','-e',q],capture_output=True,text=True,check=True);return r.stdout.strip()
def post(base,path,m,tok):
 req=urllib.request.Request(base+path,data=json.dumps(m).encode(),headers={'Authorization':'Bearer '+tok,'Content-Type':'application/json'})
 try:
  r=urllib.request.urlopen(req,timeout=20)
 except urllib.error.HTTPError as e:r=e
 return r.status,json.loads(r.read())
def token(base):
 req=urllib.request.Request(base+'/saas-auth/login',data=urllib.parse.urlencode({'username':os.environ['MIGRATION_USER'],'password':os.environ['MIGRATION_PASSWORD'],'code':'0000'}).encode());return json.load(urllib.request.urlopen(req))['data']['token']
results=[]
for base in ['http://localhost:8765','http://localhost:8799']:
 tok=token(base);out=[]
 for case in ['missing','null','normal']:
  sql('DELETE FROM bus_user_video_lc WHERE video_id=99977740');sql('DELETE FROM bus_videos WHERE id=99977740')
  if case!='missing':sql("INSERT INTO bus_videos(id,like_num,collect_num,forward_num,play_num) VALUES(99977740,"+('NULL' if case=='null' else '0')+",0,0,0)")
  res=post(base,'/saas-warehouse/busVideos/toLike',{'id':99977740},tok)
  state=sql('SELECT count(*) FROM bus_user_video_lc WHERE video_id=99977740')
  val=sql('SELECT like_num FROM bus_videos WHERE id=99977740')
  out.append((case,res,state,val))
 sql('DELETE FROM bus_user_video_lc WHERE video_id=99977740');sql('DELETE FROM bus_videos WHERE id=99977740')
 results.append(out)
print(json.dumps(results,ensure_ascii=False,indent=2));assert results[0]==results[1]
