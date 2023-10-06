package json

import (
	jsoniter "github.com/json-iterator/go"
)

var (
	encoder jsoniter.API
)

// type jsonExtension struct {
// 	jsoniter.DummyExtension
// }

// func (e jsonExtension) UpdateStructDescriptor(structDescriptor *jsoniter.StructDescriptor) {
// 	for _, bind := range structDescriptor.Fields {
// 		tag := bind.Field.Tag().Get("json")
// 		if strings.Contains(tag, "e-") {
// 			bind.ToNames = nil
// 		}
// 		if strings.Contains(tag, "d-") {
// 			bind.FromNames = nil
// 		}
// 	}
// }

func init() {
	encoder = jsoniter.ConfigCompatibleWithStandardLibrary

	// support omit field seperately while marshall(e-) or unmarshall(d-)
	// encoder.RegisterExtension(&jsonExtension{})
}

func MarshalToString(v any) (string, error) {
	return encoder.MarshalToString(v)
}

func MarshalToStringIgnoreError(v any) string {
	str, err := encoder.MarshalToString(v)
	if err != nil {
		return "ERROR"
	}
	return str
}

func Marshal(v any) ([]byte, error) {
	return encoder.Marshal(v)
}

func UnmarshalFromString(str string, v any) error {
	return encoder.UnmarshalFromString(str, v)
}

func Unmarshal(data []byte, v any) error {
	return encoder.Unmarshal(data, v)
}

func Convert(dest any, src any) (err error) {
	var inter []byte
	if inter, err = encoder.Marshal(src); err != nil {
		return err
	}

	if err = encoder.Unmarshal(inter, dest); err != nil {
		return err
	}

	return nil
}
