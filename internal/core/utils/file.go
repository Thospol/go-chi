package utils

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"os/exec"

	"saaa-api/internal/core/unique"

	"github.com/sirupsen/logrus"
)

const (
	// FormKeySingleFile key single file
	FormKeySingleFile = "file"
	// FormKeyMultiFile key multi files
	FormKeyMultiFile = "files"
)

// UploadFile upload file
func UploadFile(r *http.Request, maxMemory int64) (multipart.File, *multipart.FileHeader, error) {
	if err := r.ParseMultipartForm(maxMemory); err != nil {
		logrus.Error("[UploadFile] parses a request body as multipart/form-data error: ", err)
		return nil, nil, err
	}

	return r.FormFile(FormKeySingleFile) // file
}

// UploadFiles upload file
func UploadFiles(r *http.Request, maxMemory int64) ([]*multipart.FileHeader, error) {
	if err := r.ParseMultipartForm(maxMemory); err != nil {
		logrus.Error("[UploadFile] parses a request body as multipart/form-data error: ", err)
		return nil, err
	}

	return r.MultipartForm.File[FormKeyMultiFile], nil // files
}

// ScreenShotVideo screen shot video
func ScreenShotVideo(filePath string) (string, error) {
	outputPath := fmt.Sprintf("%s.jpeg", unique.UUID())
	cmd := exec.Command("ffmpeg", "-ss", "00:00:01.000", "-i", filePath, "-vframes", "1", outputPath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logrus.Errorf("command execute failed with %s", err)
		return outputPath, err
	}
	logrus.Infof("command execute: %s", string(out))
	return outputPath, nil
}

// RemoveFileScreenShot remove file screen shot
func RemoveFileScreenShot(filePath string) error {
	cmd := exec.Command("rm", "-f", filePath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logrus.Errorf("command execute failed with %s", err)
		return err
	}
	logrus.Infof("command execute: %s", string(out))
	return nil
}
