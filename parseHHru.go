package main

import (
	"context"
	"log"

	"github.com/chromedp/chromedp"
	// "github.com/chromedp/chromedp/kb"
)

func prepareChrome(visibleBrowser bool) (context.Context, context.CancelFunc) {
	var ctx context.Context
	var cancel context.CancelFunc
	userDir := `C:\Users\user`
	if visibleBrowser {
		opts := append(chromedp.DefaultExecAllocatorOptions[:], chromedp.DisableGPU, chromedp.Flag("headless", false),
			chromedp.Flag("no-first-run", true),
			// chromedp.Flag("no-sandbox", true),
			chromedp.Flag("disable-gpu", true),
			chromedp.Flag("enable-automation", false),
			chromedp.Flag("restore-on-startup", false),
			chromedp.UserDataDir(userDir),
		)
		ctx, cancel = chromedp.NewExecAllocator(context.Background(), opts...)
		// defer cancel()
		ctx, cancel = chromedp.NewContext(ctx, chromedp.WithLogf(log.Printf))
		// defer cancel()
	} else {
		ctx, cancel = chromedp.NewContext(context.Background())
		// defer cancel()
		_ = ctx
	}
	// err := chromedp.Run(ctx, chromedp.Navigate("https://togliatti.hh.ru"),
	// 	chromedp.Click(`/html/body/div[4]/div[1]/div/div/div[1]/div[1]/a`),
	// )
	// if err != nil {
	// 	log.Fatal(err)
	// }
	return ctx, cancel
}

//При превом запуске необходимо авторизоваться вручную
func firstRunChrome(ctx context.Context, cancel context.CancelFunc) {
	chromedp.Run(ctx,
		chromedp.Navigate("https://togliatti.hh.ru/account/login?backurl=%2F"),
		chromedp.WaitVisible(`/html/body/div[4]/div[1]/div/div/div[1]/div[1]/a`),
	)
	cancel()
}

func getResumeList() {

}
