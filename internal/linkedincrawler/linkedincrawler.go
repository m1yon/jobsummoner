package linkedincrawler

import (
	"fmt"
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
	url := fmt.Sprintf("https://www.linkedin.com/jobs/search/?keywords=%v&location=%v&f_WT=%v&f_JT=%v&f_SB2=%v&f_TPR=r%v", strings.Join(keywords, " OR "), location, strings.Join(workTypes, ","), strings.Join(jobTypes, ","), strings.Join(salaryRanges, ","), ageOfPosting.Seconds())
	err = page.Navigate(url)

	fmt.Println(url)
	// err = page.Navigate("https://bot.sannysoft.com/")

	if err != nil {
		return fmt.Errorf("failed to navigate to linkedin > %v", err)
	}

	time.Sleep(1 * time.Hour)

	return nil
}
