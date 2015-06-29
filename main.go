package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

type CookieResponse struct {
	CookieName string
	Success    bool
}

const (
	cookieNameKey = "cookie_name"
)

func setCookie(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
	}
	cName := r.FormValue(cookieNameKey)
	maxAge := 10
	c := &http.Cookie{
		Name:    cName,
		Value:   "I'm a unique id!", // Probably should be a random string verified during checkcookie.
		MaxAge:  maxAge,
		Domain:  r.URL.Host,
		Expires: time.Now().Add(time.Duration(10 * time.Second)),
	}
	http.SetCookie(w, c)
	marshalAndWrite(w, CookieResponse{CookieName: cName, Success: true})
}

func checkCookie(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
	}
	cName := r.FormValue(cookieNameKey)
	c, err := r.Cookie(cName)
	if err == nil && c != nil {
		c.MaxAge = -1
		http.SetCookie(w, c)
	}
	resp := CookieResponse{CookieName: cName, Success: true}
	if err != nil || c == nil {
		resp.Success = false
	}
	if err = marshalAndWrite(w, resp); err != nil {
		log.Println(err)
	}
}

func marshalAndWrite(w http.ResponseWriter, data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	total := 0
	for total < len(b) {
		i, err := w.Write(b)
		if err != nil {
			return err
		}
		total += i
	}
	return nil
}

func main() {
	port := flag.Int("port", 8080, "Port to listen for connections on.")
	flag.Parse()

	http.HandleFunc("/checkcookie", checkCookie)
	http.HandleFunc("/setcookie", setCookie)
	fmt.Printf("Starting server on port %d..\n", *port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		log.Fatal(err)
	}
}
