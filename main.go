package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path"
)

// FileSizeResponse defines the response send to the user
type FileSizeResponse struct {
	Size int64 `json:"size"`
}

// Declare any errors
var (
	ErrUnableToProcessForm = errors.New("Unable to process form")
	ErrNoFile              = errors.New("You must upload a file")
)

func main() {
	fmt.Println("hello world")
	port := os.Getenv("PORT")
	r := http.NewServeMux()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		temp, _ := template.ParseFiles(path.Join("index.html"))
		temp.Execute(w, temp)

		return
	})

	r.HandleFunc("/get-file-size", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			RespondErr(w, ErrUnableToProcessForm, http.StatusBadRequest)
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			RespondErr(w, ErrNoFile, http.StatusBadRequest)
			return
		}

		defer file.Close()

		response := FileSizeResponse{Size: header.Size}

		Respond(w, response, 200)
	})

	http.ListenAndServe(":"+port, r)
}

// RespondErr handles all error responses
func RespondErr(w http.ResponseWriter, err error, code int) {
	e := struct {
		Error string `json:"error"`
	}{
		Error: err.Error(),
	}

	Respond(w, e, code)
}

// Respond marshal the given data in to json and sends the response to the client
func Respond(w http.ResponseWriter, data interface{}, code int) {
	js, err := json.Marshal(data)
	if err != nil {
		fmt.Println("error", err)

		js = []byte("{}")
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	io.WriteString(w, string(js))
}
