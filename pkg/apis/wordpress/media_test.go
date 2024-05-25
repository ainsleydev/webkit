package wordpress

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_Media(t *testing.T) {
	tt := map[string]struct {
		id             int
		serverResponse string
		serverStatus   int
		wantMedia      Media
		wantErr        bool
	}{
		"Success": {
			id:             123,
			serverResponse: `{"id": 123, "title": {"rendered": "Sample Media Title"}, "link": "https://example.com/sample-media", "author": 1, "media_type": "image"}`,
			serverStatus:   http.StatusOK,
			wantMedia: Media{
				ID:        123,
				Title:     Title{Rendered: "Sample Media Title"},
				Link:      "https://example.com/sample-media",
				Author:    1,
				MediaType: "image",
			},
			wantErr: false,
		},
		"Failure": {
			id:           456,
			serverStatus: http.StatusNotFound,
			wantMedia:    Media{},
			wantErr:      true,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			client, teardown := Setup(t, test.serverResponse, test.serverStatus)
			defer teardown()

			gotMedia, err := client.Media(context.TODO(), test.id)

			assert.Equal(t, test.wantMedia, gotMedia)
			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}
