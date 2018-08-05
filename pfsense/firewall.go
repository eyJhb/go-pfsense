package pfsense

import (
    "strings"
    log "github.com/sirupsen/logrus"
    "github.com/PuerkitoBio/goquery"
)


type firewall_rule struct{
    Disabled bool
    Action string 
    States string
    Protocol string
    Src string
    SrcPort string
    Dst string
    DstPort string
    Gateway string
    Queue string
    Schedule string
    Desc string
}

func (pf *Pfsense) Rules(rif string) ([]firewall_rule, error) {
    resp, err := pf.request("get", "firewall_rules.php?if="+rif, nil)

    if err != nil {
        log.Fatalf("Error: %v, %s\n", err, pf.conf.Url)
    }

    doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp.String()))

    if err != nil {
        log.Fatal("Error: %v\n", err)
        return nil, err
    }

    var firewall_rules []firewall_rule
    first := true
    var disabled bool
    var action string
    doc.Find("table").Eq(0).Find("tr").Each(func(i int, s *goquery.Selection) {
        if first {
            first = false
            return
        }

        tds := s.Find("td")
        disabled = false

        if v, e := s.Attr("class");e == true && strings.Contains(v, "disabled") {
            disabled = true
        }

        v, exists := tds.Eq(1).Attr("title")
        if exists == false {
            log.Fatal("Title does not exists.. It should!")
        }
        action = v[11:]

        firewall_rules = append(firewall_rules, firewall_rule{
            disabled,
            action,
            eqStrip(tds, 2),
            eqStrip(tds, 3),
            eqStrip(tds, 4),
            eqStrip(tds, 5),
            eqStrip(tds, 6),
            eqStrip(tds, 7),
            eqStrip(tds, 8),
            eqStrip(tds, 9),
            eqStrip(tds, 10),
            eqStrip(tds, 11),
        })
    })

    return firewall_rules, nil
}
