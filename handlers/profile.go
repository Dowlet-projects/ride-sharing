package handlers

import (

	"net/http"
	"encoding/json"
	"ride-sharing/models"
	//"fmt"
	//"github.com/gorilla/mux"
	//"strconv"
	"ride-sharing/utils"
)

// // TaxistProfile handles GET /taxists/{taxist_id}
// // @Summary Get taxist by taxist_id
// // @Description Retrieve taxist by given taxist_id
// // @Tags Announcement
// // @Produce json
// // @Security BearerAuth
// // @Param taxist_id path string true "Taxist ID"
// // @Router /protected/taxists/{taxist_id} [get]
// func (h *App) TaxistProfile(w http.ResponseWriter, r *http.Request) {
	
// 	vars := mux.Vars(r)
// 	id, err := strconv.Atoi(vars["taxist_id"])

// 	if err != nil {
// 		fmt.Println(err, id)
// 		http.Error(w, "Invalid ID", http.StatusBadRequest)
// 		return
// 	}

// 	var taxist models.Taxist

// 	if err := h.DB.QueryRow("SELECT id, full_name, phone, car_make, car_model, car_year, car_number, rating FROM taxists WHERE id = ? ", id).Scan(&taxist.ID, &taxist.FullName, 
// 		&taxist.Phone, &taxist.CarMake, &taxist.CarModel, &taxist.CarYear, &taxist.CarNumber, &taxist.Rating); err != nil {
// 		http.Error(w, "Server error", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(taxist)
// }


// Profile handles GET /profile
// @Summary Get user profile 
// @Description Retrieve user
// @Tags Announcement
// @Produce json
// @Security BearerAuth
// @Router /protected/profile [get]
func (h *App) Profile(w http.ResponseWriter, r *http.Request) {
	
	claims, ok := r.Context().Value("claims").(*models.Claims)

	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "Invalid claims")
		return
	}



	UserID := claims.UserID
	UserType := claims.UserType

	if UserType == "passenger" {
		var passenger models.PassengerProfile
	
		if err := h.DB.QueryRow("SELECT id, full_name, phone FROM passengers WHERE id = ? ", UserID).Scan(&passenger.ID, &passenger.FullName, 
			&passenger.Phone); err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
	    passenger.UserType = "passenger"
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(passenger)
	} else {
			var taxist models.Taxist
			if err := h.DB.QueryRow("SELECT id, full_name, phone, car_make, car_model, car_year, car_number, rating FROM taxists WHERE id = ? ", UserID).Scan(&taxist.ID, &taxist.FullName, 
			&taxist.Phone, &taxist.CarMake, &taxist.CarModel, &taxist.CarYear, &taxist.CarNumber, &taxist.Rating); err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
		return
		}
		taxist.UserType = "taxist"
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(taxist)
	}

}