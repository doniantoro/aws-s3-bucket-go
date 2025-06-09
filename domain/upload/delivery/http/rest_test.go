package delivery

import (
	configMocks "aws-s3-bucket/config/interfaces/mocks"
	"aws-s3-bucket/domain/upload/interfaces"
	"aws-s3-bucket/domain/upload/interfaces/mocks"
	"aws-s3-bucket/models/document"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-playground/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func initRestUnitTest(t *testing.T) (*handler, *mocks.UsecaseInterface, *configMocks.MockValidator) {
	mockUsecase := mocks.NewUsecaseInterface(t)

	var usecase interfaces.UsecaseInterface = mockUsecase

	mockValidator := new(configMocks.MockValidator)

	return &handler{usecase: usecase, validator: mockValidator}, mockUsecase, mockValidator
}

func TestNewHandler(t *testing.T) {
	mockUsecase := mocks.NewUsecaseInterface(t)

	var usecase interfaces.UsecaseInterface = mockUsecase

	mockValidator := new(configMocks.MockValidator)
	NewHandler(fiber.New(), usecase, mockValidator)
	return
}

func TestUploadBase64(t *testing.T) {

	handler, mockUsecase, mockValidator := initRestUnitTest(t)

	type args struct {
		request string
	}
	type expected struct {
		statusCode int
	}
	tests := []struct {
		name     string
		prepare  func(args)
		args     args
		expected expected
	}{
		{
			name: "failed to parse body",
			args: args{
				request: "invalid-json",
			},
			expected: expected{
				statusCode: fiber.StatusBadRequest,
			},
		},
		{
			name: "validation failed",
			args: args{
				request: `{
					"DocumentName": "test"
				}`,
			},
			prepare: func(args args) {
				mockValidator.On("Validate", mock.Anything).Return(validator.ValidationErrors{
					configMocks.MockFieldError{Fields: "DocumentName", Tags: "required", Params: ""},
					configMocks.MockFieldError{Fields: "DocumentKey", Tags: "required", Params: ""},
				}).Once()
			},
			expected: expected{
				statusCode: fiber.StatusBadRequest,
			},
		},
		{
			name: "upload failed",
			args: args{
				request: `{
					"document_name": "test",
					"document_key": "test-key",
					"document_base64": "ZGF0YQ=="
				}`,
			},
			prepare: func(args args) {
				mockValidator.On("Validate", mock.Anything).Return(nil).Once()
				mockUsecase.On("UploadBase64", mock.Anything, mock.Anything).Return(document.ResponseUploadDocument{}, errors.New("upload failed")).Once()
			},
			expected: expected{
				statusCode: fiber.StatusInternalServerError,
			},
		},
		{
			name: "upload success",
			args: args{
				request: `{
					"document_name": "test",
					"document_key": "folder-in-s3",
					"document_base64": "ZGF0YQ=="
				}`,
			},
			prepare: func(args args) {
				mockValidator.On("Validate", mock.Anything).Return(nil).Once()
				mockUsecase.On("UploadBase64", mock.Anything, mock.Anything).Return(document.ResponseUploadDocument{
					DocumentUrl: "localhost:8090/api/v1/download/folder-in-s3/test.png",
				}, nil).Once()
			},
			expected: expected{
				statusCode: fiber.StatusCreated,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare(tt.args)
			}
			apps := fiber.New()
			apps.Post("upload/base64", handler.UploadBase64)

			req := httptest.NewRequest(http.MethodPost, "/upload/base64", bytes.NewBuffer([]byte(tt.args.request)))
			req.Header.Set("Content-Type", "application/json")

			resp, _ := apps.Test(req)

			require.Equal(t, tt.expected.statusCode, resp.StatusCode)
		})
	}

}

func TestUploadFile_Usecase(t *testing.T) {

	handler, mockUsecase, mockValidator := initRestUnitTest(t)

	type args struct {
		request []struct {
			key   string
			value string
		}
		isRequestFile bool
	}
	type expected struct {
		statusCode int
	}
	tests := []struct {
		name     string
		prepare  func(args)
		args     args
		expected expected
	}{

		{
			name: "validation file failed",
			args: args{
				request: []struct {
					key   string
					value string
				}{
					{key: "document_key", value: "test-key"},
				},
			},

			expected: expected{
				statusCode: fiber.StatusInternalServerError,
			},
		},
		{
			name: "validation failed",
			args: args{
				request: []struct {
					key   string
					value string
				}{
					{key: "document_key", value: "test-key"},
				},
				isRequestFile: true,
			},
			prepare: func(args args) {
				mockValidator.On("Validate", mock.Anything).Return(validator.ValidationErrors{
					configMocks.MockFieldError{Fields: "DocumentName", Tags: "required", Params: ""},
				}).Once()
			},
			expected: expected{
				statusCode: fiber.StatusBadRequest,
			},
		},
		{
			name: "upload failed",
			args: args{
				request: []struct {
					key   string
					value string
				}{
					{key: "document_key", value: "test-key"},
					{key: "document_name", value: "test"},
				},
				isRequestFile: true,
			},
			prepare: func(args args) {

				mockValidator.On("Validate", mock.Anything).Return(nil).Once()
				mockUsecase.On("UploadFile", mock.Anything, mock.Anything, mock.Anything).
					Return(document.ResponseUploadDocument{}, errors.New("usecase failed")).Once()
			},
			expected: expected{
				statusCode: fiber.StatusInternalServerError,
			},
		},
		{
			name: "upload success",
			args: args{
				request: []struct {
					key   string
					value string
				}{
					{key: "document_key", value: "test-key"},
					{key: "document_name", value: "test"},
				},
				isRequestFile: true,
			},
			prepare: func(args args) {

				mockValidator.On("Validate", mock.Anything).Return(nil).Once()
				mockUsecase.On("UploadFile", mock.Anything, mock.Anything, mock.Anything).
					Return(document.ResponseUploadDocument{
						DocumentUrl: "localhost:8090/api/v1/download/folder-in-s3/test.png",
					}, nil).Once()
			},
			expected: expected{
				statusCode: fiber.StatusCreated,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare(tt.args)
			}
			apps := fiber.New()
			apps.Post("upload/file", handler.UploadFile)

			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)

			for _, row := range tt.args.request {

				_ = writer.WriteField(row.key, row.value)
			}
			if tt.args.isRequestFile {
				fw, _ := writer.CreateFormFile("file", "test.txt")
				fw.Write([]byte("test"))
			}

			writer.Close()

			req := httptest.NewRequest(http.MethodPost, "/upload/file", body)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Content-Type", writer.FormDataContentType())

			resp, _ := apps.Test(req)
			log.Println("resp", resp)

			require.Equal(t, tt.expected.statusCode, resp.StatusCode)
		})
	}

}

type ErrReader struct{}

func (e ErrReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("read error")
}
func TestGetFile_SuccessInline(t *testing.T) {

	handler, mockUsecase, _ := initRestUnitTest(t)

	type args struct {
		docKey       string
		docName      string
		typeDocument string
	}
	type expected struct {
		statusCode int
		err        error
	}
	tests := []struct {
		name     string
		prepare  func(args)
		args     args
		expected expected
	}{
		{
			name: "GetFile_SuccessInline",
			args: args{
				docKey:  "abc",
				docName: "file.txt",
			},
			prepare: func(a args) {
				fileContent := "Hello Fiber"
				fileReader := io.NopCloser(bytes.NewReader([]byte(fileContent)))

				getOutput := &s3.GetObjectOutput{
					Body:        fileReader,
					ContentType: aws.String("text/plain"),
				}

				mockUsecase.On("DownloadFile", mock.Anything, mock.Anything).Return(getOutput, nil).Once()
			},
			expected: expected{
				statusCode: fiber.StatusOK,
			},
		},
		{
			name: "GetFile_SuccessDownload",
			args: args{
				docKey:       "abc",
				docName:      "file.txt",
				typeDocument: "download",
			},
			prepare: func(a args) {
				fileContent := "Hello Fiber"
				fileReader := io.NopCloser(bytes.NewReader([]byte(fileContent)))

				getOutput := &s3.GetObjectOutput{
					Body:        fileReader,
					ContentType: aws.String("text/plain"),
				}

				mockUsecase.On("DownloadFile", mock.Anything, mock.Anything).Return(getOutput, nil).Once()
			},
			expected: expected{
				statusCode: fiber.StatusOK,
			},
		},
		{
			name: "GetFile_Successbase64",
			args: args{
				docKey:       "abc",
				docName:      "file.txt",
				typeDocument: "base64",
			},
			prepare: func(a args) {
				fileContent := "Hello Fiber"
				fileReader := io.NopCloser(bytes.NewReader([]byte(fileContent)))

				getOutput := &s3.GetObjectOutput{
					Body:        fileReader,
					ContentType: aws.String("text/plain"),
				}

				mockUsecase.On("DownloadFile", mock.Anything, mock.Anything).Return(getOutput, nil).Once()
			},
			expected: expected{
				statusCode: fiber.StatusOK,
			},
		},
		{
			name: "GetFile_ErrorDownload",
			args: args{
				docKey:       "abc",
				docName:      "file.txt",
				typeDocument: "base64",
			},
			prepare: func(a args) {
				getOutput := &s3.GetObjectOutput{
					Body:        io.NopCloser(bytes.NewReader([]byte("Hello Fiber"))),
					ContentType: aws.String("text/plain"),
				}
				mockUsecase.On("DownloadFile", mock.Anything, mock.Anything).Return(getOutput, errors.New("connection error")).Once()
			},
			expected: expected{
				statusCode: fiber.StatusInternalServerError,
			},
		},
		{
			name: "GetFile_ErrorIoCopy",
			args: args{
				docKey:       "abc",
				docName:      "file.txt",
				typeDocument: "base64",
			},
			prepare: func(a args) {

				getOutput := &s3.GetObjectOutput{
					Body:        io.NopCloser(ErrReader{}),
					ContentType: aws.String("text/plain"),
				}
				mockUsecase.On("DownloadFile", mock.Anything, mock.Anything).Return(getOutput, nil).Once()
			},
			expected: expected{
				statusCode: fiber.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare(tt.args)
			}
			app := fiber.New()
			app.Get("/file/:docKey/:docName", handler.GetFile)

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/file/%s/%s?type=%s", tt.args.docKey, tt.args.docName, tt.args.typeDocument), nil)
			resp, err := app.Test(req)

			require.Equal(t, tt.expected.err, err)
			require.Equal(t, tt.expected.statusCode, resp.StatusCode)
			if tt.args.typeDocument == "download" {
				require.Contains(t, resp.Header.Get("Content-Disposition"), "attachment;")
			} else if tt.args.typeDocument != "base64" {
				require.Contains(t, resp.Header.Get("Content-Disposition"), "inline;")
			}

		})
	}

}
