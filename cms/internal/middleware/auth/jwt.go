package auth

import (
	"encoding/hex"
	"errors"
	"strconv"

	"data_backend/internal/app"
	"data_backend/internal/dao"
	"data_backend/internal/global"
	"data_backend/pkg/convert"
	"data_backend/pkg/errcode"
	"data_backend/pkg/i18n"
	"data_backend/pkg/logger"
	"data_backend/pkg/redisdb"
	"data_backend/pkg/token"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"golang.org/x/text/language"
)

type JWT struct {
	rdb       *redisdb.RedisClient
	logger    *logger.Logger
	checkPerm func(ctx *gin.Context, userID uint32, permList []string) (bool, error)
}

func NewJWT(rdb *redisdb.RedisClient, logger *logger.Logger, checkPerm func(ctx *gin.Context, userID uint32, permList []string) (bool, error)) JWT {
	return JWT{
		rdb:       rdb,
		logger:    logger,
		checkPerm: checkPerm,
	}

}

func (j JWT) JWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var tokenStr string
		var eCode *errcode.Error
		j.logger = j.logger.WithContext(ctx)
		defer app.DeferResponse(ctx, &eCode)

		// 获取请求 Token
		if s, exist := ctx.GetQuery(app.TOKEN_KEY); exist {
			tokenStr = s
		} else {
			tokenStr = ctx.GetHeader(app.TOKEN_KEY)
		}
		if tokenStr == "" {
			eCode = errcode.UnauthorizedTokenError
			return
		}

		// ParseToken
		claims, err := token.ParseToken(tokenStr)
		if err != nil {
			switch err.(*jwt.ValidationError).Errors { // Errors 是 通过 | 将多个错误拼起来的
			case jwt.ValidationErrorExpired:
				eCode = errcode.UnauthorizedTokenTimeout
			default:
				eCode = errcode.UnauthorizedTokenError
			}
			return
		}

		// decode info
		userIDByte, _ := hex.DecodeString(claims.UserID)
		ctx.Set(app.USER_ID_KEY, string(userIDByte))
		userID, err := strconv.ParseUint(string(userIDByte), 10, 32)
		if err != nil {
			eCode = errcode.UnauthorizedTokenError
			return
		}

		// 验证token
		tokenKey := token.GetRKeyByUserID(uint32(userID))
		tokenCache, err := j.rdb.Get(ctx, tokenKey).Result()
		if errors.Is(err, redis.Nil) {
			eCode = errcode.UnauthorizedTokenTimeout
			return
		} else if err != nil {
			j.logger.Errorf("RedisClient.Get %s: %s", tokenKey, err.Error())
		} else if tokenCache != tokenStr {
			eCode = errcode.UnauthorizedTokenTimeout
			return
		}
		// 若存在权限则允许其使用其他语言返回
		eCode = j.changeTranslator(ctx)
	}
}

func (j JWT) changeTranslator(ctx *gin.Context) *errcode.Error {
	if userID, ok := ctx.Get(app.USER_ID_KEY); ok {
		intUserID := convert.StrTo(userID.(string)).MustUInt32()
		hasPerm, err := j.checkPerm(ctx, intUserID, []string{dao.PERMISSION_LANG_UPDATE})
		if err != nil {
			return errcode.ServerError
		}

		if hasPerm {
			locale := ctx.GetHeader(app.LOCALE_KEY)
			if _, ok := binding.Validator.Engine().(*validator.Validate); ok {
				switch locale {
				case "zh":
					ctx.Request = ctx.Request.WithContext(i18n.WithLanguage(ctx.Request.Context(), language.Chinese))
				default:
					ctx.Request = ctx.Request.WithContext(i18n.WithLanguage(ctx.Request.Context(), language.English))
				}
			}
			if tran, ok := global.UT.GetTranslator(locale); ok {
				ctx.Set(app.TRANS_KEY, tran)
			}
		}
	}

	return nil
}
