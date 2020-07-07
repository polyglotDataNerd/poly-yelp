package main

import (
	"fmt"
	goaws "github.com/polyglotDataNerd/poly-Go-utils/aws"
	goutils "github.com/polyglotDataNerd/poly-Go-utils/utils"
	"github.com/polyglotDataNerd/poly-yelp/parser"
	uuid "github.com/satori/go.uuid"
	"math"
	"runtime"
	"strings"
	"sync"
	"time"
)

/*
 This is the main entry point of the yelp goscraper that uses the goroutine/channel design pattern to produce a review dataset by
 iterating through all review URL's concurrently parsing the HTML to get a map of reviews back.
*/

func main() {
	/*go routine that runs concurrently*/
	runtime.GOMAXPROCS(runtime.NumCPU())

	/*stopwatch start*/
	startTime := time.Now()
	var stringBuilder strings.Builder

	//var stringBuilder strings.Builder
	var WG sync.WaitGroup
	s3Bucket := "poly-testing"
	key := "yelp"
	sourceUrls := "yelp/urls"
	yelpUrls := make(chan string)
	yelpChan := make(chan map[string]string)
	baseMap := make(map[string]interface{})

	producer := parser.ObjMapper{
		Yelp:        baseMap,
		WG:          WG,
		Urls:        yelpUrls,
		YelpChanMap: yelpChan}

	go producer.Producer(s3Bucket, sourceUrls)

	for yMap := range yelpChan {
		for k, v := range yMap {
			keyArray := strings.Split(k, ":")
			stringFormat := fmt.Sprintf("\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\"%s",
				strings.ReplaceAll(keyArray[0], "\n", " "),
				strings.ReplaceAll(keyArray[1], "\n", " "),
				strings.ReplaceAll(keyArray[2], "\n", " "),
				strings.ReplaceAll(keyArray[3], "\n", " "),
				strings.ReplaceAll(keyArray[4], "\n", " "),
				v, "\n")
			stringBuilder.WriteString(stringFormat)
		}
		/* "2006-01-02" is the standard time format for go lang YYYY-MM-DD*/
		objectKey := key + "/reviews/" + time.Now().Format("2006-01-02") + "/" + uuid.NewV4().String() + "/" + time.Now().Format("2006-01-02") + "-" + uuid.NewV4().String() + ".gz"
		/*writes payload to s3*/
		goaws.S3Obj{Bucket: s3Bucket, Key: objectKey}.S3WriteGzip(stringBuilder.String(), goaws.SessionGenerator())
	}
	/*stop watch end*/
	endTimme := math.Round(time.Since(startTime).Seconds())
	goutils.Info.Println("Yelp Scraper Finish Time", endTimme, "seconds")

}
