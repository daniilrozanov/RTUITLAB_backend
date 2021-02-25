package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	UsersTransportKey = os.Getenv("USERS_TRANSPORT_KEY")
	jwtSecretKey = os.Getenv("JWT_SECRET_KEY")
)

type UserServiceConfig struct {
	Host string
	Port string
	URN string
	Scheme string
}

type AuthService struct {
	uConfs *UserServiceConfig
}

func NewAuthService(uConfs *UserServiceConfig) *AuthService {
	return &AuthService{uConfs: uConfs}
}

type CustomJWTClaims struct {
	jwt.StandardClaims
	UserId int `json:"id"`
}

func (a *AuthService) ConfirmUser(name, password string) (int, error) {
	cryptedName, _ := encrypt([]byte(name), UsersTransportKey)
	cryptedPass, _ := encrypt([]byte(password), UsersTransportKey)
	fmt.Printf(">>>%s\n>>>%s", string(cryptedName), string(cryptedPass))
	byteJSON, err := json.Marshal(map[string][]byte{
		"name": cryptedName,
		"password": cryptedPass,
	})
	if err != nil {
		return 0, errors.New("invalid JSON format")
	}
	fmt.Println(a.getUserServiceURI())
	response, err := http.Post(a.getUserServiceURI(), "application/json", bytes.NewBuffer(byteJSON))
	if err != nil { // если не ответит сервис покупок
		return 0, err
	}
	var result map[string][]byte
	if err = json.NewDecoder(response.Body).Decode(&result); err != nil { // если придет иной формат
		return 0, err
	}
	if string(result["id"]) == "" { // если не пришло поле айди
		return 0, errors.New("bad json response")
	}
	decodedBytes, err := decrypt(result["id"], UsersTransportKey)
	if err != nil { // если айли не расшифровался
		return 0, nil
	}
	strId := string(decodedBytes)
	if strId == "not found" { // если пользователся не существует
		return 0, errors.New("user not found")
	}

	userId, err := strconv.Atoi(strId)
	if err != nil { // если пришло не число
		return 0, err
	}
	return userId, nil
}

func (a *AuthService) getUserServiceURI() string {
	if a.uConfs.Port == ":" {
		return a.uConfs.Scheme + a.uConfs.Host + "/" + a.uConfs.URN
	}
	return a.uConfs.Scheme +  a.uConfs.Host + ":" + a.uConfs.Port + "/" + a.uConfs.URN
}

func (a *AuthService) GenerateToken(id int) (string, error) {
	claims := CustomJWTClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			IssuedAt: time.Now().Unix(),
		},
		UserId:         id,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecretKey))
}
func (a *AuthService) ParseToken(token string) (int, error) {
	newToken, err := jwt.ParseWithClaims(token, &CustomJWTClaims{}, func(ltoken *jwt.Token) (interface{}, error) {
		if _, ok := ltoken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token")
		}
		return []byte(jwtSecretKey), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := newToken.Claims.(*CustomJWTClaims)
	if !ok {
		return 0, errors.New("invalid token claims")
	}
	return claims.UserId, nil
}
