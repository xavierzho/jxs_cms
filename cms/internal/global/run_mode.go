package global

import (
	"github.com/gin-gonic/gin"
)

type RunMode string

func (r RunMode) String() string {
	return string(r)
}

const (
	RUN_MODE_RELEASE RunMode = gin.ReleaseMode
	RUN_MODE_DEBUG   RunMode = gin.DebugMode
	RUN_MODE_MIGRATE RunMode = "migrate"
)
