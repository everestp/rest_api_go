package main

import (
	"fmt"
	"log"
	"net/http"
)


func main(){
	port := ":3000"

	http.HandleFunc("/",  func(w http.ResponseWriter, r *http.Request){
		// fmt.Fprintf(w, "Hello Root Route")
		w.Write([]byte("Hello world"))
		fmt.Println(w, "Hello Root Route")
	})
	http.HandleFunc("/teacher",  func(w http.ResponseWriter, r *http.Request){
		// fmt.Fprintf(w, "Hello Root Route")
		w.Write([]byte("Hello TEacher"))
		fmt.Println(w, "Hello Teacher Route")
	})
	http.HandleFunc("/student",  func(w http.ResponseWriter, r *http.Request){
		// fmt.Fprintf(w, "Hello Root Route")
		w.Write([]byte("Hello Student"))
		fmt.Println(w, "Hello Student Route")
	})
	fmt.Println("Server is  running on port",port)
	err := http.ListenAndServe(port,nil)
	if err != nil{
		log.Fatalln("Error starting server",err)
	}
}