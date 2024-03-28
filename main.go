package main

import (
	"net/http"
)

func main() {
	assets := http.FileServer(http.Dir("webtemplate/assets/"))
	http.Handle("/assets/", http.StripPrefix("/assets/", assets))

	pages := http.FileServer(http.Dir("webtemplate/pages/"))
	http.Handle("/pages/", http.StripPrefix("/pages/", pages))

	http.ListenAndServe(":8181", nil)
}
