package gateway

import (
	"bytes"
	"gourd/internal/common"
	"log"
	"mime/multipart"

	"github.com/gofiber/fiber/v3"
)

// ProcessDocumentRequest is the multipart form body for ProcessFile.
type ProcessDocumentRequest struct {
	File     *multipart.FileHeader `form:"file"`
	Rules    []string              `form:"rules"`
	FileType string                `form:"fileType"`
}

// ProcessFile processes an uploaded document file
//
//	@Summary       Process a document file
//	@Description   Uploads a document file and applies the specified processing rules to it, returning the processed binary content. The file type must be given in the fileType field.
//	@Tags          documents
//	@Accept        multipart/form-data
//	@Produce       application/octet-stream
//	@Param         file      formData  file      true   "Document file to process"
//	@Param         fileType  formData  string    true   "Type of the uploaded document"
//	@Param         rules     formData  []string  false  "List of rule IDs to apply"  collectionFormat(multi)
//	@Success       200       {file}    file                "Processed document"
//	@Failure       400       {object}  map[string]string   "Invalid request body"
//	@Failure       500       {object}  map[string]string   "File read or processing error (including unsupported file type)"
//	@Router        /documents/file [post]
func (s *Service) ProcessFile(c fiber.Ctx) error {
	// parse request
	req := new(ProcessDocumentRequest)
	if err := c.Bind().Body(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// read file
	src, err := req.File.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer func(f multipart.File) {
		if cErr := f.Close(); cErr != nil {
			log.Println("failed to close file: ", cErr)
		}
	}(src)

	// process
	var buf bytes.Buffer
	if err = s.documentService.Process(c.Context(), src, &buf, common.FileType(req.FileType), req.Rules); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	c.Set(fiber.HeaderContentType, "application/octet-stream")
	return c.Send(buf.Bytes())
}

// ProcessRawRequest is the JSON body for ProcessRaw.
type ProcessRawRequest struct {
	Text     string   `json:"text,omitempty"`
	FileType string   `json:"type,omitempty"`
	Rules    []string `json:"rules,omitempty"`
}

// ProcessRaw processes raw text content as a document
//
//	@Summary      Process raw text content
//	@Description  Applies the specified processing rules to raw text content of a given file type and returns the processed output in the "output" field. Currently only "markdown" is supported as a file type.
//	@Tags         documents
//	@Accept       json
//	@Produce      json
//	@Param        request  body      ProcessRawRequest  true  "Raw content, file type and rules to apply"
//	@Success      201      {object}  map[string]string        "Processed output, returned as {\"output\": \"...\"}"
//	@Failure      400      {object}  map[string]string        "Invalid request body"
//	@Failure      500      {object}  map[string]string        "Processing error (including unsupported file type)"
//	@Router       /documents/raw [post]
func (s *Service) ProcessRaw(c fiber.Ctx) error {
	req := new(ProcessRawRequest)
	if err := c.Bind().Body(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	output, err := s.documentService.ProcessRaw(c.Context(), req.Text, req.Rules, common.FileType(req.FileType))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"output": output})
}
