// utils/utils.go
// Helper functions for the application

package utils

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"regexp"
	"time"

	"ride-sharing/config"
	"ride-sharing/models"

	"github.com/dgrijalva/jwt-go"
)

// ValidatePhone checks if the phone number is valid
func ValidatePhone(phone string) bool {
	re := regexp.MustCompile(`^\+[1-9]\d{1,14}$`)
	return re.MatchString(phone)
}

// GenerateVerificationCode generates and stores a 6-digit code with registration data
func GenerateVerificationCode(db *sql.DB, phone, fullName, userType, carMake, carModel, carNumber string, carYear int) (string, error) {
	rating := 0
	n, err := rand.Int(rand.Reader, big.NewInt(10000))
	if err != nil {
		return "", fmt.Errorf("failed to generate code: %v", err)
	}
	code := fmt.Sprintf("%04d", n)

	expiresAt := time.Now().Add(5 * time.Minute)
	_, err = db.Exec(
		`INSERT INTO verification_codes 
		 (phone, code, expires_at, full_name, user_type, car_make, car_model, car_year, car_number, rating) 
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?) 
		 ON DUPLICATE KEY UPDATE 
		 code = ?, expires_at = ?, full_name = ?, user_type = ?, car_make = ?, car_model = ?, car_year = ?, car_number = ?, rating = ?`,
		phone, code, expiresAt, fullName, userType, carMake, carModel, carYear, carNumber, rating,
		code, expiresAt, fullName, userType, carMake, carModel, carYear, carNumber, rating)
	if err != nil {
		return "", fmt.Errorf("failed to store code: %v", err)
	}

	return code, nil
}

// GenerateJWT creates a JWT token
func GenerateJWT(cfg *config.Config, userID int, userType models.UserType) (string, error) {
	claims := &models.Claims{
		UserID:   userID,
		UserType: userType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}
	return signedToken, nil
}

// RespondJSON sends a JSON response
func RespondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("Failed to encode JSON: %v", err)
	}
}

// RespondError sends an error response
func RespondError(w http.ResponseWriter, status int, message string) {
	RespondJSON(w, status, map[string]string{"error": message})
}