# Frontend Kedua di VPS Yang Sama

Panduan ini dipakai kalau ada project frontend lain dengan tema berbeda di VPS yang sama.

## Prinsip

- jangan pakai port yang sama dengan project utama
- simpan folder terpisah per project
- domain/subdomain terpisah lebih aman
- backend tetap satu dan dipakai bersama bila memang sama

## Contoh Port

- frontend utama dev: `4174`
- frontend kedua dev: `4176`
- backend: `127.0.0.1:9001`

## Contoh Struktur

```txt
/var/www/babiguling-main/
/var/www/babiguling-theme2/
```

## Contoh Build Frontend Kedua

```bash
cd /var/www/babiguling-theme2/frontend
npm install
npm run build
```

## Contoh Nginx Frontend Kedua

Jika frontend kedua memakai subdomain, misalnya `theme2.babiguling.my.id`:

```nginx
server {
    listen 80;
    server_name theme2.babiguling.my.id;

    return 301 https://$host$request_uri;
}

server {
    listen 443 ssl http2;
    server_name theme2.babiguling.my.id;

    ssl_certificate /etc/letsencrypt/live/theme2.babiguling.my.id/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/theme2.babiguling.my.id/privkey.pem;

    root /var/www/babiguling-theme2/frontend/build;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }
}
```

## Catatan

- kalau frontend kedua tetap perlu akses backend, set `VITE_API_BASE_URL` ke `https://api.babiguling.my.id`
- jangan gabungkan asset build dari dua project ke folder yang sama
- kalau pakai `npm run dev`, jalankan di port berbeda supaya tidak bentrok
