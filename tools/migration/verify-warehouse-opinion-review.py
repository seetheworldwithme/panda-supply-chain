"""Local opinion review fixture. Go only on the award branch so no provider is called."""
from pathlib import Path

exec(Path(__file__).with_name('verify-warehouse-information.py').read_text().split('def main():')[0])

marker = 99977763
assert sql(f"SELECT COUNT(*) FROM bus_opinion WHERE id={marker}") == '0'
uid = int(sql("SELECT id FROM chain_auth_user WHERE account='hyqcbjhd'"))
previous = set(sql(f"SELECT id FROM bus_red_money_clog WHERE user_id={uid} AND receive_type=7").split())
created = set()
try:
    sql(f"INSERT INTO bus_opinion(id,user_id,status,content) VALUES ({marker},{uid},0,'phase6 review fixture')")
    for port in (8765, 8799):
        status, raw = req(port, '/busOpinion/updateById', {'id': marker, 'status': 2})
        assert status == 200 and norm(raw)['code'] == 200, (port, status, raw)
        assert sql(f"SELECT status FROM bus_opinion WHERE id={marker}") == '2'
        sql(f"UPDATE bus_opinion SET status=0 WHERE id={marker}")
    status, raw = req(8799, '/busOpinion/updateById', {'id': marker, 'status': 1})
    created = set(sql(f"SELECT id FROM bus_red_money_clog WHERE user_id={uid} AND receive_type=7").split()) - previous
    assert status == 200 and norm(raw)['code'] == 200, (status, raw)
    assert sql(f"SELECT status FROM bus_opinion WHERE id={marker}") == '1'
    assert len(created) == 2
    ids = ','.join(sorted(created))
    assert sql(f"SELECT COUNT(DISTINCT create_amount), SUM(is_write_off) FROM bus_red_money_clog WHERE id IN ({ids})") == '1\t0'
    print('PASS Java/Go non-award review and Go approval reward records')
finally:
    sql(f"DELETE FROM bus_opinion WHERE id={marker}")
    if created:
        sql("DELETE FROM bus_red_money_clog WHERE id IN (" + ','.join(sorted(created)) + ")")
