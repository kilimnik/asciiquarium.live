package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/creack/pty"
)

func main() {
	http.HandleFunc("/", handle)

	s := &http.Server{
		Addr:           ":8080",
		ReadTimeout:    0,
		WriteTimeout:   0,
		MaxHeaderBytes: 0,
	}
	log.Fatal(s.ListenAndServe())
}

func parseUInt(query url.Values, key string, def uint64) (uint64, error) {
	if !query.Has(key) {
		return def, nil
	}

	val, err := strconv.ParseUint(query.Get(key), 10, 16)
	if err != nil {
		return 0, err
	}

	return val, nil
}

func handle(w http.ResponseWriter, r *http.Request) {
	userAgent := r.Header.Get("User-Agent")
	if !strings.Contains(userAgent, "curl") {
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, "Use curl to see the page ;)")

		return
	}

	query := r.URL.Query()

	cols, err := parseUInt(query, "cols", 100)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "cols: "+err.Error())

		return
	}

	rows, err := parseUInt(query, "rows", 30)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "rows: "+err.Error())

		return
	}

	c := exec.CommandContext(r.Context(), "asciiquarium")
	c.Env = []string{"TERM=xterm-256color"}
	defer func() {
		c.Process.Kill()
		c.Process.Wait()
	}()

	size := &pty.Winsize{
		Cols: uint16(cols),
		Rows: uint16(rows),
	}

	f, err := pty.StartWithSize(c, size)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "Interal Server Error: "+err.Error())

		return
	}
	defer f.Close()

	w.Header().Set("Transfer-Encoding", "chunked")
	w.WriteHeader(http.StatusOK)

	go func() {
		for {
			_, err := f.Write([]byte("."))
			if err != nil {
				return
			}

			time.Sleep(50 * time.Millisecond)
		}
	}()

	io.Copy(w, f)
}
