package payloadcms

import (
	"bytes"
	"context"
	"fmt"
	"html"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path"
	"strings"

	"github.com/gabriel-vasile/mimetype"

	"github.com/ainsleydev/webkit/pkg/util/httputil"
)

// MediaService represents a service for managing media within Payload.
type MediaService struct {
	Client *Client
}

type UploadRequest struct {
	Alt     string `json:"alt"`
	Caption string
}

func (m *MediaService) createFormData() (*bytes.Buffer, *multipart.Writer) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	return body, writer
}

func (m *MediaService) UploadFromURL(ctx context.Context, url string, altText string) error {
	// Download the file from the URL
	resp, err := m.Client.client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check if the response status is OK
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file: status code %d", resp.StatusCode)
	}

	// Create a temporary file to store the downloaded content
	tmpfile, err := os.CreateTemp("", "downloaded_file_")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpfile.Name()) // Clean up the temporary file

	// Write the downloaded content to the temporary file
	_, err = io.Copy(tmpfile, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write to temporary file: %v", err)
	}

	// Get file name from URL
	fileName := GetFileNameFromURL(url)

	// Pass the response body (file content) to the Upload method
	err = m.Upload(ctx, tmpfile, fileName)
	if err != nil {
		return err
	}

	return nil
}

func (m *MediaService) Upload(ctx context.Context, reader io.ReadSeeker, fileName string) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	_, err := reader.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	// Detect MIME type
	mime, err := mimetype.DetectReader(reader)
	if err != nil {
		return err
	}

	// Create a new file part manually with the desired headers
	partHeaders := textproto.MIMEHeader{}
	partHeaders.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, html.EscapeString(fileName)))
	partHeaders.Set("Content-Type", mime.String())
	part, err := writer.CreatePart(partHeaders)
	if err != nil {
		return err
	}

	// Copy file content to the form file part
	_, err = io.Copy(part, reader)
	if err != nil {
		return err
	}

	// Write the "alt" field to the multipart form data
	err = writer.WriteField("alt", "ALT")
	if err != nil {
		return err
	}

	// Close the writer to finalize the multipart form data
	err = writer.Close()
	if err != nil {
		return err
	}

	path := fmt.Sprintf("%s/api/%s", m.Client.baseURL, CollectionMedia)
	req, err := http.NewRequestWithContext(ctx, "POST", path, body)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "users API-KEY "+m.Client.apiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := m.Client.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(buf))

	if !httputil.Is2xx(resp.StatusCode) {
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, buf)
	}

	// Check response status here if needed
	// For now, simply returning nil (no error) assuming upload is successful
	return nil
}

// GetFileNameFromURL extracts the filename from a given URL.
func GetFileNameFromURL(url string) string {
	// Split URL by '/'
	parts := strings.Split(url, "/")

	// Get the last part of the URL which contains the filename
	filenameWithExtension := parts[len(parts)-1]

	// Extract the filename from filenameWithExtension
	filename := path.Base(filenameWithExtension)

	return filename
}
