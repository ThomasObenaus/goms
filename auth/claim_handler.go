package auth

import (
	"fmt"
	"sort"
)

type ClaimHandler func(claims Claims) error

var NoClaimHandler = func(claims Claims) error {
	return nil
}

var HasRealmRole = func(role string) ClaimHandler {
	return func(claims Claims) error {

		if sort.SearchStrings(claims.RealmAccess.Roles, role) >= len(claims.RealmAccess.Roles) {

			return fmt.Errorf("Role %s not given", role)
		}
		return nil
	}
}
