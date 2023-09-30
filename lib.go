package zerologgrpcprovider

import (
	"encoding/json"
	"fmt"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func protobufToMap(v proto.Message) (map[string]interface{}, error) {
	vBytes, err := protojson.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}

	var vMap map[string]interface{}

	err = json.Unmarshal(vBytes, &vMap)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal message from json: %w", err)
	}

	return vMap, nil
}
