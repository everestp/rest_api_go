package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
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
		FirstNAme: "jhon",
		LastName:  "Doe",
		Level:     "9A",
		Subject:   "Math",
	}
	nextID++

	teachers[nextID] = models.Teacher{
		ID:        nextID,
		FirstNAme: "Everest",
		LastName:  "Paudel",
		Level:     "Bsc",
		Subject:   "Math",
	}
	nextID++
}

func TeacherHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTeachersHandler(w, r)

	case http.MethodPost:
		addTeacherHandler(w, r)

	case http.MethodPut:
		w.Write([]byte("General Teacher Directory"))

	case http.MethodDelete:
		w.Write([]byte("General Teacher Directory"))
	}
}

func isValidSortOrder(order string) bool{
	return  order =="asc" || order =="desc"
}

func isValidSortFiled(field string) bool{
	validFiled := map[string]bool{
		"first_name": true,
		"last_name":true,
		"email":true,
		"level":true,
		"subject":true,
	}
	return  validFiled[field]
}


func getTeachersHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sqlconnect.ConnectDB()
	if err != nil {
		http.Error(w, "Error connecting to database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	path := strings.TrimPrefix(r.URL.Path, "/teacher/")
	idStr := strings.TrimSuffix(path, "")
	fmt.Println(idStr)

	// GET ALL
	if idStr == "" {
query := "SELECT id, first_name, last_name, email, level, subject FROM teacher WHERE 1=1"
var args []any

//list all the query params
query, args = addFilters(r, query, args)

//sorting
sortParams := r.URL.Query()["sortby"]
//  /teacher?sortby=name:asc&sortby=level:desc
if len(sortParams) > 0{
    query += "ORDER BY "
	for i, param := range sortParams{
	  parts :=strings.Split(param, ":")
	  if len(parts) !=2{
		continue
	  }
      field ,order :=parts[0] ,parts[1]
	  if !isValidSortFiled(field) || isValidSortFiled(order){
		continue
	  }
	  if i > 0 {
		query += ","
	  }
	  query +=" "+ field + " "+ order


	}
}
 rows , err:= db.Query(query, args...)
if err != nil {
	fmt.Println(err)
		http.Error(w, "Invalid teacher ID", http.StatusInternalServerError)
		return
	}
 defer rows.Close()

teacherList := make([]models.Teacher, 0)

for rows.Next(){
	var teacher  models.Teacher
	err :=  rows.Scan(
		&teacher.ID,
		&teacher.FirstNAme,
		&teacher.LastName,
		&teacher.Email,
		&teacher.Level,
		&teacher.Subject,
	)
	if err != nil {
	fmt.Println(err)
		http.Error(w, "Error scanning databasae", http.StatusInternalServerError)
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
		&teacher.FirstNAme,
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

func addTeacherHandler(w http.ResponseWriter, r *http.Request) {
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
			newTeacher.FirstNAme,
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
