package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func UploadFileFromFormData(field_name string, r *http.Request) ([]byte, error) {
	file, handler, err := r.FormFile(field_name)
	if err != nil {
		fmt.Println("Error Retrieving the File")
		return nil, err
	} else {
		defer file.Close()
		fmt.Printf("Uploaded File: %+v\n", handler.Filename)
		fmt.Printf("File Size: %+v\n", handler.Size)
		fmt.Printf("MIME Header: %+v\n", handler.Header)

		// Read all of the contents of our uploaded file into a byte array
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		return fileBytes, nil
	}

}
