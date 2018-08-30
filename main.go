package main

import (
	"log"
	"net/http"
)

var container CertContainer

func main() {

	container = NewCertContainer(RequestCert)

	http.HandleFunc("/cert/", certHandler)
	http.HandleFunc("/renew/", renewHandler)
	http.HandleFunc("/certs", certsHandler)
	http.HandleFunc("/expired", expiredHandler)

	http.HandleFunc("/health", healthcheck)

	log.Fatal(http.ListenAndServe(":8088", nil))
}
