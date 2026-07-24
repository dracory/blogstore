package blogstore

import (
	"context"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
)

// defaultMediaPrefix is the default URL prefix for blog media serving.
const defaultMediaPrefix = "/blog/media/"

// MediaHandler is a standalone HTTP handler that serves blog media files.
// It looks up the media record by ID and serves the content from whatever
// GetURL() the record holds (data URI, HTTP URL, file path).
//
// The URL prefix is configurable via the mediaPrefix parameter.
// For example, if the host app routes /blog/media/* to this handler,
// pass "/blog/media/" as the prefix.
//
// Usage:
//
//	handler := blogstore.NewMediaHandler(store, "/blog/media/")
//	router.AddRoute(route.SetPath("/blog/media/*").SetHTMLHandler(handler.Handler))
type MediaHandler struct {
	store       StoreInterface
	mediaPrefix string
}

// NewMediaHandler creates a new MediaHandler for the given store.
// If mediaPrefix is empty, it defaults to "/blog/media/".
func NewMediaHandler(store StoreInterface, mediaPrefix string) *MediaHandler {
	if mediaPrefix == "" {
		mediaPrefix = defaultMediaPrefix
	}
	return &MediaHandler{store: store, mediaPrefix: mediaPrefix}
}

// MediaPrefix returns the configured media URL prefix.
func (h *MediaHandler) MediaPrefix() string {
	return h.mediaPrefix
}

// ServeURL returns the clean display URL for a media item: {prefix}{mediaId}.{ext}
// This can be used by admin UIs to show copyable URLs.
func (h *MediaHandler) ServeURL(mediaID, extension string) string {
	ext := strings.TrimPrefix(extension, ".")
	if ext != "" {
		ext = "." + ext
	}
	return h.mediaPrefix + mediaID + ext
}

// Handler serves the media content. It returns a string to be compatible
// with HTMLHandler-style route handlers.
func (h *MediaHandler) Handler(w http.ResponseWriter, r *http.Request) string {
	if h.store == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return "Blog store not configured"
	}

	mediaID := h.extractMediaID(r.URL.Path)
	if mediaID == "" {
		w.WriteHeader(http.StatusNotFound)
		return "Media not found"
	}

	media, err := h.store.MediaFindByID(context.Background(), mediaID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err.Error()
	}

	if media == nil {
		w.WriteHeader(http.StatusNotFound)
		return "Media not found"
	}

	if !media.IsActive() {
		w.WriteHeader(http.StatusNotFound)
		return "Media not found"
	}

	url := media.GetURL()

	// Data URI — decode base64 and serve directly
	if strings.HasPrefix(url, "data:") {
		return h.serveDataURI(w, r, url, media)
	}

	// HTTP(S) URL — redirect
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		http.Redirect(w, r, url, http.StatusFound)
		return ""
	}

	// File path — attempt to read and serve
	return h.serveFilePath(w, r, url, media)
}

func (h *MediaHandler) serveDataURI(w http.ResponseWriter, r *http.Request, dataURI string, media MediaInterface) string {
	content, err := decodeDataURL(dataURI)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return "Failed to decode media content"
	}

	contentType := media.GetType()
	if contentType == "" {
		contentType = mimeTypeFromExtension(media.GetExtension())
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", media.GetSize())
	w.Header().Set("Cache-Control", "public, max-age=31536000")
	w.Header().Set("ETag", media.GetID())

	if _, err := w.Write(content); err != nil {
		return "Failed to write media content: " + err.Error()
	}

	return ""
}

func (h *MediaHandler) serveFilePath(w http.ResponseWriter, r *http.Request, filePath string, media MediaInterface) string {
	contentType := media.GetType()
	if contentType == "" {
		contentType = mimeTypeFromExtension(media.GetExtension())
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Cache-Control", "public, max-age=31536000")
	w.Header().Set("ETag", media.GetID())

	http.ServeFile(w, r, filePath)
	return ""
}

// IsMediaURL checks if the given path matches this handler's media URL prefix.
func (h *MediaHandler) IsMediaURL(path string) bool {
	return strings.HasPrefix(path, h.mediaPrefix)
}

// extractMediaID parses the media ID from URL paths in either format:
//   - {prefix}<id>.<ext>           (e.g. /blog/media/sbm2hdy7x10.png)
//   - {prefix}<id>/<handle>.<ext>  (e.g. /blog/media/sbm2hdy7x10/image.png)
//
// The handle in the second form is cosmetic — lookup is always by ID.
func (h *MediaHandler) extractMediaID(urlPath string) string {
	path := strings.TrimPrefix(urlPath, h.mediaPrefix)
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")

	if idx := strings.Index(path, "/"); idx > 0 {
		return path[:idx]
	}

	if idx := strings.LastIndex(path, "."); idx >= 0 {
		return path[:idx]
	}
	return path
}

// decodeDataURL decodes a base64 data URL (e.g. "data:image/png;base64,...")
// and returns the raw binary content.
func decodeDataURL(dataURL string) ([]byte, error) {
	idx := strings.Index(dataURL, "base64,")
	if idx < 0 {
		return nil, errInvalidDataURL
	}
	b64data := dataURL[idx+7:]
	return base64.StdEncoding.DecodeString(b64data)
}

// mimeTypeFromExtension returns the MIME type for a given file extension.
func mimeTypeFromExtension(ext string) string {
	switch strings.ToLower(ext) {
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".svg":
		return "image/svg+xml"
	case ".ico":
		return "image/x-icon"
	case ".pdf":
		return "application/pdf"
	default:
		return "application/octet-stream"
	}
}

var errInvalidDataURL = errors.New("invalid data URL: missing base64 prefix")
