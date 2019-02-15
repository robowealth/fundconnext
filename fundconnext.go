package fundconnext

import (
	"bytes"
	"encoding/json"
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
	var result interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	if resp.StatusCode != 200 {
		f.Error = error
		return f
	}

	f.AccessToken = "" //result.AccessToken
	return f
}
