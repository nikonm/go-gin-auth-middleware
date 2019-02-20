package user

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"strconv"
)

type LoginDTO struct {
	User         string
	Password     string
	Type         string
	PasswordHash string
}

type User struct {
	Id    int
	User  string
	Login string
	Email string
	Role  string
	Type  string
}

func (u *User) Fill(claims jwt.MapClaims) *User {
	for k, v := range claims {
		switch k {
		case "uid":
			u.Id = int(v.(float64))
			break
		case "role":
			u.Role = v.(string)
			break
		case "user":
			u.User = v.(string)
			break
		case "login":
			u.Login = v.(string)
			break
		case "email":
			u.Email = v.(string)
			break
		}
	}
	return u
}

func (u *User) FillFromSource(fromSource map[string]interface{}, mapping map[string]interface{}) (*User, error) {
	var err error
	for sourceField, targetField := range mapping {
		val, ok := fromSource[sourceField]
		if !ok {
			return nil, errors.New("No field '" + sourceField + "' find in result")
		}

		v := fmt.Sprintf("%v", val)

		switch targetField {
		case "id":
			u.Id, err = strconv.Atoi(v)
			break
		case "role":
			u.Role = v
			break
		case "user":
			u.User = v
			break
		case "login":
			u.Login = v
			break
		case "email":
			u.Email = v
			break
		}
	}

	return u, err
}
