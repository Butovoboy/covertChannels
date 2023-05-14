package utils

import (
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

// generate random data for bar chart
func generateBarItems(intervals map[int]int, sortedKeys []int) []opts.BarData {
	items := make([]opts.BarData, 0)

	for i := 0; i < sortedKeys[len(sortedKeys)-1]/250; i++ {

		value, exists := intervals[i*250]
		if exists {
			items = append(items, opts.BarData{Value: value})
		} else {
			items = append(items, opts.BarData{Value: 0})
		}

	}
	return items
}

func Show_gaps(intervals map[int]int, sortedKeys []int) error {
	// create a new bar instance
	bar := charts.NewBar()
	// set some global options like Title/Legend/ToolTip or anything else
	bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title:    "Histogram of inter-batch intervals",
		Subtitle: "All intervals are sorted in ascending order",
	}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:  "inside",
			Start: 10,
			End:   50,
		}),
	)

	// Put data into instance
	bar.SetXAxis([]string{"0.25", "0.5", "0.75", "1", "1.25", "1.5", "1.75", "2", "2.25", "2.5", "2.75", "3", "3.25", "3.5", "3.75", "4"}).
		AddSeries("Packages", generateBarItems(intervals, sortedKeys)).SetSeriesOptions(
		charts.WithLabelOpts(opts.Label{
			Show:     true,
			Position: "top",
		}),
	)
	// Where the magic happens
	f, _ := os.Create("./images/bar.html")
	bar.Render(f)
	return nil
}
