package payloadcms

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/goccy/go-json"

	"github.com/ainsleydev/webkit/pkg/util/httputil"
)

type (
	// MediaService represents a service for managing media within Payload.
	MediaService struct {
		Client *Client
	}
	// UploadRequest represents a request to the media upload endpoint.
	UploadRequest struct {
		Alt     string `json:"alt"`
		Caption string `json:"caption"`
	}
	// UploadResponse represents a response from the media upload endpoint.
	UploadResponse struct {
		ID        int         `json:"id"`
		Alt       string      `json:"alt"`
		Caption   interface{} `json:"caption"`
		UpdatedAt time.Time   `json:"updatedAt"`
		CreatedAt time.Time   `json:"createdAt"`
		Url       string      `json:"url"`
		Filename  string      `json:"filename"`
		MimeType  string      `json:"mimeType"`
		Filesize  int         `json:"filesize"`
		Width     interface{} `json:"width"`
		Height    interface{} `json:"height"`
	}
)

// Upload uploads a file to the media endpoint.
func (m MediaService) Upload(ctx context.Context, f *os.File, request UploadRequest) (*CreateResponse[UploadResponse], error) {
	values := map[string]io.Reader{
		"file":    f,
		"alt":     strings.NewReader("fuck"),
		"caption": strings.NewReader(request.Caption),
	}
	return m.upload(ctx, values)
}

func (m MediaService) UploadFromURL(ctx context.Context, url string, request UploadRequest) (*CreateResponse[UploadResponse], error) {
	// Download the file from the URL
	resp, err := m.Client.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check if the response status is OK
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download file: status code %d", resp.StatusCode)
	}

	// Create a temporary file to store the downloaded content
	tmpfile, err := os.Create(filepath.Join(os.TempDir(), GetFileNameFromURL(url)))
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpfile.Name()) // Clean up the temporary file

	// Write the downloaded content to the temporary file
	_, err = io.Copy(tmpfile, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to write to temporary file: %v", err)
	}

	// Pass the response body (file content) to the Upload method
	return m.Upload(ctx, tmpfile, request)
}

func (m MediaService) upload(ctx context.Context, values map[string]io.Reader) (*CreateResponse[UploadResponse], error) {
	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range values {
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		if x, ok := r.(*os.File); ok {
			if err := handleFileUpload(w, key, x); err != nil {
				return nil, err
			}
		} else {
			// Add other fields
			fw, err := w.CreateFormField(key)
			if err != nil {
				return nil, err
			}
			if _, err := io.Copy(fw, r); err != nil {
				return nil, err
			}
		}
	}

	// Close the multipart writer or the request may
	// be missing the terminating boundary.
	w.Close()

	// Now that you have a form, you can submit it to your handler.
	path := fmt.Sprintf("%s/api/%s", m.Client.baseURL, CollectionMedia)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, path, &b)
	if err != nil {
		return nil, err
	}

	// Set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Add("Authorization", "users API-Key "+m.Client.apiKey)

	// Submit the request
	res, err := m.Client.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if !httputil.Is2xx(res.StatusCode) {
		return nil, errors.New("failed to upload media, status code: " + res.Status)
	}

	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var response CreateResponse[UploadResponse]
	err = json.Unmarshal(buf, &response)
	return &response, err
}

// handleFileUpload adds a file to the multipart writer.
func handleFileUpload(w *multipart.Writer, key string, f *os.File) error {
	// Open the file to read its contents and detect the MIME type.
	file, err := os.Open(f.Name())
	if err != nil {
		return err
	}
	defer file.Close()

	// Read the first 512 bytes to detect the MIME type.
	mime, err := mimetype.DetectFile(file.Name())
	if err != nil {
		return err
	}

	// Create a new form part
	fw, err := w.CreatePart(textproto.MIMEHeader{
		"Content-Type": {
			mime.String(),
		},
		"Content-Disposition": {
			fmt.Sprintf(`form-data; name="%s"; filename="%s"`, key, f.Name()),
		},
	})
	if err != nil {
		return err
	}

	// Copy the remaining file contents to the form part
	_, err = io.Copy(fw, file)
	return err
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
