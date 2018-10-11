package pfsense

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/levigross/grequests"
	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Url  string
	User string
	Pass string
}

type Pfsense struct {
	conf       Config
	s          *grequests.Session
	csrf_token string
}

var (
	RequestFailed        = fmt.Errorf("Failed to make request: %s")
	RequestInvalidMethod = fmt.Errorf("Invalid method specified [get, post] availble: %s")
	RequestCSRFExpired   = errors.New("CSRF check failed")
)

func New(conf Config) (*Pfsense, error) {
	// init our session
	ro := &grequests.RequestOptions{
		UserAgent:          "pfSense GOCLI",
		InsecureSkipVerify: true,
	}
	s := grequests.NewSession(ro)

	if conf.Url == "" {
		conf.Url = "http://192.168.1.1"
	}

	if conf.User == "" {
		conf.User = "admin"
	}

	if conf.Pass == "" {
		conf.Pass = "pfsense"
	}

	pf := &Pfsense{
		conf: conf,
		s:    s,
	}

	return pf, nil
}

func init() {
	Formatter := new(log.TextFormatter)
	Formatter.TimestampFormat = "02-01-2006 15:04:05"
	Formatter.FullTimestamp = true
	log.SetFormatter(Formatter)
	//    log.SetLevel(log.DebugLevel)
}

func (pf *Pfsense) request(method string, url string, ro *grequests.RequestOptions) (*grequests.Response, error) {
	method = strings.ToLower(method)

	var req *grequests.Response
	var err error

	if method == "get" {
		log.Debug("Making get request")
		req, err = pf.s.Get(pf.conf.Url+"/"+url, ro)
	} else if method == "post" {
		log.Debug("Making post request")
		req, err = pf.s.Post(pf.conf.Url+"/"+url, ro)
	} else {
		return nil, RequestInvalidMethod
	}

	if err != nil {
		log.Fatalf("Failed to perform request: %v", err)
		return nil, RequestFailed
	}

	if strings.Contains(req.String(), "CSRF check failed") {
		return nil, RequestCSRFExpired
	}

	log.Debugf("Made %s request to: %s", method, pf.conf.Url+"/"+url)

	pf.updateCsrf(req.String())

	return req, nil
}

func (pf *Pfsense) updateCsrf(resp string) bool {
	reCsrf := regexp.MustCompile(`var csrfMagicToken = "([^"]+)";`)
	matches := reCsrf.FindStringSubmatch(resp)

	if len(matches) < 1 {
		return false
	}

	log.Debugf("Found csrf token - %s", matches[1])

	pf.csrf_token = matches[1]
	return true
}

func eqStrip(e *goquery.Selection, i int) string {
	return strings.TrimSpace(e.Eq(i).Text())
}

func MakeTable(header []string, data [][]string) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)

	for _, v := range data {
		table.Append(v)
	}
	table.Render()

	return nil
}
