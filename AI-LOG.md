# AI-LOG

Catat setiap penggunaan AI yang menghasilkan kode. Jika tidak memakai AI sama sekali, tulis satu baris: **tidak memakai AI**.

Log yang jujur tidak mengurangi nilai, sebanyak apa pun AI dipakai. Yang mengurangi nilai adalah log yang tidak sesuai dengan kode, atau isi log yang tidak bisa dijelaskan.

## Sumber daya (B1)
- Prompt: minta kode B1 (buat, daftar, detail sumber daya) dengan Go net/http
- Dipakai: struct SumberDaya, validasiSumberDaya, handler buat/daftar/detail di main.go
- Diubah: versi pertama memakai input struct dengan pointer dan unlock manual. Saya minta disederhanakan, sehingga decode langsung ke SumberDaya (karena nilai kosong sudah tidak valid, hasilnya tetap 400) dan memakai defer untuk unlock. Hasilnya saya uji dengan curl, semua sesuai kontrak.

## (Entri berikutnya)
- Prompt:
- Dipakai:
- Diubah:
