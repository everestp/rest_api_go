package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	

	"github.com/everestp/rest_api_go/internal/api/models"
	"github.com/everestp/rest_api_go/internal/api/repositories/sqlconnect"
	"github.com/everestp/rest_api_go/pkg/utils"
)



// Initialize dummy data
func init() {
	teachers[nextID] = models.Teacher{
		ID:        nextID,
		FirstName: "jhon",
		LastName:  "Doe",
		Level:     "9A",
		Subject:   "Math",
	}
	nextID++

	teachers[nextID] = models.Teacher{
		ID:        nextID,
		FirstName: "Everest",
		LastName:  "Paudel",
		Level:     "Bsc",
		Subject:   "Math",
	}
	nextID++
}



func GetStudentsHandler(w http.ResponseWriter, r *http.Request) {
	


	var teachers []models.Teacher


		teachers, err := sqlconnect.GetTeacherDBHandler(teachers ,r)
		if err != nil {
			return
		}

		response := struct {
			Status string           `json:"status"`
			Count  int              `json:"count"`
			Data   []models.Teacher `json:"data"`
		}{
			Status: "success",
			Count:  len(teachers),
			Data:   teachers,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	
}


func GetOneStudentHandler(w http.ResponseWriter, r *http.Request) {


	
	idStr := r.PathValue("id")
	fmt.Println(idStr)
	

	// GET BY ID
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid teacher ID", http.StatusBadRequest)
		return
	}
		teacher, err := sqlconnect.GetTeacherByID(id)
		if err != nil {
			log.Println(err)
			return
		}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teacher)
}





func AddStudentHandler(w http.ResponseWriter, r *http.Request) {
	
	defer r.Body.Close()

	var newStudents []models.Student

	if err := json.NewDecoder(r.Body).Decode(&newStudents); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	addStudents, err := sqlconnect.AddStudentsDBHandler(newStudents)
	if err != nil {
		log.Println(err)
		return 
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Student `json:"data"`
	}{
		Status: "success",
		Count:  len(addStudents),
		Data:   addStudents,
	}

	json.NewEncoder(w).Encode(response)
}


// PUT for teacher Route /teacher/
// PUT /teachers/{id}
func UpdateStudentHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid Teacher Id", http.StatusBadRequest)
		return
	}

	var updatedTeacher models.Teacher
	err = json.NewDecoder(r.Body).Decode(&updatedTeacher)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid Request Payload", http.StatusBadRequest)
		return
	}

	updatedTeacherFromDB, err := sqlconnect.UpdateTeacher(id, updatedTeacher)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTeacherFromDB)
}





// PATCH /teachers
func PatchStudentsHandler(w http.ResponseWriter, r *http.Request) {

	var updates []map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = sqlconnect.PatchTeachers(updates)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// PATCH /teachers/{id}
func PatchOneStudentHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid Teacher Id", http.StatusBadRequest)
		return
	}

	var updates map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid Request Payload", http.StatusBadRequest)
		return
	}

	// Apply updates using reflect
	updatedTeacher, err := sqlconnect.PatchOneTeacher(id, updates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTeacher)

}



func DeleteStudentsHandler(w http.ResponseWriter , r *http.Request){
	
	db, err := sqlconnect.ConnectDB()
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable to connect to database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

var ids []int
	err =json.NewDecoder(r.Body).Decode(&ids)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid Request Payload", http.StatusBadRequest)
		return
	}

    tx , err := db.Begin()
	if err != nil{
		log.Println(err)
		http.Error(w, "Error starting the transaction", http.StatusInternalServerError)
		return
	}
stmt, err := db.Prepare("DELETE FROM teacher WHERE id = ?")
	if err != nil {
		fmt.Println("SQL Prepare Error:", err)
		tx.Rollback()
		http.Error(w, "Error preparing delete statement", http.StatusInternalServerError)
		return
	}

	defer stmt.Close()

  deleteID := []int{}
  for _ ,id := range ids{
	result, err := stmt.Exec(id)
	if err != nil{
		tx.Rollback()
		log.Println(err)
		http.Error(w, "Error deleting the teacher ", http.StatusInternalServerError)
		return
	}
	 rowAffected , err := result.RowsAffected()
 if err != nil {
		log.Println(err)
		tx.Rollback()
		http.Error(w, "Error  retriving delete result", http.StatusInternalServerError)
		return
	}
	//if  teacher was deleted add the id to the deleted id
	if rowAffected > 0 {
		deleteID = append(deleteID, id)
	}
	if rowAffected < 1 {
		tx.Rollback()

		return 
	}

  }
	//Commit the  transaction
	err = tx.Commit()
	 if err != nil {
		log.Println(err)
		tx.Rollback()
		http.Error(w, "Error commiting transaction", http.StatusInternalServerError)
		return
	}
	if len(deleteID) < 1{
	http.Error(w, "ID doesnot exist", http.StatusNoContent)
	}
	w.Header().Set("Content-Type", "application/json")
	response := struct{
		Status string `json:"status"`
		deletedIDs []int `josn:"deletes_id"`
	}{
		Status: "Teacher Delete Sucessfully",
		deletedIDs: deleteID,
	}
	json.NewEncoder(w).Encode(response)
}
//DELETE for techer/{id}

func DeleteOneStudentHandler(w http.ResponseWriter , r *http.Request){
	idStr := strings.TrimPrefix(r.URL.Path, "/teacher/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}


	db, err := sqlconnect.ConnectDB()
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable to connect to database", http.StatusInternalServerError)
		return
	}
	defer db.Close()
	result, err := db.Exec(
		"DELETE FROM teacher WHERE id = ?",id)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error deleting  teacher teacher", http.StatusInternalServerError)
		return
	}
fmt.Println(result.RowsAffected())
 rowAffected , err := result.RowsAffected()
 if err != nil {
		log.Println(err)
		http.Error(w, "Error  retriving delete result", http.StatusInternalServerError)
		return
	}
	if rowAffected == 0 {
		http.Error(w, "No teacher found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	//Response body
	w.Header().Set("Content-Type", "application/json")
	response := struct{
		Status string `json:"status"`
		ID int `josn:"id"`
	}{
		Status: "Teacher Delete Sucessfully",
		ID: id,
	}
	json.NewEncoder(w).Encode(response)
}


















func PatchTeachers(updates []map[string]interface{}) error {
	db, err := sqlconnect.ConnectDB()
	if err != nil {
		return utils.ErrorHandler(err, "error updating data")
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return utils.ErrorHandler(err, "error updating data")
	}

	for _, update := range updates {
		idStr, ok := update["id"].(string)
		if !ok {
			tx.Rollback()
			return utils.ErrorHandler(err, "invalid Id")
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			tx.Rollback()
			return utils.ErrorHandler(err, "invalid Id")
		}

		var teacherFromDb models.Teacher
		err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id).Scan(&teacherFromDb.ID, &teacherFromDb.FirstName, &teacherFromDb.LastName, &teacherFromDb.Email, &teacherFromDb.Level, &teacherFromDb.Subject)
		if err != nil {
			log.Println("ID:", id)
			log.Printf("Type: %T", id)
			log.Println(err)
			tx.Rollback()
			if err == sql.ErrNoRows {
				return utils.ErrorHandler(err, "Teacher not found")
			}
			return utils.ErrorHandler(err, "error updating data")
		}

		teacherVal := reflect.ValueOf(&teacherFromDb).Elem()
		teacherType := teacherVal.Type()

		for k, v := range update {
			if k == "id" {
				continue // skip updating the ID field
			}
			for i := 0; i < teacherVal.NumField(); i++ {
				field := teacherType.Field(i)
				if field.Tag.Get("json") == k+",omitempty" {
					fieldVal := teacherVal.Field(i)
					if fieldVal.CanSet() {
						val := reflect.ValueOf(v)
						if val.Type().ConvertibleTo(fieldVal.Type()) {
							fieldVal.Set(val.Convert(fieldVal.Type()))
						} else {
							tx.Rollback()
							log.Printf("cannot convert %v to %v", val.Type(), fieldVal.Type())
							return utils.ErrorHandler(err, "error updating data")
						}
					}
					break
				}
			}
		}

		_, err = tx.Exec("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?", teacherFromDb.FirstName, teacherFromDb.LastName, teacherFromDb.Email, teacherFromDb.Level, teacherFromDb.Subject, teacherFromDb.ID)
		if err != nil {
			tx.Rollback()
			return utils.ErrorHandler(err, "error updating data")
		}
	}

	err = tx.Commit()
	if err != nil {
		return utils.ErrorHandler(err, "error updating data")
	}
	return nil 
}