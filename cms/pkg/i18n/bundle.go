package i18n

import (
	"context"
	"os"
	"strings"
	"sync"

	"golang.org/x/text/language"
)

type Bundle struct {
	mu               sync.Mutex
	defaultLanguage  language.Tag
	path             []string
	tag              []language.Tag
	messageTemplates map[language.Tag]map[string]*Message
	localizerMap     map[string]*Localizer
	fileMap          map[language.Tag]map[i18nFileType][]string
}

func NewI18n(defaultLanguage language.Tag) (b *Bundle) {
	b = &Bundle{mu: sync.Mutex{}}
	b.defaultLanguage = defaultLanguage
	b.messageTemplates = make(map[language.Tag]map[string]*Message)
	b.localizerMap = make(map[string]*Localizer)
	b.fileMap = make(map[language.Tag]map[i18nFileType][]string)
	return
}

func (b *Bundle) DefaultLanguage() string {
	return b.defaultLanguage.String()
}

func (b *Bundle) LanguageTags() []language.Tag {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.tag
}

func (b *Bundle) addTag(language language.Tag) {
	for _, tag := range b.tag {
		if tag == language {
			return
		}
	}
	b.tag = append(b.tag, language)
}

func (b *Bundle) FilePath() []string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.path
}

func (b *Bundle) addPath(filePath string) {
	for _, path := range b.path {
		if path == filePath {
			return
		}
	}
	b.path = append(b.path, filePath)
}

// 设置目录 用于懒加载文件
func (b *Bundle) SetPath(path string, language language.Tag, format i18nFileType) (err error) {
	_, err = os.Stat(path)
	if err != nil {
		return
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	if len(b.path) != 0 { // 设置地址时已经初始化了则重置path
		b.path = []string{}
	}

	if _, ok := b.fileMap[language]; !ok {
		b.fileMap[language] = make(map[i18nFileType][]string)
		b.fileMap[language][format] = []string{}
	}

	for _, existPath := range b.fileMap[language][format] {
		if path == existPath {
			return
		}
	}
	b.fileMap[language][format] = append(b.fileMap[language][format], path)
	return
}

// 清空文件目录目录
func (b *Bundle) UnsetPath() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.path = []string{}
	b.tag = []language.Tag{}
	b.messageTemplates = make(map[language.Tag]map[string]*Message)
	b.localizerMap = make(map[string]*Localizer)
	b.fileMap = make(map[language.Tag]map[i18nFileType][]string)
}

func (b *Bundle) addMessages(language language.Tag, messageList []*Message) {
	b.addTag(language)
	if _, ok := b.messageTemplates[language]; !ok {
		b.messageTemplates[language] = make(map[string]*Message)
	}
	for _, message := range messageList {
		b.messageTemplates[language][message.ID] = message
	}
}

func (b *Bundle) loadMessage(path string, language language.Tag, format i18nFileType) (err error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return
	}

	var messageList []*Message

	// 文件夹
	if fileInfo.IsDir() {
		messageList, err = b.loadMessageFolder(path, format)
		if err != nil {
			return
		}
	} else { // 文件
		messageList, err = b.loadMessageFile(path, format)
		if err != nil {
			return
		}
	}

	b.addMessages(language, messageList)

	return nil
}

func (b *Bundle) loadMessageFolder(dirPath string, format i18nFileType) (messageList []*Message, err error) {
	filePathList, _, err := getFilesAndDirs(dirPath)
	if err != nil {
		return
	}

	for _, filePath := range filePathList {
		buf, err := os.ReadFile(filePath)
		if err != nil {
			return nil, err
		}
		messageFile, err := parseMessageFileBytes(buf, getUnmarshalFunc(format))
		if err != nil {
			return nil, err
		}

		b.addPath(filePath)
		messageList = append(messageList, messageFile...)
	}

	return
}

func (b *Bundle) loadMessageFile(filePath string, format i18nFileType) (messageList []*Message, err error) {
	buf, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	messageFile, err := parseMessageFileBytes(buf, getUnmarshalFunc(format))
	if err != nil {
		return nil, err
	}

	b.addPath(filePath)
	messageList = append(messageList, messageFile...)

	return
}

func (b *Bundle) init() {
	b.tag = []language.Tag{}
	b.messageTemplates = make(map[language.Tag]map[string]*Message)
	b.localizerMap = make(map[string]*Localizer)
	for tag, fileMap := range b.fileMap {
		for fileType, filePathList := range fileMap {
			for _, filePath := range filePathList {
				b.loadMessage(filePath, tag, fileType)
			}
		}
		b.localizerMap[tag.String()] = NewLocalizer(tag, b) // 翻译器
	}
}

// group：表名；content：字段名
// group 为空则只使用content。允许使用{#content} 模式，将只转换定位符{#}中的内容
func (b *Bundle) ShouldT(ctx context.Context, group string, content string) (string, bool) {
	b.mu.Lock()
	// 借用 b.path 来判断是否初始化
	if len(b.path) == 0 {
		b.init()
	}
	b.mu.Unlock()

	language := b.defaultLanguage.String()
	if lang := LanguageFromCtx(ctx); lang != "" {
		language = lang
	}

	if localizer, ok := b.localizerMap[language]; ok {
		var result string
		subContentList := findContent(content)
		if len(subContentList) == 0 {
			result, ok = localizer.LocalizeMessage(getTranslateText(group, content))
			if !ok {
				return content, ok
			}
		} else {
			for _, subContent := range subContentList {
				subContentStr := subContent[0]
				subResult, ok := localizer.LocalizeMessage(getTranslateText(group, subContentStr[2:len(subContentStr)-1]))
				if !ok {
					continue
				}

				result = strings.ReplaceAll(content, subContentStr, subResult)
			}

		}
		return result, true
	} else {
		return content, false
	}
}

func (b *Bundle) T(ctx context.Context, group string, content string) string {
	result, _ := b.ShouldT(ctx, group, content)
	return result
}
