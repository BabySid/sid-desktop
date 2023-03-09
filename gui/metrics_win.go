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
	win     fyne.Window
	content *fyne.Container

	param metricsParam
	kind  int

	title  string
	images []*canvas.Image
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
	m.kind = kind
	m.param = param

	m.win = globalWin.app.NewWindow(m.title)

	m.content = container.NewAdaptiveGrid(2)
	m.win.SetContent(m.content)

	m.win.Resize(fyne.NewSize(800, 600))
	m.win.CenterOnScreen()

	m.refresh()

	return &m
}

func (m *metrics) refresh() {
	switch m.kind {
	case metricsKindSodorThomas:
		m.title = fmt.Sprintf("Top %d Metrics of %s", maxBorItems, theme.AppSodorThomasListName)
		m.images = m.loadThomasList()
	case metricsKindSodorThomasInstance:
		gobase.True(m.param.thomasIns != nil)
		m.title = fmt.Sprintf(theme.AppSodorMetricsThomasInstanceFormat,
			m.param.thomasIns.Thomas.Host+" "+strings.Join(m.param.thomasIns.Thomas.Tags, common.ArraySeparator))
		m.images = m.loadThomasInstance()
	case metricsKindSodorJob:
		gobase.True(m.param.job != nil)
		m.title = fmt.Sprintf(theme.AppSodorMetricsJobInstanceFormat, m.param.job.Name)
	}

	m.content.RemoveAll()
	for _, item := range m.images {
		m.content.Add(item)
	}
	m.content.Refresh()
}

// //////////////////////////////////////////// Thomas Metric //////////////////////////////////////////
type thomasMetric struct {
	Host  string
	value float64
}

type thomasMetrics []thomasMetric

func (t thomasMetrics) Len() int {
	return len(t)
}

func (t thomasMetrics) Less(i, j int) bool {
	return t[i].value < t[j].value
}

func (t thomasMetrics) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
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

		out, err := makeHorizontalBarRenderForThomasList(k, v)
		if err != nil {
			printErr(fmt.Errorf(theme.ProcessSodorFailedFormat, err))
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

// //////////////////////////////////////////// Thomas Instance Metric //////////////////////////////////////////
type thomasInstanceMetric struct {
	Ts    int32
	value float64
}

func (m *metrics) loadThomasInstance() []*canvas.Image {
	data := make(map[string][]thomasInstanceMetric)
	for _, ms := range m.param.thomasIns.Metrics {
		infoMetrics := ms.Metrics.AsMap()

		for k, v := range infoMetrics {
			var fv float64
			var ok bool
			if fv, ok = v.(float64); !ok {
				continue
			}
			var arr []thomasInstanceMetric
			if arr, ok = data[k]; !ok {
				arr = make([]thomasInstanceMetric, 0)
			}
			arr = append(arr, thomasInstanceMetric{
				Ts:    ms.CreateAt,
				value: fv,
			})
			data[k] = arr
		}
	}

	rs := make([]*canvas.Image, 0)
	for k, v := range data {
		out, err := makeBarRenderForThomasInstance(k, v)
		if err != nil {
			printErr(fmt.Errorf(theme.ProcessSodorFailedFormat, err))
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

func makeBarRenderForThomasInstance(name string, metrics []thomasInstanceMetric) ([]byte, error) {
	values := make([][]float64, 1)
	size := len(metrics)
	xOpts := make([]string, size)
	for i, _ := range values {
		values[i] = make([]float64, size)
		for j := 0; j < size; j++ {
			values[i][j] = metrics[j].value
		}
	}

	for i := 0; i < size; i++ {
		xOpts[i] = gobase.FormatTimeStampWithFormat(int64(metrics[i].Ts), gobase.TimeFormat)
	}

	p, err := charts.BarRender(
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

		//charts.PieSeriesShowLabel(),
		charts.MarkPointOptionFunc(0, charts.SeriesMarkDataTypeMax),
		charts.XAxisDataOptionFunc(xOpts),
		func(opt *charts.ChartOption) {
			opt.SeriesList[0].Label.FontSize = fontSize
			opt.Title.FontSize = fontSize
			opt.Legend.FontSize = fontSize
			for index := range opt.YAxisOptions {
				opt.YAxisOptions[index].FontSize = fontSize
				opt.YAxisOptions[index].DivideCount = 1
			}
			opt.XAxis.FontSize = fontSize
			opt.BarWidth = 10
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
