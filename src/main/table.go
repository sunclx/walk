package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"

	"fmt"
	"os"
	"sort"
	"time"
)

type Foo struct {
	index       int
	oldName     string
	fileType    string
	newName     string
	createdTime time.Time
	editedTime  time.Time
	checked     bool
}
type FooModel struct {
	walk.TableModelBase
	walk.SorterBase
	root       string
	sortColumn int
	sortOrder  walk.SortOrder
	image      *walk.Bitmap
	items      []*Foo
}

func NewFooModel(root string) *FooModel {
	m := &FooModel{root: root}
	m.ResetRows()
	return m
}
func (m *FooModel) RowCount() int {
	return len(m.items)
}
func (m *FooModel) Value(row, col int) interface{} {
	item := m.items[row]
	switch col {
	case 0:
		return item.index
	case 1:
		return item.oldName
	case 2:
		return item.fileType
	case 3:
		return item.newName
	case 4:
		return item.createdTime
	case 5:
		return item.editedTime
	}
	return nil
}

func (m *FooModel) Checked(row int) bool {
	return m.items[row].checked
}
func (m *FooModel) SetChecked(row int, checked bool) error {
	m.items[row].checked = checked
	return nil
}
func (m *FooModel) Sort(col int, order walk.SortOrder) error {
	m.sortColumn, m.sortOrder = col, order
	sort.Sort(m)
	return m.SorterBase.Sort(col, order)
}
func (m *FooModel) Len() int {
	return len(m.items)
}
func (m *FooModel) Less(i, j int) bool {
	a, b := m.items[i], m.items[j]
	c := func(ls bool) bool {
		if m.sortOrder == walk.SortAscending {
			return ls
		}
		return !ls
	}
	switch m.sortColumn {
	case 0:
		return c(a.index < b.index)
	case 1:
		return c(a.oldName < b.oldName)
	case 2:
		return c(a.fileType < b.fileType)
	case 3:
		return c(a.oldName < b.oldName)
	case 4:
		return c(a.createdTime.Before(b.createdTime))
	case 5:
		return c(a.editedTime.Before(b.editedTime))
	}
	panic("wrong")
}
func (m *FooModel) Swap(i, j int) {
	m.items[i], m.items[j] = m.items[j], m.items[i]
}
func (m *FooModel) ResetRows() {
	f, _ := os.Open(m.root)
	fis, _ := f.Readdir(0)
	m.items = make([]*Foo, len(fis))
	for i, fi := range fis {
		m.items[i] = &Foo{
			index:       i,
			oldName:     fi.Name(),
			createdTime: fi.ModTime(),
			editedTime:  fi.ModTime(),
		}
		m.items[i].fileType = "文件"
		if fi.IsDir() {
			m.items[i].fileType = "文件夹"
		}

	}
	m.PublishRowsReset()
	m.Sort(m.sortColumn, m.sortOrder)
}

type TableViewColumnFormat struct {
	Title      string
	Format     string
	Width      int
	Alignment  Alignment1D
	DataMember string
	Hidden     bool
	Precision  int
}

func main() {
	model := NewFooModel("img")
	foo := &TableViewColumnFormat{Title: "创建时间", Format: "2006-01-02", Width: 150}
	var tv *walk.TableView
	MainWindow{
		Title:  "TableView",
		Size:   Size{800, 600},
		Layout: VBox{},
		DataBinder: DataBinder{
			DataSource: foo,
			AutoSubmit: true,
			OnSubmitted: func() {
				fmt.Println(foo)
			},
		},
		Children: []Widget{
			Composite{
				Layout: HBox{},
				Children: []Widget{
					RadioButtonGroup{
						DataMember: "Alignment",
						Buttons: []RadioButton{
							RadioButton{
								Name:  "leftAlign",
								Text:  "左对齐",
								Value: AlignFar,
								OnClicked: func() {
									tv.Columns().Add(walk.NewTableViewColumn())
									model.PublishRowsReset()
								},
							},
							RadioButton{
								Name:  "leftAlign",
								Text:  "居中",
								Value: AlignCenter,
								OnClicked: func() {
									tv.Columns().At(4).SetAlignment(walk.AlignCenter)
									model.PublishRowsReset()
								},
							},
							RadioButton{
								Name:  "leftAlign",
								Text:  "右对齐",
								Value: AlignNear,
								OnClicked: func() {
									tv.Columns().At(4).SetAlignment(walk.AlignFar)
									//tv.Columns().Clear()
									model.PublishRowsReset()
								},
							},
						},
					},
				},
			},

			TableView{
				AssignTo:              &tv,
				AlternatingRowBGColor: walk.RGB(183, 208, 65),
				CheckBoxes:            true,
				ColumnsOrderable:      true,
				Columns: []TableViewColumn{
					{Title: "#"},
					{Title: "原文件名"},
					{Title: "类型"},
					{Title: "新文件名"},
					{Title: "创建时间", Format: "2006-01-02 15:04:05", Width: 150, Alignment: foo.Alignment, DataMember: "", Hidden: false, Precision: 2},
					{Title: "修改时间", Format: "2006-01-02 15:04:05", Width: 150, DataMember: "", Alignment: foo.Alignment},
				},
				Model: model,
			},
		},
	}.Run()
}
