package main

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	log.Printf("Start gif2web server on 8080")
	http.Handle("/convert", gifProcessor())
	http.ListenAndServe("0.0.0.0:8080", nil)
}

func gifProcessor() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if err := req.ParseMultipartForm((1 << 10) * 24); nil != err {
			log.Printf("Error while parse: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		for _, fheaders := range req.MultipartForm.File {
			for _, hdr := range fheaders {
				log.Printf("Income file len: %d", hdr.Size)

				var err error
				var infile multipart.File
				if infile, err = hdr.Open(); err != nil {
					log.Printf("[ERROR] Handle open error: %v", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				quality := req.URL.Query().Get("quality")
				if quality == "" {
					quality = "0"
				}
				q, err := strconv.ParseInt(quality, 10, 64)
				if err != nil || q < 0 || q > 100 {
					log.Printf("[ERROR] Bad quality params: %s %v", err, q)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				cmd := exec.Command("gif2webp", "-mixed", "-q", strconv.Itoa(int(q)), "-o", "-", "--", "-")

				cmd.Stdin = io.Reader(infile)

				var out bytes.Buffer
				cmd.Stdout = &out

				var errout bytes.Buffer
				cmd.Stderr = &errout

				err = cmd.Run()
				if err != nil {
					log.Printf("[ERROR] Webp Output: %v, %v, %v\n", err, out.String(), errout.String())
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				output := out.String()
				reader := strings.NewReader(output)
				imageLen := reader.Len()
				log.Printf("Outcome file len: %d", reader.Len())
				w.Header().Set("Content-Type", "image/webp")
				w.Header().Set("Content-Length", strconv.Itoa(imageLen))
				w.WriteHeader(http.StatusOK)
				io.Copy(w, reader)
				return
			}
		}

		w.WriteHeader(http.StatusInternalServerError)
	})
}
