package wordpress

type (
	// Post represents a WordPress post.
	Post struct {
		ID            int     `json:"id,omitempty"`
		Date          string  `json:"date,omitempty"`
		DateGMT       string  `json:"date_gmt,omitempty"`
		GUID          GUID    `json:"guid,omitempty"`
		Link          string  `json:"link,omitempty"`
		Modified      string  `json:"modified,omitempty"`
		ModifiedGMT   string  `json:"modifiedGMT,omitempty"` //nolint
		Password      string  `json:"password,omitempty"`
		Slug          string  `json:"slug,omitempty"`
		Status        string  `json:"status,omitempty"`
		Type          string  `json:"type,omitempty"`
		Title         Title   `json:"title,omitempty"`
		Content       Content `json:"content,omitempty"`
		Author        int     `json:"author,omitempty"`
		Excerpt       Excerpt `json:"excerpt,omitempty"`
		FeaturedImage int     `json:"featured_image,omitempty"`
		CommentStatus string  `json:"comment_status,omitempty"`
		PingStatus    string  `json:"ping_status,omitempty"`
		Format        string  `json:"format,omitempty"`
		Sticky        bool    `json:"sticky,omitempty"`
	}
	// GUID represents a WordPress GUID.
	GUID struct {
		Raw      string `json:"raw,omitempty"`
		Rendered string `json:"rendered,omitempty"`
	}
	// Title represents a WordPress title.
	Title struct {
		Raw      string `json:"raw,omitempty"`
		Rendered string `json:"rendered,omitempty"`
	}
	// Content represents a WordPress content.
	Content struct {
		Raw      string `json:"raw,omitempty"`
		Rendered string `json:"rendered,omitempty"`
	}
	// Excerpt represents a WordPress excerpt.
	Excerpt struct {
		Raw      string `json:"raw,omitempty"`
		Rendered string `json:"rendered,omitempty"`
	}
)
