package main

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func main() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn: "http://437e75d208004309b74cac7dad719c3f@tencent.danny.games:9000/2",
		// 打印过程：如 [Sentry] 2021/11/12 22:17:12 Sending info event [9c24a1511ce74ab98177b34255db0439] to tencent.danny.games project: 2
		Debug: true,
		BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			if hint.Context != nil {
				if req, ok := hint.Context.Value(sentry.RequestContextKey).(*http.Request); ok {
					// You have access to the original Request here
					fmt.Printf("请求内容%+v\n", req)
				}
			}

			return event
		},
		AttachStacktrace: true,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
		return
	}
	// Flush buffered events before the program terminates.
	defer sentry.Flush(2 * time.Second)

	app := gin.Default()

	app.Use(sentrygin.New(sentrygin.Options{
		Repanic: true,
	}))

	app.Use(func(ctx *gin.Context) {
		if hub := sentrygin.GetHubFromContext(ctx); hub != nil {
			hub.Scope().SetTag("tag_test", "maybeYouNeedIt")
		}
		ctx.Next()
	})

	app.GET("/", func(ctx *gin.Context) {
		if hub := sentrygin.GetHubFromContext(ctx); hub != nil {
			hub.WithScope(func(scope *sentry.Scope) {
				scope.SetExtra("unwantedQuery", "someQueryDataMaybe")
				hub.CaptureMessage("User provided unwanted query string, but we recovered just fine")
			})
		}
		ctx.Status(http.StatusOK)
	})

	app.GET("/foo", func(ctx *gin.Context) {
		// sentrygin handler will catch it just fine. Also, because we attached "someRandomTag"
		// in the middleware before, it will be sent through as well
		panic("y tho")
	})

	app.Run(":3000")

}
