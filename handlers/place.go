package handlers


import (

	"net/http"
	"encoding/json"
	"ride-sharing/models"
	//"fmt"
	"github.com/gorilla/mux"
	"strconv"
)

//GetAllPlaces handles GET /places
// @Summary GET car places
// @Description get all places
// @Tags Place
// @Produce json
// @Router /places [get]
func (h *App) GetAllPlaces(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query("SELECT id, name FROM places")

	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var places []models.Place = []models.Place{}
	for rows.Next() {
		var place models.Place
		if err:=rows.Scan(&place.ID, &place.Name); err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		places = append(places, place)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(places)
}

type Distance struct {
	ID int `json:"id"`
	FromPlace string `json:"from_place"`
	ToPlace string `json:"to_place"`
	Distance string `json:"distance"`
}


//GetAllDistances handles GET /distances
// @Summary GET distances
// @Description get all distances
// @Tags Place
// @Produce json
// @Router /distances [get]
func (h *App) GetAllDistances(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query("SELECT id, from_place, to_place, distance FROM distances")

	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var distances []Distance = []Distance{}
	for rows.Next() {
		var distance Distance
		if err:=rows.Scan(&distance.ID, &distance.FromPlace, &distance.ToPlace, &distance.Distance); err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		distances = append(distances, distance)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(distances)
}


// DeleteDistances handles DELETE /distances/{id}
// @Summary DELETE distance
// @Description DELETE distance by id
// @Tags Place
// @Produce json
// @Security BearerAuth
// @Param id path string true "Distance ID"
// @Router /protected/distances/{id} [DELETE]
func (h *App) DeleteDistances(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	query := "DELETE FROM place_distances WHERE id = ? "
	result, err := h.DB.Exec(query, id)
	if err != nil {
		http.Error(w, "Failed to delete distance", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "distance not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}



// DeletePlace handles DELETE /places/{id}
// @Summary DELETE place
// @Description DELETE place by id
// @Tags Place
// @Produce json
// @Security BearerAuth
// @Param id path string true "Place ID"
// @Router /protected/places/{id} [DELETE]
func (h *App) DeletePlace(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	query := "DELETE FROM places WHERE id = ? "
	result, err := h.DB.Exec(query, id)
	if err != nil {
		http.Error(w, "Failed to delete place", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "place not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}


//CreatePlace handles POST /places
// @Summary POST a new place
// @Description add a new place
// @Tags Place
// @Accept json
// @Produce json
// @Param body body handlers.PlacesRequest true "Announcement"
// @Router /places [post]
func (h *App) CreatePlace(w http.ResponseWriter, r *http.Request) {
	var place models.Place
	if err:=json.NewDecoder(r.Body).Decode(&place); err !=nil {
		http.Error(w, "Invalid Input", http.StatusBadRequest)
		return
	}

	query := "INSERT INTO places (name) VALUES (?)"
	result, err := h.DB.Exec(query, place.Name)
	if err != nil {
		http.Error(w, "Failed to create a place", http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	place.ID = int(id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(place)
}

// createDistance godoc
// @Summary Add a new distance between two places
// @Description Add a new distance between two places
// @Tags Place
// @Accept json
// @Produce json
// @Param distance body models.DistanceRequest true "Distance data"
// @Router /distances [post]
func (h *App) createDistance(w http.ResponseWriter, r *http.Request) {
	var req models.DistanceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.FromPlaceID == 0 || req.ToPlaceID == 0 || req.Distance <= 0 {
		http.Error(w, `{"error": "from_place_id, to_place_id, and distance are required and distance must be positive"}`, http.StatusBadRequest)
		return
	}

	// Validate that from_place and to_place exist
	var count int
	err := h.DB.QueryRow("SELECT COUNT(*) FROM places WHERE id = ? OR id = ?", req.FromPlaceID, req.ToPlaceID).Scan(&count)
	if err != nil || count != 2 {
		http.Error(w, `{"error": "Invalid from_place or to_place ID"}`, http.StatusBadRequest)
		return
	}

	// Insert into place_distances
	_, err1 := h.DB.Exec("INSERT INTO place_distances (to_place, from_place, distance) VALUES (?, ?, ?)",
	req.FromPlaceID, req.ToPlaceID, req.Distance)
	
	result, err := h.DB.Exec("INSERT INTO place_distances (from_place, to_place, distance) VALUES (?, ?, ?)",
		req.FromPlaceID, req.ToPlaceID, req.Distance)
	if err != nil || err1 != nil {
		http.Error(w, `{"error": "Failed to create distance"}`, http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	response := models.DistanceResponse{
		ID:          int(id),
		FromPlaceID: req.FromPlaceID,
		ToPlaceID:   req.ToPlaceID,
		Distance:    req.Distance,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}