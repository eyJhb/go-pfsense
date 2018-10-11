package pfsense

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
)

type dhcp_lease struct {
	Manual      string
	Ip          string
	Mac         string
	Manufacture string
	Hostname    string
	Desc        string
	Start       string
	End         string
	Online      string
	Ltype       string
}

func (pf *Pfsense) GetDhcp() ([]dhcp_lease, error) {
	resp, err := pf.s.Get(pf.conf.Url+"/status_dhcp_leases.php", nil)

	if err != nil {
		log.Fatal("Could not get dhcp status page")
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp.String()))

	if err != nil {
		log.Fatal("Error: %v\n", err)
		return nil, err
	}

	// get csrf magic
	first := true
	//var dhcp_leases = map[int]dhcp_lease{}
	var dhcp_leases []dhcp_lease
	doc.Find("table").Eq(0).Find("tr").Each(func(i int, s *goquery.Selection) {
		if first {
			first = false
			return
		}

		tds := s.Find("td")
		mac := eqStrip(tds, 2)
		manu := ""

		if len(mac) > 17 {
			manu = strings.TrimSpace(mac[17:])
			manu = manu[1 : len(manu)-1]
			mac = mac[:17]
		}

		dhcp_leases = append(dhcp_leases, dhcp_lease{
			"N",
			eqStrip(tds, 1),
			mac,
			manu,
			eqStrip(tds, 3),
			eqStrip(tds, 4),
			eqStrip(tds, 5),
			eqStrip(tds, 6),
			eqStrip(tds, 7),
			eqStrip(tds, 8),
		})
	})

	return dhcp_leases, nil

}
