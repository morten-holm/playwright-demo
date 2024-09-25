package main

import (
	"fmt"
	"net/http"

	"github.com/playwright-community/playwright-go"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {
	pw, err := playwright.Run()
	if err != nil {
		panic(err)
	}
	defer func() {
		err = pw.Stop()
		if err != nil {
			panic(err)
		}
	}()

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Args: []string{
			"--allow-running-insecure-content",
		},
		Headless: playwright.Bool(true),
	})
	if err != nil {
		panic(err)
	}
	defer func() {
		err = browser.Close()
		if err != nil {
			panic(err)
		}
	}()

	baseUrl := "https://playwright.dev/"
	permissions := []string{"clipboard-read", "clipboard-write", "storage-access"}

	page, err := browser.NewPage(playwright.BrowserNewPageOptions{
		BaseURL:           playwright.String(baseUrl),
		BypassCSP:         playwright.Bool(true),
		JavaScriptEnabled: playwright.Bool(true),
		Permissions:       permissions,
		Viewport: &playwright.Size{
			Width:  1240,
			Height: 800,
		},
	})
	if err != nil {
		panic(err)
	}
	defer func() {
		err = page.Close()
		if err != nil {
			panic(err)
		}
	}()

	webErrors := make([]string, 0)

	page.OnPageError(func(e error) {
		webErrors = append(webErrors, e.Error())
	})

	page.OnConsole(func(e playwright.ConsoleMessage) {
		webErrors = append(webErrors, e.Text())
	})

	_, err = page.Goto(baseUrl, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})
	if err != nil {
		panic(err)
	}

	for _, webError := range webErrors {
		fmt.Printf("%s\n", webError)
	}

	_, err = fmt.Fprintf(w, "pong")
	if err != nil {
		panic(err)
	}
}

func main() {
	err := playwright.Install(&playwright.RunOptions{
		Browsers: []string{"chromium"},
	})

	if err != nil {
		panic(err)
	}

	http.HandleFunc("/ping", pingHandler)

	fmt.Println("Starting server at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
