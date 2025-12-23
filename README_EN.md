# Backend Service Documentation

This document is intended for the **IT client team** to easily understand the architecture, operation, and deployment process of the backend you have built.

---

## 1) Project Overview
- **Language/Framework**: Go (backend)
- **Orchestration**: Docker Compose
- **API Documentation**: Swagger (`<domain>/swagger/index.html`)
- **Supporting Services**:
  - PostgreSQL
  - Redis
  - MinIO
  - Mailhog
  - PgHero

> **Note on credentials & ports**: All details are stored in the `.env` file.

---

## 2) Swagger Link
```
<domain>/swagger/index.html
```
Replace `<domain>` with your server domain or IP address.

---

## 3) Service Architecture
```
[Client] --> [Backend (Go)] --> [PostgreSQL]
                       |--> [Redis]
                       |--> [MinIO]
                       |--> [Mailhog]
                       |--> [PgHero]
```

---

## 4) Ports & Credentials
All values are taken from `.env`:
- Backend: `BACKEND_PORT=____`
- PostgreSQL: `POSTGRES_PORT=____`, etc.
- Redis, MinIO, Mailhog, PgHero → fill in according to `.env`.

---

## 5) Server Directory Structure
Since this is a **built artifact**, the server directory contains only:
```
/
├── Dockerfile
├── docker-compose.yml
├── logs/            # folder for logs
└── app              # built binary
```

---

## 6) Building the Application Locally
Before uploading to the server, build the Go binary locally:
```bash
GOOS=linux GOARCH=amd64 go build -o app ./cmd/main.go
```
- `GOOS=linux`: target OS Linux
- `GOARCH=amd64`: target architecture 64-bit
- `-o app`: output binary name
- `./cmd/main.go`: application entry point

---

## 7) Operational Commands (Docker Compose)
- When `.env` changes:
  ```bash
  docker compose up -d --force-recreate --no-deps backend
  ```
- When source code changes:
  ```bash
  docker compose up -d --build --no-deps backend
  ```
- When `docker-compose.yml` changes:
  ```bash
  docker compose up -d
  ```
- When Dockerfile changes:
  ```bash
  docker compose build --no-cache backend
  docker compose up -d --no-deps backend
  ```
- View logs:
  ```bash
  docker compose logs -f backend
  ```
- Check container status:
  ```bash
  docker compose ps
  ```

---

## 8) Health Check
- Backend: `http://<domain>:<port>/api/ping`
- Swagger: `<domain>/swagger/index.html`
- PostgreSQL: test connection
- Redis: `redis-cli ping`
- MinIO: access web UI
- Mailhog: access UI
- PgHero: access UI

---

## 9) Troubleshooting
- Container restart loop → check logs
- Swagger not accessible → verify port & container status
- Database connection failed → check `.env` & service status

---

## 10) Backup & Restore
Example for PostgreSQL:
```bash
docker exec -t <postgres-container> pg_dump -U $POSTGRES_USER -d $POSTGRES_DB > backup.sql
```

---

## 11) Security
- Do not commit `.env` to public repositories
- Rotate credentials regularly
- Restrict port access to trusted networks

---

## 12) Go-Live Checklist
- [ ] `.env` file is complete
- [ ] All services are **Up**
- [ ] Swagger is accessible
- [ ] Firewall & DNS configured

---

## 13) Step-by-Step Deployment Process
### Step 1: Build the Application Locally

Before uploading to the server, build the Go binary locally:
```bash
GOOS=linux GOARCH=amd64 go build -o app ./cmd/main.go
```

- `GOOS=linux`: file is complete
 target OS Linux
- `GOARCH=amd64`: target architecture 64-bit
- `-o app`: output binary name
- `./cmd/main.go`: application entry point 

#### Wait until the build process completes.

### Step 2: Copy Binary to Server via SFTP

1. Connect to the server using your SFTP client
2. Navigate to the backend folder on the server 
3. Optional safety step: Rename the existing app file (e.g., app.old) if present 
4. Upload the newly built app file (drag and drop from your local folder)
5. Ensure the file is placed in the correct server directory
   
#### Wait for the file transfer to complete.

### Step 3: Update Docker Container on Server

1. SSH into the server using terminal 
2. Navigate to the backend directory containing docker-compose.yml 
3. Run the update command:

``` bash
docker compose up -d --build --no-deps backend
```
This command rebuilds only the backend container without affecting other services.

#### Wait for the container to rebuild and restart.

### Step 4: Verify Deployment

Check the application logs to confirm successful deployment:

```bash
docker compose logs -f backend
```
Look for any errors and verify the application

---
