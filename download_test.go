package fundconnext_test

import (
	"fmt"
	"os"
	"testing"

	f "github.com/robowealth/fundconnext"
)

type FundProfile struct {
	FundCode         string
	AMCCode          string
	FundNameTH       string
	FundNameEN       string
	FundPolicy       string
	TaxType          string
	FIFFlag          string
	DividendFlag     string
	RegistrationDate string
	FundRiskLevel    string
	FXRiskFlag       string
	FATCAAllowFlag   string
}

func TestDownload(t *testing.T) {
	fc := &f.FundConnext{
		Username: os.Getenv("USERNAME"),
		Password: os.Getenv("PASSWORD"),
		Env:      os.Getenv("ENV"),
	}
	d := []FundProfile{}
	_, err := fc.Login().Download("20190103", f.FundProfile).Save("./test/20190103_fund_profile.zip").Extract("./test/fund_profile").One().Load(&d)
	if err != nil {
		panic(err)
	}
	fmt.Println(d)
}
