# Swagger UI

This directory contains the Swagger UI static files for the API Gateway.

- `index.html` loads the OpenAPI spec from `/swagger/openapi.yaml`.
- To update Swagger UI, replace or edit `index.html` as needed.

The Go server serves this UI at `/swagger/` and the OpenAPI spec at `/swagger/openapi.yaml`.
