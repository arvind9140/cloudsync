package main

import (
	"fmt"

	"os"

	"github.com/arvind9140/cloudsync/cmd"
	"github.com/joho/godotenv"
	
)

func main() {
	
	cmd.Execute()
	 if err := godotenv.Load(); err != nil {
        fmt.Println("Error loading .env file:", err)
        os.Exit(1)
    }

	
   
}
