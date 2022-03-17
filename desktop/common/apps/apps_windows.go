package apps

import (
	"bytes"
	"image"
	"image/png"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"sid/base"
	"strings"
	"time"

	lnk "github.com/parsiya/golnk"
	"github.com/tc-hib/winres"
)

var (
	DefaultAppPaths = make([]string, 0)
	envExpandRegexp = regexp.MustCompile(`%(.*?)%`)
)

const (
	envExpandRegexpReplace = "${$1}"
)

func init() {
	homeDir, _ := os.UserHomeDir()
	if homeDir != "" {
		DefaultAppPaths = append(DefaultAppPaths, filepath.Join(homeDir[:3], "ProgramData\\Microsoft\\Windows\\Start Menu\\Programs"))
		DefaultAppPaths = append(DefaultAppPaths, filepath.Join(homeDir, "AppData\\Roaming\\Microsoft\\Windows\\Start Menu\\Programs"))
		DefaultAppPaths = append(DefaultAppPaths, filepath.Join(homeDir, "AppData\\Roaming\\Microsoft\\Internet Explorer\\Quick Launch"))
	}

	disks := base.GetDiskPartitions()
	programFiles := []string{
		"Program Files (x86)",
		"Program Files",
	}

	for _, disk := range disks {
		for _, program := range programFiles {
			DefaultAppPaths = append(DefaultAppPaths, filepath.Join(disk, string(os.PathSeparator), program))
		}
	}
}

func InitApps(searchPath []string) (*AppList, error) {
	apps := NewAppList()
	for _, path := range searchPath {
		err := findRecursive(apps, path, 0)
		if err != nil {
			return apps, err
		}
	}

	return apps, nil
}

func findRecursive(apps *AppList, dir string, level int) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	for _, f := range entries {
		ext := filepath.Ext(f.Name())
		if ext == ".lnk" || ext == ".exe" {
			app := AppInfo{
				AppName:    f.Name(),
				FullPath:   filepath.Join(dir, f.Name()),
				CreateTime: time.Now().Unix(),
				AccessTime: -1,
			}
			if ext == ".lnk" {
				err = extractIconFromLnk(&app)
			}
			if ext == ".exe" {
				err = extractIconFromExe(&app)
			}
			apps.Append(app)
		} else if level < 5 && f.IsDir() {
			err := findRecursive(apps, filepath.Join(dir, f.Name()), level+1)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func extractIconFromExe(app *AppInfo) error {
	icon, err := extractIconFromResFile(app.FullPath, 0)
	if err != nil {
		return err
	}

	if len(icon) > 0 {
		app.Icon = icon
	}
	return nil
}

func extractIconFromLnk(app *AppInfo) error {
	f, err := lnk.File(app.FullPath)
	if err != nil {
		return err
	}

	path := f.StringData.IconLocation
	if path == "" {
		path = f.LinkInfo.LocalBasePath
	}
	if path == "" {
		path = f.LinkInfo.LocalBasePathUnicode
	}
	if path == "" {
		path = f.LinkInfo.LocalBasePathUnicode
	}
	if path == "" {
		path = strings.TrimPrefix(strings.Split(f.StringData.NameString, ",")[0], "@")
	}
	if path == "" {
		return nil
	}

	path = strings.ToLower(os.ExpandEnv(envExpandRegexp.ReplaceAllString(path, envExpandRegexpReplace)))
	ext := filepath.Ext(path)

	if ext == ".dll" && strings.Contains(path, "system32") {
		sysRes := strings.Replace(path, "system32", "systemresources", 1) + ".mun"
		if _, err := os.Stat(sysRes); err == nil {
			path = sysRes
		}
	}

	var icon []byte
	switch ext {
	case ".exe", ".dll":
		icon, err = extractIconFromResFile(path, int(f.Header.IconIndex))
	case ".ico":
		icon, err = extractIconFromIco(path)
	}

	if len(icon) > 0 {
		app.Icon = icon
	}

	return err
}

func extractIconFromResFile(path string, idx int) (icon []byte, err error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	rs, err := winres.LoadFromEXE(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	te := rs.Types[winres.RT_GROUP_ICON]
	if te == nil {
		return nil, nil
	}

	var img []byte
	if idx < 0 {
		resID := winres.ID(math.Abs(float64(idx)))
		img = getIconFromResourceSet(rs, resID)
		if img != nil {
			return img, nil
		}
	}

	te.Order()
	for _, resID := range te.OrderedKeys {
		img = getIconFromResourceSet(rs, resID)
		if img != nil {
			break
		}
	}

	return img, nil
}

func extractIconFromIco(path string) ([]byte, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	ico, err := winres.LoadICO(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	var img []byte
	for _, i := range ico.Images {
		if len(img) < len(i.Image) {
			img = i.Image
		}
	}

	return img, nil
}

func getIconFromResourceSet(rs *winres.ResourceSet, resID winres.Identifier) []byte {
	ico, err := rs.GetIcon(resID)
	if err != nil {
		return nil
	}

	var img []byte
	for _, i := range ico.Images {
		data, err := extractIconToPNG(i.Image)
		if err == nil && len(img) < len(data) {
			img = data
		}
	}
	return img
}

func extractIconToPNG(imgBytes []byte) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = png.Encode(&buf, img)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
