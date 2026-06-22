# Penugasan: Membuat Endpoint GET /users (Protected)

## 📌 Deskripsi Penugasan

Dalam tugas ini, Anda diminta untuk menambahkan fitur baru berupa endpoint **`GET /users`** pada REST API sistem perpustakaan ini. Endpoint ini digunakan untuk melihat daftar semua pengguna yang terdaftar di dalam sistem perpustakaan.

Karena data pengguna bersifat sensitif, endpoint ini **wajib dilindungi (protected)** dengan menggunakan JSON Web Token (JWT) melalui middleware autentikasi yang sudah disediakan.

---

## 🛠️ Persyaratan Teknis

Berikut adalah spesifikasi endpoint yang harus diimplementasikan:

- **Method**: `GET`
- **URL**: `/users`
- **Headers**:
  - `Authorization: Bearer <your_jwt_token>` (Wajib diisi dengan token JWT yang valid)
- **Response**:
  - **`200 OK` (Sukses)**: Mengembalikan list/array data pengguna (password wajib disembunyikan/tidak dikirim kembali).
  - **`401 Unauthorized` (Autentikasi gagal)**: Jika token tidak dikirim, format token salah, atau token tidak valid.
  - **`500 Internal Server Error` (Gagal diproses)**: Jika terjadi kesalahan pada database atau sistem.

---

## 📋 Contoh Output Response (JSON)

### 1. **200 OK (Sukses)**

```json
[
  {
    "id": 1,
    "username": "nouval",
    "name": "Nouval Arianta",
    "created_at": "2026-06-20T15:20:00Z",
    "updated_at": "2026-06-20T15:20:00Z"
  },
  {
    "id": 2,
    "username": "budi",
    "name": "Budi Santoso",
    "created_at": "2026-06-21T09:15:00Z",
    "updated_at": "2026-06-21T09:15:00Z"
  }
]
```

_(Catatan: Field password otomatis disembunyikan jika Anda menggunakan struct `domain.User` yang sudah dikonfigurasi dengan tag `json:"-"`)_.

### 2. **401 Unauthorized (Tanpa Token / Token Tidak Valid)**

```json
{
  "error": "Authorization header is required"
}
```

atau

```json
{
  "error": "Invalid token"
}
```

### 3. **500 Internal Server Error (Gagal Diproses)**

```json
{
  "error": "database connection error"
}
```

---

## 🚶‍♂️ Langkah-Langkah Panduan Pengerjaan

Untuk menyelesaikan penugasan ini, ikuti langkah-langkah terstruktur berikut:

### Langkah 1: Modifikasi Inisialisasi User Handler

Buka file [user_handler.go] (delivery/http/user_handler.go).

1. Ubah fungsi `NewUserHandler` agar dapat menerima parameter `jwtSecret string` tambahan untuk inisialisasi middleware JWT.
2. Impor middleware JWT dari package `delivery/http/middleware`.
3. Daftarkan route baru `GET /users` dengan menyisipkan middleware JWT di depannya.

_Petunjuk kode:_

```go
func NewUserHandler(app *fiber.App, us domain.UserUsecase, jwtSecret string) {
	handler := &UserHandler{
		UUsecase: us,
	}

	authMiddleware := middleware.JWTAuthMiddleware(jwtSecret)

	app.Post("/login", handler.Login)
	app.Post("/users", handler.Store)

	// Daftarkan endpoint baru di sini!
	app.Get("/users", authMiddleware, handler.Fetch)
}
```

### Langkah 2: Buat Method Handler `Fetch`

Di dalam file [user_handler.go] (delivery/http/user_handler.go) yang sama, buat fungsi method handler baru dengan nama `Fetch`:

1. Ambil context request (`c.Context()`).
2. Panggil usecase `Fetch(ctx)` dari `UUsecase` untuk mengambil data semua user.
3. Tangani _error_ jika proses pemanggilan usecase gagal (kembalikan status `500 Internal Server Error`).
4. Jika sukses, kembalikan status `200 OK` beserta data array user.

### Langkah 3: Update `main.go`

Buka file [main.go] di root folder proyek:

1. Pada bagian registrasi _handler_ (sekitar baris 71), sesuaikan pemanggilan fungsi `NewUserHandler` dengan menyertakan variabel `jwtSecret` sebagai argumen ketiga.

_Petunjuk kode:_

```go
httpDelivery.NewUserHandler(app, userUsecase, jwtSecret)
```

---

## 🧪 Cara Pengujian & Uji Coba

1. **Jalankan Aplikasi**:
   ```bash
   go run main.go
   ```
2. **Lakukan Registrasi & Login**:
   - Daftarkan user baru terlebih dahulu dengan `POST /users`.
   - Lakukan login dengan `POST /login` untuk mendapatkan JWT Token.
3. **Panggil Endpoint `/users`**:
   - Panggil `GET http://localhost:8080/users`.
   - Coba akses **tanpa** Header `Authorization` (pastikan dapat response `401 Unauthorized`).
   - Coba akses **dengan** Header `Authorization: Bearer <token_jwt_kamu>` (pastikan mendapat list users dengan status `200 OK`).

---

🌟 **Selamat Mengerjakan!** 🌟
