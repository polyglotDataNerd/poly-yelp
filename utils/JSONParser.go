package utils

import (
	"fmt"
	log "github.com/polyglotDataNerd/poly-Go-utils/utils"
	"math"
	"sort"
	"strings"
	"time"
)

func JSONtoMapYelp(data map[string]interface{}) map[string]string {

	/*stopwatch start*/
	startTime := time.Now()
	transform := make(map[string]string)
	var keybuilder strings.Builder
	/*get type of reviews payload
	log.Println("var kind:", reflect.TypeOf(data["review"]).Kind())
	log.Println("var type:", reflect.TypeOf(data["review"]))*/
	rewiews := data["review"]

	/*gets the location and address to of store to append to key*/
	address := data["address"].(map[string]interface{})
	ak := fmt.Sprintf("%s:%s:", address["addressLocality"], address["streetAddress"])
	/*gets the location and address to of store to append to key*/

	/*main parent loop that gets an array of reviews of interface type*/
	for _, v := range rewiews.([]interface{}) {
		keybuilder.WriteString(ak)
		/*sorts map keys in array type*/
		var sortedkeys []string

		/*loops through the keys to append to string array for sort order*/
		for k := range v.(map[string]interface{}) {
			sortedkeys = append(sortedkeys, k)
		}
		/*needs to sort keys in the outer Map since JSON does not do any sorting to build our output Map key*/
		sort.Strings(sortedkeys)

		/*
			child loop that takes the value (v.) of the parent loop [reviews] to build Map key and
			"description" text as the value to put into Mapper struct

			KEY => addressLocality:streetAddress:author:date published:rating
		*/
		innerMap := v.(map[string]interface{})

		for _, keys := range sortedkeys {

			if keys != "description" {
				appendValue := innerMap[keys]

				switch appendValue.(type) {
				case map[string]interface{}:
					for _, v3 := range appendValue.(map[string]interface{}) {
						keybuilder.WriteString(strings.ReplaceAll(fmt.Sprintf("%v", v3), "\n", " "))
					}
				case string:
					/*converts interface to string %v*/
					keybuilder.WriteString(strings.ReplaceAll(strings.ReplaceAll(fmt.Sprintf("%v:", appendValue), ".", ""), "\n", " "))
				}
			}
		}
		/* value of key/pair, this is the yelp review */
		transform[keybuilder.String()] = strings.ReplaceAll(strings.ReplaceAll(
			strings.TrimSpace(fmt.Sprintf("%v", innerMap["description"])),
			"\n", " "), "\"", "")
		keybuilder.Reset()
	}
	endTimme := math.Round(float64(time.Since(startTime).Nanoseconds()) * 1.0e-4)
	log.Info.Println("JSON to Map processing time", endTimme, "ms")

	if len(transform) > 0 {
		return transform
	} else {
		log.Error.Println("payload is empty")
		return nil
	}
}
