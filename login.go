package fundconnext

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

const (
	fundconnext = "https://stage.fundconnext.com/api/auth"
)

// LoginRequest structure
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse structure
type LoginResponse struct {
	Username      string `json:"username"`
	AccessToken   string `json:"access_token"`
	SACode        string `json:"saCode"`
	IsPassThrough bool   `json:"isPassthrough"`
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
	fundconnextPath, err := Endpoint(f.Env, "/api/auth")
	if err != nil {
		f.Error = err
		return f
	}
	resp, err := http.Post(fundconnextPath, "application/json", bytes.NewBuffer(requestString))
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
