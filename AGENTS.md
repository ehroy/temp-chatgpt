# AGENTS.md

Panduan agent untuk proyek aplikasi OTP Reader berbasis:

- Backend: Golang
- Frontend: Svelte

## Tujuan Sistem

Membangun aplikasi untuk membaca OTP dari email Yahoo milik user sendiri dengan aturan:

- hanya membaca dari folder tertentu
- hanya memproses email yang masuk hari ini
- hanya mengambil OTP terbaru
- OTP hanya dianggap aktif maksimal 5 menit sejak email diterima
- pencarian OTP dilakukan berdasarkan input email user
- hasil OTP tidak boleh dibroadcast ke publik
- hasil OTP hanya boleh dikembalikan ke pemilik email atau session yang sah

---

# RULES UMUM

## Larangan

Agent DILARANG:

- membuat fitur broadcast OTP ke publik
- membuat endpoint tanpa autentikasi untuk menampilkan OTP
- menampilkan seluruh isi inbox ke frontend
- menyimpan OTP lama tanpa expiry
- mengembalikan OTP dari email yang bukan milik user yang meminta
- mengambil OTP dari semua folder tanpa whitelist
- menggunakan data email tanpa validasi waktu masuk

## Wajib

Agent WAJIB:

- hanya membaca folder yang diizinkan
- hanya mengambil email hari ini
- hanya mengambil OTP terbaru
- menolak OTP yang lebih lama dari 5 menit
- memvalidasi email input user
- memisahkan logika backend dan frontend dengan jelas
- menambahkan logging aman tanpa membocorkan OTP penuh
- menyusun struktur kode rapi, modular, dan mudah diuji

---

# FLOW BISNIS

## Flow utama

1. User membuka halaman pencarian OTP
2. User memasukkan alamat email
3. Frontend mengirim request ke backend
4. Backend memvalidasi format email
5. Backend mencari email OTP hanya pada folder yang diizinkan
6. Backend memfilter:
   - tanggal harus hari ini
   - subject / sender sesuai pola OTP
   - email paling baru
   - usia email maksimal 5 menit
7. Backend mengekstrak kode OTP
8. Backend mengembalikan response ringkas ke frontend
9. Frontend menampilkan status:
   - ditemukan
   - tidak ditemukan
   - kadaluarsa
   - email tidak valid

## Jawaban atas pertanyaan flow

Ya, flow yang tepat adalah **by input email untuk pencarian OTP**.

Alasan:

- lebih aman
- lebih jelas pemilik OTP-nya
- memudahkan filtering data
- menghindari menampilkan OTP milik orang lain
- cocok untuk lookup OTP terbaru per email

---

# STRUKTUR PROYEK YANG DISARANKAN

## Backend Golang

Gunakan struktur seperti ini:

```txt
backend/
  cmd/
    api/
      main.go
  internal/
    config/
      config.go
    handler/
      otp_handler.go
      health_handler.go
    service/
      otp_service.go
      yahoo_service.go
    repository/
      otp_repository.go
    model/
      otp_result.go
      email_message.go
    middleware/
      auth.go
      logging.go
      cors.go
    utils/
      time.go
      regex.go
      response.go
  routes/
    routes.go
  go.mod

## frontend svelte

  frontend/
  src/
    lib/
      api/
        otp.ts
      types/
        otp.ts
      utils/
        formatTime.ts
    routes/
      +page.svelte
      otp/
        +page.svelte
    components/
      OtpSearchForm.svelte
      OtpResultCard.svelte
      StatusAlert.svelte
  package.json
```
