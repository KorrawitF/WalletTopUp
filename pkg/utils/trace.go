package utils

import (
	constant "WalletTopUp/pkg/const"
	"context"
)

func GetTraceId(ctx context.Context) string {
	val, ok := ctx.Value(constant.TraceIdKey).(string)
	if !ok {
		val = ""
	}
	return val
}
