// tools/generate_postman_collection.go
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type PostmanCollection struct {
	Info struct {
		PostmanID string `json:"_postman_id"`
		Name      string `json:"name"`
		Schema    string `json:"schema"`
	} `json:"info"`
	Item []PostmanItem `json:"item"`
}

type PostmanItem struct {
	Name string `json:"name"`
	Item []struct {
		Name    string        `json:"name"`
		Request PostmanRequest `json:"request"`
	} `json:"item"`
}

type PostmanRequest struct {
	Method string `json:"method"`
	Header []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"header"`
	Body struct {
		Mode string `json:"mode"`
		Raw  string `json:"raw"`
	} `json:"body"`
	URL struct {
		Raw  string   `json:"raw"`
		Host []string `json:"host"`
		Path []string `json:"path"`
	} `json:"url"`
}

func generatePostmanCollection() error {
	collection := PostmanCollection{
		// Initialize collection structure
	}

	// Populate collection based on your API routes

	// Write to file
	outputPath := filepath.Join("postman", "greeneye_collection.json")
	os.MkdirAll(filepath.Dir(outputPath), 0755)

	file, err := json.MarshalIndent(collection, "", " ")
	if err != nil {
		return err
	}

	err = os.WriteFile(outputPath, file, 0644)
	if err != nil {
		return err
	}

	fmt.Printf("Postman collection generated: %s\n", outputPath)
	return nil
}

func main() {
	if err := generatePostmanCollection(); err != nil {
		fmt.Println("Error generating Postman collection:", err)
		os.Exit(1)
	}
}