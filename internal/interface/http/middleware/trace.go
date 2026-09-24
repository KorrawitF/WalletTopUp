package middleware

import (
	constant "WalletTopUp/pkg/const"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const TraceIdHeader = "X-Trace-Id"

func GetTraceId() gin.HandlerFunc {
	return func(c *gin.Context) {
		tid := c.GetHeader(TraceIdHeader)
		if tid == "" {
			tid = uuid.NewString()
		}
		c.Set(constant.TraceIdKey, tid)
		c.Next()
	}
}
