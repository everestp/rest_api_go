package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/everestp/rest_api_go/internal/api/models"
)


var (

	 teachers = make(map[int]models.Teacher)
	  mutex = &sync.Mutex{}
	 nextID = 1
)
  // Var intitlized the Dummy data

  func init(){
	teachers[nextID]= models.Teacher{
		ID: nextID,
		FirstNAme: "jhon",
		LastName: "Doe",
		Class: "9A",
		Subject: "Math",
	}
	nextID++
	teachers[nextID]= models.Teacher{
		ID: nextID,
		FirstNAme: "Everest",
		LastName: "Paudel",
		Class: "Bsc",
		Subject: "Math",
	}
	nextID++

  }
func TeacherHandler(w http.ResponseWriter, r *http.Request) {
	//path := strings.TrimPrefix(r.URL.Path, "/teacher/")
	// userID := strings.TrimSuffix(path, "/")
	switch r.Method{
	case http.MethodGet:
		//cal get method handler fucntion
		getTeachersHandler(w, r)

		case http.MethodPost:
			postTeacherHandler(w, r)
		w.Write([]byte("General Teacher Directory"))
		case http.MethodPut:
		w.Write([]byte("General Teacher Directory"))
		case http.MethodDelete:
		w.Write([]byte("General Teacher Directory"))

	}

	
}

  func getTeachersHandler(w http.ResponseWriter , r *http.Request){
  path := strings.TrimPrefix(r.URL.Path, "/teacher/")
  idStr := strings.TrimSuffix(path, "")
  fmt.Println(idStr)

  


	if idStr =="" {firstName := r.URL.Query().Get("first_name")
	LastName := r.URL.Query().Get("last_name")
	teacherList := make([]models.Teacher,0,len(teachers))
	for _ ,teacher := range teachers{
		if(firstName=="" || teacher.FirstNAme ==firstName) &&(LastName =="" || teacher.LastName == LastName){

			teacherList = append(teacherList, teacher)
		}
	}
	response := struct{
		Status string `json:"status"`
		Count int  `json:"status"`
		Data []models.Teacher `json:"data"`
	}{Status: "sucess",
	  Count: len(teachers),
	  Data: teacherList,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

//Handle Path Parametre
id , err := strconv.Atoi(idStr)
if err != nil{
	fmt.Println(err)
	return
}

teacher ,exist := teachers[id]
if !exist {
	http.Error(w, "Teacher not found", http.StatusNoContent)
	return
}
json.NewEncoder(w).Encode(teacher)


  }

  // Post Teacher Handler
func postTeacherHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()
	defer r.Body.Close()

	var newTeachers []models.Teacher
	err := json.NewDecoder(r.Body).Decode(&newTeachers)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	for  _ , request := range newTeachers{
 teachers[len(teachers) +1] = models.Teacher{
		ID: len(teachers) +1,
		FirstNAme: request.FirstNAme,
		LastName: request.FirstNAme,
		Class: request.Class,
		Subject: request.Subject,

	  }
	}
w.Header().Set("Content-Type", "application/json")

	 
}