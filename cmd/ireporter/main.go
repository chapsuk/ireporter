package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/chapsuk/ireporter/reporter"
)

const (
	salesHelp = `Sales commands include:
	 getHelp: Returns this help message. No arguments.
	 getStatus: Returns status of Sales and Trends application. No arguments.
	 getAccounts: Returns list of available accounts. No arguments.
	 getVendors: Returns list of available vendor numbers. No arguments.
	 getReport: Downloads a report. Arguments: Vendor Number, Report Type, Report Subtype, DateType, Date.
For more details, see Reporter guide: http://help.apple.com/itc/appsreporterguide/#/itcbe21ac7db`
	financeHelp = `Finance commands include:
	 getHelp: Returns this help message. No arguments.
	 getStatus: Returns status of Financial reporting application. No arguments.
	 getAccounts: Returns list of available accounts. No arguments.
	 getVendorsAndRegions: Returns list of available vendors, regions, and report types. No arguments.
	 getReport: Downloads a report. Arguments: Vendor Number, Region Code, Report Type, Fiscal Year, Fiscal Period.
For more details, see Reporter guide in the Resources and Help section on iTunes Connect.`
)

var (
	accessToken = flag.String("accessToken", "", "Your iTunes Connect Reports Access Token")
	mode        = flag.String("mode", "Normal", `Reporter has two modes of operation: Normal and Robot.
Normal mode is intended for an actual user that executes ad-hoc commands. Messages are displayed in easily readable text.
Robot mode is intended for an automated script that’s used regularly. Messages in robot mode are displayed in XML for easy parsing.`)

	app = flag.String("app", "", "Sales or Finance")
	cmd = flag.String("cmd", "getHelp", "Command (for example, getHelp).")

	account = flag.Int("account", 0, "If your Apple ID has access to multiple accounts, you’ll need to specify the account number you’d like to use.")
	vendor  = flag.Int("vendor", 0, "Vendor number of the report to download. For a list of your vendor numbers, use the getVendors command.")

	reportType = flag.String("reportType", "", "Type of information contained in report (for example, Sales or Financial).")
	// Sales reports
	reportSubType = flag.String("reportSubtype", "", "Level of detail in the report (for example, Summary).")
	dateType      = flag.String("dateType", "", "Length of time covered by the report (for example, Daily or Weekly).")
	date          = flag.String("date", "", "Specific time covered by the report (for example, 20150201).")
	// Finance reports
	regionCode   = flag.String("regionCode", "", "Two-character code of country of report to download. For a list of country codes by vendor number, use getVendorsAndRegions command.")
	fiscalYear   = flag.Int("fiscalYear", 0, "Four-digit year of report to download. Year is specific to Apple’s [fiscal calendar](https://itunesconnect.apple.com/WebObjects/iTunesConnect.woa/wa/jumpTo?page=fiscalcalendar).")
	fiscalPeriod = flag.Int("fiscalPeriod", 0, "This is the period in fiscal year for the report you’re downloading (1–12). The period is specific to Apple’s [fiscal calendar](https://itunesconnect.apple.com/WebObjects/iTunesConnect.woa/wa/jumpTo?page=fiscalcalendar).")
)

func main() {
	flag.Parse()

	cfg := reporter.Config{
		AccessToken: *accessToken,
		Mode:        *mode,
	}

	cli, err := reporter.NewClient(cfg)
	handleError(err)

	if *app == "Sales" {
		salesCommand(cli)
	} else if *app == "Finance" {
		financeCommand(cli)
	} else {
		flag.PrintDefaults()
	}
	fmt.Print("\n")
}

func financeCommand(cli *reporter.Client) {
	switch *cmd {
	case "getStatus":
		res, err := cli.GetSalesStatus()
		handleError(err)
		fmt.Print(string(res))
	case "getAccounts":
		res, err := cli.GetSalesAccounts()
		handleError(err)
		fmt.Print(string(res))
	case "getVendors":
		res, err := cli.GetSalesVendors(*account)
		handleError(err)
		fmt.Print(string(res))
	case "getVendorsAndRegions":
		res, err := cli.GetFinanceVendorsAndRegions(*account)
		handleError(err)
		fmt.Print(string(res))
	case "getReport":
		res, err := cli.GetFinanceReport(*account, *vendor, *regionCode, *reportType, *fiscalYear, *fiscalPeriod)
		handleError(err)
		fileName := fmt.Sprintf("FinanceReport_%s.gz", *date)
		ioutil.WriteFile(fileName, res, 0644)
		fmt.Printf("Finance report saved to %s", fileName)
	default:
		fmt.Print(financeHelp)
	}
}

func salesCommand(cli *reporter.Client) {
	switch *cmd {
	case "getStatus":
		res, err := cli.GetSalesStatus()
		handleError(err)
		fmt.Print(string(res))
	case "getAccounts":
		res, err := cli.GetSalesAccounts()
		handleError(err)
		fmt.Print(string(res))
	case "getVendors":
		res, err := cli.GetSalesVendors(*account)
		handleError(err)
		fmt.Print(string(res))
	case "getReport":
		res, err := cli.GetSalesReport(*account, *vendor, *reportType, *reportSubType, *dateType, *date)
		handleError(err)
		fileName := fmt.Sprintf("SalesReport_%s.gz", *date)
		ioutil.WriteFile(fileName, res, 0644)
		fmt.Printf("Report saved to %s", fileName)
	default:
		fmt.Print(salesHelp)
	}
}

func handleError(err error) {
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
}
