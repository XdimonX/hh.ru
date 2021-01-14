package main

import (
	"context"
	"fmt"
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
		opts := append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.DisableGPU,
			chromedp.Flag("headless", true),
			chromedp.Flag("no-first-run", true),
			// chromedp.Flag("no-sandbox", true),
			chromedp.Flag("disable-gpu", true),
			chromedp.Flag("enable-automation", false),
			chromedp.Flag("restore-on-startup", false),
			chromedp.UserDataDir(userDir),
		)
		ctx, cancel = chromedp.NewExecAllocator(context.Background(), opts...)
		ctx, cancel = chromedp.NewContext(context.Background())
		// defer cancel()
		_ = ctx
	}
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
	defer cancel()
	ctx, cancel = context.WithTimeout(ctx, 25*time.Second)
	defer cancel()
	var nodes []*cdp.Node
	var children []*cdp.Node
	err := chromedp.Run(
		ctx,
		chromedp.Navigate("https://togliatti.hh.ru/applicant/resumes?from=header_new"),
		chromedp.Nodes(`div.bloko-column.bloko-column_xs-4.bloko-column_s-8.bloko-column_m-8.bloko-column_l-11`,
			&nodes),
	)
	chromedp.Run(
		ctx,
		chromedp.Nodes("div.bloko-gap.bloko-gap_top.bloko-gap_bottom", &children, chromedp.ByQueryAll, chromedp.FromNode(nodes[0])),
	)
	if err != nil {
		fmt.Println(err)
	}
	var childer2 []*cdp.Node
	_ = childer2
	var text string
	for _, n := range children {
		chromedp.Run(
			ctx,
			chromedp.Text("div>h3>a>span", &text, chromedp.ByQueryAll, chromedp.FromNode(n)),
		)
		fmt.Println("")
	}
	return nil
}
