package main

import (
	"context"
	"log"
	"strings"

	"github.com/chromedp/chromedp"
)

func runChrome(visibleBrowser bool) {
	var ctx context.Context
	if visibleBrowser {
		opts := append(chromedp.DefaultExecAllocatorOptions[:], chromedp.DisableGPU, chromedp.Flag("headless", false))
		var cancel context.CancelFunc
		ctx, cancel = chromedp.NewExecAllocator(context.Background(), opts...)
		defer cancel()
		ctx, cancel = chromedp.NewContext(ctx)
		defer cancel()
	} else {
		var cancel context.CancelFunc
		ctx, cancel = chromedp.NewContext(context.Background())
		defer cancel()
		_ = ctx
	}
	var res string
	err := chromedp.Run(ctx, chromedp.Navigate(`https://pkg.go.dev/time`),
		chromedp.Text(`#section-documentation .Documentation-overview`, &res, chromedp.NodeVisible, chromedp.ByID))
	if err != nil {
		log.Fatal(err)
	}

	log.Println(strings.TrimSpace(res))
}
