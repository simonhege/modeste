package main

import (
	"cmp"
	"context"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/simonhege/server"
)

func main() {
	ctx := context.Background()

	slog.SetDefault(slog.New(server.Wrap(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelInfo,
	}))))

	if err := godotenv.Load(); err != nil {
		slog.InfoContext(ctx, "no .env file loaded", "err", err)
	}

	client, err := minio.New(os.Getenv("MODESTE_ENDPOINT"), &minio.Options{
		Creds:  credentials.NewStaticV4(os.Getenv("MODESTE_ACCESS_KEY"), os.Getenv("MODESTE_SECRET_KEY"), ""),
		Secure: true,
	})
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create MinIO client", "error", err)
		return
	}

	a := &app{
		client:      client,
		bucketName:  os.Getenv("MODESTE_BUCKET_NAME"),
		defaultPage: cmp.Or(os.Getenv("MODESTE_DEFAULT_PAGE"), "index.html"),
	}
	slog.InfoContext(ctx, "minio client created", "bucket", a.bucketName, "endpoint", client.EndpointURL().String())

	s := server.New(false, server.RequestLogger)
	s.Handle("/", a)

	port := cmp.Or(os.Getenv("PORT"), "9022")
	address := ":" + port

	if err := s.Run(ctx, address); err != nil {
		slog.ErrorContext(ctx, "Server error", "error", err)
	}
}
