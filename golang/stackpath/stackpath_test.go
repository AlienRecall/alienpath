package stackpath

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"regexp"
	"testing"

	"golang.org/x/net/publicsuffix"
)

func TestStackapth(t *testing.T) {
	jar, _ := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	client := http.Client{Jar: jar}
	req, _ := http.NewRequest("GET", "https://www.basket4ballers.com/fr/authentification?back=my-account", nil)
	req.Header.Set("authority", "www.basket4ballers.com")
	req.Header.Set("upgrade-insecure-requests", "1")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.96 Safari/537.36")
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("sec-fetch-mode", "navigate")
	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Set("referer", "https://www.basket4ballers.com/fr/authentification?back=my-account")
	req.Header.Set("accept-language", "en-GB,en;q=0.9")
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	doc, _ := ioutil.ReadAll(res.Body)
	re := regexp.MustCompile("<title>(.*?)<")
	matches := re.FindStringSubmatch(string(doc))
	if len(matches) < 2 {
		panic("Cannot find page title")
	}
	if matches[1] == "StackPath" {
		challenge := &NewChallenge{}
		solvedBody, newJar, err := challenge.Solve(string(doc), "www.basket4ballers.com", client, req)
		if err != nil {
			panic(fmt.Sprintf("Error solving challenge: %s", err.Error()))
		}
		client.Jar = newJar
		t.Logf("%s", solvedBody)
	}
}
