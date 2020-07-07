package parser

import (
	"encoding/json"
	"fmt"
	"github.com/polyglotDataNerd/poly-Go-utils/scanner"
	log "github.com/polyglotDataNerd/poly-Go-utils/utils"
	"strconv"
	"sync"
	json "github.com/polyglotDataNerd/poly-yelp/utils"
)

type ObjMapper struct {
	Yelp *map[string]interface{}
	WG *sync.WaitGroup
}

func (receiver *ObjMapper) Producer(bucket string, key string, urls chan string, yelpChan chan map[string]string) {
	/* sender */
	defer close(yelpChan)
	baseMap := ObjMapper{}
	/* Scans all yelp URLS in object puts into a channel via go routines */
	scanner.ProcessDir(urls, bucket, key, "flat")
	log.Info.Println("start line scan")

	for url := range urls {
		l := make(map[string]interface{})
		json.Unmarshal([]byte(ReviewsJson(url)), &l)

		/*adds a +20 in the loop URL to get the last reviews if the loop count doesn't end in an even number*/
		loopcount, _ := strconv.Atoi(fmt.Sprintf("%v", l["aggregateRating"].(map[string]interface{})["reviewCount"]))
		for i := 20; i <= loopcount; i = i + 20 {

			/*Itoa turns int to primitive string and concats the pagenumer for the base url*/
			concaturl := fmt.Sprintf("%s%s%s", url, strconv.Itoa(i), "&sort_by=date_asc")
			/*for yelp reviews that only have one page*/
			if loopcount < 40 {
				concaturl = url
			}

			log.Info.Println("url", concaturl)
			json.Unmarshal([]byte(ReviewsJson(concaturl)), &baseMap.Yelp)
			yelpReview := json.JSONtoMapYelp(&baseMap.Yelp)

			receiver.WG.Add(1)
			go func(input map[string]string) {
				defer receiver.WG.Done()
				yelpChan <- yelpReview
			}(yelpReview)
		}
		receiver.WG.Wait()

	}
}
