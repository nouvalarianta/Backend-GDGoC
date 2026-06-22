# Dokumentasi API: Peminjaman & Pengembalian Buku (Borrowings)

Bagian ini mendokumentasikan endpoint untuk alur peminjaman buku oleh pengguna serta pengembalian buku.

---

## 1. Pinjam Buku (Borrow)

Membuat permohonan peminjaman buku. Status ketersediaan buku (`is_available`) akan berubah menjadi `false` secara otomatis.

- **URL**: `/borrowings/borrow`
- **Method**: `POST`
- **Autentikasi**: **Ya (Protected)**
- **Headers**:
  - `Authorization: Bearer <your_jwt_token>`
  - `Content-Type: application/json`

### Request Body (JSON)

| Field | Tipe | Wajib | Keterangan |
| :--- | :--- | :--- | :--- |
| `book_id` | Integer | Ya | ID unik buku yang ingin dipinjam |

**Contoh Request**:
```json
{
  "book_id": 1
}
```

### Response

#### **200 OK** (Sukses)
```json
{
  "message": "book borrowed successfully"
}
```

#### **400 Bad Request** (Gagal diproses)
Terjadi jika buku tidak tersedia (`is_available` bernilai `false`) atau buku tidak ditemukan.
```json
{
  "error": "book is not available for borrowing"
}
```

---

## 2. Kembalikan Buku (Return)

Mengembalikan buku yang sedang dipinjam berdasarkan ID Peminjaman. Status ketersediaan buku (`is_available`) akan kembali menjadi `true`.

- **URL**: `/borrowings/return/:id`
- **Method**: `POST`
- **Autentikasi**: **Ya (Protected)**
- **Path Parameters**:
  - `id` (integer) - ID Peminjaman (bukan ID Buku)
- **Headers**:
  - `Authorization: Bearer <your_jwt_token>`

### Response

#### **200 OK** (Sukses)
```json
{
  "message": "book returned successfully"
}
```

#### **400 Bad Request** (Gagal diproses)
Terjadi jika transaksi peminjaman sudah berstatus `returned`, atau user yang mencoba mengembalikan buku berbeda dengan user yang meminjamnya.
```json
{
  "error": "borrowing record not found or book already returned"
}
```

---

## 3. Daftar Peminjaman Saya (User Borrowings)

Mendapatkan riwayat daftar buku yang dipinjam oleh pengguna yang sedang login.

- **URL**: `/borrowings/user`
- **Method**: `GET`
- **Autentikasi**: **Ya (Protected)**
- **Headers**:
  - `Authorization: Bearer <your_jwt_token>`

### Response

#### **200 OK** (Sukses)
Mengembalikan array daftar riwayat peminjaman user bersangkutan beserta relasi data Buku.
```json
[
  {
    "id": 1,
    "user_id": 1,
    "book_id": 1,
    "borrow_date": "2026-06-20T11:00:00Z",
    "return_date": null,
    "status": "borrowed",
    "created_at": "2026-06-20T11:00:00Z",
    "updated_at": "2026-06-20T11:00:00Z",
    "book": {
      "id": 1,
      "title": "Clean Code",
      "author": "Robert C. Martin",
      "isbn": "978-0132350884",
      "is_available": false,
      "created_at": "2026-06-20T10:00:00Z",
      "updated_at": "2026-06-20T11:00:00Z"
    }
  }
]
```

---

## 4. Daftar Semua Peminjaman (Admin/Global View)

Mendapatkan seluruh data riwayat peminjaman dari semua pengguna.

- **URL**: `/borrowings`
- **Method**: `GET`
- **Autentikasi**: **Ya (Protected)**
- **Headers**:
  - `Authorization: Bearer <your_jwt_token>`

### Response

#### **200 OK** (Sukses)
Mengembalikan seluruh data peminjaman di sistem perpustakaan.
```json
[
  {
    "id": 1,
    "user_id": 1,
    "book_id": 1,
    "borrow_date": "2026-06-20T11:00:00Z",
    "return_date": "2026-06-20T14:00:00Z",
    "status": "returned",
    "created_at": "2026-06-20T11:00:00Z",
    "updated_at": "2026-06-20T14:00:00Z"
  }
]
```
