package gui

import (
	"bufio"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/BabySid/gobase"
	"github.com/montanaflynn/stats"
	"io"
	"os"
	"sid-desktop/theme"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

var _ devToolInterface = (*devToolStatistic)(nil)

type devToolStatistic struct {
	devToolAdapter
	fileSelectorBtn *widget.Button
	fileName        *widget.Label
	fileFormat      *widget.Select

	statisticResultText *widget.Entry
}

func (d *devToolStatistic) CreateView() fyne.CanvasObject {
	if d.content != nil {
		return d.content
	}

	d.fileSelectorBtn = widget.NewButtonWithIcon(theme.AppDevToolsOpenFileName, theme.ResourceOpenDirIcon, func() {
		fo := dialog.NewFileOpen(func(closer fyne.URIReadCloser, err error) {
			if closer != nil {
				name := closer.URI().Path()
				d.fileName.SetText(fmt.Sprintf(theme.AppDevToolsSelectedFileNameFormat, name))
				go d.runStatistic(name)
			}

		}, globalWin.win)
		fo.Show()
	})

	d.fileName = widget.NewLabel(fmt.Sprintf(theme.AppDevToolsSelectedFileNameFormat, ""))
	d.fileFormat = widget.NewSelect([]string{
		theme.AppDevToolsFileFormatNumber,
		theme.AppDevToolsFileFormatDuration,
	}, nil)
	d.fileFormat.SetSelectedIndex(0)

	d.statisticResultText = widget.NewMultiLineEntry()
	d.statisticResultText.Wrapping = fyne.TextWrapWord

	cont := container.NewBorder(
		container.NewHBox(d.fileName, layout.NewSpacer(), d.fileFormat, d.fileSelectorBtn),
		nil, nil, nil, d.statisticResultText)

	d.content = cont
	return d.content
}

func (d *devToolStatistic) runStatistic(fileName string) {
	progressBar := widget.NewProgressBarInfinite()

	progress := dialog.NewCustom(theme.AppDevToolsMathName+" - "+theme.AppDevToolsStatisticName,
		theme.DismissText,
		progressBar,
		globalWin.win)

	progress.Resize(fyne.NewSize(400, 100))

	var flag atomic.Bool
	progress.SetOnClosed(func() {
		flag.Store(true)
	})

	progressBar.Start()
	progress.Show()

	go func() {
		file, err := os.Open(fileName)
		if err != nil {
			printErr(fmt.Errorf(theme.AppDevToolsFailedFormat, err))
			return
		}

		defer file.Close()

		r := bufio.NewReader(file)
		rawData := make([]float64, 0)
		for {
			if flag.Load() {
				return
			}
			s, err := r.ReadString('\n')
			if err != nil && err != io.EOF {
				printErr(fmt.Errorf(theme.AppDevToolsFailedFormat, err))
				return
			}
			s = strings.TrimSpace(s)
			if s != "" {
				v, e := d.parseLine(s)
				if e != nil {
					printErr(fmt.Errorf(theme.AppDevToolsFailedFormat, err))
					return
				}
				rawData = append(rawData, v)
			}
			if err == io.EOF {
				break
			}
		}

		cont := ""
		if len(rawData) > 0 {
			rs, err := calcStats(rawData, &flag)
			if err != nil {
				printErr(fmt.Errorf(theme.AppDevToolsFailedFormat, err))
				return
			}
			cont = fmt.Sprintf(theme.AppDevToolsMathStatsResultFormat,
				rs.len, rs.min, rs.mean, rs.p50, rs.p75, rs.p80, rs.p90, rs.p95, rs.p99, rs.max)
		} else {
			cont = fmt.Sprintf(theme.AppDevToolsMathStatsResultFormat,
				len(rawData), 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0)
		}

		progressBar.Stop()
		progress.Hide()
		d.statisticResultText.SetText(cont)
	}()
}

func (d *devToolStatistic) parseLine(s string) (float64, error) {
	switch d.fileFormat.Selected {
	case theme.AppDevToolsFileFormatNumber:
		t, e := strconv.ParseFloat(s, 64)
		return t, e
	case theme.AppDevToolsFileFormatDuration:
		t, e := time.ParseDuration(s)
		return t.Seconds(), e
	default:
		gobase.AssertHere()
	}
	return 0, nil
}

type statsResult struct {
	len  int
	min  float64
	mean float64
	p50  float64
	p75  float64
	p80  float64
	p90  float64
	p95  float64
	p99  float64
	max  float64
}

func calcStats(input []float64, flag *atomic.Bool) (statsResult, error) {
	var sr statsResult
	sr.len = len(input)

	var err error
	if sr.len > 0 {
		if flag.Load() {
			return sr, nil
		}
		sr.min, err = stats.Min(input)
		if err != nil {
			return sr, err
		}
		if flag.Load() {
			return sr, nil
		}
		sr.mean, err = stats.Mean(input)
		if err != nil {
			return sr, err
		}
		if flag.Load() {
			return sr, nil
		}
		sr.p50, err = stats.Percentile(input, 50)
		if err != nil {
			return sr, err
		}
		if flag.Load() {
			return sr, nil
		}
		sr.p75, err = stats.Percentile(input, 75)
		if err != nil {
			return sr, err
		}
		if flag.Load() {
			return sr, nil
		}
		sr.p80, err = stats.Percentile(input, 80)
		if err != nil {
			return sr, err
		}
		if flag.Load() {
			return sr, nil
		}
		sr.p90, err = stats.Percentile(input, 90)
		if err != nil {
			return sr, err
		}
		if flag.Load() {
			return sr, nil
		}
		sr.p95, err = stats.Percentile(input, 95)
		if err != nil {
			return sr, err
		}
		if flag.Load() {
			return sr, nil
		}
		sr.p99, err = stats.Percentile(input, 99)
		if err != nil {
			return sr, err
		}
		if flag.Load() {
			return sr, nil
		}
		sr.max, err = stats.Max(input)
		if err != nil {
			return sr, err
		}
		if flag.Load() {
			return sr, nil
		}
	}

	return sr, nil
}
