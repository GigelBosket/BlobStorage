package blobstorage

import (
	"blobStorage/operations"
	"fmt"
	"log"
	"net/http"
)

const port = ":8080";
const storagePath = "./storage/blobs";

func main(){
	store, err := operations.NewStore(storagePath);
	if err != nil{
		log.Fatalf("Initilixation error: %v", err);
	}

	http.HandleFunc("/upload", store.HandleUpload);
	http.HandleFunc("/blobs/", store.HandleDownloads);

	fmt.Printf("BLOB storage server is running on http://localhost&s\n", port);
	fmt.Println("- ZUpload endpoint: POST/PUT /upload");
	fmt.Println(" - Download endpoint: GET /blobs/{hash}");

	errFetch := http.ListenAndServe(port, nil);
	if errFetch != nil{
		log.Fatalf("Server failed to start: %v", err);
	}
}