package main

import (
	"fmt"
)

// Token ...
type Token struct {
	AccessToken string `json:"access_token"`
	// Token lifetime in seconds
	ExpiresIn int    `json:"expires_in"`
	TokenType string `json:"token_type"`
	Scope     string `json:"scope"`
}

// AsString ...
func (t *Token) AsString() string {
	return fmt.Sprintf("Token: %s", t.AccessToken)
}
