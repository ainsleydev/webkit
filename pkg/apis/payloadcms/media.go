package payloadcms

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

// MediaService represents a service for managing media within Payload.
type MediaService struct {
	Client *Client
}

func (m *MediaService) Upload(ctx context.Context, reader io.Reader, fileName string) error {

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return err
	}

	_, err = io.Copy(part, reader)
	if err != nil {
		return err
	}
	err = writer.Close()
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/api/%s", CollectionMedia)
	req, err := http.NewRequestWithContext(ctx, "POST", path, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := m.Client.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, buf)
	}

	// Check response status here if needed
	// For now, simply returning nil (no error) assuming upload is successful
	return nil
}
