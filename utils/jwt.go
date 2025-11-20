package utils

import (
	"fmt"
	"gin/config"
	"time"
	"os"
	"crypto/rsa"
	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var (
    privateKey *rsa.PrivateKey
    publicKey  *rsa.PublicKey
)

func InitRSAKeys() error {
    return initRSAKeys()
}

func initRSAKeys() error {
    // 读取私钥文件
    privateKeyData, err := os.ReadFile("private_key.pem")
    if err != nil {
        return fmt.Errorf("读取私钥失败: %v", err)
    }

    // 解析私钥
    privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
    if err != nil {
        return fmt.Errorf("解析私钥失败: %v", err)
    }

    // 读取公钥文件
    publicKeyData, err := os.ReadFile("public_key.pem")
    if err != nil {
        return fmt.Errorf("读取公钥失败: %v", err)
    }

    // 解析公钥
    publicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicKeyData)
    if err != nil {
        return fmt.Errorf("解析公钥失败: %v", err)
    }

    fmt.Println("[JWT] RSA 密钥对初始化成功")
    return nil
}

// GetRSAPrivateKey 获取私钥
func GetRSAPrivateKey() *rsa.PrivateKey {
    if privateKey == nil {
        if err := initRSAKeys(); err != nil {
            fmt.Printf("[JWT] 错误: %v\n", err)
            return nil
        }
    }
    return privateKey
}

// GetRSAPublicKey 获取公钥
func GetRSAPublicKey() *rsa.PublicKey {
    if publicKey == nil {
        if err := initRSAKeys(); err != nil {
            fmt.Printf("[JWT] 错误: %v\n", err)
            return nil
        }
    }
    return publicKey
}

// getJWTSecretKey 获取JWT密钥（避免初始化顺序问题）
func getJWTSecretKey() []byte {
	if config.SecretKey == "" {
		fmt.Println("[JWT] 读取不到密钥，使用默认密钥")
		defaultKey := []byte("default-secret-key-at-least-256-bits-long-for-hs256-algorithm-needs-32bytes-minimum")
		fmt.Printf("[JWT] 默认密钥长度: %d 字节\n", len(defaultKey))
		return defaultKey
	}
	secretKey := []byte(config.SecretKey)
	fmt.Printf("[JWT] 读取到配置文件中的密钥，长度: %d 字节\n", len(secretKey))
	return secretKey
}

/*
生成Token
*/
func GenerateToken(username string) (string, error) {
    expirationTime := time.Now().Add(2 * time.Hour)
    claims := &Claims{
        Username: username,
        RegisteredClaims: jwt.RegisteredClaims{
            Subject:   username,                           // 主体
            Issuer:    config.AppName,                    // 签发者
            Audience:  jwt.ClaimStrings{"web", "mobile"}, // 受众
            ExpiresAt: jwt.NewNumericDate(expirationTime),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
            ID:        uuid.New().String(),               // Token 唯一 ID
        },
    }

    // 使用 PS512 算法签名
    token := jwt.NewWithClaims(jwt.SigningMethodPS512, claims)
    
    privateKey := GetRSAPrivateKey()
    if privateKey == nil {
        return "", fmt.Errorf("无法获取私钥")
    }

    signedToken, err := token.SignedString(privateKey)
    if err != nil {
        return "", fmt.Errorf("生成 Token 失败: %v", err)
    }

    fmt.Printf("[JWT] Token 生成成功 (PS512): %s...\n", signedToken[:50])
    return signedToken, nil
}
// func GenerateToken(username string) (string, error) {
// 	expirationTime := time.Now().Add(2 * time.Hour)
// 	claims := &Claims{
// 		Username: username,
// 		RegisteredClaims: jwt.RegisteredClaims{
// 			ExpiresAt: jwt.NewNumericDate(expirationTime),
// 			IssuedAt:  jwt.NewNumericDate(time.Now()),
// 			NotBefore: jwt.NewNumericDate(time.Now()),
// 			Issuer:    config.AppName,
// 			ID:        uuid.New().String(),
// 		},
// 	}
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	return token.SignedString(getJWTSecretKey())
// }

/*
* 解析Token
 */
 func ParseToken(tokenString string) (*Claims, error) {
    publicKey := GetRSAPublicKey()
    if publicKey == nil {
        return nil, fmt.Errorf("无法获取公钥")
    }

    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        // 验证算法
        if _, ok := token.Method.(*jwt.SigningMethodRSAPSS); !ok {
            return nil, fmt.Errorf("签名方法错误: 期望 PS512，实际 %v", token.Header["alg"])
        }
        return publicKey, nil
    })

    if err != nil {
        return nil, fmt.Errorf("Token 解析失败: %v", err)
    }

    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        fmt.Printf("[JWT] Token 验证成功，用户: %s\n", claims.Username)
        return claims, nil
    }

    return nil, fmt.Errorf("Token 验证失败")
}
// func ParseToken(tokenString string) (*Claims, error) {
// 	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
// 		return getJWTSecretKey(), nil
// 	})

// 	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
// 		return claims, nil
// 	}

// 	return nil, err
// }
