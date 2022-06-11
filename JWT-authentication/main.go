package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	// "golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("my_jwt_key")

var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func Login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		if loggedIn, _ := IsLoggedIn(r); loggedIn {
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write([]byte("user already logged in"))
			return
		}

		var user User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		expectedPassword, ok := users[user.Username]
		if !ok || expectedPassword != user.Password {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("username or password invalid"))
			return
		}

		expirationTime := time.Now().Add(5 * time.Minute)
		claims := &Claims{
			Username: user.Username,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		fmt.Println("generated token:", tokenString)

		cookie := &http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
		}
		http.SetCookie(w, cookie)

		w.WriteHeader(http.StatusOK)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func Home(w http.ResponseWriter, r *http.Request) {
	loggedin, claims := IsLoggedIn(r)
	if !loggedin {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hello %s", claims.Username)))
}

func Refresh(w http.ResponseWriter, r *http.Request) {
	islogged, claims := IsLoggedIn(r)
	if !islogged {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30 {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Println("generated token:", tokenString)

	cookie := &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	}
	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("/login", Login)
	http.HandleFunc("/home", Home)
	http.HandleFunc("/refresh", Refresh)

	log.Fatal(http.ListenAndServe(":9090", nil))
}

func IsLoggedIn(r *http.Request) (loggedIn bool, claims *Claims) {
	cookie, err := r.Cookie("token")
	if err != nil {
		return false, nil
	}

	tokenString := cookie.Value
	fmt.Println("Token string:", tokenString)
	if tokenString == "" {
		return false, nil
	}

	c := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, c, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		fmt.Println("error parsing claims")
		return false, nil
	}

	if !token.Valid {
		log.Println("invalid token")
		return false, nil
	}

	return true, c
}
