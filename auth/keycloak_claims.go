package auth

import "gopkg.in/square/go-jose.v2/jwt"

type Claims struct {
	jwt.Claims

	Scope             string         `json:"scope,omitempty"`
	PreferredUsername string         `json:"preferred_username,omitempty"`
	RealmAccess       realmAccess    `json:"realm_access,omitempty"`
	ResourceAccess    resourceAccess `json:"resource_access,omitempty"`
}
type realmAccess struct {
	Roles []string `json:"roles,omitempty"`
}

type resourceAccess struct {
	Account account `json:"account,omitempty"`
}
type account struct {
	Roles []string `json:"roles,omitempty"`
}
