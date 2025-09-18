/*
只支持字符串 作为结果值
支持map, array, []map 结构
支持{#...} 形式的文本替换。例如：
	有翻译文本 InvalidParams：“无效参数”
	传入 “err: {#InvalidParams}” 将得到 err: 无效参数
*/

package i18n

import (
	"context"
	"fmt"
	"os"
	"regexp"

	"golang.org/x/text/language"
)

type i18nFileType string
type i18nCtxKey string

const (
	I18N_FILE_TYPE_TOML i18nFileType = "toml"
	I18N_FILE_TYPE_JSON i18nFileType = "json"

	CTX_LANGUAGE_KEY i18nCtxKey = "I18nLanguage"

	NESTED_SEPARATOR = "."
)

var errInvalidTranslationValue = fmt.Errorf("invalid translation value")

func WithLanguage(ctx context.Context, language language.Tag) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, CTX_LANGUAGE_KEY, language.String())
}

func LanguageFromCtx(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	v := ctx.Value(CTX_LANGUAGE_KEY)
	if v != nil {
		return v.(string)
	}
	return ""
}

func getFilesAndDirs(dirPath string) (files []string, dirs []string, err error) {
	dir, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, nil, err
	}

	pathSep := string(os.PathSeparator)
	for _, fi := range dir {
		if fi.IsDir() {
			dirs = append(dirs, dirPath+pathSep+fi.Name())
			sonFiles, sonDirs, err := getFilesAndDirs(dirPath + pathSep + fi.Name())
			if err != nil {
				return nil, nil, err
			}
			files = append(files, sonFiles...)
			dirs = append(dirs, sonDirs...)
		} else {
			files = append(files, dirPath+pathSep+fi.Name())
		}
	}

	return files, dirs, nil
}

func findContent(content string) [][]string {
	return regexp.MustCompile("{#.+?}").FindAllStringSubmatch(content, -1)
}

// 若文本中间出现需要翻译的内容则需要进行装饰
func DecorateContent(content string) string {
	return fmt.Sprintf("{#%s}", content)
}

func getTranslateText(group string, content string) string {
	var text = content
	if group != "" {
		if content != "" {
			text = group + NESTED_SEPARATOR + content
		} else {
			text = group
		}
	}

	return text
}
