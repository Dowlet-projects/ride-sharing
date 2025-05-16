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

// GetAllTaxistComments handles GET /taxist-comments/{taxist_id}
// @Summary Get all specific taxist comments by taxist_id
// @Description Retrieve all taxist comments by given taxist ID
// @Tags Announcement
// @Produce json
// @Security BearerAuth
// @Param taxist_id path string true "Taxist ID"
// @Router /protected/taxist-comments/{taxist_id} [get]
func (h *App) GetAllTaxistComments(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["taxist_id"])

	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	rows, err := h.DB.Query("SELECT id, full_name, comment from comments_to_taxist where taxist_id = ?", id)

	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	defer rows.Close()
	var comments []models.Comment =[]models.Comment{}
	for rows.Next() {
		var comment models.Comment
		if err := rows.Scan(&comment.ID, &comment.FullName, &comment.Comment); err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		comments = append(comments, comment)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

type PostComment struct {
	Comment string `json:"comment" example:"" description:"passenger comment"`
}

// GetTaxistNotifications handles GET /taxist-notifications
// @Summary Get all specific taxist notifications by taxist_id
// @Description Retrieve all taxist notifications by given taxist_id
// @Tags Announcement
// @Produce json
// @Security BearerAuth
// @Router /protected/taxist-notifications [get]
func (h *App) GetAllTaxistNotifications(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("claims").(*models.Claims)
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "Invalid claims")
		return
	}

	id := claims.UserID

	rows, err := h.DB.Query(
		`SELECT id, taxist_id, full_name, package, count, created_at 
		FROM taxist_notifications WHERE taxist_id = ?`, 
		id,
	)
	
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	defer rows.Close()
	var notifications []models.Notification = []models.Notification{}

	for rows.Next() {
		var notification models.Notification
		if err := rows.Scan(&notification.ID, &notification.TaxistID, &notification.FullName, &notification.Package, &notification.Count, &notification.CreatedAt); err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		notifications = append(notifications, notification)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifications)
}

// GetTaxistAnnouncements handles GET /taxist-announcements/{departed}
// @Summary Get all specific taxist announcements by departed
// @Description Retrieve all taxist announcements by given departed
// @Tags Announcement
// @Produce json
// @Security BearerAuth
// @Param departed path string true "Departed"
// @Router /protected/taxist-announcements/{departed} [get]
func (h *App) GetTaxistAnnouncements(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	departed, err := strconv.Atoi(vars["departed"])

	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	claims, ok := r.Context().Value("claims").(*models.Claims)

	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "Invalid claims")
		return
	}

	rows, err := h.DB.Query("SELECT id, taxist_id, depart_date, depart_time, space, distance, type, full_name, car_make, car_model, car_year, car_number, from_place, to_place, rating FROM ugurlar WHERE taxist_id = ? AND departed = ? ", claims.UserID, departed)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ugurlar)

}




// GetUgurById handles GET /taxist-notifications/{not_id}
// @Summary Get notification by not_id
// @Description Retrieve notification by given not_id
// @Tags Announcement
// @Produce json
// @Security BearerAuth
// @Param not_id path string true "Notification ID"
// @Router /protected/taxist-notifications/{not_id} [get]
func (h *App) GetNotificationById(w http.ResponseWriter, r *http.Request) {
	
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["not_id"])

	if err != nil {
		fmt.Println(err, id)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	query := `SELECT id, taxist_id, full_name, phone, package, count, created_at FROM 
	taxist_notifications WHERE id = ?`
	var notification models.NotificationDetails
	if err := h.DB.QueryRow(query, id).Scan(&notification.ID, &notification.TaxistID,
		 &notification.FullName, &notification.Phone, &notification.Package, &notification.Count, 
		 &notification.CreatedAt); err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
	}


	rows, err := h.DB.Query(
		`SELECT id, full_name, phone FROM reserve_passengers_people 
		WHERE reserve_id = ?`, id,
	)

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

	notification.Passengers = passengers

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notification)
}



// UpdateTaxistAnnouncements
// @Summary Update taxist announcement department
// @Description Updates a taxist's announcement department by taxi_ann_id
// @Tags Announcement
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param taxi_ann_id path int true "Announcement ID"
// @Router /protected/taxist-departed/{taxi_ann_id} [put]
func (h *App) UpdateTaxistAnnouncements(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["taxi_ann_id"])
	if err != nil {
		http.Error(w, "Invalid taxi announcement ID", http.StatusBadRequest)
		return
	}
	
	claims, ok := r.Context().Value("claims").(*models.Claims)
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "Invalid claims")
		return
	}

	query := "UPDATE taxist_announcements SET departed = 1 WHERE taxist_id = ? AND id = ?"
	result, err := h.DB.Exec(query, claims.UserID, id)
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


type PassengerNotification struct {
	WhoReserved int `json:"who_reserved"`
	FromPlaceName string `json:"from_place_name"`
	ToPlaceName string `json:"to_place_name"`
	FullName string `json:"full_name"`
	DepartDate string `json:"depart_date"`
	DepartTime string `json:"depart_time"`
	Count int `json:"count"`
	TaxiAnnId int `json:"taxi_ann_id"`
	Package string `json:"package"`
	CreatedAt string `json:"created_at"`
}


// GetAllPassengerNotifications handles GET /passenger-notifications
// @Summary Get all specific passenger notifications
// @Description Retrieve all passenger notifications
// @Tags Announcement
// @Produce json
// @Security BearerAuth
// @Router /protected/passenger-notifications [get]
func (h *App) GetAllPassengerNotifications(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("claims").(*models.Claims)
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "Invalid claims")
		return
	}

	id := claims.UserID

	rows, err := h.DB.Query(
		`SELECT who_reserved, from_place_name, to_place_name, full_name, 
		depart_date, depart_time, count, taxi_ann_id, package, created_at
		 FROM passenger_notifications WHERE who_reserved = ?`, id,
	)
	
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	defer rows.Close()
	var notifications []PassengerNotification = []PassengerNotification{}

	for rows.Next() {
		var notification PassengerNotification
		if err := rows.Scan(&notification.WhoReserved, &notification.FromPlaceName,
			&notification.ToPlaceName, &notification.FullName, &notification.DepartDate, 
			&notification.DepartTime, &notification.Count, &notification.TaxiAnnId,
			&notification.Package, &notification.CreatedAt); err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		notifications = append(notifications, notification)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifications)
}
