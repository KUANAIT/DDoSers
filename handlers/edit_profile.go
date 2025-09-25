package handlers

import (
	"SSE/auth"
	"SSE/database"
	"SSE/models"
	"SSE/sessions"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"html/template"
	"log"
	"net/http"
	"time"
)

func EditProfile(w http.ResponseWriter, r *http.Request) {
	session, err := sessions.Get(r)
	if err != nil {
		http.Error(w, "Failed to retrieve session", http.StatusInternalServerError)
		return
	}

	userID, ok := session.Values["user_id"].(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusInternalServerError)
		return
	}

	collection, err := database.GetCollection("SSE", "users")
	if err != nil {
		http.Error(w, "Failed to get database collection", http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case http.MethodGet:
		var user models.User
		if err := collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&user); err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		tmpl, err := template.ParseFiles("templates/edit_profile.html")
		if err != nil {
			log.Printf("template parse error: %v", err)
			http.Error(w, "Failed to load template: "+err.Error(), http.StatusInternalServerError)
			return
		}

		data := struct {
			FirstName string
			LastName  string
			Birthday  string
			Address   string
			ID        string
		}{
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Birthday:  user.Birthday.Format("2006-01-02"),
			Address:   user.Address,
			ID:        userID,
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.Execute(w, data); err != nil {
			log.Printf("template execute error: %v", err)
			http.Error(w, "Failed to render page: "+err.Error(), http.StatusInternalServerError)
			return
		}

	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}

		firstName := r.FormValue("first_name")
		lastName := r.FormValue("last_name")
		birthdayStr := r.FormValue("birthday")
		address := r.FormValue("address")
		newPassword := r.FormValue("password")

		updateFields := bson.M{}

		if firstName != "" {
			updateFields["first_name"] = firstName
		}
		if lastName != "" {
			updateFields["last_name"] = lastName
		}
		if birthdayStr != "" {
			birthday, err := time.Parse("2006-01-02", birthdayStr)
			if err != nil {
				http.Error(w, "Invalid birthday format. Use YYYY-MM-DD", http.StatusBadRequest)
				return
			}
			updateFields["birthday"] = birthday
		}
		if address != "" {
			tmp := models.User{Address: address}
			if err := tmp.HashSensitiveData(); err != nil {
				http.Error(w, "Failed to hash address", http.StatusInternalServerError)
				return
			}
			updateFields["address"] = tmp.Address
		}
		if newPassword != "" {
			hashed, err := auth.HashPassword(newPassword)
			if err != nil {
				http.Error(w, "Failed to hash password", http.StatusInternalServerError)
				return
			}
			updateFields["password"] = hashed
		}

		if len(updateFields) == 0 {
			http.Redirect(w, r, "/profile", http.StatusSeeOther)
			return
		}

		updateFields["updated_at"] = time.Now()

		result, err := collection.UpdateOne(r.Context(), bson.M{"_id": objID}, bson.M{"$set": updateFields})
		if err != nil {
			http.Error(w, "Failed to update profile", http.StatusInternalServerError)
			return
		}
		if result.MatchedCount == 0 {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		http.Redirect(w, r, "/profile", http.StatusSeeOther)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
