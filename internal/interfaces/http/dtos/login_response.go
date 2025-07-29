package dtos

type LoginResponse struct {
	Profile      *ProfileResponse
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}