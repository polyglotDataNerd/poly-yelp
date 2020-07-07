package parser

//https://www.devdungeon.com/content/web-scraping-go
//https://schier.co/blog/2015/04/26/a-simple-web-scraper-in-go.html
import (
	"bytes"
	"fmt"
	Set "github.com/deckarep/golang-set"
	goutils "github.com/polyglotDataNerd/poly-Go-utils/utils"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"strings"
	"time"
)

func PayloadRequestString(url string) string {
	client := &http.Client{
		Timeout: 20 * time.Millisecond,
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		goutils.Info.Fatal(err)
	}
	request.Header.Set("User-Agent", "Chrome")

	response, err := client.Do(request)
	if err != nil {
		goutils.Info.Fatal(err)
	}

	defer response.Body.Close()

	outputbytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		goutils.Info.Fatal(err)
	}

	return string(outputbytes)

}

func PayloadRequest(url string) http.Response {
	/*stopwatch start*/
	startTime := time.Now()
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		goutils.Info.Fatal(err)
	}
	request.Header.Set("User-Agent", "Ubuntu")

	response, err := client.Do(request)
	if err != nil {
		goutils.Info.Fatal(err)
	}
	endTimme := math.Round(time.Since(startTime).Seconds())
	goutils.Info.Println("Yelp GET response time", endTimme, "seconds")
	return *response

}

func ReviewsJson(url string) (jsonpayload string) {
	payload := PayloadRequest(url)
	defer payload.Body.Close()

	l := html.NewTokenizer(payload.Body)
	for {
		tt := l.Next()
		switch {
		case tt == html.ErrorToken:
			return
		case tt == html.StartTagToken:
			token := l.Token()

			/*script tag, this condition finds and iterates the script tag
			to parse reviews from script tag in JSON format
			*/
			if "script" == token.Data {
				for _, a := range token.Attr {
					/*once within the node we filter for anytype of json type in the key/pair value*/
					if a.Val == "application/ld+json" {
						tt := l.Next()
						/*will get the json text in the text token of the html source.*/
						if tt == html.TextToken {
							jsonpayload := strings.TrimSpace(l.Token().Data)
							/*will only search and filter the text tokens for anything that has text review*/
							if strings.Contains(jsonpayload, "review") {
								return jsonpayload
							}
							break
						}
					}
				}
			}

		}
	}
}

func ReviewsText(url string) {
	payload := PayloadRequest(url)
	defer payload.Body.Close()

	l := html.NewTokenizer(payload.Body)
	for {
		tt := l.Next()
		switch {
		case tt == html.ErrorToken:
			return
		case tt == html.StartTagToken:
			token := l.Token()

			/*paragraph tag*/
			if "p" == token.Data {
				for _, a := range token.Attr {
					if a.Key == "itemprop" && a.Val == "description" {
						tt = l.Next()
						/*TextToken will search for the text in the node of the tag tokens*/
						if tt == html.TextToken {
							//report the page title and break out of the loop
							goutils.Info.Println(l.Token().Data)
							break
						}
					}
				}
			}

		}
	}

}

func ReviewsMeta(url string) {
	payload := PayloadRequest(url)
	defer payload.Body.Close()

	l := html.NewTokenizer(payload.Body)
	for {
		tt := l.Next()
		switch {
		case tt == html.ErrorToken:
			return
		case tt == html.StartTagToken:
			token := l.Token()

			/*meta tag*/
			if "meta" == token.Data {
				if token.Attr[0].Val == "ratingValue" {
					goutils.Info.Println(fmt.Sprintf("%s:%s", token.Attr[0].Val, token.Attr[1].Val))
				}
				if token.Attr[0].Val == "datePublished" {
					goutils.Info.Println(fmt.Sprintf("%s:%s", token.Attr[0].Val, token.Attr[1].Val))
				}
				if token.Attr[0].Val == "author" {
					goutils.Info.Println(fmt.Sprintf("%s:%s", token.Attr[0].Val, token.Attr[1].Val))
				}
				//for _, a := range token.Attr  {
				//	if a.Key == "itemprop" && a.Val == "ratingValue" {
				//		goutils.Info.Println(token.Attr[0].Val + " : " + token.Attr[1].Val)
				//	}
				//}
			}

		}
	}

}

func ReviewsHref(url string) (set Set.Set) {
	set = Set.NewSet()
	payload := PayloadRequest(url)
	defer payload.Body.Close()

	l := html.NewTokenizer(payload.Body)
	for {
		tt := l.Next()

		switch {
		case tt == html.ErrorToken:
			return
			// End of the document, we're done
		case tt == html.StartTagToken:
			t := l.Token()

			isAnchor := t.Data == "a"
			if isAnchor {
				for _, a := range t.Attr {
					if a.Key == "href" {
						if strings.Contains(a.Val, "/biz/sweetgreen") && strings.Contains(a.Val, "?osq=sweetgreen") {
							set.Add("https://www.yelp.com" + strings.ReplaceAll(a.Val, "?osq=sweetgreen", "?start="))
						}
					}
				}
			}
		}
	}
	return
}

func renderNode(n []*html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	for _, s := range n {
		html.Render(w, s)
	}
	return buf.String()

}
