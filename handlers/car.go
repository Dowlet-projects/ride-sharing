package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"ride-sharing/models"
	"strconv"
	"time"
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

//CreateMake handles POST /makes
// @Summary POST a new car make
// @Description add a new car make
// @Tags Car details
// @Accept json
// @Produce json
// @Param body body handlers.MakesRequest true "Car details"
// @Router /makes [post]
func (h *App) CreateMake(w http.ResponseWriter, r *http.Request) {
	var make models.Make
	if err:=json.NewDecoder(r.Body).Decode(&make); err !=nil {
		http.Error(w, "Invalid Input", http.StatusBadRequest)
		return
	}
	
	query := "INSERT INTO car_makes (name) VALUES (?)"
	result, err := h.DB.Exec(query, make.Name)
	if err != nil {
		http.Error(w, "Failed to create a make", http.StatusInternalServerError)
		return
	}
	
	id, _ := result.LastInsertId()
	make.ID = int(id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(make)
}

//GetAllMakes handles GET /makes
// @Summary GET car makes
// @Description get all car makes
// @Tags Car details
// @Produce json
// @Router /makes [get]
func (h *App) GetAllMakes(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query("SELECT id, name FROM car_makes")

	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var makes []models.Make = []models.Make{}
	for rows.Next() {
		var make models.Make
		if err:=rows.Scan(&make.ID, &make.Name); err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		makes = append(makes, make)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(makes)
}



// DeleteMake handles DELETE /makes/{id}
// @Summary DELETE make
// @Description DELETE make by id
// @Tags Car details
// @Produce json
// @Security BearerAuth
// @Param id path string true "Make ID"
// @Router /protected/makes/{id} [DELETE]
func (h *App) DeleteMake(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	query := "DELETE FROM car_makes WHERE id = ? "
	result, err := h.DB.Exec(query, id)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to delete make", http.StatusInternalServerError)
		return
	}
	
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "make not found", http.StatusNotFound)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}

//CreateModel handles POST /models
// @Summary POST a new car model
// @Description add a new car model
// @Tags Car details
// @Accept json
// @Produce json
// @Param body body handlers.ModelsRequest true "Car details"
// @Router /models [post]
func (h *App) CreateModel(w http.ResponseWriter, r *http.Request) {
	var model models.Model
	if err:=json.NewDecoder(r.Body).Decode(&model); err !=nil {
		http.Error(w, "Invalid Input", http.StatusBadRequest)
		return
	}

	query := "INSERT INTO car_models (name, make_id) VALUES (?, ?)"
	result, err := h.DB.Exec(query, model.Name, model.MakeID)
	if err != nil {
		http.Error(w, "Failed to create a model", http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	model.ID = int(id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(model)
}

// rows, err := h.DB.Query("SELECT id, name, make_id FROM car_models WHERE make_id = ?", id)

// if err != nil {
	// 	http.Error(w, "Server error", http.StatusInternalServerError)
	// 	return
	// }
	
	// defer rows.Close()
	// var modelss []models.Model = []models.Model{} 
	// for rows.Next() {
		// 	var model models.Model
		// 	if err:=rows.Scan(&model.ID, &model.Name, &model.MakeID); err != nil {
// 		http.Error(w, "Server error", http.StatusInternalServerError)
// 		return
// 	}

// 	modelss = append(modelss, model)
// }


// Redis code.


// GetAllModelsById handles GET /models/{make_id}
// @Summary Get car models by make ID
// @Description Retrieve all car models for a given make ID
// @Tags Car details
// @Produce json
// @Param make_id path string true "Make ID"
// @Router /models/{make_id} [get]
func (h *App) GetAllModelsById(w http.ResponseWriter, r *http.Request) {
	
	vars := mux.Vars(r)
	// id, err := strconv.Atoi(vars["make_id"])
	id, err := strconv.ParseInt(vars["make_id"], 10, 64)

	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}



	modelss, err := h.getCachedModels(id)
	
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch models: %v", err), http.StatusInternalServerError)
		return
	}

	if modelss == nil {
		http.Error(w, "models not found", http.StatusNotFound)
		return
	}


	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(modelss); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}


const (
	redisAddr = "127.0.0.1:6379"
	ttl = 3600 * time.Second  // Cache ttl (e.g, 1 hour)
)

// getCachedModels retrieves an models from Redis by make_id, with fallback to getModels
func (h *App) getCachedModels(id int64) (*[]models.Model, error) {

	rdb := redis.NewClient(&redis.Options {
		Addr: redisAddr,
		Password: "",
		DB: 0,
	})

	defer rdb.Close()

	ctx := context.Background()
	key := fmt.Sprintf("models:make_%d", id)

	jsonStr, err := rdb.Get(ctx, key).Result()
	if err==redis.Nil {
		// Cache, miss, fetch from Mysql
		models, err := h.getModels(id)

		if err != nil {
			return nil, fmt.Errorf("failed to fetch model %w", err)
		}

		if models == nil {
			return nil, nil
		}

		jsonBytes, err := json.Marshal(models)

		if err != nil {
			return nil, fmt.Errorf("failed to marshal model to JSON %w", err)
		}

		err = rdb.SetEX(ctx, key, jsonBytes, ttl).Err()
		if err != nil {
			return nil, fmt.Errorf("failed to cache model in Redis: %w", err)
		}

		return models, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get model from Redis: %w", err)
	}

	var models []models.Model
	
	if err := json.Unmarshal([]byte(jsonStr), &models); err != nil {
		return nil, fmt.Errorf("failed to unmarshal article JSON: %w", err)
	}

	return &models, nil

}

// getModels fetches a models from database by make_id
func (h *App) getModels(id int64) (*[]models.Model, error) {
	rows, err := h.DB.Query("SELECT id, name, make_id FROM car_models WHERE make_id = ?", id)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to Mysql %w", err)
	}
	
	defer rows.Close()
	
	var modelss []models.Model = []models.Model{} 
	for rows.Next() {
		var model models.Model
		if err:=rows.Scan(&model.ID, &model.Name, &model.MakeID); err != nil {
			return nil, fmt.Errorf("server error: %w", err)
		}

		modelss = append(modelss, model)
	}

	return &modelss, nil
}





// DeleteModel handles DELETE /models/{model_id}
// @Summary DELETE model
// @Description DELETE model by id
// @Tags Car details
// @Produce json
// @Security BearerAuth
// @Param model_id path string true "Model ID"
// @Router /protected/models/{model_id} [DELETE]
func (h *App) DeleteModel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["model_id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	query := "DELETE FROM car_models WHERE id = ? "
	result, err := h.DB.Exec(query, id)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to delete model", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "model not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}


// Accept json


// Accept json