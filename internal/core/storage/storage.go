package storage

import (
	"context"
	"io"
	"os"

	"github.com/sirupsen/logrus"

	"cloud.google.com/go/storage"
)

var (
	client = &storage.Client{}
	ctx    = context.Background()
)

// NewClient a new Google Cloud Storage client.
func NewClient() {
	var err error
	client, err = storage.NewClient(ctx)
	if err != nil {
		panic(err)
	}
}

// UploadToGCloudStorage for upload file to gcloud storage.
func UploadToGCloudStorage(bucket, filePath, fileName string, isPublic bool) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		logrus.Error("[UploadToGCloudStorage] Open file error: ", err)
		return "", err
	}

	defer func() {
		err = file.Close()
	}()

	obj := client.Bucket(bucket).Object(fileName)
	writer := obj.NewWriter(ctx)
	if _, err = io.Copy(writer, file); err != nil {
		logrus.Error("[UploadToGCloudStorage] Copy file error: ", err)
		return "", err
	}

	if err = writer.Close(); err != nil {
		logrus.Error("[UploadToGCloudStorage] close writer error: ", err)
		return "", err
	}

	if isPublic {
		acl := client.Bucket(bucket).Object(fileName).ACL()
		if err := acl.Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
			logrus.Error("[UploadToGCloudStorage] ACL set error: ", err)
			return "", err
		}
	}

	return fileName, nil
}
