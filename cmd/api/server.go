package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// User struct with Uppercase fields (Exported) so the JSON package can see them.
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
 func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	path := strings.TrimPrefix(r.URL.Path, "/teacher/")
	userID := strings.TrimSuffix(path, "/")
	fmt.Println("The USer ID is", userID)
		w.Write([]byte("Hello Route Route"))
	}


	 func  teacherHandler(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte["Hello Teacher Route"])
	}

	 func  execsHandler(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte["Hello student Route"])
	}
func main() {

	http.HandleFunc("/",rootHandler)
	http.HandleFunc("/teacher",teacherHandler)
	http.HandleFunc("/student",)
	http.HandleFunc("/exec", execsHandler)

	
	fmt.Println("Server is Running on http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

