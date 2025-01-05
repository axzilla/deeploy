package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/axzilla/deeploy/internal/app/auth"
	"github.com/axzilla/deeploy/internal/app/db"
)

func main() {
	email := flag.String("email", "", "User email")
	password := flag.String("password", "", "New password")
	flag.Parse()

	// Validate flags
	if *email == "" || *password == "" {
		fmt.Println("Usage: reset-password -email=user@example.com -password=newpassword")
		os.Exit(1)
	}

	// DB connection
	db, err := db.Init()
	if err != nil {
		fmt.Printf("Initializind DB failed: %v", err)
		os.Exit(1)
	}

	// Hash new password
	hashedPassword, err := auth.HashPassword(*password)
	if err != nil {
		fmt.Printf("Hashing password failed: %v", err)
		os.Exit(1)
	}

	// Update password
	query := `UPDATE users SET password = ? WHERE email = ?`
	result, err := db.Exec(query, hashedPassword, *email)
	if err != nil {
		fmt.Printf("Update password failed: %v", err)
		os.Exit(1)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("Update password failed: %v", err)
	}
	if rowsAffected == 0 {
		fmt.Print("No user with this email found")
		os.Exit(1)
	}

	// Success message
	fmt.Println("Password successfully reset!")
}
