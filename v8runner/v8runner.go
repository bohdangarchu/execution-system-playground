package v8runner

import (
	"fmt"
	"strings"

	v8 "rogchap.com/v8go"
)

func ExecuteJavaScript(code string) (string, string, error) {
	// for every console.log the logMessages variable get updated
	// end is returned in the end
	iso := v8.NewIsolate()

	ctx := v8.NewContext(iso)

	ctx.RunScript("var logMessages = [];", "main.js")
	ctx.RunScript("console.log = function() { logMessages.push.apply(logMessages, arguments); };", "main.js")

	resultValue, err := ctx.RunScript(code, "main.js")

	if err != nil {
		return "", "", err
	}

	output, err := ctx.RunScript("logMessages", "main.js")

	if err != nil {
		return "", "", err
	}

	return resultValue.String(), output.String(), nil
}

func executeJavaScriptRedirectOutput(code string) (string, string, error) {
	// doesnt work
	// cannot set printfn as a callback for console.log
	iso := v8.NewIsolate()
	var sb strings.Builder
	printfn := v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		sb.WriteString(fmt.Sprintf("%v", info.Args()))
		return nil
	})

	global := v8.NewObjectTemplate(iso)
	setError := global.Set("console.log", printfn)

	if setError != nil {
		fmt.Println("setError ", setError)
	}

	ctx := v8.NewContext(iso, global)

	resultValue, err := ctx.RunScript(code, "main.js")

	if err != nil {
		return "", "", err
	}

	return resultValue.String(), sb.String(), nil
}
