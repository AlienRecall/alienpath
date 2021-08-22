package stackpath

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"rogchap.com/v8go"
)

//NewChallenge main solver struct
type NewChallenge struct {
	client              http.Client
	challengeBody       string
	domain              string
	initReq             *http.Request
	secondChallengeBody string
	challengeForm       map[string]interface{}
	data                map[string]interface{}
}

var (
	errScrapeVal = errors.New("[StackPath Solver] Error scraping variants")
	errExtJS     = errors.New("[StackPath Solver] JS error")
)

func randInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(max-min+1) + min
}

func randString(n int) string {
	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func (nc *NewChallenge) getSBTSCK() (string, error) {
	re := regexp.MustCompile(`"sbtsck=(.*?);`)
	matches := re.FindStringSubmatch(nc.challengeBody)
	if len(matches) < 2 {
		return "", errScrapeVal
	}
	return matches[1], nil
}

func (nc *NewChallenge) getGPRID() (string, error) {
	re := regexp.MustCompile(`genPid\(\) {return (.*?) ;`)
	matches := re.FindStringSubmatch(nc.challengeBody)
	if len(matches) < 2 {
		return "", errScrapeVal
	}
	ctx, _ := v8go.NewContext()
	gprid, err := ctx.RunScript(string(matches[1]), "")
	if err != nil {
		return "", errExtJS
	}
	ret := gprid.String()
	ctx.Close()
	return ret, nil
}

func (nc *NewChallenge) getSBBGS() (string, error) {
	re := regexp.MustCompile(`sbbsv\("D-(.*?)"`)
	matches := re.FindStringSubmatch(nc.challengeBody)
	if len(matches) < 2 {
		return "", errScrapeVal
	}
	return matches[1], nil
}

func (nc *NewChallenge) getDDL() (string, error) {
	re := regexp.MustCompile(`'&ddl='\+(.*?)\+`)
	matches := re.FindStringSubmatch(nc.challengeBody)
	if len(matches) < 2 {
		return "", errScrapeVal
	}
	ctx, _ := v8go.NewContext()
	ddl, err := ctx.RunScript(strings.ReplaceAll(string(matches[1]), "dfx", "new Date()"), "")
	if err != nil {
		return "", errExtJS
	}
	ret := ddl.String()
	ctx.Close()
	return ret, nil
}

func (nc *NewChallenge) getADOTR() (string, error) {
	re := regexp.MustCompile(`parent\.otr = (.*?);`)
	matches := re.FindStringSubmatch(nc.secondChallengeBody)
	if len(matches) < 2 {
		return "", errScrapeVal
	}
	ctx, _ := v8go.NewContext()
	adotr, err := ctx.RunScript(string(matches[1])+".join('')", "")
	if err != nil {
		return "", errExtJS
	}
	ret := adotr.String()
	ctx.Close()
	return ret, nil
}

func (nc *NewChallenge) getTRSTR() (string, error) {
	re := regexp.MustCompile(`sbbdep\("(.*?)"`)
	matches := re.FindStringSubmatch(nc.secondChallengeBody)
	if len(matches) < 2 {
		return "", errScrapeVal
	}
	return matches[1], nil
}

func (nc *NewChallenge) getSecondChallenge() (string, error) {
	req, err := http.NewRequest("GET", nc.data["challengeURL"].(string), nil)
	if err != nil {
		return "", errors.New("[StackPath Solver] Request - Setup Error")
	}
	req.Header.Set("authority", nc.domain)
	req.Header.Set("upgrade-insecure-requests", "1")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.96 Safari/537.36")
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("sec-fetch-mode", "navigate")
	req.Header.Set("sec-fetch-dest", "iframe")
	req.Header.Set("referer", fmt.Sprintf("https://%s/", nc.domain))
	req.Header.Set("accept-language", "en-GB,en;q=0.9")
	res, err := nc.client.Do(req)
	if err != nil {
		return "", errors.New("[StackPath Solver] [GSC] Request - Error")
	}
	defer res.Body.Close()
	doc, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", errors.New("[StackPath Solver] [GSC] Error Parsing Body")
	}
	if res.StatusCode < 400 {
		return string(doc), nil
	}
	return "", errors.New("[StackPath Solver] [GSC] Bad Response")
}

func (nc *NewChallenge) submitChallenge() (string, error) {
	postData := url.Values{}
	postData.Set("cdmsg", nc.challengeForm["cdmsg"].(string))
	postData.Set("femsg", fmt.Sprintf("%d", nc.challengeForm["femsg"].(int)))
	postData.Set("bhvmsg", nc.challengeForm["bhvmsg"].(string))
	postData.Set("futgs", nc.challengeForm["futgs"].(string))
	postData.Set("jsdk", nc.challengeForm["jsdk"].(string))
	postData.Set("glv", nc.challengeForm["glv"].(string))
	postData.Set("lext", nc.challengeForm["lext"].(string))
	postData.Set("sdrv", fmt.Sprintf("%d", nc.challengeForm["sdrv"].(int)))
	req, err := http.NewRequest("POST", nc.data["challengeURL"].(string), bytes.NewBufferString(postData.Encode()))
	if err != nil {
		return "", errors.New("[StackPath Solver] Request - Setup Error")
	}
	req.Header.Set("authority", nc.domain)
	req.Header.Set("cache-control", "max-age=0")
	req.Header.Set("upgrade-insecure-requests", "1")
	req.Header.Set("origin", fmt.Sprintf("https://%s", nc.domain))
	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.96 Safari/537.36")
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("sec-fetch-mode", "navigate")
	req.Header.Set("sec-fetch-dest", "iframe")
	req.Header.Set("referer", nc.data["challengeURL"].(string))
	req.Header.Set("accept-language", "en-GB,en;q=0.9")
	res, err := nc.client.Do(req)
	if err != nil {
		return "", errors.New("[StackPath Solver] [GSC] Request - Error")
	}
	defer res.Body.Close()
	doc, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", errors.New("[StackPath Solver] [GSC] Error Parsing Body")
	}
	if res.StatusCode < 400 {
		return string(doc), nil
	}
	return "", errors.New("[StackPath Solver] [GSC] Bad Response")
}

func xrv(ketStr string, sourceStr string) string {
	var (
		keyLength = len(ketStr)
		rPos      int
		a         int
		b         int
		c         int
		d         string
		targetStr string
	)
	for i := 0; i < len(sourceStr); i++ {
		rPos = i % keyLength
		a = int([]rune(sourceStr)[i])
		b = int([]rune(ketStr)[rPos])
		c = a ^ b
		d = fmt.Sprintf("%da", c)
		targetStr = targetStr + d
	}
	return targetStr
}

//Solve init the solving process
func (nc *NewChallenge) Solve(challengeBody string, domain string, client http.Client, initReq *http.Request) (io.ReadCloser, http.CookieJar, error) {
	nc.client = client
	nc.challengeBody = challengeBody
	nc.domain = domain
	nc.initReq = initReq
	nc.challengeForm = map[string]interface{}{
		"cdmsg":  "",
		"femsg":  1,
		"bhvmsg": "",
		"futgs":  "",
		"jsdk":   "",
		"glv":    "",
		"lext":   "",
		"sdrv":   0,
	}
	nc.data = make(map[string]interface{})

	var err error
	nc.data["sbtsck"], err = nc.getSBTSCK()
	if err != nil {
		return nil, nil, err
	}
	nc.data["gprid"], err = nc.getGPRID()
	if err != nil {
		return nil, nil, err
	}
	nc.data["sbbgs"], err = nc.getSBBGS()
	if err != nil {
		return nil, nil, err
	}
	nc.data["ddl"], err = nc.getDDL()
	if err != nil {
		return nil, nil, err
	}
	u, _ := url.Parse(fmt.Sprintf("https://%s", nc.domain))
	var cs []*http.Cookie
	nc.client.Jar.SetCookies(u, append(cs, &http.Cookie{Name: "UTGv2", MaxAge: -1}))
	nc.client.Jar.SetCookies(u, append(cs, &http.Cookie{Name: "UTGv2", Value: nc.data["sbbgs"].(string)}))
	nc.client.Jar.SetCookies(u, append(cs, &http.Cookie{Name: "PRLST", Value: nc.data["gprid"].(string)}))
	nc.client.Jar.SetCookies(u, append(cs, &http.Cookie{Name: "sbtsck", Value: nc.data["sbtsck"].(string)}))
	nc.data["challengeURL"] = fmt.Sprintf("https://%s/sbbi/?sbbpg=sbbShell&gprid=%s&sbbgs=%s&ddl=%s", nc.domain, nc.data["gprid"], nc.data["sbbgs"], nc.data["ddl"])
	nc.secondChallengeBody, err = nc.getSecondChallenge()
	if err != nil {
		return nil, nil, err
	}
	nc.data["adotr"], err = nc.getADOTR()
	if err != nil {
		return nil, nil, err
	}
	nc.client.Jar.SetCookies(u, append(cs, &http.Cookie{Name: "adOtr", Value: nc.data["adotr"].(string)}))
	trstr, err := nc.getTRSTR()
	if err != nil {
		return nil, nil, err
	}
	nc.challengeForm["jsdk"] = trstr
	nc.challengeForm["glv"] = xrv(strings.ToUpper(trstr), fmt.Sprintf("%s.local", uuid.New().String()))
	nc.challengeForm["lext"] = xrv(strings.ToUpper(trstr), "[0,0]")
	nc.challengeForm["bhvmsg"] = xrv(strings.ToUpper(trstr), fmt.Sprintf("%s-%s", randString(10), randString(5)))
	nc.challengeForm["cdmsg"] = xrv(strings.ToUpper(trstr), fmt.Sprintf("%s-41-%s-%s-%s-noieo-90.%d", randString(11), randString(9), randString(11), randString(11), randInt(2000000000000000, 9999999999999999)))
	_, err = nc.submitChallenge()
	if err != nil {
		return nil, nil, err
	}
	solvedRes, err := nc.client.Do(nc.initReq)
	if err != nil {
		return nil, nil, errors.New("[StackPath Solver] [SR] Request - Error")
	}
	defer solvedRes.Body.Close()
	return solvedRes.Body, nc.client.Jar, nil
}
