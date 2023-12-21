package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db  *gorm.DB
	err error
)

// Model
type Album struct {
	gorm.Model
	Title  string `json:"title"`
	Author string `json:"author"`
}

func init() {
	var (
		host     = getEnvVariable("DB_HOST")
		port     = getEnvVariable("DB_PORT")
		user     = getEnvVariable("DB_USER")
		dbname   = getEnvVariable("DB_NAME")
		password = getEnvVariable("DB_PASSWORD")
	)

	conn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		host,
		port,
		user,
		dbname,
		password,
	)
	db, err = gorm.Open(postgres.Open(conn), &gorm.Config{})
	db.AutoMigrate(Album{})

	if err != nil {
		log.Fatal(err)
	}
}

func getEnvVariable(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	return os.Getenv(key)
}

func PostAlbum(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var album Album
	// decode json to work on db
	json.NewDecoder(r.Body).Decode(&album)
	// save on db
	db.Create(&album)
	// back to json
	json.NewEncoder(w).Encode(album)
}

func GetAlbum(w http.ResponseWriter, r *http.Request) {
	var album Album
	id := mux.Vars(r)["id"]
	db.First(&album, id)
	if album.ID == 0 {
		json.NewEncoder(w).Encode("Album not found!")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(album)
}

func GetAlbums(w http.ResponseWriter, r *http.Request) {
	var albums []*Album
	db.Find(&albums)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(albums)

}

func UpdateAlbum(w http.ResponseWriter, r *http.Request) {
	var album Album
	id := mux.Vars(r)["id"]
	db.First(&album, id)
	if album.ID == 0 {
		json.NewEncoder(w).Encode("album not found!")
		return
	}
	json.NewDecoder(r.Body).Decode(&album)
	db.Save(&album)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(album)
}

func DeleteAlbum(w http.ResponseWriter, r *http.Request) {
	var album Album
	id := mux.Vars(r)["id"]
	db.First(&album, id)
	if album.ID == 0 {
		json.NewEncoder(w).Encode("Album not found!")
		return
	}
	db.Delete(&album, id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("album deleted successfully")
}

func Home() {

}

func main() {
	r := mux.NewRouter()
	// r.HandleFunc("/", Home).Methods("GET")
	r.HandleFunc("/api/v1/albums", PostAlbum).Methods("POST")
	r.HandleFunc("/api/v1/albums", GetAlbums).Methods("GET")
	r.HandleFunc("/api/v1/albums/{id}", GetAlbum).Methods("GET")
	r.HandleFunc("/api/v1/albums/{id}", UpdateAlbum).Methods("PUT")
	r.HandleFunc("/api/v1/albums/{id}", DeleteAlbum).Methods("DELETE")

	fmt.Println("Listening and serving at :5000")

	log.Fatal(http.ListenAndServe(":5000", r))
}
