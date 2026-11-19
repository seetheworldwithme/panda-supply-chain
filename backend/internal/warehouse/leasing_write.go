package warehouse

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

func init() {
	addExtra(func(mux *http.ServeMux, s *Service) {
		mux.HandleFunc("POST /saas-warehouse/busEnterpriseLeasing/add", authGuard(s.leasingAdd))
		mux.HandleFunc("POST /saas-warehouse/busEnterpriseLeasing/updateById", authGuard(s.leasingUpdate))
	})
}

func (s *Service) leasingAdd(w http.ResponseWriter, r *http.Request) {
	s.leasingWrite(w, r, false)
}

func (s *Service) leasingUpdate(w http.ResponseWriter, r *http.Request) {
	s.leasingWrite(w, r, true)
}

func (s *Service) leasingWrite(w http.ResponseWriter, r *http.Request, update bool) {
	m := bodyMap(r)
	if !entBusEnterpriseLeasing.validateBody(m) || !validLeasingFiles(m) {
		platform.Unknown(w)
		return
	}
	d := s.pick(r)
	var previous map[string]any
	var err error
	if update {
		previous, err = platform.GetByID(d, "bus_enterprise_leasing", m["id"])
		if err != nil {
			platform.Unknown(w)
			return
		}
	}
	// 原服务没有总事务；证书批次、照片批次分别先提交，再保存主记录。
	if err = s.linkLeasingFiles(r, m); err != nil {
		platform.Unknown(w)
		return
	}
	values := pick(m, entBusEnterpriseLeasing.dbcols)
	delete(values, "id")
	if update {
		if _, err = platform.UpdateByID(d, "bus_enterprise_leasing", m["id"], values); err != nil {
			platform.Unknown(w)
			return
		}
		// 原代码在更新之后拆箱旧状态；缺行/旧状态空的异常不回滚附件和主记录。
		if previous == nil || previous["status"] == nil || (num(previous["status"]) == 0 && m["status"] == nil) {
			platform.Unknown(w)
			return
		}
		if num(previous["status"]) == 0 && num(m["status"]) == 1 {
			// 原 sendWxPubMsg 仅发三段，现役消费者却按煤种四段消息解析，先越界再失败。
			// 保留此问题，不改为成功发送租售公众号通知。
			log.Printf("企业租售后续同步失败: 原 sendWxPubMsg 消息只有三段，syncOtherInfo 读取第四段越界")
		}
	} else {
		uid := currentUser(r).UserID
		user, err := platform.GetByID(d, "chain_auth_user", uid)
		if err != nil || user == nil {
			platform.Unknown(w)
			return
		}
		values["create_user_id"], values["create_phone"] = uid, user["phone"]
		if _, err = platform.Insert(d, "bus_enterprise_leasing", values); err != nil {
			platform.Unknown(w)
			return
		}
	}
	platform.Success(w, nil)
}

func validLeasingFiles(m map[string]any) bool {
	for _, key := range []string{"cerList", "photoList"} {
		if m[key] == nil {
			continue
		}
		items, ok := m[key].([]any)
		if !ok {
			return false
		}
		for _, item := range items {
			if item == nil {
				continue // Jackson 可绑定 null 元素，服务访问该元素时才抛异常。
			}
			file, ok := item.(map[string]any)
			if !ok || !fileEnt.validateBody(file) {
				return false
			}
		}
	}
	return true
}

func (s *Service) linkLeasingFiles(r *http.Request, m map[string]any) error {
	for _, fields := range [][2]string{{"cerList", "certAssociationId"}, {"photoList", "photoAssociationId"}} {
		if m[fields[0]] == nil {
			continue
		}
		items := m[fields[0]].([]any)
		association := uuid32()
		m[fields[1]] = association
		// 先设置整个列表再调用原 Feign 批更；null 元素不能留下本批次的局部更新。
		files := make([]map[string]any, 0, len(items))
		for _, item := range items {
			file, ok := item.(map[string]any)
			if !ok || file == nil {
				return fmt.Errorf("Java 租售附件为空")
			}
			files = append(files, file)
		}
		if len(files) == 0 {
			continue
		}
		if err := platform.WithTx(s.pick(r), func(tx *sql.Tx) error {
			for _, file := range files {
				values := pick(file, fileEnt.dbcols)
				delete(values, "id")
				values["f_type"], values["owner_type"], values["owner_id"] = 4, "1", currentUser(r).UserID
				values["relation_id"], values["relation_type"] = m["id"], 1
				values["effective"], values["association_id"] = true, association
				if err := platform.TxUpdateByIDNoFill(tx, "chain_common_file", file["id"], values); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			return err
		}
	}
	return nil
}
