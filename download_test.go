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
	fundProfile := fc.Login().Download("20190103", f.FundProfile).Save("./20190103_fund_profile.zip")
	defer fundProfile.End()
	if fundProfile.Error != nil {
		panic(fundProfile.Error)
	}
	println(fundProfile.Length)
	println(fundProfile.FileType)
}
