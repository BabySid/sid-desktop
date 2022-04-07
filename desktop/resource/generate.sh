#!/bin/bash
# Get Icons from: https://www.iconfont.cn

${fyne} bundle -pkg theme -name ResourceSidLogo -o logo.go logo.png

# use tools/make_syso to generate sid.syso using sid.png
# run in desktop/
${make_syso} -o sid.syso -s ./resource/sid.png

${fyne} bundle -pkg theme -name ResourceAppIcon -o icons.go sid.png
${fyne} bundle -pkg theme -name ResourceWelIcon -a -o ../theme/icons.go wel.png
${fyne} bundle -pkg theme -name ResourceLauncherIcon -a -o ../theme/icons.go app_launcher.png
${fyne} bundle -pkg theme -name ResourceLogViewerIcon -a -o ../theme/icons.go log_viewer.png
${fyne} bundle -pkg theme -name ResourceAboutIcon -a -o ../theme/icons.go about.png
${fyne} bundle -pkg theme -name ResourceScriptRunnerIcon -a -o ../theme/icons.go script.png
${fyne} bundle -pkg theme -name ResourceDevToolsIcon -a -o ../theme/icons.go dev_tools.png

${fyne} bundle -pkg theme -name ResourceOpenDirIcon -a -o ../theme/icons.go open_dir.png
${fyne} bundle -pkg theme -name ResourceAddDirIcon -a -o ../theme/icons.go add_dir.png
${fyne} bundle -pkg theme -name ResourceRmDirIcon -a -o ../theme/icons.go rm_dir.png

${fyne} bundle -pkg theme -name ResourceConfIndexIcon -a -o ../theme/icons.go config_index.png

${fyne} bundle -pkg theme -name ResourceRunIcon -a -o ../theme/icons.go run.png
${fyne} bundle -pkg theme -name ResourceStopIcon -a -o ../theme/icons.go stop.png
${fyne} bundle -pkg theme -name ResourceSearchIcon -a -o ../theme/icons.go search.png
${fyne} bundle -pkg theme -name ResourceClearIcon -a -o ../theme/icons.go clear.png
${fyne} bundle -pkg theme -name ResourcePrettyIcon -a -o ../theme/icons.go pretty.png
${fyne} bundle -pkg theme -name ResourceCompressIcon -a -o ../theme/icons.go compress.png
${fyne} bundle -pkg theme -name ResourceEncodeIcon -a -o ../theme/icons.go encode.png
${fyne} bundle -pkg theme -name ResourceDecodeIcon -a -o ../theme/icons.go decode.png
${fyne} bundle -pkg theme -name ResourceDefAppIcon -a -o ../theme/icons.go def_app.png

${fyne} bundle -pkg theme -name ResourceFavoritesIcon -a -o ../theme/icons.go favor.png
${fyne} bundle -pkg theme -name ResourceAddIcon -a -o ../theme/icons.go add.png
${fyne} bundle -pkg theme -name ResourceRmIcon -a -o ../theme/icons.go rm.png
${fyne} bundle -pkg theme -name ResourceEditIcon -a -o ../theme/icons.go edit.png
${fyne} bundle -pkg theme -name ResourceSaveIcon -a -o ../theme/icons.go save.png
${fyne} bundle -pkg theme -name ResourceOpenUrlIcon -a -o ../theme/icons.go open_url.png
${fyne} bundle -pkg theme -name ResourceExportIcon -a -o ../theme/icons.go export.png
${fyne} bundle -pkg theme -name ResourceImportIcon -a -o ../theme/icons.go import.png

${fyne} bundle -pkg theme -name ResourceExpandDownIcon -a -o ../theme/icons.go expand_down.png
${fyne} bundle -pkg theme -name ResourceExpandUpIcon -a -o ../theme/icons.go expand_up.png

${fyne} bundle -pkg theme -name ResourceLuaIcon -a -o ../theme/icons.go lua.png
${fyne} bundle -pkg theme -name ResourceHttpIcon -a -o ../theme/icons.go http.png

# use https://github.com/lusingander/fyne-theme-generator to generate theme file