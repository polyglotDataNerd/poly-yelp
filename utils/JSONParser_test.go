package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/polyglotDataNerd/poly-Go-utils/utils"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

/*
Debug Test using Delve
	https://golangbot.com/debugging-go-delve/
		dlv test
*/
type ObjMapper struct {
	payload map[string]interface{}
}

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	os.Exit(m.Run())
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
	//fullPayload, _ := ioutil.ReadFile("./fixtures/yelpPayload.json")
	noIntPayload, _ := ioutil.ReadFile("./fixtures/yelpPayloadNoInterface.json")

	payloadErr := json.Unmarshal(noIntPayload, &mapper.payload)
	if payloadErr != nil {
		t.Errorf("JSON = %s; want non Nill", fmt.Sprintf("%v", mapper.payload))
	}

	testErr := isNilMap(mapper.payload)
	if testErr != nil {
		t.Errorf("JSON = %s; is Nil", fmt.Sprintf("%s", testErr))
	}

}
