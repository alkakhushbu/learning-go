package main

import "context"

type ctxKey string

const k ctxKey = "key"

var reqId = "123"

func main() {

	ctx := context.Background()
	// putting the value inside the context
	// key type should not be a primitive type
	// you should create a new custom type basis on any comparable type
	ctx = context.WithValue(ctx, k, reqId)
	getReqId(ctx)

}

func getReqId(ctx context.Context) {

	// fetching the value
	val := ctx.Value(k)
	//if val == nil {
	//
	//}

	// checking if value is of correct type
	// always do this
	reqIdString, ok := val.(string)
	if !ok {
		reqIdString = "unknown"
	}
	println(reqIdString)
}
