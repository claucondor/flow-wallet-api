package handlers

import (
	"net/http"
)

// swaggerUIHTML is the HTML template for Swagger UI
// It loads Swagger UI from CDN and points to our openapi.yml endpoint
const swaggerUIHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Flow Wallet API - Documentation</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5.9.0/swagger-ui.css" />
    <style>
        html { box-sizing: border-box; overflow: -moz-scrollbars-vertical; overflow-y: scroll; }
        *, *:before, *:after { box-sizing: inherit; }
        body { margin: 0; background: #fafafa; }
        .swagger-ui .topbar { display: none; }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5.9.0/swagger-ui-bundle.js"></script>
    <script src="https://unpkg.com/swagger-ui-dist@5.9.0/swagger-ui-standalone-preset.js"></script>
    <script>
        window.onload = function() {
            SwaggerUIBundle({
                url: window.location.origin + "/v1/openapi.yml",
                dom_id: '#swagger-ui',
                deepLinking: true,
                displayOperationId: true,
                defaultModelsExpandDepth: 1,
                defaultModelExpandDepth: 1,
                displayRequestDuration: true,
                docExpansion: 'list',
                filter: true,
                showExtensions: true,
                showCommonExtensions: true,
                tryItOutEnabled: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset
                ],
                plugins: [
                    SwaggerUIBundle.plugins.DownloadUrl
                ],
                layout: "StandaloneLayout"
            });
        };
    </script>
</body>
</html>`

// HandleSwaggerUI serves the Swagger UI HTML page
func HandleSwaggerUI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(swaggerUIHTML))
}

// HandleOpenAPISpec serves the OpenAPI specification YAML file
func HandleOpenAPISpec(openapiYAML []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-yaml")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		w.Write(openapiYAML)
	}
}
