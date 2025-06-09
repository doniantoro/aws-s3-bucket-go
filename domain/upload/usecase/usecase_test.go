package usecase

import (
	"bytes"
	"errors"
	"fmt"
	"mime/multipart"
	"net/textproto"
	"os"
	"testing"

	"aws-s3-bucket/domain/upload/interfaces"
	"aws-s3-bucket/domain/upload/interfaces/mocks"
	"aws-s3-bucket/models/document"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func initUseCaseUnitTest(t *testing.T) (interfaces.UsecaseInterface, *mocks.S3Interface) {
	os.Setenv("BUCKET_NAME", "test-bucket")
	os.Setenv("BASE_URL", "http://localhost:8080")

	mockS3Client := mocks.NewS3Interface(t)

	var s3Client interfaces.S3Interface = mockS3Client

	return NewUsecase(s3Client), mockS3Client
}

func createMultipartFile(content string, filename string) (multipart.File, *multipart.FileHeader, error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// Simulasi file field
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="file"; filename="`+filename+`"`)
	h.Set("Content-Type", "text/plain")

	part, err := writer.CreatePart(h)
	if err != nil {
		return nil, nil, err
	}

	part.Write([]byte(content))
	writer.Close()

	// Baca kembali untuk parsing multipart
	reqReader := multipart.NewReader(body, writer.Boundary())
	form, err := reqReader.ReadForm(1024)
	if err != nil {
		return nil, nil, err
	}

	files := form.File["file"]
	if len(files) == 0 {
		return nil, nil, err
	}

	fileHeader := files[0]
	file, err := fileHeader.Open()
	if err != nil {
		return nil, nil, err
	}

	return file, fileHeader, nil
}

func Test_DownloadFile(t *testing.T) {
	usecase, mockS3Client := initUseCaseUnitTest(t)
	type args struct {
		request string
	}
	type expected struct {
		err      error
		response *s3.GetObjectOutput
	}
	tests := []struct {
		name     string
		prepare  func(args)
		args     args
		expected expected
	}{
		{
			name: "DownloadFile_Success",
			args: args{
				request: "data/example.txt",
			},
			prepare: func(args args) {
				mockS3Client.On("GetObject", mock.Anything, mock.Anything).Return(nil, nil).Once()
			},
		},
		{
			name: "DownloadFile_Failure",
			args: args{
				request: "data/example.txt",
			},
			prepare: func(args args) {
				mockS3Client.On("GetObject", mock.Anything, mock.Anything).Return(nil, errors.New("failed to download file")).Once()
			},
			expected: expected{
				err:      fmt.Errorf("failed to download file: %w", errors.New("failed to download file")),
				response: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(tt.args)

			response, err := usecase.DownloadFile(nil, tt.args.request)

			require.Equal(t, tt.expected.err, err)
			require.Equal(t, tt.expected.response, response)

		})
	}

}

func Test_UploadFile(t *testing.T) {
	file, fileHeader, err := createMultipartFile("This is test content", "test.txt")
	require.NoError(t, err)
	require.NotNil(t, file)
	require.NotNil(t, fileHeader)

	usecase, mockS3Client := initUseCaseUnitTest(t)
	type args struct {
		request document.RequestUploadDocumentFile
		files   *multipart.FileHeader
	}
	type expected struct {
		err      error
		response document.ResponseUploadDocument
	}
	tests := []struct {
		name     string
		prepare  func(args)
		args     args
		expected expected
	}{
		{
			name: "UploadFile_Success",
			args: args{
				request: document.RequestUploadDocumentFile{
					DocumentKey:  "data",
					DocumentName: "example",
				},
				files: fileHeader,
			},
			prepare: func(args args) {
				mockS3Client.On("PutObject", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil).Once()
			},
			expected: expected{
				response: document.ResponseUploadDocument{
					DocumentUrl: fmt.Sprintf("%s/api/v1/download/%s", os.Getenv("BASE_URL"), "data/example.txt"),
				},
			},
		},
		{
			name: "UploadFile_Failure",
			args: args{
				request: document.RequestUploadDocumentFile{
					DocumentKey:  "data",
					DocumentName: "example.txt",
				},
				files: fileHeader,
			},
			prepare: func(args args) {
				mockS3Client.On("PutObject", mock.Anything, mock.Anything).Return(nil, errors.New("failed to upload file")).Once()
			},
			expected: expected{
				err: fmt.Errorf("failed to upload file: %w", errors.New("failed to upload file")),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(tt.args)

			response, err := usecase.UploadFile(nil, tt.args.request, tt.args.files)

			require.Equal(t, tt.expected.err, err)
			require.Equal(t, tt.expected.response, response)

		})
	}

}

func Test_UploadBase64(t *testing.T) {
	usecase, mockS3Client := initUseCaseUnitTest(t)
	type args struct {
		request document.RequestUploadDocumentBase64
	}
	type expected struct {
		err      error
		response document.ResponseUploadDocument
	}
	tests := []struct {
		name     string
		prepare  func(args)
		args     args
		expected expected
	}{
		{
			name: "UploadBase64_Success",
			args: args{
				request: document.RequestUploadDocumentBase64{
					DocumentKey:    "data",
					DocumentName:   "example",
					DocumentBase64: "data:text/txt;base64,VGhpcyBpcyB0ZXN0IGNvbnRlbnQ=",
				},
			},
			prepare: func(args args) {
				mockS3Client.On("PutObject", mock.Anything, mock.Anything).Return(nil, nil).Once()
			},
			expected: expected{
				response: document.ResponseUploadDocument{
					DocumentUrl: fmt.Sprintf("%s/api/v1/download/%s", os.Getenv("BASE_URL"), "data/example.txt"),
				},
			},
		},
		{
			name: "UploadBase64_Failure",
			args: args{
				request: document.RequestUploadDocumentBase64{
					DocumentKey:    "data",
					DocumentName:   "example",
					DocumentBase64: "data:text/txt;base64,VGhpcyBpcyB0ZXN0IGNvbnRlbnQ=",
				},
			},
			prepare: func(args args) {
				mockS3Client.On("PutObject", mock.Anything, mock.Anything).Return(nil, errors.New("failed to upload file")).Once()
			},
			expected: expected{
				err:      fmt.Errorf("failed to upload file: %w", errors.New("failed to upload file")),
				response: document.ResponseUploadDocument{},
			},
		},
		{
			name: "UploadBase64_InvalidFormat",
			args: args{
				request: document.RequestUploadDocumentBase64{
					DocumentKey:    "data",
					DocumentName:   "example",
					DocumentBase64: "invalid_base64_string",
				},
			},
			expected: expected{
				err:      fmt.Errorf("invalid base64 string format"),
				response: document.ResponseUploadDocument{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare the mock expectations
			if tt.prepare != nil {
				tt.prepare(tt.args)
			}

			response, err := usecase.UploadBase64(nil, tt.args.request)

			require.Equal(t, tt.expected.err, err)
			require.Equal(t, tt.expected.response, response)

		})
	}
}
