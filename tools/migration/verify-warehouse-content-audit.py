"""Positive DB assertions for local services; payment is stopped before vendor call by NULL openid."""
from pathlib import Path
exec(Path(__file__).with_name('verify-warehouse-information.py').read_text().split('def main():')[0])
marker=99977780
try:
 assert sql(f'SELECT COUNT(*) FROM bus_member_config WHERE id={marker}')=='0'
 assert sql(f'SELECT COUNT(*) FROM bus_user_coupon WHERE id={marker}')=='0'
 assert sql(f'SELECT openid IS NULL FROM chain_auth_user WHERE id=7625')=='1','refuse payment test with non-NULL openid'
 sql(f'INSERT INTO bus_member_config(id,member_type,advance_deposit,final_year_price) VALUES({marker},{marker},15.00,100.00); INSERT INTO bus_user_coupon(id,coupon_money) VALUES({marker},5.00)')
 for trans,renew,expected in [(1,None,10),(2,False,95)]:
  snapshots=[]
  for port in [8765,8799]:
   sql(f'DELETE FROM bus_my_bill WHERE user_coupon_id={marker}')
   body={'memberType':marker,'transType':trans,'userCouponId':marker}
   if renew is not None:body['isRenew']=renew
   code,b=req(port,'/busMyBill/getPayInfo',body)
   assert code==500 and norm(b)['code']==999,(port,b)
   assert sql(f'SELECT COUNT(*) FROM bus_my_bill WHERE user_coupon_id={marker}')=='0',(port,'payment exception must roll back bill')
 sql(f'INSERT INTO bus_product(id,notice_time) VALUES({marker},NULL)')
 for port in [8765,8799]:
  code,b=req(port,'/busCustManage/noticeCust',{'proId':marker,'phones':[]})
  assert code==500 and norm(b)['code']==999,(port,'NULL noticeTime must fail before write',b)
  assert sql(f'SELECT notice_time IS NULL FROM bus_product WHERE id={marker}')=='1'
 print('PASS payment failure rollback and NULL notice counter; fixtures cleaned')
finally:
 sql(f'DELETE FROM bus_member_config WHERE id={marker}; DELETE FROM bus_user_coupon WHERE id={marker}; DELETE FROM bus_my_bill WHERE user_coupon_id={marker}; DELETE FROM bus_product WHERE id={marker}')
