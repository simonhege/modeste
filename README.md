# Modeste

modeste is a proxy for single-page application on top of any storage.

## Features
- Proxy requests to object storage (e.g., MinIO, S3)
- Serve index.html for SPA routes


## Configuration
The application can be configured via environment variables:
- `MODESTE_ENDPOINT`: The endpoint of the object storage (e.g., `play.min.io:9000`)
- `MODESTE_ACCESS_KEY`: Access key for the object storage
- `MODESTE_SECRET_KEY`: Secret key for the object storage
- `MODESTE_BUCKET_NAME`: The bucket name where the SPA files are stored
- `MODESTE_DEFAULT_PAGE`: The default page to serve for SPA routes (e.g., `index.html`)



## Local development
```bash
go install
```

Create a `.env` file in the root directory with the following content:
```env
```

Then run the application with:
```bash
modeste
```

## Deployment
TODO: detail how to deploy

## TODO
- Add caching for default file
- Add caching for frequently accessed static files
- Add TLS support for secure connections
- Support for custom error pages
- Add rate limiting and security features