package usecase

import (
	"aws-s3-bucket/domain/upload/interfaces"
	"aws-s3-bucket/models/document"
	"aws-s3-bucket/shared/utils"
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type usecase struct {
	s3Client interfaces.S3Interface
}

func NewUsecase(s3Client interfaces.S3Interface) interfaces.UsecaseInterface {
	return &usecase{s3Client: s3Client}
}

func (u *usecase) UploadBase64(ctx context.Context, request document.RequestUploadDocumentBase64) (response document.ResponseUploadDocument, err error) {

	_, contentType, formatType, decodedBytes, err := utils.ExtractBase64(request.DocumentBase64)
	if err != nil {
		return
	}

	bucketName := os.Getenv("BUCKET_NAME")
	key := fmt.Sprintf("%s/%s.%s", request.DocumentKey, request.DocumentName, formatType)

	_, err = u.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(key),
		Body:        strings.NewReader(string(decodedBytes)),
		ContentType: aws.String(contentType),
		// ACL:         "public-read", //if wanna public use public read
	})
	if err != nil {
		err = fmt.Errorf("failed to upload file: %w", err)
		return
	}

	return document.ResponseUploadDocument{DocumentUrl: fmt.Sprintf("%s/api/v1/download/%s", os.Getenv("BASE_URL"), key)}, nil
}

func (u *usecase) UploadFile(ctx context.Context, request document.RequestUploadDocumentFile, files *multipart.FileHeader) (response document.ResponseUploadDocument, err error) {

	filed, err := files.Open()

	bucketName := os.Getenv("BUCKET_NAME")

	key := fmt.Sprintf("%s/%s%s", request.DocumentKey, request.DocumentName, filepath.Ext(files.Filename))
	_, err = u.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(key),
		Body:        filed,
		ContentType: aws.String(files.Header.Get("Content-Type")),
		// ACL:         "public-read", //if wanna public use public read
	})
	if err != nil {
		err = fmt.Errorf("failed to upload file: %w", err)
		return
	}

	return document.ResponseUploadDocument{DocumentUrl: fmt.Sprintf("%s/api/v1/download/%s", os.Getenv("BASE_URL"), key)}, nil
}

func (u *usecase) DownloadFile(ctx context.Context, fileIdentifier string) (response *s3.GetObjectOutput, err error) {

	response, err = u.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("BUCKET_NAME")),
		Key:    aws.String(fileIdentifier),
	})
	if err != nil {
		err = fmt.Errorf("failed to download file: %w", err)
	}

	return
}
