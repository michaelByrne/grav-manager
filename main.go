package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func main() {
	svc, err := NewAccountService()
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()

	r.HandleFunc("/upgrade/{account_id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["account_id"]

		err := svc.UpgradePlan(100, id)
		if err != nil {
			log.Fatal(err)
		}

		w.WriteHeader(http.StatusOK)
	})

	r.HandleFunc("/count/{account_id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["account_id"]

		count, err := svc.GetUserCount(id)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}

		respondWithJSON(w, http.StatusOK, map[string]int{"user_count": count})
	})

	r.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		var user User
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&user); err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}
		defer r.Body.Close()

		if user.ID == "" || user.AccountID == "" {
			respondWithError(w, http.StatusBadRequest, "Recieved empty account id or user id")
			return
		}

		err = svc.RegisterUser(user)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	})

	amw := authMiddleware{}
	amw.SetToken("shmoken")

	r.Use(amw.AuthMiddleware)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3001"},
		AllowedHeaders: []string{"Authorization"},
		AllowCredentials: true,
		//Debug: true,
		})
		
	handler := c.Handler(r)	

	fmt.Println("running server on 8443")
	err = http.ListenAndServeTLS(":8443", "ca-cert.pem", "ca-key.pem", handler)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
