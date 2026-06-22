package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
    "strconv"
    "strings"
    "time"

    "github.com/golang-jwt/jwt/v5"
    _ "github.com/lib/pq"
    "golang.org/x/crypto/bcrypt"
)

// GLOBAL VARIABLES - Ciri khas spageti code awal
var db *sql.DB
var jwtSecret = "super-secret-key-that-should-not-be-hardcoded"

// Fungsi main yang sangat panjang dan melakukan segalanya
func main() {
    fmt.Println("Starting Spaghetti API...")

    // 1. Baca Config & Koneksi DB (Tanpa error handling yang proper, langsung exit)
    host := os.Getenv("DB_HOST")
    if host == "" {
        host = "localhost"
    }
    port := os.Getenv("DB_PORT")
    if port == "" {
        port = "5432"
    }
    user := os.Getenv("DB_USER")
    if user == "" {
        user = "postgres"
    }
    pass := os.Getenv("DB_PASS")
    if pass == "" {
        pass = "postgres"
    }
    dbname := os.Getenv("DB_NAME")
    if dbname == "" {
        dbname = "library_db"
    }

    connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, dbname)
    var err error
    db, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal("Cannot connect to db", err)
    }
    err = db.Ping()
    if err != nil {
        log.Fatal("Cannot ping db", err)
    }
    defer db.Close()

    // 2. Setup Routing menggunakan standard library dengan Mux buatan sendiri (Spageti routing)
    http.HandleFunc("/api/v1/auth/register", func(w http.ResponseWriter, r *http.Request) {
        // Method check manual
        if r.Method != "POST" {
            w.WriteHeader(http.StatusMethodNotAllowed)
            w.Write([]byte(`{"error": "Method not allowed"}`))
            return
        }

        // Baca body
        body, err := io.ReadAll(r.Body)
        if err != nil {
            w.WriteHeader(http.StatusBadRequest)
            w.Write([]byte(`{"error": "Cannot read body"}`))
            return
        }

        // Parse JSON inline menggunakan map (tanpa struct)
        var req map[string]interface{}
        err = json.Unmarshal(body, &req)
        if err != nil {
            w.WriteHeader(http.StatusBadRequest)
            w.Write([]byte(`{"error": "Invalid JSON"}`))
            return
        }

        // Validasi manual yang berulang-ulang
        username, ok := req["username"].(string)
        if !ok || username == "" {
            w.WriteHeader(http.StatusBadRequest)
            w.Write([]byte(`{"error": "Username required"}`))
            return
        }

        password, ok := req["password"].(string)
        if !ok || password == "" {
            w.WriteHeader(http.StatusBadRequest)
            w.Write([]byte(`{"error": "Password required"}`))
            return
        }

        if len(password) < 6 {
            w.WriteHeader(http.StatusBadRequest)
            w.Write([]byte(`{"error": "Password too short"}`))
            return
        }

        name, ok := req["name"].(string)
        if !ok || name == "" {
            w.WriteHeader(http.StatusBadRequest)
            w.Write([]byte(`{"error": "Name required"}`))
            return
        }

        // Cek user exists (Logic DB di dalam HTTP Handler)
        var exists int
        err = db.QueryRow("SELECT 1 FROM users WHERE username = $1", username).Scan(&exists)
        if err == nil {
            w.WriteHeader(http.StatusConflict)
            w.Write([]byte(`{"error": "Username already taken"}`))
            return
        } else if err != sql.ErrNoRows {
            w.WriteHeader(http.StatusInternalServerError)
            w.Write([]byte(`{"error": "DB Error"}`))
            return
        }

        // Hash password (Logic Enkripsi di dalam HTTP Handler)
        hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            w.Write([]byte(`{"error": "Hash error"}`))
            return
        }

        // Insert (Raw query lagi)
        var newID int
        err = db.QueryRow("INSERT INTO users (username, password, name, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id",
            username, string(hashed), name, time.Now(), time.Now()).Scan(&newID)
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            w.Write([]byte(`{"error": "Failed to create user"}`))
            return
        }

        // Response
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusCreated)
        w.Write([]byte(fmt.Sprintf(`{"id": %d, "message": "success"}`, newID)))
    })

    http.HandleFunc("/api/v1/auth/login", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != "POST" {
            w.WriteHeader(http.StatusMethodNotAllowed)
            return
        }

        // Repetisi baca body (Duplikasi kode 1)
        body, err := io.ReadAll(r.Body)
        if err != nil {
            w.WriteHeader(http.StatusBadRequest)
            w.Write([]byte(`{"error": "Cannot read body"}`))
            return
        }

        var req map[string]interface{}
        err = json.Unmarshal(body, &req)
        if err != nil {
            w.WriteHeader(http.StatusBadRequest)
            w.Write([]byte(`{"error": "Invalid JSON"}`))
            return
        }

        username, _ := req["username"].(string)
        password, _ := req["password"].(string)

        // DB Query lagi
        var id int
        var hash string
        err = db.QueryRow("SELECT id, password FROM users WHERE username = $1", username).Scan(&id, &hash)
        if err != nil {
            if err == sql.ErrNoRows {
                w.WriteHeader(http.StatusUnauthorized)
                w.Write([]byte(`{"error": "Invalid credentials"}`))
                return
            }
            w.WriteHeader(http.StatusInternalServerError)
            return
        }

        // Compare
        err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
        if err != nil {
            w.WriteHeader(http.StatusUnauthorized)
            w.Write([]byte(`{"error": "Invalid credentials"}`))
            return
        }

        // Bikin JWT Inline
        token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
            "user_id": id,
            "exp":     time.Now().Add(time.Hour * 24).Unix(),
        })
        tokenStr, err := token.SignedString([]byte(jwtSecret))
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(fmt.Sprintf(`{"token": "%s"}`, tokenStr)))
    })

    // HANDLER BUKU DAN PEMINJAMAN (Semua dilokasikan dalam fungsi main)
    http.HandleFunc("/api/v1/borrow", func(w http.ResponseWriter, r *http.Request) {
        // === AUTHENTICATION INLINE (Tidak pakai middleware, manual di setiap route) ===
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            w.WriteHeader(http.StatusUnauthorized)
            w.Write([]byte(`{"error": "Unauthorized"}`))
            return
        }
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            w.WriteHeader(http.StatusUnauthorized)
            w.Write([]byte(`{"error": "Unauthorized format"}`))
            return
        }
        token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
            return []byte(jwtSecret), nil
        })
        if err != nil || !token.Valid {
            w.WriteHeader(http.StatusUnauthorized)
            w.Write([]byte(`{"error": "Invalid token"}`))
            return
        }
        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            w.WriteHeader(http.StatusUnauthorized)
            return
        }
        userID := int(claims["user_id"].(float64))
        // === SELESAI AUTHENTICATION INLINE ===

        if r.Method == "POST" {
            // DEEPLY NESTED BLOCKS (Arrow Code)
            body, err := io.ReadAll(r.Body)
            if err == nil {
                var req map[string]interface{}
                err = json.Unmarshal(body, &req)
                if err == nil {
                    bookIDFloat, ok := req["book_id"].(float64)
                    if ok {
                        bookID := int(bookIDFloat)
                        // Cek availability buku (Business logic nyampur)
                        var isAvailable bool
                        err = db.QueryRow("SELECT is_available FROM books WHERE id = $1", bookID).Scan(&isAvailable)
                        if err == nil {
                            if isAvailable {
                                // Mulai transaksi manual yang berantakan
                                tx, err := db.Begin()
                                if err == nil {
                                    _, err = tx.Exec("UPDATE books SET is_available = false WHERE id = $1", bookID)
                                    if err == nil {
                                        _, err = tx.Exec("INSERT INTO borrowings (user_id, book_id, borrow_date, status) VALUES ($1, $2, $3, $4)",
                                            userID, bookID, time.Now(), "borrowed")
                                        if err == nil {
                                            err = tx.Commit()
                                            if err == nil {
                                                w.WriteHeader(http.StatusOK)
                                                w.Write([]byte(`{"message": "Success borrow"}`))
                                                return
                                            } else {
                                                tx.Rollback()
                                                w.WriteHeader(http.StatusInternalServerError)
                                                w.Write([]byte(`{"error": "Commit failed"}`))
                                            }
                                        } else {
                                            tx.Rollback()
                                            w.WriteHeader(http.StatusInternalServerError)
                                            w.Write([]byte(`{"error": "Insert borrow failed"}`))
                                        }
                                    } else {
                                        tx.Rollback()
                                        w.WriteHeader(http.StatusInternalServerError)
                                        w.Write([]byte(`{"error": "Update book failed"}`))
                                    }
                                } else {
                                    w.WriteHeader(http.StatusInternalServerError)
                                    w.Write([]byte(`{"error": "Tx start failed"}`))
                                }
                            } else {
                                w.WriteHeader(http.StatusBadRequest)
                                w.Write([]byte(`{"error": "Book not available"}`))
                            }
                        } else {
                            w.WriteHeader(http.StatusNotFound)
                            w.Write([]byte(`{"error": "Book not found"}`))
                        }
                    } else {
                        w.WriteHeader(http.StatusBadRequest)
                        w.Write([]byte(`{"error": "Book ID required"}`))
                    }
                } else {
                    w.WriteHeader(http.StatusBadRequest)
                    w.Write([]byte(`{"error": "Invalid JSON"}`))
                }
            } else {
                w.WriteHeader(http.StatusBadRequest)
                w.Write([]byte(`{"error": "Body read error"}`))
            }
        } else {
            w.WriteHeader(http.StatusMethodNotAllowed)
        }
    })

    // Route lain dengan duplikasi kode pengecekan otentikasi
    http.HandleFunc("/api/v1/return", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != "POST" {
            w.WriteHeader(http.StatusMethodNotAllowed)
            return
        }

        // === DUPLIKASI AUTHENTICATION LAGI ===
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            w.WriteHeader(http.StatusUnauthorized)
            return
        }
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            w.WriteHeader(http.StatusUnauthorized)
            return
        }
        token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
            return []byte(jwtSecret), nil
        })
        if err != nil || !token.Valid {
            w.WriteHeader(http.StatusUnauthorized)
            return
        }
        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            w.WriteHeader(http.StatusUnauthorized)
            return
        }
        userID := int(claims["user_id"].(float64))
        // === SELESAI DUPLIKASI AUTHENTICATION ===

        body, _ := io.ReadAll(r.Body)
        var req map[string]interface{}
        json.Unmarshal(body, &req)

        // Penggunaan id dengan konversi manual yang rawan error
        borrowIDStr := fmt.Sprintf("%v", req["borrow_id"])
        borrowID, _ := strconv.Atoi(borrowIDStr)

        var bookID int
        var status string
        var actualUserID int
        err = db.QueryRow("SELECT book_id, status, user_id FROM borrowings WHERE id = $1", borrowID).Scan(&bookID, &status, &actualUserID)
        
        if err != nil || actualUserID != userID || status == "returned" {
            w.WriteHeader(http.StatusBadRequest)
            w.Write([]byte(`{"error": "Cannot return this book"}`))
            return
        }

        // Magic number dan format tanggal berantakan
        currentTime := time.Now().Format("2006-01-02 15:04:05")
        
        tx, _ := db.Begin()
        tx.Exec("UPDATE borrowings SET status = 'returned', return_date = $1 WHERE id = $2", currentTime, borrowID)
        tx.Exec("UPDATE books SET is_available = true WHERE id = $1", bookID)
        tx.Commit()

        w.Write([]byte(`{"message": "returned"}`))
    })

    log.Println("Server running at :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}

// Catatan Spageti Code:
// 1. Tidak ada struktur package (domain, usecase, delivery, dsb). Semuanya file `main.go`.
// 2. Satu fungsi `main()` raksasa mengatur konfigurasi, database, routing, logika bisnis, validasi, dan response HTTP.
// 3. Arrow Anti-Pattern: Terdapat blok `if` yang sangat menjorok ke dalam (di route /api/v1/borrow) karena error handling bersarang.
// 4. Duplikasi Kode: Logika parsing JWT (autentikasi) dan parsing JSON berulang-ulang di setiap handler.
// 5. Global State: Menggunakan variabel `db` dan `jwtSecret` secara global, yang menyulitkan testing (Unit Testing hampir mustahil dilakukan).
// 6. Magic Strings: Penamaan status seperti "borrowed" dan "returned" ditulis secara hardcode.
// 7. Pengabaian Error: Pada route "/api/v1/return", ada beberapa error `_` yang diabaikan secara paksa.
// 8. Tightly Coupled: Sangat sulit jika kita ingin memindahkan database dari PostgreSQL ke MySQL, atau jika kita ingin membuat fitur ini dapat diakses melalui CLI/gRPC, karena logika bisnis menyatu erat dengan HTTP `http.ResponseWriter`.


