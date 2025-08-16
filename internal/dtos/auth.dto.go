package dtos

type AuthResponseDTO struct {
	AccessToken  string           `json:"access_token"`
	RefreshToken string           `json:"refresh_token"`
	User         *UserResponseDTO `json:"user"`
	ExpiresIn    int64            `json:"expires_in"`
}
