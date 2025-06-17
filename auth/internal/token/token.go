package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Token struct {
	key []byte
}

func New(SecretKey string) *Token {
	return &Token{[]byte(SecretKey)}
}

func (t *Token) CreateToken(user_id int64, role int8) ([]string, error) {

	var tokenSlice []string

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userid":   user_id,
			"userRole": role,
			"exp":      time.Now().Add(time.Minute * 15).Unix(),
		})
	tokenAccess, err := token.SignedString(t.key)
	if err != nil {
		return tokenSlice, err
	}

	token = jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userid":   user_id,
			"userRole": role,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})
	tokenRefresh, err := token.SignedString(t.key)
	if err != nil {
		return tokenSlice, err
	}

	tokenSlice = append(tokenSlice, tokenAccess, tokenRefresh)

	return tokenSlice, nil
}

func (t *Token) ValidateToken(refToken string) error {
	token, err := jwt.Parse(refToken, func(token *jwt.Token) (interface{}, error) {
		return t.key, nil
	})
	if err != nil {
		return err
	}

	if !token.Valid {
		return err
	}

	return nil
}

func (t *Token) GetClaims(refToken string) (map[string]interface{}, error) {
	claims := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(refToken, claims, func(token *jwt.Token) (interface{}, error) {
		return t.key, nil
	})
	if err != nil {
		return nil, err
	}

	return claims, nil
}
