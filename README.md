# MediaDispenser

## What is it for?

Primarily to meet my needs. I have various small applications built on top of S3. 
However, I was frustrated handling uploads in them. At the same time, 
direct upload to S3 isn't feasible in some cases. 
I also wanted to try building a project in Go again after a long time (I'm a newbie, so the code might reflect that).

Additionally, I wanted, optionally, to handle unique names, convert images to WebP, 
and maybe other tasks in the future.

That's why I created MediaDispenser, a small server as a single point for S3 with a REST interface. 
The original idea was for a significantly larger and more universal project, but maybe one day...

My main need was handling uploads, but issuing files through this server is also possible — whether public or via token.

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

