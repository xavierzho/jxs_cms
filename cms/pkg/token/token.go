package token

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	Username string `json:"user_name"`
	UserID   string `json:"user_id"`
	jwt.StandardClaims
}

func getSecret() []byte {
	return []byte(JWTSetting.Secret)
}

func GetRKeyByUserID(userID uint32) string {
	return fmt.Sprintf("Token:%d", userID)
}

func GenerateToken(userName string, userID uint32) (tokenStr string, err error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(JWTSetting.Expire)
	claims := Claims{
		Username: hex.EncodeToString([]byte(userName)),
		UserID:   hex.EncodeToString([]byte(strconv.FormatUint(uint64(userID), 10))),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    JWTSetting.Issuer,
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err = tokenClaims.SignedString(getSecret())
	if err != nil {
		return "", fmt.Errorf("jwt.GenerateToken: %s", err.Error())
	}

	return tokenStr, nil
}

func ParseToken(tokenStr string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return getSecret(), nil
	})
	if err != nil {
		return nil, err
	}
	if tokenClaims != nil {
		claims, ok := tokenClaims.Claims.(*Claims)
		if ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
