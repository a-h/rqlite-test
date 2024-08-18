# github.com/a-h/rqlite-test

## Tasks

### gomod2nix-update

```bash
gomod2nix
```

### build

```bash
nix build
```

### run

```bash
nix run
```

### develop

```bash
nix develop
```

### docker-build

```bash
nix build .#docker-image
```

### docker-load

Once you've built the image, you can load it into a local Docker daemon with `docker load`.

```bash
docker load < result
```

### docker-run

```bash
docker run -p 8080:8080 app:latest
```

### sql-vec-download

```bash
#TODO: Get this done inside Nix, and use the correct architecture and OS for the download.
wget https://github.com/asg017/sqlite-vec/releases/download/v0.1.1/sqlite-vec-0.1.1-loadable-macos-aarch64.tar.gz -O sqlite-vec.tar.gz
```

### db-run

```bash
rqlited -auth=auth.json -extensions-path=sqlite-vec.tar.gz ~/node.1
```

### db-migration-create

```bash
migrate create -ext sql -dir db/migrations -seq create_documents_table
```

### go-run

```bash
go run ./cmd/app
```

### curl-get-documents

```bash
curl http://localhost:8080/documents
```

### curl-post-document

```bash
curl -X POST -d '{"name": "name1","content":"Hello, World","embedding": [1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0]}' http://localhost:8080/documents
```

### curl-get-document

```bash
curl http://localhost:8080/document/1
```
