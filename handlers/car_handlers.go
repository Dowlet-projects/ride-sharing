package handlers


// import (
// 	"database/sql"
// 	"encoding/json"
// 	"net/http"
// 	//"strconv"

// 	"ride-sharing/models"
// 	//"github.com/gorilla/mux"
// )

// type Handler struct {
// 	DB *sql.DB
// }



// //CreateMake handles POST /makes

// func (h *Handler) CreateMake(w http.ResponseWriter, r *http.Request) {
// 	var make models.Make
// 	if err:=json.NewDecoder(r.Body).Decode(&make); err !=nil {
// 		http.Error(w, "Invalid Input", http.StatusBadRequest)
// 		return
// 	}

// 	query := "INSERT INTO car_makes (name) VALUES (?)"
// 	result, err := h.DB.Exec(query, make.Name)
// 	if err != nil {
// 		http.Error(w, "Failed to create a make", http.StatusInternalServerError)
// 		return
// 	}

// 	id, _ := result.LastInsertId()
// 	make.ID = int(id)
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(make)
// }