package httpiface

import (
	"WalletTopUp/internal/interface/http/handler"
	"WalletTopUp/internal/interface/http/middleware"
	"WalletTopUp/pkg/lib"

	"github.com/gin-gonic/gin"
)

func NewRouter(logger lib.Logger, handlers handler.WalletHandler) *gin.Engine {
	r := gin.New()
	r.Use(middleware.GetTraceId())

	return r
}
