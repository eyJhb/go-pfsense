package pfsense

import (
    "strings"
    log "github.com/sirupsen/logrus"
    "github.com/PuerkitoBio/goquery"
)

type dhcp_lease struct {
    Manual bool
    Ip string
    Mac string
    Hostname string
    Desc string
    Start string
    End string
    Online string
    Ltype string
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
        dhcp_leases = append(dhcp_leases, dhcp_lease{
            false,
            eqStrip(tds, 1),
            eqStrip(tds, 2),
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

