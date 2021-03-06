package security

import (
	"crypto/md5"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/nikonm/go-gin-auth-middleware/adapters"
	"github.com/nikonm/go-gin-auth-middleware/user"
	"hash"
	"net/http"
	"reflect"
	"strings"
	"time"
)

var Security *security

type Options struct {
	Secret       string `json:"secret"`
	TokenExp     string `json:"token_exp"`
	HeaderName   string `json:"header_name"`
	PasswordAlgo string `json:"password_algo"`

	Adapters map[string]map[string]interface{} //[type]provider
}

func getAdapter(key string) (adapters.Adapter, error) {
	switch key {
	case "db":
		return &adapters.DBProvider{}, nil
	case "ipa":
		return &adapters.IpaProvider{}, nil
	default:
		return nil, errors.New("Unknown security adapter " + key)
	}
}

type security struct {
	Options  Options
	Adapters map[string]*adapters.Adapter
}

func (s *security) CheckAuth(c *gin.Context) {
	s.getToken(c)
}

func (s *security) CheckRoleAuth(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		u, _ := s.GetUser(c)
		c.Set("SecuredUser", u)

		if c.IsAborted() {
			return
		}
		var in bool
	RolesLoop:
		for _, r := range roles {
			for _, ur := range u.Roles {
				if r == ur {
					in = true
					break RolesLoop
				}
			}
		}
		if !in {
			fmt.Println(c.Writer)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You don't have permission"})
			return
		}
	}
}

func (s *security) getTokenKey(c *gin.Context) string {
	return c.Request.Header.Get(s.Options.HeaderName)
}

func (s *security) getToken(c *gin.Context) (*jwt.Token, error) {
	tokenKey := s.getTokenKey(c)
	token, err := jwt.Parse(tokenKey, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.Options.Secret), nil
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return nil, err
	}
	if !token.Valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return token, err
	}
	return token, nil
}

func (s *security) GetUser(c *gin.Context) (*user.User, error) {
	token, err := s.getToken(c)
	if err != nil {
		return nil, err
	}
	claims := token.Claims.(jwt.MapClaims)

	u := (&user.User{}).Fill(claims)
	return u, err
}

func (s *security) Login(c *gin.Context, dto user.LoginDTO) (string, error) {
	var (
		err   error
		token string
	)
	dto.PasswordHash = s.PwdHash(dto.Password)
	for t, provider := range s.Adapters {
		if dto.Type != "" {
			if dto.Type == t {
				return s.providerLogin(provider, dto)
			}
		} else {
			token, err = s.providerLogin(provider, dto)
		}
	}
	if token != "" {
		err = nil
	} else if err == nil {
		err = errors.New("Unauthorized")
	}

	return token, err
}

func (s *security) providerLogin(adapter *adapters.Adapter, dto user.LoginDTO) (string, error) {
	u, err := (*adapter).Login(dto)
	if err == nil {
		u.Type = s.getAdapterKey(adapter)
		return s.makeToken(u), nil
	}
	return "", err
}

func (s *security) getAdapterKey(adapter *adapters.Adapter) string {
	for k, adpt := range s.Adapters {
		if reflect.DeepEqual(adpt, adapter) {
			return k
		}
	}
	return ""
}

func (s *security) makeToken(user *user.User) string {
	token := jwt.New(jwt.SigningMethodHS256)
	tokenExp, _ := time.ParseDuration(s.Options.TokenExp)
	token.Claims = jwt.MapClaims{
		"uid":   user.Id,
		"role":  strings.Join(user.Roles, ","),
		"user":  user.User,
		"email": user.Email,
		"type":  user.Type,
		"exp":   time.Now().Add(tokenExp).Unix(),
	}

	tokenString, _ := token.SignedString([]byte(s.Options.Secret))
	return tokenString
}

func (s *security) getHash() hash.Hash {
	switch s.Options.PasswordAlgo {
	case "sha1":
		return sha1.New()
	default:
		return md5.New()
	}
}

func (s *security) PwdHash(password string) string {
	h := s.getHash()
	h.Write([]byte(password))

	return fmt.Sprintf("%x", h.Sum(nil))
}

func Init(options Options) error {
	_adapters := make(map[string]*adapters.Adapter)
	for adapterKey, configs := range options.Adapters {
		adapter, err := getAdapter(adapterKey)
		if err != nil {
			return err
		}
		err = adapter.Init(configs)
		if err != nil {
			return err
		}
		_adapters[adapterKey] = &adapter
	}

	Security = &security{Options: options}
	Security.Adapters = _adapters
	return nil
}
