package services

import (
	"net/http"
	"time"

	"github.com/GoProjectGroupForEducation/Go-Blog/utils"

	"github.com/GoProjectGroupForEducation/Go-Blog/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
)

var (
	SecretKey = []byte( "awesome jwt nice to use")
)

type Token struct {
	Token string `json:"token"`
}

type CustomerClaims struct {
	User *models.User
	jwt.StandardClaims
}

func GenerateAuthToken(user *models.User) Token {
	expireToken := time.Now().Add(time.Hour * 24).Unix()

	claims := CustomerClaims{
		user,
		jwt.StandardClaims{
			ExpiresAt: expireToken,
			Issuer: "test.com",
		},
	}
	 token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		return Token{""}
	}
	return Token{tokenString}
}

func GetCurrentUser(tokenStr string) *models.User {
	token, err := jwt.ParseWithClaims(tokenStr[7:], &CustomerClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return SecretKey, nil
	})
	if err == nil {
		if claims, ok := token.Claims.(*CustomerClaims); ok && token.Valid {
			username := claims.User.Username
			return models.GetUserByUsername(username)
		} else {
			return nil
		}
	}
	return nil
}


func AuthenticationGuard(w http.ResponseWriter, req *http.Request, next utils.NextFunc) error {
	token, err := request.ParseFromRequest(req, request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (i interface{}, e error) {
			return SecretKey, nil
		})
	if err == nil && token.Valid {
		return next()
	} else {
		panic(utils.Exception{"Need to login first", http.StatusUnauthorized})
	}
}
