package translations

import (
	"data_backend/internal/app"
	"data_backend/internal/global"
	"data_backend/pkg/i18n"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// 设置默认翻译器为英文
func Translations() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, ok := binding.Validator.Engine().(*validator.Validate)
		if ok {
			ctx.Request = ctx.Request.WithContext(i18n.WithLanguage(ctx.Request.Context(), global.Language))
		}
		enTran, _ := global.UT.GetTranslator("en")
		ctx.Set(app.TRANS_KEY, enTran)

		ctx.Next()
	}
}
