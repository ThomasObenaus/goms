package auth

import (
	"fmt"
)

type ClaimHandler func(claims Claims) error

var NoClaimHandler = func(claims Claims) error {
	return nil
}

var HasRealmRole = func(role string) ClaimHandler {
	return func(claims Claims) error {
		for _, roleInToken := range claims.RealmAccess.Roles {
			if role == roleInToken {
				return nil
			}
		}
		return fmt.Errorf("Role %s not given", role)
	}
}
