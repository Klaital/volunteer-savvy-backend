package users

import (
	"testing"
	"time"
)

func TestCreateJWT(t *testing.T) {
	user := User{
		Id: 1,
		Email: "kit@example.org",
		Guid: "kit@example.org",
		Roles: map[uint64][]Role{
			1: []Role{{
				OrgId: 1,
				UserId: 1,
				UserGuid: "kit@example.org",
			}},
		},
	}

	claims := CreateJWT(&user, 5 * time.Minute)
	if claims == nil {
		t.Error("Expected claims to not be nil")
		t.Fail()
	}

	if len(claims.Roles) != 1 {
		t.Errorf("Expected 1 org with roles on sample user. Got %d instead", len(claims.Roles))
	}
	if len(claims.Roles[1]) != 1 {
		t.Errorf("Expected 1 role on org 1. Got %d instead", len(claims.Roles[1]))
	}
}
