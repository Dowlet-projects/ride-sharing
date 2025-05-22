package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"ride-sharing/models"
	"ride-sharing/utils"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// CreateAnnouncement handles POST /protected/announcements
// @Summary POST a new Ugur for taxist
// @Description add a ugur
// @Tags Announcement
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body models.TaxistAnnouncmentRequest true "Ugur"
// @Router /protected/announcements [post]
func (h *App) CreateAnnouncement(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("claims").(*models.Claims)
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "Invalid claims")
		return
	}

	userID := claims.UserID
	var taxistAnn models.TaxistAnnouncement
	if err:=json.NewDecoder(r.Body).Decode(&taxistAnn); err !=nil {
		fmt.Println(err)
		http.Error(w, "Invalid Input", http.StatusBadRequest)
		return
	}

	query := "INSERT INTO taxist_announcements (taxist_id, depart_date, depart_time, from_place, to_place, full_space, space, type) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := h.DB.Exec(query, userID, taxistAnn.DepartDate, taxistAnn.DepartTime, taxistAnn.FromPlaceID, taxistAnn.ToPlaceID, taxistAnn.Space, taxistAnn.Space,  taxistAnn.Type)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to create a announcement", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message":"successfully created",
	})
}

type PassengerPeople struct {
	FullName string `json:"full_name"`
}


type  PostReservePassengers struct {
	Passengers []PassengerPeople `json:"passengers" validate:"required"`
	Package string `json:"package" validate:"required"`
	Phone string `json:"phone" validate:"required"` 
}

//CreateReservePassengers handles POST /protected/reserve-passengers
// @Summary POST a new reserving passenger for passenger user
// @Description add a new reserving passenger
// @Tags Announcement
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param taxi_ann_id path string true "Ugur ID"
// @Param body body PostReservePassengers true "Reserve Passenger"
// @Router /protected/reserve-passengers/{taxi_ann_id} [post]
func (h *App) CreateReservePassengers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	int_taxi_ann_id, err := strconv.Atoi(vars["taxi_ann_id"])

	if err != nil {
		http.Error(w, "Invalid id", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	claims, ok := r.Context().Value("claims").(*models.Claims)
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "Invalid claims")
		return
	}

	userID := claims.UserID
	
	var postReservePassengers PostReservePassengers
	if err:=json.NewDecoder(r.Body).Decode(&postReservePassengers); err !=nil {
		http.Error(w, "Invalid Input", http.StatusBadRequest)
		return
	}

	length := len(postReservePassengers.Passengers)
	query := "INSERT INTO reserve_passengers (package, phone, taxi_ann_id, who_reserved, count) VALUES (?, ?, ?, ?, ?)"
	 result, err := h.DB.Exec(query, postReservePassengers.Package, postReservePassengers.Phone, int_taxi_ann_id, userID, length)
	
	 if err != nil {
		http.Error(w, "Failed to create a reserved passenger", http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()

	reservedPassengers := postReservePassengers.Passengers
	for _, v := range reservedPassengers {
		que := "INSERT INTO reserve_passengers_people (full_name, reserve_id, taxi_ann_id) VALUES (?, ?, ?)"
		_, err := h.DB.Exec(que, v.FullName, id, int_taxi_ann_id)
		if err != nil {
			http.Error(w, "Failed to create reserved Passenger", http.StatusInternalServerError)
			return
		}
	}


	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message":"successfully created",
	})
}

//CreateReservePackages handles POST /protected/reserve-packages
// @Summary POST a new reserving package for passenger user
// @Description add a new reserving package
// @Tags Announcement
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param taxi_ann_id path string true "Ugur ID"
// @Param body body models.ReservePackages true "Reserve Packages"
// @Router /protected/reserve-packages/{taxi_ann_id} [post]
func (h *App) CreateReservePackages(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	int_taxi_ann_id, err := strconv.Atoi(vars["taxi_ann_id"])

	if err != nil {
		fmt.Println(err)
	}

	claims, ok := r.Context().Value("claims").(*models.Claims)
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "Invalid claims")
		return
	}

	userID := claims.UserID
	
	var reservePackages models.ReservePackages
	if err:=json.NewDecoder(r.Body).Decode(&reservePackages); err !=nil {
		http.Error(w, "Invalid Input", http.StatusBadRequest)
		return
	}

	query := "INSERT INTO reserve_packages (package_sender, package_reciever, sender_phone, reciever_phone, about_package, taxi_ann_id, who_reserved) VALUES (?, ?, ?, ?, ?, ?, ?)"
	 _, err = h.DB.Exec(query, reservePackages.PackageSender, reservePackages.PackageReciever, reservePackages.SenderPhone,reservePackages.RecieverPhone,reservePackages.AboutPackage, int_taxi_ann_id, userID)
	
	 if err != nil {
		http.Error(w, "Failed to create a reserve package", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message":"successfully created",
	})
}



// GetAllUgurlar handles GET /protected/ugurlar
// @Summary Get all ugurlar
// @Description Retrieve a paginated list of ugurlar filtered by optional query parameters
// @Tags Announcement
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (default: 1)" example=1
// @Param limit query int false "Number of items per page (default: 10, max: 100)" example=10
// @Param date query string false "Departure date in YYYY-MM-DD format" format=date example="2025-05-12"
// @Param from_place query int false "Departure place ID" example=1
// @Param to_place query int false "Destination place ID" example=2
// @Param car_make query string false "Car make name" example=3
// @Param car_model query string false "Car model name" example=4
// @Param space query int false "Space" example=3
// @Param rating query int false "Rating" example=3
// @Param taxist_id query int false "taxist Id" example=1
// @Param passenger_type query string false "Passenger type" example="passenger"
// @Router /protected/ugurlar [get]
func (h *App) GetAllUgurlar(w http.ResponseWriter, r *http.Request) {
	const defaultPage = 1
	const defaultLimit = 10
	const maxLimit = 100

	// Parse query parameters
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	dateStr := r.URL.Query().Get("date")
	fromPlaceStr := r.URL.Query().Get("from_place")
	toPlaceStr := r.URL.Query().Get("to_place")
	makeIDStr := r.URL.Query().Get("car_make")
	modelIDStr := r.URL.Query().Get("car_model")
	typeStr := r.URL.Query().Get("passenger_type")
	spaceStr := r.URL.Query().Get("space")
	ratingStr := r.URL.Query().Get("rating")
	taxistIdStr := r.URL.Query().Get("taxist_id")

	fmt.Println(typeStr)

	// Parse and validate page
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = defaultPage
	}

	// Parse and validate limit
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = defaultLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}

	offset := (page - 1) * limit

	// Build the WHERE clause dynamically
	query := `SELECT id, taxist_id, depart_date, depart_time, full_space, space, distance, type, full_name, car_make, car_model, 
	car_year, car_number, from_place, to_place, rating FROM ugurlar WHERE id IN (SELECT id FROM taxist_announcements_filter`
	countQuery := "SELECT COUNT(*) FROM ugurlar"
	var conditions []string
	var args []interface{}

	// Add filters based on provided parameters
	if dateStr != "" {
		if matched, _ := regexp.MatchString(`^\d{4}-\d{2}-\d{2}$`, dateStr); matched {
			conditions = append(conditions, "depart_date = ?")
			args = append(args, dateStr)
		} else {
			http.Error(w, "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest)
			return
		}
	}
	if fromPlaceStr != "" {
		fromPlace, err := strconv.Atoi(fromPlaceStr)
		if err == nil {
			conditions = append(conditions, "from_place = ?")
			args = append(args, fromPlace)
		}
	}
	if toPlaceStr != "" {
		toPlace, err := strconv.Atoi(toPlaceStr)
		if err == nil {
			conditions = append(conditions, "to_place = ?")
			args = append(args, toPlace)
		}
	}
	if makeIDStr != "" {
		conditions = append(conditions, "car_make = ?")
		args = append(args, makeIDStr)
	}

	if modelIDStr != "" {
		conditions = append(conditions, "car_model = ?")
		args = append(args, modelIDStr)
	}
	if spaceStr != "" {
		space, err := strconv.Atoi(spaceStr)
		if err == nil {
			conditions = append(conditions, "space = ?")
			args = append(args, space)
		}
	}

	if taxistIdStr != "" {
		taxist_id, err := strconv.Atoi(taxistIdStr)
		if err == nil {
			conditions = append(conditions, "taxist_id = ?")
			args = append(args, taxist_id)
		}
	}

	if typeStr != "" {
		conditions = append(conditions, "type = ?")
		args = append(args, typeStr)
	}

	if ratingStr != "" {
		rating, err := strconv.Atoi(ratingStr)
		if err == nil {
			conditions = append(conditions, "rating BETWEEN ? AND ?")
			args = append(args, rating, rating+1)
		}
	}

	// Append WHERE clause if there are conditions
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
		countQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	// Append LIMIT and OFFSET
	query += ") LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	// Execute the query
	rows, err := h.DB.Query(query, args...)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Fetch ugurlar
	var ugurlar []models.Ugur = []models.Ugur{}
	for rows.Next() {
		var ugur models.Ugur
		if err := rows.Scan(&ugur.ID, &ugur.TaxistID, &ugur.DepartDate, &ugur.DepartTime, &ugur.FullSpace, &ugur.Space,
			&ugur.Distance, &ugur.Type, &ugur.FullName, &ugur.CarMake, &ugur.CarModel, &ugur.CarYear,
			&ugur.CarNumber, &ugur.FromPlace, &ugur.ToPlace, &ugur.Rating); err != nil {
			fmt.Println(err)
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		ugurlar = append(ugurlar, ugur)
	}

	// Check for errors during row iteration
	if err := rows.Err(); err != nil {
		fmt.Println(err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Get total count with filters
	var totalAnnouncements int
	err = h.DB.QueryRow(countQuery, args[:len(args)-2]...).Scan(&totalAnnouncements)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Prepare response
	response := struct {
		Ugurlar     []models.Ugur `json:"ugurlar"`
		Page        int           `json:"page"`
		Limit       int           `json:"limit"`
		Total       int           `json:"total"`
		TotalPages  int           `json:"total_pages"`
	}{
		Ugurlar:    ugurlar,
		Page:       page,
		Limit:      limit,
		Total:      totalAnnouncements,
		TotalPages: (totalAnnouncements + limit - 1) / limit,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}


// GetUgurById handles GET /ugurlar/{ugur_id}
// @Summary Get ugur by ugur_id
// @Description Retrieve ugur by given ugur_id
// @Tags Announcement
// @Produce json
// @Security BearerAuth
// @Param ugur_id path string true "Ugur ID"
// @Router /protected/ugurlar/{ugur_id} [get]
func (h *App) GetUgurById(w http.ResponseWriter, r *http.Request) {
	
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["ugur_id"])

	if err != nil {
		fmt.Println(err, id)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	claims, ok := r.Context().Value("claims").(*models.Claims)
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "Invalid claims")
		return
	}

	userID := claims.UserID

	var ugur models.UgurDetails
	if err := h.DB.QueryRow(`SELECT id, taxist_id, depart_date, depart_time, full_space, space, distance, type, 
	full_name, taxist_phone, car_make, car_model, car_year, car_number, from_place, to_place, rating FROM ugurlar 
	WHERE id = ? `, id).Scan(&ugur.ID, &ugur.TaxistID, &ugur.DepartDate, &ugur.DepartTime, &ugur.FullSpace, &ugur.Space, 
		&ugur.Distance, &ugur.Type, &ugur.FullName, &ugur.TaxistPhone, &ugur.CarMake, &ugur.CarModel, &ugur.CarYear, 
		&ugur.CarNumber, &ugur.FromPlace, &ugur.ToPlace, &ugur.Rating); err != nil {
			fmt.Println(err)
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
	}


	rows, err := h.DB.Query("SELECT id, full_name FROM reserve_passengers_people WHERE taxi_ann_id = ?", id)

	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	
	defer rows.Close()

	query := "SELECT 1 FROM favourites WHERE taxist_id = ? AND passenger_id = ? "
	var exists bool
	
	err = h.DB.QueryRow(query, ugur.TaxistID, userID).Scan(&exists)
	if err != nil {
		ugur.TaxistLike = false
	} else {
		ugur.TaxistLike = exists
	}


	var passengers []models.ReservePassengers = []models.ReservePassengers{} 
	for rows.Next() {
		var passenger models.ReservePassengers
		if err:=rows.Scan(&passenger.ID, &passenger.FullName); err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		passengers = append(passengers, passenger)
	}

	ugur.Passengers = passengers

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ugur)
}