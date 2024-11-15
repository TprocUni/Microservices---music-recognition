package main

//imports
import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

// datatypes
// represents the database
type Repository struct {
	DB *sql.DB
}

// represents a track
type Track struct {
	Id    string `json:"id"`
	Audio string `json:"audio"`
}

// global variables
var repo Repository

//Functions
//These are all database operations and will be used by the functions that
//communicate with the cmd line

// Initialisation for database
func Init() {
	if db, err := sql.Open("sqlite3", "trackDB.db"); err == nil {
		repo = Repository{DB: db}
	} else {
		log.Fatal("Database initialisation failed.")
	}
}

// Create table in db
// With columns Id (PK) and Audio
func Create() int {
	sql := " CREATE TABLE IF NOT EXISTS Tracks " +
		" ( Id TEXT PRIMARY KEY , Audio TEXT NOT NULL) "
	if _, err := repo.DB.Exec(sql); err == nil {
		return 0
	} else {
		return -1
	}
}

// Clearing table
func Clear() int {
	sql := " DELETE FROM Tracks "
	if _, err := repo.DB.Exec(sql); err == nil {
		return 0
	} else {
		return -1
	}
}

// Inserting into DB
func Insert(c Track) int64 {
	sql := " INSERT INTO Tracks ( Id , Audio ) " +
		" VALUES (? , ?) "
	if stmt, err := repo.DB.Prepare(sql); err == nil {
		defer stmt.Close()
		if res, err := stmt.Exec(c.Id, c.Audio); err == nil {
			if n, err := res.RowsAffected(); err == nil {
				return n
			}
		}
	}
	return -1
}

// Reading from DB
func Read(id string) (Track, int64) {
	sql := " SELECT * FROM Tracks WHERE Id = ? "
	if stmt, err := repo.DB.Prepare(sql); err == nil {
		defer stmt.Close()
		var c Track
		row := stmt.QueryRow(id)
		if err := row.Scan(&c.Id, &c.Audio); err == nil {
			return c, 1
		} else {
			return Track{}, 0
		}
	}
	return Track{}, -1
}

// Updating DB
func Update(c Track) int64 {
	sql := " UPDATE Tracks SET Audio = ? " +
		" WHERE id = ? "
	if stmt, err := repo.DB.Prepare(sql); err == nil {
		defer stmt.Close()
		if res, err := stmt.Exec(c.Audio, c.Id); err == nil {
			if n, err := res.RowsAffected(); err == nil {
				return n
			}
		}
	}
	return -1
}

// Deleting from the DB
func Delete(Id string) int64 {
	sql := " DELETE FROM Tracks WHERE Id = ? "
	if stmt, err := repo.DB.Prepare(sql); err == nil {
		defer stmt.Close()
		if res, err := stmt.Exec(Id); err == nil {
			if n, err := res.RowsAffected(); err == nil {
				return n
			}
		}
	}
	return -1
}

// Handler functions
// These are the functions that will be called by the cmd line
// ----
// PUT handler
// Input: Id
// Output:
//
//	    201 - track inserted successfully
//		204 - track updated successfully
//		400 - bad request
//		500 - operation failed
func putHandler(w http.ResponseWriter, r *http.Request) {
	//get the id from the url
	Id, exists := mux.Vars(r)["id"]

	//check if the id exists
	if !exists {
		http.Error(w, "ID not found", http.StatusBadRequest)
		return
	}

	//create track structure to store the audio
	var trackC Track
	//get the audio from the request body
	//check if there was an error
	if err := json.NewDecoder(r.Body).Decode(&trackC); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//update the values of the track
	a := trackC.Audio

	trackC.Id = Id
	trackC.Audio = string(a)

	updated := false

	//attempt to insert the track into the database
	n := Insert(trackC)
	//if track already exists, update the track
	if n == 1 {
		//pass
	} else if n == -1 {
		n := Update(trackC)
		//if successful, track was updated and therefore must have existed
		if n == 1 {
			updated = true
		} else {
			//track couldn't update -> appropriate error
			http.Error(w, " Track update failed", http.StatusInternalServerError)
			return
		}
	} else {
		//track couldn't insert -> appropriate error
		http.Error(w, " Track insert failed", http.StatusInternalServerError)
	}

	//if the track was inserted, return 201
	if !updated {
		w.WriteHeader(http.StatusCreated)
		return
	}

	//if the track was updated, return 204
	if updated {
		w.WriteHeader(http.StatusNoContent)
		return
	}
}

// DELETE handler
// Input: Id
// Output:
//
//	204 - track deleted successfully
//	404 - track not found
//	500 - server errror
func deleteHandler(w http.ResponseWriter, r *http.Request) {
	//get the id from the url
	Id := mux.Vars(r)["id"]
	//delete the track from the database using Id
	n := Delete(Id)
	//if the track was deleted, return 204
	if n == 1 {
		w.WriteHeader(http.StatusNoContent)
		return
	} else if n == 0 {
		//if the track wasn't deleted, return 404
		http.Error(w, " Track not found", http.StatusNotFound)
	} else {
		//if the track couldn't be deleted, return 500
		http.Error(w, " Track delete failed", http.StatusInternalServerError)
	}

}

// GET all handler
// Input: none
// Output:
//
//	200 - track read successfully
//	400 - bad request
//	500 - operation failed
func getAllHandler(w http.ResponseWriter, r *http.Request) {
	//get all the tracks from the database
	rows, err := repo.DB.Query("SELECT * FROM Tracks")
	//check if there was an error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//close the rows
	defer rows.Close()

	//create a list of tracks - using slices since they accept variable length args
	tracks := []Track{}

	//iterate through the rows
	for rows.Next() {
		//add the track to the list
		track := Track{}
		if err := rows.Scan(&track.Id, &track.Audio); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError) //500
			return
		}
		tracks = append(tracks, track)
	}

	//check if there was an error - 500
	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//response creation
	//set response header to 200
	w.WriteHeader(http.StatusOK)
	//formatting the return string
	w.Write([]byte("\n["))
	for _, track := range tracks {
		w.Write([]byte("\"" + track.Id + "\""))
		if track.Id != tracks[len(tracks)-1].Id {
			w.Write([]byte("; "))
		} else {
			w.Write([]byte("]\n"))
		}
	}

}

// GET handler
// Input: Id
// Output:
//
//	200 - track read successfully
//	400 - bad request
//	500 - operation failed
func getHandler(w http.ResponseWriter, r *http.Request) {
	//get the id from the url
	Id := mux.Vars(r)["id"]

	//get the track from the database using Id
	track, n := Read(Id)

	//if the track was read, return 200
	if n == 1 {
		//decode track into json form
		audio, err := json.Marshal(track)
		//check if there was an error
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(audio)
		w.Write([]byte("\n"))
	} else if n == 0 {
		//return track not found error
		http.Error(w, " Track not found", http.StatusNotFound)
		return
	} else {
		//return server error
		http.Error(w, " Track read failed", http.StatusInternalServerError)
		return
	}

}

// Main function
func main() {
	//Initialise the router
	router := mux.NewRouter()
	//Initialise the database
	Init()
	//Create the table
	Create()
	//Clear the table
	Clear()

	//Adding handlers
	router.HandleFunc("/tracks/{id}", putHandler).Methods("PUT")
	router.HandleFunc("/tracks/{id}", deleteHandler).Methods("DELETE")
	router.HandleFunc("/tracks", getAllHandler).Methods("GET")
	router.HandleFunc("/tracks/{id}", getHandler).Methods("GET")

	//Start the server
	log.Fatal(http.ListenAndServe(":3000", router))

}

/*TODO
  - Add functionality to stay open and recieve requests
       - do for each:
			         - delete one
					 - get one
                     - get all
					 - put one
*/
