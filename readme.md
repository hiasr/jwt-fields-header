# JWT Fields Header Plugin

A Traefik middleware plugin that extracts fields from JWT tokens and creates custom headers from the claim values.

[![Build Status](https://github.com/traefik/plugindemo/workflows/Main/badge.svg?branch=master)](https://github.com/traefik/plugindemo/actions)

## What it does

This plugin extracts JWT claims from the Authorization Bearer token and creates a custom header by concatenating the specified claim values with '-' separators. This is useful for:

- Passing user information to downstream services via headers
- Creating custom identifiers from multiple JWT claims
- Enabling downstream services to access JWT data without parsing tokens
- Simplifying authentication data flow in microservices architectures

The plugin safely handles missing tokens, invalid JWTs, or missing claims by passing requests through unchanged.

## How it works

1. **Token Extraction**: Extracts the Bearer token from the `Authorization` header
2. **JWT Parsing**: Parses the JWT payload (without signature verification)
3. **Claim Extraction**: Extracts the specified claims from the JWT
4. **Header Creation**: Concatenates claim values with '-' and sets the custom header
5. **Graceful Handling**: Passes requests through unchanged if any step fails

## Configuration

The plugin requires two configuration parameters:

- `HeaderName`: The name of the header to set with the concatenated claim values
- `JwtClaims`: An array of JWT claim names to extract and concatenate

### Static Configuration

```yaml
experimental:
  plugins:
    jwt-fields-header:
      moduleName: github.com/hiasr/jwt-fields-header
      version: v0.1.0
```

### Dynamic Configuration

```yaml
# Dynamic configuration

http:
  routers:
    my-router:
      rule: host(`demo.localhost`)
      service: service-foo
      entryPoints:
        - web
      middlewares:
        - jwt-headers

  services:
   service-foo:
      loadBalancer:
        servers:
          - url: http://127.0.0.1:5000
  
  middlewares:
    jwt-headers:
      plugin:
        jwt-fields-header:
          HeaderName: X-User-Info
          JwtClaims:
            - sub
            - role
            - tenant
```

## Example

Given a JWT token with the following claims:
```json
{
  "sub": "user123",
  "role": "admin", 
  "tenant": "company-a",
  "iat": 1516239022
}
```

And the configuration:
```yaml
HeaderName: X-User-Info
JwtClaims:
  - sub
  - role
  - tenant
```

The plugin will create a header: `X-User-Info: user123-admin-company-a`

## Error Handling

The plugin gracefully handles various error conditions:

- **Missing Authorization header**: Request passes through unchanged
- **Invalid Bearer token format**: Request passes through unchanged  
- **Malformed JWT**: Request passes through unchanged
- **Missing claims**: Only available claims are included in the header
- **Non-string claim values**: Converted to JSON string representation
- **No claims configured**: Plugin is bypassed entirely

## Security Considerations

- The plugin does **not** verify JWT signatures - it only extracts claims
- JWT validation should be handled by other middleware or upstream services
- Sensitive claim data will be exposed in headers to downstream services
- Consider which claims to expose and configure accordingly

## Development

### Testing

```bash
go test
```

### Building

```bash
go mod tidy
go build
```