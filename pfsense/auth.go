package pfsense

import (
    "errors"
    "strings"
    log "github.com/sirupsen/logrus"
    "github.com/levigross/grequests"
)

var (
    InvalidLogin    = errors.New("Invalid login")
    CouldNotSignOut = errors.New("Could not signout")
)

func (pf *Pfsense) Login() error {
    log.Debugf("Logging into pfsense (%s), username: %s, password: %s.", pf.conf.Url, pf.conf.User, pf.conf.Pass)

    // init our variables we will use
    var err error
    var resp *grequests.Response

    resp, err = pf.request("get", "", nil)

    if err != nil {
        log.Fatalf("Error: %v, %s\n", err, pf.conf.Url)
        return err
    }

    ro := &grequests.RequestOptions{
        Data: map[string]string{
            "__csrf_magic": pf.csrf_token,
            "usernamefld": pf.conf.User,
            "passwordfld": pf.conf.Pass,
            "login": "Sign In",
        },
    }

    resp, err = pf.request("post", "", ro)

    if err != nil {
        log.Fatalf("Error: %v\n", err)
        return err
    }

    if strings.Contains(resp.String(), "Username or Password incorrect") {
        log.Error("Failed to login to pfSense")
        return InvalidLogin
    }

    log.Debugf("Login to pfSense successfull")
    return nil
}

func (pf *Pfsense) Logout() error {
    ro := &grequests.RequestOptions{
        Data: map[string]string{
            "logout": "",
            "__csrf_magic": pf.csrf_token,
        },
    }

    resp, err := pf.request("post", "index.php?logout", ro)

    if err != nil {
        log.Fatalf("Error: %v\n", err)
        return err
    }

    log.Warn(resp.String())
    if !strings.Contains(resp.String(), "usernamefld") {
        log.Warn("Failed to logout of pfSense")
        return CouldNotSignOut
    }

    return nil
}
