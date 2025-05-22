package handlers 


import (

	"net/http"
	"encoding/json"
	"ride-sharing/models"
	"fmt"
	// "github.com/gorilla/mux"
	"strconv"
	"ride-sharing/utils"
	"strings"
)

type Taxist struct {
	ID int `json:"id"`
	FullName string `json:"full_name"`
	Phone string `json:"phone"`
}


// GetAllTaxists handles GET /protected/taxists
// @Summary Get all taxists
// @Description Retrieve a paginated list of taxists filtered by optional query parameters
// @Tags Announcement
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (default: 1)" example=1
// @Param limit query int false "Number of items per page (default: 10, max: 100)" example=10
// @Param search query string false "Search through taxist name" example="Amanow Aman"
// @Router /protected/taxists [get]
func (h *App) GetAllTaxists(w http.ResponseWriter, r *http.Request) {
	const defaultPage = 1
	const defaultLimit = 10
	const maxLimit = 100

	// Parse query parameters
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	search := r.URL.Query().Get("search")

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
	query := `SELECT id, full_name, phone FROM taxists `
	countQuery := "SELECT COUNT(*) FROM taxists"
	var conditions []string
	var args []interface{}

	
	if search != "" {
		conditions = append(conditions, "full_name LIKE ? ")
		args = append(args, "%"+search+"%")
	}

	// Append WHERE clause if there are conditions
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
		countQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	// Append LIMIT and OFFSET
	query += " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	// Execute the query
	rows, err := h.DB.Query(query, args...)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	
	var taxists []Taxist = []Taxist{}
	for rows.Next() {
		var taxist Taxist
		if err := rows.Scan(&taxist.ID, &taxist.FullName, &taxist.Phone); err != nil {
			fmt.Println(err)
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		taxists = append(taxists, taxist)
	}

	// Check for errors during row iteration
	if err := rows.Err(); err != nil {
		fmt.Println(err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Get total count with filters
	var totalTaxists int
	err = h.DB.QueryRow(countQuery, args[:len(args)-2]...).Scan(&totalTaxists)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Prepare response
	response := struct {
		Taxists     []Taxist `json:"taxists"`
		Page        int           `json:"page"`
		Limit       int           `json:"limit"`
		Total       int           `json:"total"`
		TotalPages  int           `json:"total_pages"`
	}{
		Taxists:    taxists,
		Page:       page,
		Limit:      limit,
		Total:      totalTaxists,
		TotalPages: (totalTaxists + limit - 1) / limit,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Message


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

type PasMessage struct {
	ID int `json:"id"`
	PassengerID int `json:"passenger_id"`
	FullName string `json:"full_name"`
	Phone string `json:"phone"`
	Message string `json:"message"`
}

// GetAllPassengerMessage handles GET /protected/passenger-messages
// @Summary Get all passenger-messages
// @Description Retrieve all passenger-messages
// @Tags Message
// @Produce json
// @Security BearerAuth
// @Router /protected/passenger-messages [get]
func (h *App) GetAllPassengerMessage(w http.ResponseWriter, r *http.Request) {

	rows, err := h.DB.Query(`SELECT id, passenger_id, full_name, phone, message  
								FROM passenger_messages_full `)

	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	defer rows.Close()
    
	var passenger_messages []PasMessage = []PasMessage{}

	for rows.Next() {
		var message PasMessage
		if err := rows.Scan(&message.ID, &message.PassengerID, &message.FullName, 
			&message.Phone, &message.Message); err != nil {
			fmt.Println(err)
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		passenger_messages = append(passenger_messages, message)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(passenger_messages)
}



type TaxistMessage struct {
	Message string `json:"message"`
}


//CreateMessageTaxist handles POST /protected/taxist-message
// @Summary POST a taxist message
// @Description add a taxist message
// @Tags Message
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body TaxistMessage true "Message"
// @Router /protected/taxist-message [post]
func (h *App) CreateMessageTaxist(w http.ResponseWriter, r *http.Request) {
	var message TaxistMessage
	if err:=json.NewDecoder(r.Body).Decode(&message); err !=nil {
		http.Error(w, "Invalid Input", http.StatusBadRequest)
		return
	}

	claims, ok := r.Context().Value("claims").(*models.Claims)
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "Invalid claims")
		return
	}

	taxist_id := claims.UserID

	query := "INSERT INTO taxist_messages (taxist_id, message) VALUES (?, ?)"
	_, err := h.DB.Exec(query, taxist_id, message.Message)
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

type TaxMessage struct {
	ID int `json:"id"`
	TaxistID int `json:"taxist_id"`
	FullName string `json:"full_name"`
	Phone string `json:"phone"`
	Message string `json:"message"`
}

// GetAllTaxistMessage handles GET /protected/taxist-messages
// @Summary Get all taxist-messages
// @Description Retrieve all taxist-messages
// @Tags Message
// @Produce json
// @Security BearerAuth
// @Router /protected/taxist-messages [get]
func (h *App) GetAllTaxistMessage(w http.ResponseWriter, r *http.Request) {

	rows, err := h.DB.Query(`SELECT id, taxist_id, full_name, phone, message  
								FROM taxist_messages_full `)

	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	defer rows.Close()
    
	var taxist_messages []TaxMessage = []TaxMessage{}

	for rows.Next() {
		var message TaxMessage
		if err := rows.Scan(&message.ID, &message.TaxistID, &message.FullName, 
			&message.Phone, &message.Message); err != nil {
			fmt.Println(err)
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		taxist_messages = append(taxist_messages, message)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(taxist_messages)
}
