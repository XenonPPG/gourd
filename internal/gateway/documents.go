package gateway

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"gourd/internal/common"
	"log"
	"mime/multipart"

	"github.com/gofiber/fiber/v3"
)

// ProcessDocumentRequest is the multipart form body for ProcessFile.
type ProcessDocumentRequest struct {
	File     *multipart.FileHeader `form:"file"`
	Rules    []string              `form:"rules"`
	FileType string                `form:"file_type"`
}

// ProcessDocumentResponseMeta is the processing metadata returned in the
// X-Process-Meta response header (JSON, base64-encoded).
type ProcessDocumentResponseMeta struct {
	FileType string `json:"file_type"`
}

const xProcessMeta = "X-Process-Meta"

// ProcessFile processes an uploaded document file
//
//	@Summary       Process a document file
//	@Description   Uploads a document file and applies the specified processing rules to it, returning the processed binary content.
//	@Description   The file type must be given in the file_type field.
//	@Description   At least one rule must be specified in the rules field.
//	@Description   Processing metadata is returned in the X-Process-Meta response header: a base64-encoded JSON object, e.g. {"file_type":"docx"}.
//	@Tags          documents
//	@Accept        multipart/form-data
//	@Produce       application/octet-stream
//	@Param         file       formData  file      true   "Document file to process"
//	@Param         file_type  formData  string    true   "Type of the uploaded document"
//	@Param         rules      formData  []string  true   "List of rule IDs to apply (at least one required)"  collectionFormat(multi)
//	@Success       200        {file}    file                "Processed document"
//	@Header        200        {string}  X-Process-Meta      "Base64-encoded JSON with processing metadata (file_type)"
//	@Failure       400        {object}  map[string]string   "Invalid request body or no rules specified"
//	@Failure       500        {object}  map[string]string   "File read or processing error (including unsupported file type)"
//	@Router        /documents/file [post]
func (s *Service) ProcessFile(c fiber.Ctx) error {
	// parse request
	req := new(ProcessDocumentRequest)
	if err := c.Bind().Body(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if len(req.Rules) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "no rules"})
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
	meta := ProcessDocumentResponseMeta{
		FileType: req.FileType,
	}
	headerMeta, err := json.Marshal(meta)
	if err != nil {
		log.Println("failed to marshal header meta")
	}

	c.Set(fiber.HeaderContentType, "application/octet-stream")
	c.Set(xProcessMeta, base64.StdEncoding.EncodeToString(headerMeta))

	return c.Send(buf.Bytes())
}

// ProcessRawRequest is the JSON body for ProcessRaw.
type ProcessRawRequest struct {
	Text     string   `json:"text,omitempty"`
	FileType string   `json:"type,omitempty" binding:"required"`
	Rules    []string `json:"rules,omitempty"`
}

// ProcessRaw processes raw text content as a document
//
//	@Summary      Process raw text content
//	@Description  Applies the specified processing rules to raw text content of a given file type and returns the processed output in the "output" field. Currently only "markdown" is supported as a file type.
//	@Description  At least one rule must be specified in the rules field.
//	@Tags         documents
//	@Accept       json
//	@Produce      json
//	@Param        request  body      ProcessRawRequest  true  "Raw content, file type and rules to apply (rules must not be empty)"
//	@Success      201      {object}  map[string]string        "Processed output, returned as {\"output\": \"...\"}"
//	@Failure      400      {object}  map[string]string        "Invalid request body or no rules specified"
//	@Failure      500      {object}  map[string]string        "Processing error (including unsupported file type)"
//	@Router       /documents/raw [post]
func (s *Service) ProcessRaw(c fiber.Ctx) error {
	req := new(ProcessRawRequest)
	if err := c.Bind().Body(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if len(req.Rules) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "no rules"})
	}

	output, err := s.documentService.ProcessRaw(c.Context(), req.Text, req.Rules, common.FileType(req.FileType))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"output": output})
}
