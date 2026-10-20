# AI-LOG

Catat setiap penggunaan AI yang menghasilkan kode. Jika tidak memakai AI sama sekali, tulis satu baris: **tidak memakai AI**.

Log yang jujur tidak mengurangi nilai, sebanyak apa pun AI dipakai. Yang mengurangi nilai adalah log yang tidak sesuai dengan kode, atau isi log yang tidak bisa dijelaskan.

## Sumber daya (B1)
- Prompt: minta kode B1 (buat, daftar, detail sumber daya) dengan Go net/http
- Dipakai: struct SumberDaya, validasiSumberDaya, handler buat/daftar/detail di main.go
- Diubah: versi pertama memakai input struct dengan pointer dan unlock manual. Saya minta disederhanakan, sehingga decode langsung ke SumberDaya (karena nilai kosong sudah tidak valid, hasilnya tetap 400) dan memakai defer untuk unlock. Hasilnya saya uji dengan curl, semua sesuai kontrak.

## Buat peminjaman (B2)
- Prompt: minta kode B2 (POST /peminjaman) dengan Go net/http, hanya menyimpan data dan memberi id serta status aktif.
- Dipakai: struct Peminjaman, tambahPeminjaman, buatPeminjaman, dan route POST /peminjaman di main.go.
- Diubah: versi awal punya field nextPaketID di store sementara fitur paket belum dikerjakan. Saya hapus dan akan menambahkannya saat mengerjakan fitur paket. Hasilnya saya uji dengan curl: 201, id 1, status "aktif", paket_id 0.

## Lihat peminjaman (B3)
- Prompt: minta kode B3 (GET /peminjaman dengan filter sumber_daya_id, tanggal, status dan GET /peminjaman/{id}) dengan Go net/http.
- Dipakai: daftarPeminjaman, detailPeminjaman, dan dua route GET di main.go.
- Diubah: tidak ada perubahan yang dilakukan, kode digunakan apa adanya.

## Validasi peminjaman (B4)
- Prompt: minta kode B4 (validasi POST /peminjaman: peminjam, tanggal, slot, sumber_daya_id, jumlah) dengan Go net/http.
- Dipakai: validasiPeminjaman, tambahPeminjaman, buatPeminjaman di main.go.
- Diubah: versi awal punya dua pengecekan slot terpisah (batas 1 sampai 8, lalu slot_mulai tidak boleh lebih besar dari slot_selesai). Saya minta disederhanakan menjadi satu kondisi: slot_mulai >= 1, slot_selesai <= 8, dan slot_mulai <= slot_selesai, karena dua batas lainnya otomatis terpenuhi. Tanpa pointer, sehingga field yang tidak dikirim bernilai 0 atau kosong dan tetap gagal validasi (400). Hasilnya saya uji dengan curl.