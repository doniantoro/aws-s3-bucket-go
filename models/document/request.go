package document

type RequestUploadDocumentBase64 struct {
	DocumentKey    string `json:"document_key" validate:"required"`
	DocumentName   string `json:"document_name" validate:"required"`
	DocumentBase64 string `json:"document_base64" validate:"required"`
}

type RequestUploadDocumentFile struct {
	DocumentKey  string `json:"document_key" validate:"required"`
	DocumentName string `json:"document_name" validate:"required"`
}
