# Dokumentasi API Perpustakaan

Dokumentasi ini menjelaskan semua endpoint API yang tersedia pada aplikasi sistem perpustakaan berbasis Go Fiber.

## Persiapan (Prerequisites)

Sebelum menjalankan aplikasi, pastikan Anda telah menginstal:
- **Go** (versi 1.25.5 atau yang lebih baru)
- **PostgreSQL** (sebagai database utama)
- **Git**

## Instalasi & Cara Menjalankan

1. **Clone repositori ini:**
   ```bash
   git clone <url-repo-anda>
   cd gdgoc-backend
   ```

2. **Siapkan Database:**
   - Pastikan layanan PostgreSQL sudah berjalan.
   - Buat database baru, misalnya `backend-gdgoc`.
   - Jalankan query yang ada di dalam file `schema.sql` ke database tersebut untuk membuat tabel-tabel yang diperlukan (`users`, `books`, `borrowings`).

3. **Konfigurasi Environment:**
   Buat file `.env` di direktori utama proyek (jika belum ada) dan sesuaikan dengan konfigurasi database Anda. Contoh isi `.env`:
   ```env
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASS=password_anda
   DB_NAME=backend-gdgoc
   JWT_SECRET=super-secret-key
   ```

4. **Unduh Dependensi:**
   Jalankan perintah berikut untuk mengunduh semua modul Go yang dibutuhkan:
   ```bash
   go mod tidy
   ```

5. **Jalankan Aplikasi:**
   Mulai server backend dengan perintah:
   ```bash
   go run main.go
   ```
   Aplikasi akan berjalan secara default di `http://localhost:8080`.

## Struktur Dokumentasi

Dokumentasi dibagi menjadi beberapa bagian:

1. **[Autentikasi & Pengguna (Auth & Users)]**
   - Pendaftaran pengguna baru (`POST /users`)
   - Login untuk mendapatkan token JWT (`POST /login`)

2. **[Manajemen Buku (Books)]**
   - Mendapatkan daftar buku (`GET /books`)
   - Mendapatkan detail buku (`GET /books/:id`)
   - Menambahkan buku baru (`POST /books` - *Protected*)
   - Memperbarui data buku (`PUT /books/:id` - *Protected*)
   - Menghapus buku (`DELETE /books/:id` - *Protected*)

3. **[Peminjaman & Pengembalian Buku (Borrowings)]**
   - Meminjam buku (`POST /borrowings/borrow` - *Protected*)
   - Mengembalikan buku (`POST /borrowings/return/:id` - *Protected*)
   - Melihat daftar peminjaman saya (`GET /borrowings/user` - *Protected*)
   - Melihat semua riwayat peminjaman (`GET /borrowings/` - *Protected*)

---

## Informasi Umum

- **Base URL**: `http://localhost:8080`
- **Format Respons**: Semua respons menggunakan format JSON.
- **Autentikasi**: Endpoint dengan label **Protected** memerlukan header Authorization bertipe Bearer Token JWT.
  ```http
  Authorization: Bearer <your_jwt_token>
  ```
