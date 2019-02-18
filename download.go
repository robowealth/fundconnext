package fundconnext

import (
	"archive/zip"
	"crypto/md5"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	path "path/filepath"
)

const (
	// FundProfile  T
	FundProfile string = "FundProfile.zip"
	// FundMapping  T
	FundMapping string = "FundMapping.zip"
	// SwitchingMatrix  T
	SwitchingMatrix string = "SwitchingMatrix.zip"
	// FundHoliday  T
	FundHoliday string = "FundHoliday.zip"
	// TradeCalendar  T
	TradeCalendar string = "TradeCalendar.zip"
	// AccountProfile  T
	AccountProfile string = "AccountProfile.zip"
	// UnitholderMapping  T
	UnitholderMapping string = "UnitholderMapping.zip"
	// BankAccountUnitholder  T
	BankAccountUnitholder string = "BankAccountUnitholder.zip"
	// CustomerProfile  T
	CustomerProfile string = "CustomerProfile.zip"
	// NAV  T-1
	NAV string = "Nav.zip"
	// UnitholderBalance  T-1
	UnitholderBalance string = "UnitholderBalance.zip"
	// AllottedTransactions  T-1
	AllottedTransactions string = "AllottedTransactions.zip"
	// DividendNews   T-1
	DividendNews string = "DividendNews.zip"
	// DividendTransactions  T-1
	DividendTransactions string = "DividendTransactions.zip"
)

// DownloadedFile is
type DownloadedFile struct {
	Error        error
	FileType     string
	Reader       *io.ReadCloser
	Length       int64
	UnCompressed bool
}

// SavedFile structure
type SavedFile struct {
	DownloadedFile
	Location string
}

// DataFile structure
type DataFile struct {
	SavedFile
}

// Extract is
func (d *SavedFile) Extract(dst string) *DataFile {
	r, err := zip.OpenReader(d.Location)
	if err != nil {
		d.Error = err
		return &DataFile{
			SavedFile: *d,
		}
	}

	rs := make([]string, len(r.Reader.File))
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		os.MkdirAll(dst, 0755)
	}
	for i, file := range r.Reader.File {
		zippedFile, err := file.Open()
		if err != nil {
			d.Error = err
			return &DataFile{
				SavedFile: *d,
			}
		}
		defer zippedFile.Close()
		extractedFilePath := filepath.Join(
			dst,
			file.Name,
		)
		rs[i] = extractedFilePath
		if file.FileInfo().IsDir() {
			os.MkdirAll(extractedFilePath, file.Mode())
		} else {
			outputFile, err := os.OpenFile(
				extractedFilePath,
				os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
				file.Mode(),
			)
			if err != nil {
				d.Error = err
				return &DataFile{
					SavedFile: *d,
				}
			}
			defer outputFile.Close()
			_, err = io.Copy(outputFile, zippedFile)
			if err != nil {
				d.Error = err
				return &DataFile{
					SavedFile: *d,
				}
			}
		}
	}
	d.Location = rs[0]
	return &DataFile{
		SavedFile: *d,
	}
}

// Hash is
func (d *SavedFile) Hash() ([]byte, error) {
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

// Struct is
func (d *DownloadedFile) Struct(T interface{}) *DownloadedFile {
	return d
}

// SetLocation path
func (d *DownloadedFile) SetLocation(filepath string) *SavedFile {
	abspath, err := path.Abs(filepath)
	if err != nil {
		d.Error = err
		return &SavedFile{
			DownloadedFile: *d,
		}
	}
	return &SavedFile{
		DownloadedFile: *d,
		Location:       abspath,
	}
}

// Save filepath
func (d *DownloadedFile) Save(filepath string) *SavedFile {
	if d.Error != nil {
		return &SavedFile{
			DownloadedFile: *d,
		}
	}
	out, err := os.Create(filepath)
	if err != nil {
		d.Error = err
		return &SavedFile{
			DownloadedFile: *d,
		}
	}
	defer out.Close()
	if _, err = io.Copy(out, *d.Reader); err != nil {
		d.Error = err
		return &SavedFile{
			DownloadedFile: *d,
		}
	}
	abspath, err := path.Abs(filepath)
	if err != nil {
		d.Error = err
		return &SavedFile{
			DownloadedFile: *d,
		}
	}
	return &SavedFile{
		DownloadedFile: *d,
		Location:       abspath,
	}
}

// End is
func (d *DownloadedFile) End() error {
	if d.Error != nil {
		return d.Error
	}
	defer (*d.Reader).Close()
	return nil
}

// Download is
func (f *FundConnext) Download(date, file string) *DownloadedFile {
	if f.Error != nil {
		return &DownloadedFile{
			Error: f.Error,
		}
	}
	fundconnextPath, err := endpoint(f.Env, "/api/files/"+date+"/"+file)
	if err != nil {
		return &DownloadedFile{
			Error: err,
		}
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", fundconnextPath, nil)
	if err != nil {
		return &DownloadedFile{
			Error: err,
		}
	}
	req.Header.Set("X-Auth-Token", f.AccessToken)
	resp, err := client.Do(req)
	if err != nil {
		return &DownloadedFile{
			Error: err,
		}
	}
	// Check Error
	if resp.StatusCode != 200 {
		var errorResponse map[string]map[string]string
		json.NewDecoder(resp.Body).Decode(&errorResponse)
		message, code := errorResponse["errMsg"]["message"], errorResponse["errMsg"]["code"]
		return &DownloadedFile{
			Error: errors.New(code + " " + message),
		}
	}

	return &DownloadedFile{
		Error:        nil,
		Reader:       &resp.Body,
		FileType:     resp.Header.Get("Content-Type"),
		Length:       resp.ContentLength,
		UnCompressed: resp.Uncompressed,
	}
}
