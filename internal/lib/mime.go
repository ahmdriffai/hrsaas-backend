package lib

import "mime"

func extensionFromMime(mimeType string) string {
	exts, _ := mime.ExtensionsByType(mimeType)

	if len(exts) > 0 {
		return exts[0]
	}

	return ""
}
