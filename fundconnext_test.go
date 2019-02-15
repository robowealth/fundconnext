package fundconnext

import (
	"testing"

	"gotest.tools/assert"
)

func TestLogin(t *testing.T) {
	fc := (&FundConnext{
		Username: "API_ROBO01",
		Password: "password3",
	})
	fc.Login()
	println(fc.AccessToken)
	assert.Equal(t, "A", "A")
}
