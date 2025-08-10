# MediaDispenser

## What is it for?

A lightweight S3 gateway with a REST API for uploading, serving, and processing media files.

It sits between your apps and S3, taking care of uploads, downloads, and optional processing like automatic file renaming or converting images to WebP. Access can be public or token-based, and extra logic can be added through plugins.

Designed as a single, easy-to-deploy service, it makes S3 integration simpler, more secure, and more flexible. Originally built to solve my own upload frustrations, it’s flexible enough to work as a drop-in media gateway for almost any project.

## Build

Go version: 1.23.2

```
go build
./MediaDispenserGo
```

## Konfigurace

```yaml
# List of tokens
tokens:
  # Token with full access
  - name: full-perms
    token: token-string
  # Token with limited access to the /pics prefix only
  - name: pics
    token: pics-token
    isolation: pics  # Uploaded objects will always be within this prefix; objects must start with this prefix for reading
  # This token is prohibited from all operations
  - name: useless
    token: useless
    disallow:
      - list
      - get
      - upload
      - delete
# S3 credentials
s3:
  access_key: access
  secret_key: secret
  bucket: media-dispenser
  host: host
  secure: true
# Object dispensing mode
dispensing: public
dispensingRoot: none  # The root won't be accessible
# We can adjust dispensing for specific prefixes
dispensingPrefix:
  # The `profiles` prefix will only be accessible with a token
  - prefix: profiles
    mode: private
  # Since rules are processed top to bottom, the previous rule matches, and this one won't be applied
  - prefix: profiles/abc
    mode: public
# Plugins that may modify upload/get actions
plugins:
  - WebPConverter
  - KSUIDRenamer
```

## Configuration

```bash
# Upload
curl -X POST "http://localhost:8080/u \
  -H "Authorization: pics"
  -F "file=@/path/to/file.png" 
# The response is the path relative to the dispenser: `pics/<ksuid>.webp`

# Get
curl "http://localhost:8080/g/pics/<ksuid>.webp

# Delete
curl -X DELETE "http://localhost:8080/g/pics/<ksuid>.webp" \
  -H "Authorization: pics"
  
# List
# -- not implemented yet --
```

## Plugins

| Název           | Operace   | Popis                                           |
|-----------------|-----------|-------------------------------------------------|
| `KSUIDRenamer`  | `upload`  | Replaces the names of uploaded files with KSUID |
| `WebPConverter` | `upload`  | Converts `image/jpeg` and `image/png` to WebP.  |

