package operations

import (
	"log"
	"os"
	"testing"
)

func TestBlobLifeCycvle(t *testing.T){
	tempDir, err := os.Mkdir("", "blobtest")
	if err != nil{
		t.Fatalf("failed to create temp dir :%v", err);
	}

	defer os.RemoveAll(tempDir);

	store, err := NewStore(tempDir)
	if err != nil {
		t.Fatalf("Failed to init store:%v", err)
	}

	payload := []byte("TEsting pure file-backend Go BLOB storage!")

	log.Print(store, payload)// so it doesn't throw an error as I can't code any longer today
}