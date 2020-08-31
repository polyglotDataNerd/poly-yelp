package utils

import (
	"fmt"
	"github.com/k3a/html2text"
	"github.com/polyglotDataNerd/poly-Go-utils/utils"
	"math"
	"sort"
	"strings"
	"time"
)

/*
@Deprecated use JSONtoMapYelpV2
*/
func JSONtoMapYelpV1(data map[string]interface{}) map[string]string {

	/*stopwatch start*/
	startTime := time.Now()
	transform := make(map[string]string)
	var keybuilder strings.Builder
	/*get type of reviews payload
	utils.Println("var kind:", reflect.TypeOf(data["review"]).Kind())
	utils.Println("var type:", reflect.TypeOf(data["review"]))*/
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
		for k, _ := range v.(map[string]interface{}) {
			//println(reflect.TypeOf(k).Kind())
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
		endTimme := math.Round(float64(time.Since(startTime).Nanoseconds()) * 1.0e-4)
		utils.Info.Println("JSON to Map processing time", endTimme, "ms")

	}
	if len(transform) > 0 {
		return transform
	} else {
		utils.Error.Println("payload is empty")
		return nil
	}

}

func JSONtoMapYelpV2(data map[string]interface{}) map[string]string {
	/*stopwatch start*/
	startTime := time.Now()
	transform := make(map[string]string)
	var keybuilder strings.Builder
	/* yelp may or may not have business address */
	var businessAddress string

	if data["bizDetailsPageProps"].(map[string]interface{})["bizContactInfoProps"].(map[string]interface{})["businessAddress"] == "no value" {
		businessAddress = "no address"
	} else {
		businessAddress = data["bizDetailsPageProps"].(map[string]interface{})["bizContactInfoProps"].(map[string]interface{})["businessAddress"].(string)
	}

	/* html2text.HTML2Text converts all encoded characters eg. &amp; into plain text */
	ak := fmt.Sprintf("%s:%s:", html2text.HTML2Text(businessAddress), html2text.HTML2Text(data["bizDetailsPageProps"].(map[string]interface{})["businessName"].(string)))

	for _, z := range
		data["bizDetailsPageProps"].(map[string]interface{})["reviewFeedQueryProps"].(map[string]interface{})["reviews"].([]interface{}) {
		keybuilder.WriteString(ak)
		//reviewDate, _ := time.Parse("2006-01-02", z.(map[string]interface{})["localizedDate"].(string))

		keyreviews := fmt.Sprintf("%s:%s:%s:%d:",
			z.(map[string]interface{})["user"].(map[string]interface{})["markupDisplayName"].(string),
			z.(map[string]interface{})["user"].(map[string]interface{})["displayLocation"].(string),
			z.(map[string]interface{})["localizedDate"].(string),
			int(z.(map[string]interface{})["rating"].(float64)))
		keybuilder.WriteString(keyreviews)

		/* value of key/pair, this is the yelp review */
		transform[keybuilder.String()] = strings.ReplaceAll(strings.ReplaceAll(
			strings.TrimSpace(fmt.Sprintf("%v", html2text.HTML2Text(z.(map[string]interface{})["comment"].(map[string]interface{})["text"].(string)))),
			"\n", " "), "\"", "")

		keybuilder.Reset()
		endTimme := math.Round(float64(time.Since(startTime).Nanoseconds()) * 1.0e-4)
		utils.Info.Println("JSON to Map processing time", endTimme, "ms")
	}
	if len(transform) > 0 {
		return transform
	} else {
		utils.Error.Println("payload is empty")
		return nil
	}
}
