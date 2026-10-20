# PinjamKampus API

- Nama dan NRP: Rafian Dany Azadirahta (5053241008)
- Persona: D
- Cara menjalankan: `go run .` lalu buka http://localhost:8080

## Fitur
- [x] B1 Sumber daya
- [x] B2 Buat peminjaman
- [x] B3 Lihat peminjaman (daftar, filter, detail)
- [x] B4 Validasi dan error JSON
- [x] B5 Deteksi bentrok (409)
- [ ] B6 Ubah dan hapus
- [ ] B7 Aman untuk permintaan bersamaan
- [ ] Fitur khas persona
- [ ] Change request

## Yang belum selesai
- B6 Ubah dan hapus: belum dikerjakan (PUT dan DELETE /peminjaman/{id} belum ada).
- B7 Aman untuk permintaan bersamaan: semua akses data sudah dilindungi satu sync.Mutex, dan pengecekan bentrok serta penyimpanan berada dalam satu kunci. Tetapi belum saya uji dengan banyak request bersamaan, jadi belum saya centang.
- Fitur khas persona D (POST dan DELETE /paket-peminjaman): belum dikerjakan. Field paket_id sudah ada di objek peminjaman (selalu 0).
- Change request: belum dikerjakan.
