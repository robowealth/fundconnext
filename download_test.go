package fundconnext_test

import (
	"encoding/json"
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
	BuyCutOffTime    string
}

func TestDownload(t *testing.T) {
	fc := &f.FundConnext{
		Username: os.Getenv("USERNAME"),
		Password: os.Getenv("PASSWORD"),
		Env:      os.Getenv("ENV"),
	}
	data := []FundProfile{}
	_, err := fc.Login().Download("20190103", f.FundProfile).Save("./test/20190103_fund_profile.zip").Extract("./test/fund_profile").One().All(&data)
	if err != nil {
		panic(err)
	}
	for _, v := range data {
		b, _ := json.MarshalIndent(v, "", "  ")
		fmt.Println(string(b))
	}
}
