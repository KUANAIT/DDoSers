package handlers

import (
	"SSE/database"
	"SSE/models"
	"SSE/sessions"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

func LoginCustomer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var credentials struct {
		IIN      string `json:"iin"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	collection, err := database.GetCollection("SSE", "users")
	if err != nil {
		http.Error(w, "Failed to get database collection", http.StatusInternalServerError)
		return
	}

	var users []models.User
	cursor, err := collection.Find(r.Context(), bson.M{})
	if err != nil {
		http.Error(w, "Invalid IIN or password", http.StatusUnauthorized)
		return
	}
	defer cursor.Close(r.Context())

	if err = cursor.All(r.Context(), &users); err != nil {
		http.Error(w, "Invalid IIN or password", http.StatusUnauthorized)
		return
	}

	var user models.User
	found := false
	for _, u := range users {
		if u.CheckIIN(credentials.IIN) && u.CheckPassword(credentials.Password) {
			user = u
			found = true
			break
		}
	}

	if !found {
		http.Error(w, "Invalid IIN or password", http.StatusUnauthorized)
		return
	}

	userID := user.ID.Hex()

	if err := sessions.SetUserSession(w, r, userID); err != nil {
		http.Error(w, "Failed to set session", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Login successful"})
}
