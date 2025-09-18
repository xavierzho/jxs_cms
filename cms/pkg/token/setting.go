package token

import (
	"time"
)

type JWTConfig struct {
	Secret string
	Issuer string
	Expire time.Duration
}

var JWTSetting *JWTConfig
