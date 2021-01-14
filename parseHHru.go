package main

import (
	"context"
	"log"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

//Подготовить контекст для запуска браузера
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
	ctx, cancel = context.WithTimeout(ctx, 1*time.Minute)
	return ctx, cancel
}

//При первом запуске необходимо авторизоваться вручную
func firstRunChrome(ctx context.Context, cancel context.CancelFunc) {
	chromedp.Run(ctx,
		chromedp.Navigate("https://togliatti.hh.ru/account/login?backurl=%2F"),
		chromedp.WaitVisible(`/html/body/div[4]/div[1]/div/div/div[1]/div[1]/a`),
	)
	cancel()
}

//Получить список резюме
func getResumeList(ctx context.Context, cancel context.CancelFunc) []string {
	//*[@id="HH-React-Root"]/div/div/div/div[1]/div[2]
	var nodes []*cdp.Node
	chromedp.Run(
		ctx,
		chromedp.Navigate("https://togliatti.hh.ru/applicant/resumes?from=header_new"),
		chromedp.Nodes(`*[@id="HH-React-Root"]/div/div/div/div[1]/div[2]`, &nodes),
	)
	cancel()
	return nil
}
