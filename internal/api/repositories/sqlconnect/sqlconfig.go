package sqlconnect

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)


func ConnectDB()(*sql.DB , error){
	
	user :=os.Getenv("DB_USER")
	password :=os.Getenv("DB_PASSWORD")
	dbname :=os.Getenv("DB_NAME")
	connectionString := fmt.Sprintf(
	"%s:%s@tcp(127.0.0.1:3306)/%s",
	user,
	password,
	dbname,
)


	 db ,err := sql.Open("mysql", connectionString)
	 if err != nil{
		// panic(err)
		return nil, err
	 }
 fmt.Println("Connected to mariaDB")
	 return db ,nil
}
