package main

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
)

func runChrome(visibleBrowser bool) {
	var ctx context.Context
	if visibleBrowser {
		opts := append(chromedp.DefaultExecAllocatorOptions[:], chromedp.DisableGPU, chromedp.Flag("headless", false))
		var cancel context.CancelFunc
		ctx, cancel = chromedp.NewExecAllocator(context.Background(), opts...)
		defer cancel()
		ctx, cancel = chromedp.NewContext(ctx, chromedp.WithLogf(log.Printf))
		defer cancel()
	} else {
		var cancel context.CancelFunc
		ctx, cancel = chromedp.NewContext(context.Background())
		defer cancel()
		_ = ctx
	}
	var res string
	err := chromedp.Run(ctx, chromedp.Navigate("https://togliatti.hh.ru"), //"https://togliatti.hh.ru/account/login?backurl=%2F"
		//chromedp.WaitReady("supernova-button "),
		// chromedp.Click("supernova-button"),
		//chromedp.WaitVisible("supernova-navi-item_button"),
		// chromedp.Click("supernova-navi-item_button"),
		chromedp.WaitVisible("HH-Supernova-RegionClarification-Confirm"),
		// chromedp.SendKeys("HH-Supernova-RegionClarification-Confirm", kb.Enter),
		chromedp.Click("HH-Supernova-RegionClarification-Confirm"),
		chromedp.Click("supernova-button", chromedp.NodeVisible),
		chromedp.WaitReady("bloko-input"),
		chromedp.SendKeys("bloko-input", "xdimon777@gmail.com"),
		chromedp.Sleep((3 * time.Second)),
		chromedp.SendKeys("bloko-input_password", "72345223"),
		chromedp.Sleep((2 * time.Second)),
		// chromedp.Click("bloko-button bloko-button_primary bloko-button_stretched"),
		chromedp.SendKeys("bloko-button bloko-button_primary bloko-button_stretched", kb.Enter),
	)
	//chromedp.Text("#section-documentation .Documentation-overview", &res, chromedp.NodeVisible, chromedp.ByID))
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(1000 * time.Second)
	log.Println(strings.TrimSpace(res))
}
