package operations

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Store struct{
	BaseDir string;
}
func NewStore(baseDir string) (*Store, error){
	if err := os.MkdirAll(baseDir, 0755); err != nil{
		return nil, fmt.Errorf("failed to create base storage directory: %w", err);
	}
	return &Store{BaseDir: baseDir}, nil;
}

func(s *Store) ResolvePath(hash string) (string, error) {
	if len(hash) < 4{
		return "", fmt.Errorf("invalid hash length")
	}
	dir := filepath.Join(s.BaseDir, hash[0:2], hash[2:4]);
	return filepath.Join(dir, hash), nil;
}

func (s *Store) HandleUpload(writer http.ResponseWriter, request *http.Request){
	if request.Method != http.MethodPut && request.Method != http.MethodPost{
		http.Error(writer, "Method not allowed. Use PUT or POST apis.", http.StatusMethodNotAllowed);
		return;
	}

	tempFile, err := os.CreateTemp(s.BaseDir, "upload-*")
	if err != nil{
		http.Error(writer, "Failed to create temp file.", http.StatusInternalServerError);
		return;
	}

	defer os.Remove(tempFile.Name());
	defer tempFile.Close();

	hasher := sha256.New();
	multiWriter := io.MultiWriter(tempFile, hasher);

	if _, err := io.Copy(multiWriter, request.Body); err != nil {
		http.Error(writer, "Error reading request payload", http.StatusInternalServerError);
		return;
	}
	
	hashString := hex.EncodeToString(hasher.Sum(nil));

	targetPath, err := s.ResolvePath(hashString);
	if err != nil{
		http.Error(writer, "Failed to create storage key", http.StatusInternalServerError);
		return;
	}

	errMk := os.MkdirAll(filepath.Dir(targetPath), 0755);
	if errMk != nil{
		http.Error(writer, "Failed to reserve storage path", http.StatusInternalServerError);
		return;
	}

	writer.WriteHeader(http.StatusCreated);
	fmt.Fprintf(writer, "%s\n", hashString);
}

func (s *Store) HandleDownloads(writer http.ResponseWriter, request *http.Request){
	if request.Method != http.MethodGet {
		http.Error(writer, "Method not allowed! Use Get method.", http.StatusMethodNotAllowed);
		return;
	}

	hash := strings.TrimPrefix(request.URL.Path, "/blobs/");
	if hash == ""{
		http.Error(writer, "Hash key is required", http.StatusBadRequest);
		return;
	}

	filePath, err := s.ResolvePath(filepath.Clean(hash))
	if err != nil{ 
		http.Error(writer, "Invalid key format", http.StatusBadRequest);
		return;
	}

	file, err := os.Open(filePath)
	if os.IsNotExist(err) {
		http.Error(writer, "Blob not found this time", http.StatusNotFound);
		return; 
	}else if err != nil {
		http.Error(writer, "Failed to acess blob", http.StatusInternalServerError);
		return;
	}

	defer file.Close();

	io.Copy(writer, file);

}