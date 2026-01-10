package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// User struct with Uppercase fields (Exported) so the JSON package can see them.
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	// --- ROUTES (Go 1.22+ Syntax) ---

	// 1. Process Path Parameters (Wildcards)
	http.HandleFunc("GET /user/{id}", handleRequest)

	// 2. Process Everything Else (Body, Maps, Loops, Forms)
	http.HandleFunc("POST /process", handleRequest)

	fmt.Println("🚀 Server is flying on http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n--- New Request Received ---")

	// ==========================================
	// 1. PROCESSING PATH PARAMETERS
	// ==========================================
	id := r.PathValue("id")
	if id != "" {
		fmt.Printf("Path ID: %s\n", id)
	}

	// ==========================================
	// 2. PROCESSING QUERY PARAMETERS (?name=test)
	// ==========================================
	queryParams := r.URL.Query()
	search := queryParams.Get("search") // Gets specific key
	fmt.Printf("Query Search: %s\n", search)

	// ==========================================
	// 3. PROCESSING HEADERS (Metadata)
	// ==========================================
	token := r.Header.Get("Authorization")
	fmt.Printf("Auth Token: %s\n", token)

	// ==========================================
	// 4. PROCESSING BODY (JSON into MAP + LOOP)
	// ==========================================
	// We read the raw body bytes first so we can use them multiple times
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Can't read body", 400)
		return
	}
	defer r.Body.Close()

	// A. UNMARSHAL INTO MAP (For dynamic/unknown data)
	var dataMap map[string]interface{}
	json.Unmarshal(bodyBytes, &dataMap)

	// B. LOOPING OVER THE MAP
	fmt.Println("Looping over JSON Map:")
	for key, value := range dataMap {
		fmt.Printf("  -> %s : %v\n", key, value)
	}

	// C. UNMARSHAL INTO STRUCT (Best for fixed data)
	var userInstance User
	json.Unmarshal(bodyBytes, &userInstance)
	fmt.Printf("Struct Result: Name=%s, Email=%s\n", userInstance.Name, userInstance.Email)

	// ==========================================
	// 5. PROCESSING FORM DATA (Looping over r.Form)
	// ==========================================
	r.ParseForm() // Essential step to populate r.Form
	fmt.Println("Looping over Form/URL Values:")
	for key, values := range r.Form {
		// Form values are slices []string, so we loop again
		for _, v := range values {
			fmt.Printf("  -> %s = %s\n", key, v)
		}
	}
	//Accesr the request details
	// 1. BASIC NETWORK INFO
	fmt.Println("Method:", r.Method)         // GET, POST, etc.
	fmt.Println("Protocol:", r.Proto)        // HTTP/1.1 or HTTP/2
	fmt.Println("Address:", r.RemoteAddr)    // IP address of the user
	fmt.Println("Host:", r.Host)             // The domain (e.g., localhost:3000)

	// 2. URL DETAILS
	fmt.Println("Path:", r.URL.Path)         // The part after the domain (/user/profile)
	fmt.Println("Raw Query:", r.URL.RawQuery)// The string after '?' (name=john&age=20)
    fmt.Println("Scheme:", r.URL.Scheme)
	// 3. HEADERS (The 'Metadata')
	// Use .Get() for specific headers, or range to see all
	fmt.Println("User-Agent:", r.Header.Get("User-Agent")) 
	fmt.Println("Content-Type:", r.Header.Get("Content-Type"))

	// 4. FORM & QUERY DATA
	// You MUST call ParseForm() before you can see r.Form data
	r.ParseForm() 
	fmt.Println("All Form Data:", r.Form)           // Map of both URL query + Form body
	fmt.Println("Post Form Only:", r.PostForm)      // Map of ONLY Form body data
	
	// 5. COOKIES
	fmt.Println("Cookies:", r.Cookies())

	// 6. PATH VALUES (Go 1.22+)
	// If your route is "/user/{id}", this gets the ID
	fmt.Println("Path ID:", r.PathValue("id"))

	// 7. REQUEST CONTEXT
	// Useful for checking if a user cancelled the request (timeout)
	fmt.Println("Context:", r.Context())



	// ==========================================
	// 6. SENDING THE RESPONSE (MARSHAL)
	// ==========================================
	response := map[string]string{
		"status":  "success",
		"message": "We processed everything!",
	}
	
	// Convert Go map to JSON bytes (Marshal)
	finalJson, _ := json.Marshal(response)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(finalJson)
}