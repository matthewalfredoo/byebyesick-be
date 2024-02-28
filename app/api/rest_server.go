package api

import (
	"context"
	"errors"
	"fmt"
	"halodeksik-be/app/appconfig"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/applogger"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func NewServer(router *gin.Engine) *http.Server {
	uri := appconfig.Config.AppUri
	port := appconfig.Config.AppRestPort
	addr := fmt.Sprintf("%s:%s", uri, port)

	return &http.Server{
		Addr:    addr,
		Handler: router,
	}
}

func StartAndSetupGracefulShutdown(server *http.Server) {
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			applogger.Log.Infof("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	applogger.Log.Info("Shutting down server...")

	serverTimeout, err := strconv.Atoi(appconfig.Config.ServerShutdownTimeout)
	if err != nil {
		serverTimeout = appconstant.DefaultServerShutdownTimeout
	}
	shutdownTimeoutDuration := time.Duration(serverTimeout * 1e9)

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeoutDuration)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		applogger.Log.Info("Server forced to shutdown:", err)
	}
	applogger.Log.Info("Server exiting")
}
