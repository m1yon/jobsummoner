package linkedincrawler

import (
	"fmt"
	"net/url"
	"os"
	"strings"
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

	keywords := []string{"typescript", "react"}
	location := "United States"
	workTypes := []string{"2"}    // 2 = remote
	jobTypes := []string{"F"}     // F = fulltime
	salaryRanges := []string{"5"} // 5 = $120,000+
	ageOfPosting := 1 * time.Hour
	url, err := url.Parse("https://linkedin.com/jobs/search/")

	if err != nil {
		return fmt.Errorf("failed to build URL > %v", err)
	}

	q := url.Query()
	q.Set("keywords", strings.Join(keywords, " OR "))
	q.Set("location", location)
	q.Set("f_WT", strings.Join(workTypes, ","))
	q.Set("f_JT", strings.Join(jobTypes, ","))
	q.Set("f_SB2", strings.Join(salaryRanges, ","))
	q.Set("f_TPR", fmt.Sprintf("r%v", ageOfPosting.Seconds()))

	url.RawQuery = q.Encode()

	err = page.Navigate(url.String())

	fmt.Println(url)
	// err = page.Navigate("https://bot.sannysoft.com/")

	if err != nil {
		return fmt.Errorf("failed to navigate to linkedin > %v", err)
	}

	time.Sleep(1 * time.Hour)

	return nil
}
