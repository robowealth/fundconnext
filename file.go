package fundconnext

import (
	"archive/zip"
	"bytes"
	"crypto/md5"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	path "path/filepath"
	"strconv"
	"strings"
)

// FileReader is
type FileReader struct {
	Error              error
	Read               *io.ReadCloser
	ContentDisposition string
	ContentType        string
}

// File structure
type File struct {
	Error       error
	Name        string
	ContentType string
	Location    string
	Length      int64
}

// Files structure
type Files struct {
	Error    error
	Location string
	Files    []string
}

// TextFileMeta structure
type TextFileMeta struct {
	Date    string
	Count   int
	Version string
}

// One File
func (d *Files) One() (f *File) {
	defer func() {
		if p := recover(); p != nil {
			e, _ := p.(error)
			f = &File{
				Error: e,
			}
		}
	}()
	if d.Error != nil {
		panic(d.Error)
	}
	if len(d.Files) == 0 {
		return nil
	}
	file := d.Files[0]
	fi, err := os.Stat(file)
	if err != nil {
		panic(err)
	}
	ct, err := getFileContentTypeFromSrc(file)
	if err != nil {
		panic(err)
	}
	return &File{
		Error:       nil,
		Name:        fi.Name(),
		Location:    file,
		ContentType: ct,
	}
}

// Load is
func (d *File) Load(t interface{}) (m *TextFileMeta, er error) {
	defer func() {
		if p := recover(); p != nil {
			e, _ := p.(error)
			er = e
		}
	}()
	dat, err := ioutil.ReadFile(d.Location)
	if err != nil {
		panic(err)
	}
	s := bytes.NewBuffer(dat).String()
	var date, count, version string
	for k, v := range strings.Split(s, "\n") {
		if k == 0 {
			row := strings.Split(v, "|")
			date, count, version = row[0], row[1], row[2]
		} else {

		}
	}
	c, err := strconv.Atoi(count)
	if err != nil {
		panic(err)
	}
	return &TextFileMeta{
		Date:    date,
		Count:   c,
		Version: version,
	}, nil
}

// Extract is
func (d *File) Extract(dst string) (f *Files) {
	defer func() {
		if p := recover(); p != nil {
			e, _ := p.(error)
			f = &Files{
				Error: e,
			}
		}
	}()
	r, err := zip.OpenReader(d.Location)
	if err != nil {
		panic(err)
	}

	rs := make([]string, len(r.Reader.File))
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		os.MkdirAll(dst, 0755)
	}
	for i, file := range r.Reader.File {
		zippedFile, err := file.Open()
		if err != nil {
			panic(err)
		}
		defer zippedFile.Close()
		extractedFilePath := filepath.Join(
			dst,
			file.Name,
		)
		rs[i], err = path.Abs(extractedFilePath)
		if err != nil {
			panic(err)
		}
		if file.FileInfo().IsDir() {
			os.MkdirAll(extractedFilePath, file.Mode())
		} else {
			outputFile, err := os.OpenFile(
				extractedFilePath,
				os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
				file.Mode(),
			)
			if err != nil {
				panic(err)
			}
			defer outputFile.Close()
			_, err = io.Copy(outputFile, zippedFile)
			if err != nil {
				panic(err)
			}
		}
	}
	absDst, err := path.Abs(filepath.Join(dst))
	if err != nil {
		panic(err)
	}
	return &Files{
		Location: absDst,
		Files:    rs,
	}
}

// Hash is
func (d *File) Hash() ([]byte, error) {
	if d.Error != nil {
		return nil, d.Error
	}
	f, err := os.Open(d.Location)
	if err != nil {
		d.Error = err
		return nil, err
	}
	defer f.Close()
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		d.Error = err
		return nil, err
	}

	return h.Sum(nil), nil
}
