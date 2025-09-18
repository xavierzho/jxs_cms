package database

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"gorm.io/gorm"
)

func PaginateParams(page, pageSize int) (int, int) {
	offset := (page - 1) * pageSize
	return offset, pageSize
}

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset, pageSize := PaginateParams(page, pageSize)
		return db.Offset(offset).Limit(pageSize)
	}
}

func Order(field, order string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if field == "" {
			return db
		}
		return db.Order(fmt.Sprintf("%s %s", field, order))
	}
}

type QueryWhere struct {
	Prefix string
	Value  []interface{}
}

func ScopeQuery(params []QueryWhere) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(params) == 0 {
			return db
		}
		for _, param := range params {
			db = db.Where(param.Prefix, param.Value...)
		}
		return db
	}
}

type QueryWhereGroup []QueryWhere

func (q QueryWhereGroup) GetQuerySqlParams() (string, []interface{}) {
	whereParams, sqlParam := q.GetQueryWhereParams()
	return strings.Join(whereParams, " AND "), sqlParam
}

func (q QueryWhereGroup) GetQueryWhereParams() (whereParams []string, sqlParam []interface{}) {
	whereArr := make([]string, 0, len(q))    // Where sql语句拼接数组
	params := make([]interface{}, 0, len(q)) // sql语句值参数
	for _, param := range q {
		whereArr = append(whereArr, param.Prefix)
		params = append(params, param.Value...)
	}
	return whereArr, params
}

type QueryJoin struct {
	Prefix string
	Value  []interface{}
}

func ScopeJoinQuery(params []QueryJoin) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(params) == 0 {
			return db
		}
		for _, param := range params {
			db = db.Joins(param.Prefix, param.Value...)
		}
		return db
	}
}

func ParseResultSet(rows *sql.Rows) (*ProcResultSet, error) {
	defer rows.Close()
	var ret [][]map[string]interface{} // 多个recordSet
	var retFields [][]string
	recordSet, recordSetFields, err := ParseRowsRet(rows)
	if err != nil {
		return nil, err
	}
	retFields = append(retFields, recordSetFields)
	ret = append(ret, recordSet)
	for rows.NextResultSet() {
		nextRecordSet, nextRecordSetFields, err := ParseRowsRet(rows)
		if err != nil {
			return nil, err
		}
		retFields = append(retFields, nextRecordSetFields)
		ret = append(ret, nextRecordSet)
	}
	return &ProcResultSet{RetData: ret, RetFields: retFields}, nil
}

func ParseRowsRet(rows *sql.Rows) ([]map[string]interface{}, []string, error) {
	var recordSet []map[string]interface{}
	var recordSetFields []string
	for rows.Next() {
		record := make(map[string]interface{})
		c, err := rows.Columns()
		if err != nil || c == nil {
			return nil, nil, err
		}
		if len(recordSetFields) == 0 {
			recordSetFields = c
		}

		if len(c) == 2 && strings.ToLower(c[0]) == "msg" && strings.ToLower(c[1]) == "code" {
			var msg interface{}
			var code interface{}
			err = rows.Scan(&msg, &code)
			if err != nil {
				return nil, nil, err
			}
			return nil, nil, fmt.Errorf("code: %s - msg: %s", code, msg)
		}

		tt, err := rows.ColumnTypes()
		if err != nil {
			return nil, nil, fmt.Errorf("rows.ColumnTypes(): %w", err)
		}
		types := make([]reflect.Type, len(tt))
		for i, tp := range tt {
			st := tp.ScanType()
			if st == nil {
				return nil, nil, fmt.Errorf("scantype is null for column %q", tp.Name())
			}
			types[i] = st
		}
		data := make([]interface{}, len(tt))
		for i := range data {
			data[i] = reflect.New(types[i]).Interface()
		}
		err = rows.Scan(data...)
		if err != nil {
			return nil, nil, fmt.Errorf("rows Scan %s", err.Error())
		}
		for i := 0; i < len(data); i++ {
			record[c[i]] = parseSQLType(data[i])
		}
		recordSet = append(recordSet, record)
	}
	return recordSet, recordSetFields, nil
}

// 转换数据库原始类型
func parseSQLType(value interface{}) interface{} {
	switch val := value.(type) {
	case *sql.NullBool:
		v, _ := val.Value()
		return v
	case *sql.NullByte:
		v, _ := val.Value()
		return v
	case *sql.NullFloat64:
		v, _ := val.Value()
		return v
	case *sql.NullInt16:
		v, _ := val.Value()
		return v
	case *sql.NullInt32:
		v, _ := val.Value()
		return v
	case *sql.NullInt64:
		v, _ := val.Value()
		return v
	case *sql.NullString:
		v, _ := val.Value()
		return v
	case *sql.NullTime:
		v, _ := val.Value()
		return v
	case *sql.RawBytes:
		if val != nil {
			return string(*val)
		}
		return ""
	default:
		return val
	}
}
