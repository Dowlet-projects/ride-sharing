package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"ride-sharing/models"
	"ride-sharing/utils"
	"strconv"
	"strings"
	"sync"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
	"googlemaps.github.io/maps"
)

// TaxistConn represents a taxi driver
type TaxistConn struct {
	ID        int             `json:"id"`
	Phone     string          `json:"phone"`
	FullName  string          `json:"full_name"`
	CarMake   string          `json:"car_make"`
	CarModel  string          `json:"car_model"`
	CarYear   int             `json:"car_year"`
	CarNumber string          `json:"car_number"`
	Rating    float64         `json:"rating"`
	Latitude  float64         `json:"latitude"`
	Longitude float64         `json:"longitude"`
	Conn      *websocket.Conn // WebSocket connection
}

// PassengerConn represents a client
type PassengerConn struct {
	ID              int             `json:"id"`
	Phone           string          `json:"phone"`
	FullName        string          `json:"full_name"`
	Latitude        float64         `json:"latitude"`
	Longitude       float64         `json:"longitude"`
	Money           float64         `json:"money"`
	DestinationName string          `json:"destination_name"`
	Conn            *websocket.Conn // WebSocket connection
}

// Message types for WebSocket communication
type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// LocationUpdate from taxist
type LocationUpdate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// RideRequest for passenger ride requests
type RideRequest struct {
	PassengerID   int     `json:"passenger_id"`
	PassengerName string  `json:"passenger_name"`
	FromLat       float64 `json:"from_lat"`
	FromLon       float64 `json:"from_lon"`
	ToLat         float64 `json:"to_lat"`
	ToLon         float64 `json:"to_lon"`
	Destination   string  `json:"destination"` // e.g., "Central Market"
}

// RideAcceptance for taxist acceptance
type RideAcceptance struct {
	TaxistID   int         `json:"taxist_id"`
	TaxistConn *TaxistConn `json:"taxist_conn"`
}

// Manager handles taxists, passengers, and ride requests
type Manager struct {
	Taxists      map[int]*TaxistConn
	Passengers   map[int]*PassengerConn
	RideRequests map[int]*RideRequest // Map passenger ID to active ride request
	Mutex        sync.RWMutex
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		allowedOrigins := []string{"http://localhost:3000", "http://192.168.55.42:8000", "null"}
		origin := r.Header.Get("Origin")
		for _, allowed := range allowedOrigins {
			if origin == allowed || (origin == "" && allowed == "null") {
				return true
			}
		}
		log.Printf("Rejected WebSocket connection from origin: %s", origin)
		return false
	},
}

var manager = Manager{
	Taxists:      make(map[int]*TaxistConn),
	Passengers:   make(map[int]*PassengerConn),
	RideRequests: make(map[int]*RideRequest),
}

// Initialize Google Maps client
var mapsClient *maps.Client

func init() {
	var err error
	mapsClient, err = maps.NewClient(maps.WithAPIKey("YOUR_GOOGLE_API_KEY"))
	if err != nil {
		log.Fatalf("Failed to initialize Google Maps client: %v", err)
	}
}

// SearchLocationByName queries Google Maps Places API for location coordinates
func SearchLocationByName(locationName string) (float64, float64, string, error) {
	if locationName == "" {
		return 0, 0, "", fmt.Errorf("location name cannot be empty")
	}
	req := &maps.TextSearchRequest{
		Query: locationName,
	}
	resp, err := mapsClient.TextSearch(context.Background(), req)
	if err != nil {
		return 0, 0, "", fmt.Errorf("failed to search location: %v", err)
	}
	if len(resp.Results) == 0 {
		return 0, 0, "", fmt.Errorf("no results found for %s", locationName)
	}
	result := resp.Results[0]
	lat := result.Geometry.Location.Lat
	lng := result.Geometry.Location.Lng
	formattedAddress := result.FormattedAddress
	return lat, lng, formattedAddress, nil
}

// HandleTaxistConnection handles WebSocket connections for taxists
func (app *App) HandleTaxistConnection(w http.ResponseWriter, r *http.Request) {
	// Extract token from query parameter or Authorization header
	tokenStr := r.URL.Query().Get("token")
	if tokenStr == "" {
		tokenStr = r.Header.Get("Authorization")
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
	}
	if tokenStr == "" {
		log.Println("Authorization missing")
		utils.RespondError(w, http.StatusUnauthorized, "Authorization missing")
		return
	}

	// Parse and validate JWT
	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(app.Config.JWTSecret), nil
	})
	if err != nil || !token.Valid {
		log.Printf("Invalid token: %v", err)
		utils.RespondError(w, http.StatusUnauthorized, "Invalid token")
		return
	}
	id := claims.UserID

	taxist, err := app.getTaxistFromDB(id)
	if err != nil {
		log.Printf("Taxist not found: %v", err)
		utils.RespondError(w, http.StatusUnauthorized, "Taxist not found")
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade error: %v, Headers: %v", err, r.Header)
		utils.RespondError(w, http.StatusBadRequest, fmt.Sprintf("Failed to upgrade to WebSocket: %v", err))
		return
	}

	taxist.Conn = conn
	manager.Mutex.Lock()
	manager.Taxists[taxist.ID] = taxist
	manager.Mutex.Unlock()

	defer func() {
		manager.Mutex.Lock()
		delete(manager.Taxists, taxist.ID)
		manager.Mutex.Unlock()
		conn.Close()
	}()

	app.sendNearbyPassengers(taxist)

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		if msg.Type == "location" {
			if data, ok := msg.Data.(map[string]interface{}); ok {
				lat, ok := data["latitude"].(float64)
				if !ok {
					log.Println("Invalid latitude type")
					continue
				}
				lon, ok := data["longitude"].(float64)
				if !ok {
					log.Println("Invalid longitude type")
					continue
				}
				if lat < -90 || lat > 90 || lon < -180 || lon > 180 {
					log.Printf("Invalid coordinates: lat=%f, lon=%f", lat, lon)
					continue
				}
				manager.Mutex.Lock()
				taxist.Latitude = lat
				taxist.Longitude = lon
				manager.Mutex.Unlock()
				app.broadcastTaxistsToPassengers()
			}
		}
	}
}

// HandlePassengerConnection handles WebSocket connections for passengers
func (app *App) HandlePassengerConnection(w http.ResponseWriter, r *http.Request) {
	// Extract token from query parameter or Authorization header
	tokenStr := r.URL.Query().Get("token")
	if tokenStr == "" {
		tokenStr = r.Header.Get("Authorization")
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
	}
	if tokenStr == "" {
		log.Println("Authorization missing")
		utils.RespondError(w, http.StatusUnauthorized, "Authorization missing")
		return
	}

	// Parse and validate JWT
	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(app.Config.JWTSecret), nil
	})
	if err != nil || !token.Valid {
		log.Printf("Invalid token: %v", err)
		utils.RespondError(w, http.StatusUnauthorized, "Invalid token")
		return
	}
	id := claims.UserID

	passenger, err := app.getPassengerFromDB(id)
	if err != nil {
		log.Printf("Passenger not found: %v", err)
		utils.RespondError(w, http.StatusUnauthorized, "Passenger not found")
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade error: %v, Headers: %v", err, r.Header)
		utils.RespondError(w, http.StatusBadRequest, fmt.Sprintf("Failed to upgrade to WebSocket: %v", err))
		return
	}
	// ws://192.168.55.42:8000/protected/ws/passenger?lat=kjhh&long=jkjhgl

	lat, long, err := parseLocation(r.URL.Query().Get("lat"), r.URL.Query().Get("long"))
	moneyStr := r.URL.Query().Get("money")

	if err != nil {
		log.Printf("Parse location error: %v", err)
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		conn.Close()
		return
	}

	money, err := strconv.ParseFloat(moneyStr, 64)

	if err != nil {
		fmt.Printf("Error converting to float64: %v\n", err)
		return
	}

	destination_name := r.URL.Query().Get("destination_name")

	passenger.Latitude = lat
	passenger.Longitude = long
	passenger.Money = money
	passenger.DestinationName = destination_name
	passenger.Conn = conn

	manager.Mutex.Lock()
	manager.Passengers[passenger.ID] = passenger
	manager.Mutex.Unlock()

	defer func() {
		manager.Mutex.Lock()
		delete(manager.Passengers, passenger.ID)
		manager.Mutex.Unlock()
		conn.Close()
	}()

	app.sendNearbyTaxists(passenger)

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		if msg.Type == "location" {
			if data, ok := msg.Data.(map[string]interface{}); ok {
				lat, ok := data["latitude"].(float64)
				if !ok {
					log.Println("Invalid latitude type")
					continue
				}
				lon, ok := data["longitude"].(float64)
				if !ok {
					log.Println("Invalid longitude type")
					continue
				}
				if lat < -90 || lat > 90 || lon < -180 || lon > 180 {
					log.Printf("Invalid coordinates: lat=%f, lon=%f", lat, lon)
					continue
				}
				manager.Mutex.Lock()
				passenger.Latitude = lat
				passenger.Longitude = lon
				manager.Mutex.Unlock()
				app.broadcastPassengersToTaxists()
			}
		}
	}
}

// HandlePassengerRideRequest handles a passenger's ride request with a destination name
func (app *App) HandlePassengerRideRequest(w http.ResponseWriter, r *http.Request) {
	// Extract token from Authorization header
	tokenStr := r.Header.Get("Authorization")
	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
	if tokenStr == "" {
		log.Println("Authorization missing")
		utils.RespondError(w, http.StatusUnauthorized, "Authorization missing")
		return
	}

	// Parse and validate JWT
	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(app.Config.JWTSecret), nil
	})
	if err != nil || !token.Valid {
		log.Printf("Invalid token: %v", err)
		utils.RespondError(w, http.StatusUnauthorized, "Invalid token")
		return
	}
	passengerID := claims.UserID

	var req struct {
		Destination string `json:"destination"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Invalid request body: %v", err)
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Search for the destination using Google Maps Places API
	toLat, toLon, formattedAddress, err := SearchLocationByName(req.Destination)
	if err != nil {
		log.Printf("Failed to find location: %v", err)
		utils.RespondError(w, http.StatusBadRequest, fmt.Sprintf("Failed to find location: %v", err))
		return
	}

	manager.Mutex.RLock()
	passenger, exists := manager.Passengers[passengerID]
	manager.Mutex.RUnlock()
	if !exists {
		log.Printf("Passenger %d not connected", passengerID)
		utils.RespondError(w, http.StatusNotFound, "Passenger not connected")
		return
	}

	// Create a ride request
	rideRequest := &RideRequest{
		PassengerID:   passengerID,
		PassengerName: passenger.FullName,
		FromLat:       passenger.Latitude,
		FromLon:       passenger.Longitude,
		ToLat:         toLat,
		ToLon:         toLon,
		Destination:   formattedAddress,
	}

	manager.Mutex.Lock()
	manager.RideRequests[passengerID] = rideRequest
	manager.Mutex.Unlock()

	// Notify nearby taxists
	app.notifyNearbyTaxists(rideRequest)

	// Respond to passenger
	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Ride request sent to nearby taxists"})
}

// Notify nearby taxists about a passenger's ride request
func (app *App) notifyNearbyTaxists(rideRequest *RideRequest) {
	const maxDistance = 10.0 // kilometers
	var nearbyTaxists []*TaxistConn

	manager.Mutex.RLock()
	for _, taxist := range manager.Taxists {
		if distance(rideRequest.FromLat, rideRequest.FromLon, taxist.Latitude, taxist.Longitude) <= maxDistance {
			nearbyTaxists = append(nearbyTaxists, taxist)
		}
	}
	manager.Mutex.RUnlock()

	// Send notification to nearby taxists
	msg := Message{
		Type: "ride_request",
		Data: rideRequest,
	}
	for _, taxist := range nearbyTaxists {
		if err := taxist.Conn.WriteJSON(msg); err != nil {
			log.Printf("Failed to notify taxist %d: %v", taxist.ID, err)
		}
	}
}

// HandleTaxistAccept handles a taxist's acceptance of a ride request
func (app *App) HandleTaxistAccept(w http.ResponseWriter, r *http.Request) {
	// Extract token from Authorization header
	tokenStr := r.Header.Get("Authorization")
	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
	if tokenStr == "" {
		log.Println("Authorization missing")
		utils.RespondError(w, http.StatusUnauthorized, "Authorization missing")
		return
	}

	// Parse and validate JWT
	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(app.Config.JWTSecret), nil
	})
	if err != nil || !token.Valid {
		log.Printf("Invalid token: %v", err)
		utils.RespondError(w, http.StatusUnauthorized, "Invalid token")
		return
	}
	taxistID := claims.UserID

	var req struct {
		PassengerID int `json:"passenger_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Invalid request body: %v", err)
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	manager.Mutex.RLock()
	taxist, taxistExists := manager.Taxists[taxistID]
	passenger, passengerExists := manager.Passengers[req.PassengerID]
	_, requestExists := manager.RideRequests[req.PassengerID]
	manager.Mutex.RUnlock()

	if !taxistExists || !passengerExists || !requestExists {
		log.Printf("Taxist %d, passenger %d, or ride request not found", taxistID, req.PassengerID)
		utils.RespondError(w, http.StatusNotFound, "Taxist, passenger, or ride request not found")
		return
	}

	// Notify passenger of acceptance
	acceptance := RideAcceptance{
		TaxistID:   taxistID,
		TaxistConn: taxist,
	}
	msg := Message{
		Type: "ride_accepted",
		Data: acceptance,
	}
	if err := passenger.Conn.WriteJSON(msg); err != nil {
		log.Printf("Failed to notify passenger %d: %v", req.PassengerID, err)
		utils.RespondError(w, http.StatusInternalServerError, "Failed to notify passenger")
		return
	}

	// Clear the ride request
	manager.Mutex.Lock()
	delete(manager.RideRequests, req.PassengerID)
	manager.Mutex.Unlock()

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Ride request accepted"})
}

// getTaxistFromDB retrieves taxist details from the database
func (app *App) getTaxistFromDB(id int) (*TaxistConn, error) {
	var taxist TaxistConn
	err := app.DB.QueryRow(`
		SELECT id, phone, full_name, car_make, car_model, car_year, car_number, rating
		FROM taxists WHERE id = ?`, id).Scan(&taxist.ID, &taxist.Phone, &taxist.FullName, &taxist.CarMake,
		&taxist.CarModel, &taxist.CarYear, &taxist.CarNumber, &taxist.Rating)
	if err != nil {
		return nil, err
	}
	return &taxist, nil
}

// getPassengerFromDB retrieves passenger details from the database
func (app *App) getPassengerFromDB(id int) (*PassengerConn, error) {
	var passenger PassengerConn
	err := app.DB.QueryRow(`SELECT id, phone, full_name FROM passengers WHERE id = ?`, id).
		Scan(&passenger.ID, &passenger.Phone, &passenger.FullName)
	if err != nil {
		return nil, err
	}
	return &passenger, nil
}

// broadcastTaxistsToPassengers sends nearby taxists to all passengers
func (app *App) broadcastTaxistsToPassengers() {
	manager.Mutex.RLock()
	defer manager.Mutex.RUnlock()

	for _, passenger := range manager.Passengers {
		app.sendNearbyTaxists(passenger)
	}
}

// broadcastTaxistsToPassengers sends nearby taxists to all passengers
func (app *App) broadcastPassengersToTaxists() {
	manager.Mutex.RLock()
	defer manager.Mutex.RUnlock()

	for _, taxist := range manager.Taxists {
		app.sendNearbyPassengers(taxist)
	}
}

// sendNearbyTaxists sends nearby taxists to a specific passenger
func (app *App) sendNearbyTaxists(passenger *PassengerConn) {
	var nearbyTaxists []TaxistConn
	const maxDistance = 10.0

	manager.Mutex.RLock()
	defer manager.Mutex.RUnlock()

	for _, taxist := range manager.Taxists {
		if distance(passenger.Latitude, passenger.Longitude, taxist.Latitude, taxist.Longitude) <= maxDistance {
			nearbyTaxists = append(nearbyTaxists, *taxist)
		}
	}

	msg := Message{
		Type: "taxists",
		Data: nearbyTaxists,
	}

	if err := passenger.Conn.WriteJSON(msg); err != nil {
		log.Println("Write error:", err)
	}
}


// sendNearbyPassengers sends nearby passengers to a specific taxist
func (app *App) sendNearbyPassengers(taxist *TaxistConn) {
	var nearbyPassengers []PassengerConn
	const maxDistance = 10.0

	manager.Mutex.RLock()
	defer manager.Mutex.RUnlock()

	for _, passenger := range manager.Passengers {
		if distance(taxist.Latitude, taxist.Longitude, passenger.Latitude, passenger.Longitude) <= maxDistance {
			nearbyPassengers = append(nearbyPassengers, *passenger)
		}
	}

	msg := Message{
		Type: "passengers",
		Data: nearbyPassengers,
	}

	if err := taxist.Conn.WriteJSON(msg); err != nil {
		log.Println("Write error:", err)
	}
}

// distance calculates the distance between two points using the Haversine formula
func distance(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371 // Earth radius in kilometers
	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}

// parseLocation parses latitude and longitude from strings
func parseLocation(latStr, longStr string) (float64, float64, error) {
	var lat, long float64
	if _, err := fmt.Sscanf(latStr, "%f", &lat); err != nil {
		return 0, 0, fmt.Errorf("invalid latitude: %v", err)
	}
	if _, err := fmt.Sscanf(longStr, "%f", &long); err != nil {
		return 0, 0, fmt.Errorf("invalid longitude: %v", err)
	}
	if lat < -90 || lat > 90 || long < -180 || long > 180 {
		return 0, 0, fmt.Errorf("coordinates out of range: lat=%f, long=%f", lat, long)
	}
	return lat, long, nil
}
