package global

import (
	"time"

	"data_backend/pkg/util"
)

type ServerConfig struct {
	ServerName   string
	Version      string
	RunMode      RunMode
	HttpPort     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	TimeZone     string
	Language     string
	WhiteList    []string
}

type APPConfig struct {
	DefaultPageSize       int
	MaxPageSize           int
	DefaultContextTimeout time.Duration
}

var (
	ServerSetting *ServerConfig
	APPSetting    *APPConfig
	TelSetting    *util.TelConfig
)
