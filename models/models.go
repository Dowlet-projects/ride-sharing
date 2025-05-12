// models/models.go
// Data models for the ride-sharing application

package models

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// UserType defines the type of user
type UserType string

const (
	RolePassenger UserType = "passenger"
	RoleTaxist    UserType = "taxist"
)

// Passenger represents a passenger user
type Passenger struct {
	ID        int       `json:"id"`
	FullName  string    `json:"full_name"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
}

// Taxist represents a taxist (driver) user
type Taxist struct {
	ID        int       `json:"id"`
	FullName  string    `json:"full_name"`
	Phone     string    `json:"phone"`
	CarMake   string    `json:"car_make"`
	CarModel  string    `json:"car_model"`
	CarYear   int       `json:"car_year"`
	CarNumber string    `json:"car_number"`
	Rating    float32   `json:"rating"`
	UserType  string `json:"user_type"`
	CreatedAt time.Time `json:"created_at"`
}

// VerificationCode represents a verification code entry
type VerificationCode struct {
	Phone     string    `json:"phone"`
	Code      string    `json:"code"`
	ExpiresAt time.Time `json:"expires_at"`
	FullName  string    `json:"full_name"`
	UserType  string    `json:"user_type"` // "passenger" or "taxist"
	CarMake   string    `json:"car_make"`
	CarModel  string    `json:"car_model"`
	CarYear   int       `json:"car_year"`
	CarNumber string    `json:"car_number"`
}

// Claims for JWT
type Claims struct {
	UserID   int      `json:"user_id"`
	UserType UserType `json:"user_type"`
	jwt.StandardClaims
}
