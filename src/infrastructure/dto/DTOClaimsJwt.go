package dto

type DTOClaimsJwt struct {
	Username   string `json:"username"`
	Email      string `json:"email"`
	IssuedAt   int64  `json:"firstname"`
	Expiration int64  `json:"expiration"`
}
