# Sistem Authentication (Go + React + PostgreSQL)

Aplikasi authentication untuk aplikasi internal tim kecil dengan dua peran: **user** dan **admin**. Backend ditulis dengan Go, frontend dengan React (Vite), database PostgreSQL, dan seluruhnya dijalankan dengan Docker Compose.

> Catatan: bagian bertanda **[SESUAIKAN]** harus diperiksa dan disesuaikan dengan kode yang sebenarnya sebelum dikumpulkan. Hapus tanda ini setelah selesai.

## Daftar Isi

1. [Fitur](#fitur)
2. [Cara Menjalankan](#cara-menjalankan)
3. [Login Admin Pertama](#login-admin-pertama)
4. [Daftar Endpoint](#daftar-endpoint)
5. [Keputusan Teknis](#keputusan-teknis)
6. [Penggunaan AI](#penggunaan-ai)
7. [Catatan Implementasi](#catatan-implementasi)

## Fitur

**Kebutuhan minimum**

- Registrasi dengan email, username, dan kata sandi.
- Login dengan email **atau** username beserta kata sandi, dan logout.
- Halaman profil (hanya setelah login) yang menampilkan data milik pengguna sendiri.
- Dua peran (user dan admin). Hanya admin yang dapat melihat daftar seluruh pengguna.
- Akun admin pertama dibuat otomatis saat aplikasi dijalankan.
- Maksimal 2 sesi aktif per akun. Login ketiga mengakhiri sesi paling lama.
- Admin dapat melakukan force logout terhadap pengguna lain.

**Fitur tambahan** **[SESUAIKAN]**

- (contoh) Validasi input di backend dan frontend
- (contoh) Rate limiting pada endpoint login
- (contoh) Daftar sesi aktif di halaman profil
- (contoh) Migrasi database otomatis saat startup

## Cara Menjalankan

Prasyarat: Docker dan Docker Compose.

```bash
git clone https://github.com/arvian-vicky/goreact-management.git
cd goreact-management
cp .env.example .env
docker compose up -d --build
```

Setelah selesai:

| Layanan     | URL                   |
|-------------|-----------------------|
| Frontend    | http://localhost:3000 |
| Backend API | http://localhost:8080 |

Migrasi database dan pembuatan akun admin berjalan otomatis saat backend start. Tidak ada langkah manual tambahan.

Menghentikan aplikasi:

```bash
docker compose down        # data tetap tersimpan
docker compose down -v     # data database ikut dihapus
```

### Variabel environment (`.env.example`)

| Variabel         | Keterangan                          |
|------------------|-------------------------------------|
| `DB_PASSWORD`    | Password PostgreSQL                 |
| `JWT_SECRET`     | Secret untuk penandatanganan token **[SESUAIKAN]** |
| `ADMIN_EMAIL`    | Email admin pertama                 |
| `ADMIN_USERNAME` | Username admin pertama              |
| `ADMIN_PASSWORD` | Kata sandi admin pertama            |

## Login Admin Pertama

Saat startup, backend memeriksa apakah sudah ada akun admin. Jika belum, backend membuat akun admin dari nilai `ADMIN_EMAIL`, `ADMIN_USERNAME`, dan `ADMIN_PASSWORD` di `.env`. Proses ini idempotent, sehingga restart tidak membuat akun ganda. **[SESUAIKAN: sebutkan file/fungsi seed]**

Kredensial default dari `.env.example`:

- Email: `admin@example.com`
- Username: `admin`
- Kata sandi:

Ganti nilai ini untuk penggunaan di luar development.

## Daftar Endpoint

**[SESUAIKAN dengan route di kode]**

| Method | Endpoint                          | Akses  | Deskripsi                                   |
|--------|-----------------------------------|--------|---------------------------------------------|
| POST   | `/auth/register`                  | Publik | Registrasi (email, username, kata sandi)    |
| POST   | `/auth/login`                     | Publik | Login dengan email atau username            |
| POST   | `/auth/logout`                    | Login  | Mengakhiri sesi saat ini                    |
| GET    | `/me`                             | Login  | Profil pengguna yang sedang login           |
| GET    | `/admin/users`                    | Admin  | Daftar seluruh pengguna                     |
| POST   | `/admin/users/:id/force-logout`   | Admin  | Mengakhiri seluruh sesi aktif pengguna      |

## Keputusan Teknis

**[SESUAIKAN]**

- **Sesi disimpan di database** (tabel `sessions`), bukan hanya JWT stateless, agar batas 2 sesi dan force logout bisa ditegakkan seketika.
- **Kata sandi di-hash dengan bcrypt**, tidak pernah disimpan atau dikembalikan dalam bentuk asli.
- **Migrasi SQL dijalankan otomatis** saat backend start, sehingga cukup satu perintah `docker compose up`.
- **Healthcheck PostgreSQL** dan `depends_on: service_healthy` memastikan backend baru start setelah database siap.
- **Host database memakai nama service** (`psql`), bukan `localhost`, karena komunikasi antar-container.
- **Port database tidak dipublikasikan** ke host.

## Penggunaan AI

**[SESUAIKAN dengan penggunaan yang sebenarnya]**

- **Claude (Anthropic)**: membantu menyusun dan memperbaiki konfigurasi Docker Compose (healthcheck, environment, urutan startup), debugging error build dan Git, serta menyusun kerangka README.
- (tambahkan tools lain jika ada, misalnya GitHub Copilot, ChatGPT, dan bagian pekerjaan yang dibantu)

Seluruh kode telah saya baca, uji, dan pahami, dan saya bertanggung jawab atas hasil akhirnya.

## Catatan Implementasi

> Jawaban di bawah adalah draf berdasarkan desain umum. **[SESUAIKAN]** nama file dan fungsi dengan kode sebenarnya, dan ubah isinya bila implementasimu berbeda.

**1. Penyimpanan kata sandi.**
Kata sandi di-hash dengan bcrypt (`golang.org/x/crypto/bcrypt`) di `internal/auth/password.go` (fungsi `HashPassword` dan `CheckPassword`). Hash disimpan di kolom `users.password_hash`, dan kata sandi asli tidak pernah disimpan maupun di-log. Bcrypt dipilih karena lambat secara disengaja, memakai salt acak per kata sandi, dan cost factor-nya bisa dinaikkan seiring perkembangan perangkat keras, sehingga serangan brute force dan rainbow table menjadi mahal.

**2. Definisi sesi dan penegakan batas 2 sesi.**
Satu sesi adalah satu baris di tabel `sessions` (id, user_id, created_at, expires_at, revoked) yang dibuat setiap login berhasil. Token yang diberikan ke klien memuat id sesi, dan setiap request terautentikasi memeriksa bahwa sesi tersebut masih ada, belum dicabut, dan belum kedaluwarsa. Batas 2 sesi ditegakkan di fungsi `Login` pada `internal/auth/service.go`: dalam satu transaksi, jika jumlah sesi aktif pengguna sudah 2, sesi dengan `created_at` paling lama dihapus sebelum sesi baru dibuat.

**3. Jika token atau sesi dicuri.**
Penyerang dapat bertindak sebagai korban selama sesi tersebut masih valid, sampai token kedaluwarsa atau sesi dicabut. Karena sesi disimpan di database, pemilik akun dapat logout, dan jika penyerang login ulang sehingga melewati batas 2 sesi, sesi paling lama (yang kemungkinan milik korban) akan terkena eviction. Admin dapat memakai endpoint force logout (`POST /admin/users/:id/force-logout`) untuk mencabut seluruh sesi pengguna tersebut seketika, lalu korban perlu mengganti kata sandi.

**4. Dua kekurangan yang disadari.**
(a) Belum ada rate limiting atau penguncian akun pada endpoint login, sehingga rentan terhadap brute force. (b) Token disimpan di sisi klien (localStorage atau cookie tanpa proteksi lengkap), sehingga rentan terhadap XSS atau CSRF, dan belum ada rotasi token atau deteksi perangkat mencurigakan. **[SESUAIKAN dengan kekurangan yang benar-benar ada]**

**5. Fitur "keluar dari seluruh perangkat".**
Saya akan menambahkan endpoint `POST /auth/logout-all` di `internal/handler/auth.go`, yang memanggil fungsi baru `RevokeAllSessions(userID)` di `internal/auth/service.go`. Fungsi itu sudah hampir sama dengan logika force logout admin, yaitu menghapus atau menandai `revoked` semua baris `sessions` milik `user_id` yang sedang login, sehingga cukup memakai ulang query yang sama.
