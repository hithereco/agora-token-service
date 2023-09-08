package main

import (
	"context"
	"os"

	"github.com/AgoraIO-Community/agora-token-service/service"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

var ginLambda *ginadapter.GinLambda

func main() {
	s := service.NewService()

	env := os.Getenv("GIN_MODE")
	if env == "release" {
		ginEngine, ok := s.Server.Handler.(*gin.Engine)
		if ok {
			ginLambda = ginadapter.New(ginEngine)
			lambda.Start(Handler)
		}
	} else {
		// Stop is called on another thread, but waits for an interrupt
		go s.Stop()
		s.Start()
	}
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, request)
}
