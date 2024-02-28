package main

import (
	"fmt"
	"halodeksik-be/app/api"
	"halodeksik-be/app/appcloud"
	"halodeksik-be/app/appconfig"
	"halodeksik-be/app/appdb"
	"halodeksik-be/app/appencoder"
	"halodeksik-be/app/applogger"
	"halodeksik-be/app/appvalidator"
	"halodeksik-be/app/ws"
	"os"
)

func main() {
	err := appconfig.LoadConfig()
	if err != nil {
		applogger.Log.Error(err)
		return
	}

	logger, logFile := applogger.NewLogger()
	if logFile != nil {
		defer logFile.Close()
	}
	applogger.SetLogger(logger)

	pwd, err := os.Getwd()
	if err != nil {
		applogger.Log.Error(err)
	}
	tmpdir := fmt.Sprintf("%s/%s", pwd, appconfig.Config.Tmpdir)
	_, err = os.Stat(tmpdir)
	if os.IsNotExist(err) {
		err = os.Mkdir(tmpdir, os.ModePerm)
		if err != nil {
			applogger.Log.Errorf("failed to create temp directory: %v", err)
			return
		}
	}
	err = os.Setenv(appconfig.Config.LinuxEnvTmpDir, tmpdir)
	if err != nil {
		applogger.Log.Error(err)
	}

	jsonEncoder := appencoder.NewAppJsonEncoderImpl()
	appencoder.SetAppJsonEncoder(jsonEncoder)

	validator := appvalidator.NewAppValidatorImpl()
	if err := appvalidator.AddCustomValidators(validator); err != nil {
		applogger.Log.Error(err)
	}
	appvalidator.SetValidator(validator)

	fileUploader := appcloud.NewFileUploaderImpl()
	appcloud.SetAppFileUploader(fileUploader)

	db, dbErr := appdb.Connect()
	if dbErr != nil {
		applogger.Log.Error(dbErr)
	}

	allRepositories := api.InitializeRepositories(db)
	allUtil := api.InitializeUtil()
	allUseCases := api.InitializeUseCases(allRepositories, allUtil)
	hub := ws.NewHub()
	routerOpts := api.InitializeAllRouterOpts(allUseCases, hub)

	err = allUseCases.CronUseCase.StartCron()
	if err != nil {
		applogger.Log.Error(err)
	}

	go hub.Run()
	ginMode := api.GetGinMode()
	router := api.NewRouter(routerOpts, ginMode)
	server := api.NewServer(router)
	api.StartAndSetupGracefulShutdown(server)
}
