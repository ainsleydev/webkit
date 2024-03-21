package seo

import "github.com/a-h/templ"

type Schema interface {
	ToHTML() templ.Component
}
