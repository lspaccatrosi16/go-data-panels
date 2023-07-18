# Go-Data-Panels

> Quick tree data panels for your go project with windows

## Why

Ever had some data produced by your go program that you wanted to display nicely? Fed up of going through console output trying to pick out what you need? This is the tool for you.

Organise related data in hierarchical trees so you don't have to do see everything at once, and organise these trees into different views so you only see what you want to see, with the rest a few keystrokes away.

## How do I get started?

### Create your menu items

```go
items := []*panels.MenuItem{
    {Name: "Foo", Details: "Foo Bar Baz", Shortcut: 'f'},
    {Name: "Other Menu Item", Details: "Some other string", 'o'},
    // etc...
}
```

In each item, there is a `Name`, `Details`, and a `Shortcut`. When the Menu is in focus, pressing the shortcut key has the same effect as hitting enter when the option is highlighted.

### Create your data trees

For each menu item you create, you must create a corresponding `DataTree` for it. It will be associated with the menu item visible in `items`.

```go
func makeFooTree() *panels.DataTree {
    tree := panels.NewDataTree("The foo tree")

    for i := 0; i < 10; i++ {
        tree.AddChild(fmt.Sprintf("%d", i))
    }

    return &tree
}

func makeOtherTree() *panels.DataTree {
    tree := panels.NewDataTree("Another tree")

    level1 := tree.AddChild("1st Level")
    level2 := level1.AddChild("2nd Level")

    return &tree
    // you can go down to as many levels as you want - not that you should
}

// ....

fooTree := makeFooTree()
otherTree := makeOtherTree()

trees := []*panels.DataTree{fooTree, otherTree}
```

### See your data

```go
gui := panels.MakeGui({
    MenuItems: items,
    DataViews: trees,
})

gui.Run()
```

You should now have a menu consisting of 2 items. Clicking on either of them shows the respective data tree. The component in focus is denoted by a red outline. Initially, the menu is always in focus. To cycle focus, press `Tab` and the focus switch from the menu to the data tree. You can naviage the tree with arrow keys, and if a given node has children, enter will hide / unhide them.

### That's good, what about seing multiple things

You can simply add another data tree to the current screen (up to a max of 4) by pressing `Ctrl-N` and selecting the name of the data view you want to see (yes you can add duplicates of the same widget, and they will share the same state).

To close a window, press `Ctrl-C` which will close the most recently opened one.

That's it, basic, functional data visualization with very little code needed.

### Full Example

```go
package main

import (
    "fmt"

    "github.com/lspaccatrosi16/go-data-panels"
)

func main() {
    items := []*panels.MenuItem{
        {Name: "Foo", Details: "Foo Bar Baz", Shortcut: 'f'},
        {Name: "Other Menu Item", Details: "Some other string", 'o'},
    }

    fooTree := makeFooTree()
    otherTree := makeOtherTree()

    trees := []*panels.DataTree{fooTree, otherTree}

    gui := panels.MakeGui({
        MenuItems: items,
        DataViews: trees,
    })

    gui.Run()
}

func makeFooTree() *panels.DataTree {
    tree := panels.NewDataTree("The foo tree")

    for i := 0; i < 10; i++ {
        tree.AddChild(fmt.Sprintf("%d", i))
    }

    return &tree
}

func makeOtherTree() *panels.DataTree {
    tree := panels.NewDataTree("Another tree")

    level1 := tree.AddChild("1st Level")
    level2 := level1.AddChild("2nd Level")

    return &tree
}

```

## Credits

This is based off of [tview](https://github.com/rivo/tview/), and uses it under the hood. For more functionally complex guis, it should be used instead. This project is essentialy a series of tview presets combined with basic window and state management in a nice wrapper. In the future, I may allow custom tview components to be used as well, but that would likely increase the complexity significantly.

## Licence

This project is availible under the Apache-2.0 licence. See [LICENCE](./LICENCE) for details.
