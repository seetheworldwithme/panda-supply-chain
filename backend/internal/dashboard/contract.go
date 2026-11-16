// saas-contract 仪表盘接口(合同数据/产品交易份额)。
// Java 位置:saas-contract/.../ChainContractController → ChainContractServiceImpl
package dashboard

import (
	"net/http"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

// cttSelectQuantity POST /saas-contract/chainContract/selectQuantity 合同数量/金额按类型。
// Java ChainContractServiceImpl.selectQuat:431 + selectCount:751 → mapper selectQuat:57/selectCount:69
// (拦截器租户在本环境被 ignore-tables 跳过,全表)。
// 响应 {htllist:[{contract_type,sumNum}], fslist:[{contract_type,contractNum}]}
func (s *Service) cttSelectQuantity(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if m["timeType"] == nil {
		platform.Unknown(w) // Java `1 == timeType` Integer 拆箱 NPE(行为保持)
		return
	}
	d := s.pick(r)
	cond := timeTypeCond(m)
	var args []any

	htllist := []map[string]any{}
	if rows, err := d.Query("SELECT contract_type, sum(contract_num) sumNum FROM chain_contract WHERE 1=1"+cond+" GROUP BY contract_type", args...); err == nil {
		if ms, _ := scanMaps(rows); ms != nil {
			for _, row := range ms {
				htllist = append(htllist, map[string]any{"contract_type": num(row["contract_type"]), "sumNum": num(row["sumNum"])})
			}
		}
	}
	fslist := []map[string]any{}
	if rows, err := d.Query("SELECT contract_type, count(1) contractNum FROM chain_contract WHERE 1=1"+cond+" GROUP BY contract_type", args...); err == nil {
		if ms, _ := scanMaps(rows); ms != nil {
			for _, row := range ms {
				fslist = append(fslist, map[string]any{"contract_type": num(row["contract_type"]), "contractNum": num(row["contractNum"])})
			}
		}
	}
	platform.Success(w, map[string]any{"htllist": htllist, "fslist": fslist})
}

// cttCountPro POST /saas-contract/chainContract/countPro 产品交易份额(filter=true 无租户)。
// Java ChainContractServiceImpl.countPro:784 → mapper countPro:76
func (s *Service) cttCountPro(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	q := "SELECT d.d_key, count(c.id) count FROM chain_contract c, chain_common_dict d WHERE 1=1"
	var args []any
	if c := f64Ptr(m, "companyId"); c != nil {
		q += " AND c.company_id LIKE ?"
		args = append(args, str(m, "companyId"))
	}
	q += " AND d.d_type='pro__middle_type' AND c.pro_sma_type=d.d_value GROUP BY d.d_key"
	rows, err := d.Query(q, args...)
	if err != nil {
		platform.Unknown(w)
		return
	}
	ms, _ := scanMaps(rows)
	out := []map[string]any{}
	for _, row := range ms {
		out = append(out, map[string]any{"d_key": bts(row["d_key"]), "count": num(row["count"])})
	}
	platform.Success(w, out)
}
