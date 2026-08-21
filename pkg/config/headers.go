package config

import "net/http"

// CopyPreservedHeaders copies from src into dst the headers named in names.
// Header name matching is case-insensitive. Values are copied (not shared).
// If names is empty, this function is a no-op.
func CopyPreservedHeaders(dst, src http.Header, names []string) {
	if dst == nil || src == nil || len(names) == 0 {
		return
	}

	allowed := make(map[string]struct{}, len(names))
	for _, name := range names {
		if name == "" {
			continue
		}
		allowed[http.CanonicalHeaderKey(name)] = struct{}{}
	}
	if len(allowed) == 0 {
		return
	}

	for key, values := range src {
		if _, ok := allowed[http.CanonicalHeaderKey(key)]; !ok {
			continue
		}
		dst[key] = append([]string(nil), values...)
	}
}
