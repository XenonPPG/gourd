package gateway

import (
	"bytes"
	"gourd/internal/document"
	"gourd/internal/domain"
	"mime/multipart"

	"github.com/gofiber/fiber/v3"
)

type ProcessDocumentRequest struct {
	File  *multipart.FileHeader `form:"file"`
	Rules []int                 `form:"rule"`
}

// ProcessFile processes an uploaded document file
//
//	@Summary		Process a document file
//	@Description	Uploads a document file and applies the specified processing rules to it, returning the processed binary content
//	@Tags			documents
//	@Accept			multipart/form-data
//	@Produce		application/octet-stream
//	@Param			file	formData	file	true	"Document file to process"
//	@Param			rule	formData	[]int	false	"List of rule IDs to apply"	collectionFormat(multi)
//	@Success		200		{file}		file				"Processed document"
//	@Failure		400		{object}	map[string]string	"Invalid request body"
//	@Failure		500		{object}	map[string]string	"File processing error"
//	@Router			/documents/file [post]
func ProcessFile(c fiber.Ctx) error {
	req := new(ProcessDocumentRequest)
	if err := c.Bind().Body(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	src, err := req.File.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer src.Close()

	var buf bytes.Buffer
	if err := document.Process(src, &buf, req.Rules, req.File.Filename); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	c.Set(fiber.HeaderContentType, "application/octet-stream")
	return c.Send(buf.Bytes())
}

type ProcessRawRequest struct {
	Text     string `json:"text,omitempty"`
	FileType string `json:"type,omitempty"`
	Rules    []int  `json:"rules,omitempty"`
}

// ProcessRaw processes raw text content as a document
//
//	@Summary Process raw text content
//	@Description Applies the specified processing rules to raw text content of a given file type and returns the processed output. Currently only "markdown" is supported as a file type
//	@Tags documents
//	@Accept json
//	@Produce json
//	@Param request body ProcessRawRequest true "Raw content, file type and rules to apply"
//	@Success 201 {object} map[string]string "Processed output"
//	@Failure 201 {object} map[string]string "Unsupported file type"
//	@Failure 400 {object} map[string]string "Invalid request body"
//	@Failure 500 {object} map[string]string "Processing error"
//	@Router /documents/raw [post]
func ProcessRaw(c fiber.Ctx) error {
	req := new(ProcessRawRequest)
	if err := c.Bind().Body(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	fileType := domain.FileType(req.FileType)
	if fileType != domain.Markdown {
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"error": "unsupported file type"})
	}

	output, err := document.ProcessRaw(req.Text, req.Rules, fileType)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"output": output})
}
