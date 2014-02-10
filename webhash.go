package main

import (
	"fmt"
	"github.com/gorilla/schema"
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
		panic(err)
	}
	err = decoder.Decode(hasher, r.Form)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, "<html>"+"Query: %s"+"<br>"+"Format: %s"+"</html>", hasher.Query, hasher.Format)
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
