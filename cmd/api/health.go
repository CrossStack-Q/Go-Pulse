package main

import "net/http"

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Write([]byte("Health is good \n"))
	}

	w.Write([]byte("Bye  Bye "))
}
