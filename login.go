package fundconnext

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
