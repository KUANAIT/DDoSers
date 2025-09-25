package routes

import (
	"SSE/handlers"
	"net/http"
)

func RegisterRoutes() {
	http.HandleFunc("/users", handlers.CreateUser)
	http.HandleFunc("/users/get", handlers.GetUser)
	http.HandleFunc("/users/update", handlers.UpdateUser)
	http.HandleFunc("/users/delete", handlers.DeleteUser)
	http.HandleFunc("/loginuser", handlers.LoginCustomer)
	http.HandleFunc("/profile", handlers.Profile)

}
