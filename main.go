package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

func uploadFile(w http.ResponseWriter, r *http.Request) {
	// 10 << 20 specifies a maximum upload of 10 MB files.
	var MAX_SIZE int64 = 10 << 20

	fmt.Println("File Upload Endpoint Hit")
	if r.ContentLength > MAX_SIZE {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	l, err := strconv.ParseInt(r.Header.Get("Content-Length"), 10, 64)
	if err != nil || l > MAX_SIZE {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Parse our multipart form
	r.ParseMultipartForm(MAX_SIZE)
	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern
	tempFile, err := ioutil.TempFile("temp-images", "upload-*.png")
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)
	// return that we have successfully uploaded our file!
	fmt.Fprintf(w, "Successfully Uploaded File\n")
}

func setupRoutes() {
	http.HandleFunc("/upload", uploadFile)
	http.Handle("/", http.FileServer(http.Dir(".")))
	http.ListenAndServe(":8080", nil)
}

func main() {
	fmt.Println("Hello World")
	setupRoutes()
}
