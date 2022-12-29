package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/BabySid/gobase"
	"github.com/PuerkitoBio/goquery"
	"log"
	"math"
	"net/http"
	sidTheme "sid-desktop/theme"
	"strconv"
	"sync"
)

var _ toyInterface = (*toyHotSearch)(nil)

type toyHotSearch struct {
	toyAdapter
	hotList        *widget.List
	hotListBinding binding.UntypedList
	hotLinksMutex  *sync.Mutex
	hotLinks       []hotSearchURL

	titleLabel *widget.Label
	curPage    int
	nextPage   *widget.Button
}

func (ths *toyHotSearch) Init() error {
	ths.titleLabel = widget.NewLabel(sidTheme.ToyHotSearchTitle)
	ths.nextPage = widget.NewButton(sidTheme.ToyHotSearchRefreshing, func() {
		ths.setNextPageText()
		ths.refreshHotLinkList()
	})
	ths.nextPage.Disable()

	ths.hotLinksMutex = &sync.Mutex{}

	ths.hotListBinding = binding.NewUntypedList()
	ths.hotList = widget.NewListWithData(
		ths.hotListBinding,
		func() fyne.CanvasObject {
			return widget.NewHyperlink("", nil)
		},
		func(data binding.DataItem, item fyne.CanvasObject) {
			o, _ := data.(binding.Untyped).Get()
			url := o.(hotSearchURL)
			src := strconv.Itoa(url.id) + "." + url.text
			txt := gobase.CutUTF8(src, 0, 16, "...")
			item.(*widget.Hyperlink).SetText(txt)
			_ = item.(*widget.Hyperlink).SetURLFromString(url.url)
		},
	)

	ths.curPage = -1
	ths.widget = widget.NewCard("", "",
		container.NewBorder(container.NewHBox(ths.titleLabel, layout.NewSpacer(), ths.nextPage),
			nil, nil, nil, ths.hotList),
	)
	ths.widget.Resize(fyne.NewSize(ToyWidth, 250))

	// for init
	go ths.Run()
	_ = gobase.GlobalScheduler.AddJob("toy_hot_search", "0 */30 * * * *", ths)

	return nil
}

var (
	urlNumOfPerPage = 5
)

type hotSearchURL struct {
	id   int
	text string
	url  string
}

func (ths *toyHotSearch) Run() {
	tmpLinks := make([]hotSearchURL, 0)

	res, err := http.Get("https://top.baidu.com/board?tab=realtime")
	if err != nil {
		printErr(fmt.Errorf(sidTheme.NetWorkErrorFormat, err))
		return
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		err := fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
		printErr(fmt.Errorf(sidTheme.NetWorkErrorFormat, err))
		return
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		printErr(fmt.Errorf(sidTheme.InvalidContentFormat, err))
		return
	}

	doc.Find(".category-wrap_iQLoo").Each(func(i int, s *goquery.Selection) {
		url, _ := s.Find("a").Attr("href")
		txt := s.Find(".title_dIF3B .c-single-text-ellipsis").Text()
		link := hotSearchURL{
			id:   i + 1,
			text: txt,
			url:  url,
		}

		tmpLinks = append(tmpLinks, link)
	})

	ths.updateHotLinks(tmpLinks)

	ths.setNextPageText()
	ths.refreshHotLinkList()
	if ths.nextPage.Disabled() {
		ths.nextPage.Enable()
	}

	log.Printf(sidTheme.ToyHotSearchUpdateComplete)
}

func (ths *toyHotSearch) updateHotLinks(urls []hotSearchURL) {
	ths.hotLinksMutex.Lock()
	defer ths.hotLinksMutex.Unlock()
	ths.hotLinks = urls
}

func (ths *toyHotSearch) getHotLinks(start, end int) []hotSearchURL {
	ths.hotLinksMutex.Lock()
	defer ths.hotLinksMutex.Unlock()
	return ths.hotLinks[start:end]
}

func (ths *toyHotSearch) setNextPageText() {
	ths.hotLinksMutex.Lock()
	defer ths.hotLinksMutex.Unlock()

	totalPage := int(math.Ceil(float64(len(ths.hotLinks)) / float64(urlNumOfPerPage)))
	ths.curPage = (ths.curPage + 1) % totalPage
	ths.nextPage.SetText(fmt.Sprintf(sidTheme.ToyHotSearchRefreshFormat, ths.curPage+1, totalPage))
}

func (ths *toyHotSearch) refreshHotLinkList() {
	begin := 0
	if urlNumOfPerPage*ths.curPage > begin {
		begin = urlNumOfPerPage * ths.curPage
	}

	end := begin + urlNumOfPerPage
	if end > len(ths.hotLinks) {
		end = len(ths.hotLinks)
	}

	curPage := ths.getHotLinks(begin, end)

	_ = ths.hotListBinding.Set(convertHotURLsToInterfaceArray(curPage))

	// https://github.com/fyne-io/fyne/issues/2843
	ths.hotList.Refresh()
}

func convertHotURLsToInterfaceArray(src []hotSearchURL) []interface{} {
	rs := make([]interface{}, len(src), len(src))
	for i := range src {
		rs[i] = src[i]
	}
	return rs
}
