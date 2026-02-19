package service

import (
	"errors"
	"time"

	"red-packet/model"
	"red-packet/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	jwtSecret      []byte
	jwtExpireHours int
)

func InitUserService(secret string, expireHours int) {
	jwtSecret = []byte(secret)
	jwtExpireHours = expireHours
}

type Claims struct {
	UserID uint64 `json:"user_id"`
	jwt.RegisteredClaims
}

func Register(username, password string) (*model.User, error) {
	_, err := repository.GetUserByUsername(username)
	if err == nil {
		return nil, errors.New("username already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username:     username,
		PasswordHash: string(hash),
	}
	if err := repository.CreateUser(user); err != nil {
		return nil, err
	}
	return user, nil
}

func Login(username, password string) (string, error) {
	user, err := repository.GetUserByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("username or password incorrect")
		}
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("username or password incorrect")
	}

	claims := Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(jwtExpireHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

func GetProfile(userID uint64) (*model.User, error) {
	return repository.GetUserByID(userID)
}
