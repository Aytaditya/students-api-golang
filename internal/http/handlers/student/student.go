package student

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Aytaditya/students-api-golang/internal/storage"
	"github.com/Aytaditya/students-api-golang/internal/types"
	"github.com/Aytaditya/students-api-golang/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

// this returns a hhtp function which is used  (method to create a new student)
// our sb instance has name, email, age fields hence we can accept it here to function in sotrage.storage
func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var student types.Student
		// decode json body to student struct
		err := json.NewDecoder(r.Body).Decode(&student) // The JSON data from the body is read (consumed) from r.Body and then decoded directly into your Go struct (student).

		// checking if error is empty body
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// if there is some other error (apart from empty body)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		// validating the request
		er := validator.New().Struct(student)
		if er != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(er)) // http.StatusBadRequest=400
			return
		}

		id, er := storage.CreateStudent(student.Name, student.Email, student.Age) // implementing the CreateStudent method of storage interface

		slog.Info("Student created successfully", "id", id, "name", student.Name, "email", student.Email, "age", student.Age)

		if er != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(er))
			return
		}

		response.WriteJson(w, http.StatusCreated, map[string]int64{"id": id})
	}
}

func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// to be implemented
		id := r.PathValue("id")
		slog.Info("Get student by id called", "id", id)
		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		name, email, age, err := storage.GetStudent(intId)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		student := types.Student{
			Id:    intId,
			Name:  name,
			Email: email,
			Age:   age,
		}
		response.WriteJson(w, http.StatusOK, student)
	}
}
