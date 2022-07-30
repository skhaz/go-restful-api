package model

import (
	"encoding/json"
)

func Last[E any](arr []E) (E, bool) {
	if len(arr) == 0 {
		var zero E
		return zero, false
	}
	return arr[len(arr)-1], true
}

func NoError(err error) bool {
	return err == nil
}

func JSONRemarshal(bytes []byte) ([]byte, error) {
	var ifce interface{}
	err := json.Unmarshal(bytes, &ifce)
	if err != nil {
		return nil, err
	}
	return json.Marshal(ifce)
}
