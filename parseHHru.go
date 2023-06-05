package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os/user"
	"strconv"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

// Получить домашнюю директорию пользователя
func getUsrHomeDir() string {
	usr, err := user.Current()
	if err != nil {
		fmt.Println(err)
	}
	return usr.HomeDir
}

// Подготовить контекст для запуска браузера
func prepareChrome(visibleBrowser bool) (context.Context, context.CancelFunc) {
	var ctx context.Context
	var cancel context.CancelFunc
	// userDir := `C:\Users\user`
	userDir := getUsrHomeDir()
	if visibleBrowser {
		opts := append(chromedp.DefaultExecAllocatorOptions[:],
			// chromedp.DisableGPU,
			chromedp.Flag("start-maximized", true),
			chromedp.Flag("headless", false),
			chromedp.Flag("no-first-run", true),
			chromedp.Flag("no-sandbox", true),
			chromedp.Flag("disable-gpu", true),
			chromedp.Flag("enable-automation", true),
			chromedp.Flag("restore-on-startup", true),
			chromedp.UserDataDir(userDir),
		)
		ctx, cancel = chromedp.NewExecAllocator(context.Background(), opts...)
		ctx, cancel = chromedp.NewContext(ctx)
	} else {
		opts := append(chromedp.DefaultExecAllocatorOptions[:],
			// chromedp.DisableGPU,
			chromedp.Flag("headless", false),
			chromedp.Flag("no-first-run", true),
			chromedp.Flag("no-sandbox", true),
			chromedp.Flag("disable-gpu", true),
			chromedp.Flag("enable-automation", true),
			chromedp.Flag("restore-on-startup", true),
			chromedp.UserDataDir(userDir),
			chromedp.WindowSize(1280, 1024),
			// chromedp.Flag("minimal", true),
			chromedp.Flag("window-position", "-10000,-10000"),
		)
		ctx, cancel = chromedp.NewExecAllocator(context.Background(), opts...)
		ctx, cancel = chromedp.NewContext(ctx)
		_ = ctx
	}
	return ctx, cancel
}

// При первом запуске необходимо авторизоваться вручную
func firstRunChrome(ctx context.Context, cancel context.CancelFunc) {
	chromedp.Run(ctx,
		chromedp.Navigate("https://togliatti.hh.ru/account/login?backurl=%2F"),
		chromedp.WaitVisible(`/html/body/div[4]/div[1]/div/div/div[1]/div[1]/a`),
	)
	cancel()
}

// Получить список резюме
func getResumeList(ctx context.Context, cancel context.CancelFunc) (result []string) {
	if ctx == nil {
		return
	}
	if ctx.Err() != nil {
		fmt.Println(ctx.Err())
		return
	}
	defer cancel()
	ctx, cancel = context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	var nodes, children []*cdp.Node
	var resume string

	err := chromedp.Run(
		ctx,
		chromedp.Navigate("https://togliatti.hh.ru/applicant/resumes?from=header_new"),
		chromedp.Nodes(`div.bloko-column.bloko-column_xs-4.bloko-column_s-8.bloko-column_m-8.bloko-column_l-10`,
			&nodes),
	)

	if err != nil {
		fmt.Println(err)
		return
	}
	err = chromedp.Run(
		ctx,
		chromedp.Nodes("div.applicant-resumes-card-wrapper",
			&children, chromedp.ByQueryAll, chromedp.FromNode(nodes[4])),
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, n := range children {
		chromedp.Run(
			ctx,
			chromedp.Text("div>h3>a>span", &resume, chromedp.ByQueryAll, chromedp.FromNode(n)))
		result = append(result, resume)
	}
	return
}

// Монитор обновления резюме
func goUpdateMonitor(visibleBrowser bool) {
	timeout := 0
	timeUntilUpdate := 0
	defer log.Println("Exit from goUpdateMonitor")
	for {
		time.Sleep(1 * time.Second)
		lock.Lock()
		if !working {
			timeout = 0
			timeUntilUpdate = 0
			lock.Unlock()
			continue
		}
		tmp := timeoutResumeUpdate
		lock.Unlock()
		if timeout != tmp {
			lock.Lock()
			timeout = tmp
			timeUntilUpdate = 0
			lock.Unlock()
		}
		timeUntilUpdate++
		if timeUntilUpdate >= (tmp * 60) {
			var ctx context.Context
			var cancel context.CancelFunc
			if visibleBrowser {
				ctx, cancel = prepareChrome(true)
			} else {
				ctx, cancel = prepareChrome(false)
			}
			lock.Lock()
			for _, v := range resumeForUpdates {
				log.Println("Run update resume from goUpdateMonitor")
				done := make(chan bool)
				go func(done chan bool) {
					updateResume(ctx, v)
					done <- true
				}(done)
				select {
				case <-done:
					log.Println("Done update resume from goUpdateMonitor")
				case <-time.After((timeOutContextUpdateResumeInSeconds + 30) * time.Second):
					log.Println("Done update resume from goUpdateMonitor by TIMEOUT")
				}
				cancel()
			}
			lock.Unlock()
			cancel()
			timeUntilUpdate = 0
		}
	}
}

// Функция обновления резюме
func updateResume(ctx context.Context, resume string) {
	log.Println("Run update resume from updateResume")
	defer log.Println("Exit from updateResume")
	var (
		nodes, children []*cdp.Node
	)
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, timeOutContextUpdateResumeInSeconds*time.Second)
	log.Println("Getting context with timeout")
	defer cancel()
	err := chromedp.Run(
		ctx,
		chromedp.Navigate("https://togliatti.hh.ru/applicant/resumes?from=header_new"),
	)
	log.Println("Open hh.ru")
	if err != nil {
		fmt.Println(err)
		log.Println(err)
		return
	}
	err = chromedp.Run(
		ctx,
		chromedp.Nodes(`div.bloko-column.bloko-column_xs-4.bloko-column_s-8.bloko-column_m-8.bloko-column_l-10`,
			&nodes),
	)
	log.Println("Get first nodes")
	if err != nil {
		fmt.Println(err)
		log.Println(err)
		return
	}
	err = chromedp.Run(
		ctx,
		chromedp.Nodes("div.applicant-resumes-card-wrapper",
			&children, chromedp.ByQueryAll, chromedp.FromNode(nodes[4])),
	)
	log.Println("Getting next nodes array")
	if err != nil {
		fmt.Println(err)
		log.Println(err)
		return
	}
	for i, n := range children {
		resumeInt, _ := strconv.Atoi(resume)
		if (i + 1) == resumeInt {
			err = chromedp.Run(
				ctx,
				chromedp.Click("div>div.applicant-resumes-recommendations.applicant-resumes-recommendations_full-width>div>div:nth-child(1)>button",
					chromedp.ByQueryAll, chromedp.FromNode(n), chromedp.NodeVisible),
				chromedp.DoubleClick("div>div.applicant-resumes-recommendations.applicant-resumes-recommendations_full-width>div>div:nth-child(1)>button",
					chromedp.ByQueryAll, chromedp.FromNode(n), chromedp.NodeVisible),
				chromedp.DoubleClick("div>div.applicant-resumes-recommendations.applicant-resumes-recommendations_full-width>div>div:nth-child(1)>button",
					chromedp.ByQueryAll, chromedp.FromNode(n), chromedp.NodeVisible),
				chromedp.Sleep((5 * time.Second)),
			)
			log.Println("Click by element")
			if err != nil {
				log.Println(err)
				return
			}
			log.Println("Successful resume update")
		}
	}
}

func fullScreenshot(filename string, quality int64, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			// get layout metrics
			_, _, _, _, _, contentSize, err := page.GetLayoutMetrics().Do(ctx)
			if err != nil {
				return err
			}

			width, height := int64(math.Ceil(contentSize.Width)), int64(math.Ceil(contentSize.Height))

			// force viewport emulation
			err = emulation.SetDeviceMetricsOverride(width, height, 1, false).
				WithScreenOrientation(&emulation.ScreenOrientation{
					Type:  emulation.OrientationTypePortraitPrimary,
					Angle: 0,
				}).
				Do(ctx)
			if err != nil {
				return err
			}

			// capture screenshot
			*res, err = page.CaptureScreenshot().
				WithQuality(quality).
				WithClip(&page.Viewport{
					X:      contentSize.X,
					Y:      contentSize.Y,
					Width:  contentSize.Width,
					Height: contentSize.Height,
					Scale:  1,
				}).Do(ctx)
			if err != nil {
				return err
			}
			ioutil.WriteFile(filename, *res, 0644)
			return nil
		}),
	}
}
