package handlers 

import (

	"net/http"
	"encoding/json"
	"ride-sharing/models"
	"fmt"
	"github.com/gorilla/mux"
	"strconv"
	"ride-sharing/utils"
)


type PassengerDetail struct {
	ID int `json:"id"`
	FullName string `json:"full_name"`
	Phone string `json:"phone"`
}

type PassengerDetails struct {
	Passengers []PassengerDetail `json:"passengers"`
	Package string `json:"package"`
	SubmitterName string `json:"submitter_name"`
	SubmitterPhone string `json:"submitter_phone"`
	CreatedAt string `json:"created_at"`
}

// ReverseDetails handles GET /reverse-details/{reverse_id}
// @Summary Get all specific taxist reversed Details by reverse_id
// @Description Retrieve all taxist reversed Details by given reverse ID
// @Tags Announcement
// @Produce json
// @Security BearerAuth
// @Param reverse_id path string true "Reverse ID"
// @Router /protected/reverse-details/{reverse_id} [get]
func (h *App) ReverseDetails(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["reverse_id"])

	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	rows, err := h.DB.Query("SELECT id, full_name, phone from reserve_passengers_people where reserve_id = ? ", id)

	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	defer rows.Close()
	var passengers []PassengerDetail = []PassengerDetail{}
	for rows.Next() {
		var passenger PassengerDetail
		if err := rows.Scan(&passenger.ID, &passenger.FullName, &passenger.Phone); err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		passengers = append(passengers, passenger)
	}

	var passengerDetails PassengerDetails

	if err := h.DB.QueryRow("SELECT package, full_name, phone, created_at FROM view_reverse_passengers where id = ? ", id).Scan(&passengerDetails.Package, &passengerDetails.SubmitterName, &passengerDetails.SubmitterPhone, &passengerDetails.CreatedAt); err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	passengerDetails.Passengers = passengers

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(passengerDetails)
}


//CreateComment handles POST /protected/taxist-comments
// @Summary POST a new comment for passenger user
// @Description add a new comment
// @Tags Announcement
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param taxist_id path string true "Taxist ID"
// @Param body body PostComment true "Passenger Comment"
// @Router /protected/taxist-comments/{taxist_id} [post]
func (h *App) CreateComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["taxist_id"])

	if err != nil {
		http.Error(w, "Invalid ID", http.StatusInternalServerError)
		return
	}

	claims, ok := r.Context().Value("claims").(*models.Claims)

	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "Invalid claims")
		return
	}

	passengerID := claims.UserID

	var postComment PostComment
	if err := json.NewDecoder(r.Body).Decode(&postComment); err != nil {
		http.Error(w, "Invalid Input", http.StatusBadRequest)
		return
	}

	query := "INSERT INTO taxist_comments (taxist_id, passenger_id, comment) VALUES (?, ?, ?)"
	_, err = h.DB.Exec(query, id, passengerID, postComment.Comment)

	if err != nil {
		http.Error(w, "Failed to create a comment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message":"successfully created",
	})
}

type Favourite struct {
	TaxistID int `json:"taxist_id"`
}

// CreateFavourites handles POST /protected/favourites
// @Summary POST a new passenger favourite
// @Description add a new passenger favourite
// @Tags Announcement
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body Favourite true "Announcement"
// @Router /protected/favourites [post]
func (h *App) CreateFavourites(w http.ResponseWriter, r *http.Request) {
	
	var favourite Favourite

	if err := json.NewDecoder(r.Body).Decode(&favourite); err != nil {
		http.Error(w, "Invalid Input", http.StatusBadRequest)
		return
	}

	claims, ok := r.Context().Value("claims").(*models.Claims)
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "Invalid claims")
		return
	}
	
	passenger_id := claims.UserID

	query := "CALL insert_halanlarym(?, ?)"
	_, err := h.DB.Exec(query, favourite.TaxistID, passenger_id)

	if err != nil {
		http.Error(w, "Failed to add or delete action", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message":"favourite taxist deleted or added",
	})
}

type Halanym struct {
	ID int `json:"id"`
	TaxistID int `json:"taxist_id"`
	FullName string `json:"full_name"`
	CarMake string `json:"car_make"`
	CarModel string `json:"car_model"`
	CarYear string `json:"car_year"`
	CarNumber string `json:"car_number"`
	Rating float32 `json:"rating"`
}

// GetAllFavourites handles GET /favourites
// @Summary Get all specific passenger favourites
// @Description Retrieve all passenger favourites
// @Tags Announcement
// @Produce json
// @Security BearerAuth
// @Router /protected/favourites [get]
func (h *App) GetAllFavourites(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("claims").(*models.Claims)
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "Invalid claims")
		return
	}
	
	passenger_id := claims.UserID
	rows, err := h.DB.Query("SELECT id, taxist_id, full_name, car_make, car_model, car_year, car_number, rating FROM halanlarym WHERE passenger_id = ? ", passenger_id)

	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	defer rows.Close()
    
	var halanlarym []Halanym = []Halanym{}

	for rows.Next() {
		var halanym Halanym
		if err := rows.Scan(&halanym.ID, &halanym.TaxistID, &halanym.FullName, &halanym.CarMake, &halanym.CarModel, &halanym.CarYear, &halanym.CarNumber, &halanym.Rating); err != nil {
			fmt.Println(err)
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		halanlarym = append(halanlarym, halanym)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(halanlarym)
}
