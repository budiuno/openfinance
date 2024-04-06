package main

import (
	"fmt"
	"net/http"

	accountController "github.com/budiuno/openfinance/controller"
	auth "github.com/budiuno/openfinance/middlewares/auth"
	"github.com/budiuno/openfinance/repo"
	"github.com/gorilla/mux"
)

func main() {

	user := auth.User{UserID: "e7cf3bed-32c2-4d9a-8ea4-3ab8a2e58a93", Username: "Budi Setyawan", Password: "secret_password"}

	// Generate a token for the sample user
	token, err := auth.GenerateToken(user)
	if err != nil {
		fmt.Println("Error generating token:", err)
		return
	}
	fmt.Println("token: ", token)

	r := mux.NewRouter()

	r.HandleFunc("/healthcheck", healthCheckHandler).Methods("GET")

	validateHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accountController.ValidateAccount(w, r, repo.AccountRepo{})
	})

	r.Handle("/v1/account/{bank_code}/{account_number}", auth.Authenticate(validateHandler)).Methods("GET")

	port := 8000
	fmt.Printf("Server is running on http://localhost:%d\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	_, err := fmt.Fprintf(w, `{"status": "ok"}`)
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}
