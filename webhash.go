package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"github.com/gorilla/schema" //For populating struct with multiple url values
	"html/template"
	"io"
	"net/http"
)

// Hasher type contains the multiple url values, Query and Format.
// Currently the formats supported are only md5, sha1, sha256.
//
// IMPORTANT: Make sure that fields are exported, i.e. Capitalized first letters.
type Hasher struct {
	Query  string `schema:"q"`      //query
	Format string `schema:"format"` //format of hash (md5, sha1, sha256)
}

var hasher = new(Hasher) //Returns a pointer to a new Hasher type
var decoder = schema.NewDecoder()
var templates = template.Must(template.ParseFiles("index.html"))

// HashHandler parses through the url values and determines the query string and
// format type parameters. It checks hasher.Format and writes the corresponding hash
// to the http.ResponseWriter.
// Currently, It suports only MD5, SHA1, SHA256.
func HashHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = decoder.Decode(hasher, r.Form)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if hasher.Query == "" || hasher.Format == "" {
		http.Redirect(w, r, "/", http.StatusFound) //If the query or the format is empty, redirect to the home page.
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
		http.Redirect(w, r, "/", http.StatusFound) //If the format is not supported, redirect to the home page.
		return
	}
}

// InputHandler returns a HTML form for query string input and selecting hash type
func InputHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index")
}

// renderTemplate renders the template and handles errors.
// It takes http.Response Writer and the template filename as inputs.
func renderTemplate(w http.ResponseWriter, tmpl string) {
	err := templates.ExecuteTemplate(w, tmpl+".html", "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", InputHandler)
	http.HandleFunc("/hash", HashHandler)
	http.ListenAndServe(":8080", nil)
}
