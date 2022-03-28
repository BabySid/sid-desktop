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

${fyne} bundle -pkg theme -name ResourceOpenDirIcon -a -o ../theme/icons.go open_dir.png
${fyne} bundle -pkg theme -name ResourceAddDirIcon -a -o ../theme/icons.go add_dir.png
${fyne} bundle -pkg theme -name ResourceRmDirIcon -a -o ../theme/icons.go rm_dir.png

${fyne} bundle -pkg theme -name ResourceConfIndexIcon -a -o ../theme/icons.go config_index.png

${fyne} bundle -pkg theme -name ResourceRunAppIcon -a -o ../theme/icons.go run_app.png
${fyne} bundle -pkg theme -name ResourceDefAppIcon -a -o ../theme/icons.go def_app.png

${fyne} bundle -pkg theme -name ResourceFavoritesIcon -a -o ../theme/icons.go favor.png
${fyne} bundle -pkg theme -name ResourceAddFavorIcon -a -o ../theme/icons.go add_favor.png
${fyne} bundle -pkg theme -name ResourceRmFavorIcon -a -o ../theme/icons.go rm_favor.png
${fyne} bundle -pkg theme -name ResourceEditFavorIcon -a -o ../theme/icons.go edit_favor.png
${fyne} bundle -pkg theme -name ResourceOpenFavorIcon -a -o ../theme/icons.go open_favor.png
${fyne} bundle -pkg theme -name ResourceExportFavorIcon -a -o ../theme/icons.go export_favor.png
${fyne} bundle -pkg theme -name ResourceImportFavorIcon -a -o ../theme/icons.go import_favor.png

${fyne} bundle -pkg theme -name ResourceExpandDownIcon -a -o ../theme/icons.go expand_down.png
${fyne} bundle -pkg theme -name ResourceExpandUpIcon -a -o ../theme/icons.go expand_up.png



# use https://github.com/lusingander/fyne-theme-generator to generate theme file