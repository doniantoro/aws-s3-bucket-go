package delivery

import (
	configApp "aws-s3-bucket/config/interfaces"
	"aws-s3-bucket/domain/upload/interfaces"
	"aws-s3-bucket/models/document"
	"aws-s3-bucket/models/dto"
	"aws-s3-bucket/shared/constant"
	"aws-s3-bucket/shared/utils"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type handler struct {
	usecase   interfaces.UsecaseInterface
	validator configApp.Validator
}

func NewHandler(route fiber.Router, usecase interfaces.UsecaseInterface, validator configApp.Validator) {
	handler := handler{
		usecase:   usecase,
		validator: validator,
	}

	route.Post("upload/base64", handler.UploadBase64)
	route.Post("upload/file", handler.UploadFile)
	route.Get("download/:docKey/:docName", handler.GetFile)

}

// Integrator godoc
// @Description  orchestrator to upload base64 to s3
// @Produce json
// @Param body body document.RequestUploadDocumentBase64 true "Body payload"
// @Success 200 {object} dto.ApiResponse{data=document.ResponseUploadDocument}
// @Failure 400 {object} dto.ApiResponse{error=dto.ErrorValidation}
// @Failure 500 {object} dto.ApiResponse{}
// @Router /api/v1/upload/base64 [post]
func (h *handler) UploadBase64(c *fiber.Ctx) error {
	var request document.RequestUploadDocumentBase64

	if err := c.BodyParser(&request); err != nil {
		log.Error("Error parsing request body")
		return c.Status(http.StatusBadRequest).JSON(dto.ApiResponse{
			Code:       constant.STATUS_CODE_PARSING_REQUEST,
			Message:    "Failed to parse request body",
			ServerTime: time.Now().Format(time.RFC3339),
		})
	}

	err := h.validator.Validate(&request)
	if err != nil {
		log.Error("Validation error", err)
		return c.Status(http.StatusBadRequest).JSON(dto.ApiResponse{
			Code:       constant.STATUS_CODE_VALIDATION_ERROR,
			Errors:     utils.UnwrapValidation(err),
			Message:    "Validation failed",
			ServerTime: time.Now().Format(time.RFC3339),
		})
	}

	response, err := h.usecase.UploadBase64(c.Context(), request)
	if err != nil {
		log.Error("Error to upload base64")
		return c.Status(http.StatusInternalServerError).JSON(dto.ApiResponse{
			Code:       constant.STATUS_CODE_GENERAL_ERROR,
			Message:    "Failed upload document",
			ServerTime: time.Now().Format(time.RFC3339),
		})
	}

	defer log.Info("Document uploaded successfully", "document_key", request.DocumentKey, "document_name", request.DocumentName)

	return c.Status(http.StatusCreated).JSON(dto.ApiResponse{
		Code:       constant.STATUS_CODE_GENERAL_SUCCESS,
		Message:    "Document uploaded successfully",
		Data:       response,
		ServerTime: time.Now().Format(time.RFC3339),
	})

}

// Integrator godoc
// @Description  orchestrator to upload base64 to s3
// @Produce json
// @Param file formData file true "file document"
// @Param document_key formData string true "key document" default(folder-in-s3)
// @Param document_name formData string true "name document" default(example)
// @Success 200 {object} dto.ApiResponse{data=document.ResponseUploadDocument}
// @Failure 400 {object} dto.ApiResponse{error=dto.ErrorValidation}
// @Failure 500 {object} dto.ApiResponse{}
// @Router /api/v1/upload/file [post]
func (h *handler) UploadFile(c *fiber.Ctx) error {

	file, err := c.FormFile("file")
	if err != nil {
		log.Error("Error Get file from form")
		return c.Status(http.StatusInternalServerError).JSON(dto.ApiResponse{
			Code:       constant.STATUS_CODE_PARSING_REQUEST,
			Message:    "Error get file from form",
			ServerTime: time.Now().Format(time.RFC3339),
		})
	}

	documentName := c.FormValue("document_key")
	documenKey := c.FormValue("document_name")
	request := document.RequestUploadDocumentFile{
		DocumentKey:  documentName,
		DocumentName: documenKey,
	}

	err = h.validator.Validate(request)
	if err != nil {
		log.Error("Validation error", err)
		return c.Status(http.StatusBadRequest).JSON(dto.ApiResponse{
			Code:       constant.STATUS_CODE_VALIDATION_ERROR,
			Errors:     utils.UnwrapValidation(err),
			Message:    "Validation failed",
			ServerTime: time.Now().Format(time.RFC3339),
		})
	}

	response, err := h.usecase.UploadFile(c.Context(), request, file)
	if err != nil {
		log.Error("Error usecase to upload file")
		return c.Status(http.StatusInternalServerError).JSON(dto.ApiResponse{
			Code:       constant.STATUS_CODE_GENERAL_ERROR,
			Message:    "Failed to upload document",
			ServerTime: time.Now().Format(time.RFC3339),
		})
	}
	defer log.Info("Document uploaded successfully", "document_key", request.DocumentKey, "document_name", request.DocumentName)

	return c.Status(http.StatusCreated).JSON(dto.ApiResponse{
		Code:       constant.STATUS_CODE_GENERAL_SUCCESS,
		Message:    "Document uploaded successfully",
		Data:       response,
		ServerTime: time.Now().Format(time.RFC3339),
	})
}

// Integrator godoc
// @Description  orchestrator to get base64 to s3
// @Produce json
// @Param type query string  false "type downloaded can be empty(file),downloaded and base64" default(base64)
// @Param docKey path string true "document key" default(folder-in-s3)
// @Param docName path string true "document name" default(example.png)
// @Failure 200 {object} dto.ApiResponse{}
// @Failure 500 {object} dto.ApiResponse{}
// @Failure 400 {object} dto.ApiResponse{}
// @Router /api/v1/download/{docKey}/{docName} [get]
func (h *handler) GetFile(c *fiber.Ctx) error {

	typeResponse := c.Query("type")
	response, err := h.usecase.DownloadFile(c.Context(), fmt.Sprintf("%s/%s", c.Params("docKey"), c.Params("docName")))
	if err != nil {
		log.Error("Error to get file :%s", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(dto.ApiResponse{
			Code:       constant.STATUS_CODE_GENERAL_ERROR,
			Message:    "Failed to get document",
			ServerTime: time.Now().Format(time.RFC3339),
		})
	}
	defer response.Body.Close()

	var buf bytes.Buffer
	tee := io.TeeReader(response.Body, &buf)
	_, err = io.Copy(c.Response().BodyWriter(), tee)
	if err != nil {
		log.Error("Error to copy document")
		return c.Status(http.StatusInternalServerError).JSON(dto.ApiResponse{
			Code:       "5002",
			Message:    "Failed to get document",
			ServerTime: time.Now().Format(time.RFC3339),
		})
	}

	if typeResponse != "base64" {
		diposition := "inline"
		if typeResponse == "download" {
			diposition = "attachment"
		}
		c.Status(http.StatusOK)
		c.Append("Content-Type", *response.ContentType)
		c.Append("Content-Length", fmt.Sprintf("%d", response.ContentLength))
		c.Append("Content-Disposition", fmt.Sprintf("%s; filename=\"%s\"", diposition, c.Params("docName")))

		return nil
	} else {
		return c.Status(http.StatusOK).JSON(dto.ApiResponse{
			Code:    constant.STATUS_CODE_GENERAL_SUCCESS,
			Message: "Document Get Data successfully",
			Data: map[string]interface{}{
				"document_base64": base64.StdEncoding.EncodeToString(buf.Bytes()),
			},
			ServerTime: time.Now().Format(time.RFC3339),
		})
	}

}
