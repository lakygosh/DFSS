package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func uploadFile(serverURL, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		return fmt.Errorf("failed to create form file: %v", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return fmt.Errorf("failed to copy file: %v", err)
	}
	writer.Close()

	resp, err := http.Post(serverURL+"/upload", writer.FormDataContentType(), body)
	if err != nil {
		return fmt.Errorf("failed to upload file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to upload file, server responded with status: %s", resp.Status)
	}

	fmt.Println("File uploaded successfully")
	return nil
}

func downloadFile(serverURL, key, destPath string) error {
	resp, err := http.Get(serverURL + "/download?key=" + key)
	if err != nil {
		return fmt.Errorf("failed to download file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file, server responded with status: %s", resp.Status)
	}

	out, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save file: %v", err)
	}

	fmt.Println("File downloaded successfully")
	return nil
}

func main() {
	// Example usage
	serverURL := "http://localhost:3000"
	uploadFile(serverURL, "path/to/your/file.txt")
	downloadFile(serverURL, "file.txt", "path/to/save/file.txt")
}