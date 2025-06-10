package document

type RequestUploadDocumentBase64 struct {
	DocumentKey    string `json:"document_key" validate:"required" example:"folder-in-s3"`
	DocumentName   string `json:"document_name" validate:"required" example:"example"`
	DocumentBase64 string `json:"document_base64" validate:"required" example:"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAIAAACQd1PeAAAAA3NCSVQICAjb4U/gAAAADElEQVQImWMwMTEBAAE8AJ1tV5pJAAAAAElFTkSuQmCC"`
}

type RequestUploadDocumentFile struct {
	DocumentKey  string `json:"document_key" validate:"required" example:"folder-in-s3"`
	DocumentName string `json:"document_name" validate:"required" example:"example"`
}
