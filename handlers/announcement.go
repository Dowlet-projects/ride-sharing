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




//CreateAnnouncement handles POST /protected/announcements
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
	fmt.Println(userID)
	var taxistAnn models.TaxistAnnouncement
	if err:=json.NewDecoder(r.Body).Decode(&taxistAnn); err !=nil {
		fmt.Println(err)
		http.Error(w, "Invalid Input", http.StatusBadRequest)
		return
	}

	query := "INSERT INTO taxist_announcements (taxist_id, depart_date, depart_time, from_place, to_place, space, type) VALUES (?, ?, ?, ?, ?, ?, ?)"
	_, err := h.DB.Exec(query, userID, taxistAnn.DepartDate, taxistAnn.DepartTime, taxistAnn.FromPlaceID, taxistAnn.ToPlaceID, taxistAnn.Space,  taxistAnn.Type)
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
	Phone string `json:"phone"`
}


type  PostReservePassengers struct {
	Passengers []PassengerPeople `json:"passengers" validate:"required"`
	Package string `json:"package" validate:"required"`
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
	query := "INSERT INTO reserve_passengers (package, taxi_ann_id, who_reserved, count) VALUES (?, ?, ?, ?)"
	 result, err := h.DB.Exec(query, postReservePassengers.Package, int_taxi_ann_id, userID, length)
	
	 if err != nil {
		http.Error(w, "Failed to create a reserved passenger", http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()

	reservedPassengers := postReservePassengers.Passengers
	for _, v := range reservedPassengers {
		que := "INSERT INTO reserve_passengers_people (full_name, phone, reserve_id, taxi_ann_id) VALUES (?, ?, ?, ?)"
		_, err := h.DB.Exec(que, v.FullName, v.Phone, id, int_taxi_ann_id)
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


// Success 200 {object} object{places=[]models.Place,page=int,limit=int,total=int,total_pages=int} "Paginated list of places"
// Failure 500 {string} string "Server error"

// GetAllUgurlar handles GET /protected/ugurlar
// @Summary Get all ugurlar
// @Description Retrieve a paginated list of all ugurlar
// @Tags Announcement
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of items per page (default: 10, max: 100)"
// @Router /protected/ugurlar [get]
func (h *App) GetAllUgurlar(w http.ResponseWriter, r *http.Request) {

	const defaultPage = 1
	const defaultLimit = 10
	const maxLimit = 100

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = defaultPage
	}

	limit, err := strconv.Atoi(limitStr)

	if err != nil || limit < 1 {
		limit = defaultLimit
	}

	if limit > maxLimit {
		limit = maxLimit
	}

	offset := ( page - 1 ) * limit

	rows, err := h.DB.Query("SELECT id, taxist_id, depart_date, depart_time, space, distance, type, full_name, car_make, car_model, car_year, car_number, from_place, to_place, rating FROM ugurlar LIMIT ? OFFSET ?", limit, offset)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var totalAnnouncements int

	err = h.DB.QueryRow("SELECT COUNT(*) FROM ugurlar").Scan(&totalAnnouncements)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	var ugurlar []models.Ugur = []models.Ugur{}

	for rows.Next() {
		var ugur models.Ugur

		if err := rows.Scan(&ugur.ID, &ugur.TaxistID, &ugur.DepartDate, &ugur.DepartTime, &ugur.Space, 
			&ugur.Distance, &ugur.Type, &ugur.FullName, &ugur.CarMake, &ugur.CarModel, &ugur.CarYear, 
			&ugur.CarNumber, &ugur.FromPlace, &ugur.ToPlace, &ugur.Rating); err != nil {
				http.Error(w, "Server error", http.StatusInternalServerError)
				return
		}

		ugurlar = append(ugurlar, ugur)
	}

	response := struct {
		Ugurlar []models.Ugur `json:"ugurlar"`
		Page int `json:"page"`
		Limit int `json:"limit"`
		Total int `json:"total"`
		TotalPages int `json:"total_pages"`
	}{
		Ugurlar: ugurlar,
		Page: page,
		Limit: limit,
		Total: totalAnnouncements,
		TotalPages: (totalAnnouncements +limit -1) / limit,
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

	var ugur models.UgurDetails
	if err := h.DB.QueryRow("SELECT id, taxist_id, depart_date, depart_time, space, distance, type, full_name, car_make, car_model, car_year, car_number, from_place, to_place, rating FROM ugurlar WHERE id = ? ", id).Scan(&ugur.ID, &ugur.TaxistID, &ugur.DepartDate, &ugur.DepartTime, &ugur.Space, 
		&ugur.Distance, &ugur.Type, &ugur.FullName, &ugur.CarMake, &ugur.CarModel, &ugur.CarYear, 
		&ugur.CarNumber, &ugur.FromPlace, &ugur.ToPlace, &ugur.Rating); err != nil {
			fmt.Println(err)
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
	}


	rows, err := h.DB.Query("SELECT id, full_name, phone FROM reserve_passengers_people WHERE taxi_ann_id = ?", id)

	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	
	defer rows.Close()

	var passengers []models.ReservePassengers = []models.ReservePassengers{} 
	for rows.Next() {
		var passenger models.ReservePassengers
		if err:=rows.Scan(&passenger.ID, &passenger.FullName, &passenger.Phone); err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		passengers = append(passengers, passenger)
	}

	ugur.Passengers = passengers

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ugur)
}