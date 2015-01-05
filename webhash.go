// Copyright 2014 Hari haran. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/gorilla/schema"
)

// Hasher type contains the multiple url values, Query and Format.
// Make sure that fields are exported, i.e. capitalized first letters.
type Hasher struct {
	Query  string `schema:"q"`
	Format string `schema:"format"` // format of hash e.g. md5, sha1, etc.
}

var templates = template.Must(template.ParseFiles("index.html"))

// HashHandler parses through the url values and determines the query string and
// format type parameters. It checks hasher.Format and writes the corresponding hash
// to the http.ResponseWriter.
// Currently, It suports only MD5, SHA1, SHA256.
func HashHandler(w http.ResponseWriter, r *http.Request) {
	hasher := new(Hasher)
	decoder := schema.NewDecoder()

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := decoder.Decode(hasher, r.Form); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if hasher.Query == "" || hasher.Format == "" {
		// If the query or the format is empty, redirect to the home page.
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	switch hasher.Format {
	case "md5":
		h := md5.New()
		io.WriteString(h, hasher.Query)
		fmt.Fprintf(w, "%x", h.Sum(nil))
		return
	case "sha1":
		h := sha1.New()
		io.WriteString(h, hasher.Query)
		fmt.Fprintf(w, "%x", h.Sum(nil))
		return
	case "sha256":
		h := sha256.New()
		io.WriteString(h, hasher.Query)
		fmt.Fprintf(w, "%x", h.Sum(nil))
		return
	default:
		// If the format is not supported, redirect to the home page.
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
}

// InputHandler returns a HTML form for input.
func InputHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index")
}

// renderTemplate renders the template and handles errors.
func renderTemplate(w http.ResponseWriter, tmpl string) {
	buf := new(bytes.Buffer)
	err := templates.ExecuteTemplate(buf, tmpl+".html", "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.Copy(w, buf)
}

func init() {
	http.HandleFunc("/", InputHandler)
	http.HandleFunc("/hash", HashHandler)
}
