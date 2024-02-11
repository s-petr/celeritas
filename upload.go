package celeritas

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/gabriel-vasile/mimetype"
	"github.com/s-petr/celeritas/filesystems"
)

func (c *Celeritas) UploadFile(r *http.Request, destination, field string, fs filesystems.FS) error {
	fileName, err := c.getFileToUpload(r, field)
	if err != nil {
		c.ErrorLog.Println(err)
		return err
	}

	if fs != nil {
		if err := fs.Put(fileName, destination); err != nil {
			c.ErrorLog.Println(err)
			return err
		}
	} else {
		if err := os.Rename(fileName, fmt.Sprintf("%s/%s",
			destination, path.Base(fileName))); err != nil {
			c.ErrorLog.Println(err)
			return err
		}
	}

	defer func() {
		_ = os.Remove(fileName)
	}()

	return nil
}

func (c *Celeritas) getFileToUpload(r *http.Request, fieldName string) (string, error) {
	_ = r.ParseMultipartForm(c.config.upload.maxUploadSize)

	file, header, err := r.FormFile(fieldName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	mimeType, err := mimetype.DetectReader(file)

	if _, err := file.Seek(0, 0); err != nil {
		return "", err
	}

	if !includes(c.config.upload.allowedMimeTypes, mimeType.String()) {
		return "", errors.New(fmt.Sprintf("invalid file type (%s)",
			mimeType.String()))
	}

	dst, err := os.Create(fmt.Sprintf("./tmp/%s", header.Filename))
	if err != nil {
		return "", err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("./tmp/%s", header.Filename), nil
}

func includes(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
