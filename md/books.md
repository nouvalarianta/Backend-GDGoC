# Dokumentasi API: Manajemen Buku (Books)

Bagian ini mendokumentasikan endpoint untuk mengelola katalog buku perpustakaan.

---

## 1. Dapatkan Semua Buku

Mendapatkan daftar seluruh buku yang ada di perpustakaan.

- **URL**: `/books`
- **Method**: `GET`
- **Autentikasi**: Tidak (Public)

### Response

#### **200 OK** (Sukses)
Mengembalikan array daftar buku.
```json
[
  {
    "id": 1,
    "title": "Clean Code",
    "author": "Robert C. Martin",
    "isbn": "978-0132350884",
    "is_available": true,
    "created_at": "2026-06-20T10:00:00Z",
    "updated_at": "2026-06-20T10:00:00Z"
  },
  {
    "id": 2,
    "title": "The Clean Coder",
    "author": "Robert C. Martin",
    "isbn": "978-0137081073",
    "is_available": false,
    "created_at": "2026-06-20T10:05:00Z",
    "updated_at": "2026-06-20T10:05:00Z"
  }
]
```

---

## 2. Dapatkan Detail Buku

Mencari satu buku berdasarkan ID.

- **URL**: `/books/:id`
- **Method**: `GET`
- **Autentikasi**: Tidak (Public)
- **Path Parameters**:
  - `id` (integer) - ID unik buku

### Response

#### **200 OK** (Sukses)
Mengembalikan detail data buku.
```json
{
  "id": 1,
  "title": "Clean Code",
  "author": "Robert C. Martin",
  "isbn": "978-0132350884",
  "is_available": true,
  "created_at": "2026-06-20T10:00:00Z",
  "updated_at": "2026-06-20T10:00:00Z"
}
```

#### **400 Bad Request** (ID tidak valid)
Terjadi jika ID bukan angka.
```json
{
  "error": "Invalid book ID"
}
```

#### **404 Not Found** (Buku tidak ditemukan)
```json
{
  "error": "book not found"
}
```

---

## 3. Tambah Buku Baru

Menyimpan buku baru ke dalam katalog perpustakaan.

- **URL**: `/books`
- **Method**: `POST`
- **Autentikasi**: **Ya (Protected)**
- **Headers**:
  - `Authorization: Bearer <your_jwt_token>`
  - `Content-Type: application/json`

### Request Body (JSON)

| Field | Tipe | Wajib | Keterangan |
| :--- | :--- | :--- | :--- |
| `title` | String | Ya | Judul buku |
| `author` | String | Ya | Penulis buku |
| `isbn` | String | Ya | Kode ISBN unik buku |
| `is_available` | Boolean | Opsional | Ketersediaan buku (Default: `true`) |

**Contoh Request**:
```json
{
  "title": "Refactoring",
  "author": "Martin Fowler",
  "isbn": "978-0134757599",
  "is_available": true
}
```

### Response

#### **201 Created** (Sukses)
Mengembalikan data buku yang berhasil disimpan.
```json
{
  "id": 3,
  "title": "Refactoring",
  "author": "Martin Fowler",
  "isbn": "978-0134757599",
  "is_available": true,
  "created_at": "2026-06-20T15:25:00Z",
  "updated_at": "2026-06-20T15:25:00Z"
}
```

#### **401 Unauthorized** (Belum login/token salah)
```json
{
  "error": "Missing or malformed JWT"
}
```

---

## 4. Perbarui Buku (Update)

Mengubah data buku yang sudah ada berdasarkan ID.

- **URL**: `/books/:id`
- **Method**: `PUT`
- **Autentikasi**: **Ya (Protected)**
- **Path Parameters**:
  - `id` (integer) - ID unik buku
- **Headers**:
  - `Authorization: Bearer <your_jwt_token>`
  - `Content-Type: application/json`

### Request Body (JSON)

Sama seperti request body pada "Tambah Buku Baru".

### Response

#### **200 OK** (Sukses)
Mengembalikan data buku terbaru setelah berhasil diperbarui.
```json
{
  "id": 3,
  "title": "Refactoring (2nd Edition)",
  "author": "Martin Fowler",
  "isbn": "978-0134757599",
  "is_available": true,
  "created_at": "2026-06-20T15:25:00Z",
  "updated_at": "2026-06-20T15:30:00Z"
}
```

---

## 5. Hapus Buku (Delete)

Menghapus data buku dari katalog perpustakaan.

- **URL**: `/books/:id`
- **Method**: `DELETE`
- **Autentikasi**: **Ya (Protected)**
- **Path Parameters**:
  - `id` (integer) - ID unik buku
- **Headers**:
  - `Authorization: Bearer <your_jwt_token>`

### Response

#### **204 No Content** (Sukses)
Buku berhasil dihapus, tidak mengembalikan konten/body.

#### **404 Not Found** / **400 Bad Request**
Sama seperti pada endpoint pencarian detail buku.
```json
{
  "error": "book not found"
}
```
