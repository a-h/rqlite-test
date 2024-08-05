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

### db-run

```bash
rqlited -auth=auth.json ~/node.1
```

### db-migration-create

```bash
migrate create -ext sql -dir db/migrations -seq create_documents_table
```

### sqlc-generate

```bash
sqlc generate
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
curl -X POST -d '{"name": "name1","content":"Hello, World"}' http://localhost:8080/documents
```

### curl-get-document

```bash
curl http://localhost:8080/document/1
```
