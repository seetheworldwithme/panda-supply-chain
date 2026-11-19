package commonrest

import (
	"net/http"
	"strconv"
	"time"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

// contractImport mirrors ChainContractServiceImpl.importExcel/saveImportResult.
// Prefix is chosen by the endpoint, while contractType remains the spreadsheet's
// value. Java's supplier writes happen before the transactional contract batch.
func (s *Service) contractImport(w http.ResponseWriter, r *http.Request, prefix string) {
	fh, ok := formFile(r)
	if !ok {
		platform.Unknown(w)
		return
	}
	d := s.pick(r)
	dict, err := s.contractDictionary(r)
	if err != nil {
		platform.Unknown(w)
		return
	}
	left := prefix + currentUser(r).EnterpriseCode + time.Now().Format("20060102")
	var code string
	if err := d.QueryRow("SELECT IFNULL(max(system_code),'') FROM chain_contract WHERE system_code LIKE ?", left+"%").Scan(&code); err != nil {
		platform.Unknown(w)
		return
	}
	if code == "" {
		code = left + "0000000"
	}
	e := authRestEnts["ChainContract"]
	irt := readExcelFile(fh, e.excel, code, false, dict)
	if len(irt.ErrInfos) > 0 {
		platform.FailData(w, 999, "失败", irt.ErrInfos)
		return
	}
	if err := s.saveContractRows(r, irt.DataList); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, "导入成功，共"+strconv.Itoa(len(irt.DataList))+"条")
}
