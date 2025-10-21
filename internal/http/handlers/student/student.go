package student

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/Aytaditya/students-api-golang/internal/types"
	"github.com/Aytaditya/students-api-golang/internal/utils/response"
)

// this returns a hhtp function which is used
func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var student types.Student
		// decode json body to student struct
		err := json.NewDecoder(r.Body).Decode(&student)
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

		response.WriteJson(w, http.StatusCreated, map[string]string{"success": "student created"})
	}
}
