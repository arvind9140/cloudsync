package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/spf13/cobra"
)

var (
	accessKey  string
	secretKey  string
	bucketName string
	region     string
	localDir   string
)

var syncCmd = &cobra.Command{
	Use:   "sync_remote",
	Short: "Sync files from S3 bucket to local directory",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Enter AWS Access Key ID: ")
		accessKey, _ = reader.ReadString('\n')

		fmt.Print("Enter AWS Secret Access Key: ")
		secretKey, _ = reader.ReadString('\n')

		fmt.Print("Enter S3 Bucket Name: ")
		bucketName, _ = reader.ReadString('\n')

		fmt.Print("Enter AWS Region: ")
		region, _ = reader.ReadString('\n')

		fmt.Print("Enter Local Directory Path: ")
		localDir, _ = reader.ReadString('\n')

		// Remove newline characters
		accessKey = strings.TrimSpace(accessKey)
		secretKey = strings.TrimSpace(secretKey)
		bucketName = strings.TrimSpace(bucketName)
		region = strings.TrimSpace(region)
		localDir = strings.TrimSpace(localDir)
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := syncS3Files(accessKey, secretKey, bucketName, region, localDir)
		if err != nil {
			fmt.Println("Error syncing files:", err)
			return
		}
		fmt.Println("Files synced successfully!")
	},
}

func init() {
  rootCmd.AddCommand(syncCmd)
}

func syncS3Files(accessKey, secretKey, bucketName, region, localDir string) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
		Credentials: credentials.NewStaticCredentials(
			accessKey,
			secretKey,
			"", // a session token is not needed
		),
	})
	if err != nil {
		return err
	}

	svc := s3.New(sess)

	params := &s3.ListObjectsInput{
		Bucket: aws.String(bucketName),
	}

	resp, err := svc.ListObjects(params)
	if err != nil {
		return err
	}

	for _, item := range resp.Contents {
		key := aws.StringValue(item.Key)
		downloadParams := &s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
		}

		file, err := os.Create(localDir + "/" + key)
		if err != nil {
			return err
		}

		defer file.Close()

		resp, err := svc.GetObject(downloadParams)
		if err != nil {
			return err
		}

		_, err = io.Copy(file, resp.Body)
		if err != nil {
			return err
		}
	}

	return nil
}
