package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"os"
	"strings"
)

// BucketBasics encapsulates the Amazon Simple storage Service (Amazon S3) actions
type BucketBasics struct {
	S3Client *s3.Client
}

type ImageDataS3 struct {
	Base64Data []byte
	Type       string
}

// S3ParamsFromBase64 converts a base64 string with MIME type to ImageDataS3
func S3ParamsFromBase64(base64Str string) (ImageDataS3, error) {
	// Split the base64 string by commas to separate the data prefix from the actual base64 data
	parts := strings.Split(base64Str, ",")
	if len(parts) != 2 {
		return ImageDataS3{
			Base64Data: []byte(""),
			Type:       "",
		}, errors.New("invalid base64 string")
	}

	// Extract the file type from the data prefix
	fileTypeParts := strings.Split(parts[0], ";")
	if len(fileTypeParts) != 2 {
		return ImageDataS3{
			Base64Data: []byte(""),
			Type:       "",
		}, errors.New("invalid base64 string")
	}
	fileType := strings.Split(fileTypeParts[0], "/")[1]

	// Decode the base64 string into bytes
	decodedData, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return ImageDataS3{
			Base64Data: []byte(""),
			Type:       "",
		}, err
	}

	return ImageDataS3{
		Base64Data: decodedData,
		Type:       fileType,
	}, nil
}

// UploadFile reads from a file and puts the data into an object in a bucket.
func (basics BucketBasics) UploadFile(base64 string) (string, error) {
	bucketName := os.Getenv("AWS_S3_BUCKET_NAME")
	bucketUrl := os.Getenv("AWS_S3_BUCKET_URL")
	objectKey := uuid.NewString()

	params, err := S3ParamsFromBase64(base64)
	if err != nil {
		return "", err
	}
	reader := bytes.NewReader(params.Base64Data)
	contentType := fmt.Sprintf("image/%s", params.Type)

	_, err = basics.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		GrantRead:       aws.String("uri=http://acs.amazonaws.com/groups/global/AllUsers"),
		Bucket:          aws.String(bucketName),
		Key:             aws.String(objectKey),
		Body:            reader,
		ContentEncoding: aws.String("base64"),
		ContentType:     aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("Couldn't upload file %v to %v:%v. Here's why: %v\n",
			base64, bucketName, objectKey, err)
	}

	imageUrl := fmt.Sprintf("%s/%s", bucketUrl, objectKey)

	return imageUrl, err
}
