# Project Title

Base Project Go

## Description

Aplikasi web menggunakan framework [ITS Go](https://github.com/dptsi/its-go) yang dapat digunakan sebagai rujukan atau basis pengembangan aplikasi web yang lain.

## Getting Started

### Dependencies

- Disarankan menggunakan WSL2 dengan distro Ubuntu 22.04
- [Go](https://go.dev/doc/install)

### Installing

1. Clone repository

   Repository ini menggunakan git submodules sehingga terdapat sedikit perbedaan pada cara clone pada umumnya. Perbedaan tersebut berupa adanya argumen `--recurse-submodules`:

   ```bash
   git clone --recurse-submodules git@bitbucket.org:dptsi/base-go.git
   ```

### Executing program

1. Copy .env.example ke .env
   ```bash
   cp .env.example .env
   ```
2. Set up environment variable untuk OpenID Connect (Bisa didapatkan dari Secman)
   ```bash
   # Dapat dari Secman
   OIDC_PROVIDER=
   OIDC_CLIENT_ID=
   OIDC_CLIENT_SECRET=
   OIDC_REDIRECT_URL=
   OIDC_POST_LOGOUT_REDIRECT_URI=
   # Scope minimal yang diperlukan adalah "openid email profile resource role"
   OIDC_SCOPES="openid email profile resource role"
   ```
3. Generate key
   ```bash
   go run script/script.go key:generate
   ```
4. Atur konfigurasi koneksi database default
   ```bash
   DB_DRIVER=
   DB_HOST=
   DB_PORT=
   DB_DATABASE=
   DB_USER=
   DB_PASSWORD=
   # Jika belum ada database dan ingin mencoba2 dapat menggunakan driver sqlite
   DB_DRIVER=sqlite
   DB_HOST=
   DB_PORT=
   DB_DATABASE=./db.sqlite # File ini akan dibuat secara otomatis
   DB_USER=
   DB_PASSWORD=
   ```
5. Jalankan server.
   ```bash
   # Menjalankan server pada port 8080 (secara default)
   go run main.go
   # Menjalankan server pada port lain
   PORT=1111 go run main.go
   ```
6. Dokumentasi API dapat diakses melalui `/swagger/index.html` (Pastikan `APP_ENV` bernilai `local`).

## Help

1. Jika mendapatkan error `no such table: sessions`,

   ```json
   {
     "code": 9005,
     "data": {
       "error": "no such table: sessions",
       "request_id": "4d7fc0dd-2ae9-4a71-9512-6512398c8e19"
     },
     "message": "internal_server_error"
   }
   ```

   maka pastikan terdapat tabel database untuk menyimpan session sesuai konfigurasi `config/sessions.go` (secara default, framework akan membaca dari tabel sessions). Jika tidak ada, maka perlu tabel dengan struktur berikut:

   ```text
   id          -> text
   data        -> text
   expired_at  -> timestamp
   csrf_token  -> text
   ```

   Namun **jika administrator database mengizinkan aplikasi untuk memigrasikan tabel session sendiri**, maka cukup aktifkan auto migration pada `config/sessions.go`.

   ```go
   package config

   import "github.com/dptsi/its-go/sessions"

   func sessionsConfig() sessions.Config {
   	return sessions.Config{
   		Storage:    "database",
   		Connection: "default",
   		Table:      "sessions",
   		Cookie: sessions.CookieConfig{
   			Name:           "myits_academics_session",
   			CsrfCookieName: "CSRF-TOKEN",
   			Path:           "/",
   			Domain:         "",
   			Secure:         false,
   			Lifetime:       60,
   		},
   		AutoMigrate: true, // tambahkan line ini
   	}
   }
   ```

## Authors

- [@zydhanlinnar11](https://github.com/zydhanlinnar11)

## License

This project is licensed under the MIT License - see the [LICENSE file](./LICENSE) for details

## Acknowledgments

- [Laravel](https://laravel.com)
- [Base Laravel DPTSI](https://bitbucket.org/dptsi/base-laravel)
- many others...
