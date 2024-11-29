package main

import (
	"fmt"
	"math"
	"os"
	"strconv"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
	"github.com/pkg/browser"
)

const (
	N                  = 2
	minDeposit float64 = 250
	alpha_up           = .1
	alpha_down         = .1
	k          float64 = 1
)

func depositThrottling() error {
	f, err := os.CreateTemp("", "chart*.html")
	if err != nil {
		return err
	}
	defer f.Close()
	page := components.NewPage()
	page.PageTitle = "Deposit throttling"
	page.AddCharts(
		newLineChart(100),
	)
	page.Render(f)
	fmt.Printf("Charts rendered in %s\n", f.Name())
	browser.OpenFile(f.Name())
	return nil
}

func newLineChart(nbBlocks int) components.Charter {
	// create a new line instance
	line := charts.NewLine()
	// set some global options like Title/Legend/ToolTip or anything else
	line.SetGlobalOptions(
		charts.WithDataZoomOpts(opts.DataZoom{XAxisIndex: 0, Start: 0, End: 100, Type: "slider"}),
		charts.WithInitializationOpts(opts.Initialization{
			Theme:  types.ThemeWesteros,
			Height: "700px",
		}),
		charts.WithTitleOpts(opts.Title{
			Title: "Deposit throttling evolution",
			Subtitle: fmt.Sprintf(
				"N=%d\nmin_deposit=%.f\nalpha_up=%.2f\nalpha_down=%.2f\nk=%.2f",
				N, minDeposit, alpha_up, alpha_down, k),
		}),
		charts.WithGridOpts(opts.Grid{
			Top: "150px",
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Name: "Blocks",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: "Deposit",
		}),
		/*
			charts.WithVisualMapOpts(opts.VisualMap{
				Pieces: []opts.Piece{
					{
						Gte:   0,
						Lte:   float32(minDeposit),
						Color: "#112233",
					},
					{
						Gt:    float32(minDeposit),
						Color: "#666666",
					},
				},
			}),
		*/
	)
	line.ExtendYAxis(opts.YAxis{
		Name: "Num proposals",
		Show: true,
		Min:  0,
		Max:  5,
	})

	// Put data into instance
	xaxis := make([]string, nbBlocks)
	for i := 0; i < nbBlocks; i++ {
		xaxis[i] = strconv.Itoa(i + 1)
	}
	numProposals := generateNumProposal(nbBlocks)
	line.SetXAxis(xaxis).
		AddSeries("Num proposals", numProposals, charts.WithLineChartOpts(opts.LineChart{YAxisIndex: 1})).
		AddSeries("Deposit", generateDeposits(numProposals))
	return line
}

func generateNumProposal(nbBlocks int) []opts.LineData {
	items := make([]opts.LineData, nbBlocks)
	for i := 0; i < nbBlocks; i++ {
		var n int
		switch {
		case i < 10:
			n = 0
		case i < 20:
			n = 1
		case i < 30:
			n = 2
		case i < 40:
			n = 3
		case i < 50:
			n = 4
		case i < 60:
			n = 3
		case i < 70:
			n = 2
		case i < 80:
			n = 1
		default:
			n = 0
		}
		items[i].Value = n
	}
	return items
}

func generateDeposits(numProposals []opts.LineData) []opts.LineData {
	items := make([]opts.LineData, len(numProposals))
	for i := 0; i < len(items); i++ {
		var (
			n           = numProposals[i].Value.(int)
			alpha       = alpha_up
			beta        = 0
			lastDeposit = minDeposit
		)
		if n <= N {
			alpha = alpha_down
			beta = -1
		}
		if i > 0 {
			lastDeposit = items[i-1].Value.(float64)
		}
		v := lastDeposit * (1 + alpha*math.Pow(float64(n-N+beta), k))
		if v < minDeposit {
			v = minDeposit
		}
		items[i].Value = v
	}
	return items
}
