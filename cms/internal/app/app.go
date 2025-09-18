/*
格式化返回 http 请求
*/
package app

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"data_backend/internal/global"
	"data_backend/pkg/errcode"
	"data_backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type Response struct {
	Ctx *gin.Context
}

func NewResponse(ctx *gin.Context) *Response {
	return &Response{
		Ctx: ctx,
	}
}

func (r *Response) ToResponse(data gin.H) {
	r.Ctx.JSON(http.StatusOK, data)
}

func (r *Response) ToResponseData(data interface{}) {
	if data == nil {
		data = gin.H{}
	}
	r.Ctx.JSON(http.StatusOK, gin.H{"data": data})
}

func (r *Response) ToResponseOK() {
	r.Ctx.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func (r *Response) ToResponseList(data interface{}, totalRows int64) {
	r.Ctx.JSON(http.StatusOK, gin.H{
		"data": data,
		"headers": Pager{
			Page:      GetPage(r.Ctx),
			PageSize:  GetPageSize(r.Ctx),
			TotalRows: totalRows,
		},
	})
}

func (r *Response) ToResponseDetail(code int, data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		r.Ctx.JSON(http.StatusOK, gin.H{
			"code": code,
			"data": data,
		})
		return
	}
	if string(b) == "null" {
		data = make(map[string]interface{})
	}
	r.Ctx.JSON(http.StatusOK, gin.H{
		"code": code,
		"data": data,
	})
}

func (r *Response) ToResponseBase64(data interface{}) {
	if data == nil {
		data = gin.H{}
	}
	b, err := json.Marshal(data)
	if err != nil {
		r.ToErrorResponse(errcode.TransformFail.WithDetails("json.Marshal error: " + err.Error()))
		return
	}
	r.Ctx.JSON(http.StatusOK, base64.StdEncoding.EncodeToString(b))
}

func (r *Response) ToErrorResponse(err *errcode.Error) {
	response := gin.H{
		"code": err.Code(),
		"msg":  global.I18n.T(r.Ctx.Request.Context(), "error", err.Msg()),
	}
	details := err.Details()
	if len(details) > 0 {
		tranDetails := make([]string, 0, len(details))
		// 将错误组进行国际化
		for _, detail := range details {
			tranDetails = append(tranDetails, global.I18n.T(r.Ctx.Request.Context(), "error", detail))
		}
		response["details"] = tranDetails
	}

	r.Ctx.JSON(err.HtmlCode(), response)
}

func (r *Response) BindAndValid(ctx *gin.Context, params interface{}, log *logger.Logger) bool {
	valid, errs := BindAndValid(ctx, params)
	if !valid {
		log.Errorf("app.BindAndValid errs: %v", errs)
		r.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return false
	}

	if global.ServerSetting.RunMode == global.RUN_MODE_DEBUG {
		fmt.Printf("Bind params: %+v\n", params)
	}

	return true
}

// 用于 返回文件
func (r *Response) ExportFile(ctx *gin.Context, excelModel *excelize.File, fileName string) (err error) {
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Access-Control-Expose-Headers", "Content-Disposition")
	ctx.Header("Content-Disposition", "attachment; filename="+fileName+".xlsx")
	ctx.Header("Content-Transfer-Encoding", "binary")
	if err = excelModel.Write(ctx.Writer); err != nil {
		r.ToErrorResponse(errcode.ExportFail.WithDetails(err.Error()))
		return
	}

	return nil
}

// 用于中间件处理
func DeferResponse(ctx *gin.Context, eCode **errcode.Error) {
	if *eCode != nil {
		response := NewResponse(ctx)
		response.ToErrorResponse(*eCode)
		ctx.Abort()
		return
	}

	ctx.Next()
}
