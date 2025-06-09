package interfaces

import (
	"aws-s3-bucket/models/document"
	"context"
	"mime/multipart"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type UsecaseInterface interface {
	UploadBase64(ctx context.Context, request document.RequestUploadDocumentBase64) (response document.ResponseUploadDocument, err error)
	UploadFile(ctx context.Context, request document.RequestUploadDocumentFile, file *multipart.FileHeader) (response document.ResponseUploadDocument, err error)
	DownloadFile(ctx context.Context, fileIdentifier string) (response *s3.GetObjectOutput, err error)
}
