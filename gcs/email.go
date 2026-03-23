package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Pepipost API structures
type From struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type Content struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Attachment struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Content string `json:"content"`
}

type ToRecipient struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type Personalization struct {
	To []ToRecipient `json:"to"`
}

type PepipostEmailRequest struct {
	From             From              `json:"from"`
	Subject          string            `json:"subject"`
	Content          []Content         `json:"content"`
	Attachments      []Attachment      `json:"attachments"`
	Personalizations []Personalization `json:"personalizations"`
}

// Function to create Pepipost email request with GCS file attachment
func CreatePepipostEmailRequest(fileName string, fileData []byte, mimeType string) *PepipostEmailRequest {
	// Encode the file data to base64
	encodedData := base64.StdEncoding.EncodeToString(fileData)

	return &PepipostEmailRequest{
		From: From{
			Email: "info@crm.purplleletters.com",
			Name:  "Purplle_Letters_CRM",
		},
		Subject: "Test email with PDF from Pepipost",
		Content: []Content{
			{
				Type:  "html",
				Value: "Hello, Please find the attached file.",
			},
		},
		Attachments: []Attachment{
			{
				Name:    fileName,
				Type:    mimeType,
				Content: encodedData,
			},
		},
		Personalizations: []Personalization{
			{
				To: []ToRecipient{
					{
						Email: "sanil.s@purplle.com",
						Name:  "Recipient Name",
					},
				},
			},
		},
	}
}

// Function to send email via Pepipost API
func SendPepipostEmail(apiKey string, emailRequest *PepipostEmailRequest) error {
	// Convert request to JSON
	jsonData, err := json.Marshal(emailRequest)
	if err != nil {
		return fmt.Errorf("failed to marshal email request: %v", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", "https://api.pepipost.com/v5/mail/send", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api_key", apiKey)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("pepipost API returned status: %d", resp.StatusCode)
	}

	log.Printf("Email sent successfully via Pepipost. Status: %d", resp.StatusCode)
	return nil
}

