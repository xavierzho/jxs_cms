package i18n

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"golang.org/x/text/language"
)

func TestDemo(t *testing.T) {
	ctx := context.Background()
	bundle := NewI18n(language.English)
	err := bundle.SetPath("../../internal/i18n/en", language.English, FileTypeToml)
	if err != nil {
		fmt.Println("en:", err)
	}
	err = bundle.SetPath("../../internal/i18n/zh-CN", language.Chinese, FileTypeToml)
	if err != nil {
		fmt.Println("zh:", err)
	}

	fmt.Println(bundle.DefaultLanguage())
	fmt.Println(bundle.FilePath())
	fmt.Println(bundle.LanguageTags())

	ctx = WithLanguage(ctx, language.Chinese)
	fmt.Println("trans")
	fmt.Println(bundle.T(ctx, "error", "ServerError"))

	fmt.Println(bundle.FilePath())
	fmt.Println(bundle.LanguageTags())

	fmt.Println(bundle.T(ctx, "enum", "SMSLogStateEnum"))

	bundle.UnsetPath()

	fmt.Println(bundle.T(ctx, "enum", "SMSLogStateEnum"))

	fmt.Println(bundle.T(ctx, "enum", "aaaa"))
	transContent := bundle.T(ctx, "enum", "SMSLogStateEnum")
	var data []string
	err = json.Unmarshal([]byte(transContent), &data)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(data)
}

func TestErrorCodeT(t *testing.T) {
	ctx := context.Background()
	bundle := NewI18n(language.Chinese)
	err := bundle.SetPath("../../cmd/cms_backend/storage/i18n/zh", language.Chinese, FileTypeToml)
	if err != nil {
		fmt.Println("zh:", err)
	}

	ctx = WithLanguage(ctx, language.Chinese)
	fmt.Println("trans")
	// fmt.Println(bundle.T(ctx, "error", errcode.InvalidParams.Error()))
	fmt.Println(bundle.T(ctx, "error", "123 {#tasgdfasf} {#InvalidParams}"))
	fmt.Println(bundle.T(ctx, "error", "InvalidParams"))
}

func TestAllType(t *testing.T) {
	ctx := context.Background()
	bundle := NewI18n(language.English)
	err := bundle.SetPath("./test.toml", language.Chinese, FileTypeToml)
	if err != nil {
		fmt.Println("zh:", err)
	}

	ctx = WithLanguage(ctx, language.Chinese)
	fmt.Println("trans")
	fmt.Println(bundle.T(ctx, "", "one"))
	fmt.Println(bundle.T(ctx, "", "arr.0"))
	fmt.Println(bundle.T(ctx, "error", "InvalidParams"))
	fmt.Println(bundle.T(ctx, "error", "error: {#InvalidParams}"))
	fmt.Println(bundle.T(ctx, "nestMap", "firstMap.Map_1.name"))
	fmt.Println(bundle.T(ctx, "arrMap", "0.name"))
}

func TestAllTypeJson(t *testing.T) {
	ctx := context.Background()
	bundle := NewI18n(language.English)
	err := bundle.SetPath("./test.json", language.Chinese, FileTypeJson)
	if err != nil {
		fmt.Println("zh:", err)
	}

	ctx = WithLanguage(ctx, language.Chinese)
	fmt.Println("trans")
	fmt.Println(bundle.T(ctx, "", "one"))
	fmt.Println(bundle.T(ctx, "", "arr.0"))
	fmt.Println(bundle.T(ctx, "error", "InvalidParams"))
	fmt.Println(bundle.T(ctx, "error", "error: {#InvalidParams}"))
	fmt.Println(bundle.T(ctx, "nestMap", "firstMap.Map_1.name"))
	fmt.Println(bundle.T(ctx, "arrMap", "0.name"))
}

func TestNestedMenu(t *testing.T) {
	ctx := context.Background()
	bundle := NewI18n(language.Chinese)
	err := bundle.SetPath("../../cmd/cms_backend/storage/i18n/zh", language.Chinese, FileTypeToml)
	if err != nil {
		fmt.Println("error zh:", err)
	}

	ctx = WithLanguage(ctx, language.Chinese)
	fmt.Println("trans")
	fmt.Println(bundle.T(ctx, "menu", "Report"))
	fmt.Println(bundle.T(ctx, "menu", "Report.SummaryRevenue"))
}

func TestFindContent(t *testing.T) {
	fmt.Println(findContent("{#test}{#test}{#test2}{#test1}{#test2}"))
}
