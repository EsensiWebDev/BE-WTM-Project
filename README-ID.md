# Dokumentasi Backend Layanan

Dokumen ini ditujukan untuk tim **IT client** agar mudah memahami arsitektur, cara menjalankan, serta proses operasi backend yang Anda bangun.

---

## 1) Ringkasan Proyek
- **Bahasa/Framework**: Go (backend)
- **Orkestrasi**: Docker Compose
- **Dokumentasi API**: Swagger (`<domain>/swagger/index.html`)
- **Service pendukung**:
    - PostgreSQL
    - Redis
    - MinIO
    - Mailhog
    - PgHero

> **Catatan kredensial & port**: Semua ada di file `.env`.

---

## 2) Link Swagger
```
<domain>/swagger/index.html
```
Ganti `<domain>` dengan domain atau IP server.

---

## 3) Arsitektur Layanan
```
[Client] --> [Backend (Go)] --> [PostgreSQL]
                       |--> [Redis]
                       |--> [MinIO]
                       |--> [Mailhog]
                       |--> [PgHero]
```

---

## 4) Port & Kredensial
Semua diambil dari `.env`:
- Backend: `BACKEND_PORT=____`
- PostgreSQL: `POSTGRES_PORT=____`, dll.
- Redis, MinIO, Mailhog, PgHero → isi sesuai `.env`.

---

## 5) Struktur Direktori di Server
Karena ini **hasil build**, struktur di server hanya:
```
/
├── Dockerfile
├── docker-compose.yml
├── logs/            # folder untuk log
└── app              # binary hasil build
```

---

## 6) Proses Build Aplikasi (di Lokal)
Sebelum upload ke server, build binary Go:
```bash
GOOS=linux GOARCH=amd64 go build -o app ./cmd/main.go
```
- `GOOS=linux`: target OS Linux
- `GOARCH=amd64`: target arsitektur 64-bit
- `-o app`: nama output binary
- `./cmd/main.go`: entry point aplikasi

---

## 7) Perintah Operasional (Docker Compose)
- Perubahan `.env`:
  ```bash
  docker compose up -d --force-recreate --no-deps backend
  ```
- Perubahan kode:
  ```bash
  docker compose up -d --build --no-deps backend
  ```
- Perubahan `docker-compose.yml`:
  ```bash
  docker compose up -d
  ```
- Perubahan Dockerfile:
  ```bash
  docker compose build --no-cache backend
  docker compose up -d --no-deps backend
  ```
- Logs:
  ```bash
  docker compose logs -f backend
  ```
- Status:
  ```bash
  docker compose ps
  ```

---

## 8) Health Check
- Backend: `http://<domain>:<port>/api/ping`
- Swagger: `<domain>/swagger/index.html`
- PostgreSQL: tes koneksi
- Redis: `redis-cli ping`
- MinIO: akses web UI
- Mailhog: akses UI
- PgHero: akses UI

---

## 9) Troubleshooting
- Container restart loop → cek logs
- Swagger tidak tampil → cek port & container
- DB gagal → cek `.env` & service status

---

## 10) Backup & Restore
Contoh PostgreSQL:
```bash
docker exec -t <postgres-container> pg_dump -U $POSTGRES_USER -d $POSTGRES_DB > backup.sql
```

---

## 11) Keamanan
- Jangan commit `.env`
- Rotasi kredensial
- Batasi akses port

---

## 12) Checklist Go-Live
- [ ] `.env` lengkap
- [ ] Semua service **Up**
- [ ] Swagger dapat diakses
- [ ] Firewall & DNS OK

## 13) Proses Deployment Langkah demi Langkah

### Langkah 1: Bangun Aplikasi Secara Lokal

Sebelum mengunggah ke server, bangun binary Go secara lokal:

```bash
GOOS=linux GOARCH=amd64 go build -o app ./cmd/main.go
```
- `GOOS=linux`: target OS Linux
- `GOARCH=amd64`: target arsitektur 64-bit
- `-o app`: nama output binary
- `./cmd/main.go`: titik masuk aplikasi
#### Tunggu hingga proses build selesai.

### Langkah 2: Salin Binary ke Server via SFTP

1. Hubungkan ke server menggunakan klien SFTP Anda 
2. Navigasi ke folder backend di server 
3. Langkah keamanan opsional: Ganti nama file app yang ada (misalnya, app.old) jika ada 
4. Unggah file app yang baru dibangun (drag and drop dari folder lokal Anda)
5. Pastikan file ditempatkan di direktori server yang benar
#### Tunggu hingga transfer file selesai.

### Langkah 3: Perbarui Container Docker di Server

1. SSH ke server menggunakan terminal 
2. Navigasi ke direktori backend yang berisi docker-compose.yml 
3. Jalankan perintah pembaruan:
```bash
docker compose up -d --build --no-deps backend
```
Perintah ini hanya membangun ulang container backend tanpa mempengaruhi layanan lainnya.

#### Tunggu container untuk dibangun ulang dan restart.

### Langkah 4: Verifikasi Deployment

Periksa log aplikasi untuk memastikan deployment berhasil:

```bash
docker compose logs -f backend
```

Cari error apa pun dan verifikasi aplikasi dimulai dengan benar.

---
