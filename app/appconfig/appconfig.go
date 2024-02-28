package appconfig

import (
	"halodeksik-be/app/env"
	"os"
)

var Config *AppConfig

type AppConfig struct {
	LinuxEnvTmpDir string

	DbHost     string
	DbPort     string
	DbUser     string
	DbPassword string
	DbName     string

	AppName     string
	AppUri      string
	AppRestPort string
	AppMode     string
	AppClient   string
	Tmpdir      string

	MailAddress  string
	MailSender   string
	MailPassword string
	MailSmtpHost string
	MailSmtpPort string

	FrontendUrl string

	RegisterTokenExpired string
	LoginTokenExpired    string
	ForgotTokenExpired   string

	GcloudCredentialFile                    string
	GcloudStorageProjectId                  string
	GcloudStorageBucketName                 string
	GcloudStorageUrl                        string
	GcloudStorageCdn                        string
	GcloudStorageFolderProducts             string
	GcloudStorageFolderCertificates         string
	GcloudStorageFolderProfiles             string
	GcloudStoragePaymentProofs              string
	GcloudStorageFolderConsultationSessions string
	GmapUrl                                 string
	GmapKey                                 string

	RajaongkirUrl string
	RajaongkirKey string

	JwtSecret string

	RequestTimeout        string
	ServerShutdownTimeout string
}

func LoadConfig() error {
	err := env.LoadEnv()
	if err != nil {
		return err
	}
	Config = &AppConfig{
		LinuxEnvTmpDir:                          "TMPDIR",
		DbHost:                                  os.Getenv("DB_HOST"),
		DbPort:                                  os.Getenv("DB_PORT"),
		DbUser:                                  os.Getenv("DB_USER"),
		DbPassword:                              os.Getenv("DB_PASSWORD"),
		DbName:                                  os.Getenv("DB_NAME"),
		AppName:                                 os.Getenv("APP_NAME"),
		AppUri:                                  os.Getenv("APP_URI"),
		AppRestPort:                             os.Getenv("APP_REST_PORT"),
		AppMode:                                 os.Getenv("APP_MODE"),
		AppClient:                               os.Getenv("APP_CLIENT"),
		Tmpdir:                                  os.Getenv("APP_TMPDIR"),
		MailAddress:                             os.Getenv("MAIL_EMAIL"),
		MailSender:                              os.Getenv("MAIL_SENDER"),
		MailPassword:                            os.Getenv("MAIL_PASSWORD"),
		MailSmtpHost:                            os.Getenv("MAIL_SMTP_HOST"),
		MailSmtpPort:                            os.Getenv("MAIL_SMTP_PORT"),
		FrontendUrl:                             os.Getenv("FRONTEND_URL"),
		RegisterTokenExpired:                    os.Getenv("REGISTER_TOKEN_EXPIRED_MINUTE"),
		LoginTokenExpired:                       os.Getenv("LOGIN_TOKEN_EXPIRED_MINUTE"),
		ForgotTokenExpired:                      os.Getenv("FORGOT_TOKEN_EXPIRED_MINUTE"),
		GcloudCredentialFile:                    os.Getenv("GCLOUD_CREDENTIAL_FILE"),
		GcloudStorageProjectId:                  os.Getenv("GCLOUD_STORAGE_PROJECT_ID"),
		GcloudStorageBucketName:                 os.Getenv("GCLOUD_STORAGE_BUCKET_NAME"),
		GcloudStorageUrl:                        os.Getenv("GCLOUD_STORAGE_URL"),
		GcloudStorageCdn:                        os.Getenv("GCLOUD_STORAGE_CDN"),
		GcloudStorageFolderProducts:             os.Getenv("GCLOUD_STORAGE_FOLDER_PRODUCTS"),
		GcloudStorageFolderCertificates:         os.Getenv("GCLOUD_STORAGE_FOLDER_CERTIFICATES"),
		GcloudStorageFolderProfiles:             os.Getenv("GCLOUD_STORAGE_FOLDER_PROFILES"),
		GcloudStoragePaymentProofs:              os.Getenv("GCLOUD_STORAGE_FOLDER_PAYMENT_PROOFS"),
		GcloudStorageFolderConsultationSessions: os.Getenv("GCLOUD_STORAGE_FOLDER_CONSULTATION_SESSIONS"),
		GmapUrl:                                 os.Getenv("GMAP_URL"),
		GmapKey:                                 os.Getenv("GMAP_KEY"),
		RajaongkirUrl:                           os.Getenv("RAJAONGKIR_URL"),
		RajaongkirKey:                           os.Getenv("RAJAONGKIR_API_KEY"),
		JwtSecret:                               os.Getenv("SECRET_JWT_KEY"),
		RequestTimeout:                          os.Getenv("REQUEST_TIMEOUT"),
		ServerShutdownTimeout:                   os.Getenv("SERVER_SHUTDOWN_TIMEOUT"),
	}
	return nil
}
