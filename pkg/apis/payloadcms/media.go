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

	"github.com/ainsleydev/webkit/pkg/util/httputil"
)

// MediaService represents a service for managing media within Payload.
type MediaService struct {
	Client *Client
}

func (m *MediaService) Upload(ctx context.Context, reader io.Reader, fileName string) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Create a new file part manually with the desired headers
	partHeaders := textproto.MIMEHeader{}
	partHeaders.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, html.EscapeString(fileName)))
	partHeaders.Set("Content-Type", "image/jpeg") // or the appropriate MIME type for the file being uploaded
	part, err := writer.CreatePart(partHeaders)
	if err != nil {
		return err
	}

	// Copy file content to the form file part
	_, err = io.Copy(part, reader)
	if err != nil {
		return err
	}

	// Detect MIME type
	//mime, err := mimetype.DetectReader(reader)
	//if err != nil {
	//	return err
	//}

	// Write the "alt" field to the multipart form data
	err = writer.WriteField("alt", "ALT")
	if err != nil {
		return err
	}

	fmt.Println(writer.FormDataContentType())

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

	req.Header.Add("Authorization", "users API-KEY 46aabd6a-7303-4db3-a4ce-40625f47fd93")

	fmt.Println(writer.FormDataContentType())

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
