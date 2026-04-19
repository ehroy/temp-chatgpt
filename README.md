# OTP Reader

Scaffold aplikasi OTP reader sesuai `AGENTS.md`.

## Backend

```bash
cd backend
go test ./...
go run ./cmd/api
```

Env utama:

- `APP_ADDR` default `127.0.0.1:9001`
- `APP_AUTH_TOKEN` default `dev-token`
- `YAHOO_ALLOWED_FOLDERS` default `INBOX,OTP`
- `OTP_MAX_AGE_MINUTES` default `5`
- `YAHOO_IMAP_HOST` default `imap.mail.yahoo.com`
- `YAHOO_IMAP_PORT` default `993`
- `YAHOO_EMAIL` email Yahoo pemilik mailbox
- `YAHOO_APP_PASSWORD` app password Yahoo
- `YAHOO_MAX_SCAN` jumlah pesan terakhir yang dipindai per folder

## Frontend

```bash
cd frontend
npm install
npm run dev
```

Dev server frontend berjalan di port `4174`.

Env utama:

- `VITE_API_BASE_URL` default `http://localhost:9001`
- `VITE_API_TOKEN` default `dev-token`

Catatan:

- `APP_AUTH_TOKEN` hanya untuk mengamankan endpoint HTTP backend
- akses mailbox langsung memakai credential IMAP Yahoo di backend
