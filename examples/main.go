package main

import (
	"context"
	"fmt"
	"log"
	"os"

	s21client "github.com/arseniisemenow/s21auto-client-go"
	"github.com/arseniisemenow/s21auto-client-go/auth"
	"github.com/arseniisemenow/s21auto-client-go/requests"
)

func main() {
	// TODO: Add your credentials as environment variables
	username := os.Getenv("S21_USERNAME")
	password := os.Getenv("S21_PASSWORD")

	if username == "" || password == "" {
		log.Fatal("Please set S21_USERNAME and S21_PASSWORD environment variables")
	}

	ctx := context.Background()

	// Test auth first
	fmt.Println("Testing auth...")
	token, err := auth.RequestToken(username, password, ctx)
	if err != nil {
		log.Fatalf("Failed to get token: %v", err)
	}
	fmt.Printf("Got access token (first 20 chars): %s...\n", token.AccessToken[:20])

	// Test getting user data for school ID
	fmt.Println("\nFetching user data for school ID...")
	user, err := auth.RequestUserData(token, ctx)
	if err != nil {
		log.Fatalf("Failed to get user data: %v", err)
	}
	fmt.Printf("User ID: %s, Roles count: %d\n", user.Data.ID, len(user.Roles))

	// Test getting context headers (new feature)
	fmt.Println("\nFetching edu-context headers...")
	contextHeaders, err := auth.RequestContextHeaders(token, ctx)
	if err != nil {
		log.Fatalf("Failed to get context headers: %v", err)
	}
	fmt.Printf("X-EDU-SCHOOL-ID: %s\n", contextHeaders.XEDUSchoolID)
	fmt.Printf("X-EDU-PRODUCT-ID: %s\n", contextHeaders.XEDUProductID)
	fmt.Printf("X-EDU-ORG-UNIT-ID: %s\n", contextHeaders.XEDUOrgUnitID)
	fmt.Printf("X-EDU-ROUTE-INFO: %s\n", contextHeaders.XEDURouteInfo)

	// Now test the full client flow
	// The client automatically fetches and applies context headers to all requests
	fmt.Println("\nTesting full client flow (headers applied automatically)...")
	client := s21client.New(
		s21client.DefaultAuth(username, password),
	)

	currentUser, err := client.R().GetCurrentUser(requests.GetCurrentUser_Variables{})
	if err != nil {
		log.Fatalf("Failed to get current user: %v", err)
	}

	fmt.Printf("Current user: %s %s\n", currentUser.User.GetCurrentUser.FirstName, currentUser.User.GetCurrentUser.LastName)
	fmt.Println("Success! Context headers are now automatically added to all requests.")
}
