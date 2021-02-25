package buisness

import (
	"github.com/dgrijalva/jwt-go"
	templates "purchases/pkg"
	"purchases/pkg/repository"
	"crypto/sha1"
	"errors"
	"fmt"
	"os"
	"time"
)

type AuthService struct {
	repo repository.Authorization
}

type CustomJWTClaims struct {
	jwt.StandardClaims
	UserId int `json:"id"`
}

var (
	jwtSecretKey = os.Getenv("JWT_SECRET_KEY")
)

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (as *AuthService) CreateUser(user templates.User) (int, error) {
	user.Password = as.getPasswordHash(user.Password)
	return as.repo.CreateUser(user)
}

func (as *AuthService) getPasswordHash(password string) string{
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(os.Getenv("PASS_GEN_SALT"))))
}

func (as *AuthService) GetUserId(username, password string) (int, error) {
	user, err := as.repo.GetUser(username, as.getPasswordHash(password))
	if err != nil {
		return 0, err
	}
	return user.Id, nil
}

func (as *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := as.repo.GetUser(username, as.getPasswordHash(password))
	if err != nil {
		return "", err
	}
	claims := CustomJWTClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			IssuedAt: time.Now().Unix(),
		},
		UserId:         user.Id,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecretKey))
}

func (as *AuthService) ParseToken(token string) (int, error){
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