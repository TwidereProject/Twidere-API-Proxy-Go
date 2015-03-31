package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

var DOMAIN_FORMAT_PREFIX = "/domain."

func main() {
	http.HandleFunc("/", handleRequest)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	host, port := os.Getenv("HOST"), os.Getenv("PORT")
	var bind string
	if host == "" || port == "" {
		bind = fmt.Sprintf("%s:%s", "127.0.0.1", "18080")
	} else {
		bind = fmt.Sprintf("%s:%s", host, port)
	}
	err := http.ListenAndServe(bind, nil)
	if err != nil {
		panic(err)
	}
}

func handleRequest(res http.ResponseWriter, req *http.Request) {
	url, uri := req.URL, req.RequestURI
	if url.Path == "/" {
		handleWelcome(res, req)
	} else if strings.HasPrefix(uri, DOMAIN_FORMAT_PREFIX) {
		regulatedUri := uri[len(DOMAIN_FORMAT_PREFIX):]
		slashIdx := strings.Index(regulatedUri, "/")
		if slashIdx < 0 {
			handleTwitterRequest(res, req, regulatedUri, "/")
		} else {
			handleTwitterRequest(res, req, regulatedUri[:slashIdx], regulatedUri[slashIdx:])
		}
	} else {
		handleUnimplementedRequest(res, req)
	}
}

func handleTwitterRequest(res http.ResponseWriter, req *http.Request, domain string, uri string) {
	client := http.Client{}
	var twitterUri string
	if domain == "" {
		twitterUri = fmt.Sprintf("https://twitter.com%s", uri)
	} else {
		twitterUri = fmt.Sprintf("https://%s.twitter.com%s", domain, uri)
	}
	twitterReq, err := http.NewRequest(req.Method, twitterUri, req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	reqHeader := req.Header
	copyHeader(reqHeader, twitterReq.Header)
	reqHeader.Del("Cookie")
	cookieDomain := fmt.Sprintf(".%s", req.Host)
	for _, v := range req.Cookies() {
		if strings.EqualFold(v.Domain, cookieDomain) {
			v.Domain = ".twitter.com"
		}
		reqHeader.Add("Cookie", v.String())
	}
	twitterRes, err := client.Do(twitterReq)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	resHeader := res.Header()
	copyHeader(twitterRes.Header, res.Header())
	resHeader.Del("Set-Cookie")
	for _, v := range twitterRes.Cookies() {
		if strings.EqualFold(v.Domain, ".twitter.com") {
			v.Domain = cookieDomain
		}
		resHeader.Add("Set-Cookie", v.String())
	}
	res.WriteHeader(twitterRes.StatusCode)
	io.Copy(res, twitterRes.Body)
}

func handleUnimplementedRequest(res http.ResponseWriter, req *http.Request) {
	errMessage := fmt.Sprintf("Unable to handle: %s, not implemented.", req.URL.Path)
	http.Error(res, errMessage, http.StatusNotImplemented)
}

func handleWelcome(res http.ResponseWriter, req *http.Request) {
	info := WelcomeInfo{"Twidere API Proxy", "0.9", req.Host}
	fp := path.Join("templates", "welcome.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(res, info); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func copyHeader(in http.Header, out http.Header) {
	for k, vs := range in {
		for _, v := range vs {
			out.Add(k, v)
		}
	}
}

type WelcomeInfo struct {
	Name    string
	Version string
	Host    string
}
