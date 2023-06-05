package helper

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

func GetFilename(name string) (string, error) {
	ext := strings.Split(name, ".")

	filename := strings.Replace(uuid.NewString(), "-", "", -1)

	switch ext[len(ext)-1] {
	case "jpeg":
		filename += ".jpeg"
	case "jpg":
		filename += ".jpg"
	case "png":
		filename += ".png"
	case "mp4":
		filename += ".mp4"
	case "mkv":
		filename += ".mkv"
	case "pdf":
		filename += ".pdf"
	default:
		return "", errors.New("not supported file")
	}

	return filename, nil
}
