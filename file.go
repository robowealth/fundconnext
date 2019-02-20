package fundconnext

import (
	"archive/zip"
	"bytes"
	"crypto/md5"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	path "path/filepath"
	"reflect"
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

// All is
func (d *File) All(t interface{}) (m *TextFileMeta, er error) {
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
	lines := strings.Split(s, "\n")
	dest := reflect.ValueOf(t)
	if dest.Kind() != reflect.Ptr {
		panic(errors.New("some: check must be a pointer"))
	}
	item := dest.Elem()

	if item.Kind() != reflect.Slice {
		panic(errors.New("Input interface is not a slice"))
	}
	ElemType := item.Type()
	stc := ElemType.Elem()
	if stc.Kind() != reflect.Struct {
		panic(errors.New("Input interface is not a slice of struct/array"))
	}
	item.Set(reflect.MakeSlice(ElemType, len(lines)-1, len(lines)-1))
	var c int
	for k, v := range lines {
		if k == 0 {
			row := strings.Split(v, "|")
			date, count, version = row[0], row[1], row[2]
			c, err = strconv.Atoi(count)
			if err != nil {
				panic(err)
			}
			if len(lines)-1 != c {
				panic(errors.New("Data Length is invalid"))
			}
		} else {
			row := strings.Split(v, "|")
			size := stc.NumField()
			el := reflect.New(stc).Elem()
			for i := 0; i < size; i++ {
				typeField := stc.Field(i)
				el.FieldByName(typeField.Name).SetString(row[i])
			}
			item.Index(k - 1).Set(el)
		}
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
