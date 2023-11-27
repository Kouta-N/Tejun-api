package models

type AuthorizedUser struct {
	Role          string `json:"role"`
	Iss           string `json:"iss"`
	Azp           string `json:"azp"`
	Aud           string `json:"aud"`
	Sub           string `json:"sub"`
	HD            string `json:"hd"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	AtHash        string `json:"at_hash"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Locale        string `json:"locale"`
	Iat           int64  `json:"iat"`
	Exp           int64  `json:"exp"`
}

type User struct {
  Email         string `json:"email"`
	Name          string `json:"name"`
}