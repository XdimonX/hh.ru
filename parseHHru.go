package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
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
		opts := append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.DisableGPU,
			chromedp.Flag("headless", false),
			chromedp.Flag("no-first-run", true),
			// chromedp.Flag("no-sandbox", true),
			chromedp.Flag("disable-gpu", true),
			chromedp.Flag("enable-automation", true),
			chromedp.Flag("restore-on-startup", false),
			chromedp.UserDataDir(userDir),
		)
		ctx, cancel = chromedp.NewExecAllocator(context.Background(), opts...)
		// defer cancel()
		ctx, cancel = chromedp.NewContext(ctx)
		// defer cancel()
	} else {
		opts := append(chromedp.DefaultExecAllocatorOptions[:],
			// chromedp.DisableGPU,
			chromedp.Flag("headless", false),
			chromedp.Flag("no-first-run", true),
			chromedp.Flag("no-sandbox", true),
			// chromedp.Flag("disable-gpu", true),
			chromedp.Flag("enable-automation", true),
			chromedp.Flag("restore-on-startup", true),
			chromedp.UserDataDir(userDir),
			// chromedp.WindowSize(10, 10),
			chromedp.Flag("minimal", true),
			chromedp.Flag("window-position", "-1000,-1000"),
		)
		ctx, cancel = chromedp.NewExecAllocator(context.Background(), opts...)
		ctx, cancel = chromedp.NewContext(ctx)

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
func getResumeList(ctx context.Context, cancel context.CancelFunc) (result []string) {
	defer cancel()
	ctx, cancel = context.WithTimeout(ctx, 25*time.Second)
	defer cancel()
	var nodes, children []*cdp.Node
	var resume, status string

	err := chromedp.Run(
		ctx,
		chromedp.Navigate("https://togliatti.hh.ru/applicant/resumes?from=header_new"),
		chromedp.Nodes(`div.bloko-column.bloko-column_xs-4.bloko-column_s-8.bloko-column_m-8.bloko-column_l-11`,
			&nodes),
	)

	if err != nil {
		fmt.Println(err)
	}
	err = chromedp.Run(
		ctx,
		chromedp.Nodes("div.bloko-gap.bloko-gap_top.bloko-gap_bottom",
			&children, chromedp.ByQueryAll, chromedp.FromNode(nodes[0])),
	)
	if err != nil {
		fmt.Println(err)
	}
	for _, n := range children {
		chromedp.Run(
			ctx,
			chromedp.Text("div>h3>a>span", &resume, chromedp.ByQueryAll, chromedp.FromNode(n)),
			chromedp.Text("div>div.applicant-resumes-status",
				&status, chromedp.ByQueryAll, chromedp.FromNode(n)),
		)
		if strings.ToLower(status) != "не видно никому" {
			result = append(result, resume)
		}
	}
	return
}

//Монитор обновлния резюме
func goUpdateMonitor(visibleBrowser bool) {
	timeout := 0
	timeUntilUpdate := 0
	for {
		if timeout != timeoutResumeUpdate {
			lock.Lock()
			timeout = timeoutResumeUpdate
			timeUntilUpdate = 0
			lock.Unlock()
		}
		time.Sleep(1 * time.Second)
		timeUntilUpdate++
		lock.Lock()
		tmp := timeoutResumeUpdate
		lock.Unlock()
		if timeUntilUpdate >= (tmp * 60) {
			// TODO...
			var ctx context.Context
			var cancel context.CancelFunc
			if visibleBrowser {
				ctx, cancel = prepareChrome(true)
			} else {
				ctx, cancel = prepareChrome(false)
			}
			for _, v := range resumeForUpdates {
				updateResume(ctx, v)
			}
			cancel()
			timeUntilUpdate = 0
		}
	}
}

func updateResume(ctx context.Context, resume string) {
	var nodes, children []*cdp.Node
	ctx, _ = context.WithTimeout(ctx, 25*time.Second)
	err := chromedp.Run(
		ctx,
		chromedp.Navigate("https://togliatti.hh.ru/applicant/resumes?from=header_new"),
		chromedp.Nodes(`div.bloko-column.bloko-column_xs-4.bloko-column_s-8.bloko-column_m-8.bloko-column_l-11`,
			&nodes),
	)
	if err != nil {
		fmt.Println(err)
	}
	err = chromedp.Run(
		ctx,
		chromedp.Nodes("div.bloko-gap.bloko-gap_top.bloko-gap_bottom",
			&children, chromedp.ByQueryAll, chromedp.FromNode(nodes[0])),
	)
	if err != nil {
		fmt.Println(err)
	}
	for i, n := range children {
		resumeInt, _ := strconv.Atoi(resume)
		if (i + 1) == resumeInt {
			chromedp.Run(
				ctx,
				chromedp.Click("div>div.bloko-gap.bloko-gap_top>div>div>div>div:nth-child(1)>span>button",chromedp.ByQueryAll, chromedp.FromNode(n)),
			)
		}
	}
}
