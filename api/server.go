package api

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Item struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type Server struct {
	*mux.Router

	restaurantItems []Item
}

func NewServer() *Server {
	s := &Server{
		Router:          mux.NewRouter(),
		restaurantItems: []Item{},
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.HandleFunc("/restaurant-items", s.listRestaurantItems()).Methods("GET")
	s.HandleFunc("/restaurant-items", s.createRestaurantItem()).Methods("POST")
	s.HandleFunc("/restaurant-items/{id}", s.removeRestaurantItem()).Methods("DELETE")
}

func (s *Server) createRestaurantItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var i Item
		if err := json.NewDecoder(r.Body).Decode(&i); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		i.ID = uuid.New()
		s.restaurantItems = append(s.restaurantItems, i)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(i); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) listRestaurantItems() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(s.restaurantItems); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) removeRestaurantItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr, _ := mux.Vars(r)["id"]
		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		for i, item := range s.restaurantItems {
			if item.ID == id {
				s.restaurantItems = append(s.restaurantItems[:i], s.restaurantItems[i+1:]...)
				break
			}
		}
	}
}
