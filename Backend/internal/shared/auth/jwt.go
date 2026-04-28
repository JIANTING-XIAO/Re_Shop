package auth

import (
	"Re_Shop/Backend/internal/modules/user/model"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"time"
)

var ErrInvalidToken = errors.New("invalid token")

type UserClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Role     int8   `json:"role"`
	Exp      int64  `json:"exp"`
	Iat      int64  `json:"iat"`
}

func GenerateToken(user *model.User) (string, error) {
	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}

	now := time.Now()
	claims := UserClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		Iat:      now.Unix(),
		Exp:      now.Add(2 * time.Hour).Unix(),
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", err
	}

	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	headerPart := encodeSegment(headerJSON)
	claimsPart := encodeSegment(claimsJSON)
	unsignedToken := headerPart + "." + claimsPart
	signature := sign(unsignedToken)

	return unsignedToken + "." + signature, nil
}

func ParseToken(token string) (*UserClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, ErrInvalidToken
	}

	unsignedToken := parts[0] + "." + parts[1]
	expectedSignature := sign(unsignedToken)
	if !hmac.Equal([]byte(expectedSignature), []byte(parts[2])) {
		return nil, ErrInvalidToken
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, ErrInvalidToken
	}

	var claims UserClaims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, ErrInvalidToken
	}

	if claims.Exp < time.Now().Unix() {
		return nil, ErrInvalidToken
	}

	return &claims, nil
}

func encodeSegment(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}

func sign(data string) string {
	mac := hmac.New(sha256.New, []byte(jwtSecret()))
	mac.Write([]byte(data))
	return encodeSegment(mac.Sum(nil))
}

func jwtSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "reshop-dev-secret"
	}

	return secret
}
