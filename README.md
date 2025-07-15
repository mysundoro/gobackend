# GoBackend

GoBackend adalah backend API service menggunakan **Golang + Fiber**, dirancang modular dan siap produksi, dengan fitur keamanan dan otentikasi yang lengkap:

- 🔐 Login dengan perlindungan brute-force
- ✉️ Reset Password via Email
- 🧠 Akun diblokir sementara jika gagal login berulang
- 📧 Email notifikasi: reset, blokir, unlock akun
- 🧱 Modular: middleware, seeder, utils
- 📚 Pengelolaan konfigurasi lewat database (`settings` table)
- 🔌 Support RESTful API, JWT, dan CORS

---

## 🚀 Fitur Utama

- ✅ **Login Auth**: login dengan validasi email dan password
- 🔐 **Keamanan Otomatis**:
  - Akun diblokir jika gagal login lebih dari `max_failed_login_user`
  - IP diblokir sementara jika gagal login berulang (`max_failed_login_ip`)
  - Semua pengaturan (durasi blokir, jumlah maksimum, dll) dinamis dari database
- 🔄 **Reset Password**:
  - Token dikirim ke email
  - Link reset berlaku selama 15 menit
- 📬 **Email Notifikasi Otomatis**:
  - Reset Password
  - Akun Terkunci
  - Link Unlock Akun
- 📦 **Active Session**:
  - Login hanya bisa aktif di satu device (dengan `device_id`)
- ⚙️ **Pengaturan Dinamis**:
  - Semua konfigurasi seperti SMTP, blokir akun/IP disimpan di tabel `settings`
- 🔧 **Seeder Otomatis**:
  - Seeder akan mengisi nilai default `site_name`, `logo`, `smtp`, dll
- 👤 **User Management**:
  - CRUD user (Create, Read, Update, Delete)
  - Pagination, pencarian, dan sorting

---

## 🔐 Endpoint Keamanan & Auth

| Method | Endpoint                 | Deskripsi                                 |
|--------|--------------------------|-------------------------------------------|
| POST   | `/api/login`             | Login user dengan perlindungan brute-force |
| POST   | `/api/forgot-password`   | Request reset password (email token)      |
| POST   | `/api/reset-password`    | Setel ulang password menggunakan token    |
| GET    | `/api/unlock-account`    | Buka akun terkunci dengan token dari email |

---

## 👤 Endpoint Manajemen User

| Method | Endpoint         | Deskripsi                                          |
|--------|------------------|----------------------------------------------------|
| GET    | `/api/users`     | Ambil semua user (dengan pagination, search, sort) |
| POST   | `/api/users`     | Tambah user baru                                   |
| GET    | `/api/users/:id` | Ambil detail user berdasarkan ID                   |
| PUT    | `/api/users/:id` | Update user berdasarkan ID                         |
| DELETE | `/api/users/:id` | Hapus user berdasarkan ID                          |

> 🔐 Beberapa endpoint sebaiknya hanya diakses oleh admin (pastikan middleware otorisasi ditambahkan jika perlu).

---

## ⚙️ Endpoint Settings

| Method | Endpoint           | Deskripsi                                 |
|--------|--------------------|-------------------------------------------|
| GET    | `/api/settings`    | Ambil semua setting                       |
| GET    | `/api/settings/:k` | Ambil setting berdasarkan key             |
| POST   | `/api/settings`    | Buat setting baru                         |
| PUT    | `/api/settings/:k` | Update setting berdasarkan key            |

---

## 🛠️ Instalasi & Setup

### 1. Clone repo dan install dependency

```bash
git clone https://github.com/mysundoro/gobackend.git
cd gobackend
go mod tidy
cp env.local .env
air / go run main.go