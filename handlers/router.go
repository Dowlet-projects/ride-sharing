package handlers

import (
	"database/sql"
	//"net/http"

	"ride-sharing/config"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// App holds application dependencies
type App struct {
	DB     *sql.DB
	Config *config.Config
}

// SetupRouter configures the HTTP router
// @Summary Setup API router
// @Description Configures the HTTP router with CORS and routes
// @Tags Internal
// @Produce json
func SetupRouter(db *sql.DB, cfg *config.Config) *mux.Router {
	app := &App{DB: db, Config: cfg}
	router := mux.NewRouter()

	// Configure CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
		AllowedHeaders:   []string{
			"Content-Type",
			"Authorization",
			"X-Requested-With",
			"Accept",
			"Origin",
		},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           86400,
	})
	router.Use(c.Handler)

	// Public routes
	router.HandleFunc("/passenger/register", app.handlePassengerRegister).Methods("POST", "OPTIONS")
	router.HandleFunc("/passenger/login", app.handlePassengerLogin).Methods("POST", "OPTIONS")
	router.HandleFunc("/verify", app.handleVerifyCode).Methods("POST", "OPTIONS")
	router.HandleFunc("/taxist/register", app.handleTaxistRegister).Methods("POST", "OPTIONS")
	router.HandleFunc("/taxist/login", app.handleTaxistLogin).Methods("POST", "OPTIONS")
	router.HandleFunc("/taxist/verify", app.handleVerifyCode).Methods("POST", "OPTIONS")
	router.HandleFunc("/makes", app.CreateMake).Methods("POST", "OPTIONS")
	router.HandleFunc("/makes", app.GetAllMakes).Methods("GET", "OPTIONS")
	router.HandleFunc("/models", app.CreateModel).Methods("POST", "OPTIONS")
	router.HandleFunc("/models/{make_id}", app.GetAllModelsById).Methods("GET", "OPTIONS")
	router.HandleFunc("/places", app.GetAllPlaces).Methods("GET", "OPTIONS")
	router.HandleFunc("/places", app.CreatePlace).Methods("POST", "OPTIONS")
	router.HandleFunc("/distances", app.createDistance).Methods("POST", "OPTIONS")
	router.HandleFunc("/distances", app.GetAllDistances).Methods("GET", "OPTIONS")

	// Protected routes
	protected := router.PathPrefix("/protected").Subrouter()
	protected.Use(app.authMiddleware)
	protected.HandleFunc("", app.handleProtected).Methods("GET", "OPTIONS")
	protected.HandleFunc("/announcements", app.CreateAnnouncement).Methods("POST", "OPTIONS")
	protected.HandleFunc("/ugurlar", app.GetAllUgurlar).Methods("GET", "OPTIONS")
	protected.HandleFunc("/ugurlar/{ugur_id}", app.GetUgurById).Methods("GET", "OPTIONS")
	//protected.HandleFunc("/taxists/{taxist_id}", app.TaxistProfile).Methods("GET", "OPTIONS")
	protected.HandleFunc("/reserve-passengers/{taxi_ann_id}", app.CreateReservePassengers).Methods("POST", "OPTIONS")
	protected.HandleFunc("/reserve-packages/{taxi_ann_id}", app.CreateReservePackages).Methods("POST", "OPTIONS")
	protected.HandleFunc("/taxist-rating/{taxist_id}", app.UpdateRatingTaxist).Methods("PUT", "OPTIONS")
	protected.HandleFunc("/taxist-comments/{taxist_id}", app.GetAllTaxistComments).Methods("GET", "OPTIONS")
	protected.HandleFunc("/taxist-comments/{taxist_id}", app.CreateComment).Methods("POST", "OPTIONS")
	protected.HandleFunc("/profile", app.Profile).Methods("GET", "OPTIONS")
	protected.HandleFunc("/taxist-notifications/{taxist_id}", app.GetAllTaxistNotifications).Methods("GET", "OPTIONS")
	protected.HandleFunc("/reverse-details/{reverse_id}", app.ReverseDetails).Methods("GET", "OPTIONS")
	protected.HandleFunc("/favourites", app.CreateFavourites).Methods("POST", "OPTIONS")
	protected.HandleFunc("/favourites", app.GetAllFavourites).Methods("GET", "OPTIONS")
	protected.HandleFunc("/taxist-departed/{taxi_ann_id}", app.UpdateTaxistAnnouncements).Methods("PUT", "OPTIONS")
	protected.HandleFunc("/taxist-announcements/{departed}", app.GetTaxistAnnouncements).Methods("GET", "OPTIONS")
	protected.HandleFunc("/distances/{id}", app.DeleteDistances).Methods("DELETE", "OPTIONS")
	protected.HandleFunc("/places/{id}", app.DeletePlace).Methods("DELETE", "OPTIONS")
	protected.HandleFunc("/makes/{id}", app.DeleteMake).Methods("DELETE", "OPTIONS")
	protected.HandleFunc("/models/{model_id}", app.DeleteModel).Methods("DELETE", "OPTIONS")
	return router
}

