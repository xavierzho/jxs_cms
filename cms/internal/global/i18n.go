package global

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"data_backend/pkg/i18n"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	validatorPkg "github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"golang.org/x/text/language"
)

var (
	I18n *i18n.Bundle
	UT   *ut.UniversalTranslator
)

// 配置i18n
func SetupI18n() (err error) {
	I18n = i18n.NewI18n(Language)
	err = I18n.SetPath(filepath.Join(StoragePath, "i18n", "en"), language.English, i18n.I18N_FILE_TYPE_TOML)
	if err != nil {
		if _, ok := err.(*fs.PathError); ok {
			Logger.Warn("i18n/en is not exist")
		} else {
			return err
		}
	}
	err = I18n.SetPath(filepath.Join(StoragePath, "i18n", "zh"), language.Chinese, i18n.I18N_FILE_TYPE_TOML)
	if err != nil {
		if _, ok := err.(*fs.PathError); ok {
			Logger.Warn("i18n/zh is not exist")
		} else {
			return err
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
		err = en_translations.RegisterDefaultTranslations(v, enTrans)
		if err != nil {
			return err
		}
		err = zh_translations.RegisterDefaultTranslations(v, zhTrans)
		if err != nil {
			return err
		}
	}
	UT = uni
	return nil
}
