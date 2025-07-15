# GoBackend

GoBackend adalah backend API service menggunakan **Golang + Fiber**, dirancang modular dan siap produksi, dengan fitur lengkap seperti:

- 🔐 Login
- ✉️ Reset Password via Email
- 🧠 Akun diblokir sementara jika gagal login berulang
- 📧 Email notifikasi: reset, blokir, unlock
- 🧱 Modular: middleware, seeder, utils
- 📚 Pengelolaan konfigurasi lewat database (`settings` table)
- 🔌 Support RESTful API, JWT, CORS, dan OAuth2

---

## 🚀 Fitur Utama

- ✅ Auth: Login
- 🔐 Keamanan:
  - Blokir akun sementara jika gagal login >3x
  - Blokir IP jika mencoba brute force
- 🔄 Reset Password: via email + link token
- 📬 Email: kirim email via SMTP (semua setting disimpan di DB)
- ⚙️ Setting dinamis: diambil dari DB (bukan .env)
- 📈 Auto seeding data awal (`settings` dan lainnya)

---

## 🛠️ Instalasi & Setup

### 1. Clone repo dan install dependencies

```bash
git clone https://github.com/mysundoro/gobackend.git
cd gobackend
go mod tidy
.env (env.local)
buat database mysql misal gobackend
air (dev)
