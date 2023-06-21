package lib

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/geveit/go-api/src/config"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		config, err := config.GetAuthConfig()
		if err != nil {
			log.Println("Failed to load auth config")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		tokenString := strings.Split(authHeader, "Bearer ")
		if len(tokenString) != 2 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token, _ := jwt.Parse(tokenString[1], func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
			}

			kid, ok := t.Header["kid"].(string)
			if !ok {
				return nil, fmt.Errorf("Unexpected kid type: %v", t.Header["kid"])
			}

			publicKey, err := fetchPublicKey(kid, config)
			if err != nil {
				return nil, fmt.Errorf("Error fetching public key")
			}

			return publicKey, nil
		},
			jwt.WithIssuer(config.Iss),
			jwt.WithAudience(config.Aud),
		)

		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		sub, err := token.Claims.GetSubject()
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", sub)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func fetchPublicKey(kid string, config *config.AuthConfig) (interface{}, error) {
	response, err := http.Get(config.JwkUri)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var certs map[string]string
	if err := json.Unmarshal(body, &certs); err != nil {
		return nil, err
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(certs[kid]))
	if err != nil {
		return nil, err
	}

	return publicKey, nil
}
