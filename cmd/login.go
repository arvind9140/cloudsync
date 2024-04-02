package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

// UserData represents the structure of data to be stored
type UserData struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	AccessToken string `json:"access_token"`
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login with GitHub",
	Run: func(cmd *cobra.Command, args []string) {
		config := oauth2.Config{
			ClientID:     viper.GetString("github.client_id"),
			ClientSecret: viper.GetString("github.client_secret"),
			Scopes:       []string{"user:email"},
			Endpoint:     github.Endpoint,
		}

		authURL := config.AuthCodeURL("state")
		fmt.Println("Visit the following URL to authorize access:")
		color.Blue(authURL) // Print URL in blue

		var code string
		shutdown := make(chan struct{})

		// Start an HTTP server to receive the authorization code
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			code = r.URL.Query().Get("code")

			fmt.Printf("Enter the authorization code: %s \n", code)

			// Now you can use the code for further processing
			// For example, exchange it for an access token
			token, err := config.Exchange(oauth2.NoContext, code)
			if err != nil {
				fmt.Println("Failed to exchange token:", err)
				return
			}

			fmt.Println("Access token:", token.AccessToken)

			// Get user information
			user, err := getUserInfo(token.AccessToken)
			if err != nil {
				fmt.Println("Failed to get user info:", err)
				return
			}

			// Print "Login successfully" in green
			color.Green("Login successfully")

			// Store the user data in JSON format
			userData := UserData{
				Username:    user.Username,
				Email:       user.Email,
				AccessToken: token.AccessToken,
			}

			filePath := "userdata.json"

			// Write the user data to the file in JSON format
			err = writeUserDataToFile(userData, filePath)
			if err != nil {
				fmt.Println("Error writing user data to file:", err)
				return
			}

			// Shutdown the HTTP server
			close(shutdown)
		})

		// Start the HTTP server
		port := ":8000"
		server := &http.Server{Addr: port}
		go func() {
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("Failed to start server: %v", err)
			}
		}()

		// Wait for the shutdown signal
		<-shutdown

		if err := server.Shutdown(context.Background()); err != nil {
			log.Fatalf("Failed to shutdown server: %v", err)
		}

		// Print all available commands
		color.Cyan("Available commands:")
		err := rootCmd.Usage()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

// Function to write user data to a file in JSON format
func writeUserDataToFile(userData UserData, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Encode user data to JSON
	encoder := json.NewEncoder(file)
	err = encoder.Encode(userData)
	if err != nil {
		return err
	}

	// Print a success message
	fmt.Println("User data written to file:", filePath)
	return nil
}

// Function to get user information using the access token
func getUserInfo(accessToken string) (UserData, error) {
	// Example implementation to get user information from GitHub
	// You need to implement this function based on your authentication provider

	// Here, we assume that you have a function called "getGitHubUserInfo" which gets user information from GitHub
	username, email, err := getGitHubUserInfo(accessToken)
	if err != nil {
		return UserData{}, err
	}
	fmt.Println(username)

	return UserData{
		Username:    username,
		Email:       email,
		AccessToken: accessToken,
	}, nil
}

// Example implementation of getting user information from GitHub
// Replace this function with your actual implementation
func getGitHubUserInfo(accessToken string) (string, string, error) {
	// Create an HTTP client
	client := &http.Client{}

	// Create a request to fetch user information from GitHub
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return "", "", err
	}

	// Set the Authorization header with the access token
	req.Header.Set("Authorization", "token "+accessToken)

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	// Decode the response JSON
	var user UserData
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return "", "", err
	}
	fmt.Println(req)

	// Return the username and email
	return user.Username, user.Email, nil
}
func init() {
	rootCmd.AddCommand(loginCmd)

	// Load configuration
	viper.SetConfigFile("config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("Error reading config file:", err)
		os.Exit(1)
	}
}
