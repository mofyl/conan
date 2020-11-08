package jwt

import (
	"errors"
	"gopkg.in/dgrijalva/jwt-go.v3"
	"time"
	"conan/log"
	"conan/utils"
)

const (
	AUTH_NAME = "bee"
)

type TokenInfo struct {
	ChannelName string `json:"channel_name,omitempty"`
	ExternInfo  string `json:"extern_info,omitempty"`
}

type CustomClaims struct {
	TokenInfo
	ChannelKey string `json:"channel_key,omitempty"`
	RandomKey  string `json:"random_key,omitempty"`
	jwt.StandardClaims
}

type Jwt struct {
	key *PrivateKey
}

func NewJwt() *Jwt {
	key, err := NewPrivateKey(nil)

	if err != nil {
		panic("NewJwt Get PrivateKey Fail," + err.Error())
	}
	j := &Jwt{
		key: key,
	}
	//if ecode != nil {
	//	panic("init Jwt Fail ecode" + ecode.Error())
	//}

	return j
}

func NewJwtWithPrivateKey(key *PrivateKey) *Jwt {
	return &Jwt{
		key: key,
	}
}

func (j *Jwt) getRandomStr() string {
	uid, err := utils.GetUUID()
	if err != nil {
		return ""
	}
	return uid.String()
}

func (j *Jwt) decode(token string) error {

	to, err := jwt.ParseWithClaims(token, &CustomClaims{}, j.doVaild)

	if err != nil {
		return err
	}

	customClaims, ok := to.Claims.(*CustomClaims)

	if !ok {
		return errors.New("CustomClaims Convert Fail")
	}

	if !to.Valid {
		return errors.New("Token is not vaild")
	}
	expTime := customClaims.ExpiresAt
	auth := customClaims.Issuer

	log.Info("Decode expTime is %d , secret is %s", expTime, auth)
	return nil
}

func (j *Jwt) doVaild(token *jwt.Token) (interface{}, error) {
	return j.key.GetKey(), nil
}

func (j *Jwt) Encode(tokenInfo TokenInfo) (string, error) {

	timeNow := time.Now()
	beforeTime := timeNow.Unix() - 2*60
	afterTime := timeNow.Add(2 * time.Hour).Unix()
	claims := CustomClaims{}

	claims.TokenInfo = tokenInfo
	claims.ChannelKey = "qweqwe"
	claims.RandomKey = j.getRandomStr()
	claims.NotBefore = beforeTime
	claims.ExpiresAt = afterTime
	claims.IssuedAt = timeNow.Unix()
	claims.Issuer = AUTH_NAME

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)

	return token.SignedString(j.key.GetKey())
}
