package middleware

import (
	"fmt"
	"halodeksik-be/app/applogger"
	"time"

	"github.com/gin-gonic/gin"
)

func LogHandler(ctx *gin.Context) {
	if gin.Mode() == gin.TestMode {
		ctx.Next()
		return
	}

	args := make(map[string]interface{})
	args["client_ip"] = ctx.ClientIP()
	args["type"] = "REQUEST REST"
	args["method"] = ctx.Request.Method
	args["uri"] = ctx.Request.RequestURI

	applogger.Log.WithFields(args).Info("request received")

	startingTime := time.Now()

	ctx.Next()

	args["type"] = "RESPONSE REST"
	args["latency"] = fmt.Sprintf("%d%s", time.Since(startingTime).Microseconds(), "μs")
	args["response_status"] = ctx.Writer.Status()

	err := ctx.Errors.Last()
	if err == nil {
		applogger.Log.WithFields(args).Info("response sent")
		return
	}

	applogger.Log.WithFields(args).Error(err)
}
