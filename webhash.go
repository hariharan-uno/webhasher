package main

import (
	"fmt"
	"net/http"
)

func MyHandler(w http.ResponseWriter, r *http.Request) {
	q := r.FormValue("q")
	fmt.Fprintf(w, "%s", q)
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
</form>
</body>
</html>`

func main() {
	http.HandleFunc("/", InputHandler)
	http.HandleFunc("/hash", MyHandler)
	http.ListenAndServe(":8080", nil)
}
