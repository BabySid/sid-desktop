#!/bin/bash
# Get Icons from: https://www.iconfont.cn

${fyne} bundle -pkg theme -name ResourceSidLogo -o logo.go logo.png

# use tools/make_syso to generate sid.syso using sid.png
# run in desktop/
${make_syso} -o sid.syso -s ./resource/sid.png

${fyne} bundle -pkg theme -name ResourceAppIcon -o icons.go sid.png
${fyne} bundle -pkg theme -name ResourceWelIcon -a -o icons.go wel.png
${fyne} bundle -pkg theme -name ResourceLauncherIcon -a -o icons.go app_launcher.png
${fyne} bundle -pkg theme -name ResourceLogViewerIcon -a -o icons.go log_viewer.png
${fyne} bundle -pkg theme -name ResourceAboutIcon -a -o icons.go about.png

${fyne} bundle -pkg theme -name ResourceOpenDirIcon -a -o icons.go open_dir.png
${fyne} bundle -pkg theme -name ResourceAddDirIcon -a -o icons.go add_dir.png
${fyne} bundle -pkg theme -name ResourceRmDirIcon -a -o icons.go rm_dir.png

${fyne} bundle -pkg theme -name ResourceConfIndexIcon -a -o icons.go config_index.png

${fyne} bundle -pkg theme -name ResourceRunAppIcon -a -o icons.go run_app.png
${fyne} bundle -pkg theme -name ResourceDefAppIcon -a -o icons.go def_app.png