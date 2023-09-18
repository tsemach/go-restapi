package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func main() {
	// Load AWS configuration from the default config file or environment variables
	// cfg, err := config.LoadDefaultConfig(context.TODO())
	// if err != nil {
	// 	fmt.Println("Error loading AWS configuration:", err)
	// 	return
	// }

	// // Create an S3 client
	// client := s3.NewFromConfig(cfg)

	const defaultRegion = "us-east-1"
	hostAddress := "http://localhost:9000"

	resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...any) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:       "aws",
			SigningRegion:     defaultRegion,
			URL:               hostAddress,
			HostnameImmutable: true,
		}, nil
	})

	cfg := aws.Config{
		Region:                      defaultRegion,
		EndpointResolverWithOptions: resolver,
		Credentials:                 credentials.NewStaticCredentialsProvider("minioadmin", "minioadmin", ""),
	}

	client := s3.NewFromConfig(cfg)

	// Create an HTTP server to handle the incoming request
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		// Replace with your S3 bucket name and object key
		bucketName := "upload-test"
		objectKey := "file.dat"

		// Create a new multipart upload
		createOutput, err := client.CreateMultipartUpload(context.TODO(), &s3.CreateMultipartUploadInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectKey),
		})
		if err != nil {
			fmt.Println("Error creating multipart upload:", err)
			http.Error(w, "Failed to create multipart upload", http.StatusInternalServerError)
			return
		}

		uploadID := createOutput.UploadId

		// Read the request body (file)
		file, _, err := r.FormFile("file")
		if err != nil {
			fmt.Println("Error reading request body:", err)
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		defer file.Close()

		partNumber := int32(1)
		partSize := int64(5 * 1024 * 1024) // 5 MB (adjust as needed)

		// Initialize variables to keep track of uploaded parts
		var completedParts []types.CompletedPart

		// Upload parts of the file
		for {
			partBuffer := make([]byte, partSize)
			n, err := file.Read(partBuffer)
			if err != nil && err != io.EOF {
				fmt.Println("Error reading part of the file:", err)
				http.Error(w, "Failed to read part of the file", http.StatusInternalServerError)
				return
			}

			if n == 0 {
				break // No more data to read
			}

			// Upload the part to S3
			uploadPartOutput, err := client.UploadPart(context.TODO(), &s3.UploadPartInput{
				Bucket:     aws.String(bucketName),
				Key:        aws.String(objectKey),
				UploadId:   uploadID,
				PartNumber: partNumber,
				Body:       bytes.NewReader(partBuffer[:n]),
			})
			if err != nil {
				fmt.Println("Error uploading part to S3:", err)
				http.Error(w, "Failed to upload part to S3", http.StatusInternalServerError)
				return
			}

			completedParts = append(completedParts, types.CompletedPart{
				ETag:       uploadPartOutput.ETag,
				PartNumber: partNumber,
			})

			partNumber++
		}

		// Complete the multipart upload
		_, err = client.CompleteMultipartUpload(context.TODO(), &s3.CompleteMultipartUploadInput{
			Bucket:   aws.String(bucketName),
			Key:      aws.String(objectKey),
			UploadId: uploadID,
			MultipartUpload: &types.CompletedMultipartUpload{
				Parts: completedParts,
			},
		})
		if err != nil {
			fmt.Println("Error completing multipart upload:", err)
			http.Error(w, "Failed to complete multipart upload", http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "File uploaded successfully to S3")
	})

	// Start the HTTP server
	port := ":8080"
	fmt.Printf("Server started on port %s\n", port)
	http.ListenAndServe(port, nil)
}
