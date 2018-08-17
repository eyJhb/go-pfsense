package cmd

import (
	"fmt"
    "strings"
	"pfsense/pfsense"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func initPfsense() (*pfsense.Pfsense, error) {
    pf, err := pfsense.New(
        pfsense.Config{
            Url: viper.GetString("url"),
            User: viper.GetString("username"),
            Pass: viper.GetString("password"),
    })

    if err != nil {
        fmt.Println("Could not init pfsense object: ", err)
        return nil, err
    }

    err = pf.Login()
    if err != nil {
        fmt.Println("Could not login to pfsense: ", err)
        return nil, err
    }

    return pf, nil
}

var get = &cobra.Command{
	Use:   "get",
	Short: "Get some information from pfsense",
	Long:  `Get some information from pfsense`,
}

var dhcp_leases = &cobra.Command{
	Use:   "dhcp",
	Short: "get dhcp leases",
	Long:  `get dhcp leases`,
	Run: func(cmd *cobra.Command, args []string) {
        pf, err := initPfsense()

        if err != nil {
            fmt.Println("Could not init pfsense: ", err)
            return
        }

        leases, _ := pf.GetDhcp()

        header := strings.Split(cmd.Flag("info").Value.String(), ",")

        // generate
        var data [][]string
        for _, v := range leases {
            temp := []string{}
            for _, h := range header {
                switch h {
                    case "manual":
                        temp = append(temp, v.Manual)
                    case "ip":
                        temp = append(temp, v.Ip)
                    case "mac":
                        temp = append(temp, v.Mac)
                    case "manufacture":
                        temp = append(temp, v.Manufacture)
                    case "hostname":
                        temp = append(temp, v.Hostname)
                    case "desc":
                        temp = append(temp, v.Desc)
                    case "start":
                        temp = append(temp, v.Start)
                    case "end":
                        temp = append(temp, v.End)
                    case "online":
                        temp = append(temp, v.Online)
                    case "ltype":
                        temp = append(temp, v.Ltype)
                }
            }
            data = append(data, temp)
        }
        pfsense.MakeTable(header, data)
        return
	},
}

var rules = &cobra.Command{
	Use:   "rules [lan,wan,etc]",
	Short: "get rules - can be a comma seperated list",
	Long:  `get rules - can be a comma seperated list`,
    Args: cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        pf, err := initPfsense()

        if err != nil {
            fmt.Println("Could not init pfsense: ", err)
            return
        }

        header := []string{
            "disabled",
            "action",
            "states",
            "src",
            "src-port",
            "dst",
            "dst-port",
            "gateway",
            "desc",
        }

        for _, rif := range strings.Split(args[0], ",") {
            fmt.Println(rif)
            rules, err := pf.Rules(rif)

            if err != nil {
                fmt.Println("Could not get rules for interface..")
                return
            }

            var data [][]string
            for _, v := range rules {
                disabled := "no"
                if v.Disabled {
                    disabled = "yes"
                }
                data = append(data, []string{
                    disabled,
                    v.Action,
                    v.States,
                    v.Src,
                    v.SrcPort,
                    v.Dst,
                    v.DstPort,
                    v.Gateway,
                    v.Desc,
                })

            }
            pfsense.MakeTable(header, data)

        }
        return
	},
}

var backup = &cobra.Command{
	Use:   "backup",
	Short: "get the backup.xml as string",
	Long:  `get rules - can be a comma seperated list`,
    Run: func(cmd *cobra.Command, args []string) {
        pf, err := initPfsense()

        if err != nil {
            fmt.Println("Could not init pfsense: ", err)
            return
        }

        backup, err := pf.Backup()

        if err != nil {
            fmt.Println("Could not get backup: ", err)
        }

        fmt.Println(backup)

        return
	},
}

func init() {
	rootCmd.AddCommand(get)
	get.AddCommand(dhcp_leases)
	get.AddCommand(rules)
	get.AddCommand(backup)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	dhcp_leases.PersistentFlags().String("info", "ip,mac,hostname,desc,start,end,online,ltype", "lease information to display - manual,ip,mac,manufacture,hostname,desc,start,end,online,ltype")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// helloCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
