package fundconnext

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

// FundConnext is
type FundConnext struct {
	Username    string
	Password    string
	Error       error
	AccessToken string
}

// Login is
func (f *FundConnext) Login() *FundConnext {
	if f.Error != nil {
		return f
	}
	requestString, err := json.Marshal(&LoginRequest{Username: f.Username, Password: f.Password})
	if err != nil {
		f.Error = err
		return f
	}
	resp, err := http.Post(fundconnext, "application/json", bytes.NewBuffer(requestString))
	if err != nil {
		f.Error = err
		return f
	}

	decoder := json.NewDecoder(resp.Body)

	if resp.StatusCode != 200 {
		var errorResponse map[string]map[string]string
		decoder.Decode(&errorResponse)

		message, code := errorResponse["errMsg"]["message"], errorResponse["errMsg"]["code"]

		f.Error = errors.New(code + " " + message)
		return f
	}

	var result *LoginResponse
	decoder.Decode(&result)
	f.AccessToken = result.AccessToken
	return f
}
