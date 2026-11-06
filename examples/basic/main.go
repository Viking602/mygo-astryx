// A small consumer of the public library, independent of Gallery internals.
package main

import (
	kit "github.com/Viking602/mygo-astryx"
	"github.com/Viking602/mygo-astryx/theme"
	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"log"
)

type app struct {
	name string
	sync bool
}

func (a *app) view(c *ui.Context) {
	kit.Scope(c, theme.System(c), func(k *kit.UI) {
		bar := c.TitleBar()
		ui.Row(c).Height(max(bar.Height, 28)).Padding(0, bar.Right, 0, bar.Left).FillWidth().DragWindow()
		k.Container().Padding(12, 24, 24, 24).Children(func() {
			k.Heading("欢迎使用 MyGo UI", 2)
			k.Field("名称", func() { k.Key("name").TextInput(&a.name, kit.InputOptions{Placeholder: "你的名称"}).FillWidth() })
			k.Key("sync").Switch(&a.sync, "启用同步")
			k.Button("保存", kit.ButtonOptions{Variant: kit.Primary}).AlignSelf(ui.Start).OnClick(func() { c.AddToast(ui.Toast{Title: "已保存", Description: a.name}) })
		})
		k.ToastViewport()
	}).Fill()
}
func main() {
	a := &app{}
	mygo.App.WhenReady(func() {
		mygo.NewWindow(mygo.WindowOptions{Title: "MyGo UI — Basic", Width: 560, Height: 360, TitleBarStyle: mygo.TitleBarHiddenInset, Content: ui.View(a.view)})
	})
	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}
