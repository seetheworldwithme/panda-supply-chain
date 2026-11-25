// Package platform - dbutil.go MyBatis-Plus 内置方法的 Go 等价物。
// 二期组织基础接口大量依赖 MP 的 getById/save/updateById/removeById/removeByIds/
// list(QueryWrapper.allEq)/saveBatch/updateBatchById/removeByMap,此处统一翻译。
// 约定:列名只接受 ^[a-z_][a-zA-Z0-9_]*$(防列名注入;chain_capital_bkcall.msgText、
// chain_capital_transinfo.orderSubtype 等库内驼峰列放行,与 Java 端 allEq 直拼键名的
// 行为差异是刻意收紧);值一律参数化。时间列由调用方按 Java MP 自动填充行为显式带上
// (create_time/update_time)。
package platform

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"
	"time"
)

var colRe = regexp.MustCompile(`^[a-z_][a-zA-Z0-9_]*$`)

// CheckCols 校验列名合法(非法直接 panic 语义错误,编程期即暴露)
func CheckCols(cols []string) {
	for _, c := range cols {
		if !colRe.MatchString(c) {
			panic(fmt.Sprintf("非法列名: %q", c))
		}
	}
}

// GetByID 等价 MP getById:SELECT * FROM t WHERE id=?
func GetByID(d *sql.DB, table string, id any) (map[string]any, error) {
	rows, err := d.Query("SELECT * FROM "+table+" WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	ms, err := ScanMaps(rows)
	if err != nil || len(ms) == 0 {
		return nil, err
	}
	return ms[0], nil
}

// ListWhere 等价 MP list(QueryWrapper.allEq(maps)):全部等值 AND 条件。
// cond 键为下划线列名,空 map 查全表;显式 null → IS NULL(单参 allEq null2IsNull=true)。
// 注:allEq(maps,false) 调用方应先在 cond 构造时剔除 null 键。
func ListWhere(d *sql.DB, table string, cond map[string]any) ([]map[string]any, error) {
	q := "SELECT * FROM " + table
	var args []any
	if len(cond) > 0 {
		cols := SortedKeys(cond)
		CheckCols(cols)
		var conds []string
		for _, k := range cols {
			if cond[k] == nil {
				conds = append(conds, k+" IS NULL")
			} else {
				conds = append(conds, k+" = ?")
				args = append(args, cond[k])
			}
		}
		q += " WHERE " + strings.Join(conds, " AND ")
	}
	rows, err := d.Query(q, args...)
	if err != nil {
		return nil, err
	}
	return ScanMaps(rows)
}

// autoTimeVals 对照 MP MateMetaObjectHandler:insert 自动填 create_time/update_time,
// update 自动填 update_time(strictFill 语义:调用方已带则不覆盖)。
// 实体时间字段无 @TableField(fill) 注解时 Java 不自动填(autoTime=false):
// create_time 落库默认,update_time 恒 NULL/不变。skipNull 恒真,对齐 MP
// 字段策略 FieldStrategy.NOT_NULL:null 字段不进 SQL(库默认值生效)
func autoTimeVals(vals map[string]any, insert, autoTime bool) map[string]any {
	out := make(map[string]any, len(vals)+2)
	for k, v := range vals {
		if v == nil {
			continue
		}
		out[k] = v
	}
	if !autoTime {
		return out
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	if _, has := out["create_time"]; !has && insert {
		out["create_time"] = now
	}
	if _, has := out["update_time"]; !has {
		out["update_time"] = now
	}
	return out
}

// Insert 等价 MP save:INSERT 全列(自动补时间戳)
func Insert(d *sql.DB, table string, vals map[string]any) (int64, error) {
	return insertRow(d, table, vals, true)
}

// InsertNoFill 等价无 @TableField(fill) 注解实体的 save:时间列交给库默认
// (create_time DEFAULT CURRENT_TIMESTAMP,update_time 恒 NULL),如 chain_eqb_config/
// chain_contract_supplier/chain_basic_enterprise/chain_common_file
func InsertNoFill(d *sql.DB, table string, vals map[string]any) (int64, error) {
	return insertRow(d, table, vals, false)
}

func insertRow(d *sql.DB, table string, vals map[string]any, autoTime bool) (int64, error) {
	vals = autoTimeVals(vals, true, autoTime)
	cols := SortedKeys(vals)
	CheckCols(cols)
	ph := strings.TrimSuffix(strings.Repeat("?,", len(cols)), ",")
	var args []any
	for _, k := range cols {
		args = append(args, vals[k])
	}
	q := "INSERT INTO " + table + " (" + strings.Join(cols, ",") + ") VALUES (" + ph + ")"
	res, err := d.Exec(q, args...)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// UpdateByID 等价 MP updateById:按主键更新非空子集(自动补 update_time)
func UpdateByID(d *sql.DB, table string, id any, vals map[string]any) (int64, error) {
	return updateRow(d, table, id, vals, true)
}

// UpdateByIDNoFill 等价无 fill 注解实体的 updateById:update_time 不自动刷新
// (Java strictUpdateFill 只作用于带注解字段),仅当调用方显式携带时才写
func UpdateByIDNoFill(d *sql.DB, table string, id any, vals map[string]any) (int64, error) {
	return updateRow(d, table, id, vals, false)
}

func updateRow(d *sql.DB, table string, id any, vals map[string]any, autoTime bool) (int64, error) {
	vals = autoTimeVals(vals, false, autoTime)
	cols := SortedKeys(vals)
	CheckCols(cols)
	var sets []string
	var args []any
	for _, k := range cols {
		sets = append(sets, k+" = ?")
		args = append(args, vals[k])
	}
	args = append(args, id)
	res, err := d.Exec("UPDATE "+table+" SET "+strings.Join(sets, ",")+" WHERE id = ?", args...)
	if err != nil {
		return 0, err
	}
	n, _ := res.RowsAffected()
	return n, nil
}

// DeleteByID 等价 MP removeById
func DeleteByID(d *sql.DB, table string, id any) error {
	_, err := d.Exec("DELETE FROM "+table+" WHERE id = ?", id)
	return err
}

// DeleteByIDs 等价 MP removeByIds
func DeleteByIDs(d *sql.DB, table string, ids []any) error {
	if len(ids) == 0 {
		return nil
	}
	ph := strings.TrimSuffix(strings.Repeat("?,", len(ids)), ",")
	var args []any
	args = append(args, ids...)
	_, err := d.Exec("DELETE FROM "+table+" WHERE id IN ("+ph+")", args...)
	return err
}

// TxDeleteByIDs 事务版 removeByIds
func TxDeleteByIDs(tx *sql.Tx, table string, ids []any) error {
	if len(ids) == 0 {
		return nil
	}
	ph := strings.TrimSuffix(strings.Repeat("?,", len(ids)), ",")
	var args []any
	args = append(args, ids...)
	_, err := tx.Exec("DELETE FROM "+table+" WHERE id IN ("+ph+")", args...)
	return err
}

// DeleteByMap 等价 MP removeByMap:全部等值 AND 删除
func DeleteByMap(d *sql.DB, table string, cond map[string]any) error {
	cols := SortedKeys(cond)
	CheckCols(cols)
	var conds []string
	var args []any
	for _, k := range cols {
		if cond[k] == nil {
			conds = append(conds, k+" IS NULL")
		} else {
			conds = append(conds, k+" = ?")
			args = append(args, cond[k])
		}
	}
	_, err := d.Exec("DELETE FROM "+table+" WHERE "+strings.Join(conds, " AND "), args...)
	return err
}

// WithTx 事务包裹(Java 侧多表写普遍无 @Transactional,Go 统一收进事务防半写)
func WithTx(d *sql.DB, fn func(tx *sql.Tx) error) error {
	tx, err := d.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() // 提交后再回滚只返回 ErrTxDone;fn 内 panic 时也把连接还回池里(池仅 20 条)
	if err := fn(tx); err != nil {
		return err
	}
	return tx.Commit()
}

// Tx 版本:事务内复用(列名校验同全局)
func TxGetByID(tx *sql.Tx, table string, id any) (map[string]any, error) {
	rows, err := tx.Query("SELECT * FROM "+table+" WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	ms, err := ScanMaps(rows)
	if err != nil || len(ms) == 0 {
		return nil, err
	}
	return ms[0], nil
}

func TxInsert(tx *sql.Tx, table string, vals map[string]any) (int64, error) {
	return txInsert(tx, table, vals, true)
}

func TxInsertNoFill(tx *sql.Tx, table string, vals map[string]any) (int64, error) {
	return txInsert(tx, table, vals, false)
}

func txInsert(tx *sql.Tx, table string, vals map[string]any, fill bool) (int64, error) {
	vals = autoTimeVals(vals, true, fill)
	cols := SortedKeys(vals)
	CheckCols(cols)
	ph := strings.TrimSuffix(strings.Repeat("?,", len(cols)), ",")
	var args []any
	for _, k := range cols {
		args = append(args, vals[k])
	}
	res, err := tx.Exec("INSERT INTO "+table+" ("+strings.Join(cols, ",")+") VALUES ("+ph+")", args...)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func TxUpdateByID(tx *sql.Tx, table string, id any, vals map[string]any) error {
	return txUpdate(tx, table, id, vals, true)
}

func TxUpdateByIDNoFill(tx *sql.Tx, table string, id any, vals map[string]any) error {
	return txUpdate(tx, table, id, vals, false)
}

func txUpdate(tx *sql.Tx, table string, id any, vals map[string]any, fill bool) error {
	vals = autoTimeVals(vals, false, fill)
	cols := SortedKeys(vals)
	CheckCols(cols)
	var sets []string
	var args []any
	for _, k := range cols {
		sets = append(sets, k+" = ?")
		args = append(args, vals[k])
	}
	args = append(args, id)
	_, err := tx.Exec("UPDATE "+table+" SET "+strings.Join(sets, ",")+" WHERE id = ?", args...)
	return err
}

func TxDeleteByMap(tx *sql.Tx, table string, cond map[string]any) error {
	cols := SortedKeys(cond)
	CheckCols(cols)
	var conds []string
	var args []any
	for _, k := range cols {
		if cond[k] == nil {
			conds = append(conds, k+" IS NULL")
		} else {
			conds = append(conds, k+" = ?")
			args = append(args, cond[k])
		}
	}
	_, err := tx.Exec("DELETE FROM "+table+" WHERE "+strings.Join(conds, " AND "), args...)
	return err
}
