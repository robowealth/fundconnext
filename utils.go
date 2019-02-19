package fundconnext

import (
	"net/http"
	"os"
)

func getFileContentTypeFromSrc(src string) (string, error) {
	f, err := os.Open(src)
	if err != nil {
		return "", err
	}
	defer f.Close()

	contentType, err := getFileContentType(f)
	if err != nil {
		return "", err
	}
	return contentType, nil
}

func getFileContentType(out *os.File) (string, error) {

	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}
