package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"github.com/BabySid/gobase"
	"github.com/BabySid/proto/sodor"
	"github.com/vicanso/go-charts/v2"
	"sid-desktop/common"
	"sid-desktop/theme"
	"sort"
	"strings"
)

type metrics struct {
	win fyne.Window

	param metricsParam
}

const (
	metricsKindSodorThomas = iota
	metricsKindSodorThomasInstance
	metricsKindSodorJob
)

type metricsParam struct {
	thomasIns *sodor.ThomasInstance
	job       *sodor.Job
}

func newMetrics(kind int, param metricsParam) *metrics {
	m := metrics{}

	title := ""
	var out []*canvas.Image
	switch kind {
	case metricsKindSodorThomas:
		title = theme.AppSodorThomasListName
		out = m.loadThomasList()
	case metricsKindSodorThomasInstance:
		gobase.True(param.thomasIns != nil)
		title = fmt.Sprintf(theme.AppSodorMetricsThomasInstanceFormat,
			param.thomasIns.Thomas.Host+" "+strings.Join(param.thomasIns.Thomas.Tags, common.ArraySeparator))
	case metricsKindSodorJob:
		gobase.True(param.job != nil)
		title = fmt.Sprintf(theme.AppSodorMetricsJobInstanceFormat, param.job.Name)
	}

	m.win = globalWin.app.NewWindow(title)

	cont := container.NewAdaptiveGrid(2)
	for _, item := range out {
		cont.Add(item)
	}
	m.win.SetContent(cont)

	m.win.Resize(fyne.NewSize(800, 600))
	m.win.CenterOnScreen()
	return &m
}

type thomasMetric struct {
	Host  string
	value float64
}

type thomasMetrics []thomasMetric

func (t thomasMetrics) Len() int {
	return len(t)
}

func (t thomasMetrics) Less(i, j int) bool {
	return t[i].value > t[j].value
}

func (t thomasMetrics) Swap(i, j int) {
	t[i] = t[j]
}

func (m *metrics) loadThomasList() []*canvas.Image {
	infos := common.GetSodorCache().GetThomasInfos()
	if infos == nil {
		return nil
	}

	data := make(map[string]thomasMetrics)
	for _, info := range infos.ThomasInfos {
		infoMetrics := info.LatestMetrics.AsMap()

		for k, v := range infoMetrics {
			var fv float64
			var ok bool
			if fv, ok = v.(float64); !ok {
				continue
			}
			var arr thomasMetrics
			if arr, ok = data[k]; !ok {
				arr = make([]thomasMetric, 0)
			}
			arr = append(arr, thomasMetric{
				Host:  info.Host,
				value: fv,
			})
			data[k] = arr
		}
	}

	rs := make([]*canvas.Image, 0)
	for k, v := range data {
		sort.Sort(v)
		data[k] = v

		out, err := makeHorizontalBarRenderForThomasList(k, v)
		if err != nil {
			continue
		}

		chart := canvas.NewImageFromResource(fyne.NewStaticResource(k, out))
		chart.FillMode = canvas.ImageFillContain
		chart.ScaleMode = canvas.ImageScaleFastest
		chart.SetMinSize(fyne.NewSize(300, 300))

		rs = append(rs, chart)
	}

	return rs
}

const maxBorItems = 10
const fontSize = 8

func makeHorizontalBarRenderForThomasList(name string, metrics thomasMetrics) ([]byte, error) {
	values := make([][]float64, 1)
	size := maxBorItems
	if metrics.Len() < maxBorItems {
		size = metrics.Len()
	}
	yOpts := make([]string, size)
	for i, _ := range values {
		values[i] = make([]float64, size)
		for j := 0; j < size; j++ {
			values[i][j] = metrics[j].value
		}
	}

	for i := 0; i < size; i++ {
		yOpts[i] = metrics[i].Host
	}

	p, err := charts.HorizontalBarRender(
		values,
		charts.TitleTextOptionFunc(name),
		charts.PaddingOptionFunc(charts.Box{
			Top:    20,
			Right:  40,
			Bottom: 20,
			Left:   20,
		}),
		charts.ThemeOptionFunc(charts.ThemeAnt),
		charts.LegendLabelsOptionFunc([]string{
			name,
		}),

		charts.PieSeriesShowLabel(),
		charts.YAxisDataOptionFunc(yOpts),
		func(opt *charts.ChartOption) {
			opt.SeriesList[0].Label.FontSize = fontSize
			opt.Title.FontSize = fontSize
			opt.Legend.FontSize = fontSize
			for index := range opt.YAxisOptions {
				opt.YAxisOptions[index].FontSize = fontSize
				opt.YAxisOptions[index].DivideCount = 1
			}
			opt.XAxis.FontSize = fontSize
			opt.BarHeight = 10
		},
	)
	if err != nil {
		return nil, err
	}

	buf, err := p.Bytes()
	if err != nil {
		return nil, err
	}
	return buf, nil
}
