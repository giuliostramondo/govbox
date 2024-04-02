package main

import (
	"os"

	"github.com/olekukonko/tablewriter"
)

func newMarkdownTable(headers ...string) *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.SetHeader(headers)
	return table
}
