package theme

const (
	AppTitle   = "SID Desktop"
	WelComeMsg = "Welcome to Sid Desktop"

	QuitAppTitle = "Want To Quit?"
	QuitAppMsg   = "Are you sure to want to quit?"

	DismissText = "Cancel"
	ConfirmText = "OK"

	CannotCloseTitle = "Cannot Close"

	// AppWelcomeName begins for built-in apps
	AppWelcomeName = "WelCome"

	AppLauncherName                      = "App Launcher"
	AppLauncherSearchText                = "Type the name of application or `enter` to run directly"
	AppLauncherExplorerText              = "Explorer"
	AppLauncherConfigBtnText             = "Config"
	AppLauncherCannotCloseMsg            = "Config Window is running now. Close it first, please"
	AppLauncherAppListHeader1            = "AppIcon/Name"
	AppLauncherAppListHeader2            = "Last Access Time"
	AppLauncherAppListOp1                = "Open File Location"
	AppLauncherAppListOp2                = "Run"
	AppLauncherNeedInitTitle             = "Need Init Index"
	AppLauncherNeedInitMsg               = "App index is not built, want to initialize the index?"
	AppLauncherConfigTitle               = AppLauncherName + "-" + AppLauncherConfigBtnText
	AppLauncherConfigAddDirBtn           = "Add Dir"
	AppLauncherConfigRmDirBtn            = "Remove"
	AppLauncherConfigFileFilter          = "File Mask:"
	AppLauncherConfigBuildBtn            = "Build"
	AppLauncherConfigIndexBuilding       = "Index is building"
	AppLauncherConfigStartScanApp        = "Begin to scan app..."
	AppLauncherConfigFinishScanAppFormat = "Scan App finished. Found %d apps"
	AppLauncherConfigCannotCloseMsg      = "Index is building. Cannot close"

	AppFavoritesName             = "Favorites"
	AppFavoritesSearchText       = "Type the name, url or tags of favorites"
	AppFavoritesImportBtnText    = "Import"
	AppFavoritesExportBtnText    = "Export"
	AppFavoritesAddFavorBtnText  = "Add"
	AppFavoritesRmFavorBtnText   = "Remove"
	AppFavoritesFavorListOp1     = "Edit"
	AppFavoritesFavorListOp2     = "Open"
	AppFavoritesFavorListHeader1 = "Name(Tags)"
	AppFavoritesFavorListHeader2 = "Url"
	AppFavoritesAddFavorTitle    = "Add Favor"
	AppFavoritesAddFavorName     = "Name"
	AppFavoritesAddFavorUrl      = "Url"
	AppFavoritesAddFavorTags     = "Tags"
	AppFavoritesAddFavorExpand   = "Expand"
	AppFavoritesAddFavorShrink   = "Shrink"

	LogViewerRefreshBtn = "Refresh"
	LogViewerTitle      = "Show Log"

	AboutTitle = "About"
	AboutIntro = `## Sid Desktop  

Sid desktop is a desktop software based on [Fyne](https://fyne.io/),  

which is purely built by personal interests.`

	OpenAppLocationFailedFormat = "open location for %s failed. %s"
	RunAppFailedFormat          = "run %s failed. %s"
	UpdateAppIndexFailedFormat  = "update index for %s failed. %s"
	RunCommandFailedFormat      = "run command %s failed. %s"
	RunExplorerFailedFormat     = "run explorer to %s failed. %s"
	RunAppIndexFailedFormat     = "app index run failed. %s"
	OpenFavoritesFailedFormat   = "open favorites failed. %s"
	ExportFavoritesFailedFormat = "export favorites failed. %s"
	ImportFavoritesFailedFormat = "import favorites failed. %s"

	// MenuSys begins for menus
	MenuSys     = "System"
	MenuSysQuit = "Quit"

	MenuOption        = "Option"
	MenuOptTheme      = "Theme"
	MenuOptThemeDark  = "Dark"
	MenuOptThemeLight = "Light"
	MenuOptFullScreen = "FullScreen"

	MenuHelp      = "Help"
	MenuHelpLog   = "Show Log"
	MenuHelpAbout = "About Sid"

	// ToyDateTimeTitle begins for toys
	ToyDateTimeTitle               = "DateTime"
	ToyResourceMonitorTitle        = "Monitor"
	ToyResourceMonitorItem1        = "CPU:"
	ToyResourceMonitorItem2        = "MEM:"
	ToyResourceMonitorUpTimeFormat = "Up time %d:%d:%d"
	ToyHotSearchTitle              = "Hot Search(BD)"
	ToyHotSearchRefreshing         = "Content is Updating"
	ToyHotSearchRefreshFormat      = "Refresh (%d/%d)"

	NetWorkErrorFormat   = "network error. %s"
	InvalidContentFormat = "invalid content. %s"
	InternalErrorFormat  = "there is a internal error: %s"
)
