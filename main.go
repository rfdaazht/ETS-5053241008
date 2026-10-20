package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

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
}

func newStore() *store {
	return &store{
		sumberDaya:       make([]SumberDaya, 0),
		peminjaman:       make([]Peminjaman, 0),
		nextPeminjamanID: 1,
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, pesan string) {
	writeJSON(w, status, map[string]string{"error": pesan})
}

func validasiSumberDaya(sd *SumberDaya) string {
	sd.Nama = strings.TrimSpace(sd.Nama)
	sd.Jenis = strings.TrimSpace(sd.Jenis)
	if sd.Nama == "" {
		return "nama wajib diisi"
	}
	if utf8.RuneCountInString(sd.Nama) > 50 {
		return "nama maksimal 50 karakter"
	}
	if sd.Jenis != "ruang" && sd.Jenis != "alat" {
		return "jenis harus ruang atau alat"
	}
	if sd.Stok < 1 {
		return "stok minimal 1"
	}
	return ""
}

func validasiPeminjaman(p *Peminjaman) string {
	p.Peminjam = strings.TrimSpace(p.Peminjam)
	p.Tanggal = strings.TrimSpace(p.Tanggal)
	if p.Peminjam == "" {
		return "peminjam wajib diisi"
	}
	if utf8.RuneCountInString(p.Peminjam) > 50 {
		return "peminjam maksimal 50 karakter"
	}
	if _, err := time.Parse("2006-01-02", p.Tanggal); err != nil {
		return "tanggal harus berformat YYYY-MM-DD dan tanggal yang valid"
	}
	if p.SlotMulai < 1 || p.SlotSelesai > 8 || p.SlotMulai > p.SlotSelesai {
		return "slot harus 1 sampai 8 dan slot_mulai tidak boleh lebih besar dari slot_selesai"
	}
	return ""
}

func (s *store) terpakai(sumberDayaID int, tanggal string, slot int) int {
	total := 0
	for _, p := range s.peminjaman {
		if p.Status == "aktif" && p.SumberDayaID == sumberDayaID &&
			p.Tanggal == tanggal && p.SlotMulai <= slot && slot <= p.SlotSelesai {
			total += p.Jumlah
		}
	}
	return total
}

func (s *store) tambahSumberDaya(sd SumberDaya) SumberDaya {
	s.mu.Lock()
	defer s.mu.Unlock()
	sd.ID = len(s.sumberDaya) + 1
	s.sumberDaya = append(s.sumberDaya, sd)
	return sd
}

func (s *store) buatSumberDaya(w http.ResponseWriter, r *http.Request) {
	var sd SumberDaya
	if err := json.NewDecoder(r.Body).Decode(&sd); err != nil {
		writeError(w, http.StatusBadRequest, "body tidak valid")
		return
	}
	if pesan := validasiSumberDaya(&sd); pesan != "" {
		writeError(w, http.StatusBadRequest, pesan)
		return
	}
	writeJSON(w, http.StatusCreated, s.tambahSumberDaya(sd))
}

func (s *store) daftarSumberDaya(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	writeJSON(w, http.StatusOK, s.sumberDaya)
}

func (s *store) detailSumberDaya(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "id harus berupa angka")
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if id < 1 || id > len(s.sumberDaya) {
		writeError(w, http.StatusNotFound, "sumber daya tidak ditemukan")
		return
	}
	writeJSON(w, http.StatusOK, s.sumberDaya[id-1])
}

func (s *store) tambahPeminjaman(p Peminjaman) (Peminjaman, int, string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if p.SumberDayaID < 1 || p.SumberDayaID > len(s.sumberDaya) {
		return Peminjaman{}, http.StatusBadRequest, "sumber_daya_id tidak ditemukan"
	}
	stok := s.sumberDaya[p.SumberDayaID-1].Stok
	if p.Jumlah < 1 || p.Jumlah > stok {
		return Peminjaman{}, http.StatusBadRequest, "jumlah harus antara 1 dan stok sumber daya"
	}
	for slot := p.SlotMulai; slot <= p.SlotSelesai; slot++ {
		if s.terpakai(p.SumberDayaID, p.Tanggal, slot)+p.Jumlah > stok {
			return Peminjaman{}, http.StatusConflict, "kapasitas terlampaui pada salah satu slot"
		}
	}
	p.ID = s.nextPeminjamanID
	s.nextPeminjamanID++
	p.Status = "aktif"
	p.PaketID = 0
	s.peminjaman = append(s.peminjaman, p)
	return p, http.StatusCreated, ""
}

func (s *store) buatPeminjaman(w http.ResponseWriter, r *http.Request) {
	var p Peminjaman
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeError(w, http.StatusBadRequest, "body tidak valid")
		return
	}
	if pesan := validasiPeminjaman(&p); pesan != "" {
		writeError(w, http.StatusBadRequest, pesan)
		return
	}
	hasil, status, pesan := s.tambahPeminjaman(p)
	if pesan != "" {
		writeError(w, status, pesan)
		return
	}
	writeJSON(w, status, hasil)
}

func (s *store) daftarPeminjaman(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	adaSumberDaya := q.Has("sumber_daya_id")
	sumberDayaID := 0
	if adaSumberDaya {
		id, err := strconv.Atoi(q.Get("sumber_daya_id"))
		if err != nil {
			writeError(w, http.StatusBadRequest, "sumber_daya_id harus berupa angka")
			return
		}
		sumberDayaID = id
	}
	tanggal := q.Get("tanggal")
	status := q.Get("status")

	s.mu.Lock()
	defer s.mu.Unlock()
	hasil := make([]Peminjaman, 0)
	for _, p := range s.peminjaman {
		if adaSumberDaya && p.SumberDayaID != sumberDayaID {
			continue
		}
		if tanggal != "" && p.Tanggal != tanggal {
			continue
		}
		if status != "" && p.Status != status {
			continue
		}
		hasil = append(hasil, p)
	}
	writeJSON(w, http.StatusOK, hasil)
}

func (s *store) detailPeminjaman(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "id harus berupa angka")
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, p := range s.peminjaman {
		if p.ID == id {
			writeJSON(w, http.StatusOK, p)
			return
		}
	}
	writeError(w, http.StatusNotFound, "peminjaman tidak ditemukan")
}

func NewHandler() http.Handler {
	s := newStore()
	mux := http.NewServeMux()
	mux.HandleFunc("POST /sumber-daya", s.buatSumberDaya)
	mux.HandleFunc("GET /sumber-daya", s.daftarSumberDaya)
	mux.HandleFunc("GET /sumber-daya/{id}", s.detailSumberDaya)
	mux.HandleFunc("POST /peminjaman", s.buatPeminjaman)
	mux.HandleFunc("GET /peminjaman", s.daftarPeminjaman)
	mux.HandleFunc("GET /peminjaman/{id}", s.detailPeminjaman)
	return mux
}

func main() {
	log.Fatal(http.ListenAndServe(":8080", NewHandler()))
}