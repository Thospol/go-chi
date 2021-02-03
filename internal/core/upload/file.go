package upload

import (
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"saaa-api/internal/core/config"

	"github.com/disintegration/imaging"
)

type UploadUtil interface {
	Download(url string) (*http.Response, error)
	SaveFile() (string, error)
	Resize(imageSize, thumbImageSize int) (*Image, error)
}

// File file
type File struct {
	cf      *config.Configs
	rr      *config.ReturnResult
	tempDir string
	src     io.Reader
}

// New new file
func New(src io.Reader) (*File, error) {
	temp, err := ioutil.TempDir(os.TempDir(), "workkami--")
	if err != nil {
		return nil, err
	}
	return &File{
		cf:      config.CF,
		rr:      config.RR,
		tempDir: temp,
		src:     src,
	}, nil
}

// Close for end of save file
func (file *File) Close() error {
	return os.RemoveAll(file.tempDir)
}

// Download for download file from url
func (file *File) Download(url string) (*http.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// SaveFile file after upload from client.
func (file *File) SaveFile() (string, error) {
	err := file.saveFile("original")
	if err != nil {
		return "", err
	}
	filePath := fmt.Sprintf("%s/original", file.tempDir)

	return filePath, nil
}

// Image filename for resize images
type Image struct {
	OriginalImageFilename string
	ThumbImageFilename    string
	TempDir               string
	Size                  image.Point
}

// Resize new model resize
func (file *File) Resize(imageSize, thumbImageSize int) (*Image, error) {
	err := file.saveFile("original")
	if err != nil {
		return nil, err
	}

	originalPath := fmt.Sprintf("%s/original", file.tempDir)
	resize, err := file.resizeImage(imageSize, thumbImageSize, originalPath)
	if err != nil {
		return nil, err
	}

	return resize, nil
}

// saveFile save file
func (file *File) saveFile(filename string) error {
	if file.src == nil {
		return file.rr.FileNotFound
	}

	filePath := fmt.Sprintf("%s/%s", file.tempDir, filename)

	f, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer func() {
		err = f.Close()
	}()

	_, err = io.Copy(f, file.src)
	if err != nil {
		return err
	}

	return nil
}

// saveImage save image
func saveImage(src image.Image, tempDir, filename string) error {
	f, err := os.Create(fmt.Sprintf("%s/%s", tempDir, filename))
	if err != nil {
		return err
	}

	defer func() {
		err = f.Close()
	}()

	err = jpeg.Encode(f, src, &jpeg.Options{Quality: 100})
	if err != nil {
		return err
	}

	return nil
}

// resizeImage resize image
func (file *File) resizeImage(imageSize, thumbImageSize int, originalFilePath string) (*Image, error) {
	originalImage, err := imaging.Open(originalFilePath)
	if err != nil {
		return nil, err
	}

	original, thumb := "original.jpg", "thumb.jpg"
	originalResizeImage := file.image(originalImage, imageSize)
	err = saveImage(originalResizeImage, file.tempDir, original)
	if err != nil {
		return nil, err
	}

	thumbResizeImage := file.image(originalImage, thumbImageSize)
	err = saveImage(thumbResizeImage, file.tempDir, thumb)
	if err != nil {
		return nil, err
	}

	return &Image{
		OriginalImageFilename: original,
		ThumbImageFilename:    thumb,
		TempDir:               file.tempDir,
		Size:                  originalImage.Bounds().Size(),
	}, nil
}

// Image resize image
func (file *File) image(originalImage image.Image, maxWidth int) *image.NRGBA {
	if originalImage.Bounds().Size().X > maxWidth {
		return imaging.Resize(originalImage, maxWidth, 0, imaging.Lanczos)
	}

	return imaging.Resize(originalImage, originalImage.Bounds().Size().X, 0, imaging.Lanczos)
}
