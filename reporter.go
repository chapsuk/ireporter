package ireporter

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

var version = 1.0
var salesEndpoint = "https://reportingitc-reporter.apple.com/reportservice/sales/v1"
var financeEndpoint = "https://reportingitc-reporter.apple.com/reportservice/finance/v1"

// Client is reporter client
type Client struct {
	cfg Config
}

// Config base properties
type Config struct {
	UserID   string
	Password string
	Mode     string
}

// Request to Reporter endpoints
type Request struct {
	UserID     string `json:"userid"`
	Password   string `json:"password"`
	Version    string `json:"version"`
	Mode       string `json:"mode"`
	SalesURL   string `json:"salesurl"`
	FinanceURL string `json:"financeurl"`
	QueryInput string `json:"queryInput"`
}

// NewClient yield a new Client instance
func NewClient(cfg Config) (*Client, error) {
	err := checkConfig(cfg)
	if err != nil {
		return nil, err
	}
	return &Client{
		cfg: cfg,
	}, nil
}

// GetSalesStatus return Sales.getStatus response
func (c Client) GetSalesStatus() ([]byte, error) {
	req := c.getBaseRequest()
	req.QueryInput = "%5Bp%3DReporter.properties%2C+Sales.getStatus%5D"
	return c.send(salesEndpoint, req)
}

// GetFinanceStatus return Finance.getStatus response
func (c Client) GetFinanceStatus() ([]byte, error) {
	req := c.getBaseRequest()
	req.QueryInput = "%5Bp%3DReporter.properties%2C+Finance.getStatus%5D"
	return c.send(financeEndpoint, req)
}

// GetSalesAccounts return Sales.getAccounts response
func (c Client) GetSalesAccounts() ([]byte, error) {
	req := c.getBaseRequest()
	req.QueryInput = "%5Bp%3DReporter.properties%2C+Sales.getAccounts%5D"
	return c.send(salesEndpoint, req)
}

// GetFinanceAccounts return Finance.getAccounts response
func (c Client) GetFinanceAccounts() ([]byte, error) {
	req := c.getBaseRequest()
	req.QueryInput = "%5Bp%3DReporter.properties%2C+Finance.getAccounts%5D"
	return c.send(financeEndpoint, req)
}

// GetSalesVendors return Sales.getVendors response
func (c Client) GetSalesVendors(account int) ([]byte, error) {
	if account <= 0 {
		return nil, errors.New("Wrong vendor number")
	}
	req := c.getBaseRequest()
	req.QueryInput = fmt.Sprintf("%%5Bp%%3DReporter.properties%%2C+a%%3D%d%%2C+Sales.getVendors%%5D", account)
	return c.send(salesEndpoint, req)
}

// GetFinanceVendorsAndRegions return Finance.getVendors response
func (c Client) GetFinanceVendorsAndRegions(account int) ([]byte, error) {
	if account <= 0 {
		return nil, errors.New("Wrong vendor number")
	}
	req := c.getBaseRequest()
	req.QueryInput = fmt.Sprintf("%%5Bp%%3DReporter.properties%%2C+a%%3D%d%%2C+Finance.getVendorsAndRegions%%5D", account)
	return c.send(financeEndpoint, req)
}

// GetSalesReport return Sales.getReport response (is report file or error)
func (c Client) GetSalesReport(account, vendor int, reportType, reportSubType, dateType, date string) ([]byte, error) {
	err := validateSalesReportArgs(account, vendor, reportType, reportSubType, dateType, date)
	if err != nil {
		return nil, err
	}
	req := c.getBaseRequest()
	qI := "%%5Bp%%3DReporter.properties%%2C+a%%3D%d%%2C+Sales.getReport%%2C+%d%%2C%s%%2C%s%%2C%s%%2C%s%%5D"
	req.QueryInput = fmt.Sprintf(qI, account, vendor, reportType, reportSubType, dateType, date)
	return c.send(salesEndpoint, req)
}

// GetFinanceReport return Finance.getReport response (is report file or error)
// func (c Client) GetFinanceReport() ([]byte, error) {
// TODO implement me
// }

func (c Client) send(endpoint string, r Request) ([]byte, error) {
	q, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	log.Print(string(q))
	query := fmt.Sprintf("jsonRequest=%s", string(q))
	resp, err := http.Post(endpoint, "application/x-www-form-urlencoded", strings.NewReader(query))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(body))
	}
	return body, nil
}

func (c Client) getBaseRequest() Request {
	return Request{
		UserID:     url.QueryEscape(c.cfg.UserID),
		Password:   url.QueryEscape(c.cfg.Password),
		Version:    url.QueryEscape(fmt.Sprintf("%.1f", version)),
		Mode:       url.QueryEscape(c.cfg.Mode),
		SalesURL:   url.QueryEscape(salesEndpoint),
		FinanceURL: url.QueryEscape(financeEndpoint),
	}
}

func checkConfig(cfg Config) error {
	if cfg.Mode != "Normal" && cfg.Mode != "Robot.xml" {
		return errors.New("Undefined mode. Use available modes: Normal or Robot.xml")
	}
	if cfg.UserID == "" {
		return errors.New("UserID not set")
	}
	if cfg.Password == "" {
		return errors.New("Password not set")
	}
	return nil
}

func validateSalesReportArgs(account, vendor int, reportType, reportSubType, dateType, date string) error {
	if vendor <= 0 {
		return errors.New("Wrong vendor value")
	}
	if account <= 0 {
		return errors.New("Wrong account value")
	}
	switch reportSubType {
	case "Summary",
		"Detailed",
		"Opt-In":
		break
	default:
		return errors.New("Wrong ReportSubType, use: Summary, Detailed or Opt-In")
	}

	switch dateType {
	case "Daily":
		if len(date) != 8 {
			return errors.New("Wrong DateType format for Daily Report, use: YYYYMMDD")
		}
		break
	case "Weekly":
		if len(date) != 8 {
			return errors.New("Wrong DateType format for Weekly Report, use: YYYYMMDD")
		}
		break
	case "Monthly":
		if len(date) != 6 {
			return errors.New("Wrong DateType format for Monthly Report, use: YYYYMM")
		}
		break
	case "Yearly":
		if len(date) != 4 {
			return errors.New("Wrong DateType format for Yearly Report, use: YYYY")
		}
		break
	default:
		return errors.New("Wrong DateType, use: Daily, Weekly, Monthly or Yearly")
	}

	return nil
}