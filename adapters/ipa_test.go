package adapters

import (
	"github.com/nikonm/go-gin-auth-middleware/user"
	"github.com/ubccr/goipa"
	"testing"
)

type IpaClientMock struct {
}

func (i *IpaClientMock) RemoteLogin(uid, passwd string) error {
	return nil
}

func (i IpaClientMock) UserShow(uid string) (*ipa.UserRecord, error) {
	return &ipa.UserRecord{
		UidNumber:   "123",
		DisplayName: "test",
		Email:       "test@test.local",
		Uid:         "test-123",
	}, nil
}

func getIpaProvider() (*IpaProvider, error) {
	m := &IpaClientMock{}
	p := &IpaProvider{}
	config := map[string]interface{}{
		"host":    "host.local",
		"timeout": "1m",
		"secured": false,
		"source_target_fields": map[string]interface{}{
			"id":       "id",
			"login":    "login",
			"username": "user",
		},
	}
	err := p.Init(config)
	p.SetCustomClient(m)
	return p, err
}

func TestIpaProvider_Login(t *testing.T) {
	p, err := getIpaProvider()

	if err != nil {
		t.Fail()
	}
	dto := user.LoginDTO{User: "test", Password: "test", Type: "ipa"}
	u, err := p.Login(dto)

	if u == nil {
		t.Fail()
	}
	if err != nil {
		t.Fail()
	}
}
