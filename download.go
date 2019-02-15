package fundconnext

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"hash"
	"io"
	"net/http"
	"os"
)

const (
	// FundProfileFileName filename
	FundProfileFileName string = "FundProfile.zip"
	// FundMappingFileName filename
	FundMappingFileName string = ""
	// FundSwitchingFileName filename
	FundSwitchingFileName string = ""
	// FundHolidayFileName filename
	FundHolidayFileName string = ""
	// TradeCalendarFileName filename
	TradeCalendarFileName string = ""
)

// DownloadedFile is
type DownloadedFile struct {
	Error    error
	Hash     hash.Hash
	FileType string
	Reader   *io.ReadCloser
}

// Extract file
func (d *DownloadedFile) Extract() {

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

	h := md5.New()
	if _, err := io.Copy(h, resp.Body); err != nil {
		return &DownloadedFile{
			Error: err,
		}
	}
	return &DownloadedFile{
		Error:    nil,
		Reader:   &resp.Body,
		Hash:     h,
		FileType: resp.Header.Get("Content-Type"),
	}
}
