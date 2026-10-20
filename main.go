package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"unicode/utf8"
)

type SumberDaya struct {
	ID    int    `json:"id"`
	Nama  string `json:"nama"`
	Jenis string `json:"jenis"`
	Stok  int    `json:"stok"`
}

type store struct {
	mu         sync.Mutex
	sumberDaya []SumberDaya
}

func newStore() *store {
	return &store{sumberDaya: make([]SumberDaya, 0)}
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

func NewHandler() http.Handler {
	s := newStore()
	mux := http.NewServeMux()
	mux.HandleFunc("POST /sumber-daya", s.buatSumberDaya)
	mux.HandleFunc("GET /sumber-daya", s.daftarSumberDaya)
	mux.HandleFunc("GET /sumber-daya/{id}", s.detailSumberDaya)
	return mux
}

func main() {
	log.Fatal(http.ListenAndServe(":8080", NewHandler()))
}