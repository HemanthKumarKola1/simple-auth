package models

type Token struct {
	TokenType string `json:"token_type"`
	Revoked   bool   `json:"revoked"`
}
