package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// HealthCheck used to return health check responses
type HealthCheck struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message,omitempty"`
}

func sendAsJSON(rw http.ResponseWriter, obj interface{}) {

	rw.Header().Set("Content-Type", "application/json")

	var encoder = json.NewEncoder(rw)

	encoder.Encode(obj)
}

func sendFailure(rw http.ResponseWriter, code int, msg string) {

	rw.WriteHeader(code)
	rw.Header().Set("Content-Type", "application/json")
	rw.Write([]byte(msg))
}

func certHandler(w http.ResponseWriter, r *http.Request) {

	var domain = r.URL.Path[len("/cert/"):]

	var cert, err = container.GetCert(domain)
	if err != nil {
		log.Printf("Error fetching cert, review and take action: %s", err)
		sendFailure(w, http.StatusInternalServerError, "Failed to obtain cert, try again later")
		return
	}

	sendAsJSON(w, cert)
}

func renewHandler(w http.ResponseWriter, r *http.Request) {

	var domain = r.URL.Path[len("/renew/"):]

	_, ok := container.certs[domain]
	if !ok {
		sendFailure(w, http.StatusBadRequest, "no domain to renew")
		return
	}

	var cert, err = container.RenewCert(domain)
	if err != nil {
		log.Printf("Error renewing cert, review and take action: %s", err)
		sendFailure(w, http.StatusInternalServerError, "Failed to obtain cert, try again later")
		return
	}

	sendAsJSON(w, cert)
}

func certsHandler(w http.ResponseWriter, r *http.Request) {

	var certs = []Cert{}

	for _, c := range container.certs {
		certs = append(certs, c)
	}

	sendAsJSON(w, certs)
}

func healthcheck(w http.ResponseWriter, r *http.Request) {
	sendAsJSON(w, HealthCheck{Ok: true})
}
