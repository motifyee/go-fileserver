package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"

	"crypto/rand"
	// "github.com/nu7hatch/gouuid"
)

func main() {
	http.HandleFunc("/files/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "files/"+r.URL.Path[len("/files/"):])
	})
	http.HandleFunc("/files/list", func(w http.ResponseWriter, r *http.Request) {
		entries, err := os.ReadDir("./files")
		if err != nil {
			jsonResp(w, "error", "Error listing files: "+err.Error())

			return
		}

		var files []string
		for _, entry := range entries {
			files = append(files, entry.Name())
		}

		jsonResp(w, "success", files)
	})

	http.HandleFunc("/upload", uploadFile)

	http.HandleFunc("/inc", incCounter)

	log.Fatal(http.ListenAndServe(":8989", nil))

}

func uploadFile(w http.ResponseWriter, r *http.Request) {

	fmt.Println("uploadFile")
	r.ParseMultipartForm(10 << 20) // 10 MB max upload files

	file, handler, err := r.FormFile("file")
	if err != nil {
		jsonResp(w, "error", "Error Retrieving the File\n"+err.Error())

		return
	}

	defer file.Close()

	// fmt.Fprintf(w, "%v", handler.Header)
	// fmt.Fprintf(w, "%v", handler.Size)
	// fmt.Fprintf(w, "%v", handler.Filename)

	uuid, err := createUUID()
	if err != nil {
		jsonResp(w, "error", "Error creating UUID\n"+err.Error())

		return
	}

	newFileName := uuid + "__" + handler.Filename

	tempFile, err := os.Create("./files/" + newFileName)
	if err != nil {
		jsonResp(w, "error", "Error creating File\n"+err.Error())

		return
	}

	defer tempFile.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		jsonResp(w, "error", "Error reading File\n"+err.Error())

		return
	}

	tempFile.Write(fileBytes)

	jsonResp(w, "success", newFileName)
}

func jsonResp(w http.ResponseWriter, status string, message interface{}) {
	resp := make(map[string]interface{})

	resp["status"] = status
	resp["message"] = message

	jsonR, err := json.Marshal(resp)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "{ \"status\": \"error\", \"message\": \"Error happened in JSON marshal. Err: %s\" }", err)

		return
	}

	// fmt.Println(w, string(jsonResp))
	w.WriteHeader(http.StatusOK)
	w.Write(jsonR)
}

func createUUID() (string, error) {
	b := make([]byte, 16)

	s, err := rand.Read(b)
	if err != nil || s != len(b) {
		fmt.Println("createUUID: error reading random bytes:", err)
		return "", err
	}

	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid, nil
}

var counter int
var mutex = &sync.Mutex{}

func incCounter(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	{
		counter++
		fmt.Fprintf(w, "Counter = %d\n", counter)
	}
	mutex.Unlock()
}
