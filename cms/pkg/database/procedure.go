package database

import (
	"fmt"
	"strings"
)

type ProcParamsType string

const (
	PROC_PARAMS_INPUT    ProcParamsType = "input"
	PROC_PARAMS_INOUTPUT ProcParamsType = "inoutput"
	PROC_PARAMS_OUTPUT   ProcParamsType = "output"
)

type ProcValueType string

const (
	PROC_VALUE_INT    ProcValueType = "int"
	PROC_VALUE_STRING ProcValueType = "string"
	PROC_VALUE_BOOl   ProcValueType = "bool"
)

//sql参数
type SqlParameter struct {
	Name       string         // 参数名
	Value      interface{}    // 参数值
	ValueType  ProcValueType  // 值类型
	ParamsType ProcParamsType // 参数类型
	// Size  int32
}

func (p SqlParameter) ValueString() string {
	switch p.ValueType {
	case PROC_VALUE_INT, PROC_VALUE_BOOl:
		return fmt.Sprintf("%v", p.Value)
	case PROC_VALUE_STRING:
		value := strings.ReplaceAll(fmt.Sprintf("%s", p.Value), "'", `\'`)
		return fmt.Sprintf("'%s'", value)
	}
	return ""
}

func newSqlIntParameter(name string, value int64, procParamsType ProcParamsType) SqlParameter {
	return SqlParameter{
		Name:       name,
		Value:      value,
		ValueType:  PROC_VALUE_INT,
		ParamsType: procParamsType,
	}
}

func newSqlStringParameter(name string, value string, procParamsType ProcParamsType) SqlParameter {
	return SqlParameter{
		Name:       name,
		Value:      value,
		ValueType:  PROC_VALUE_STRING,
		ParamsType: procParamsType,
	}
}

func newSqlBoolParameter(name string, value bool, procParamsType ProcParamsType) SqlParameter {
	return SqlParameter{
		Name:       name,
		Value:      value,
		ValueType:  PROC_VALUE_BOOl,
		ParamsType: procParamsType,
	}
}

//参数数组
type Parameters struct {
	innerParameters []*SqlParameter
}

//创建Sql参数集合
func NewSqlParameters() *Parameters {
	return new(Parameters)
}

//添加参数
func (ps *Parameters) addParameter(p SqlParameter) {
	ps.innerParameters = append(ps.innerParameters, &p)
}

//添加int input 参数
func (ps *Parameters) AddIntInputParameter(name string, value int64) {
	p := newSqlIntParameter(name, value, PROC_PARAMS_INPUT)
	ps.addParameter(p)
}

//添加int inOutput 参数
func (ps *Parameters) AddIntInOutputParameter(name string, value int64) {
	p := newSqlIntParameter(name, value, PROC_PARAMS_INOUTPUT)
	ps.addParameter(p)
}

//添加int output 参数
func (ps *Parameters) AddIntOutputParameter(name string) {
	p := newSqlIntParameter(name, 0, PROC_PARAMS_OUTPUT)
	ps.addParameter(p)
}

//添加string input 参数
func (ps *Parameters) AddStringInputParameter(name string, value string) {
	p := newSqlStringParameter(name, value, PROC_PARAMS_INPUT)
	ps.addParameter(p)
}

//添加string inOutput 参数
func (ps *Parameters) AddStringInOutputParameter(name string, value string) {
	p := newSqlStringParameter(name, value, PROC_PARAMS_INOUTPUT)
	ps.addParameter(p)
}

//添加string output 参数
func (ps *Parameters) AddStringOutputParameter(name string) {
	p := newSqlStringParameter(name, "", PROC_PARAMS_OUTPUT)
	ps.addParameter(p)
}

//添加bool input 参数
func (ps *Parameters) AddBoolInputParameter(name string, value bool) {
	p := newSqlBoolParameter(name, value, PROC_PARAMS_INPUT)
	ps.addParameter(p)
}

//添加bool inOutput 参数
func (ps *Parameters) AddBoolInOutputParameter(name string, value bool) {
	p := newSqlBoolParameter(name, value, PROC_PARAMS_INOUTPUT)
	ps.addParameter(p)
}

//添加bool output 参数
func (ps *Parameters) AddBoolOutputParameter(name string) {
	p := newSqlBoolParameter(name, false, PROC_PARAMS_OUTPUT)
	ps.addParameter(p)
}

//存储过程
type Procedure struct {
	Name       string //存储过程名称
	Parameters *Parameters
}

//创建存储过程
func NewProcedure(name string, ps *Parameters) *Procedure {
	// 规定所有的存储过程都有该参数
	ps.AddIntInOutputParameter("intRet", 0)
	ps.AddIntInOutputParameter("intEnCode", 0)
	ps.AddStringInOutputParameter("chvErrMsg", "")

	p := new(Procedure)
	p.Name = name
	p.Parameters = ps

	return p
}

// 获取执行SQL
func (p *Procedure) GetProcedureSql() (sqlString string) {
	var setSql string    // 设置用户变量
	var callSql string   // 调用存储过程
	var selectSql string // select结果

	callSql += fmt.Sprintf("call %s(", p.Name)
	for _, pr := range p.Parameters.innerParameters {
		switch pr.ParamsType {
		case PROC_PARAMS_INPUT:
			callSql += fmt.Sprintf("%s, ", pr.ValueString())
		case PROC_PARAMS_OUTPUT:
			callSql += fmt.Sprintf("@%s, ", pr.Name)
			selectSql += fmt.Sprintf("@%s, ", pr.Name)
		case PROC_PARAMS_INOUTPUT:
			setSql += fmt.Sprintf("set @%s=%s; ", pr.Name, pr.ValueString())
			callSql += fmt.Sprintf("@%s, ", pr.Name)
			selectSql += fmt.Sprintf("@%s, ", pr.Name)
		}
	}
	callSql = strings.TrimSuffix(callSql, ", ") + "); "

	if setSql != "" {
		sqlString += setSql + "\n"
	}

	sqlString += callSql + "\n"

	if selectSql != "" {
		sqlString += "select " + strings.TrimSuffix(selectSql, ", ") + ";"
	}
	return
}

// 存储过程结果集
type ProcResultSet struct {
	RetData   [][]map[string]interface{}
	RetFields [][]string
}

// 判断存储过程结果集是否有output参数之外的结果集
func (p *ProcResultSet) HaveResultSet() bool {
	return len(p.RetData) > 1
}

// 获取存储过程结果集的output参数结果
func (p *ProcResultSet) GetOutputParam(field string) interface{} {
	retLen := len(p.RetData)
	return p.RetData[retLen-1][0]["@"+field]
}

// 存储过程执行结果信息
type ProcReturnResult struct {
	IntRet    int
	IntEnCode int
	ChvErrMsg string
}
