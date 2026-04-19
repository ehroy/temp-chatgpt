# Dokumentasi Deploy

Dokumen ini menjelaskan cara deploy aplikasi OTP Reader dengan:

- backend Golang dijalankan lewat `systemd`
- Nginx sebagai reverse proxy
- backend bind ke port lokal lain, misalnya `127.0.0.1:9001`
- domain utama `babiguling.my.id` dan opsi subdomain API
- SSL/TLS dengan Let’s Encrypt

## Arsitektur

Rekomendasi struktur produksi:

- `babiguling.my.id` untuk frontend
- `api.babiguling.my.id` untuk backend, atau backend diproxy lewat `/api`
- backend tidak dibuka ke publik langsung
- backend hanya listen di localhost

Contoh alur:

`browser -> nginx -> backend systemd -> Yahoo IMAP`

## Komponen yang Dibutuhkan

- Ubuntu/Debian server
- `nginx`
- `golang` versi stabil
- `nodejs` dan `npm` jika frontend ikut dibuild di server
- `certbot` untuk SSL

## 1. Persiapan Server

Update package:

```bash
sudo apt update
sudo apt upgrade -y
```

Install dependency:

```bash
sudo apt install -y nginx git curl
```

Install Go dan Node.js sesuai kebutuhan proyek.

## 2. Struktur Folder Deploy

Contoh lokasi deploy:

```txt
/var/www/email-chatgpt/
  backend/
  frontend/
  logs/
```

Contoh clone repo:

```bash
sudo mkdir -p /var/www/email-chatgpt
sudo chown -R $USER:$USER /var/www/email-chatgpt
git clone <REPO_URL> /var/www/email-chatgpt
```

## 3. Konfigurasi Backend

File env backend ada di `backend/.env.example`. Untuk production, buat file:

`/var/www/email-chatgpt/backend/.env`

Contoh:

```env
APP_ADDR=127.0.0.1:9001
APP_AUTH_TOKEN=ganti-token-yang-kuat
YAHOO_ALLOWED_FOLDERS=OTP,INBOX
OTP_MAX_AGE_MINUTES=5
YAHOO_IMAP_HOST=imap.mail.yahoo.com
YAHOO_IMAP_PORT=993
YAHOO_EMAIL=your.yahoo@email.com
YAHOO_APP_PASSWORD=app-password-yahoo
YAHOO_MAX_SCAN=25
```

### Catatan env

- `APP_ADDR` harus diarahkan ke localhost, bukan port publik
- `APP_AUTH_TOKEN` dipakai untuk mengamankan endpoint backend
- `YAHOO_ALLOWED_FOLDERS` hanya folder whitelist
- `OTP_MAX_AGE_MINUTES` sebaiknya tetap `5`
- `YAHOO_EMAIL` dan `YAHOO_APP_PASSWORD` wajib diisi

### Build backend

Masuk ke folder backend lalu build binary:

```bash
cd /var/www/email-chatgpt/backend
go mod download
go build -o emailchatgpt ./cmd/api
```

Hasil binary:

```txt
/var/www/email-chatgpt/backend/emailchatgpt
```

## 4. Systemd Service Backend

Buat file service:

`/etc/systemd/system/emailchatgpt-backend.service`

Contoh isi:

```ini
[Unit]
Description=Email ChatGPT Backend
After=network.target

[Service]
Type=simple
WorkingDirectory=/var/www/email-chatgpt/backend
ExecStart=/var/www/email-chatgpt/backend/emailchatgpt
Restart=always
RestartSec=5
User=www-data
Group=www-data
EnvironmentFile=/var/www/email-chatgpt/backend/.env
NoNewPrivileges=true
PrivateTmp=true

[Install]
WantedBy=multi-user.target
```

### Aktifkan service

```bash
sudo systemctl daemon-reload
sudo systemctl enable emailchatgpt-backend
sudo systemctl start emailchatgpt-backend
sudo systemctl status emailchatgpt-backend
```

### Cek log backend

```bash
sudo journalctl -u emailchatgpt-backend -f
```

## 5. Konfigurasi Nginx

Ada 2 pola umum.

### Opsi A: 1 domain utama, backend diproxy lewat `/api`

Frontend diakses dari `https://babiguling.my.id`, backend diakses lewat `https://babiguling.my.id/api/...`.

Untuk development lokal, frontend Vite berjalan di port `4174`.

Contoh konfigurasi Nginx:

`/etc/nginx/sites-available/emailchatgpt`

```nginx
server {
    listen 80;
    server_name babiguling.my.id www.babiguling.my.id;

    root /var/www/email-chatgpt/frontend/build;
    index index.html;

    location /api/ {
        proxy_pass http://127.0.0.1:9001/;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    location /health {
        proxy_pass http://127.0.0.1:9001/health;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    location / {
        try_files $uri $uri/ /index.html;
    }
}
```

Catatan:

- gunakan `root` hanya jika frontend sudah dibuild menjadi file statis
- jika frontend memakai server Node sendiri, ganti `location /` menjadi proxy ke port frontend

### Opsi B: backend di subdomain khusus

Contoh:

- frontend: `https://babiguling.my.id`
- backend: `https://api.babiguling.my.id`

Contoh Nginx backend:

```nginx
server {
    listen 80;
    server_name api.babiguling.my.id;

    location / {
        proxy_pass http://127.0.0.1:9001;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

## 6. Aktifkan Site Nginx

```bash
sudo ln -s /etc/nginx/sites-available/emailchatgpt /etc/nginx/sites-enabled/emailchatgpt
sudo nginx -t
sudo systemctl reload nginx
```

## 7. SSL dengan Let’s Encrypt

Install certbot:

```bash
sudo apt install -y certbot python3-certbot-nginx
```

Jika pakai domain utama:

```bash
sudo certbot --nginx -d babiguling.my.id -d www.babiguling.my.id
```

Jika pakai subdomain API:

```bash
sudo certbot --nginx -d babiguling.my.id -d www.babiguling.my.id -d api.babiguling.my.id
```

Pastikan redirect HTTP ke HTTPS aktif.

## 8. Firewall

Buka port yang dibutuhkan saja:

```bash
sudo ufw allow OpenSSH
sudo ufw allow 80
sudo ufw allow 443
sudo ufw enable
```

Port backend `9001` tidak perlu dibuka karena hanya untuk localhost.

## 9. Konfigurasi Domain

### DNS

Tambahkan record DNS:

- `A babiguling.my.id -> IP_SERVER`
- `A www.babiguling.my.id -> IP_SERVER` jika dipakai
- `A api.babiguling.my.id -> IP_SERVER` jika memakai subdomain API

### Contoh skema domain

#### Skema 1

- `babiguling.my.id` = frontend
- `babiguling.my.id/api` = backend

#### Skema 2

- `babiguling.my.id` = frontend
- `api.babiguling.my.id` = backend

Skema 1 lebih sederhana jika frontend dan backend berada pada server yang sama.

## 10. Endpoint Backend

Endpoint yang tersedia:

- `GET /health`
- `POST /api/otp/lookup`

Header yang wajib:

```http
Authorization: Bearer <APP_AUTH_TOKEN>
Content-Type: application/json
```

Body lookup OTP:

```json
{
  "email": "user@yahoo.com"
}
```

## 11. Frontend Env Produksi

Di frontend, set env production sesuai domain backend.

Contoh jika backend diproxy dari domain yang sama:

```env
VITE_API_BASE_URL=https://babiguling.my.id
VITE_API_TOKEN=ganti-token-yang-sama-dengan-backend
```

Contoh jika backend pakai subdomain:

```env
VITE_API_BASE_URL=https://api.babiguling.my.id
VITE_API_TOKEN=ganti-token-yang-sama-dengan-backend
```

## 12. Verifikasi Setelah Deploy

### Cek backend

```bash
curl -i http://127.0.0.1:9001/health
```

### Cek lewat Nginx

```bash
curl -i https://babiguling.my.id/health
```

### Cek lookup OTP

```bash
curl -X POST https://babiguling.my.id/api/otp/lookup \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer ganti-token-yang-sama-dengan-backend' \
  -d '{"email":"user@yahoo.com"}'
```

## 13. Troubleshooting

### Service backend tidak jalan

```bash
sudo systemctl status emailchatgpt-backend
sudo journalctl -u emailchatgpt-backend -n 100 --no-pager
```

### Nginx error 502

Penyebab umum:

- backend belum jalan
- `APP_ADDR` tidak sesuai
- Nginx proxy ke port yang salah

### SSL gagal terpasang

Pastikan:

- DNS sudah mengarah ke IP server
- port 80 dan 443 terbuka
- tidak ada server block Nginx yang bentrok

## 14. Update Versi Aplikasi

Langkah update:

```bash
cd /var/www/email-chatgpt
git pull
cd backend
go build -o emailchatgpt ./cmd/api
sudo systemctl restart emailchatgpt-backend
```

Jika frontend berubah dan dibuild statis:

```bash
cd /var/www/email-chatgpt/frontend
npm install
npm run build
sudo systemctl reload nginx
```

## 15. Catatan Keamanan

- jangan expose backend langsung ke publik
- jangan simpan token di repo
- gunakan `APP_AUTH_TOKEN` yang kuat
- pastikan folder Yahoo yang dibaca hanya whitelist
- tetap pertahankan limit usia OTP 5 menit
