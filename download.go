package fundconnext

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
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

// Download is
func (f *FundConnext) Download(date, file string) (d *FileReader) {
	if f.Error != nil {
		return &FileReader{
			Error: f.Error,
		}
	}
	fundconnextPath, err := endpoint(f.Env, "/api/files/"+date+"/"+file)
	if err != nil {
		return &FileReader{
			Error: err,
		}
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", fundconnextPath, nil)
	if err != nil {
		return &FileReader{
			Error: err,
		}
	}
	req.Header.Set("X-Auth-Token", f.AccessToken)
	resp, err := client.Do(req)
	if err != nil {
		return &FileReader{
			Error: err,
		}
	}
	// Check Error
	if resp.StatusCode != 200 {
		var errorResponse map[string]map[string]string
		json.NewDecoder(resp.Body).Decode(&errorResponse)
		message, code := errorResponse["errMsg"]["message"], errorResponse["errMsg"]["code"]
		return &FileReader{
			Error: errors.New(code + " " + message),
		}
	}

	return &FileReader{
		Error:              nil,
		Read:               &resp.Body,
		ContentType:        resp.Header.Get("Content-Type"),
		ContentDisposition: resp.Header.Get("Content-Disposition"),
	}
}

// Save filepath
func (d *FileReader) Save(filepath string) (s *File) {
	defer func() {
		if p := recover(); p != nil {
			e, _ := p.(error)
			s = &File{
				Error: e,
			}
		}
	}()
	defer (*d.Read).Close()
	if d.Error != nil {
		panic(d.Error)
	}

	out, err := os.Create(filepath)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	if _, err = io.Copy(out, *d.Read); err != nil {
		panic(err)
	}
	abspath, err := path.Abs(filepath)
	if err != nil {
		panic(err)
	}
	fi, err := os.Stat(abspath)
	if err != nil {
		panic(err)
	}
	return &File{
		Error:       nil,
		Name:        fi.Name(),
		ContentType: d.ContentType,
		Location:    abspath,
		Length:      fi.Size(),
	}
}
