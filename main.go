package main

import (
	"io"
	"log"
	"net/http"
	"strings"
	"os/exec"
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

func handle(w http.ResponseWriter, r *http.Request) {
	userAgent := r.Header.Get("User-Agent")
	if !strings.Contains(userAgent, "curl") {
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, "Use curl to see the page ;)")

		return
	}

	c := exec.CommandContext(r.Context(), "asciiquarium")
	c.Env = []string{"TERM=xterm-256color"}
	defer func ()  {
		c.Process.Kill()
		c.Process.Wait()
	}()

	size := &pty.Winsize{
		Cols: 100,
		Rows: 30,
	}

	f, err := pty.StartWithSize(c, size)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "Interal Server Error: " + err.Error())

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
