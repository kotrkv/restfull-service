package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type User struct {
	Username string `json:"user"`
	Password string `json:"password"`
}

var user User
var PORT = ":1234"
var DATA = make(map[string]string)

func main() {
	fmt.Println("Hello Restfull server....")
	arguments := os.Args
	if len(arguments) != 1 {
		PORT = ":" + arguments[1]
	}
	mux := http.NewServeMux()
	s := &http.Server{
		Addr:         PORT,
		Handler:      mux,
		IdleTimeout:  10 * time.Second,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
	}
	mux.Handle("/", http.HandlerFunc(defaultHandler))
	mux.Handle("/time", http.HandlerFunc(timeHandler))
	mux.Handle("/users", http.HandlerFunc(addHandler))

	fmt.Println("Ready to serve at", PORT)
	err := s.ListenAndServe()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving:", r.URL.Path, "from", r.Host)
	w.WriteHeader(http.StatusNotFound)
	Body := "Thanks for visiting!\n"
	fmt.Fprintf(w, "%s", Body)
}

func timeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving: ", r.URL.Path, "from", r.Host)
	t := time.Now().Format(time.RFC1123)
	Body := "The current time is: " + t + "\n"
	fmt.Fprintf(w, "%s", Body)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving: ", r.URL.Path, "from", r.Host, r.Method)
	if r.Method != http.MethodPost {
		http.Error(w, "Error: ", http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "%s\n", "Method not allowed!!!")
		return
	}
	d, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error:", http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(d, &user)
	if err != nil {
		log.Println("Error on unmarshalling:", err)
		http.Error(w, "Error:", http.StatusBadRequest)
		return
	}
	if user.Username != "" {
		DATA[user.Username] = user.Password
		log.Println(DATA)
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Error:", http.StatusBadRequest)
		return
	}
}
