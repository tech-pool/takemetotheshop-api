package registration

import (
	"fmt"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "%s s-print %q q-print %v v-print \n", "Hello Registration", "Hallo Registration", "Hola Registration")

	fmt.Fprintf(w, "URL.Path is: %q\n", r.URL.Path)
	fmt.Fprintf(w, "r.RequestURI is: %s\n", r.RequestURI)
	fmt.Fprintf(w, "r.URL.Path is: %s\n", r.URL.Path)
	fmt.Fprintf(w, "r.URL is: %s\n", r.URL)
	fmt.Fprintf(w, "r.URL.User is: %s\n", r.URL.User)
	fmt.Fprintf(w, "r.URL.Host is: %s\n", r.URL.Host)
}
