package auth

type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"email,required"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"email,required"`
	Password string `json:"password" validate:"required"`
}

type SetBannedRequest struct {
	IsBanned *bool `json:"isBanned" validate:"required"`
}

type SetDeletedRequest struct {
	IsDeleted *bool `json:"isDeleted" validate:"required"`
}

type SetAdminRequest struct {
	IsAdmin *bool `json:"IsAdmin" validate:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RegisterResponse struct {
	UserId int64 `json:"userId"`
}

type RefreshResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type SetBannedResponse struct {
	IsBanned bool `json:"isBanned"`
}

type SetDeletedResponse struct {
	IsDeleted bool `json:"isDeleted"`
}

type SetAdminResponse struct {
	IsAdmin bool `json:"IsAdmin"`
}
