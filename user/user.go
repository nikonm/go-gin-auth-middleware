package user

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"reflect"
	"strconv"
	"strings"
)

type LoginDTO struct {
	User         string
	Password     string
	Type         string
	PasswordHash string `json:"-"`
}

type User struct {
	Id    int
	User  string
	Login string
	Email string
	Roles []string
	Type  string
}

func (u *User) Fill(claims jwt.MapClaims) *User {
	for k, v := range claims {
		switch k {
		case "uid":
			u.Id = int(v.(float64))
		case "role":
			u.Roles = strings.Split(v.(string), ",")
		case "user":
			u.User = v.(string)
		case "login":
			u.Login = v.(string)
		case "email":
			u.Email = v.(string)
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
		case "roles":
			if reflect.TypeOf(val).Kind() == reflect.Slice {
				u.Roles = val.([]string)
				break
			}
			u.Roles = []string{v}
		case "user":
			u.User = v
		case "login":
			u.Login = v
		case "email":
			u.Email = v
		}
	}

	return u, err
}
