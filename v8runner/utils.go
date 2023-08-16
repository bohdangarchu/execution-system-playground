package v8runner

import (
	"app/types"
	"fmt"

	v8 "rogchap.com/v8go"
)

func v8ValueToArgument(value v8.Value) types.Argument {
	// not used because we use JSON types
	switch {
	case value.IsString():
		return types.Argument{
			Value: value.String(),
			Type:  "string",
		}
	case value.IsNumber():
		// Assuming it's a float64 value for simplicity, but you can modify it accordingly
		return types.Argument{
			Value: value.Number(),
			Type:  "number",
		}
	case value.IsBoolean():
		return types.Argument{
			Value: value.Boolean(),
			Type:  "boolean",
		}
	case value.IsArray():
		fmt.Printf("length of return value %v", value.Object().InternalFieldCount())
		// Assuming it's an array of V8 values
		var elements []types.Argument
		for i := 0; i < int(value.Object().InternalFieldCount()); i++ {
			elem, err := value.Object().GetIdx(uint32(i))
			if err != nil {
				panic(err)
			}
			elements[i] = types.Argument{
				Value: v8ValueToArgument(*elem),
				Type:  "array",
			}
		}
		return types.Argument{
			Value: elements,
			Type:  "array",
		}
	// case value.IsObject():
	// 	// Assuming it's a plain object
	// 	object := make(map[string]interface{})
	// 	value.ForEach(func(key string, prop v8go.Value) {
	// 		object[key] = v8ValueToArgument(prop).Value
	// 	})
	// 	return types.Argument{
	// 		Value: object,
	// 		Type:  "object",
	// 	}
	default:
		// Returning nil for unsupported types or fallback case
		return types.Argument{
			Value: nil,
			Type:  "unknown",
		}
	}
}
