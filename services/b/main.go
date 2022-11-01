package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/tayalone/go-otel-jaeger/services/b/tracing"
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

	r.Use(otelgin.Middleware("b-service"))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong srv-b",
		})
	})

	r.GET("/todo", func(ctx *gin.Context) {
		tracer := otel.Tracer("gin-server")
		_, span := tracer.Start(ctx.Request.Context(), "getTodo")
		defer span.End()

		type todoStuct struct {
			UserID uint   `json:"userId"`
			ID     uint   `json:"id"`
			Title  string `json:"title"`
			Body   string `json:"body"`
		}

		todo := todoStuct{
			UserID: 1,
			ID:     1,
			Title:  "sunt aut facere repellat provident occaecati excepturi optio reprehenderit",
			Body:   "quia et suscipit suscipit recusandae consequuntur expedita et cum reprehenderit molestiae ut ut quas totam nostrum rerum est autem sunt rem eveniet architecto",
		}
		ctx.JSON(http.StatusOK, gin.H{
			"userId": todo.UserID,
			"id":     todo.ID,
			"title":  todo.Title,
			"body":   todo.Body,
		})
	})

	p := "3000"

	pEnv := os.Getenv("PORT")

	if pEnv != "" {
		p = pEnv
	}

	r.Run(":" + p) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
