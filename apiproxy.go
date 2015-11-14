package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

var DOMAIN_FORMAT_PREFIX = "/domain."

func main() {
	http.HandleFunc("/", handleRequest)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	host, port := os.Getenv("HOST"), os.Getenv("PORT")
	var bind string
	if host == "" {
		bind = fmt.Sprintf(":%s", port)
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
	roundTripper := http.RoundTripper{}
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
		reqHeader.Add("Cookie", cookieToString(v))
	}
	twitterRes, err := roundTripper.RoundTrip(twitterReq)
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
		resHeader.Add("Set-Cookie", cookieToString(v))
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

func cookieToString(c *http.Cookie) string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "%s=%s", c.Name, c.Value)
	if len(c.Path) > 0 {
		fmt.Fprintf(&b, "; Path=%s", c.Path)
	}
	if len(c.Domain) > 0 {
		fmt.Fprintf(&b, "; Domain=%s", c.Domain)
	}
	if c.Expires.Unix() > 0 {
		fmt.Fprintf(&b, "; Expires=%s", c.Expires.UTC().Format(time.RFC1123))
	}
	if c.MaxAge > 0 {
		fmt.Fprintf(&b, "; Max-Age=%d", c.MaxAge)
	} else if c.MaxAge < 0 {
		fmt.Fprintf(&b, "; Max-Age=0")
	}
	if c.HttpOnly {
		fmt.Fprintf(&b, "; HttpOnly")
	}
	if c.Secure {
		fmt.Fprintf(&b, "; Secure")
	}
	return b.String()
}

type WelcomeInfo struct {
	Name    string
	Version string
	Host    string
}
