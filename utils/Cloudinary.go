package utils

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
)

// UploadPublicFileToCloudinary Func for uploading public file to cloudinary
func UploadPublicFileToCloudinary(fileBytes []byte) (string, string, error) {
	// Initialize Cloudinary configuration
	cloudinaryURL := os.Getenv("CLOUDINARY_URL") // Get Cloudinary URL from environment variable
	if cloudinaryURL == "" {
		return "", "", fmt.Errorf("cloudinary_url environment variable is not set")
	}

	cld, err := cloudinary.NewFromURL(cloudinaryURL)
	if err != nil {
		return "", "", fmt.Errorf("error creating Cloudinary client: %v", err)
	}

	// Upload file to Cloudinary
	uploadResult, err := cld.Upload.Upload(context.Background(), bytes.NewReader(fileBytes), uploader.UploadParams{})
	if err != nil {
		return "", "", fmt.Errorf("error uploading file to Cloudinary: %v", err)
	}

	// Print the URL of the uploaded file
	fmt.Println("Uploaded URL:", uploadResult.SecureURL)
	fmt.Println("Uploaded Public ID:", uploadResult.PublicID)
	return uploadResult.SecureURL, uploadResult.PublicID, nil
}

// DeleteFileFromCloudinary Func for deleting file from cloudinary
func DeleteFileFromCloudinary(publicID string) (string, error) {
	// Initialize Cloudinary configuration
	cloudinaryURL := os.Getenv("CLOUDINARY_URL") // Get Cloudinary URL from environment variable
	if cloudinaryURL == "" {
		return "", fmt.Errorf("cloudinary_url environment variable is not set")
	}

	cld, err := cloudinary.NewFromURL(cloudinaryURL)
	if err != nil {
		return "", fmt.Errorf("error creating Cloudinary client: %v", err)
	}

	// Delete file from Cloudinary
	_, err = cld.Upload.Destroy(context.Background(), uploader.DestroyParams{PublicID: publicID})
	if err != nil {
		return "", fmt.Errorf("error deleting file from Cloudinary: %v", err)
	}

	result := "file deleted successfully from Cloudinary"
	return result, nil
}
