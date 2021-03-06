package parser

import (
	"encoding/json"
	"fmt"
	"github.com/polyglotDataNerd/poly-Go-utils/scanner"
	log "github.com/polyglotDataNerd/poly-Go-utils/utils"
	jsonYelp "github.com/polyglotDataNerd/poly-yelp/utils"
	"strconv"
	"sync"
	"time"
)

type ObjMapper struct {
	Yelp        map[string]interface{}
	WG          sync.WaitGroup
	Urls        chan string
	YelpChanMap chan map[string]string
}

func (receiver *ObjMapper) Producer(bucket string, urls string, loadType string) {
	/* sender */
	defer close(receiver.YelpChanMap)

	if loadType == "files" {
		/*
			args passes an s3 object that has many urls
			Scans all yelp URLS in object puts into a channel via go routines
		*/
		go scanner.ProcessDir(receiver.Urls, bucket, urls, "flat")
		log.Info.Println("start line scan")
	} else if loadType == "url" {
		/* args passes a single url */
		go func(url string) {
			defer close(receiver.Urls)
			receiver.Urls <- urls
		}(urls)
	}

	/* Yelp URLS channel coming from an s3 bucket list of YELP urls */
	for url := range receiver.Urls {
		log.Info.Println("main url", url)
		l := make(map[string]interface{})
		count := ReviewsJsonV1(url)

		/* checks to see if string is empty */
		if len(count) > 0 {
			paginate := json.Unmarshal([]byte(count), &l)
			if paginate != nil {
				log.Error.Println("empty count", paginate)
			}

			/*adds a +20 in the loop URL to get the last reviews if the loop count doesn't end in an even number*/
			loopcount, _ := strconv.Atoi(fmt.Sprintf("%v", l["aggregateRating"].(map[string]interface{})["reviewCount"]))
			log.Info.Println("Number of reviews", loopcount)
			for i := 20; i <= loopcount; i = i + 20 {

				/*Itoa turns int to primitive string and concats the pagenumer for the base url*/
				concaturl := fmt.Sprintf("%s%s%s", url, strconv.Itoa(i), "&sort_by=date_asc")

				/*for yelp reviews that only have one page*/
				if loopcount < 40 {
					concaturl = url
				}

				/* runs the parser in parallel using Waitgroup on the go routines passing it to another channel */
				receiver.WG.Add(1)
				/* controls throttling, API calls are to fast when ran in parallel and it's causing HTTP call to fail,
				works fine when number of parallel calls are below 100 */
				time.Sleep(1 * time.Millisecond)

				go func(url string) {
					defer receiver.WG.Done()
					log.Info.Println("url", concaturl)
					payloadString := ReviewsJsonV2(concaturl)
					/* checks to see if string is empty */
					if len(payloadString) > 0 {
						payloaderr := json.Unmarshal([]byte(payloadString), &receiver.Yelp)
						if payloaderr != nil {
							log.Error.Println("Yelp payload error", payloaderr)
						}
						yelpReview := jsonYelp.JSONtoMapYelpV2(receiver.Yelp)
						receiver.YelpChanMap <- yelpReview
					}
				}(concaturl)
			}
			receiver.WG.Wait()
		} else {
			log.Warning.Println("Empty results for url", url)
			continue
		}

	}
}
