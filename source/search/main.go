package main

//imports
import (
	"bytes"
	"encoding/json"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/gorilla/mux"
)

// Structures
type Track struct {
	Id string `json:"id"`
}

// Constants for program
const (
	//The port number to listen on
	port = ":3001"
	//API key for audd.io
	APIKey = "YOUR_API_KEY_HERE"
	//URL for audd.io
	URL = "https://api.audd.io/"
)

// SearchHandler handles requests to /search
// Inputs: HTTP request with audio data
// Outputs: HTTP response with track name/ID
func searchHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body as JSON
	var requestBody struct {
		Audio string `json:"Audio"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	//error checking
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//---------------------------------------------------------------------
	// Decode the audio data from base64
	/*
		audio, err := base64.StdEncoding.DecodeString(requestBody.Audio)
		//error checking
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	*/
	audio := requestBody.Audio

	// Send the HTTP request to audd.io
	// Build form data with API token and base64-encoded audio
	//convert to base64
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("api_token", APIKey)
	writer.WriteField("audio", audio)
	writer.Close()

	// Build request with data
	request, err := http.NewRequest("POST", URL, body)
	//check for errors
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Set the request headers
	request.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request
	client := &http.Client{}
	response, err := client.Do(request)
	//check for errors
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	//close connection
	defer response.Body.Close()

	var decoded map[string]interface{}
	// Read the response from audd.io
	if response.StatusCode == http.StatusOK {
		// Try to map response to struct
		// Try to decode response, if it fails, return error
		if err := json.NewDecoder(response.Body).Decode(&decoded); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// result is decoded

	//------------------------------------------------------------------------
	// Parse the response as a Track object
	var track Track
	track.Id = decoded["result"].(map[string]interface{})["title"].(string)

	//Error checking
	if track.Id == "" {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Write the track object as JSON to the response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(track)

}

func main() {
	// Start router
	router := mux.NewRouter()

	// Handlers
	router.HandleFunc("/search", searchHandler).Methods("POST")

	log.Fatal(http.ListenAndServe(":3001", router))

}
