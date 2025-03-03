package auths

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateAccountActivationToken(expirationDurationHours int) (string, string, time.Time, error) {
	tokenBytes := make([]byte, 20)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", "", time.Time{}, err
	}

	generatedToken := hex.EncodeToString(tokenBytes)

	hashedToken := sha256.New()
	hashedToken.Write([]byte(generatedToken))
	hashedTokenHex := hex.EncodeToString(hashedToken.Sum(nil))

	expirationDate := time.Now().Add(time.Duration(expirationDurationHours) * time.Hour)

	return generatedToken, hashedTokenHex, expirationDate, nil
}

func VerifyActivationToken(candidateToken, hashedToken string, tokenExpire time.Time) (bool, error) {
	hash := sha256.New()
	hash.Write([]byte(candidateToken))
	hashedCandidateToken := hex.EncodeToString(hash.Sum(nil))

	if hashedCandidateToken != hashedToken {
		return false, errors.New("invalid activation token")
	}

	if time.Now().After(tokenExpire) {
		return false, errors.New("activation token has expired")
	}

	return true, nil
}

func GenerateResetPasswordToken() (string, string, time.Time, error) {
	tokenBytes := make([]byte, 16)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", "", time.Time{}, fmt.Errorf("failed to generate random token: %v", err)
	}

	token := base64.URLEncoding.EncodeToString(tokenBytes)

	hashedToken := sha256.Sum256(tokenBytes)

	expirationDate := time.Now().Add(1 * time.Hour)

	return token, fmt.Sprintf("%x", hashedToken), expirationDate, nil
}

func VerifyResetPasswordToken(candidateToken, storedToken string, expireDate time.Time) (bool, error) {
	rawCandidateToken, err := base64.URLEncoding.DecodeString(candidateToken)
	if err != nil {
		fmt.Printf("Error decoding candidate token: %v\n", err)
		return false, fmt.Errorf("failed to decode token: %v", err)
	}

	hashedCandidateToken := sha256.Sum256(rawCandidateToken)

	if fmt.Sprintf("%x", hashedCandidateToken) == storedToken && time.Now().Before(expireDate) {
		return true, nil
	}

	if fmt.Sprintf("%x", hashedCandidateToken) != storedToken {
		fmt.Println("Token mismatch.")
	}
	if !time.Now().Before(expireDate) {
		fmt.Println("Token expired.")
	}

	return false, nil
}

func GenerateAccessToken(userID int, userInitials string) (string, error) {
	service := os.Getenv("ISSUE_NAME")
	recipient := os.Getenv("RECIPIENT_NAME")
	accessSecret := os.Getenv("ACCESS_SECRET")

	claims := jwt.MapClaims{
		"sub":  userID,
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(time.Hour * 1).Unix(),
		"iss":  service,
		"aud":  recipient,
		"init": userInitials,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(accessSecret))
	if err != nil {
		return "", fmt.Errorf("failed to generate access token: %v", err)
	}

	return accessToken, nil
}

func GenerateRefreshToken(userID int) (string, error) {
	service := os.Getenv("ISSUE_NAME")
	recipient := os.Getenv("RECIPIENT_NAME")
	refreshSecret := os.Getenv("REFRESH_SECRET")

	claims := jwt.MapClaims{
		"sub": userID,
		"iat": time.Now().Unix(),
		"nbf": time.Now().Add(time.Second * 30).Unix(),
		"exp": time.Now().Add(time.Hour * 5).Unix(),
		"iss": service,
		"aud": recipient,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshToken, err := token.SignedString([]byte(refreshSecret))
	if err != nil {
		return "", fmt.Errorf("failed to generate refresh token: %v", err)
	}

	return refreshToken, nil
}

func VerifyRefreshToken(tokenString string) (jwt.MapClaims, error) {
	refreshSecret := os.Getenv("REFRESH_SECRET")
	service := os.Getenv("ISSUE_NAME")
	recipient := os.Getenv("RECIPIENT_NAME")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(refreshSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid or expired token")
	}

	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return nil, errors.New("token has expired")
		}
	}

	if nbf, ok := claims["nbf"].(float64); ok {
		if time.Now().Unix() < int64(nbf) {
			return nil, errors.New("token is not valid yet")
		}
	}

	if iss, ok := claims["iss"].(string); !ok || iss != service {
		return nil, errors.New("invalid issuer")
	}

	if aud, ok := claims["aud"].(string); !ok || aud != recipient {
		return nil, errors.New("invalid audience")
	}

	return claims, nil
}

func CheckPassword(providedPassword, hashedPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(providedPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
