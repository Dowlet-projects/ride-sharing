// handlers/handlers.go
// HTTP handlers for the ride-sharing application

package handlers

// import (
// 	"database/sql"
// 	"encoding/json"
// 	"fmt"
// 	"strconv"

// 	//"fmt"
// 	"log"
// 	"net/http"
// 	"regexp"
// 	"time"

// 	//"ride-sharing/config"
// 	"ride-sharing/models"
// 	"ride-sharing/utils"

// 	//"github.com/go-chi/jwtauth/v5"
// 	//"github.com/gorilla/handlers"
// 	"github.com/gorilla/mux"
// 	//"github.com/rs/cors"
// )

// type App struct {
// 	DB     *sql.DB
// 	Config *config.Config
// }

// App holds application dependencies
// type App struct {
// 	DB     *sql.DB
// 	Config *config.Config
// }

// // SetupRouter configures the HTTP router
// // @Summary Setup API router
// // @Description Configures the HTTP router with CORS and routes for passenger, taxist, and protected endpoints.
// // @Tags Internal
// // @Produce json
// func SetupRouter(db *sql.DB, cfg *config.Config) *mux.Router {
//     app := &App{DB: db, Config: cfg}
//     router := mux.NewRouter()

//     // Set up CORS with more permissive settings
//     // cors := handlers.CORS(
//     //     handlers.AllowedOrigins([]string{"*"}),  // Allow all origins in development
//     //     handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"}),
//     //     handlers.AllowedHeaders([]string{
//     //         "Content-Type",
//     //         "Authorization",
//     //         "X-Requested-With",
//     //         "Accept",
//     //         "Origin",
//     //     }),
//     //     handlers.ExposedHeaders([]string{"Content-Length"}),
//     //     handlers.AllowCredentials(),
//     //     handlers.MaxAge(86400), // 24 hours
//     // )

//     // // Apply CORS middleware first
//     // router.Use(cors)

// 	c := cors.New(cors.Options{
// 		AllowedOrigins:   []string{"*"}, // Allow all origins in development
// 		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
// 		AllowedHeaders:   []string{
// 			"Content-Type",
// 			"Authorization",
// 			"X-Requested-With",
// 			"Accept",
// 			"Origin",
// 		},
// 		ExposedHeaders:   []string{"Content-Length"},
// 		AllowCredentials: true,
// 		MaxAge:           86400, // 24 hours
// 	})

// 	// Apply CORS middleware
// 	router.Use(c.Handler)

//     // Define routes
//     router.HandleFunc("/passenger/register", app.handlePassengerRegister).Methods("POST", "OPTIONS")
//     router.HandleFunc("/passenger/login", app.handlePassengerLogin).Methods("POST", "OPTIONS")
//     router.HandleFunc("/verify", app.handleVerifyCode).Methods("POST", "OPTIONS")
//     router.HandleFunc("/taxist/register", app.handleTaxistRegister).Methods("POST", "OPTIONS")
//     router.HandleFunc("/taxist/login", app.handleTaxistLogin).Methods("POST", "OPTIONS")
//     router.HandleFunc("/taxist/verify", app.handleVerifyCode).Methods("POST", "OPTIONS")
//     router.HandleFunc("/makes", app.CreateMake).Methods("POST", "OPTIONS")
//     router.HandleFunc("/makes", app.GetAllMakes).Methods("GET", "OPTIONS")
//     router.HandleFunc("/models", app.CreateModel).Methods("POST", "OPTIONS")
//     router.HandleFunc("/models/{make_id}", app.GetAllModelsById).Methods("GET", "OPTIONS")
// 	router.HandleFunc("/places", app.GetAllPlaces).Methods("GET", "OPTIONS")
//     router.HandleFunc("/places", app.CreatePlace).Methods("POST", "OPTIONS")
// 	router.HandleFunc("/distances", app.createDistance).Methods("POST", "OPTIONS")
// 	router.HandleFunc("/distances", app.GetAllDistances).Methods("GET", "OPTIONS")
//     protected := router.PathPrefix("/protected").Subrouter()
//     protected.Use(app.authMiddleware)
//     protected.HandleFunc("", app.handleProtected).Methods("GET", "OPTIONS")
// 	protected.HandleFunc("/announcements", app.CreateAnnouncement).Methods("POST", "OPTIONS")
// 	protected.HandleFunc("/ugurlar", app.GetAllUgurlar).Methods("GET", "OPTIONS")
// 	protected.HandleFunc("/ugurlar/{ugur_id}", app.GetUgurById).Methods("GET", "OPTIONS")
// 	protected.HandleFunc("/taxists/{taxist_id}", app.TaxistProfile).Methods("GET", "OPTIONS")
// 	protected.HandleFunc("/reserve-passengers/{taxi_ann_id}", app.CreateReservePassengers).Methods("POST", "OPTIONS")
// 	protected.HandleFunc("/reserve-packages/{taxi_ann_id}", app.CreateReservePackages).Methods("POST", "OPTIONS")
// 	protected.HandleFunc("/taxist-rating/{taxist_id}", app.UpdateRatingTaxist).Methods("PUT", "OPTIONS")
// 	protected.HandleFunc("/taxist-comments/{taxist_id}", app.GetAllTaxistComments).Methods("GET","OPTIONS")
// 	protected.HandleFunc("/taxist-comments/{taxist_id}", app.CreateComment).Methods("POST", "OPTIONS")
// 	protected.HandleFunc("/profile", app.Profile).Methods("GET", "OPTIONS")
// 	protected.HandleFunc("/taxist-notifications/{taxist_id}", app.GetAllTaxistNotifications).Methods("GET", "OPTIONS")
// 	protected.HandleFunc("/reverse-details/{reverse_id}", app.ReverseDetails).Methods("GET", "OPTIONS");
// 	protected.HandleFunc("/favourites", app.CreateFavourites).Methods("POST", "OPTIONS")
// 	protected.HandleFunc("/favourites", app.GetAllFavourites).Methods("GET", "OPTIONS")
// 	protected.HandleFunc("/taxist-departed/{taxi_ann_id}", app.UpdateTaxistAnnouncements).Methods("PUT","OPTIONS")
// 	protected.HandleFunc("/taxist-announcements/{departed}", app.GetTaxistAnnouncements).Methods("GET","OPTIONS")
// 	protected.HandleFunc("/distances/{id}", app.DeleteDistances).Methods("DELETE", "OPTIONS")
// 	protected.HandleFunc("/places/{id}", app.DeletePlace).Methods("DELETE", "OPTIONS")
// 	protected.HandleFunc("/makes/{id}", app.DeleteMake).Methods("DELETE", "OPTIONS")
// 	protected.HandleFunc("/models/{model_id}", app.DeleteModel).Methods("DELETE", "OPTIONS")

// 	return router
// }


// //CreateMake handles POST /makes
// // @Summary POST a new car make
// // @Description add a new car make
// // @Tags Car details
// // @Accept json
// // @Produce json
// // @Param body body handlers.MakesRequest true "Car details"
// // @Router /makes [post]
// func (h *App) CreateMake(w http.ResponseWriter, r *http.Request) {
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


// // DeleteMake handles DELETE /makes/{id}
// // @Summary DELETE make
// // @Description DELETE make by id
// // @Tags Car details
// // @Produce json
// // @Security BearerAuth
// // @Param id path string true "Make ID"
// // @Router /protected/makes/{id} [DELETE]
// func (h *App) DeleteMake(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	id, err := strconv.Atoi(vars["id"])
// 	if err != nil {
// 		http.Error(w, "Invalid ID", http.StatusBadRequest)
// 		return
// 	}

// 	query := "DELETE FROM car_makes WHERE id = ? "
// 	result, err := h.DB.Exec(query, id)
// 	if err != nil {
// 		fmt.Println(err)
// 		http.Error(w, "Failed to delete make", http.StatusInternalServerError)
// 		return
// 	}

// 	rowsAffected, _ := result.RowsAffected()
// 	if rowsAffected == 0 {
// 		http.Error(w, "make not found", http.StatusNotFound)
// 		return
// 	}

// 	w.WriteHeader(http.StatusNoContent)
// }

// // DeleteModel handles DELETE /models/{model_id}
// // @Summary DELETE model
// // @Description DELETE model by id
// // @Tags Car details
// // @Produce json
// // @Security BearerAuth
// // @Param model_id path string true "Model ID"
// // @Router /protected/models/{model_id} [DELETE]
// func (h *App) DeleteModel(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	id, err := strconv.Atoi(vars["model_id"])
// 	if err != nil {
// 		http.Error(w, "Invalid ID", http.StatusBadRequest)
// 		return
// 	}

// 	query := "DELETE FROM car_models WHERE id = ? "
// 	result, err := h.DB.Exec(query, id)
// 	if err != nil {
// 		fmt.Println(err)
// 		http.Error(w, "Failed to delete model", http.StatusInternalServerError)
// 		return
// 	}

// 	rowsAffected, _ := result.RowsAffected()
// 	if rowsAffected == 0 {
// 		http.Error(w, "model not found", http.StatusNotFound)
// 		return
// 	}

// 	w.WriteHeader(http.StatusNoContent)
// }


// // Accept json

// //GetAllMakes handles GET /makes
// // @Summary GET car makes
// // @Description get all car makes
// // @Tags Car details
// // @Produce json
// // @Router /makes [get]
// func (h *App) GetAllMakes(w http.ResponseWriter, r *http.Request) {
// 	rows, err := h.DB.Query("SELECT id, name FROM car_makes")

// 	if err != nil {
// 		http.Error(w, "Server error", http.StatusInternalServerError)
// 		return
// 	}

// 	defer rows.Close()

// 	var makes []models.Make = []models.Make{}
// 	for rows.Next() {
// 		var make models.Make
// 		if err:=rows.Scan(&make.ID, &make.Name); err != nil {
// 			http.Error(w, "Server error", http.StatusInternalServerError)
// 			return
// 		}

// 		makes = append(makes, make)
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(makes)
// }


// //CreateModel handles POST /models
// // @Summary POST a new car model
// // @Description add a new car model
// // @Tags Car details
// // @Accept json
// // @Produce json
// // @Param body body handlers.ModelsRequest true "Car details"
// // @Router /models [post]
// func (h *App) CreateModel(w http.ResponseWriter, r *http.Request) {
// 	var model models.Model
// 	if err:=json.NewDecoder(r.Body).Decode(&model); err !=nil {
// 		http.Error(w, "Invalid Input", http.StatusBadRequest)
// 		return
// 	}

// 	query := "INSERT INTO car_models (name, make_id) VALUES (?, ?)"
// 	result, err := h.DB.Exec(query, model.Name, model.MakeID)
// 	if err != nil {
// 		http.Error(w, "Failed to create a model", http.StatusInternalServerError)
// 		return
// 	}

// 	id, _ := result.LastInsertId()
// 	model.ID = int(id)
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(model)
// }
// // Accept json

// // GetAllModelsById handles GET /models/{make_id}
// // @Summary Get car models by make ID
// // @Description Retrieve all car models for a given make ID
// // @Tags Car details
// // @Produce json
// // @Param make_id path string true "Make ID"
// // @Router /models/{make_id} [get]
// func (h *App) GetAllModelsById(w http.ResponseWriter, r *http.Request) {
	
// 	vars := mux.Vars(r)
// 	id, err := strconv.Atoi(vars["make_id"])

// 	if err != nil {
// 		http.Error(w, "Invalid ID", http.StatusBadRequest)
// 		return
// 	}


// 	rows, err := h.DB.Query("SELECT id, name, make_id FROM car_models WHERE make_id = ?", id)

// 	if err != nil {
// 		http.Error(w, "Server error", http.StatusInternalServerError)
// 		return
// 	}
	
// 	defer rows.Close()
// 	var modelss []models.Model = []models.Model{} 
// 	for rows.Next() {
// 		var model models.Model
// 		if err:=rows.Scan(&model.ID, &model.Name, &model.MakeID); err != nil {
// 			http.Error(w, "Server error", http.StatusInternalServerError)
// 			return
// 		}

// 		modelss = append(modelss, model)
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(modelss)
// }


// Accept json

// //GetAllPlaces handles GET /places
// // @Summary GET car places
// // @Description get all places
// // @Tags Place
// // @Produce json
// // @Router /places [get]
// func (h *App) GetAllPlaces(w http.ResponseWriter, r *http.Request) {
// 	rows, err := h.DB.Query("SELECT id, name FROM places")

// 	if err != nil {
// 		http.Error(w, "Server error", http.StatusInternalServerError)
// 		return
// 	}

// 	defer rows.Close()

// 	var places []models.Place = []models.Place{}
// 	for rows.Next() {
// 		var place models.Place
// 		if err:=rows.Scan(&place.ID, &place.Name); err != nil {
// 			http.Error(w, "Server error", http.StatusInternalServerError)
// 			return
// 		}

// 		places = append(places, place)
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(places)
// }

// type Distance struct {
// 	ID int `json:"id"`
// 	FromPlace string `json:"from_place"`
// 	ToPlace string `json:"to_place"`
// 	Distance string `json:"distance"`
// }


// //GetAllDistances handles GET /distances
// // @Summary GET distances
// // @Description get all distances
// // @Tags Place
// // @Produce json
// // @Router /distances [get]
// func (h *App) GetAllDistances(w http.ResponseWriter, r *http.Request) {
// 	rows, err := h.DB.Query("SELECT id, from_place, to_place, distance FROM distances")

// 	if err != nil {
// 		http.Error(w, "Server error", http.StatusInternalServerError)
// 		return
// 	}

// 	defer rows.Close()

// 	var distances []Distance = []Distance{}
// 	for rows.Next() {
// 		var distance Distance
// 		if err:=rows.Scan(&distance.ID, &distance.FromPlace, &distance.ToPlace, &distance.Distance); err != nil {
// 			http.Error(w, "Server error", http.StatusInternalServerError)
// 			return
// 		}

// 		distances = append(distances, distance)
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(distances)
// }


// // DeleteDistances handles DELETE /distances/{id}
// // @Summary DELETE distance
// // @Description DELETE distance by id
// // @Tags Place
// // @Produce json
// // @Security BearerAuth
// // @Param id path string true "Distance ID"
// // @Router /protected/distances/{id} [DELETE]
// func (h *App) DeleteDistances(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	id, err := strconv.Atoi(vars["id"])
// 	if err != nil {
// 		http.Error(w, "Invalid ID", http.StatusBadRequest)
// 		return
// 	}

// 	query := "DELETE FROM place_distances WHERE id = ? "
// 	result, err := h.DB.Exec(query, id)
// 	if err != nil {
// 		http.Error(w, "Failed to delete distance", http.StatusInternalServerError)
// 		return
// 	}

// 	rowsAffected, _ := result.RowsAffected()
// 	if rowsAffected == 0 {
// 		http.Error(w, "distance not found", http.StatusNotFound)
// 		return
// 	}

// 	w.WriteHeader(http.StatusNoContent)
// }



// // DeletePlace handles DELETE /places/{id}
// // @Summary DELETE place
// // @Description DELETE place by id
// // @Tags Place
// // @Produce json
// // @Security BearerAuth
// // @Param id path string true "Place ID"
// // @Router /protected/places/{id} [DELETE]
// func (h *App) DeletePlace(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	id, err := strconv.Atoi(vars["id"])
// 	if err != nil {
// 		http.Error(w, "Invalid ID", http.StatusBadRequest)
// 		return
// 	}

// 	query := "DELETE FROM places WHERE id = ? "
// 	result, err := h.DB.Exec(query, id)
// 	if err != nil {
// 		http.Error(w, "Failed to delete place", http.StatusInternalServerError)
// 		return
// 	}

// 	rowsAffected, _ := result.RowsAffected()
// 	if rowsAffected == 0 {
// 		http.Error(w, "place not found", http.StatusNotFound)
// 		return
// 	}

// 	w.WriteHeader(http.StatusNoContent)
// }


// //CreatePlace handles POST /places
// // @Summary POST a new place
// // @Description add a new place
// // @Tags Place
// // @Accept json
// // @Produce json
// // @Param body body handlers.PlacesRequest true "Announcement"
// // @Router /places [post]
// func (h *App) CreatePlace(w http.ResponseWriter, r *http.Request) {
// 	var place models.Place
// 	if err:=json.NewDecoder(r.Body).Decode(&place); err !=nil {
// 		http.Error(w, "Invalid Input", http.StatusBadRequest)
// 		return
// 	}

// 	query := "INSERT INTO places (name) VALUES (?)"
// 	result, err := h.DB.Exec(query, place.Name)
// 	if err != nil {
// 		http.Error(w, "Failed to create a place", http.StatusInternalServerError)
// 		return
// 	}

// 	id, _ := result.LastInsertId()
// 	place.ID = int(id)
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(place)
// }

// // createDistance godoc
// // @Summary Add a new distance between two places
// // @Description Add a new distance between two places
// // @Tags Place
// // @Accept json
// // @Produce json
// // @Param distance body models.DistanceRequest true "Distance data"
// // @Router /distances [post]
// func (h *App) createDistance(w http.ResponseWriter, r *http.Request) {
// 	var req models.DistanceRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
// 		return
// 	}

// 	// Validate required fields
// 	if req.FromPlaceID == 0 || req.ToPlaceID == 0 || req.Distance <= 0 {
// 		http.Error(w, `{"error": "from_place_id, to_place_id, and distance are required and distance must be positive"}`, http.StatusBadRequest)
// 		return
// 	}

// 	// Validate that from_place and to_place exist
// 	var count int
// 	err := h.DB.QueryRow("SELECT COUNT(*) FROM places WHERE id = ? OR id = ?", req.FromPlaceID, req.ToPlaceID).Scan(&count)
// 	if err != nil || count != 2 {
// 		http.Error(w, `{"error": "Invalid from_place or to_place ID"}`, http.StatusBadRequest)
// 		return
// 	}

// 	// Insert into place_distances
// 	_, err1 := h.DB.Exec("INSERT INTO place_distances (to_place, from_place, distance) VALUES (?, ?, ?)",
// 	req.FromPlaceID, req.ToPlaceID, req.Distance)
	
// 	result, err := h.DB.Exec("INSERT INTO place_distances (from_place, to_place, distance) VALUES (?, ?, ?)",
// 		req.FromPlaceID, req.ToPlaceID, req.Distance)
// 	if err != nil || err1 != nil {
// 		http.Error(w, `{"error": "Failed to create distance"}`, http.StatusInternalServerError)
// 		return
// 	}

// 	id, _ := result.LastInsertId()
// 	response := models.DistanceResponse{
// 		ID:          int(id),
// 		FromPlaceID: req.FromPlaceID,
// 		ToPlaceID:   req.ToPlaceID,
// 		Distance:    req.Distance,
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(response)
// }

//CreateAnnouncement handles POST /protected/announcements
// @Summary POST a new Ugur for taxist
// @Description add a ugur
// @Tags Announcement
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body models.TaxistAnnouncmentRequest true "Ugur"
// @Router /protected/announcements [post]
// func (h *App) CreateAnnouncement(w http.ResponseWriter, r *http.Request) {
// 	claims, ok := r.Context().Value("claims").(*models.Claims)
// 	if !ok {
// 		utils.RespondError(w, http.StatusUnauthorized, "Invalid claims")
// 		return
// 	}

// 	userID := claims.UserID
// 	fmt.Println(userID)
// 	var taxistAnn models.TaxistAnnouncement
// 	if err:=json.NewDecoder(r.Body).Decode(&taxistAnn); err !=nil {
// 		fmt.Println(err)
// 		http.Error(w, "Invalid Input", http.StatusBadRequest)
// 		return
// 	}

// 	query := "INSERT INTO taxist_announcements (taxist_id, depart_date, depart_time, from_place, to_place, space, type) VALUES (?, ?, ?, ?, ?, ?, ?)"
// 	_, err := h.DB.Exec(query, userID, taxistAnn.DepartDate, taxistAnn.DepartTime, taxistAnn.FromPlaceID, taxistAnn.ToPlaceID, taxistAnn.Space,  taxistAnn.Type)
// 	if err != nil {
// 		fmt.Println(err)
// 		http.Error(w, "Failed to create a announcement", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(map[string]string{
// 		"message":"successfully created",
// 	})
// }

// type PassengerPeople struct {
// 	FullName string `json:"full_name"`
// 	Phone string `json:"phone"`
// }


// type  PostReservePassengers struct {
// 	Passengers []PassengerPeople `json:"passengers" validate:"required"`
// 	Package string `json:"package" validate:"required"`
// }

// //CreateReservePassengers handles POST /protected/reserve-passengers
// // @Summary POST a new reserving passenger for passenger user
// // @Description add a new reserving passenger
// // @Tags Announcement
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param taxi_ann_id path string true "Ugur ID"
// // @Param body body PostReservePassengers true "Reserve Passenger"
// // @Router /protected/reserve-passengers/{taxi_ann_id} [post]
// func (h *App) CreateReservePassengers(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	int_taxi_ann_id, err := strconv.Atoi(vars["taxi_ann_id"])

// 	if err != nil {
// 		http.Error(w, "Invalid id", http.StatusInternalServerError)
// 		fmt.Println(err)
// 		return
// 	}

// 	claims, ok := r.Context().Value("claims").(*models.Claims)
// 	if !ok {
// 		utils.RespondError(w, http.StatusUnauthorized, "Invalid claims")
// 		return
// 	}

// 	userID := claims.UserID
	
// 	var postReservePassengers PostReservePassengers
// 	if err:=json.NewDecoder(r.Body).Decode(&postReservePassengers); err !=nil {
// 		http.Error(w, "Invalid Input", http.StatusBadRequest)
// 		return
// 	}

// 	length := len(postReservePassengers.Passengers)
// 	query := "INSERT INTO reserve_passengers (package, taxi_ann_id, who_reserved, count) VALUES (?, ?, ?, ?)"
// 	 result, err := h.DB.Exec(query, postReservePassengers.Package, int_taxi_ann_id, userID, length)
	
// 	 if err != nil {
// 		http.Error(w, "Failed to create a reserved passenger", http.StatusInternalServerError)
// 		return
// 	}

// 	id, _ := result.LastInsertId()

// 	reservedPassengers := postReservePassengers.Passengers
// 	for _, v := range reservedPassengers {
// 		que := "INSERT INTO reserve_passengers_people (full_name, phone, reserve_id, taxi_ann_id) VALUES (?, ?, ?, ?)"
// 		_, err := h.DB.Exec(que, v.FullName, v.Phone, id, int_taxi_ann_id)
// 		if err != nil {
// 			http.Error(w, "Failed to create reserved Passenger", http.StatusInternalServerError)
// 			return
// 		}
// 	}


// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(map[string]string{
// 		"message":"successfully created",
// 	})
// }

// //CreateReservePackages handles POST /protected/reserve-packages
// // @Summary POST a new reserving package for passenger user
// // @Description add a new reserving package
// // @Tags Announcement
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param taxi_ann_id path string true "Ugur ID"
// // @Param body body models.ReservePackages true "Reserve Packages"
// // @Router /protected/reserve-packages/{taxi_ann_id} [post]
// func (h *App) CreateReservePackages(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	int_taxi_ann_id, err := strconv.Atoi(vars["taxi_ann_id"])

// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	claims, ok := r.Context().Value("claims").(*models.Claims)
// 	if !ok {
// 		utils.RespondError(w, http.StatusUnauthorized, "Invalid claims")
// 		return
// 	}

// 	userID := claims.UserID
	
// 	var reservePackages models.ReservePackages
// 	if err:=json.NewDecoder(r.Body).Decode(&reservePackages); err !=nil {
// 		http.Error(w, "Invalid Input", http.StatusBadRequest)
// 		return
// 	}

// 	query := "INSERT INTO reserve_packages (package_sender, package_reciever, sender_phone, reciever_phone, about_package, taxi_ann_id, who_reserved) VALUES (?, ?, ?, ?, ?, ?, ?)"
// 	 _, err = h.DB.Exec(query, reservePackages.PackageSender, reservePackages.PackageReciever, reservePackages.SenderPhone,reservePackages.RecieverPhone,reservePackages.AboutPackage, int_taxi_ann_id, userID)
	
// 	 if err != nil {
// 		http.Error(w, "Failed to create a reserve package", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(map[string]string{
// 		"message":"successfully created",
// 	})
// }


// // Success 200 {object} object{places=[]models.Place,page=int,limit=int,total=int,total_pages=int} "Paginated list of places"
// // Failure 500 {string} string "Server error"

// // GetAllUgurlar handles GET /protected/ugurlar
// // @Summary Get all ugurlar
// // @Description Retrieve a paginated list of all ugurlar
// // @Tags Announcement
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param page query int false "Page number (default: 1)"
// // @Param limit query int false "Number of items per page (default: 10, max: 100)"
// // @Router /protected/ugurlar [get]
// func (h *App) GetAllUgurlar(w http.ResponseWriter, r *http.Request) {

// 	const defaultPage = 1
// 	const defaultLimit = 10
// 	const maxLimit = 100

// 	pageStr := r.URL.Query().Get("page")
// 	limitStr := r.URL.Query().Get("limit")

// 	page, err := strconv.Atoi(pageStr)
// 	if err != nil || page < 1 {
// 		page = defaultPage
// 	}

// 	limit, err := strconv.Atoi(limitStr)

// 	if err != nil || limit < 1 {
// 		limit = defaultLimit
// 	}

// 	if limit > maxLimit {
// 		limit = maxLimit
// 	}

// 	offset := ( page - 1 ) * limit

// 	rows, err := h.DB.Query("SELECT id, taxist_id, depart_date, depart_time, space, distance, type, full_name, car_make, car_model, car_year, car_number, from_place, to_place, rating FROM ugurlar LIMIT ? OFFSET ?", limit, offset)

// 	if err != nil {
// 		fmt.Println(err)
// 		http.Error(w, "Server error", http.StatusInternalServerError)
// 		return
// 	}

// 	defer rows.Close()

// 	var totalAnnouncements int

// 	err = h.DB.QueryRow("SELECT COUNT(*) FROM ugurlar").Scan(&totalAnnouncements)

// 	if err != nil {
// 		fmt.Println(err)
// 		http.Error(w, "Server error", http.StatusInternalServerError)
// 		return
// 	}

// 	var ugurlar []models.Ugur = []models.Ugur{}

// 	for rows.Next() {
// 		var ugur models.Ugur

// 		if err := rows.Scan(&ugur.ID, &ugur.TaxistID, &ugur.DepartDate, &ugur.DepartTime, &ugur.Space, 
// 			&ugur.Distance, &ugur.Type, &ugur.FullName, &ugur.CarMake, &ugur.CarModel, &ugur.CarYear, 
// 			&ugur.CarNumber, &ugur.FromPlace, &ugur.ToPlace, &ugur.Rating); err != nil {
// 				http.Error(w, "Server error", http.StatusInternalServerError)
// 				return
// 		}

// 		ugurlar = append(ugurlar, ugur)
// 	}

// 	response := struct {
// 		Ugurlar []models.Ugur `json:"ugurlar"`
// 		Page int `json:"page"`
// 		Limit int `json:"limit"`
// 		Total int `json:"total"`
// 		TotalPages int `json:"total_pages"`
// 	}{
// 		Ugurlar: ugurlar,
// 		Page: page,
// 		Limit: limit,
// 		Total: totalAnnouncements,
// 		TotalPages: (totalAnnouncements +limit -1) / limit,
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(response)
// }


// // GetUgurById handles GET /ugurlar/{ugur_id}
// // @Summary Get ugur by ugur_id
// // @Description Retrieve ugur by given ugur_id
// // @Tags Announcement
// // @Produce json
// // @Security BearerAuth
// // @Param ugur_id path string true "Ugur ID"
// // @Router /protected/ugurlar/{ugur_id} [get]
// func (h *App) GetUgurById(w http.ResponseWriter, r *http.Request) {
	
// 	vars := mux.Vars(r)
// 	id, err := strconv.Atoi(vars["ugur_id"])

// 	if err != nil {
// 		fmt.Println(err, id)
// 		http.Error(w, "Invalid ID", http.StatusBadRequest)
// 		return
// 	}

// 	var ugur models.UgurDetails
// 	if err := h.DB.QueryRow("SELECT id, taxist_id, depart_date, depart_time, space, distance, type, full_name, car_make, car_model, car_year, car_number, from_place, to_place, rating FROM ugurlar WHERE id = ? ", id).Scan(&ugur.ID, &ugur.TaxistID, &ugur.DepartDate, &ugur.DepartTime, &ugur.Space, 
// 		&ugur.Distance, &ugur.Type, &ugur.FullName, &ugur.CarMake, &ugur.CarModel, &ugur.CarYear, 
// 		&ugur.CarNumber, &ugur.FromPlace, &ugur.ToPlace, &ugur.Rating); err != nil {
// 			fmt.Println(err)
// 			http.Error(w, "Server error", http.StatusInternalServerError)
// 			return
// 	}


// 	rows, err := h.DB.Query("SELECT id, full_name, phone FROM reserve_passengers_people WHERE taxi_ann_id = ?", id)

// 	if err != nil {
// 		http.Error(w, "Server error", http.StatusInternalServerError)
// 		return
// 	}
	
// 	defer rows.Close()

// 	var passengers []models.ReservePassengers = []models.ReservePassengers{} 
// 	for rows.Next() {
// 		var passenger models.ReservePassengers
// 		if err:=rows.Scan(&passenger.ID, &passenger.FullName, &passenger.Phone); err != nil {
// 			http.Error(w, "Server error", http.StatusInternalServerError)
// 			return
// 		}

// 		passengers = append(passengers, passenger)
// 	}

// 	ugur.Passengers = passengers

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(ugur)
// }


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


// // Profile handles GET /profile
// // @Summary Get user profile 
// // @Description Retrieve user
// // @Tags Announcement
// // @Produce json
// // @Security BearerAuth
// // @Router /protected/profile [get]
// func (h *App) Profile(w http.ResponseWriter, r *http.Request) {
	
// 	claims, ok := r.Context().Value("claims").(*models.Claims)

// 	if !ok {
// 		utils.RespondError(w, http.StatusUnauthorized, "Invalid claims")
// 		return
// 	}



// 	UserID := claims.UserID
// 	UserType := claims.UserType

// 	if UserType == "passenger" {
// 		var passenger models.PassengerProfile
	
// 		if err := h.DB.QueryRow("SELECT id, full_name, phone FROM passengers WHERE id = ? ", UserID).Scan(&passenger.ID, &passenger.FullName, 
// 			&passenger.Phone); err != nil {
// 			http.Error(w, "Server error", http.StatusInternalServerError)
// 			return
// 		}
// 	    passenger.UserType = "passenger"
// 		w.Header().Set("Content-Type", "application/json")
// 		json.NewEncoder(w).Encode(passenger)
// 	} else {
// 			var taxist models.Taxist
// 			if err := h.DB.QueryRow("SELECT id, full_name, phone, car_make, car_model, car_year, car_number, rating FROM taxists WHERE id = ? ", UserID).Scan(&taxist.ID, &taxist.FullName, 
// 			&taxist.Phone, &taxist.CarMake, &taxist.CarModel, &taxist.CarYear, &taxist.CarNumber, &taxist.Rating); err != nil {
// 			http.Error(w, "Server error", http.StatusInternalServerError)
// 		return
// 		}
// 		taxist.UserType = "taxist"
// 		w.Header().Set("Content-Type", "application/json")
// 		json.NewEncoder(w).Encode(taxist)
// 	}

// }


// // PassengerRegisterRequest defines the request body for passenger registration
// type PassengerRegisterRequest struct {
// 	FullName string `json:"full_name" example:"John Doe" description:"Full name of the passenger" validate:"required"`
// 	Phone    string `json:"phone" example:"+12345678901" description:"Phone number in international format" validate:"required"`
// }

// // handlePassengerRegister handles passenger registration
// // @Summary Register a new passenger
// // @Description Registers a new passenger with a full name and phone number, sending a verification code. Validates input, checks for duplicate phone numbers, and logs the code (replace with SMS in production).
// // @Tags Passenger
// // @Accept json
// // @Produce json
// // @Param body body handlers.PassengerRegisterRequest true "Passenger registration details"
// // @Router /passenger/register [post]
// func (app *App) handlePassengerRegister(w http.ResponseWriter, r *http.Request) {
// 	var req struct {
// 		FullName string `json:"full_name"`
// 		Phone    string `json:"phone"`
// 	}
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
// 		return
// 	}

// 	// Validate input
// 	if req.FullName == "" {
// 		utils.RespondError(w, http.StatusBadRequest, "Full name is required")
// 		return
// 	}
// 	if !utils.ValidatePhone(req.Phone) {
// 		utils.RespondError(w, http.StatusBadRequest, "Invalid phone number format")
// 		return
// 	}

// 	// Check if phone exists
// 	var exists bool
// 	err := app.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM passengers WHERE phone = ?)", req.Phone).Scan(&exists)
// 	if err != nil {
// 		utils.RespondError(w, http.StatusInternalServerError, "Database error")
// 		return
// 	}
// 	if exists {
// 		utils.RespondError(w, http.StatusConflict, "Phone number already registered")
// 		return
// 	}

// 	// Generate verification code and store with registration data
// 	code, err := utils.GenerateVerificationCode(app.DB, req.Phone, req.FullName, "passenger", "", "", "", 0)
// 	if err != nil {
// 		utils.RespondError(w, http.StatusInternalServerError, "Failed to generate verification code")
// 		return
// 	}

// 	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
// 		"message": "Verification code generated",
// 		"code":    code,
// 	})

// 	// // Insert passenger
// 	// result, err := app.DB.Exec("INSERT INTO passengers (full_name, phone, created_at) VALUES (?, ?, ?)",
// 	// 	req.FullName, req.Phone, time.Now())
// 	// if err != nil {
// 	// 	utils.RespondError(w, http.StatusInternalServerError, "Failed to register passenger")
// 	// 	return
// 	// }
// 	// userID, _ := result.LastInsertId()

// 	// // Generate and store verification code
// 	// code, err := utils.GenerateVerificationCode(app.DB, req.Phone)
// 	// if err != nil {
// 	// 	utils.RespondError(w, http.StatusInternalServerError, "Failed to generate verification code")
// 	// 	return
// 	// }

// 	// // Simulate sending code (replace with Twilio in production)
// 	// log.Printf("Verification code for %s: %s", req.Phone, code)

// 	// utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
// 	// 	"message": "Verification code sent",
// 	// 	"user_id": userID,
// 	// })
// }


// GetAllTaxistComments handles GET /taxist-comments/{taxist_id}
// @Summary Get all specific taxist comments by taxist_id
// @Description Retrieve all taxist comments by given taxist ID
// @Tags Announcement
// @Produce json
// @Security BearerAuth
// @Param taxist_id path string true "Taxist ID"
// @Router /protected/taxist-comments/{taxist_id} [get]
// func (h *App) GetAllTaxistComments(w http.ResponseWriter, r *http.Request) {

// 	vars := mux.Vars(r)
// 	id, err := strconv.Atoi(vars["taxist_id"])

// 	if err != nil {
// 		http.Error(w, "Invalid ID", http.StatusBadRequest)
// 		return
// 	}

// 	rows, err := h.DB.Query("SELECT id, full_name, comment from comments_to_taxist where taxist_id = ?", id)

// 	if err != nil {
// 		http.Error(w, "Server error", http.StatusInternalServerError)
// 		return
// 	}

// 	defer rows.Close()
// 	var comments []models.Comment =[]models.Comment{}
// 	for rows.Next() {
// 		var comment models.Comment
// 		if err := rows.Scan(&comment.ID, &comment.FullName, &comment.Comment); err != nil {
// 			http.Error(w, "Server error", http.StatusInternalServerError)
// 			return
// 		}

// 		comments = append(comments, comment)
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(comments)
// }

// type PostComment struct {
// 	Comment string `json:"comment" example:"" description:"passenger comment"`
// }

// // GetTaxistNotifications handles GET /taxist-notifications/{taxist_id}
// // @Summary Get all specific taxist notifications by taxist_id
// // @Description Retrieve all taxist notifications by given taxist_id
// // @Tags Announcement
// // @Produce json
// // @Security BearerAuth
// // @Param taxist_id path string true "Taxist ID"
// // @Router /protected/taxist-notifications/{taxist_id} [get]
// func (h *App) GetAllTaxistNotifications(w http.ResponseWriter, r *http.Request) {

// 	vars := mux.Vars(r)
// 	id, err := strconv.Atoi(vars["taxist_id"])

// 	if err != nil {
// 		http.Error(w, "Invalid ID", http.StatusBadRequest)
// 		return
// 	}

// 	rows, err := h.DB.Query("SELECT id, taxist_id, full_name, count, created_at FROM taxist_notifications WHERE taxist_id = ?", id)
	
// 	if err != nil {
// 		http.Error(w, "Server error", http.StatusInternalServerError)
// 		return
// 	}

// 	defer rows.Close()
// 	var notifications []models.Notification = []models.Notification{}

// 	for rows.Next() {
// 		var notification models.Notification
// 		if err := rows.Scan(&notification.ID, &notification.TaxistID, &notification.FullName, &notification.Count, &notification.CreatedAt); err != nil {
// 			http.Error(w, "Server error", http.StatusInternalServerError)
// 			return
// 		}

// 		notifications = append(notifications, notification)
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(notifications)

// }

// // GetTaxistAnnouncements handles GET /taxist-announcements/{departed}
// // @Summary Get all specific taxist announcements by departed
// // @Description Retrieve all taxist announcements by given departed
// // @Tags Announcement
// // @Produce json
// // @Security BearerAuth
// // @Param departed path string true "Departed"
// // @Router /protected/taxist-announcements/{departed} [get]
// func (h *App) GetTaxistAnnouncements(w http.ResponseWriter, r *http.Request) {

// 	vars := mux.Vars(r)
// 	departed, err := strconv.Atoi(vars["departed"])

// 	if err != nil {
// 		http.Error(w, "Invalid ID", http.StatusBadRequest)
// 		return
// 	}

// 	claims, ok := r.Context().Value("claims").(*models.Claims)

// 	if !ok {
// 		utils.RespondError(w, http.StatusUnauthorized, "Invalid claims")
// 		return
// 	}

// 	rows, err := h.DB.Query("SELECT id, taxist_id, depart_date, depart_time, space, distance, type, full_name, car_make, car_model, car_year, car_number, from_place, to_place, rating FROM ugurlar WHERE taxist_id = ? AND departed = ? ", claims.UserID, departed)

// 	if err != nil {
// 		fmt.Println(err)
// 		http.Error(w, "Server error", http.StatusInternalServerError)
// 		return
// 	}

// 	defer rows.Close()

// 	var ugurlar []models.Ugur = []models.Ugur{}

// 	for rows.Next() {
// 		var ugur models.Ugur

// 		if err := rows.Scan(&ugur.ID, &ugur.TaxistID, &ugur.DepartDate, &ugur.DepartTime, &ugur.Space, 
// 			&ugur.Distance, &ugur.Type, &ugur.FullName, &ugur.CarMake, &ugur.CarModel, &ugur.CarYear, 
// 			&ugur.CarNumber, &ugur.FromPlace, &ugur.ToPlace, &ugur.Rating); err != nil {
// 				http.Error(w, "Server error", http.StatusInternalServerError)
// 				return
// 		}

// 		ugurlar = append(ugurlar, ugur)
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(ugurlar)

// }

// type PassengerDetail struct {
// 	ID int `json:"id"`
// 	FullName string `json:"full_name"`
// 	Phone string `json:"phone"`
// }

// type PassengerDetails struct {
// 	Passengers []PassengerDetail `json:"passengers"`
// 	Package string `json:"package"`
// 	SubmitterName string `json:"submitter_name"`
// 	SubmitterPhone string `json:"submitter_phone"`
// 	CreatedAt string `json:"created_at"`
// }

// // ReverseDetails handles GET /reverse-details/{reverse_id}
// // @Summary Get all specific taxist reversed Details by reverse_id
// // @Description Retrieve all taxist reversed Details by given reverse ID
// // @Tags Announcement
// // @Produce json
// // @Security BearerAuth
// // @Param reverse_id path string true "Reverse ID"
// // @Router /protected/reverse-details/{reverse_id} [get]
// func (h *App) ReverseDetails(w http.ResponseWriter, r *http.Request) {

// 	vars := mux.Vars(r)
// 	id, err := strconv.Atoi(vars["reverse_id"])

// 	if err != nil {
// 		http.Error(w, "Invalid ID", http.StatusBadRequest)
// 		return
// 	}

// 	rows, err := h.DB.Query("SELECT id, full_name, phone from reserve_passengers_people where reserve_id = ? ", id)

// 	if err != nil {
// 		http.Error(w, "Server error", http.StatusInternalServerError)
// 		return
// 	}

// 	defer rows.Close()
// 	var passengers []PassengerDetail = []PassengerDetail{}
// 	for rows.Next() {
// 		var passenger PassengerDetail
// 		if err := rows.Scan(&passenger.ID, &passenger.FullName, &passenger.Phone); err != nil {
// 			http.Error(w, "Server error", http.StatusInternalServerError)
// 			return
// 		}
// 		passengers = append(passengers, passenger)
// 	}

// 	var passengerDetails PassengerDetails

// 	if err := h.DB.QueryRow("SELECT package, full_name, phone, created_at FROM view_reverse_passengers where id = ? ", id).Scan(&passengerDetails.Package, &passengerDetails.SubmitterName, &passengerDetails.SubmitterPhone, &passengerDetails.CreatedAt); err != nil {
// 		http.Error(w, "Server error", http.StatusInternalServerError)
// 		return
// 	}

// 	passengerDetails.Passengers = passengers

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(passengerDetails)
// }


// //CreateComment handles POST /protected/taxist-comments
// // @Summary POST a new comment for passenger user
// // @Description add a new comment
// // @Tags Announcement
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param taxist_id path string true "Taxist ID"
// // @Param body body PostComment true "Passenger Comment"
// // @Router /protected/taxist-comments/{taxist_id} [post]
// func (h *App) CreateComment(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	id, err := strconv.Atoi(vars["taxist_id"])

// 	if err != nil {
// 		http.Error(w, "Invalid ID", http.StatusInternalServerError)
// 		return
// 	}

// 	claims, ok := r.Context().Value("claims").(*models.Claims)

// 	if !ok {
// 		utils.RespondError(w, http.StatusUnauthorized, "Invalid claims")
// 		return
// 	}

// 	passengerID := claims.UserID

// 	var postComment PostComment
// 	if err := json.NewDecoder(r.Body).Decode(&postComment); err != nil {
// 		http.Error(w, "Invalid Input", http.StatusBadRequest)
// 		return
// 	}

// 	query := "INSERT INTO taxist_comments (taxist_id, passenger_id, comment) VALUES (?, ?, ?)"
// 	_, err = h.DB.Exec(query, id, passengerID, postComment.Comment)

// 	if err != nil {
// 		http.Error(w, "Failed to create a comment", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(map[string]string{
// 		"message":"successfully created",
// 	})
// }

// type Favourite struct {
// 	TaxistID int `json:"taxist_id"`
// }

// // CreateFavourites handles POST /protected/favourites
// // @Summary POST a new passenger favourite
// // @Description add a new passenger favourite
// // @Tags Announcement
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param body body Favourite true "Announcement"
// // @Router /protected/favourites [post]
// func (h *App) CreateFavourites(w http.ResponseWriter, r *http.Request) {
	
// 	var favourite Favourite

// 	if err := json.NewDecoder(r.Body).Decode(&favourite); err != nil {
// 		http.Error(w, "Invalid Input", http.StatusBadRequest)
// 		return
// 	}

// 	claims, ok := r.Context().Value("claims").(*models.Claims)
// 	if !ok {
// 		utils.RespondError(w, http.StatusUnauthorized, "Invalid claims")
// 		return
// 	}
	
// 	passenger_id := claims.UserID

// 	query := "CALL insert_halanlarym(?, ?)"
// 	_, err := h.DB.Exec(query, favourite.TaxistID, passenger_id)

// 	if err != nil {
// 		http.Error(w, "Failed to add or delete action", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(map[string]string{
// 		"message":"favourite taxist deleted or added",
// 	})
// }

// type Halanym struct {
// 	ID int `json:"id"`
// 	TaxistID int `json:"taxist_id"`
// 	FullName string `json:"full_name"`
// 	CarMake string `json:"car_make"`
// 	CarModel string `json:"car_model"`
// 	CarYear string `json:"car_year"`
// 	CarNumber string `json:"car_number"`
// 	Rating float32 `json:"rating"`
// }

// // GetAllFavourites handles GET /favourites
// // @Summary Get all specific passenger favourites
// // @Description Retrieve all passenger favourites
// // @Tags Announcement
// // @Produce json
// // @Security BearerAuth
// // @Router /protected/favourites [get]
// func (h *App) GetAllFavourites(w http.ResponseWriter, r *http.Request) {
// 	claims, ok := r.Context().Value("claims").(*models.Claims)
// 	if !ok {
// 		utils.RespondError(w, http.StatusUnauthorized, "Invalid claims")
// 		return
// 	}
	
// 	passenger_id := claims.UserID
// 	rows, err := h.DB.Query("SELECT id, taxist_id, full_name, car_make, car_model, car_year, car_number, rating FROM halanlarym WHERE passenger_id = ? ", passenger_id)

// 	if err != nil {
// 		http.Error(w, "Server error", http.StatusInternalServerError)
// 		return
// 	}

// 	defer rows.Close()
    
// 	var halanlarym []Halanym = []Halanym{}

// 	for rows.Next() {
// 		var halanym Halanym
// 		if err := rows.Scan(&halanym.ID, &halanym.TaxistID, &halanym.FullName, &halanym.CarMake, &halanym.CarModel, &halanym.CarYear, &halanym.CarNumber, &halanym.Rating); err != nil {
// 			fmt.Println(err)
// 			http.Error(w, "Server error", http.StatusInternalServerError)
// 			return
// 		}

// 		halanlarym = append(halanlarym, halanym)
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(halanlarym)
// }


// // UpdateTaxistAnnouncements
// // @Summary Update taxist announcement department
// // @Description Updates a taxist's announcement department by taxi_ann_id
// // @Tags Announcement
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param taxi_ann_id path int true "Announcement ID"
// // @Router /protected/taxist-departed/{taxi_ann_id} [put]
// func (h *App) UpdateTaxistAnnouncements(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	id, err := strconv.Atoi(vars["taxi_ann_id"])
// 	if err != nil {
// 		http.Error(w, "Invalid taxi announcement ID", http.StatusBadRequest)
// 		return
// 	}
	
// 	claims, ok := r.Context().Value("claims").(*models.Claims)
// 	if !ok {
// 		utils.RespondError(w, http.StatusUnauthorized, "Invalid claims")
// 		return
// 	}

// 	query := "UPDATE taxist_announcements SET departed = 1 WHERE taxist_id = ? AND id = ?"
// 	result, err := h.DB.Exec(query, claims.UserID, id)
// 	if err != nil {
// 		http.Error(w, "Failed to update taxist rating", http.StatusInternalServerError)
// 		return
// 	}

// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		http.Error(w, "Error checking update", http.StatusInternalServerError)
// 		return
// 	}

// 	if rowsAffected == 0 {
// 		http.Error(w, "User not found", http.StatusNotFound)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(map[string]string{
// 		"message":"successfully updated",
// 	})

// }





// // TaxistRegisterRequest defines the request body for taxist registration
// type TaxistRegisterRequest struct {
// 	FullName  string `json:"full_name" example:"Jane Smith" description:"Full name of the taxist" validate:"required"`
// 	Phone     string `json:"phone" example:"+12345678902" description:"Phone number in international format" validate:"required,phone"`
// 	CarMake   string `json:"car_make" example:"Toyota" description:"Vehicle manufacturer" validate:"required"`
// 	CarModel  string `json:"car_model" example:"Camry" description:"Vehicle model" validate:"required"`
// 	CarYear   int    `json:"car_year" example:"2020" description:"Vehicle manufacturing year" validate:"required,gte=1900,lte=2026"`
// 	CarNumber string `json:"car_number" example:"ABC123" description:"Vehicle license plate number" validate:"required"`
// }

// type MakesRequest struct {
// 	Name  string `json:"name" example:"Toyota" description:"name of car make" validate:"required"`
// }

// type PlacesRequest struct {
// 	Name  string `json:"name" example:"Ashgabat" description:"name of place" validate:"required"`
// }

// type ModelsRequest struct {
// 	Name  string `json:"name" example:"Camry" description:"name of car model" validate:"required"`
// 	MakeId int `json:"make_id" example:"0" description:"car make id" validate:"required"`
// }

// type TaxistRating struct {
// 	Rating float32 `json:"rating" example:"4" description:"update rating of taxist" validate:"required"`
// }

// // UpdateRatingTaxist godoc
// // @Summary Update taxist rating
// // @Description Updates a taxist's rating by taxist_id
// // @Tags Announcement
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param taxist_id path int true "Taxist ID"
// // @Param taxist_rating body handlers.TaxistRating true "Taxist rating details"
// // @Router /protected/taxist-rating/{taxist_id} [put]
// func (h *App) UpdateRatingTaxist(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	id, err := strconv.Atoi(vars["taxist_id"])
// 	if err != nil {
// 		http.Error(w, "Invalid taxist ID", http.StatusBadRequest)
// 		return
// 	}

// 	var rating TaxistRating
// 	if err:= json.NewDecoder(r.Body).Decode(&rating); err != nil {
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}

// 	defer r.Body.Close()

// 	query := "CALL ratingPut(?, ?)"
// 	result, err := h.DB.Exec(query, id, rating.Rating)
// 	if err != nil {
// 		http.Error(w, "Failed to update taxist rating", http.StatusInternalServerError)
// 		return
// 	}

// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		http.Error(w, "Error checking update", http.StatusInternalServerError)
// 		return
// 	}

// 	if rowsAffected == 0 {
// 		http.Error(w, "User not found", http.StatusNotFound)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(map[string]string{
// 		"message":"successfully updated",
// 	})

// }


// // handleTaxistRegister handles taxist registration
// // @Summary Register a new taxist
// // @Description Registers a new taxist with personal and vehicle details, sending a verification code. Validates input, checks for duplicate phone numbers, and logs the code.
// // @Tags Taxist
// // @Accept json
// // @Produce json
// // @Param body body handlers.TaxistRegisterRequest true "Taxist registration details"
// // @Router /taxist/register [post]
// func (app *App) handleTaxistRegister(w http.ResponseWriter, r *http.Request) {
// 	var req struct {
// 		FullName  string `json:"full_name"`
// 		Phone     string `json:"phone"`
// 		CarMake   string `json:"car_make"`
// 		CarModel  string `json:"car_model"`
// 		CarYear   int    `json:"car_year"`
// 		CarNumber string `json:"car_number"`
// 	}
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
// 		return
// 	}

// 	// Validate input
// 	if req.FullName == "" || req.CarMake == "" || req.CarModel == "" ||
// 		req.CarNumber == "" || req.CarYear < 1900 || req.CarYear > time.Now().Year()+1 {
// 		utils.RespondError(w, http.StatusBadRequest, "All fields are required and must be valid")
// 		return
// 	}
// 	if !utils.ValidatePhone(req.Phone) {
// 		utils.RespondError(w, http.StatusBadRequest, "Invalid phone number format")
// 		return
// 	}

// 	// Check if phone exists
// 	var exists bool
// 	err := app.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM taxists WHERE phone = ?)", req.Phone).Scan(&exists)
// 	if err != nil {
// 		utils.RespondError(w, http.StatusInternalServerError, "Database error")
// 		return
// 	}
// 	if exists {
// 		utils.RespondError(w, http.StatusConflict, "Phone number already registered")
// 		return
// 	}

// 	// Generate verification code and store with registration data
// 	code, err := utils.GenerateVerificationCode(app.DB, req.Phone, req.FullName, "taxist",
// 		req.CarMake, req.CarModel, req.CarNumber, req.CarYear)
// 	if err != nil {
// 		utils.RespondError(w, http.StatusInternalServerError, "Failed to generate verification code")
// 		return
// 	}

// 	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
// 		"message": "Verification code generated",
// 		"code":    code,
// 	})

// 	// // Insert taxist
// 	// result, err := app.DB.Exec(
// 	// 	"INSERT INTO taxists (full_name, phone, car_make, car_model, car_year, car_number, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
// 	// 	req.FullName, req.Phone, req.CarMake, req.CarModel, req.CarYear, req.CarNumber, time.Now())
// 	// if err != nil {
// 	// 	utils.RespondError(w, http.StatusInternalServerError, "Failed to register taxist")
// 	// 	return
// 	// }
// 	// userID, _ := result.LastInsertId()

// 	// // Generate and store verification code
// 	// code, err := utils.GenerateVerificationCode(app.DB, req.Phone)
// 	// if err != nil {
// 	// 	utils.RespondError(w, http.StatusInternalServerError, "Failed to generate verification code")
// 	// 	return
// 	// }

// 	// // Simulate sending code
// 	// log.Printf("Verification code for %s: %s", req.Phone, code)

// 	// utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
// 	// 	"message": "Verification code sent",
// 	// 	"user_id": userID,
// 	// })
// }

// // LoginRequest defines the request body for login
// type LoginRequest struct {
// 	Phone string `json:"phone" example:"+12345678901" description:"Phone number in international format" validate:"required,phone"`
// }

// // handlePassengerLogin initiates passenger login
// // @Summary Initiate passenger login
// // @Description Initiates login by sending a verification code to the passenger's phone number. Checks if the phone is registered and logs the code.
// // @Tags Passenger
// // @Accept json
// // @Produce json
// // @Param body body handlers.LoginRequest true "Passenger login details"
// // @Router /passenger/login [post]
// func (app *App) handlePassengerLogin(w http.ResponseWriter, r *http.Request) {
// 	var req struct {
// 		Phone string `json:"phone"`
// 	}
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
// 		return
// 	}

// 	// Validate input
// 	if !utils.ValidatePhone(req.Phone) {
// 		utils.RespondError(w, http.StatusBadRequest, "Invalid phone number format")
// 		return
// 	}

// 	// Check if phone exists
// 	var userID int
// 	err := app.DB.QueryRow("SELECT id FROM passengers WHERE phone = ?", req.Phone).Scan(&userID)
// 	if err == sql.ErrNoRows {
// 		utils.RespondError(w, http.StatusNotFound, "Phone number not registered")
// 		return
// 	}
// 	if err != nil {
// 		utils.RespondError(w, http.StatusInternalServerError, "Database error")
// 		return
// 	}

// 	// Generate verification code
// 	code, err := utils.GenerateVerificationCode(app.DB, req.Phone, "", "passenger", "", "", "", 0)
// 	if err != nil {
// 		utils.RespondError(w, http.StatusInternalServerError, "Failed to generate verification code")
// 		return
// 	}

// 	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
// 		"message": "Verification code generated",
// 		"code":    code,
// 		"user_id": userID,
// 	})

// 	// // Generate and store verification code
// 	// code, err := utils.GenerateVerificationCode(app.DB, req.Phone)
// 	// if err != nil {
// 	// 	utils.RespondError(w, http.StatusInternalServerError, "Failed to generate verification code")
// 	// 	return
// 	// }

// 	// // Simulate sending code
// 	// log.Printf("Verification code for %s: %s", req.Phone, code)

// 	// utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
// 	// 	"message": "Verification code sent",
// 	// 	"user_id": userID,
// 	// })
// }

// // handleTaxistLogin initiates taxist login
// // @Summary Initiate taxist login
// // @Description Initiates login by sending a verification code to the taxist's phone number. Checks if the phone is registered and logs the code.
// // @Tags Taxist
// // @Accept json
// // @Produce json
// // @Param body body handlers.LoginRequest true "Taxist login details"
// // @Router /taxist/login [post]
// func (app *App) handleTaxistLogin(w http.ResponseWriter, r *http.Request) {
// 	var req struct {
// 		Phone string `json:"phone"`
// 	}
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
// 		return
// 	}

// 	// Validate input
// 	if !utils.ValidatePhone(req.Phone) {
// 		utils.RespondError(w, http.StatusBadRequest, "Invalid phone number format")
// 		return
// 	}

// 	// Check if phone exists
// 	var userID int
// 	err := app.DB.QueryRow("SELECT id FROM taxists WHERE phone = ?", req.Phone).Scan(&userID)
// 	if err == sql.ErrNoRows {
// 		utils.RespondError(w, http.StatusNotFound, "Phone number not registered")
// 		return
// 	}
// 	if err != nil {
// 		utils.RespondError(w, http.StatusInternalServerError, "Database error")
// 		return
// 	}

// 	// Generate verification code
// 	code, err := utils.GenerateVerificationCode(app.DB, req.Phone, "", "taxist", "", "", "", 0)
// 	if err != nil {
// 		utils.RespondError(w, http.StatusInternalServerError, "Failed to generate verification code")
// 		return
// 	}

// 	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
// 		"message": "Verification code generated",
// 		"code":    code,
// 		"user_id": userID,
// 	})

// 	// // Generate and store verification code
// 	// code, err := utils.GenerateVerificationCode(app.DB, req.Phone)
// 	// if err != nil {
// 	// 	utils.RespondError(w, http.StatusInternalServerError, "Failed to generate verification code")
// 	// 	return
// 	// }

// 	// // Simulate sending code
// 	// log.Printf("Verification code for %s: %s", req.Phone, code)

// 	// utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
// 	// 	"message": "Verification code sent",
// 	// 	"user_id": userID,
// 	// })
// }

// // VerifyCodeRequest defines the request body for verification
// type VerifyCodeRequest struct {
// 	Phone string `json:"phone" example:"+12345678901" description:"Phone number in international format" validate:"required,phone"`
// 	Code  string `json:"code" example:"1234" description:"4-digit verification code" validate:"required,len=4"`
// }

// // handleVerifyCode verifies the code and completes registration/login
// // @Summary Verify phone number
// // @Description Verifies a phone number using a 4-digit code for passengers or taxists, issuing a JWT token upon success. Deletes the used code.
// // @Tags Passenger,Taxist
// // @Accept json
// // @Produce json
// // @Param body body handlers.VerifyCodeRequest true "Verification details"
// // @Router /verify [post]
// func (app *App) handleVerifyCode(w http.ResponseWriter, r *http.Request) {
// 	var req struct {
// 		Phone string `json:"phone"`
// 		Code  string `json:"code"`
// 	}
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
// 		return
// 	}

// 	// Validate input
// 	if !utils.ValidatePhone(req.Phone) {
// 		utils.RespondError(w, http.StatusBadRequest, "Invalid phone number format")
// 		return
// 	}
// 	if len(req.Code) != 4 || !regexp.MustCompile(`^\d{4}$`).MatchString(req.Code) {
// 		utils.RespondError(w, http.StatusBadRequest, "Invalid verification code")
// 		return
// 	}

// 	// Validate code and retrieve registration data
// 	var storedCode, fullName, userType, carMake, carModel, carNumber string
// 	var expiresAt time.Time
// 	var carYear int
// 	var rating float32
// 	err := app.DB.QueryRow(
// 		`SELECT code, expires_at, full_name, user_type, car_make, car_model, car_year, car_number, rating 
// 		 FROM verification_codes WHERE phone = ?`,
// 		req.Phone).Scan(&storedCode, &expiresAt, &fullName, &userType, &carMake, &carModel, &carYear, &carNumber, &rating)
// 	if err == sql.ErrNoRows {
// 		utils.RespondError(w, http.StatusBadRequest, "No verification code found")
// 		return
// 	}
// 	if err != nil {
// 		utils.RespondError(w, http.StatusInternalServerError, "Database error")
// 		return
// 	}

// 	if time.Now().After(expiresAt) {
// 		utils.RespondError(w, http.StatusBadRequest, "Verification code expired")
// 		return
// 	}

// 	if storedCode != req.Code {
// 		utils.RespondError(w, http.StatusBadRequest, "Invalid verification code")
// 		return
// 	}

// 	// Check if user is already registered
// 	var userID int64
// 	var isRegistered bool
// 	var registeredType models.UserType
// 	err = app.DB.QueryRow("SELECT id FROM passengers WHERE phone = ?", req.Phone).Scan(&userID)
// 	if err == nil {
// 		isRegistered = true
// 		registeredType = models.RolePassenger
// 	} else if err == sql.ErrNoRows {
// 		err = app.DB.QueryRow("SELECT id FROM taxists WHERE phone = ?", req.Phone).Scan(&userID)
// 		if err == nil {
// 			isRegistered = true
// 			registeredType = models.RoleTaxist
// 		}
// 	}
// 	if err != nil && err != sql.ErrNoRows {
// 		utils.RespondError(w, http.StatusInternalServerError, "Database error")
// 		return
// 	}

// 	// If not registered, complete registration
// 	if !isRegistered {
// 		if userType == "passenger" {
// 			if fullName == "" {
// 				utils.RespondError(w, http.StatusBadRequest, "Missing registration data")
// 				return
// 			}
// 			result, err := app.DB.Exec(
// 				"INSERT INTO passengers (full_name, phone, created_at) VALUES (?, ?, ?)",
// 				fullName, req.Phone, time.Now())
// 			if err != nil {
// 				utils.RespondError(w, http.StatusInternalServerError, "Failed to register passenger")
// 				return
// 			}
// 			userID, _ = result.LastInsertId()
// 			registeredType = models.RolePassenger
// 		} else if userType == "taxist" {
// 			if fullName == "" || carMake == "" || carModel == "" || carNumber == "" || carYear == 0 || rating < 0 {
// 				utils.RespondError(w, http.StatusBadRequest, "Missing registration data")
// 				return
// 			}
// 			result, err := app.DB.Exec(
// 				"INSERT INTO taxists (full_name, phone, car_make, car_model, car_year, car_number, rating, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
// 				fullName, req.Phone, carMake, carModel, carYear, carNumber, rating, time.Now())
// 			if err != nil {
// 				utils.RespondError(w, http.StatusInternalServerError, "Failed to register taxist")
// 				return
// 			}
// 			userID, _ = result.LastInsertId()
// 			registeredType = models.RoleTaxist
// 		} else {
// 			utils.RespondError(w, http.StatusBadRequest, "Invalid user type")
// 			return
// 		}
// 	}

// 	// Delete used code
// 	_, err = app.DB.Exec("DELETE FROM verification_codes WHERE phone = ?", req.Phone)
// 	if err != nil {
// 		log.Printf("Failed to delete verification code: %v", err)
// 	}

// 	// Generate JWT
// 	token, err := utils.GenerateJWT(app.Config, int(userID), registeredType)
// 	if err != nil {
// 		utils.RespondError(w, http.StatusInternalServerError, "Failed to generate token")
// 		return
// 	}

// 	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
// 		"userID":userID,
// 		"token": token,
// 	})
// }

// // ProtectedResponse defines the response for the protected endpoint
// type ProtectedResponse struct {
// 	UserID   int            `json:"user_id" example:"1" description:"ID of the authenticated user"`
// 	UserType models.UserType `json:"user_type" example:"passenger" description:"Type of user (passenger or taxist)"`
// }

// // handleProtected is an example protected route
// // @Summary Access protected resource
// // @Description Returns user information for an authenticated user. Requires a valid JWT token in the Authorization header.
// // @Tags Authentication
// // @Produce json
// // @Security BearerAuth
// // @Router /protected [get]
// func (app *App) handleProtected(w http.ResponseWriter, r *http.Request) {
// 	claims, ok := r.Context().Value("claims").(*models.Claims)
// 	if !ok {
// 		utils.RespondError(w, http.StatusUnauthorized, "Invalid claims")
// 		return
// 	}

// 	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
// 		"user_id":   claims.UserID,
// 		"user_type": claims.UserType,
// 	})
// }


// func (app *App) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		tokenStr := r.Header.Get("Authorization")
// 		if tokenStr == "" {
// 			utils.RespondError(w, http.StatusUnauthorized, "Authorization header missing")
// 			return
// 		}
// 		// Remove "Bearer " prefix
// 		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

// 		claims := &models.Claims{}
// 		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
// 			return []byte(app.Config.JWTSecret), nil
// 		})
// 		if err != nil || !token.Valid {
// 			utils.RespondError(w, http.StatusUnauthorized, "Invalid token")
// 			return
// 		}

// 		// Add claims to context
// 		ctx := context.WithValue(r.Context(), "claims", claims)
// 		next(w, r.WithContext(ctx))
// 	}
// }