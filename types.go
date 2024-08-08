package panels

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// A part of the data tree. Each tree can contain as many children as desired.
type TreeNode struct {
	name     string
	children []*TreeNode
}

// Add a child of a node to the data tree
func (t *TreeNode) AddChild(text string) *TreeNode {
	child := &TreeNode{name: text}
	t.children = append(t.children, child)

	return child
}

// Add a sub-tree to the node of the data tree
func (t *TreeNode) InheritTree(subTree *DataTree) *TreeNode {
	child := &TreeNode{name: subTree.root.name}
	child.children = subTree.root.children

	t.children = append(t.children, child)

	return child
}

func (t *TreeNode) generateNode() *tview.TreeNode {
	node := tview.NewTreeNode(t.name)

	for _, child := range t.children {
		node.AddChild(child.generateNode())
	}

	node.SetExpanded(false)

	return node
}

// The overall data tree. Has one (implicit) root child and many sub children
type DataTree struct {
	root TreeNode
}

// Add a child of the root node
func (t *DataTree) AddChild(text string) *TreeNode {
	return t.root.AddChild(text)
}

// Add a sub-tree to the root node
func (t *DataTree) InheritTree(subTree *DataTree) *TreeNode {
	return t.root.InheritTree(subTree)
}

func (t *DataTree) generateTree() *tview.TreeView {
	root := t.root.generateNode().SetColor(tcell.ColorRed)

	root.SetExpanded(true)

	tree := tview.NewTreeView().SetRoot(root).SetCurrentNode(root)

	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		level := node.GetLevel()

		if level == 0 {
			return
		}

		expanded := node.IsExpanded()
		children := node.GetChildren()

		if len(children) == 0 {
			return
		}

		if expanded {
			node.SetColor(tcell.ColorWhite)
			node.SetExpanded(false)
		} else {
			node.SetColor(tcell.ColorGreen)
			node.SetExpanded(true)
		}

	})

	return tree
}

// Make a new data tree. The name is the name of the root node
func NewDataTree(name string) *DataTree {
	rootNode := TreeNode{name: name}

	tree := DataTree{root: rootNode}
	return &tree
}

type listItem struct {
	name     string
	pageName string
	details  string
	char     rune
	selected func()
}

// The Gui container and handler
type GuiContext interface {
	// Runs the Gui when it has been createed
	Run() error
	// Stops the Gui (causes Run to return)
	Stop()
}

type context struct {
	app               *tview.Application
	pages             *tview.Pages
	screens           *[]listItem
	components        *[]tview.Primitive
	currentComponents []*tview.Primitive
	currentPage       string
	back              func()
	frame             *tview.Frame
	dataTrees         []*DataTree
}

func (c *context) switchPage(name string) {
	pageIdx := -1

	for i, sc := range *c.screens {
		if name == sc.pageName {
			pageIdx = i
			break
		}
	}

	if pageIdx == -1 {
		panic("Could not find page")
	}

	if len(*c.components) == 0 {
		panic("Test results not submitted")
	}

	component := &(*c.components)[pageIdx]

	if component == nil {
		panic("component is nil ptr")
	}

	c.currentComponents = []*tview.Primitive{component}

	c.currentPage = name

	c.paintView()

	c.pages.SwitchToPage(name)
}

func (c *context) addView(idx int) {
	activeComps := []*tview.Primitive{}

	if idx == -1 {
		c.refreshPage()
		return
	} else if idx < len(*c.components) {
		activeComps = append(activeComps, c.currentComponents...)
		activeComps = append(activeComps, &(*c.components)[idx])
	} else {
		for i := range *c.components {
			c := &(*c.components)[i]
			activeComps = append(activeComps, c)
		}
	}

	c.currentComponents = activeComps

	c.paintView()
	c.refreshPage()
}

func (c *context) popView() {
	n := len(c.currentComponents)

	if n > 0 {
		activeComps := (c.currentComponents)[:n-1]
		c.currentComponents = activeComps
	}

	c.paintView()
	c.refreshPage()
}

func (c *context) refreshPage() {

	if len(c.currentComponents) == 0 {
		c.pages.SwitchToPage("Menu")
	} else if c.currentPage == "" {
		c.pages.SwitchToPage((*c.screens)[0].pageName)
		c.currentPage = "Laps"
	} else {
		c.pages.SwitchToPage(c.currentPage)
	}
}

func (c *context) paintView() {
	name := c.currentPage
	c.pages.RemovePage(name)
	page := c.makePage()
	c.pages.AddPage(name, page, true, false)
}

func (c *context) makePage() *tview.Grid {
	pages := c.currentComponents

	if len(pages) > 4 {
		pages = pages[:4]
	}

	list := append(*c.screens, listItem{name: "Menu", details: "Go back to the menu", char: 'b', selected: c.back})

	menuInt := makeBaseList(list)
	menu := tview.NewFrame(menuInt)
	menu.SetBorder(true)

	numPages := len(pages)

	cols := []int{}
	rows := []int{}

	switch numPages {
	case 1:
		cols = []int{40, 0}
		rows = []int{0}
	case 2:
		cols = []int{40, 0, 0}
		rows = []int{0}
	case 3:
		cols = []int{40, 0}
		rows = []int{0, 0, 0}
	case 4:
		cols = []int{40, 0, 0}
		rows = []int{0, 0}
	}

	numCol := len(cols)
	numRow := len(rows)

	grid := tview.NewGrid().
		SetRows(rows...).
		SetColumns(cols...).
		SetBorders(false).
		AddItem(menu, 0, 0, numRow, 1, 0, 0, true)

	frames := []*tview.Frame{}

	pageCtr := 0

	for i := 0; i < numRow; i++ {
		for j := 0; j < numCol; j++ {
			if j == 0 {
				//menu col
				continue
			}
			framed := tview.NewFrame(*pages[pageCtr])
			frames = append(frames, framed)
			framed.SetBorderColor(tcell.ColorWhite).SetBorder(true)

			grid.AddItem(framed, i, j, 1, 1, 0, 0, false)
			pageCtr++
		}
	}
	menu.SetBorderColor(tcell.ColorRed)

	focusedIdx := -1

	capturefn := func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == '\t' {
			menu.SetBorderColor(tcell.ColorWhite)
			for i := 0; i < numPages; i++ {
				frames[i].SetBorderColor(tcell.ColorWhite)
			}

			if focusedIdx+1 >= numPages {
				focusedIdx = -1
				c.app.SetFocus(menu)
				menu.SetBorderColor(tcell.ColorRed)
			} else {
				focusedIdx++
				c.app.SetFocus(frames[focusedIdx])
				frames[focusedIdx].SetBorderColor(tcell.ColorRed)
			}
			return nil
		}

		return event
	}

	grid.SetInputCapture(capturefn)

	return grid
}

func (c *context) paintWidgets() {
	for idx, tree := range c.dataTrees {
		t := tree.generateTree()
		(*c.components)[idx] = t
	}
}

// Run the
func (c *context) Run() error {

	c.paintWidgets()

	if err := c.app.SetRoot(c.frame, true).SetFocus(c.pages).Run(); err != nil {
		return err
	}

	return nil
}

func (c *context) Stop() {
	c.app.Stop()
}

// A menu item
type MenuItem struct {
	// The name of the item
	Name string

	// Details attached to the item
	Details string

	// A shortcut key that will be active when the menu is in focus
	Shortcut rune
}

// Data that the gui needs
type GuiData struct {
	// Text that can be displayed above the widgets at all times
	TopFrameText string

	// Text that can be displayed below the widgets at all times
	BottomFrameText string

	// A list of menu items
	MenuItems []*MenuItem

	// A list of data views
	DataViews []*DataTree
}
