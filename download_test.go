package fundconnext_test

import (
	"fmt"
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
	f := fc.Login().Download("20190103", f.FundProfile).Save("./20190103_fund_profile.zip")
	defer f.End()
	if f.Error != nil {
		panic(f.Error)
	}
	b, _ := f.Hash()
	fmt.Printf("%x", b)
	err := f.Extract("./sal").Error
	if err != nil {
		panic(err)
	}
}
