package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/dchaykin/mygolib/log"
	"github.com/golang-jwt/jwt/v4"
)

type SimpleUserIdentity interface {
	FirstName() string
	SurName() string
	Email() string
	Username() string
	IsAdmin() bool
	IsDeveloper() bool
	Set(req *http.Request) error
}

type simpleUserToken struct {
	Claims jwt.MapClaims `json:"claims"`
}

func (j simpleUserToken) FirstName() string {
	claim, ok := j.Claims["firstName"]
	if !ok {
		log.Warn("Claim 'firstName' not found")
		return ""
	}
	return claim.(string)
}

func (j simpleUserToken) IsAdmin() bool {
	claim, ok := j.Claims["admin"]
	if !ok {
		return false
	}
	return claim.(bool)
}

func (j simpleUserToken) IsDeveloper() bool {
	if j.Username() == "dchaykin" { // TODO
		return true
	}
	claim, ok := j.Claims["developer"]
	if !ok {
		return false
	}
	return claim.(bool)
}

func (j simpleUserToken) SurName() string {
	claim, ok := j.Claims["surName"]
	if !ok {
		log.Warn("Claim 'surName' not found")
		return ""
	}
	return claim.(string)
}

func (j simpleUserToken) Email() string {
	claim, ok := j.Claims["eMail"]
	if !ok {
		log.Warn("Claim 'eMail' not found")
		return ""
	}
	return claim.(string)
}

func (j simpleUserToken) Username() string {
	claim, ok := j.Claims["userName"]
	if !ok {
		log.Warn("Claim 'userName' not found")
		return ""
	}
	return claim.(string)
}

func GetSimpleUserIdentityFromRequest(r http.Request) (SimpleUserIdentity, error) {
	userInfo := r.Header.Get("X-User-Info")
	if userInfo == "" {
		return nil, fmt.Errorf("no user info in the request found")
	}
	ui := simpleUserToken{}
	err := json.Unmarshal([]byte(userInfo), &ui)
	return ui, err
}

func (j simpleUserToken) Set(req *http.Request) error {
	authorization, err := CreateAuthorizationToken(j.Claims, os.Getenv("AUTH_SECRET"))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+string(authorization))
	return nil
}

func GetUserIdentity(authorization, secret string) (SimpleUserIdentity, error) {
	claims, err := parseToken(authorization, secret)
	if err != nil {
		return nil, err
	}
	return &simpleUserToken{
		Claims: claims,
	}, nil
}

func CreateAuthorizationToken(claims jwt.MapClaims, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
