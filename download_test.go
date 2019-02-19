package fundconnext_test

import (
	"fmt"
	"os"
	"testing"

	f "github.com/robowealth/fundconnext"
)

type FundProfile struct {
	FundCode string
	AMCCode  string
}

func TestDownload(t *testing.T) {
	fc := &f.FundConnext{
		Username: os.Getenv("USERNAME"),
		Password: os.Getenv("PASSWORD"),
		Env:      os.Getenv("ENV"),
	}
	d := &FundProfile{}
	fs, _ := fc.Login().Download("20190103", f.FundProfile).Save("./test/20190103_fund_profile.zip").Extract("./test/fund_profile").One().Load(d)
	fmt.Println(fs)
}
