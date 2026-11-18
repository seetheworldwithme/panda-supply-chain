"""Verify only Go's positive video path so the Java SMS vendor is never contacted."""
from pathlib import Path

exec(Path(__file__).with_name('verify-warehouse-information.py').read_text().split('def main():')[0])

fixture = 99977762
label = 'phase6-video-99977762'
assert sql(f"SELECT count(*) FROM chain_common_file WHERE id={fixture}") == '0'
assert sql(f"SELECT count(*) FROM bus_videos WHERE video_id={fixture}") == '0'
assert sql(f"SELECT count(*) FROM bus_opinion WHERE content='{label}'") == '0'
uid = int(sql("SELECT id FROM chain_auth_user WHERE account='hyqcbjhd'"))
self_care = sql(f"SELECT count(*) FROM bus_user_care WHERE user_id={uid} AND care_user_id={uid}")
try:
    sql(f"INSERT INTO chain_common_file(id, file_name) VALUES ({fixture}, '{label}')")
    status, raw = req(8799, '/busVideos/add', {'videoId': fixture, 'videoTitle': label})
    assert status == 200 and norm(raw)['code'] == 200, (status, raw)
    assert sql(f"SELECT count(*) FROM bus_videos WHERE video_id={fixture} AND user_id={uid}") == '1'
    assert sql(f"SELECT count(*) FROM chain_common_file WHERE id={fixture} AND owner_id={uid} AND effective=1 AND association_id IS NOT NULL") == '1'
    assert int(sql(f"SELECT count(*) FROM bus_user_care WHERE user_id={uid} AND care_user_id={uid}")) == max(1, int(self_care))
    status, raw = req(8799, '/busOpinion/add', {'content': label})
    assert status == 200 and norm(raw)['code'] == 200, (status, raw)
    assert sql(f"SELECT count(*) FROM bus_opinion o JOIN chain_auth_user u ON o.user_id=u.id AND o.phone=u.phone WHERE o.content='{label}' AND u.id={uid}") == '1'
    print('PASS Go video, attachment, self-follow and feedback; fixture cleanup follows')
finally:
    sql(f"DELETE FROM bus_videos WHERE video_id={fixture}; DELETE FROM chain_common_file WHERE id={fixture}; DELETE FROM bus_opinion WHERE content='{label}' AND user_id={uid}")
    if self_care == '0':
        sql(f"DELETE FROM bus_user_care WHERE user_id={uid} AND care_user_id={uid}")
