package payloadcms

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/gabriel-vasile/mimetype"
)

// MediaService represents a service for managing media within Payload.
type MediaService struct {
	Client *Client
}

// Upload uploads a file to the media endpoint.
func (s MediaService) Upload(ctx context.Context, f *os.File, collection string, in, out any) (Response, error) {
	values, err := getUploadValues(f, in)
	if err != nil {
		return Response{}, err
	}
	return s.uploadFile(ctx, collection, values, out)
}

func (s MediaService) UploadFromURL(ctx context.Context, url string, collection string, in, out any) (Response, error) {
	// Download the file from the URL
	resp, err := s.Client.client.Get(url)
	if err != nil {
		return Response{}, err
	}
	defer resp.Body.Close()

	// Check if the response status is OK
	if resp.StatusCode != http.StatusOK {
		return Response{}, fmt.Errorf("failed to download file: status code %d", resp.StatusCode)
	}

	// Create a temporary file to store the downloaded content
	tmpfile, err := os.Create(filepath.Join(os.TempDir(), fileNameFromURL(url)))
	if err != nil {
		return Response{}, fmt.Errorf("failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpfile.Name()) // Clean up the temporary file

	// Write the downloaded content to the temporary file
	_, err = io.Copy(tmpfile, resp.Body)
	if err != nil {
		return Response{}, fmt.Errorf("failed to write to temporary file: %v", err)
	}

	values, err := getUploadValues(tmpfile, in)
	if err != nil {
		return Response{}, err
	}

	return s.uploadFile(ctx, collection, values, out)
}

// uploadFile prepares a multipart form and performs the upload request
//   - Takes the context, collection name, map of form values (including the file), and optional output struct
//   - Returns a Response object and any errors encountered
func (s MediaService) uploadFile(ctx context.Context, collection string, values map[string]io.Reader, out any) (Response, error) {
	// Prepare a multipart form
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range values {
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		if x, ok := r.(*os.File); ok {
			if err := handleFileUpload(w, key, x); err != nil {
				return Response{}, err
			}
		} else {
			// Add other fields
			fw, err := w.CreateFormField(key)
			if err != nil {
				return Response{}, err
			}
			if _, err := io.Copy(fw, r); err != nil {
				return Response{}, err
			}
		}
	}

	// Close the multipart writer
	if err := w.Close(); err != nil {
		return Response{}, fmt.Errorf("failed to close multipart writer: %v", err)
	}

	p := fmt.Sprintf("/api/%s", collection)
	req, err := s.Client.newFormRequest(ctx, http.MethodPost, p, &b, w.FormDataContentType())
	if err != nil {
		return Response{}, err
	}

	return s.Client.DoWithRequest(ctx, req, &out)
}

func getUploadValues(f *os.File, v any) (map[string]io.Reader, error) {
	if f == nil {
		return nil, fmt.Errorf("file is required")
	}

	values := map[string]io.Reader{
		"file": f,
	}

	// If 'in' is a struct, iterate over its fields and get the JSON tags
	m := reflect.ValueOf(v)
	if m.Kind() == reflect.Struct {
		for i := 0; i < m.NumField(); i++ {
			field := m.Type().Field(i)
			tag := field.Tag.Get("json")
			if tag != "" {
				values[tag] = strings.NewReader(fmt.Sprintf("%v", m.Field(i).Interface()))
			}
		}
	}

	return values, nil
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

func fileNameFromURL(url string) string {
	parts := strings.Split(url, "/")

	// Get the last part of the URL which contains the filename
	filenameWithExtension := parts[len(parts)-1]

	// Extract the filename from filenameWithExtension
	filename := path.Base(filenameWithExtension)

	return filename
}
