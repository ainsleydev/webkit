package templates

import "fmt"

// secret wraps a secret name in GitHub Actions syntax
func secret(name string) string {
	return fmt.Sprintf("{{ secrets.%s }}", name)
}
