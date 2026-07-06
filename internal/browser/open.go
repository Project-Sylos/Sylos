package browser

import (
	"fmt"

	pkgbrowser "github.com/pkg/browser"
)

func OpenURL(url string) error {
	if err := pkgbrowser.OpenURL(url); err != nil {
		return fmt.Errorf("open browser: %w", err)
	}
	return nil
}
