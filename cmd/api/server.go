package main

import (
	"crypto/tls"

	"fmt"

	"log"
	"net/http"
	"strings"

	 mw "github.com/everestp/rest_api_go/internal/api/middlewares"
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
		w.Write([]byte("Hello Teacher Route"))
	}

	 func  execsHandler(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello student Route"))
	}
func main() {
    cert :="cert.pem"
	key := "key.pem"
	port :=":3000"

	mux := http.NewServeMux()

	mux.HandleFunc("/",rootHandler)
	mux.HandleFunc("/teacher",teacherHandler)
	mux.HandleFunc("/student",)
	mux.HandleFunc("/exec", execsHandler)


	  tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	  }

	  //Create a custom server\
	  server := &http.Server{
		Addr: port,
		Handler: mw.SecurityHeaders(mw.Cors(mux)),
		TLSConfig: tlsConfig,
	  }
	
	fmt.Println("Server is Running on http://localhost:3000")
	 err :=server.ListenAndServeTLS(cert, key)
	 if err != nil{
		log.Fatalln("Error starting the server",err)
	 }
}

