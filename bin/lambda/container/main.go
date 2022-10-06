package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	awscont "github.com/pip-services3-gox/pip-services3-aws-gox/test/container"
)

func main() {
	ctx := context.Background()
	var container *awscont.DummyLambdaFunction

	container = awscont.NewDummyLambdaFunction()

	defer container.Close(ctx, "")
	err := container.Run(ctx)
	if err != nil {
		panic(err)
	}
	lambda.Start(container.GetHandler())
}
