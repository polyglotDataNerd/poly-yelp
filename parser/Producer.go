package parser

import (
	"encoding/json"
	"fmt"
	"github.com/polyglotDataNerd/poly-Go-utils/scanner"
	log "github.com/polyglotDataNerd/poly-Go-utils/utils"
	jsonYelp "github.com/polyglotDataNerd/poly-yelp/utils"
	"strconv"
	"sync"
)

type ObjMapper struct {
	Yelp        map[string]interface{}
	WG          sync.WaitGroup
	Urls        chan string
	YelpChanMap chan map[string]string
}

func (receiver *ObjMapper) Producer(bucket string, key string) {
	/* sender */
	defer close(receiver.YelpChanMap)
	/* Scans all yelp URLS in object puts into a channel via go routines */
	go scanner.ProcessDir(receiver.Urls, bucket, key, "flat")
	log.Info.Println("start line scan")

	/* Yelp URLS channel coming from an s3 bucket list of YELP urls */
	for url := range receiver.Urls {
		log.Info.Println("main url", url)
		l := make(map[string]interface{})
		count := ReviewsJson(url)

		/* checks to see if string is empty */
		if len(count) > 0 {
			paginate := json.Unmarshal([]byte(count), &l)
			if paginate != nil {
				log.Warning.Println("empty count", paginate)
			}

			/*adds a +20 in the loop URL to get the last reviews if the loop count doesn't end in an even number*/
			loopcount, _ := strconv.Atoi(fmt.Sprintf("%v", l["aggregateRating"].(map[string]interface{})["reviewCount"]))
			for i := 20; i <= loopcount; i = i + 20 {

				/*Itoa turns int to primitive string and concats the pagenumer for the base url*/
				concaturl := fmt.Sprintf("%s%s%s", url, strconv.Itoa(i), "&sort_by=date_asc")

				/*for yelp reviews that only have one page*/
				if loopcount < 40 {
					concaturl = url
				}

				/* runs the parser in parallel using Waitgroup on the go routines passing it to another channel */
				receiver.WG.Add(1)
				go func(loop int) {
					log.Info.Println("url", concaturl)
					payloadString := ReviewsJson(concaturl)

					/* checks to see if string is empty */
					if len(payloadString) > 0 {
						payload := json.Unmarshal([]byte(payloadString), &receiver.Yelp)
						if payload != nil {
							log.Warning.Println("Yelp payload error", payload)
						}
						yelpReview := jsonYelp.JSONtoMapYelp(receiver.Yelp)
						defer receiver.WG.Done()
						receiver.YelpChanMap <- yelpReview
					}
				}(loopcount)
			}
			receiver.WG.Wait()
		} else {
			break
		}

	}
}
