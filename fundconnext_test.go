package fundconnext

import (
	"os"
	"testing"

	"gotest.tools/assert"
)

func TestLogin(t *testing.T) {
	Username, Password := os.Getenv("USERNAME"), os.Getenv("PASSWORD")
	fc := (&FundConnext{
		Username: Username,
		Password: Password,
		Env:      "demo",
	})

	if err := fc.Login().Error; err != nil {
		panic(err)
	}
	assert.Equal(t, fc.AccessToken != "", true)
}
