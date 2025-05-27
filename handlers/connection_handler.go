package handlers

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"ride-sharing/models"
	"ride-sharing/utils"
	"sync"
	"github.com/gorilla/websocket"
)

// Taxist represents a taxi driver
type TaxistConn struct {
	ID int `json:"id"`
	FullName string `json:"full_name"`
	CarMake string `json:"car_make"`
	CarModel string `json:"car_model"`
	CarYear int `json:"car_year"`
	CarNumber string `json:"car_number"`
	Rating float64 `json:"rating"`
	Latitude float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Conn *websocket.Conn // WebSocket connection
}


// Passenger represents a client
type PassengerConn struct {
	ID int `json:"id"`
	FullName string `json:"full_name"`
	Latitude float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Conn *websocket.Conn // WebSocket connection
}

// Message types for WebSocket communication
type Message struct {
	Type string `json:"type"`
	Data interface{} `json:"data"`
}



// Location update from taxist
type LocationUpdate struct {
	Latitude float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}


// Manager handles taxists and passengers
type Manager struct {
	Taxists map[int]*TaxistConn
	Passengers map[int]*PassengerConn
	Mutex sync.RWMutex
}

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool { return true },  // Restrict in production
}


var manager = Manager{
	Taxists: make(map[int]*TaxistConn),
	Passengers: make(map[int]*PassengerConn),
}


// HandleTaxistConnection handles WebSocket connections for taxists
// @Summary Establish WebSocket connection for taxist
// @Description Upgrades an HTTP connection to a WebSocket connection for a taxist to send location updates. Requires Bearer token authentication. Expects WebSocket messages with type "location" and data containing latitude and longitude.
// @Tags WebSocket
// @Produce json
// @Security BearerAuth
// @Router /protected/ws/taxist [get]
func (app *App) HandleTaxistConnection(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("claims").(*models.Claims)
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "Invalid claims")
		return
	}
	id := claims.UserID

	taxist, err := app.getTaxistFromDB(id)
	if err != nil {
		utils.RespondError(w, http.StatusUnauthorized, "Taxist not found")
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error: ", err)
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

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		if msg.Type == "location" {
			//var loc LocationUpdate
			if data, ok := msg.Data.(map[string]interface{}); ok {
				lat, _ := data["latitude"].(float64)
				lon, _ := data["longitude"].(float64)
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
// @Summary Establish WebSocket connection for passenger
// @Description Upgrades an HTTP connection to a WebSocket connection for a passenger to receive nearby taxists. Requires Bearer token authentication. Expects query parameters lat and long for initial passenger location.
// @Tags WebSocket
// @Produce json
// @Security BearerAuth
// @Param lat query float64 true "Passenger's latitude"
// @Param long query float64 true "Passenger's longitude"
// @Router /protected/ws/passenger [get]
func (app *App) HandlePassengerConnection(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("claims").(*models.Claims)
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "Invalid claims")
		return
	}
	id := claims.UserID

	passenger, err := app.getPassengerFromDB(id)
	if err != nil {
		utils.RespondError(w, http.StatusUnauthorized, "Passenger not found")
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	lat, long := parseLocation(r.URL.Query().Get("lat"), r.URL.Query().Get("long"))
	passenger.Latitude = lat
	passenger.Longitude = long
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
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}
	}
}


func (app *App) getTaxistFromDB(id int) (*TaxistConn, error) {
	var taxist TaxistConn
	err := app.DB.QueryRow(`
		SELECT id, full_name, car_make, car_model, car_year, car_number, rating
		FROM taxists WHERE id = ?`, id).Scan(&taxist.ID, &taxist.FullName, &taxist.CarMake,
		&taxist.CarModel, &taxist.CarYear, &taxist.CarNumber, &taxist.Rating)
	
	if err != nil {
		return nil, err
	}
	return &taxist, nil
}


func (app *App) getPassengerFromDB(id int) (*PassengerConn, error) {
	var passenger PassengerConn
	err := app.DB.QueryRow(`SELECT id, full_name FROM passengers WHERE id = ?`, id).
	Scan(&passenger.ID, &passenger.FullName)
	
	if err != nil {
		return nil, err
	}
	return &passenger, nil
}

func (app *App) broadcastTaxistsToPassengers() {
	manager.Mutex.RLock()
	defer manager.Mutex.RUnlock()


	for _, passenger := range manager.Passengers {
		app.sendNearbyTaxists(passenger)
	}
}

func (app *App) sendNearbyTaxists(passenger *PassengerConn) {
	var nearbyTaxists []TaxistConn
	const maxDistance = 10.0

	for _, taxist := range manager.Taxists {
		if distance(passenger.Latitude, passenger.Longitude, taxist.Latitude,
		taxist.Longitude) <= maxDistance {
			nearbyTaxists = append(nearbyTaxists, *taxist)
		}
	}

	msg := Message{
		Type: "taxists",
		Data: nearbyTaxists,
	}

	if err := passenger.Conn.WriteJSON(msg); err != nil {
		log.Println("Write error:", err)
		return
	}
}

func distance(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371
	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}

func parseLocation(latStr, longStr string) (float64, float64) {
	var lat, long float64
	fmt.Sscanf(latStr, "%f", &lat)
	fmt.Sscanf(longStr, "%f", &long)
	return lat, long
}





















