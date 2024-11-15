package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
)

type RequestBody struct {
	Audio string `json:"Audio"`
}

type SearchRequest struct {
	Audio string `json:"Audio"`
}

type SearchResponse struct {
	Id string `json:"id"`
}

const (
	port       = ":3002"
	localHostS = "http://localhost:3001/search"
)

func cooltownHandler(w http.ResponseWriter, r *http.Request) {
	//get fragment of audio from client
	// Parse the request body as JSON
	var rBody RequestBody
	err := json.NewDecoder(r.Body).Decode(&rBody)
	//check for errors
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//Get the audio data from base64
	audio := rBody.Audio

	//create request
	searchRequest := SearchRequest{Audio: audio}
	searchBytes, err := json.Marshal(searchRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//send a fragment of audio to audd.io using search mircoservice
	//set up client connection
	response1, err := http.Post(localHostS, "application/json", bytes.NewBuffer(searchBytes))
	//error checking
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//end connection
	defer response1.Body.Close()

	//get the data from the response
	response2, err := io.ReadAll(response1.Body)
	//check for errors
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if the response is Not Found - 404
	if string(response2) == "Not Found\n" {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	//---------------------------------------------------------------------------------------
	// Convert the byte slice to JSON
	var searchResponse SearchResponse
	err = json.Unmarshal(response2, &searchResponse)
	//error checking
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	idValue := searchResponse.Id

	//fmt.Println("Type of Id:", idType.Name())

	//ALter the response to be in the correct format
	// Convert the JSON to a string and replace spaces with "+"
	id := idValue

	id = strings.ReplaceAll(id, " ", "+")
	// Escape the ID value using url.PathEscape

	escaped := url.PathEscape(id)

	//create request
	theTrack := fmt.Sprintf("http://localhost:3000/tracks/%s", escaped)

	//send the track id to the database using tracks
	// Make the GET request
	trackResponse, err := http.Get(theTrack)
	//error checking
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//end connection
	defer trackResponse.Body.Close()

	/// Read the response body into a byte form, then string
	trackBody, err := io.ReadAll(trackResponse.Body)
	//erorr
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	trackBodyString := string(trackBody)
	//send the track object to the client
	// Print track body
	w.Write([]byte(trackBodyString))
}

func main() {
	// Start router
	router := mux.NewRouter()

	// Handlers
	router.HandleFunc("/cooltown", cooltownHandler).Methods("POST")

	// Start server
	http.ListenAndServe(":3002", router)
}
