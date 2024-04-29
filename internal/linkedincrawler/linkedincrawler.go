package linkedincrawler

import (
	"fmt"
	"os"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/launcher/flags"
	"github.com/go-rod/stealth"
)

func Crawl() error {

	PROXY_HOSTNAME := os.Getenv("PROXY_HOSTNAME")
	PROXY_USERNAME := os.Getenv("PROXY_USERNAME")
	PROXY_PASSWORD := os.Getenv("PROXY_PASSWORD")

	// proxy setup
	l := launcher.New()
	l = l.Set(flags.ProxyServer, PROXY_HOSTNAME)
	controlURL, _ := l.Launch()
	browser := rod.New()
	err := browser.ControlURL(controlURL).Connect()

	if err != nil {
		return fmt.Errorf("browser connection failed > %v", err)
	}

	go browser.MustHandleAuth(PROXY_USERNAME, PROXY_PASSWORD)()

	page, err := stealth.Page(browser)

	if err != nil {
		return fmt.Errorf("failed to create stealth page > %v", err)
	}

	err = page.Navigate("https://linkedin.com")
	// err = page.Navigate("https://bot.sannysoft.com/")

	if err != nil {
		return fmt.Errorf("failed to navigate to linkedin > %v", err)
	}

	time.Sleep(1 * time.Hour)

	return nil
}
