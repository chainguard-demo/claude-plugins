package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Component struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	License  string `json:"license"`
	Supplier string `json:"supplier"`
}

var components = []Component{
	{Name: "openssl", Version: "3.0.8", License: "Apache-2.0", Supplier: "OpenSSL Foundation"},
	{Name: "curl", Version: "8.4.0", License: "MIT", Supplier: "curl contributors"},
	{Name: "zlib", Version: "1.3.1", License: "Zlib", Supplier: "zlib authors"},
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"service": "sbom-scanner",
			"status":  "ok",
			"version": "1.0.0",
		})
	})

	http.HandleFunc("/components", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(components)
	})

	log.Println("sbom-scanner listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
