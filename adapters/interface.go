package adapters

import "github.com/nikonm/go-gin-auth-middleware/user"

type Adapter interface {
	Init(config map[string]interface{}) error
	Login(dto user.LoginDTO) (*user.User, error) // user, error
}
