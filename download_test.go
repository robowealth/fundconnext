package fundconnext_test

import (
	"os"
	"testing"

	f "github.com/robowealth/fundconnext"
)

func TestDownload(t *testing.T) {
	fc := &f.FundConnext{
		Username: os.Getenv("USERNAME"),
		Password: os.Getenv("PASSWORD"),
		Env:      os.Getenv("ENV"),
	}
	err := fc.Login().Download("20190103", f.DividendNews).Save("./20190103_fund.zip").End()
	if err != nil {
		panic(err)
	}
}
