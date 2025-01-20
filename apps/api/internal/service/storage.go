package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type StorageService struct {
	bucketName string
}

func NewStorageService(bucketName string) *StorageService {
	return &StorageService{bucketName: bucketName}
}

func (s *StorageService) UploadMarkdownHandler(w http.ResponseWriter, r *http.Request, organization_id string, folderPath string) error {
	if err := r.ParseMultipartForm(10 << 20); err != nil { // Limit upload size to 10MB
		return fmt.Errorf("Error parsing form data")
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		return fmt.Errorf("Error retrieving file")
	}
	defer file.Close()

	if header.Header.Get("Content-Type") != "text/markdown" && header.Filename[len(header.Filename)-3:] != ".md" {
		return fmt.Errorf("Only .md files are allowed")
	}

	path := organization_id + folderPath
	ctx := context.Background()
	if err := s.UploadFileToBucket(ctx, file, header.Filename, path); err != nil {
		return fmt.Errorf("Error uploading file: %v", err)

	}

	return nil
}

func (s *StorageService) UploadFileToBucket(ctx context.Context, file io.Reader, filename string, folderName string) error {
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	fileReader := bytes.NewReader(fileBytes)

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to load AWS configuration: %w", err)
	}

	client := s3.NewFromConfig(cfg)

	objectPath := filename
	if folderName != "" {
		objectPath = fmt.Sprintf("%s/%s", folderName, filename)
	}

	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(objectPath),
		Body:   fileReader,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file to S3: %w", err)
	}

	fmt.Printf("File %s uploaded to bucket %s in folder %s\n", filename, s.bucketName, folderName)
	return nil
}
