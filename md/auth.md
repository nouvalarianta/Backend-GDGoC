# Dokumentasi API: Autentikasi & Pengguna (Auth & Users)

Bagian ini mendokumentasikan endpoint yang terkait dengan registrasi akun baru dan autentikasi (login) pengguna.

---

## 1. Pendaftaran Pengguna Baru (Register)

Membuat akun pengguna baru untuk dapat masuk ke sistem perpustakaan.

- **URL**: `/users`
- **Method**: `POST`
- **Headers**:
  - `Content-Type: application/json`

### Request Body (JSON)

| Field | Tipe | Wajib | Keterangan |
| :--- | :--- | :--- | :--- |
| `username` | String | Ya | Nama unik untuk login |
| `password` | String | Ya | Kata sandi minimal 6 karakter (akan di-hash menggunakan bcrypt) |
| `name` | String | Ya | Nama lengkap pengguna |

**Contoh Request**:
```json
{
  "username": "nouval",
  "password": "mysecurepassword",
  "name": "Nouval"
}
```

### Response

#### **201 Created** (Sukses)
Mengembalikan data pengguna yang telah berhasil didaftarkan (password disembunyikan/tidak dikirim kembali).
```json
{
  "id": 1,
  "username": "nouval",
  "name": "Nouval",
  "created_at": "2026-06-20T15:20:00Z",
  "updated_at": "2026-06-20T15:20:00Z"
}
```

#### **422 Unprocessable Entity** (Payload salah)
Terjadi jika format JSON salah atau gagal diparsing.
```json
{
  "error": "invalid character '...' looking for beginning of value"
}
```

#### **500 Internal Server Error** (Gagal diproses)
Terjadi jika ada kesalahan database (misalnya username sudah terdaftar).
```json
{
  "error": "username already taken"
}
```

---

## 2. Masuk Log (Login)

Melakukan autentikasi menggunakan username dan password untuk mendapatkan token JWT.

- **URL**: `/login`
- **Method**: `POST`
- **Headers**:
  - `Content-Type: application/json`

### Request Body (JSON)

| Field | Tipe | Wajib | Keterangan |
| :--- | :--- | :--- | :--- |
| `username` | String | Ya | Username pengguna |
| `password` | String | Ya | Password pengguna |

**Contoh Request**:
```json
{
  "username": "nouval",
  "password": "mysecurepassword"
}
```

### Response

#### **200 OK** (Sukses)
Mengembalikan token JWT yang valid selama 24 jam.
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzA5NzI4MDAsInVzZXJfaWQiOjF9..."
}
```

#### **400 Bad Request** (Validasi gagal)
Username atau password tidak diisi.
```json
{
  "error": "username and password are required"
}
```

#### **401 Unauthorized** (Autentikasi gagal)
Username tidak ditemukan atau password salah.
```json
{
  "error": "invalid username or password"
}
```
