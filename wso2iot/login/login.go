package login

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	//	"io/ioutil"
	"crypto/tls"
	"net/http/cookiejar"

	"golang.org/x/net/html"
)

var loginUrl string

func parse(n *html.Node, postValues url.Values) {
	if n.Type == html.ElementNode {
		if n.Data == "form" {
			for _, f := range n.Attr {
				if f.Key == "action" {
					loginUrl = f.Val
					break
				}
			}
		}
		if n.Data == "input" {
			var key, val string
			for _, f := range n.Attr {
				if f.Key == "name" {
					key = f.Val
				}
				if f.Key == "value" {
					val = f.Val
				}
			}
			postValues.Add(key, val)
			//fmt.Println("+++++++++++++++++++++++++++ key, val = ", key, val)
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		parse(c, postValues)
	}
}

func httpRequest(method string, client *http.Client, postValues *url.Values) *url.Values {
	result := url.Values{}
	var resp *http.Response
	var err error
	if method == "POST" {
		//resp, err = client.PostForm(loginUrl, *postValues)
		data := postValues.Encode()
		req, _ := http.NewRequest("POST", loginUrl, strings.NewReader(data))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		resp, err = client.Do(req)
		postValues = &url.Values{}
	} else {
		//resp, err = client.Get(loginUrl)
		req, _ := http.NewRequest("GET", loginUrl, nil)
		resp, err = client.Do(req)
	}
	if err != nil {
		// handle error
		log.Fatal(err)
	}
	status := resp.StatusCode
	fmt.Println(status)
	/*
			fmt.Println(resp.Header)
			fmt.Println("##################    cookies")
			u, _ := url.Parse(loginUrl)
			for _, cookie := range client.Jar.Cookies(u) {
		        fmt.Printf("  %s: %s\n", cookie.Name, cookie.Value)
			}
	*/
	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Panic(err)
	}
	parse(doc, result)
	/*
		fmt.Println("%%%%%%%%%%%%%%%%%     begin")
		for k, v := range result {
			fmt.Println(k, v)
		}
		fmt.Println("%%%%%%%%%%%%%%%%%     end")
	*/
	return &result
}

func Login(host string) *http.Client {
	loginUrl = host + "/devicemgt"
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{
		Transport: tr,
		Jar:       jar,
	}

	// step 1 ----------------------------------------------------
	fmt.Println("step 1 ------- ", loginUrl)
	postValues := httpRequest("GET", client, nil)

	// step 2 ----------------------------------------------------
	fmt.Println("step 2 ------- ", loginUrl)
	postValues = httpRequest("POST", client, postValues)

	// step 3 ----------------------------------------------------
	//loginUrl = "https://localhost:9443" + loginUrl
	loginUrl = host + loginUrl
	fmt.Println("step 3 ------- ", loginUrl)
	postValues.Set("username", "admin")
	postValues.Set("password", "admin")
	postValues = httpRequest("POST", client, postValues)

	// step 4 ----------------------------------------------------
	fmt.Println("step 4 ------- ", loginUrl)
	postValues = httpRequest("POST", client, postValues)

	return client
}
