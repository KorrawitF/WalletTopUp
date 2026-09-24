package middleware

import (
	constant "WalletTopUp/pkg/const"
	"WalletTopUp/pkg/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func runTrace(header string) (ginTid, ctxTid string) {
	r := gin.New()
	r.Use(GetTraceId())
	r.GET("/", func(c *gin.Context) {
		ginTid = c.GetString(constant.TraceIdKey)
		ctxTid = utils.GetTraceId(c.Request.Context())
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	if header != "" {
		req.Header.Set(constant.TraceIdHeader, header)
	}
	r.ServeHTTP(httptest.NewRecorder(), req)
	return ginTid, ctxTid
}

func TestGetTraceId_UsesHeader(t *testing.T) {
	ginTid, ctxTid := runTrace("trace-123")

	assert.Equal(t, "trace-123", ginTid)
	assert.Equal(t, "trace-123", ctxTid)
}

func TestGetTraceId_GeneratesWhenMissing(t *testing.T) {
	ginTid, ctxTid := runTrace("")

	_, err := uuid.Parse(ginTid)
	assert.NoError(t, err)
	assert.Equal(t, ginTid, ctxTid)
}
