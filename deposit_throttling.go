package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

const (
	N                  = 2
	minDeposit float64 = 250
	alpha_up           = .1
	alpha_down         = .1
	k          float64 = 1
)

func depositThrottling() error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		boundaryName := "x5789h"
		// send multi-part header
		w.Header().Set("Content-Type", fmt.Sprintf("multipart/x-mixed-replace; boundary=%s", boundaryName))
		w.WriteHeader(http.StatusOK)

		nbBlocks := 100
		for {
			select {
			case <-r.Context().Done():
				fmt.Println("STOP")
				return
			default:
				fmt.Println("RENDER BLOCK", nbBlocks)
				filename, err := genChartFile(nbBlocks)
				if err != nil {
					panic(err)
				}

				bz, err := toPNG(r.Context(), filename)
				if err != nil {
					panic(err)
				}

				// start boundary
				io.WriteString(w, fmt.Sprintf("--%s\n", boundaryName))
				io.WriteString(w, "Content-Type: image/png\n")
				io.WriteString(w, fmt.Sprintf("Content-Length: %d\n\n", len(bz)))

				if _, err := w.Write(bz); err != nil {
					log.Printf("failed to write mjpeg image: %s", err)
					return
				}

				// close boundary
				if _, err := io.WriteString(w, "\n"); err != nil {
					log.Printf("failed to write boundary: %s", err)
					return
				}
			}
			time.Sleep(time.Second)
			nbBlocks++
		}
	})

	http.ListenAndServe(":8080", nil)
	return nil
}

func generateLineData(data []int) []opts.LineData {
	items := make([]opts.LineData, len(data))
	for i, v := range data {
		items[i] = opts.LineData{Value: v}
	}
	return items
}

func genChartFile(nbBlocks int) (string, error) {
	f, err := os.CreateTemp("", "chart*.html")
	if err != nil {
		return "", err
	}
	defer f.Close()
	page := components.NewPage()
	page.PageTitle = "Deposit throttling"
	page.AddCharts(
		newLineChart(nbBlocks),
	)
	page.Render(f)
	// browser.OpenFile(f.Name())
	return f.Name(), nil
}

func newLineChart(nbBlocks int) components.Charter {
	// create a new line instance
	line := charts.NewLine()
	line.Animation = false
	// set some global options like Title/Legend/ToolTip or anything else
	line.SetGlobalOptions(
		// charts.WithDataZoomOpts(opts.DataZoom{XAxisIndex: 0, Start: 0, End: 100, Type: "slider"}),
		charts.WithInitializationOpts(opts.Initialization{
			// Theme:  types.ThemeWesteros,
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

func toPNG(ctx context.Context, file string) ([]byte, error) {
	// get full path of file
	abs, _ := filepath.Abs(file)

	// create context
	ctx, cancel := chromedp.NewContext(
		ctx,
		// chromedp.WithDebugf(log.Printf),
	)
	defer cancel()

	// screenshot buffer
	var buf []byte

	// capture entire browser viewport, returning png with quality=100
	task := chromedp.Tasks{
		// use file:// for a local file
		chromedp.Navigate("file://" + abs),
		// set resolution for screenshot
		chromedp.EmulateViewport(1200, 800),
		// take screenshot
		chromedp.FullScreenshot(&buf, 100),
	}

	// run tasks
	err := chromedp.Run(ctx, task)
	return buf, err
}
