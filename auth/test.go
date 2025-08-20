package auth

import (
	"fmt"
	"net/http"
	"os"
)

type TestUser struct {
	Claims        map[string]any
	CurrentTenant string
}

func (u TestUser) Partner() string {
	return u.Claims["partner"].(string)
}

func (u TestUser) Tenant() string {
	return u.CurrentTenant
}

func (u TestUser) RoleByApp(appName string) string {
	roles := u.Claims["roles"].(map[string]any)
	return fmt.Sprintf("%v", roles[appName])
}

func (u TestUser) Apps() []string {
	roles := u.Claims["roles"].(map[string]any)
	apps := []string{}
	for app := range roles {
		apps = append(apps, fmt.Sprintf("%v", app))
	}
	return apps
}

func (u TestUser) FirstName() string {
	return u.Claims["firstName"].(string)
}

func (u TestUser) SurName() string {
	return u.Claims["surName"].(string)
}

func (u TestUser) Email() string {
	return u.Claims["eMail"].(string)
}

func (u TestUser) Username() string {
	return u.Claims["userName"].(string)
}

func (u TestUser) IsAdmin() bool {
	return u.Claims["admin"].(bool)
}

func (u TestUser) IsDeveloper() bool {
	return u.Claims["developer"].(bool)
}

func (u TestUser) Set(req *http.Request) error {
	if os.Getenv("AUTH_SECRET") == "" {
		return nil
	}
	authorization, err := CreateAuthorizationToken(u.Claims, os.Getenv("AUTH_SECRET"))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+string(authorization))
	return nil
}

func GetTestUserIdentity() TestUser {
	return TestUser{
		Claims: map[string]any{
			"partner": "PARTNER-X",
			"tenant":  []string{"default"},
			"roles": map[string]any{
				"testCase": "customer",
			},
			"firstName": "John",
			"surName":   "Rocket",
			"eMail":     "j.rocket@example.com",
			"userName":  "jrocket",
			"admin":     false,
			"developer": false,
		},
		CurrentTenant: "default",
	}
}
