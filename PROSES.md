# PROSES

## 1. Persona dan masalahnya
Pengguna adalah panitia acara yang menyelenggarakan acara dan butuh aula, proyektor, dan sound system sekaligus. Jika salah satu saja dari peralatan tersebut sudah dipinjam orang lain, maka acara tidak bisa jalan. Fitur khas memastikan saat ada peminjaman, ketiga alat tersebut dipinjam bersamaan atau tidak dipinjam sama sekali (peminjaman paket).

## 2. Rancangan struktur data
type SumberDaya struct {
    ID    int    `json:"id"`
    Nama  string `json:"nama"`
    Jenis string `json:"jenis"`
    Stok  int    `json:"stok"`
}

type Peminjaman struct {
    ID           int    `json:"id"`
    SumberDayaID int    `json:"sumber_daya_id"`
    Peminjam     string `json:"peminjam"`
    Tanggal      string `json:"tanggal"`
    SlotMulai    int    `json:"slot_mulai"`
    SlotSelesai  int    `json:"slot_selesai"`
    Jumlah       int    `json:"jumlah"`
    Status       string `json:"status"`
    PaketID      int    `json:"paket_id"`
}

type store struct {
    mu               sync.Mutex
    sumberDaya       []SumberDaya
    peminjaman       []Peminjaman
    nextPeminjamanID int
    nextPaketID      int
}

Alasan:
- Data sumber daya dan peminjaman disimpan dalam slice. Id selalu naik, jadi urutannya otomatis dari id terkecil. Untuk daftar kosong bisa dengan `make` agar hasilnya [], bukan null.
- Id sumber daya dibuat dengan len(sumberDaya)+1. Karena sumber daya tidak pernah dihapus seharusnya aman.
- Id peminjaman dan id paket menggunakan counter masing-masing. Karena peminjaman bisa dihapus, maka len+1 id lama bisa terpakai lagi. Counter baru dinaikkan setelah semua pengecekan lolos, jadi paket yang gagal tidak akan menghabiskan nomor id.
- Jika ingin mengetahui berapa slot yang sudah terpakai, bisa hitung dari peminjaman yang berstatus aktif. Dengan begitu paket dan peminjaman biasa memakai data yang sama.
- Peminjaman milik paket dikenali dari PaketID (0 berarti peminjaman biasa). Saat menghapus paket, cari semua peminjaman dengan PaketID itu, lalu dihapus. Kalau tidak ada satu pun, hasilnya 404.
- Satu mutex di dalam store untuk mengunci semua data. Store dibuat baru setiap NewHandler() dipanggil, jadi selalu mulai dari kosong.

## 3. Satu keputusan desain
(Satu keputusan yang kamu ambil, pilihan lain yang kamu pertimbangkan, dan alasanmu memilih yang ini. Contoh: kenapa cek bentrok per slot, bukan per rentang.)

## 4. Bug atau kendala (minimal dua)
### Bug 1
- Gejala:
- Cara menemukan:
- Penyebab dan solusi:

### Bug 2
- Gejala:
- Cara menemukan:
- Penyebab dan solusi:

## 5. Change request
(Apa yang diminta, bagian kode mana yang kamu ubah, dan satu keputusan desain yang kamu ambil.)
