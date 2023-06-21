package executor

import (
	"docker/types"
	"fmt"

	v8 "rogchap.com/v8go"
)

func RunFunctionWithInputs(submission types.FunctionSubmission) []types.ExecutionOutput {

	fmt.Println("function to be executed: \n", submission.Code)
	// creates a new V8 context with a new Isolate aka VM
	iso := v8.NewIsolate()
	ctx := v8.NewContext(iso)
	// executes a script on the global context
	_, err := ctx.RunScript(submission.Code, "main.js")
	if err != nil {
		// TODO add error handling
		panic(err)
	}

	fnVal, err := ctx.Global().Get(submission.FunctionName)

	if err != nil {
		// TODO add error handling
		panic(err)
	}

	function, err := fnVal.AsFunction()

	if err != nil {
		// TODO add error handling
		panic(err)
	}

	results := make([]types.ExecutionOutput, len(submission.TestCases))
	for i, testCase := range submission.TestCases {
		// TODO add error handling
		values := make([]v8.Valuer, len(testCase.InputArray))
		for i, input := range testCase.InputArray {
			value, _ := v8.NewValue(iso, input.Value)
			values[i] = value
		}

		val, err := function.Call(ctx.Global(), values...)
		if err != nil {
			// If an error occurs, create an ExecutionOutput object with the error message
			results[i] = types.ExecutionOutput{
				Value: types.Argument{},
				Error: err.Error(),
			}
			continue
		}

		// If no error occurs, create an ExecutionOutput object with the actual output value
		results[i] = types.ExecutionOutput{
			Value: v8ValueToArgument(*val),
			Error: "",
		}

	}
	return results
}

func v8ValueToArgument(value v8.Value) types.Argument {
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
