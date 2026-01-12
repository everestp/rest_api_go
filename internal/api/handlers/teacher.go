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
	"sync"

	"github.com/everestp/rest_api_go/internal/api/models"
	"github.com/everestp/rest_api_go/internal/api/repositories/sqlconnect"
)

var (
	teachers = make(map[int]models.Teacher)
	mutex    = &sync.Mutex{}
	nextID   = 1
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


func isValidSortOrder(order string) bool {
	return order == "asc" || order == "desc"
}

func isValidSortField(field string) bool {
	validField := map[string]bool{
		"first_name": true,
		"last_name":  true,
		"email":      true,
		"level":      true,
		"subject":    true,
	}
	return validField[field]
}

func GetTeachersHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sqlconnect.ConnectDB()
	if err != nil {
		http.Error(w, "Error connecting to database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

		query := "SELECT id, first_name, last_name, email, level, subject FROM teacher WHERE 1=1"
		var args []any

		// Add filters
		query, args = addFilters(r, query, args)

		// Add sorting
		query = addSorting(r, query)

		rows, err := db.Query(query, args...)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Invalid teacher ID", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		teacherList := make([]models.Teacher, 0)

		for rows.Next() {
			var teacher models.Teacher
			err := rows.Scan(
				&teacher.ID,
				&teacher.FirstName,
				&teacher.LastName,
				&teacher.Email,
				&teacher.Level,
				&teacher.Subject,
			)
			if err != nil {
				fmt.Println(err)
				http.Error(w, "Error scanning database", http.StatusInternalServerError)
				return
			}
			teacherList = append(teacherList, teacher)
		}

		response := struct {
			Status string           `json:"status"`
			Count  int              `json:"count"`
			Data   []models.Teacher `json:"data"`
		}{
			Status: "success",
			Count:  len(teacherList),
			Data:   teacherList,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	
}
func GetTeacherHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sqlconnect.ConnectDB()
	if err != nil {
		http.Error(w, "Error connecting to database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	
	idStr := r.PathValue("id")
	fmt.Println(idStr)
	

	// GET BY ID
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid teacher ID", http.StatusBadRequest)
		return
	}

	var teacher models.Teacher

	err = db.QueryRow(
		"SELECT id, first_name, last_name, email, level, subject FROM teacher WHERE id = ?",
		id,
	).Scan(
		&teacher.ID,
		&teacher.FirstName,
		&teacher.LastName,
		&teacher.Email,
		&teacher.Level,
		&teacher.Subject,
	)

	if err == sql.ErrNoRows {
		http.Error(w, "Teacher not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Database query error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teacher)
}

func addSorting(r *http.Request, query string) string {
	sortParams := r.URL.Query()["sortby"]
	// /teacher?sortby=name:asc&sortby=level:desc
	if len(sortParams) > 0 {
		query += " ORDER BY"
		for i, param := range sortParams {
			parts := strings.Split(param, ":")
			if len(parts) != 2 {
				continue
			}
			field, order := parts[0], parts[1]
			if !isValidSortOrder(order) || !isValidSortField(field) {
				continue
			}
			if i > 0 {
				query += ","
			}
			query += " " + field + " " + strings.ToUpper(order)
		}
	}
	return query
}

func addFilters(r *http.Request, query string, args []any) (string, []any) {
	params := map[string]string{
		"first_name": "first_name",
		"last_name":  "last_name",
		"email":      "email",
		"level":      "level",
		"subject":    "subject",
	}
	for param, dbField := range params {
		value := r.URL.Query().Get(param)
		if value != "" {
			query += " AND " + dbField + " = ? "
			args = append(args, value)
		}
	}
	return query, args
}

func AddTeacherHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sqlconnect.ConnectDB()
	if err != nil {
		http.Error(w, "Error connecting to database", http.StatusInternalServerError)
		return
	}
	defer db.Close()
	defer r.Body.Close()

	var newTeachers []models.Teacher

	if err := json.NewDecoder(r.Body).Decode(&newTeachers); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	stmt, err := db.Prepare(`
    INSERT INTO teacher
    (first_name, last_name, email, level, subject)
    VALUES (?,?,?,?,?)
`)
	if err != nil {
		fmt.Println("SQL Prepare Error:", err)
		http.Error(w, "Error preparing SQL statement", http.StatusInternalServerError)
		return
	}

	defer stmt.Close()

	addedTeachers := make([]models.Teacher, len(newTeachers))

	for i, newTeacher := range newTeachers {
		res, err := stmt.Exec(
			newTeacher.FirstName,
			newTeacher.LastName,
			newTeacher.Email,
			newTeacher.Level,
			newTeacher.Subject,
		)
		if err != nil {
			http.Error(w, "Error inserting teacher", http.StatusInternalServerError)
			return
		}

		lastID, err := res.LastInsertId()
		if err != nil {
			http.Error(w, "Error getting inserted ID", http.StatusInternalServerError)
			return
		}

		newTeacher.ID = int(lastID)
		addedTeachers[i] = newTeacher
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Teacher `json:"data"`
	}{
		Status: "success",
		Count:  len(addedTeachers),
		Data:   addedTeachers,
	}

	json.NewEncoder(w).Encode(response)
}

// PUT for teacher Route /teacher/
func UpdateTeacherHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/teacher/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var updateTeacher models.Teacher
	err = json.NewDecoder(r.Body).Decode(&updateTeacher)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid Request Payload", http.StatusBadRequest)
		return
	}

	db, err := sqlconnect.ConnectDB()
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable to connect to database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var existingTeacher models.Teacher
	err = db.QueryRow("SELECT id, first_name, last_name, email, level, subject FROM teacher WHERE id = ?", id).Scan(
		&existingTeacher.ID,
		&existingTeacher.FirstName,
		&existingTeacher.LastName,
		&existingTeacher.Email,
		&existingTeacher.Level,
		&existingTeacher.Subject,
	)
	if err == sql.ErrNoRows {
		log.Println(err)
		http.Error(w, "Teacher not found", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable to retrieve data", http.StatusInternalServerError)
		return
	}

	updateTeacher.ID = existingTeacher.ID
	_, err = db.Exec(
		"UPDATE teacher SET first_name = ?, last_name = ?, email = ?, level = ?, subject = ? WHERE id = ?",
		updateTeacher.FirstName,
		updateTeacher.LastName,
		updateTeacher.Email,
		updateTeacher.Level,
		updateTeacher.Subject,
		updateTeacher.ID,
	)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error updating teacher", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updateTeacher)
}

//Patch for teacher/{id}
func PatchTeacherHandler(w http.ResponseWriter , r *http.Request ){
	idStr := strings.TrimPrefix(r.URL.Path, "/teacher/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	 var updates map[string]any
	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid Request Payload", http.StatusBadRequest)
		return
	}

	db, err := sqlconnect.ConnectDB()
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable to connect to database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var existingTeacher models.Teacher
	err = db.QueryRow("SELECT id, first_name, last_name, email, level, subject FROM teacher WHERE id = ?", id).Scan(
		&existingTeacher.ID,
		&existingTeacher.FirstName,
		&existingTeacher.LastName,
		&existingTeacher.Email,
		&existingTeacher.Level,
		&existingTeacher.Subject,
	)
	if err == sql.ErrNoRows {
		log.Println(err)
		http.Error(w, "Teacher not found", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable to retrieve data", http.StatusInternalServerError)
		return
	}
  // Apply updates
//   for k , v := range updates {
// 	switch k{
// 	case "first_name":
// 		existingTeacher.FirstName = v.(string)
// 		case "last_name":
// 		existingTeacher.LastName = v.(string)
// 		case "email":
// 		existingTeacher.Email = v.(string)
// 		case "subject":
// 		existingTeacher.Subject = v.(string)
// 	}
//   }
  //Apply update using Reflect
// updates is map[string]interface{}
teacherVal := reflect.ValueOf(&existingTeacher).Elem()
teacherType := teacherVal.Type()

for k, v := range updates {
    for i := 0; i < teacherVal.NumField(); i++ {
        fieldStruct := teacherType.Field(i)
        fieldVal := teacherVal.Field(i)
        tagName := strings.Split(fieldStruct.Tag.Get("json"), ",")[0]

        if tagName == k {
            if fieldVal.CanSet() {
                val := reflect.ValueOf(v)

                // HANDLE JSON NUMBERS: Convert float64 to int/int64 if needed
                if val.Kind() == reflect.Float64 {
                    switch fieldVal.Kind() {
                    case reflect.Int, reflect.Int64:
                        fieldVal.SetInt(int64(v.(float64)))
                        goto NextKey // Move to next update key
                    case reflect.Int32:
                        fieldVal.SetInt(int64(v.(float64)))
                        goto NextKey
                    }
                }

                if val.Type().ConvertibleTo(fieldVal.Type()) {
                    fieldVal.Set(val.Convert(fieldVal.Type()))
                }
            }
            break
        }
    }
NextKey:
}
	_, err = db.Exec(
		"UPDATE teacher SET first_name = ?, last_name = ?, email = ?, level = ?, subject = ? WHERE id = ?",
		existingTeacher.FirstName,
		existingTeacher.LastName,
		existingTeacher.Email,
		existingTeacher.Level,
		existingTeacher.Subject,
		existingTeacher.ID,
	)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error updating teacher", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existingTeacher)
}
func DeleteTeachersHandler(w http.ResponseWriter , r *http.Request){
	
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

func DeleteOneTeacherHandler(w http.ResponseWriter , r *http.Request){
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