package setting

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

const NESTED_SEPARATOR = "." // 显式地设置分割符

type Config struct {
	vp *viper.Viper
}

var sections = make(map[string]interface{})

func NewSetting(configs ...string) (*Config, error) {
	opt := viper.KeyDelimiter(NESTED_SEPARATOR)
	vp := viper.NewWithOptions(opt)
	vp.SetConfigName("config")
	for _, config := range configs {
		if config != "" {
			vp.AddConfigPath(config)
		}
	}
	vp.SetConfigType("yaml")
	err := vp.ReadInConfig()
	if err != nil {
		return nil, err
	}

	s := &Config{vp}
	return s, nil
}

func (s *Config) WatchSettingChange(fn func() error, logFn func(logStr string)) {
	go func() {
		s.vp.WatchConfig()
		s.vp.OnConfigChange(func(in fsnotify.Event) {
			err := s.ReloadAllSection()
			if err != nil {
				logFn(fmt.Sprintf("setting.OnConfigChange error: %s", err.Error()))
			}
			err = fn()
			if err != nil {
				logFn(fmt.Sprintf("setting.OnConfigChange error: %s", err.Error()))
			}
		})
	}()
}

func (s *Config) ReadSection(k string, v interface{}) (err error) {
	err = s.vp.UnmarshalKey(k, v)
	if err != nil {
		return err
	}

	if _, ok := sections[k]; !ok {
		sections[k] = v
	}
	return nil
}

func (s *Config) ReloadAllSection() (err error) {
	for k, v := range sections {
		err = s.ReadSection(k, v)
		if err != nil {
			return err
		}
	}

	return nil
}
