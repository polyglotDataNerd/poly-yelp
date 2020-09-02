package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/polyglotDataNerd/poly-Go-utils/utils"
	"io/ioutil"
	"reflect"
	"testing"
)

type ObjMapper struct {
	payload map[string]interface{}
}

func isNilMap(inMap map[string]interface{}) (err error) {
	val := reflect.ValueOf(inMap)

	if val.Kind() == reflect.Map {
		for _, k := range val.MapKeys() {
			v := val.MapIndex(k)
			if reflect.ValueOf(v).IsNil() {
				utils.Error.Println("interface{} is null")
				err = errors.New(fmt.Sprintf("%s%s", "interface{} is null for ", k))
			}
		}
	}
	return err
}

func TestForNil_YelpInterface(t *testing.T) {
	var mapper ObjMapper
	//fullPayload, _ := ioutil.ReadFile("../fixtures/yelpPayload.json")
	noIntPayload, _ := ioutil.ReadFile("../fixtures/yelpPayloadNoInterface.json")

	payloadErr := json.Unmarshal(noIntPayload, mapper.payload)
	if payloadErr != nil {
		t.Errorf("JSON = %s; want non Nill", fmt.Sprintf("%v", mapper.payload))
	}

	testErr := isNilMap(mapper.payload); if testErr != nil {
		t.Errorf("JSON = %s; is Nil", fmt.Sprintf("%s", testErr))
	}

}
