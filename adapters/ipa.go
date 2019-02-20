package adapters

import (
	"crypto/tls"
	"github.com/nikonm/go-gin-auth-middleware/user"
	"github.com/ubccr/goipa"
	"net/http"
	"strconv"
	"time"
)

type IpaClientInterface interface {
	RemoteLogin(uid, passwd string) error
	UserShow(uid string) (*ipa.UserRecord, error)
}

type IpaProvider struct {
	client IpaClientInterface
	config map[string]interface{}
}

/**
config = map[string]interface{}{
	"host": "host.local",
	"timeout": "1m",
	"secured": false,
	"source_target_fields": map[string]interface{}
}
*/
func (i *IpaProvider) Init(config map[string]interface{}) error {
	i.config = config
	var err error

	timeout, err := time.ParseDuration(config["timeout"].(string))
	if err != nil {
		return err
	}
	httpClient := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: !config["secured"].(bool)},
		},
	}

	i.client = ipa.NewClientCustomHttp(config["host"].(string), "", httpClient)

	return err
}

func (i *IpaProvider) SetCustomClient(c IpaClientInterface) {
	i.client = c
}

func (i *IpaProvider) Login(dto user.LoginDTO) (*user.User, error) {
	err := i.client.RemoteLogin(dto.User, dto.Password)
	if err != nil {
		return nil, err
	}
	userInfo, err := i.client.UserShow(dto.User)
	if err != nil {
		return nil, err
	}
	res := make(map[string]interface{})
	id, err := strconv.Atoi(userInfo.UidNumber.String())

	res["id"] = id
	res["login"] = userInfo.Uid.String()
	res["username"] = userInfo.DisplayName.String()
	res["email"] = userInfo.Email.String()

	return (&user.User{}).FillFromSource(res, i.config["source_target_fields"].(map[string]interface{}))
}
