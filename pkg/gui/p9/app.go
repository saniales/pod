package p9

import (
	"fmt"

	l "gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"golang.org/x/exp/shiny/materialdesign/icons"

	"github.com/p9c/pod/pkg/gui/wallet/ico"
)

// App defines an application with a header, sidebar/menu, right side button bar, changeable body page widget and
// pop-over layers
type App struct {
	*Theme
	activePage         string
	bodyBackground     string
	bodyColor          string
	buttonBar          []l.Widget
	hideSideBar        bool
	hideTitleBar       bool
	layers             []l.Widget
	pages              map[string]l.Widget
	root               *Stack
	sideBar            []l.Widget
	sideBarSize        unit.Value
	sideBarColor       string
	sideBarBackground  string
	logo               []byte
	logoClickable      *Clickable
	title              string
	titleBarBackground string
	titleBarColor      string
	titleFont          string
	menuClickable      *Clickable
	menuButton         *IconButton
	menuIcon           []byte
	menuColor          string
	menuBackground     string
	MenuOpen           bool
	responsive         *Responsive
	Size               *int
}

func (th *Theme) App(size int) *App {
	mc := th.Clickable()
	return &App{
		Theme:              th,
		activePage:         "main",
		bodyBackground:     "PanelBg",
		bodyColor:          "PanelText",
		buttonBar:          nil,
		hideSideBar:        false,
		hideTitleBar:       false,
		layers:             nil,
		pages:              make(map[string]l.Widget),
		root:               th.Stack(),
		sideBar:            nil,
		sideBarSize:        th.TextSize.Scale(20),
		sideBarColor:       "DocText",
		sideBarBackground:  "DocBg",
		logo:               ico.ParallelCoin,
		logoClickable:      th.Clickable(),
		title:              "parallelcoin",
		titleBarBackground: "Primary",
		titleBarColor:      "DocBg",
		titleFont:          "plan9",
		menuIcon:           icons.NavigationMenu,
		menuClickable:      mc,
		menuButton:         th.IconButton(mc),
		menuColor:          "Light",
		MenuOpen:           false,
		Size:               &size,
	}
}

// Fn renders the app widget
func (a *App) Fn() func(gtx l.Context) l.Dimensions {
	return func(gtx l.Context) l.Dimensions {
		x := gtx.Constraints.Max.X
		a.Size = &x
		return a.Flex().Rigid(
			a.VFlex().Rigid(
				a.Flex().Flexed(1,
					a.Fill(a.titleBarBackground,
						a.Flex().
							Rigid(
								a.Responsive(*a.Size,
									Widgets{
										{Widget: a.MenuButton},
										{Size: 800, Widget: a.MenuButtonAction}}).
									Fn,
							).
							Rigid(a.LogoAndTitle).
							Flexed(1,
								EmptyMinWidth(),
							).
							// Rigid(
							// 	a.DimensionCaption,
							// ).
							Rigid(
								a.RenderButtonBar,
							).
							Fn,
					).Fn,
				).Fn,
			).
				Flexed(1, a.MainFrame()).Fn,
		).Fn(gtx)
	}
}

func (a *App) RenderButtonBar(gtx l.Context) l.Dimensions {
	out := a.Flex()
	for i := range a.buttonBar {
		out.Rigid(a.buttonBar[i])
	}
	dims := out.Fn(gtx)
	gtx.Constraints.Min = dims.Size
	gtx.Constraints.Max = dims.Size
	return dims
}

func (a *App) MainFrame() func(gtx l.Context) l.Dimensions {
	return func(gtx l.Context) l.Dimensions {
		return a.Flex().
			Rigid(
				a.Flex().
					Rigid(
						a.Fill(a.sideBarBackground,
							a.Responsive(*a.Size, Widgets{
								{
									Widget: func(gtx l.Context) l.Dimensions {
										return If(a.MenuOpen,
											a.Fill(a.sideBarBackground,
												a.renderSideBar(),
											).Fn,
											EmptySpace(0, 0),
										)(gtx)
									},
								},
								{Size: 800,
									Widget: a.Fill(a.sideBarBackground,
										a.renderSideBar(),
									).Fn,
								},
							},
							).Fn,
						).Fn,
					).Fn,
			).
			Flexed(1,
				a.RenderPage,
			).
			Fn(gtx)
	}
}

func (a *App) MenuButton(gtx l.Context) l.Dimensions {
	bg := a.titleBarBackground
	color := a.menuColor
	if a.MenuOpen {
		color = "DocText"
	}
	return a.Flex().Rigid(
		a.Inset(0.25,
			a.ButtonLayout(a.menuClickable).
				CornerRadius(0).
				Embed(
					a.Inset(0.25,
						a.Icon().
							Scale(Scales["H5"]).
							Color(color).
							Src(icons.NavigationMenu).
							Fn,
					).Fn,
				).
				Background(bg).
				SetClick(
					func() {
						a.MenuOpen = !a.MenuOpen
					}).
				Fn,
		).Fn,
	).Fn(gtx)
}

func (a *App) MenuButtonAction(gtx l.Context) l.Dimensions {
	a.MenuOpen = false
	return l.Dimensions{}
}

func (a *App) LogoAndTitle(gtx l.Context) l.Dimensions {
	return a.Flex().
		Rigid(
			a.Responsive(*a.Size, Widgets{
				{
					Widget: EmptySpace(0, 0),
				},
				{Size: 800,
					Widget: a.Inset(0.125,
						a.Inset(0.125,
							a.IconButton(
								a.logoClickable.SetClick(
									func() {
										Debug("clicked logo")
										a.Dark = !a.Dark
										a.Theme.Colors.SetTheme(a.Dark)
									}),
							).
								Icon(
									a.Icon().
										Scale(Scales["H5"]).
										Color("Light").
										Src(a.logo)).
								Background("Dark").Color("Light").
								Inset(0.25).
								Fn,
						).Fn,
					).Fn,
				},
			},
			).Fn,
		).
		Rigid(
			a.Responsive(*a.Size, Widgets{
				{Size: 800,
					Widget:
					a.Inset(0.5,
						a.H5(a.title).Color("Light").Fn,
					).Fn,
				},
				{
					Widget:
					a.ButtonLayout(a.logoClickable).Embed(
						a.Inset(0.5,
							a.H5(a.title).Color("Light").Fn,
						).Fn,
					).Background("Transparent").Fn,
				},
			}).Fn,
		).Fn(gtx)
}

func (a *App) RenderPage(gtx l.Context) l.Dimensions {
	return a.Fill(a.bodyBackground,
		func(gtx l.Context) l.Dimensions {
			if page, ok := a.pages[a.activePage]; !ok {
				return a.Flex().
					Flexed(1,
						a.Inset(0.5,
							a.VFlex().SpaceEvenly().
								Rigid(
									a.H1("404").
										Alignment(text.Middle).
										Fn,
								).
								Rigid(
									a.Body1("page not found").
										Alignment(text.Middle).
										Fn,
								).
								Fn,
						).Fn,
					).Fn(gtx)
			} else {
				return page(gtx)
			}
		},
	).Fn(gtx)
}

func (a *App) DimensionCaption(gtx l.Context) l.Dimensions {
	return a.Caption(fmt.Sprintf("%dx%d", gtx.Constraints.Max.X, gtx.Constraints.Max.Y)).Fn(gtx)
}

func (a *App) renderSideBar() l.Widget {
	return func(gtx l.Context) l.Dimensions {
		gtx.Constraints.Max.X = 200 // a.scrollBarSize
		// gtx.Constraints.Min.X = a.scrollBarSize
		out := a.VFlex()
		for i := range a.sideBar {
			out.Rigid(a.sideBar[i])
		}
		// out.Rigid(EmptySpace(int(a.sideBarSize.V), 0))
		return out.Fn(gtx)
	}
}

func (a *App) ActivePage(activePage string) *App {
	a.activePage = activePage
	return a
}
func (a *App) ActivePageGet() string {
	return a.activePage
}

func (a *App) BodyBackground(bodyBackground string) *App {
	a.bodyBackground = bodyBackground
	return a
}
func (a *App) BodyBackgroundGet() string {
	return a.bodyBackground
}

func (a *App) BodyColor(bodyColor string) *App {
	a.bodyColor = bodyColor
	return a
}
func (a *App) BodyColorGet() string {
	return a.bodyColor
}

func (a *App) ButtonBar(bar []l.Widget) *App {
	a.buttonBar = bar
	return a
}
func (a *App) ButtonBarGet() (bar []l.Widget) {
	return a.buttonBar
}

func (a *App) HideSideBar(hideSideBar bool) *App {
	a.hideSideBar = hideSideBar
	return a
}
func (a *App) HideSideBarGet() bool {
	return a.hideSideBar
}

func (a *App) HideTitleBar(hideTitleBar bool) *App {
	a.hideTitleBar = hideTitleBar
	return a
}
func (a *App) HideTitleBarGet() bool {
	return a.hideTitleBar
}

func (a *App) Layers(widgets []l.Widget) *App {
	a.layers = widgets
	return a
}
func (a *App) LayersGet() []l.Widget {
	return a.layers
}

func (a *App) MenuBackground(menuBackground string) *App {
	a.menuBackground = menuBackground
	return a
}
func (a *App) MenuBackgroundGet() string {
	return a.menuBackground
}

func (a *App) MenuColor(menuColor string) *App {
	a.menuColor = menuColor
	return a
}
func (a *App) MenuColorGet() string {
	return a.menuColor
}

func (a *App) MenuIcon(menuIcon []byte) *App {
	a.menuIcon = menuIcon
	return a
}
func (a *App) MenuIconGet() []byte {
	return a.menuIcon
}

func (a *App) Pages(widgets map[string]l.Widget) *App {
	a.pages = widgets
	return a
}
func (a *App) PagesGet() map[string]l.Widget {
	return a.pages
}

func (a *App) Root(root *Stack) *App {
	a.root = root
	return a
}
func (a *App) RootGet() *Stack {
	return a.root
}

func (a *App) SideBar(widgets []l.Widget) *App {
	a.sideBar = widgets
	return a
}
func (a *App) SideBarGet() []l.Widget {
	return a.sideBar
}

func (a *App) SideBarBackground(sideBarBackground string) *App {
	a.sideBarBackground = sideBarBackground
	return a
}
func (a *App) SideBarBackgroundGet() string {
	return a.sideBarBackground
}

func (a *App) SideBarColor(sideBarColor string) *App {
	a.sideBarColor = sideBarColor
	return a
}
func (a *App) SideBarColorGet() string {
	return a.sideBarColor
}

func (a *App) Title(title string) *App {
	a.title = title
	return a
}
func (a *App) TitleGet() string {
	return a.title
}

func (a *App) TitleBarBackground(TitleBarBackground string) *App {
	a.bodyBackground = TitleBarBackground
	return a
}
func (a *App) TitleBarBackgroundGet() string {
	return a.titleBarBackground
}

func (a *App) TitleBarColor(titleBarColor string) *App {
	a.titleBarColor = titleBarColor
	return a
}
func (a *App) TitleBarColorGet() string {
	return a.titleBarColor
}

func (a *App) TitleFont(font string) *App {
	a.titleFont = font
	return a
}
func (a *App) TitleFontGet() string {
	return a.titleFont
}
