package main

import (
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

func CheckAvatar(file multipart.File) (string, error) {
	fileHeader := make([]byte, 1024*1024*10)
	ContentType := ""
	if _, err := file.Read(fileHeader); err != nil {
		return ContentType, err
	}

	if _, err := file.Seek(0, 0); err != nil {
		return ContentType, err
	}

	count, err := file.Seek(0, 2)
	if err != nil {
		return ContentType, err
	}
	if count > 1024*1024*10 {
		return ContentType, err
	}
	if _, err := file.Seek(0, 0); err != nil {
		return ContentType, err
	}
	ContentType = http.DetectContentType(fileHeader)

	if ContentType != "image/jpg" && ContentType != "image/png" && ContentType != "image/jpeg" {
		return ContentType, err
	}

	return ContentType, nil
}

func main() {
	log.Println("Starting images microservice")

	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(1024 * 1024 * 10); err != nil {
			log.Printf("ERROR: %v", err)
			http.Error(w, err.Error(), 500)
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, 1024*1024*10)
		defer r.Body.Close()

		file, _, err := r.FormFile("avatar")
		if err != nil {
			log.Printf("ERROR: %v", err)
			http.Error(w, err.Error(), 500)
			return
		}
		// log.Printf("HEADER: %-v", header)

		fileType, err := CheckAvatar(file)
		if err != nil {
			log.Printf("ERROR: %v", err)
			http.Error(w, err.Error(), 500)
			return
		}
		log.Printf("fileType: %-v", fileType)

		f, err := os.Create("image.png")
		if err != nil {
			log.Printf("ERROR: %v", err)
			http.Error(w, err.Error(), 500)
			return
		}

		defer f.Close()

		written, err := io.Copy(f, file)
		if err != nil {
			log.Printf("ERROR: %v", err)
			http.Error(w, err.Error(), 500)
			return
		}

		log.Printf("written: %-v", written)

		w.WriteHeader(200)
		w.Write([]byte(fileType))

	})
	http.ListenAndServe(":5007", nil)
}
