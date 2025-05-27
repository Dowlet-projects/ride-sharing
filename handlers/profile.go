package handlers

import (

	"net/http"
	"encoding/json"
	"ride-sharing/models"
	"fmt"
	//"github.com/gorilla/mux"
	//"strconv"
	"ride-sharing/utils"
	//"strings"
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

// // Define struct for portial updates
// type PatchInput2 struct {
// 	FullName string `json:"full_name"`
// 	Phone string `json:"phone"`
// 	CarMake string `json:"car_make"`
// 	CarModel string `json:"car_model"`
// 	CarYear int `json:"car_year"`
// 	CarNumber string `json:"car_number"`
// }


// // PatchUser handles protected/taxist-profile PATCH
// // @Summary Update taxist profile
// // @Description Partially update taxist profile fields
// // @Tags Announcement
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param profile body PatchInput2 true "Profile fields to update"
// // @Router /protected/taxist-profile [patch]
// func (h *App) PatchUser(w http.ResponseWriter, r *http.Request) {
// 	claims, ok := r.Context().Value("claims").(*models.Claims)

// 	if !ok {
// 		utils.RespondError(w, http.StatusUnauthorized, "Invalid claims")
// 		return
// 	}

// 	taxist_id := claims.UserID;

// 	// Check if user exists
// 	var exists bool
// 	err := h.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM taxists WHERE id = ?)", taxist_id).Scan(&exists)
// 	if err != nil || !exists {
// 		http.Error(w, `{"error": "User not found"}`, http.StatusNotFound)
// 		return
// 	}
// 	type PatchInput struct {
// 	FullName *string `json:"full_name"`
// 	Phone *string `json:"phone"`
// 	CarMake *string `json:"car_make"`
// 	CarModel *string `json:"car_model"`
// 	CarYear *int `json:"car_year"`
// 	CarNumber *string `json:"car_number"`
// }

// 	var input PatchInput
// 	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
// 		fmt.Println(err)
// 		http.Error(w, `{"error": "Invalid input"}`, http.StatusBadRequest)
// 		return
// 	}

// 	// Build dynamic update query
// query := "UPDATE taxists SET "
// args := []interface{}{}
// fields := []string{} // Track fields to avoid trailing comma issues

// if input.FullName != nil && *input.FullName != "" {
//     fields = append(fields, "full_name = ?")
//     args = append(args, *input.FullName)
// }
// if input.Phone != nil && *input.Phone != "" {
//     fields = append(fields, "phone = ?")
//     args = append(args, *input.Phone)
// }
// if input.CarMake != nil && *input.CarMake != "" {
//     fields = append(fields, "car_make = ?")
//     args = append(args, *input.CarMake)
// }
// if input.CarModel != nil && *input.CarModel != "" {
//     fields = append(fields, "car_model = ?")
//     args = append(args, *input.CarModel)
// }
// if input.CarYear != nil {
//     fields = append(fields, "car_year = ?")
//     args = append(args, *input.CarYear)
// }
// if input.CarNumber != nil && *input.CarNumber != "" {
//     fields = append(fields, "car_number = ?")
//     args = append(args, *input.CarNumber)
// }

// // If no fields to update, return error
// if len(fields) == 0 {
//     http.Error(w, `{"error":"No valid fields provided for update"}`, http.StatusBadRequest)
//     return
// }

// // Join fields with commas and add WHERE clause
// query += strings.Join(fields, ", ") + " WHERE id = ?"
// args = append(args, taxist_id)

// // Log query for debugging
// fmt.Printf("Query: %s, Args: %v\n", query, args)

// // Perform partial update
// _, err = h.DB.Exec(query, args...)
// if err != nil {
//     fmt.Println(err)
//     http.Error(w, `{"error": "Failed to update user"}`, http.StatusInternalServerError)
//     return
// }

// 	var taxist models.Taxist
// 			if err := h.DB.QueryRow("SELECT id, full_name, phone, car_make, car_model, car_year, car_number, rating FROM taxists WHERE id = ? ", taxist_id).Scan(&taxist.ID, &taxist.FullName, 
// 			&taxist.Phone, &taxist.CarMake, &taxist.CarModel, &taxist.CarYear, &taxist.CarNumber, &taxist.Rating); err != nil {
// 			http.Error(w, "Server error", http.StatusInternalServerError)
// 		return
// 		}
// 		taxist.UserType = "taxist"
// 		w.Header().Set("Content-Type", "application/json")
// 		json.NewEncoder(w).Encode(taxist)

// }


// Define struct for full updates
type PutInput struct {
    FullName string `json:"full_name"`
    CarMake string `json:"car_make"`
    CarModel string `json:"car_model"`
    CarYear int `json:"car_year"`
    CarNumber string `json:"car_number"`
}

// PutUser handles protected/taxist-profile PUT
// @Summary Update taxist profile
// @Description Fully update taxist profile fields
// @Tags Announcement
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param profile body PutInput true "Complete profile to update"
// @Router /protected/taxist-profile [put]
func (h *App) PutUser(w http.ResponseWriter, r *http.Request) {
    claims, ok := r.Context().Value("claims").(*models.Claims)
    
    if !ok {
        utils.RespondError(w, http.StatusUnauthorized, "Invalid claims")
        return
    }
    
    taxist_id := claims.UserID
    
    // Check if user exists
    var exists bool
    err := h.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM taxists WHERE id = ?)", taxist_id).Scan(&exists)
    if err != nil || !exists {
        http.Error(w, `{"error": "User not found"}`, http.StatusNotFound)
        return
    }
    
    var input PutInput
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        fmt.Println(err)
        http.Error(w, `{"error": "Invalid input"}`, http.StatusBadRequest)
        return
    }
    
    // Validate required fields
    if input.FullName == "" || input.CarMake == "" || 
       input.CarModel == "" || input.CarYear == 0 || input.CarNumber == "" {
        http.Error(w, `{"error": "All fields are required"}`, http.StatusBadRequest)
        return
    }
    
    // Build update query with all fields
    query := `UPDATE taxists SET 
        full_name = ?,
        car_make = ?,
        car_model = ?,
        car_year = ?,
        car_number = ?
        WHERE id = ?`
    
    args := []interface{}{
        input.FullName,
        input.CarMake,
        input.CarModel,
        input.CarYear,
        input.CarNumber,
        taxist_id,
    }
    
    // Log query for debugging
    fmt.Printf("Query: %s, Args: %v\n", query, args)
    
    // Perform update
    _, err = h.DB.Exec(query, args...)
    if err != nil {
        fmt.Println(err)
        http.Error(w, `{"error": "Failed to update user"}`, http.StatusInternalServerError)
        return
    }
    
    var taxist models.Taxist
    if err := h.DB.QueryRow("SELECT id, full_name, phone, car_make, car_model, car_year, car_number, rating FROM taxists WHERE id = ?", taxist_id).Scan(
        &taxist.ID, 
        &taxist.FullName,
        &taxist.Phone, 
        &taxist.CarMake, 
        &taxist.CarModel, 
        &taxist.CarYear, 
        &taxist.CarNumber, 
        &taxist.Rating,
    ); err != nil {
        http.Error(w, `{"error": "Server error"}`, http.StatusInternalServerError)
        return
    }
    
    taxist.UserType = "taxist"
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(taxist)
}



// Define struct for full updates
type PutInputPassenger struct {
    FullName string `json:"full_name"`
}



// PutUserPassenger handles protected/passenger-profile PUT
// @Summary Update passenger profile
// @Description Fully update passenger profile fields
// @Tags Announcement
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param profile body PutInputPassenger true "Complete profile to update"
// @Router /protected/passenger-profile [put]
func (h *App) PutUserPassenger(w http.ResponseWriter, r *http.Request) {
    claims, ok := r.Context().Value("claims").(*models.Claims)
    
    if !ok {
        utils.RespondError(w, http.StatusUnauthorized, "Invalid claims")
        return
    }
    
    passenger_id := claims.UserID
    
    // Check if user exists
    var exists bool
    err := h.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM passengers WHERE id = ?)", passenger_id).Scan(&exists)
    if err != nil || !exists {
        http.Error(w, `{"error": "User not found"}`, http.StatusNotFound)
        return
    }
    
    var input PutInputPassenger
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        fmt.Println(err)
        http.Error(w, `{"error": "Invalid input"}`, http.StatusBadRequest)
        return
    }
    
    // Validate required fields
    if input.FullName == ""  {
        http.Error(w, `{"error": "All fields are required"}`, http.StatusBadRequest)
        return
    }
    
    // Build update query with all fields
    query := `UPDATE passengers SET 
        full_name = ? 
        WHERE id = ?`
    
    args := []interface{}{
        input.FullName,
        passenger_id,
    }
    
    // Log query for debugging
    fmt.Printf("Query: %s, Args: %v\n", query, args)
    
    // Perform update
    _, err = h.DB.Exec(query, args...)
    if err != nil {
        fmt.Println(err)
        http.Error(w, `{"error": "Failed to update user"}`, http.StatusInternalServerError)
        return
    }
    
    var passenger models.Passenger
    if err := h.DB.QueryRow("SELECT id, full_name, phone, created_at FROM passengers WHERE id = ?", passenger_id).Scan(
        &passenger.ID, 
        &passenger.FullName,
        &passenger.Phone, 
        &passenger.CreatedAt,
    ); err != nil {
        http.Error(w, `{"error": "Server error"}`, http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(passenger)
}




