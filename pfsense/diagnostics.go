package pfsense

import (
    log "github.com/sirupsen/logrus"
    "github.com/levigross/grequests"
)

func (pf *Pfsense) Backup() (string, error) {
    resp, err := pf.request("get", "diag_backup.php", nil)

    if err != nil {
        log.Fatalf("Error: %v, %s\n", err, pf.conf.Url)
    }

    ro := &grequests.RequestOptions{
        Data: map[string]string{
            "__csrf_magic": pf.csrf_token,
            "download": "download",
            "donotbackup": "yes",
        },
    }

    resp, err = pf.request("post", "diag_backup.php", ro)

    if err != nil {
        log.Fatalf("Error: %v\n", err)
        return "", err
    }

    return resp.String(), nil
}
