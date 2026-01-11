package handlers

import "net/http"

func ExcesHandler(w http.ResponseWriter, r *http.Request) {
	//path := strings.TrimPrefix(r.URL.Path, "/teacher/")
	// userID := strings.TrimSuffix(path, "/")
	switch r.Method{
	case http.MethodGet:
		//cal get method handler fucntion
		// getTeachersHandler(w, r)

		case http.MethodPost:
			 // postTeacherHandler(w, r)
		w.Write([]byte("General Teacher Directory"))
		case http.MethodPut:
		w.Write([]byte("General Teacher Directory"))
		case http.MethodDelete:
		w.Write([]byte("General Teacher Directory"))

	}

	
}