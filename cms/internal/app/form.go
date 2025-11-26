package app

/*
表单验证
*/
import (
	"strings"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	val "github.com/go-playground/validator/v10"
)

type ValidError struct {
	Key     string
	Message string
}

func (v *ValidError) Error() string {
	return v.Message
}

type ValidErrors []*ValidError

func (v ValidErrors) Error() string {
	return strings.Join(v.Errors(), ",")
}

func (v ValidErrors) Errors() (errs []string) {
	for _, err := range v {
		errs = append(errs, err.Error())
	}

	return
}

func BindAndValid(ctx *gin.Context, v interface{}) (bool, ValidErrors) {
	var errs ValidErrors
	err := ctx.ShouldBind(v)
	if err != nil {
		v := ctx.Value(TRANS_KEY)
		trans, _ := v.(ut.Translator)
		validErrs, ok := err.(val.ValidationErrors)
		if !ok {
			return false, errs
		}

		for key, value := range validErrs.Translate(trans) {
			errs = append(errs, &ValidError{
				Key:     key,
				Message: value,
			})
		}

		return false, errs
	}

	// TODO 添加 v 定义的 valid

	return true, nil
}
