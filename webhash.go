package main

import (
	"crypto/md5"
	"fmt"
	"github.com/gorilla/schema"
	"io"
	"net/http"
)

type Hasher struct {
	Query  string `schema:"q"`      //query
	Format string `schema:"format"` //format of hash
}

var hasher = new(Hasher)
var decoder = schema.NewDecoder()

func MyHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = decoder.Decode(hasher, r.Form)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if hasher.Query == "" || hasher.Format == "" {
		fmt.Fprint(w, InputForm)
		return
	}
	switch hasher.Format {
	case "md5":
		h := md5.New()
		io.WriteString(h, hasher.Query)
		fmt.Fprintf(w, "%x", h.Sum(nil))
		return
	case "sha1":
		fmt.Fprintf(w, "sha1")
		return
	case "sha256":
		fmt.Fprintf(w, "sha256")
		return
	default:
		fmt.Fprintf(w, "Shit, not supported")
		return
	}
}

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
	http.HandleFunc("/hash", MyHandler)
	http.ListenAndServe(":8080", nil)
}
