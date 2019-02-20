package security

import (
	"github.com/gin-gonic/gin"
	"github.com/nikonm/go-gin-auth-middleware/user"
	"net/http"
	"testing"
)

func initSecurity(t *testing.T) {
	options := Options{
		Secret:       "verysecret",
		PasswordAlgo: "sha1",
		Adapters:     make(map[string]map[string]interface{}),
		HeaderName:   "AUTH",
		TokenExp:     "30m",
	}
	err := Init(options)
	if err != nil {
		t.Fail()
	}
}

func TestSecurity_PwdHash(t *testing.T) {
	initSecurity(t)
	h := Security.PwdHash("test")
	if h != "a94a8fe5ccb19ba61c4c0873d391e987982fbbd3" {
		t.Fail()
	}
}

func TestSecurity_GetUser(t *testing.T) {
	initSecurity(t)
	u := &user.User{Id: 41, Login: "test", Email: "test@test.local", User: "Name", Role: "role"}

	token := Security.makeToken(u)

	headers := http.Header{}
	headers.Add("AUTH", token)
	r := &http.Request{Header: headers}
	c := &gin.Context{Request: r}
	securedUser, err := Security.GetUser(c)

	if err != nil {
		t.Fail()
	}
	if securedUser == nil {
		t.Fail()
	}
}
