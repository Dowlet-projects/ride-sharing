package handlers 


import (

	"net/http"
	"encoding/json"
	"ride-sharing/models"
	//"fmt"
	// "github.com/gorilla/mux"
	//"strconv"
	"ride-sharing/utils"
	//"strings"
)


type PassengerMessage struct {
	Message string `json:"message"`
}


//CreateMessagePassenger handles POST /protected/passenger-message
// @Summary POST a passenger message
// @Description add a passenger message
// @Tags Message
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body PassengerMessage true "Message"
// @Router /protected/passenger-message [post]

func (h *App) CreateMessagePassenger(w http.ResponseWriter, r *http.Request) {
	
	var message PassengerMessage
	if err:=json.NewDecoder(r.Body).Decode(&message); err !=nil {
		http.Error(w, "Invalid Input", http.StatusBadRequest)
		return
	}

	claims, ok := r.Context().Value("claims").(*models.Claims)
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "Invalid claims")
		return
	}

	passenger_id := claims.UserID

	query := "INSERT INTO passenger_messages (passenger_id, message) VALUES (?, ?)"
	_, err := h.DB.Exec(query, passenger_id, message.Message)
	if err != nil {
		http.Error(w, "Failed to create a place", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message":"successfully created",
	})
}