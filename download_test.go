package fundconnext_test

import (
	"os"
	"testing"

	f "github.com/robowealth/fundconnext"
)

func TestDownload(t *testing.T) {
	fc := &f.FundConnext{
		Username: os.Getenv("DEMO_USERNAME"),
		Password: os.Getenv("DEMO_PASSWORD"),
		Env:      "demo",
	}
	err := fc.Login().Download("20190103", f.FundProfileFileName).Save("./20190103_fund.zip").End()
	if err != nil {
		panic(err)
	}
}
