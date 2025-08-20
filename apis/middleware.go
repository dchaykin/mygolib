package apis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/dchaykin/mygolib/auth"
	log "github.com/dchaykin/mygolib/log"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Request: %s %s", r.Method, r.URL.Path)
		userIdentity, err := getUserIdentityFromToken(r.Header.Get("Authorization"))
		if err != nil {
			log.Info("Invalid user token: %v", err)
		} else {
			userData, err := json.Marshal(userIdentity)
			if err != nil {
				log.Info("Invalid user structure: %v", err)
				return
			}
			reqClone := r.Clone(r.Context())
			reqClone.Header.Set("X-User-Info", string(userData))
			log.Debug("User: %s", userData)
			next.ServeHTTP(w, reqClone)
			return
		}
		http.Error(w, fmt.Sprintf("%v", err), http.StatusUnauthorized)
	})
}

func getUserIdentityFromToken(token string) (auth.SimpleUserIdentity, error) {
	if token != "" {
		token := strings.TrimPrefix(token, "Bearer ")
		return auth.GetUserIdentity(token, os.Getenv("AUTH_SECRET"))
	}
	return nil, fmt.Errorf("no token found")
}
