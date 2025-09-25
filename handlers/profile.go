package handlers

import (
	"SSE/database"
	"SSE/models"
	"SSE/sessions"
	"bytes"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"html/template"
	"log"
	"net/http"
)

func Profile(w http.ResponseWriter, r *http.Request) {
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

	var user models.User
	err = collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	tmpl, err := template.ParseFiles("templates/profile.html")
	if err != nil {
		log.Printf("template parse error: %v", err)
		http.Error(w, "Failed to load template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	iinSet := len(user.IIN) > 0 && user.IIN != ""
	identityCardSet := len(user.IdentityCard) > 0 && user.IdentityCard != ""
	addressSet := len(user.Address) > 0 && user.Address != ""

	data := struct {
		FirstName       string
		LastName        string
		Birthday        string
		Address         string
		AddressSet      bool
		ID              string
		Admin           bool
		IIN             string
		IdentityCard    string
		IINSet          bool
		IdentityCardSet bool
		CreatedAt       string
		UpdatedAt       string
	}{
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		Birthday:        user.Birthday.Format("January 2, 2006"),
		Address:         user.Address,
		AddressSet:      addressSet,
		ID:              userID,
		Admin:           user.Admin,
		IIN:             user.IIN,
		IdentityCard:    user.IdentityCard,
		IINSet:          iinSet,
		IdentityCardSet: identityCardSet,
		CreatedAt:       user.CreatedAt.Format("January 2, 2006 15:04:05"),
		UpdatedAt:       user.UpdatedAt.Format("January 2, 2006 15:04:05"),
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		log.Printf("template execute error: %v", err)
		http.Error(w, "Failed to render page: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = buf.WriteTo(w)
}
