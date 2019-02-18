package fundconnext

import (
	"encoding/json"
	"errors"
	"hash"
	"io"
	"net/http"
	"os"
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
	Error    error
	FileType string
	Reader   *io.ReadCloser
}

// DataFile structure
type DataFile struct {
	Error error
}

// Extract is
func (d *DownloadedFile) Extract() *DownloadedFile {
	return d
}

// Hash is
func (d *DownloadedFile) Hash(H *hash.Hash) *DownloadedFile {
	return d
}

// Struct is
func (d *DownloadedFile) Struct(T interface{}) *DownloadedFile {
	return d
}

// Save filepath
func (d *DownloadedFile) Save(filepath string) *DownloadedFile {
	if d.Error != nil {
		return d
	}
	out, err := os.Create(filepath)
	if err != nil {
		d.Error = err
		return d
	}
	defer out.Close()
	if _, err = io.Copy(out, *d.Reader); err != nil {
		d.Error = err
		return d
	}
	return d
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
	// h := md5.New()
	// if _, err := io.Copy(h, resp.Body); err != nil {
	// 	return &DownloadedFile{
	// 		Error: err,
	// 	}
	// }
	return &DownloadedFile{
		Error:    nil,
		Reader:   &resp.Body,
		FileType: resp.Header.Get("Content-Type"),
	}
}
