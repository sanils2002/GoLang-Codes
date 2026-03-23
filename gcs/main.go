package main

import (
	"context"
	"io"
	"log"

	"cloud.google.com/go/storage"
)

func ReadGCSFile() []byte {
	ctx := context.Background()

	// Create a GCS client
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Replace with your bucket and object names
	bucketName := "martech-playground"
	objectName := "object_1/Email Attachments POC.pdf"

	// Create a new reader for the object
	rc, err := client.Bucket(bucketName).Object(objectName).NewReader(ctx)
	if err != nil {
		log.Fatalf("Failed to create object reader: %v", err)
	}
	defer rc.Close() // Ensure the reader is closed after use

	// Read the content of the file
	data, err := io.ReadAll(rc)
	if err != nil {
		log.Fatalf("Failed to read object data: %v", err)
	}

	log.Printf("Successfully read file %s/%s, size: %d bytes", bucketName, objectName, len(data))

	return data
}

func main() {
	// Get Pepipost API key from environment variable
	apiKey := "5542279cc9d93abd5912f3592c66c0ae"
	if apiKey == "" {
		log.Fatal("PEPIPOST_API_KEY environment variable is required")
	}

	// Read file from GCS
	log.Println("Reading file from Google Cloud Storage...")
	fileData := ReadGCSFile()

	// Create Pepipost email request
	log.Println("Creating Pepipost email request...")
	emailRequest := CreatePepipostEmailRequest(
		"Email Attachments POC.pdf", // File name
		fileData,                    // File data
		"application/pdf",           // MIME type
	)

	// Send email via Pepipost
	log.Println("Sending email via Pepipost...")
	err := SendPepipostEmail(apiKey, emailRequest)
	if err != nil {
		log.Fatalf("Failed to send email: %v", err)
	}

	log.Println("Email sent successfully!")
}
