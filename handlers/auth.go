package handlers

import (

	"net/http"
	"encoding/json"
	"ride-sharing/models"
	//"fmt"
	"time"
	"github.com/gorilla/mux"
	"strconv"
	"ride-sharing/utils"
	"database/sql"
	"regexp"
	"log"
)

// PassengerRegisterRequest defines the request body for passenger registration
type PassengerRegisterRequest struct {
	FullName string `json:"full_name" example:"John Doe" description:"Full name of the passenger" validate:"required"`
	Phone    string `json:"phone" example:"+12345678901" description:"Phone number in international format" validate:"required"`
}

// handlePassengerRegister handles passenger registration
// @Summary Register a new passenger
// @Description Registers a new passenger with a full name and phone number, sending a verification code. Validates input, checks for duplicate phone numbers, and logs the code (replace with SMS in production).
// @Tags Passenger
// @Accept json
// @Produce json
// @Param body body handlers.PassengerRegisterRequest true "Passenger registration details"
// @Router /passenger/register [post]
func (app *App) handlePassengerRegister(w http.ResponseWriter, r *http.Request) {
	var req struct {
		FullName string `json:"full_name"`
		Phone    string `json:"phone"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate input
	if req.FullName == "" {
		utils.RespondError(w, http.StatusBadRequest, "Full name is required")
		return
	}
	if !utils.ValidatePhone(req.Phone) {
		utils.RespondError(w, http.StatusBadRequest, "Invalid phone number format")
		return
	}

	// Check if phone exists
	var exists bool
	err := app.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM passengers WHERE phone = ?)", req.Phone).Scan(&exists)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Database error")
		return
	}
	if exists {
		utils.RespondError(w, http.StatusConflict, "Phone number already registered")
		return
	}

	// Generate verification code and store with registration data
	code, err := utils.GenerateVerificationCode(app.DB, req.Phone, req.FullName, "passenger", "", "", "", 0)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to generate verification code")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Verification code generated",
		"code":    code,
	})

}



// TaxistRegisterRequest defines the request body for taxist registration
type TaxistRegisterRequest struct {
	FullName  string `json:"full_name" example:"Jane Smith" description:"Full name of the taxist" validate:"required"`
	Phone     string `json:"phone" example:"+12345678902" description:"Phone number in international format" validate:"required,phone"`
	CarMake   string `json:"car_make" example:"Toyota" description:"Vehicle manufacturer" validate:"required"`
	CarModel  string `json:"car_model" example:"Camry" description:"Vehicle model" validate:"required"`
	CarYear   int    `json:"car_year" example:"2020" description:"Vehicle manufacturing year" validate:"required,gte=1900,lte=2026"`
	CarNumber string `json:"car_number" example:"ABC123" description:"Vehicle license plate number" validate:"required"`
}

type MakesRequest struct {
	Name  string `json:"name" example:"Toyota" description:"name of car make" validate:"required"`
}

type PlacesRequest struct {
	Name  string `json:"name" example:"Ashgabat" description:"name of place" validate:"required"`
}

type ModelsRequest struct {
	Name  string `json:"name" example:"Camry" description:"name of car model" validate:"required"`
	MakeId int `json:"make_id" example:"0" description:"car make id" validate:"required"`
}

type TaxistRating struct {
	Rating float32 `json:"rating" example:"4" description:"update rating of taxist" validate:"required"`
}


// handleTaxistRegister handles taxist registration
// @Summary Register a new taxist
// @Description Registers a new taxist with personal and vehicle details, sending a verification code. Validates input, checks for duplicate phone numbers, and logs the code.
// @Tags Taxist
// @Accept json
// @Produce json
// @Param body body handlers.TaxistRegisterRequest true "Taxist registration details"
// @Router /taxist/register [post]
func (app *App) handleTaxistRegister(w http.ResponseWriter, r *http.Request) {
	var req struct {
		FullName  string `json:"full_name"`
		Phone     string `json:"phone"`
		CarMake   string `json:"car_make"`
		CarModel  string `json:"car_model"`
		CarYear   int    `json:"car_year"`
		CarNumber string `json:"car_number"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate input
	if req.FullName == "" || req.CarMake == "" || req.CarModel == "" ||
		req.CarNumber == "" || req.CarYear < 1900 || req.CarYear > time.Now().Year()+1 {
		utils.RespondError(w, http.StatusBadRequest, "All fields are required and must be valid")
		return
	}
	if !utils.ValidatePhone(req.Phone) {
		utils.RespondError(w, http.StatusBadRequest, "Invalid phone number format")
		return
	}

	// Check if phone exists
	var exists bool
	err := app.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM taxists WHERE phone = ?)", req.Phone).Scan(&exists)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Database error")
		return
	}
	if exists {
		utils.RespondError(w, http.StatusConflict, "Phone number already registered")
		return
	}

	// Generate verification code and store with registration data
	code, err := utils.GenerateVerificationCode(app.DB, req.Phone, req.FullName, "taxist",
		req.CarMake, req.CarModel, req.CarNumber, req.CarYear)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to generate verification code")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Verification code generated",
		"code":    code,
	})
	
	// // Insert taxist
	// result, err := app.DB.Exec(
	// 	"INSERT INTO taxists (full_name, phone, car_make, car_model, car_year, car_number, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
	// 	req.FullName, req.Phone, req.CarMake, req.CarModel, req.CarYear, req.CarNumber, time.Now())
	// if err != nil {
	// 	utils.RespondError(w, http.StatusInternalServerError, "Failed to register taxist")
	// 	return
	// }
	// userID, _ := result.LastInsertId()
	
	// // Generate and store verification code
	// code, err := utils.GenerateVerificationCode(app.DB, req.Phone)
	// if err != nil {
	// 	utils.RespondError(w, http.StatusInternalServerError, "Failed to generate verification code")
	// 	return
	// }

	// // Simulate sending code
	// log.Printf("Verification code for %s: %s", req.Phone, code)
	
	// utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		// 	"message": "Verification code sent",
		// 	"user_id": userID,
	// })
}

// LoginRequest defines the request body for login
type LoginRequest struct {
	Phone string `json:"phone" example:"+12345678901" description:"Phone number in international format" validate:"required,phone"`
}

// handlePassengerLogin initiates passenger login
// @Summary Initiate passenger login
// @Description Initiates login by sending a verification code to the passenger's phone number. Checks if the phone is registered and logs the code.
// @Tags Passenger
// @Accept json
// @Produce json
// @Param body body handlers.LoginRequest true "Passenger login details"
// @Router /passenger/login [post]
func (app *App) handlePassengerLogin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Phone string `json:"phone"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate input
	if !utils.ValidatePhone(req.Phone) {
		utils.RespondError(w, http.StatusBadRequest, "Invalid phone number format")
		return
	}

	// Check if phone exists
	var userID int
	err := app.DB.QueryRow("SELECT id FROM passengers WHERE phone = ?", req.Phone).Scan(&userID)
	if err == sql.ErrNoRows {
		utils.RespondError(w, http.StatusNotFound, "Phone number not registered")
		return
	}
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Database error")
		return
	}
	
	// Generate verification code
	code, err := utils.GenerateVerificationCode(app.DB, req.Phone, "", "passenger", "", "", "", 0)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to generate verification code")
		return
	}
	
	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Verification code generated",
		"code":    code,
		"user_id": userID,
	})

	// // Generate and store verification code
	// code, err := utils.GenerateVerificationCode(app.DB, req.Phone)
	// if err != nil {
	// 	utils.RespondError(w, http.StatusInternalServerError, "Failed to generate verification code")
	// 	return
	// }

	// // Simulate sending code
	// log.Printf("Verification code for %s: %s", req.Phone, code)
	
	// utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
	// 	"message": "Verification code sent",
	// 	"user_id": userID,
	// })
}

// handleTaxistLogin initiates taxist login
// @Summary Initiate taxist login
// @Description Initiates login by sending a verification code to the taxist's phone number. Checks if the phone is registered and logs the code.
// @Tags Taxist
// @Accept json
// @Produce json
// @Param body body handlers.LoginRequest true "Taxist login details"
// @Router /taxist/login [post]
func (app *App) handleTaxistLogin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Phone string `json:"phone"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate input
	if !utils.ValidatePhone(req.Phone) {
		utils.RespondError(w, http.StatusBadRequest, "Invalid phone number format")
		return
	}

	// Check if phone exists
	var userID int
	err := app.DB.QueryRow("SELECT id FROM taxists WHERE phone = ?", req.Phone).Scan(&userID)
	if err == sql.ErrNoRows {
		utils.RespondError(w, http.StatusNotFound, "Phone number not registered")
		return
	}
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Database error")
		return
	}

	// Generate verification code
	code, err := utils.GenerateVerificationCode(app.DB, req.Phone, "", "taxist", "", "", "", 0)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to generate verification code")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Verification code generated",
		"code":    code,
		"user_id": userID,
	})

	// // Generate and store verification code
	// code, err := utils.GenerateVerificationCode(app.DB, req.Phone)
	// if err != nil {
	// 	utils.RespondError(w, http.StatusInternalServerError, "Failed to generate verification code")
	// 	return
	// }
	
	// // Simulate sending code
	// log.Printf("Verification code for %s: %s", req.Phone, code)

	// utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
	// 	"message": "Verification code sent",
	// 	"user_id": userID,
	// })
}

// VerifyCodeRequest defines the request body for verification
type VerifyCodeRequest struct {
	Phone string `json:"phone" example:"+12345678901" description:"Phone number in international format" validate:"required,phone"`
	Code  string `json:"code" example:"1234" description:"4-digit verification code" validate:"required,len=4"`
}

// handleVerifyCode verifies the code and completes registration/login
// @Summary Verify phone number
// @Description Verifies a phone number using a 4-digit code for passengers or taxists, issuing a JWT token upon success. Deletes the used code.
// @Tags Passenger,Taxist
// @Accept json
// @Produce json
// @Param body body handlers.VerifyCodeRequest true "Verification details"
// @Router /verify [post]
func (app *App) handleVerifyCode(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate input
	if !utils.ValidatePhone(req.Phone) {
		utils.RespondError(w, http.StatusBadRequest, "Invalid phone number format")
		return
	}
	if len(req.Code) != 4 || !regexp.MustCompile(`^\d{4}$`).MatchString(req.Code) {
		utils.RespondError(w, http.StatusBadRequest, "Invalid verification code")
		return
	}

	// Validate code and retrieve registration data
	var storedCode, fullName, userType, carMake, carModel, carNumber string
	var expiresAt time.Time
	var carYear int
	var rating float32
	err := app.DB.QueryRow(
		`SELECT code, expires_at, full_name, user_type, car_make, car_model, car_year, car_number, rating 
		 FROM verification_codes WHERE phone = ?`,
		 req.Phone).Scan(&storedCode, &expiresAt, &fullName, &userType, &carMake, &carModel, &carYear, &carNumber, &rating)
	if err == sql.ErrNoRows {
		utils.RespondError(w, http.StatusBadRequest, "No verification code found")
		return
	}
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Database error")
		return
	}

	if time.Now().After(expiresAt) {
		utils.RespondError(w, http.StatusBadRequest, "Verification code expired")
		return
	}

	if storedCode != req.Code {
		utils.RespondError(w, http.StatusBadRequest, "Invalid verification code")
		return
	}
	
	// Check if user is already registered
	var userID int64
	var isRegistered bool
	var registeredType models.UserType
	err = app.DB.QueryRow("SELECT id FROM passengers WHERE phone = ?", req.Phone).Scan(&userID)
	if err == nil {
		isRegistered = true
		registeredType = models.RolePassenger
	} else if err == sql.ErrNoRows {
		err = app.DB.QueryRow("SELECT id FROM taxists WHERE phone = ?", req.Phone).Scan(&userID)
		if err == nil {
			isRegistered = true
			registeredType = models.RoleTaxist
		}
	}
	if err != nil && err != sql.ErrNoRows {
		utils.RespondError(w, http.StatusInternalServerError, "Database error")
		return
	}

	// If not registered, complete registration
	if !isRegistered {
		if userType == "passenger" {
			if fullName == "" {
				utils.RespondError(w, http.StatusBadRequest, "Missing registration data")
				return
			}
			result, err := app.DB.Exec(
				"INSERT INTO passengers (full_name, phone, created_at) VALUES (?, ?, ?)",
				fullName, req.Phone, time.Now())
			if err != nil {
				utils.RespondError(w, http.StatusInternalServerError, "Failed to register passenger")
				return
			}
			userID, _ = result.LastInsertId()
			registeredType = models.RolePassenger
		} else if userType == "taxist" {
			if fullName == "" || carMake == "" || carModel == "" || carNumber == "" || carYear == 0 || rating < 0 {
				utils.RespondError(w, http.StatusBadRequest, "Missing registration data")
				return
			}
			result, err := app.DB.Exec(
				"INSERT INTO taxists (full_name, phone, car_make, car_model, car_year, car_number, rating, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
				fullName, req.Phone, carMake, carModel, carYear, carNumber, rating, time.Now())
			if err != nil {
				utils.RespondError(w, http.StatusInternalServerError, "Failed to register taxist")
				return
			}
			userID, _ = result.LastInsertId()
			registeredType = models.RoleTaxist
		} else {
			utils.RespondError(w, http.StatusBadRequest, "Invalid user type")
			return
		}
	}

	// Delete used code
	_, err = app.DB.Exec("DELETE FROM verification_codes WHERE phone = ?", req.Phone)
	if err != nil {
		log.Printf("Failed to delete verification code: %v", err)
	}

	// Generate JWT
	token, err := utils.GenerateJWT(app.Config, int(userID), registeredType)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"userID":userID,
		"token": token,
	})
}

// ProtectedResponse defines the response for the protected endpoint
type ProtectedResponse struct {
	UserID   int            `json:"user_id" example:"1" description:"ID of the authenticated user"`
	UserType models.UserType `json:"user_type" example:"passenger" description:"Type of user (passenger or taxist)"`
}

// handleProtected is an example protected route
// @Summary Access protected resource
// @Description Returns user information for an authenticated user. Requires a valid JWT token in the Authorization header.
// @Tags Authentication
// @Produce json
// @Security BearerAuth
// @Router /protected [get]
func (app *App) handleProtected(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("claims").(*models.Claims)
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "Invalid claims")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"user_id":   claims.UserID,
		"user_type": claims.UserType,
	})
}
	

// UpdateRatingTaxist godoc
// @Summary Update taxist rating
// @Description Updates a taxist's rating by taxist_id
// @Tags Announcement
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param taxist_id path int true "Taxist ID"
// @Param taxist_rating body handlers.TaxistRating true "Taxist rating details"
// @Router /protected/taxist-rating/{taxist_id} [put]
func (h *App) UpdateRatingTaxist(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["taxist_id"])
	if err != nil {
		http.Error(w, "Invalid taxist ID", http.StatusBadRequest)
		return
	}

	var rating TaxistRating
	if err:= json.NewDecoder(r.Body).Decode(&rating); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	query := "CALL ratingPut(?, ?)"
	result, err := h.DB.Exec(query, id, rating.Rating)
	if err != nil {
		http.Error(w, "Failed to update taxist rating", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Error checking update", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message":"successfully updated",
	})

}
