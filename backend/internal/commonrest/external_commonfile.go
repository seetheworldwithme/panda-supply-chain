package commonrest

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

func commonFileBody(w http.ResponseWriter, r *http.Request) (map[string]any, bool) {
	b, e := io.ReadAll(r.Body)
	var m map[string]any
	if e != nil || json.Unmarshal(b, &m) != nil || m == nil {
		platform.Unknown(w)
		return nil, false
	}
	m = normalizeFileBody(m)
	ent := entityForTable("chain_common_file")
	if !ent.validateBody(m) {
		platform.Unknown(w)
		return nil, false
	}
	return m, true
}
func registerCommonFile(mux *http.ServeMux, s *Service) {
	reg := func(p string, h http.HandlerFunc) {
		mux.HandleFunc("POST /saas-common/chainCommonFile/"+p, authGuard(h))
	}
	reg("deleteById", s.commonFileDelete)
	reg("updateById", s.commonFileUpdate)
	reg("delBatch", s.commonFileDelBatch)
	reg("updBthById", s.commonFileUpdBatch)
	reg("delByRelId", s.commonFileDelRel)
	reg("compressedFile", func(w http.ResponseWriter, r *http.Request) {
		m, ok := commonFileBody(w, r)
		if !ok {
			return
		}
		platform.Success(w, s.compressFile(r, int64(num(m["id"]))))
	})
}

// Java removes the DB row before getById: original file and *_img directory survive.
// Converted preview files are deleted on disk, but chain_common_img rows remain.
func (s *Service) removeCommonFile(r *http.Request, id any) error {
	d := s.pick(r)
	if e := platform.DeleteByID(d, "chain_common_file", id); e != nil {
		return e
	}
	if id == nil {
		_, e := d.Exec("UPDATE chain_basic_enterprise SET advertising_video_id=NULL WHERE advertising_video_id=")
		return e
	}
	if _, e := d.Exec("UPDATE chain_basic_enterprise SET advertising_video_id=NULL WHERE advertising_video_id=?", id); e != nil {
		return e
	}
	if _, e := d.Exec("UPDATE chain_basic_enterprise SET website_img_id=NULL,website_ced_img_id=NULL WHERE website_img_id=? OR website_ced_img_id=?", id, id); e != nil {
		return e
	}
	rows, e := queryRows(d, "SELECT file_path,file_name FROM chain_common_img WHERE file_id=?", id)
	if e != nil {
		return e
	}
	for _, row := range rows {
		os.Remove(filepath.Join(str(row, "file_path"), str(row, "file_name")))
	}
	_, e = d.Exec("DELETE FROM chain_common_file WHERE source_id=?", id)
	return e
}
func (s *Service) commonFileDelete(w http.ResponseWriter, r *http.Request) {
	m, ok := commonFileBody(w, r)
	if !ok {
		return
	}
	if e := s.removeCommonFile(r, m["id"]); e != nil {
		platform.Unknown(w)
		return
	}
	id := "null"
	if m["id"] != nil {
		id = javaStr(m["id"])
	}
	platform.Success(w, "删除成功，id为："+id)
}
func fileWriteVals(e *ent, m map[string]any) map[string]any {
	vals := pick(m, e.dbcols)
	delete(vals, "id")
	for k, v := range vals {
		if v == nil {
			delete(vals, k)
		}
	}
	return vals
}
func (s *Service) commonFileUpdate(w http.ResponseWriter, r *http.Request) {
	m, ok := commonFileBody(w, r)
	if !ok {
		return
	}
	e := entityForTable("chain_common_file")
	vals := fileWriteVals(e, m)
	d := s.pick(r)
	if m["fileName"] != nil {
		ks := platform.SortedKeys(vals)
		sets := []string{}
		args := []any{}
		for _, k := range ks {
			sets = append(sets, k+"=?")
			args = append(args, vals[k])
		}
		args = append(args, m["fileName"])
		if _, err := d.Exec("UPDATE chain_common_file SET "+strings.Join(sets, ",")+" WHERE file_name=?", args...); err != nil {
			platform.Unknown(w)
			return
		}
	} else {
		if _, err := platform.UpdateByIDNoFill(d, "chain_common_file", m["id"], vals); err != nil {
			platform.Unknown(w)
			return
		}
		q := "SELECT * FROM chain_common_file WHERE source_id=?"
		args := []any{m["id"]}
		if m["id"] == nil {
			q = "SELECT * FROM chain_common_file WHERE source_id IS NULL"
			args = nil
		}
		rows, err := queryRows(d, q, args...)
		if err != nil || len(rows) > 1 {
			platform.Unknown(w)
			return
		}
		if len(rows) == 1 {
			f := rows[0]
			v := map[string]any{}
			for _, c := range e.dbcols {
				if c != "id" && f[c] != nil {
					v[c] = f[c]
				}
			}
			delete(v, "association_id")
			delete(v, "f_type")
			if m["associationId"] != nil {
				v["association_id"] = m["associationId"]
			}
			ft := m["ftype"]
			if ft == nil {
				ft = m["fType"]
			}
			if ft != nil {
				v["f_type"] = ft
			}
			if _, err = platform.UpdateByIDNoFill(d, "chain_common_file", f["id"], v); err != nil {
				platform.Unknown(w)
				return
			}
		}
	}
	platform.Success(w, e.echoBody(m))
}
func (s *Service) commonFileDelBatch(w http.ResponseWriter, r *http.Request) {
	m, ok := commonFileBody(w, r)
	if !ok {
		return
	}
	ids, ok := m["ids"].([]any)
	if !ok || len(ids) == 0 {
		platform.Unknown(w)
		return
	}
	ph := strings.TrimSuffix(strings.Repeat("?,", len(ids)), ",")
	rows, e := queryRows(s.pick(r), "SELECT id FROM chain_common_file WHERE id IN ("+ph+")", ids...)
	if e != nil {
		platform.Unknown(w)
		return
	}
	for _, row := range rows {
		if e := s.removeCommonFile(r, row["id"]); e != nil {
			platform.Unknown(w)
			return
		}
	}
	platform.Success(w, "删除成功")
}
func (s *Service) commonFileDelRel(w http.ResponseWriter, r *http.Request) {
	m, ok := commonFileBody(w, r)
	if !ok {
		return
	}
	cond := "1=1"
	args := []any{}
	for _, k := range platform.SortedKeys(m) {
		if !externalColumn(k) {
			platform.Unknown(w)
			return
		}
		if m[k] == nil {
			cond += " AND " + k + " IS NULL"
		} else {
			cond += " AND " + k + "=?"
			args = append(args, m[k])
		}
	}
	rows, e := queryRows(s.pick(r), "SELECT id FROM chain_common_file WHERE "+cond, args...)
	if e != nil {
		platform.Unknown(w)
		return
	}
	for _, row := range rows {
		if e := s.removeCommonFile(r, row["id"]); e != nil {
			platform.Unknown(w)
			return
		}
	}
	platform.Success(w, "删除成功")
}
func (s *Service) commonFileUpdBatch(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	var ms []map[string]any
	if err != nil || json.Unmarshal(b, &ms) != nil || bytes.Equal(bytes.TrimSpace(b), []byte("null")) {
		platform.Unknown(w)
		return
	}
	e := entityForTable("chain_common_file")
	d := s.pick(r)
	tx, err := d.Begin()
	if err != nil {
		platform.Unknown(w)
		return
	}
	defer tx.Rollback()
	for _, m := range ms {
		m = normalizeFileBody(m)
		if !e.validateBody(m) {
			platform.Unknown(w)
			return
		}
		vals := fileWriteVals(e, m)
		ks := platform.SortedKeys(vals)
		sets := []string{}
		args := []any{}
		for _, k := range ks {
			sets = append(sets, k+"=?")
			args = append(args, vals[k])
		}
		args = append(args, m["id"])
		if _, err = tx.Exec("UPDATE chain_common_file SET "+strings.Join(sets, ",")+" WHERE id=?", args...); err != nil {
			platform.Unknown(w)
			return
		}
	}
	if err = tx.Commit(); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, "批量修改成功")
}

func normalizeFileBody(m map[string]any) map[string]any {
	out := map[string]any{}
	for k, v := range m {
		if k != "fType" && k != "ftype" {
			out[k] = v
		}
	}
	if v, ok := m["ftype"]; ok {
		out["fType"] = v
	}
	return out
}
