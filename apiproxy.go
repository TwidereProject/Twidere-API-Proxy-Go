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

const DOMAIN_FORMAT_PREFIX = "/domain."

func main() {
	http.HandleFunc("/", handleRequest)
	bind := fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))
	err := http.ListenAndServe(bind, nil)
	if err != nil {
		panic(err)
	}
}

func handleRequest(res http.ResponseWriter, req *http.Request) {
	uri := req.RequestURI
	if uri == "/" {
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

	}
	for k, vs := range req.Header {
		for _, v := range vs {
			twitterReq.Header.Add(k, v)
		}
	}
	twitterRes, err := client.Do(twitterReq)
	if err != nil {

	}
	resHeader := res.Header()
	for k, vs := range twitterRes.Header {
		for _, v := range vs {
			resHeader.Add(k, v)
		}
	}
	res.WriteHeader(twitterRes.StatusCode)
	io.Copy(res, twitterRes.Body)
}

func handleUnimplementedRequest(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(501)
	header := res.Header()
	header.Set("Content-Type", "text/plain")
	fmt.Fprintf(res, "Not implemented\n")
}

func handleWelcome(res http.ResponseWriter, req *http.Request) {
	info := WelcomeInfo{"0.9"}
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

type WelcomeInfo struct {
	Version string
}
