package middleware

import (
	constant "WalletTopUp/pkg/const"
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetTraceId() gin.HandlerFunc {
	return func(c *gin.Context) {
		tid := c.GetHeader(constant.TraceIdHeader)
		if tid == "" {
			tid = uuid.NewString()
		}
		c.Set(constant.TraceIdKey, tid)
		ctx := context.WithValue(c.Request.Context(), constant.TraceIdKey, tid)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
