package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/mail"
	"net/smtp"
	"strings"
	"time"
)

func logRequestMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		// Log the request information
		log.Printf("Started %s %s", r.Method, r.URL.Path)

		// Call the next handler in the chain
		handler.ServeHTTP(w, r)

		// Log the request duration
		duration := time.Since(startTime)
		log.Printf("Completed %s in %v", r.URL.Path, duration)
	})
}

type jsonResponse struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		email = "Guest@guejkknkst.com"
	}

	emailToVerify := "achuulimagbama@gmail.com"

	var result string
	var status int
	err := verifyEmail(emailToVerify)
	if err != nil {
		status = 400
		w.WriteHeader(http.StatusPreconditionFailed)
		result = fmt.Sprintf("Email verification failed: %v", err)
		fmt.Println("Email verification failed:", err)
	} else {
		status = 200
		w.WriteHeader(http.StatusOK)
		result = fmt.Sprintf("Email verification successful %v", "yeah!")
		fmt.Println("Email verification successful!")
	}
	response := jsonResponse{Message: result, StatusCode: status}

	w.Header().Set("Content-Type", "application/json")
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Write(jsonResponse)
}

func main() {
	// Replace with the email address you want to verify

	http.Handle("/", logRequestMiddleware(http.HandlerFunc(helloHandler)))

	fmt.Println("Server listening on :8080...")

	// Start the server with the loggedHandler
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func verifyEmail(email string) error {
	// Step 1: Validate email syntax
	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("invalid email syntax: %v", err)
	}

	// Step 2: Extract domain from email
	parts := strings.Split(email, "@")
	domain := parts[1]

	// Step 3: Verify domain MX records
	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		return fmt.Errorf("failed to lookup MX records for domain: %v", err)
	}

	// Step 4: Connect to the first MX server and check if the email address exists
	err = verifySMTP(email, mxRecords[0].Host)
	if err != nil {
		return fmt.Errorf("SMTP verification failed: %v", err)
	}

	return nil
}

func verifySMTP(email, mxHost string) error {
	// Step 5: Connect to the SMTP server
	client, err := smtp.Dial(mxHost + ":25")
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %v", err)
	}
	defer client.Close()

	// Step 6: Set sender and recipient
	err = client.Mail("achuulimagbama@gmail.com")
	if err != nil {
		return fmt.Errorf("failed to set sender address: %v", err)
	}

	err = client.Rcpt(email)
	if err != nil {
		return fmt.Errorf("failed to set recipient address: %v", err)
	}

	return nil
}
