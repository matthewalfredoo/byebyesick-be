package appcloud

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/api/option"
	"halodeksik-be/app/appconfig"
	"halodeksik-be/app/applogger"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

var AppFileUploader FileUploader

type FileUploader interface {
	SendToBucket(ctx context.Context, file multipart.File, object, path string) error
	UploadFromFile(ctx context.Context, file *os.File, path, name string) (string, error)
	UploadFromFileHeader(ctx context.Context, fileHeader any, folderName string) (string, error)
}

func SetAppFileUploader(uploader FileUploader) {
	AppFileUploader = uploader
}

type FileUploaderImpl struct {
	client     *storage.Client
	projectId  string
	bucketName string
	cloudUrl   string
}

func NewFileUploaderImpl() *FileUploaderImpl {
	credentialFile := appconfig.Config.GcloudCredentialFile
	client, err := storage.NewClient(context.Background(), option.WithCredentialsFile(credentialFile))
	if err != nil {
		applogger.Log.Errorf("failed to create file uploader client: %v", err)
	}

	projectId := appconfig.Config.GcloudStorageProjectId
	bucketName := appconfig.Config.GcloudStorageBucketName
	cloudUrl := appconfig.Config.GcloudStorageCdn

	return &FileUploaderImpl{
		client:     client,
		projectId:  projectId,
		bucketName: bucketName,
		cloudUrl:   cloudUrl,
	}
}

func (f *FileUploaderImpl) SendToBucket(ctx context.Context, file multipart.File, path, name string) error {
	bucketObject := f.client.Bucket(f.bucketName).Object(path + name)
	wc := bucketObject.NewWriter(ctx)
	wc.ACL = []storage.ACLRule{
		{
			Entity: storage.AllUsers,
			Role:   storage.RoleReader,
		},
	}

	if _, err := io.Copy(wc, file); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}
	return nil
}

func (f *FileUploaderImpl) UploadFromFile(ctx context.Context, file *os.File, folderName, fileName string) (string, error) {
	bucketObject := f.client.Bucket(f.bucketName).Object(
		fmt.Sprintf("%s/%s", folderName, fileName),
	)
	wc := bucketObject.NewWriter(ctx)
	wc.ACL = []storage.ACLRule{
		{
			Entity: storage.AllUsers,
			Role:   storage.RoleReader,
		},
	}

	if _, err := io.Copy(wc, file); err != nil {
		return "", fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("Writer.Close: %v", err)
	}

	url := fmt.Sprintf("%s/%s/%s", appconfig.Config.GcloudStorageCdn, folderName, fileName)
	return url, nil
}

func (f *FileUploaderImpl) UploadFromFileHeader(ctx context.Context, fileHeader any, folderName string) (string, error) {
	parsedFileHeader := fileHeader.(*multipart.FileHeader)
	file, err := parsedFileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	extension := filepath.Ext(parsedFileHeader.Filename)
	updateUUID, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	fileName := fmt.Sprintf("%s%s", updateUUID.String(), extension)

	err = f.SendToBucket(ctx, file, fmt.Sprintf("%s/", folderName), fileName)
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf("%s/%s/%s", f.cloudUrl, folderName, fileName)
	return url, nil
}
