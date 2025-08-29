# Swagger UI

This directory contains the static files for Swagger UI, used to visualize and interact with the API Gateway's OpenAPI specification.

## Usage

- `index.html` is configured to load the OpenAPI spec from `/swagger/openapi.yaml`.
- To update Swagger UI, you can:
  - Replace or edit `index.html` as needed.
  - Update the OpenAPI spec at `/swagger/openapi.yaml` to reflect API changes.

## Serving

- The Go server serves the Swagger UI at [`/swagger/`](http://localhost:PORT/swagger/).
- The OpenAPI spec is available at [`/swagger/openapi.yaml`](http://localhost:PORT/swagger/openapi.yaml).

## Updating Swagger UI

To upgrade Swagger UI to a newer version:

1. Download the latest release from [Swagger UI GitHub Releases](https://github.com/swagger-api/swagger-ui/releases).
2. Replace the contents of this directory with the new static files.
3. Ensure `index.html` is configured to load `/swagger/openapi.yaml`.

## Notes

- Make sure the OpenAPI spec is kept up to date with your API changes.
- For local development, access the Swagger UI at `http://localhost:PORT/swagger/`.
