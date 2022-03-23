package common

import (
	"bufio"
	"github.com/BabySid/gobase"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

var _ io.Writer = (*LogWriter)(nil)

type LogWriterOption struct {
	CacheCapacity int
	LogPath       string
	OnMessage     func(string)
}

type LogWriter struct {
	opt      LogWriterOption
	logCache *gobase.Queue
}

func NewLogWriter(opt LogWriterOption) *LogWriter {
	var lw LogWriter
	lw.opt = opt
	lw.logCache = gobase.NewQueue()
	lw.logCache.SetCapacity(opt.CacheCapacity)

	_ = lw.loadLogFromLogFile()

	return &lw
}

func (lw *LogWriter) Size() int {
	return lw.logCache.Size()
}

func (lw *LogWriter) getLogFileName() string {
	return filepath.Join(lw.opt.LogPath, "sid.log")
}

func (lw *LogWriter) loadLogFromLogFile() error {
	f, e := os.OpenFile(lw.getLogFileName(), os.O_RDWR|os.O_CREATE, 0666)
	if e != nil {
		return e
	}
	defer f.Close()

	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		lw.logCache.PushBack(line)
	}

	return nil
}

func (lw *LogWriter) flush() error {
	contStr := ""
	lw.logCache.Traversal(func(i interface{}) error {
		contStr += i.(string)
		return nil
	})

	return ioutil.WriteFile(lw.getLogFileName(), []byte(contStr), 0666)
}

func (lw *LogWriter) Write(p []byte) (int, error) {
	cont := string(p)
	if f := lw.opt.OnMessage; f != nil {
		f(cont)
	}
	lw.logCache.PushBack(cont)
	_ = lw.flush()
	return len(p), nil
}

func (lw *LogWriter) Traversal(handle func(string)) {
	lw.logCache.Traversal(func(i interface{}) error {
		handle(i.(string))
		return nil
	})
}
