package handler

import "github.com/gin-gonic/gin"

type WalletHandler interface {
	VerifyTx(c *gin.Context)
}
