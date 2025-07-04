# Sistem Aplikasi Pemesanan Makanan Digital untuk Meningkatkan Efisiensi Layanan Dine-In di Restoran

| Nama                          | NRP        | Role     |
| ----------------------------- | ---------- | -------- |
| Muhammad Budhi Salmanjannah   | 5025201084 | Frontend |
| Ainun Nadhifah Syamsiyah      | 5025221053 | Backend  |
| Fawwas Aldy Nurramdhan Kaisar | 5025221179 | Backend  |
| Muammar Bahalwan              | 5053231020 | Frontend |

## 🍽️ Gambaran Umum

Aplikasi sistem pemesanan makanan untuk restoran dine-in dirancang untuk meningkatkan efisiensi layanan dan kenyamanan pelanggan. Melalui pemindaian QR code di meja, pelanggan dapat langsung mengakses menu digital, memilih makanan berdasarkan ketersediaan dan kategori, serta melakukan pembayaran tanpa perlu menunggu pelayan. Sistem ini juga mendukung alur kerja dapur dan pelayan dengan fitur pelacakan pesanan secara real-time dan pengelolaan antrian pesanan.

Tujuan utama dari sistem ini adalah menyediakan pengalaman pemesanan yang cepat, akurat, dan terintegrasi antar peran pengguna seperti pelanggan, staf dapur, dan pelayan. Arsitektur perangkat lunak yang dibangun harus mengutamakan keandalan, skalabilitas, dan kemudahan penggunaan, dengan harapan dapat mempercepat proses layanan, mengurangi kesalahan operasional, dan meningkatkan kepuasan semua pihak yang terlibat.

Laporan selengkapnya dapat dilihat di [sini](https://drive.google.com/file/d/1pvrBNacUcCM_Vs9nj4ThdDpkH-ErhTYz/view?usp=sharing).

## 🏗️ Arsitektur

Proyek ini mengikuti prinsip **Clean Architecture** dengan pemisahan tanggung jawab yang jelas:

```
gacoan-backend/
├── application/          # Lapisan aplikasi (layanan, permintaan, respons)
├── domain/              # Lapisan domain (entitas, logika bisnis)
├── infrastructure/      # Lapisan infrastruktur (database, layanan eksternal)
├── presentation/        # Lapisan presentasi (controller, rute, middleware)
├── platform/           # Utilitas platform bersama
├── command/            # Perintah CLI
└── test/              # Suite pengujian komprehensif
```

### Pola Desain Utama

- **Domain-Driven Design (DDD)**
- **Repository Pattern**
- **Dependency Injection**
- **Middleware Pattern**
- **Factory Pattern**

## 🚀 Fitur

### 🔐 Autentikasi & Otorisasi

- Autentikasi berbasis JWT
- Kontrol akses berbasis peran (RBAC)
- Berbagai peran pengguna: Pelanggan, Dapur, Pelayan, Super Admin

### 📋 Manajemen Pesanan

- Manajemen siklus hidup pesanan lengkap
- Pelacakan status pesanan real-time
- Manajemen antrian dengan kode antrian unik
- Riwayat pesanan dan pagination

### 👨‍🍳 Operasi Dapur

- **Mulai Memasak**: Memulai persiapan makanan
- **Selesai Memasak**: Menandai pesanan siap disajikan
- **Pesanan Berikutnya**: Mendapatkan pesanan berikutnya dalam antrian
- Pelacakan waktu memasak dan deteksi keterlambatan

### 🍽️ Operasi Pelayan

- **Siap Disajikan**: Melihat pesanan siap untuk pengiriman
- **Mulai Mengantar**: Memulai pengiriman makanan
- **Selesai Mengantar**: Menyelesaikan pengiriman pesanan
- Pembaruan status pesanan real-time

### 💳 Pemrosesan Pembayaran

- Integrasi gateway pembayaran Midtrans
- Pelacakan status pembayaran
- Penanganan webhook untuk pembaruan pembayaran
- Pemrosesan transaksi yang aman

### 📊 Manajemen Menu

- Organisasi menu berbasis kategori
- Manajemen ketersediaan menu
- Manajemen harga dengan presisi desimal
- Operasi CRUD item menu

### 🏢 Manajemen Restoran

- Manajemen meja
- Manajemen pengguna
- Riwayat transaksi
- Pelaporan komprehensif

## 🛠️ Stack Teknologi

- **Bahasa**: Go 1.24
- **Framework**: Gin (framework web HTTP)
- **Database**: PostgreSQL dengan GORM
- **Autentikasi**: JWT
- **Gateway Pembayaran**: Midtrans
- **Pengujian**: Testify
- **Validasi**: Validasi domain kustom
- **Penanganan Desimal**: ShopSpring Decimal

## 📋 Prasyarat

- Go 1.24 atau lebih tinggi
- PostgreSQL 12 atau lebih tinggi
- Git

## 🚀 Instalasi & Setup

### 1. Clone Repository

```bash
git clone https://github.com/ainunns/gacoan-backend.git
cd gacoan-backend
```

### 2. Instalasi Dependensi

```bash
go mod download
```

### 3. Konfigurasi Environment

Salin dan isi file `.env.example` menjadi `.env` di direktori root:

```bash
cp .env.example .env
```

### 4. Setup Database

```bash
# Jalankan migrasi database
go run main.go --migrate

# Seed data awal
go run main.go --seed
```

### 5. Jalankan Aplikasi

```bash
go run main.go
```

Server akan berjalan di `http://localhost:8888`

## 🧪 Pengujian

### Jalankan Semua Test

```bash
go test ./test/... -v
```

### Jalankan File Test Tertentu

```bash
go test ./test/start_delivering_test.go -v
```

### Jalankan Fungsi Test Tertentu

```bash
go test ./test/ -run TestStartDelivering_Success -v
```

### Cakupan Test

```bash
go test ./test/... -cover
```

## 📚 Dokumentasi API

### Base URL

```
http://localhost:8888/api
```

### API Documentation

Dokumentasi API dapat dilihat di [sini](https://documenter.getpostman.com/view/31404175/2sB2xFgUDM).

### Autentikasi

Semua endpoint yang dilindungi memerlukan token JWT di header Authorization:

```
Authorization: Bearer <your_jwt_token>
```

### Endpoint Utama

#### 🔐 Autentikasi

- `POST /user/register` - Registrasi pengguna
- `POST /user/login` - Login pengguna

#### 📋 Transaksi

- `POST /transaction/` - Buat transaksi baru
- `GET /transaction/` - Dapatkan semua transaksi (dengan pagination)
- `GET /transaction/:id` - Dapatkan transaksi berdasarkan ID
- `POST /transaction/hook` - Webhook pembayaran

#### 👨‍🍳 Operasi Dapur

- `GET /transaction/next-order` - Dapatkan pesanan berikutnya dalam antrian
- `POST /transaction/start-cooking` - Mulai memasak pesanan
- `POST /transaction/finish-cooking` - Selesai memasak pesanan

#### 🍽️ Operasi Pelayan

- `GET /transaction/ready-to-serve` - Dapatkan pesanan siap disajikan
- `POST /transaction/start-delivering` - Mulai mengantar pesanan
- `POST /transaction/finish-delivering` - Selesai mengantar pesanan

#### 📊 Manajemen Menu

- `GET /menu/` - Dapatkan semua menu
- `POST /menu/` - Buat item menu baru
- `PUT /menu/:id/availability` - Perbarui ketersediaan menu

#### 🏢 Manajemen Restoran

- `GET /table/` - Dapatkan semua meja
- `GET /category/` - Dapatkan semua kategori
- `GET /user/` - Dapatkan semua pengguna

## 👥 Peran Pengguna & Izin

### 🛒 Pelanggan

- Membuat transaksi
- Melihat riwayat transaksi sendiri
- Melakukan pembayaran

### 👨‍🍳 Dapur

- Melihat pesanan berikutnya dalam antrian
- Mulai/selesai memasak pesanan
- Melihat detail pesanan

### 🍽️ Pelayan

- Melihat pesanan siap disajikan
- Mulai/selesai mengantar pesanan
- Memperbarui status pesanan

## 📊 Alur Status Pesanan

```
Pending → Preparing → Ready to Serve → Delivering → Served
   ↓         ↓            ↓              ↓          ↓
Pelanggan  Dapur       Dapur         Pelayan    Pelayan
Membuat   Mulai       Selesai       Mulai      Selesai
Pesanan   Memasak     Memasak       Mengantar  Mengantar
```

## 🔧 Pengembangan

### Struktur Proyek

```
├── application/
│   ├── request/          # DTO Permintaan
│   ├── response/         # DTO Respons
│   └── service/          # Layanan aplikasi
├── domain/
│   ├── identity/         # Objek nilai ID
│   ├── menu/             # Domain menu
│   ├── order/            # Domain pesanan
│   ├── shared/           # Objek domain bersama
│   ├── table/            # Domain meja
│   ├── transaction/      # Domain transaksi
│   └── user/             # Domain pengguna
├── infrastructure/
│   ├── adapter/          # Adapter layanan eksternal
│   ├── database/         # Lapisan database
│   └── validation/       # Validasi database
├── presentation/
│   ├── controller/       # Controller HTTP
│   ├── middleware/       # Middleware HTTP
│   ├── route/            # Definisi rute
│   └── message/          # Pesan error
└── test/                 # File pengujian
```

## 🚀 Deployment

1. Set `APP_ENV=production` di variabel environment
2. Konfigurasi database produksi
3. Atur secret JWT yang tepat
4. Konfigurasi kunci produksi Midtrans
5. Atur reverse proxy (nginx/apache)
6. Gunakan process manager (systemd/pm2)

## Lainnya

Repository frontend dapat dilihat di [sini](https://github.com/salmanhermana/gacoan-frontend/).
