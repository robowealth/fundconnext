package fundconnext_test

import (
	"os"
	"testing"

	"github.com/robowealth/fundconnext"
	"gotest.tools/assert"
)

func TestLoginSuccess(t *testing.T) {
	Fc := &fundconnext.FundConnext{
		Username: os.Getenv("USERNAME"),
		Password: os.Getenv("PASSWORD"),
		Env:      "demo",
	}
	if err := Fc.Login().Error; err != nil {
		panic(err)
	}
	assert.Equal(t, Fc.AccessToken != "", true)
}

func TestLoginFail(t *testing.T) {
	Fc := &fundconnext.FundConnext{
		Username: "foo",
		Password: "bar",
		Env:      "demo",
	}
	if err := Fc.Login().Error; err != nil {
		assert.ErrorContains(t, err, "E000 Unauthorized access")
	}
}
