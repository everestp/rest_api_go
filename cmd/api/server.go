package main

import (
	"fmt"
	"log"
	"net/http"

	"time"

	mw "github.com/everestp/rest_api_go/internal/api/middlewares"
	"github.com/everestp/rest_api_go/internal/api/repositories/sqlconnect"
	"github.com/everestp/rest_api_go/internal/api/router"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// teacherHandler demonstrates dynamic path parsing.

func main() {
	err:= godotenv.Load()
	if err != nil {
		return 
	}
	 _ ,err = sqlconnect.ConnectDB()
	  if err != nil{
		// panic(err)
		fmt.Println("Error ------ ",err)
		return 
	 }
	// Logic Block: Configuration
	// We no longer need certFile or keyFile constants.
	const port = ":3000"


   rl := mw.NewRateLimiter(5, time.Minute)
     hppOptions := mw.HPPOptions{
		CheckQuery: true,
		CheckBody: true,
		CheckBodyOnlyForContentType: "application/x-www-form-urlencoded",
		Whitelist: []string{"name"},
	 }
	 fmt.Println(rl ,hppOptions)
	// Logic Block: Middleware Onion
	// The order remains the same: Timing -> Compression -> Security -> CORS -> App
	// secureMux1 := applyMiddlewares(mux, mw.Hpp(hppOptions) ,mw.Compression , mw.SecurityHeaders , mw.ResponseTimeMiddleware , rl.Middleware ,mw.Cors)
   secureMux := mw.SecurityHeaders(router.MainRouter())
	// Logic Block: Server Initialization
	// We removed the TLSConfig field.
	server := &http.Server{
		Addr:         port,
		Handler:      secureMux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Future Reference: HTTP vs HTTPS
	// Without TLS, data is sent in 'Plaintext'. Anyone on the network 
	// can see the traffic. This is fine for local dev but never for production.
	fmt.Printf("🚀 Server running on http://localhost%s\n", port)

	// Logic Block: Standard Execution
	// We use ListenAndServe() instead of ListenAndServeTLS()
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Critical Server Failure: %v", err)
	}
}

