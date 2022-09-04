package tabs

import (
	"fmt"
	"net"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type Tab struct {
	Name       string
	Records    binding.ExternalStringList
	CreateItem func() fyne.CanvasObject
	UpdateItem func(i binding.DataItem, o fyne.CanvasObject)
	GetResult  func(addr string) []string
}

var tabsForApp = []*Tab{
	{
		Name: "MX",
		Records: binding.BindStringList(
			&[]string{},
		),
		CreateItem: func() fyne.CanvasObject {
			return widget.NewLabel("MX records")
		},
		UpdateItem: func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		},
		GetResult: func(addr string) []string {
			hosts := []string{}

			mxRecs, err := net.LookupMX(addr)
			if err != nil {
				return []string{err.Error()}
			}

			for _, rec := range mxRecs {
				hosts = append(hosts, fmt.Sprintf("%s", rec.Host))
			}

			return hosts
		},
	},
	{
		Name: "NS",
		Records: binding.BindStringList(
			&[]string{},
		),
		CreateItem: func() fyne.CanvasObject {
			return widget.NewLabel("NS records")
		},
		UpdateItem: func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		},
		GetResult: func(addr string) []string {
			hosts := []string{}

			nsRecs, err := net.LookupNS(addr)
			if err != nil {
				return []string{err.Error()}
			}

			for _, rec := range nsRecs {
				hosts = append(hosts, rec.Host)
			}

			return hosts
		},
	},
	{
		Name: "TXT",
		Records: binding.BindStringList(
			&[]string{},
		),
		CreateItem: func() fyne.CanvasObject {
			return widget.NewLabel("NS records")
		},
		UpdateItem: func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		},
		GetResult: func(addr string) []string {
			records, err := net.LookupTXT(addr)
			if err != nil {
				return []string{err.Error()}
			}

			return records
		},
	},
	{
		Name: "IP",
		Records: binding.BindStringList(
			&[]string{},
		),
		CreateItem: func() fyne.CanvasObject {
			return widget.NewLabel("IP records")
		},
		UpdateItem: func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		},
		GetResult: func(addr string) []string {
			result := []string{}

			ips, err := net.LookupIP(addr)
			if err != nil {
				return []string{err.Error()}
			}

			for _, ip := range ips {
				result = append(result, ip.String())
			}

			return result
		},
	},
	{
		Name: "Other",
		Records: binding.BindStringList(
			&[]string{},
		),
		CreateItem: func() fyne.CanvasObject {
			return widget.NewLabel("A records")
		},
		UpdateItem: func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		},
		GetResult: func(addr string) []string {
			result := []string{}

			cname, err := net.LookupCNAME(addr)
			if err != nil {
				result = append(result, fmt.Sprintf("Error getting CNAME: %s", err))
			} else {
				result = append(result, fmt.Sprintf("CNAME: %s", cname))
			}

			return result
		},
	},
}

func Get() []*Tab {
	return tabsForApp
}
