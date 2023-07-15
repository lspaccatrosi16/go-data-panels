package panels

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func MakeGui(data GuiData) *GuiContext {
	app := tview.NewApplication()
	pages := tview.NewPages()

	frame := tview.NewFrame(pages).
		SetBorders(1, 1, 1, 1, 0, 0).
		AddText(data.TopFrameText, true, tview.AlignCenter, tcell.ColorWhite).
		AddText(data.BottomFrameText, false, tview.AlignCenter, tcell.ColorWhite)

	frame.SetBorder(true)

	comps := make([]tview.Primitive, len(data.DataViews))
	ctx := GuiContext{
		app:               app,
		pages:             pages,
		components:        &comps,
		currentComponents: []*tview.Primitive{},
		back:              func() { pages.SwitchToPage("Menu") },
		frame:             frame,
		dataTrees:         data.DataViews,
	}

	screens := []listItem{}

	for idx, item := range data.MenuItems {
		pName := fmt.Sprintf("page__%d", idx)
		screens = append(screens, listItem{
			Name:     item.Name,
			PageName: pName,
			Details:  item.Details,
			Char:     item.Shortcut,
			Selected: func() { ctx.switchPage(pName) },
		})
	}

	ctx.screens = &screens

	list := makeList(screens)
	modal := makeModal(&ctx)

	pages.AddPage("Menu", list, true, true)
	pages.AddPage("Modal", modal, true, false)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlQ {
			ctx.Stop()
		} else if event.Key() == tcell.KeyCtrlC {
			ctx.popView()
			return nil
		} else if event.Key() == tcell.KeyCtrlN {
			ctx.pages.SwitchToPage("Modal")
			// split window
			return nil
		} else if event.Key() == tcell.KeyCtrlR {
			ctx.paintWidgets()
			ctx.paintView()
			ctx.refreshPage()
			return nil
		} else if event.Key() == tcell.KeyCtrlA {
			numComps := len(*ctx.components)
			ctx.addView(numComps)
			return nil
		}

		return event
	})

	return &ctx
}

func makeModal(ctx *GuiContext) *tview.Modal {
	modal := tview.NewModal()

	modal.SetText("Pick a new view to open")

	buttons := []string{}

	for _, item := range *ctx.screens {
		buttons = append(buttons, item.Name)
	}

	buttons = append(buttons, "All")

	modal.AddButtons(buttons)
	modal.SetDoneFunc(func(idx int, label string) {
		ctx.addView(idx)
	})

	return modal

}

func makeBaseList(items []listItem) *tview.Flex {
	list := tview.NewList()

	for _, item := range items {
		list.AddItem(item.Name, item.Details, item.Char, item.Selected)
	}

	shortcutText := []string{
		"Ctrl-A to open all fields",
		"Ctrl-C to close window",
		"Ctrl-N to split window",
		"Ctrl-Q to close gui",
		"Ctrl-R to refresh window",
	}

	joinedText := strings.Join(shortcutText, "\n")

	shortcuts := tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText(joinedText)

	flex := tview.NewFlex()

	flex.SetDirection(tview.FlexRow).
		AddItem(list, len(items)*2, 0, true).
		AddItem(nil, 0, 1, false).
		AddItem(shortcuts, 0, 1, false)

	return flex
}

func makeList(items []listItem) *tview.Flex {
	list := makeBaseList(items)

	blank := tview.NewBox()

	flex := tview.NewFlex().
		AddItem(blank, 0, 1, false).
		AddItem(list, 0, 1, true).
		AddItem(blank, 0, 1, false)

	return flex
}
