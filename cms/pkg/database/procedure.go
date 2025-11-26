package database

import (
	"fmt"
	"strings"
)

type ProcParamsType string

const (
	ProcParamsInput    ProcParamsType = "input"
	ProcParamsInoutput ProcParamsType = "inoutput"
	ProcParamsOutput   ProcParamsType = "output"
)

type ProcValueType string

const (
	ProcValueInt    ProcValueType = "int"
	ProcValueString ProcValueType = "string"
	ProcValueBool   ProcValueType = "bool"
)

// SqlParameter sql参数
type SqlParameter struct {
	Name       string         // 参数名
	Value      interface{}    // 参数值
	ValueType  ProcValueType  // 值类型
	ParamsType ProcParamsType // 参数类型
	// Size  int32
}

func (p SqlParameter) ValueString() string {
	switch p.ValueType {
	case ProcValueInt, ProcValueBool:
		return fmt.Sprintf("%v", p.Value)
	case ProcValueString:
		value := strings.ReplaceAll(fmt.Sprintf("%s", p.Value), "'", `\'`)
		return fmt.Sprintf("'%s'", value)
	}
	return ""
}

func newSqlIntParameter(name string, value int64, procParamsType ProcParamsType) SqlParameter {
	return SqlParameter{
		Name:       name,
		Value:      value,
		ValueType:  ProcValueInt,
		ParamsType: procParamsType,
	}
}

func newSqlStringParameter(name string, value string, procParamsType ProcParamsType) SqlParameter {
	return SqlParameter{
		Name:       name,
		Value:      value,
		ValueType:  ProcValueString,
		ParamsType: procParamsType,
	}
}

func newSqlBoolParameter(name string, value bool, procParamsType ProcParamsType) SqlParameter {
	return SqlParameter{
		Name:       name,
		Value:      value,
		ValueType:  ProcValueBool,
		ParamsType: procParamsType,
	}
}

// 参数数组
type Parameters struct {
	innerParameters []*SqlParameter
}

// 创建Sql参数集合
func NewSqlParameters() *Parameters {
	return new(Parameters)
}

// 添加参数
func (ps *Parameters) addParameter(p SqlParameter) {
	ps.innerParameters = append(ps.innerParameters, &p)
}

// 添加int input 参数
func (ps *Parameters) AddIntInputParameter(name string, value int64) {
	p := newSqlIntParameter(name, value, ProcParamsInput)
	ps.addParameter(p)
}

// 添加int inOutput 参数
func (ps *Parameters) AddIntInOutputParameter(name string, value int64) {
	p := newSqlIntParameter(name, value, ProcParamsInoutput)
	ps.addParameter(p)
}

// 添加int output 参数
func (ps *Parameters) AddIntOutputParameter(name string) {
	p := newSqlIntParameter(name, 0, ProcParamsOutput)
	ps.addParameter(p)
}

// 添加string input 参数
func (ps *Parameters) AddStringInputParameter(name string, value string) {
	p := newSqlStringParameter(name, value, ProcParamsInput)
	ps.addParameter(p)
}

// 添加string inOutput 参数
func (ps *Parameters) AddStringInOutputParameter(name string, value string) {
	p := newSqlStringParameter(name, value, ProcParamsInoutput)
	ps.addParameter(p)
}

// 添加string output 参数
func (ps *Parameters) AddStringOutputParameter(name string) {
	p := newSqlStringParameter(name, "", ProcParamsOutput)
	ps.addParameter(p)
}

// 添加bool input 参数
func (ps *Parameters) AddBoolInputParameter(name string, value bool) {
	p := newSqlBoolParameter(name, value, ProcParamsInput)
	ps.addParameter(p)
}

// 添加bool inOutput 参数
func (ps *Parameters) AddBoolInOutputParameter(name string, value bool) {
	p := newSqlBoolParameter(name, value, ProcParamsInoutput)
	ps.addParameter(p)
}

// 添加bool output 参数
func (ps *Parameters) AddBoolOutputParameter(name string) {
	p := newSqlBoolParameter(name, false, ProcParamsOutput)
	ps.addParameter(p)
}

// 存储过程
type Procedure struct {
	Name       string //存储过程名称
	Parameters *Parameters
}

// 创建存储过程
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
		case ProcParamsInput:
			callSql += fmt.Sprintf("%s, ", pr.ValueString())
		case ProcParamsOutput:
			callSql += fmt.Sprintf("@%s, ", pr.Name)
			selectSql += fmt.Sprintf("@%s, ", pr.Name)
		case ProcParamsInoutput:
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
