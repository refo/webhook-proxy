package main

import (
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/go-chi/chi/v5"
	"gopkg.in/yaml.v2"
)

type Service struct {
	Name string   `yaml:"name"`
	URLs []string `yaml:"urls"`
}

type Config struct {
	Services []Service `yaml:"services"`
}

func loadConfig(filename string) (Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	config := Config{}
	if err := yaml.NewDecoder(file).Decode(&config); err != nil {
		return Config{}, err
	}
	return config, nil
}

func forwardRequest(serviceURLs []string, w http.ResponseWriter, r *http.Request) {

	var wg sync.WaitGroup

	for _, serviceURL := range serviceURLs {
		wg.Add(1)
		go func(serviceURL string) {
			defer wg.Done()
			handleForwarding(serviceURL, w, r)
		}(serviceURL)
	}
	wg.Wait()

	// Respond with a 200 OK
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Request broadcasted to all targets"))
}

func handleForwarding(serviceURL string, w http.ResponseWriter, r *http.Request) {

	var targetURL string = serviceURL + r.URL.Path
	log.Printf("Forwarding request from %s to %s", r.URL.Path, targetURL)

	req, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	copyHeaders(req, r)
	req.URL.RawQuery = r.URL.RawQuery

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		log.Printf("Error forwarding request to %s: %v", serviceURL, err)
		http.Error(w, "Failed to forward request", http.StatusBadGateway)
		return
	}
}

func copyHeaders(dest *http.Request, src *http.Request) {
	for key, values := range src.Header {
		for _, value := range values {
			dest.Header.Add(key, value)
		}
	}
}

func main() {
	config, err := loadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	r := chi.NewRouter()

	// Define routes for each service
	for _, service := range config.Services {
		log.Printf("Mounting service %s", service.Name)
		r.Mount("/"+service.Name, http.StripPrefix("/"+service.Name, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			forwardRequest(service.URLs, w, r)
		})))
	}

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
