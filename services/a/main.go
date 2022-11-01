package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/tayalone/go-otel-jaeger/services/a/tracing"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func main() {
	tp, tpErr := tracing.JaegerProvider()
	if tpErr != nil {
		log.Fatal(tpErr)
	}

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	r := gin.Default()

	r.Use(otelgin.Middleware("a-service"))

	r.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/get-todo", func(ctx *gin.Context) {
		response, err := http.Get("https://jsonplaceholder.typicode.com/posts/1")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "call api fail",
			})
		}

		responseData, err := ioutil.ReadAll(response.Body)

		type respStuct struct {
			UserID uint   `json:"userId"`
			ID     uint   `json:"id"`
			Title  string `json:"title"`
			Body   string `json:"body"`
		}

		var respObj respStuct
		json.Unmarshal(responseData, &respObj)

		ctx.JSON(http.StatusOK, gin.H{
			"post": respObj,
		})
	})

	p := "3000"

	pEnv := os.Getenv("PORT")

	if pEnv != "" {
		p = pEnv
	}

	r.Run(":" + p) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
