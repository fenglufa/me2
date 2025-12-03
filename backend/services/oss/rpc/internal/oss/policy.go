package oss

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// PolicyConditions OSS 上传策略条件
type PolicyConditions struct {
	Expiration string        `json:"expiration"`
	Conditions []interface{} `json:"conditions"`
}

// GeneratePolicy 生成 OSS 上传策略
func GeneratePolicy(dir string, maxSize int64, expireTime int64) (string, int64, error) {
	now := time.Now()
	expireEnd := now.Unix() + expireTime
	expireTimeStr := time.Unix(expireEnd, 0).UTC().Format("2006-01-02T15:04:05Z")

	// 构建策略条件
	policy := PolicyConditions{
		Expiration: expireTimeStr,
		Conditions: []interface{}{
			[]interface{}{"starts-with", "$key", dir},
			[]interface{}{"content-length-range", 0, maxSize},
		},
	}

	// 转为 JSON
	policyJSON, err := json.Marshal(policy)
	if err != nil {
		return "", 0, fmt.Errorf("marshal policy failed: %w", err)
	}

	// Base64 编码
	policyBase64 := base64.StdEncoding.EncodeToString(policyJSON)
	return policyBase64, expireEnd, nil
}

// GenerateSignature 生成 OSS 签名
func GenerateSignature(policyBase64, accessKeySecret string) string {
	h := hmac.New(sha1.New, []byte(accessKeySecret))
	h.Write([]byte(policyBase64))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return signature
}

// CompleteTokenClaims JWT Token Claims
type CompleteTokenClaims struct {
	ServiceName string `json:"service_name"`
	Dir         string `json:"dir"`
	UserID      int64  `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateCompleteToken 生成完成上传的验证 Token
func GenerateCompleteToken(serviceName, dir string, userID int64, jwtSecret string, expireTime int64) (string, error) {
	claims := CompleteTokenClaims{
		ServiceName: serviceName,
		Dir:         dir,
		UserID:      userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expireTime) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

// VerifyCompleteToken 验证完成上传的 Token
func VerifyCompleteToken(tokenString, jwtSecret string) (*CompleteTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CompleteTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("parse token failed: %w", err)
	}

	if claims, ok := token.Claims.(*CompleteTokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
