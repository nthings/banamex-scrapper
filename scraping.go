package main

import (
	"log"
	"os"
	"time"

	"github.com/mxschmitt/playwright-go"
)

func scrape() {
	username, username_present := os.LookupEnv("CLIENT_NUMBER")
	password, password_present := os.LookupEnv("CLIENT_PASSWORD")
	if !(username_present && password_present) {
		log.Fatalln("One of the required env vars is not set!")
	}
	// Initializing
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not launch playwright: %v", err)
	}
	browser, err := pw.Chromium.Launch()
	if err != nil {
		log.Fatalf("could not launch Chromium: %v", err)
	}
	page, err := browser.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}
	// Go to banamex
	if _, err = page.Goto("http://bancanet.banamex.com/MXGCB/JPS/portal/LocaleSwitch.do?locale=es_MX/", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
		Timeout:   playwright.Float(60000),
	}); err != nil {
		log.Fatalf("could not goto: %v", err)
	}
	// Close popup
	if err = page.Click("#splash-207555-close-button"); err != nil {
		log.Println("Popup didn't opened")
	}
	// Fill login form
	if err = page.Click("#textCliente"); err != nil {
		log.Fatalf("could not click on #textCliente: %v", err)
	}
	if err = page.Type("#textCliente", username); err != nil {
		log.Fatalf("could not type on #textCliente: %v", err)
	}
	if err = page.Press("#textCliente", "Enter"); err != nil {
		log.Fatalf("could not press enter on #textCliente: %v", err)
	}
	if err = page.Click("#textFirma"); err != nil {
		log.Fatalf("could not click on #textFirma: %v", err)
	}
	if err = page.Type("#textFirma", password); err != nil {
		log.Fatalf("could not type on #textFirma: %v", err)
	}
	if err = page.Press("#textFirma", "Enter"); err != nil {
		log.Fatalf("could not press enter on #textFirma: %v", err)
	}

	time.Sleep(5 * time.Second)
	if err = page.Click(".closeSvgImg"); err != nil {
		log.Println("Popup didn't opened")
	}

	entries, err := page.QuerySelectorAll(".account-list-content")
	if err != nil {
		log.Fatalf("could not get entries: %v", err)
	}

	if len(entries) > 0 {
		log.Fatalln("Couldn't find accounts")
	}

	for i, entry := range entries {
		link, err := entry.QuerySelector(".account-mask-link")
		if err != nil {
			log.Fatalf("could not get title element: %v", err)
		}
		// title, err := titleElement.TextContent()
		// if err != nil {
		// 	log.Fatalf("could not get text content: %v", err)
		// }
		log.Printf("%d: %s\n", i+1, link)
	}

	// Screenshot to see where we at and close
	if _, err = page.Screenshot(playwright.PageScreenshotOptions{
		Path: playwright.String("foo.png"),
	}); err != nil {
		log.Fatalf("could not create screenshot: %v", err)
	}
	if err = browser.Close(); err != nil {
		log.Fatalf("could not close browser: %v", err)
	}
	if err = pw.Stop(); err != nil {
		log.Fatalf("could not stop Playwright: %v", err)
	}
}
