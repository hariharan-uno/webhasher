package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"github.com/gorilla/schema" //For populating struct with multiple url values
	"io"
	"net/http"
)

type Hasher struct {
	Query  string `schema:"q"`      //query
	Format string `schema:"format"` //format of hash
}

var hasher = new(Hasher) //Returns a pointer to a new Hasher type
var decoder = schema.NewDecoder()

//HashHandler parses through the url values and determines the query string and
//format type parameters. It checks hasher.Format and writes the corresponding hash
//to the http.ResponseWriter.
//Currently, It suports only MD5, SHA1, SHA256.
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
		fmt.Fprint(w, InputForm) //If the query string or format is empty, it writes the input form
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
		fmt.Fprintf(w, "Shit, not supported")
		return
	}
}

//InputHandler returns a HTML form for query string input and selecting hash type
func InputHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, InputForm)
}

const InputForm = `<html>
<body>
<form method="GET" action="/hash">
<label>
Type the text you want to convert: 
<input type="text" name="q" />
</label>
<select name="format">
<option value="md5">MD5</option>
<option value="sha1">SHA1</option>
<option value="sha256">SHA256</option>
</select>
<button type="submit">Go</button>
</form>
</body>
</html>`

func main() {
	http.HandleFunc("/", InputHandler)
	http.HandleFunc("/hash", HashHandler)
	http.ListenAndServe(":8080", nil)
}
