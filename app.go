package main

import (
	"log/slog"
	"net/http"

	"github.com/minio/minio-go/v7"
)

type app struct {
	client      *minio.Client
	bucketName  string
	defaultPage string
}

func (a *app) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slog.DebugContext(r.Context(), "Serving", "path", r.URL.Path)
	// If path contains a dot, we forward to the S3 bucket
	// Else we serve the index page
	if isStaticResource(r.URL.Path) {
		objectName := r.URL.Path[1:] // remove leading slash
		a.serve(w, r, objectName)
	} else {
		a.serve(w, r, a.defaultPage)
	}
}

func (a *app) serve(w http.ResponseWriter, r *http.Request, objectName string) {
	opts := minio.GetObjectOptions{}
	obj, err := a.client.GetObject(r.Context(), a.bucketName, objectName, opts)
	if err != nil {
		slog.ErrorContext(r.Context(), "Failed to get object", "object", objectName, "error", err)
		http.Error(w, "Failed to get object", http.StatusInternalServerError)
		return
	}
	defer obj.Close()
	slog.DebugContext(r.Context(), "Serving object", "object", objectName)
	stat, err := obj.Stat()
	if err != nil {
		slog.ErrorContext(r.Context(), "Failed to stat object", "object", objectName, "error", err)
		http.Error(w, "Object not found", http.StatusNotFound)
		return
	}
	http.ServeContent(w, r, objectName, stat.LastModified, obj)
}

func isStaticResource(path string) bool {
	if len(path) > 1 {
		for i := len(path) - 1; i >= 0; i-- {
			if path[i] == '/' {
				return false
			}
			if path[i] == '.' {
				return true
			}
		}
	}
	return false
}
