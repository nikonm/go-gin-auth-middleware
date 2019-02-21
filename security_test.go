package security

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/nikonm/go-gin-auth-middleware/user"
	"net/http"
	"net/http/httptest"
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

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString("{\"foo\":\"bar\", \"bar\":\"foo\"}"))
	c.Request.Header.Add("AUTH", token)

	securedUser, err := Security.GetUser(c)

	if err != nil {
		t.Fail()
	}
	if securedUser == nil {
		t.Fail()
	}
}

func TestSecurity_CheckRoleAuth(t *testing.T) {
	initSecurity(t)
	u := &user.User{Id: 41, Login: "test", Email: "test@test.local", User: "Name", Role: "admin"}

	token := Security.makeToken(u)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString("{\"foo\":\"bar\", \"bar\":\"foo\"}"))
	c.Request.Header.Add("AUTH", token)

	Security.CheckRoleAuth("user")(c)

	if !c.IsAborted() {
		t.Fail()
	}
}
