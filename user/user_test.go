package user

import (
	"github.com/dgrijalva/jwt-go"
	"testing"
)

func TestUser_FillFromSource(t *testing.T) {
	u := &User{}
	source := map[string]interface{}{
		"a": "username-val",
		"b": "email-val",
		"c": "login-val",
		"r": "role-val",
		"i": 123,
	}
	mapping := map[string]interface{}{
		"a": "user",
		"b": "email",
		"c": "login",
		"r": "role",
		"i": "id",
	}
	user, err := u.FillFromSource(source, mapping)
	if err != nil {
		t.Fail()
	}
	if user.Login != "login-val" {
		t.Fail()
	}
	t.Log(u)
}

func TestUser_Fill(t *testing.T) {
	u := &User{}
	c := jwt.MapClaims{}
	c["uid"] = 123.0
	c["role"] = "role-val"
	c["user"] = "user-val"
	c["login"] = "login-val"
	c["email"] = "email-val"
	user := u.Fill(c)
	if user.Login != "login-val" {
		t.Fail()
	}
	t.Log(u)
}
