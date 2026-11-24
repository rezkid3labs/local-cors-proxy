# Local CORS proxy

A small development-only reverse proxy that adds CORS headers to an upstream API response. Do not expose it to the internet or use it as a production security boundary.

## Install and run

Install Go, then install the command:

```sh
go install github.com/rezkid3labs/local-cors-proxy@latest
```

Start the proxy by giving it the upstream API URL:

```sh
local-cors-proxy -target https://api.example.com
```

Or specify the origin

```sh
local-cors-proxy -target https://api.example.com -origin http://localhost:3000
```

It listens on `http://localhost:8080`. Requests to that address are forwarded to your upstream, preserving the path and query string. For example, `http://localhost:8080/api/health` is forwarded to `https://api.example.com/api/health`.

Alternatively, without installing, run directly from this project directory:

```sh
go run . -target https://api.example.com
```

The proxy listens on `http://localhost:8080`; requests retain their path and query string. The allowed frontend origin defaults to `http://localhost:5500`:

```sh
go run . -target https://api.example.com -origin http://localhost:3000
```

Serve `index.html` on the allowed origin (for example with VS Code Live Server on port 5500), then use its button to request `/` from the upstream through the proxy. Replace the path field to try another upstream endpoint. The example page requires a server; opening it directly as a `file://` URL will not match the configured origin.
