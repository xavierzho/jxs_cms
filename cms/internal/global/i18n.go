package global

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"

	"data_backend/pkg/i18n"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	validatorPkg "github.com/go-playground/validator/v10"
	entranslations "github.com/go-playground/validator/v10/translations/en"
	zhtranslations "github.com/go-playground/validator/v10/translations/zh"
	"golang.org/x/text/language"
)

var (
	I18n *i18n.Bundle
	UT   *ut.UniversalTranslator
)

// SetupI18n 配置i18n
func SetupI18n() (err error) {
	I18n = i18n.NewI18n(Language)
	err = I18n.SetPath(filepath.Join(StoragePath, "i18n", "en"), language.English, i18n.FileTypeToml)
	if err != nil {
		var pathError *fs.PathError
		if errors.As(err, &pathError) {
			Logger.Warn("i18n/en is not exist")
		}
	}
	err = I18n.SetPath(filepath.Join(StoragePath, "i18n", "zh"), language.Chinese, i18n.FileTypeToml)
	if err != nil {
		var pathError *fs.PathError
		if errors.As(err, &pathError) {
			Logger.Warn("i18n/zh is not exist")
		}
	}
	return nil
}

func SetupValidator() (err error) {
	// 注册验证器的i18n
	uni := ut.New(en.New(), en.New(), zh.New())
	enTrans, found := uni.GetTranslator("en")
	if !found {
		return fmt.Errorf("not found en translator")
	}
	zhTrans, found := uni.GetTranslator("zh")
	if !found {
		return fmt.Errorf("not found zh translator")
	}
	v, ok := binding.Validator.Engine().(*validatorPkg.Validate)
	if ok {
		err = entranslations.RegisterDefaultTranslations(v, enTrans)
		if err != nil {
			return err
		}
		err = zhtranslations.RegisterDefaultTranslations(v, zhTrans)
		if err != nil {
			return err
		}
	}
	UT = uni
	return nil
}
