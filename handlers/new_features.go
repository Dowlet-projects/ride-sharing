package handlers 


import (

	"net/http"
	"encoding/json"
	//"ride-sharing/models"
	"fmt"
	// "github.com/gorilla/mux"
	"strconv"
	// "ride-sharing/utils"
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
