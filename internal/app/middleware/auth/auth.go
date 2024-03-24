package auth

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	tokenExp  = time.Hour
	secretKey = "secret_key"
)

type claims struct {
	jwt.RegisteredClaims
	UserID int
}

func Auth() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			rw := w

			var needAuthString bool
			var userID int

			cookie, err := r.Cookie("Authorization")
			if err != nil {
				needAuthString = true
			} else {
				userID = getUserID(cookie.Value)
				if userID < 0 {
					needAuthString = true
				}
			}

			cookie = &http.Cookie{
				Name: "Authorization",
			}

			if needAuthString {
				log.Println("auth", cookie)
				userID = 123
				token, err := buildJWTString(userID)
				log.Println("token", token)
				if err == nil {
					cookie.Value = token
					http.SetCookie(w, cookie)
				}
			}

			next.ServeHTTP(rw, r)
		}

		return http.HandlerFunc(fn)
	}
}

func buildJWTString(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
		},
		UserID: userID,
	})

	tokenStr, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

func getUserID(tokenString string) int {
	claims := &claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return -1
	}

	if !token.Valid {
		return -1
	}

	return claims.UserID
}

func CheckAuth(r *http.Request) int {
	var userID int

	token, err := r.Cookie("Authorization")
	if err != nil {
		return -1
	}

	userID = getUserID(token.Value)
	if userID < 0 {
		return -1
	}

	return userID
}
