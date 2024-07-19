package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/skryvvara/gossht/internal/clear"
	"github.com/skryvvara/gossht/internal/ssh"
)

var (
	flex    *tview.Flex
	table   *tview.Table
	Version string // This is set during build time

	AccentColor tcell.Color = tcell.NewHexColor(0x324191)
)

func main() {
	StartTUI()
}

func StartTUI() {
	// Create a new application
	app := tview.NewApplication()

	// Create a flex container
	flex = tview.NewFlex().SetDirection(tview.FlexRow)

	// Create a box for the title
	title := tview.NewTextView().
		SetText("SSH Config " + Version).
		SetTextAlign(tview.AlignCenter).
		SetTextColor(tview.Styles.PrimaryTextColor).
		SetDynamicColors(true) // Optional: enable dynamic colors

	// Create a new table
	table = tview.NewTable().
		SetSeparator('|').
		SetSelectable(true, false).
		SetFixed(1, 1).
		SetEvaluateAllRows(true)

	table.SetTitle("Connections").SetBorder(true)

	// Customize the selected cell style
	selectedStyle := tcell.StyleDefault.
		Background(tcell.ColorDarkGray).
		Foreground(tcell.ColorWhite).
		Attributes(tcell.AttrBold)

	table.SetSelectedStyle(selectedStyle)

	//table.SetBorderPadding(0, 0, 1, 0)

	// Stop the application if ESCAPE has been pressed
	table.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			app.Stop()
		}
	})

	// Register key events
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		k := event.Key()

		switch k {
		case tcell.KeyCtrlE: // Edit Entry
			loadForm(app, true)
		case tcell.KeyCtrlN: // New Entry
			loadForm(app, false)
		case tcell.KeyCtrlD: // Delete Entry
			app.Stop()
		case tcell.KeyCtrlU: // Duplicate Entry
			app.Stop()
		}

		return event
	})

	loadSSHConfig()

	// Set selection handler for the table
	table.SetSelectable(true, false).
		SetSelectedFunc(func(row, column int) {
			// Example: Print selected host details
			//host := table.GetCell(row, 0).Text
			name := table.GetCell(row, 1).Text
			user := table.GetCell(row, 2).Text

			app.Stop()

			clear.CallClear()
			ssh.SSHConnect(name+":22", user)
			clear.CallClear()

			StartTUI()
		})

	// Add headers with styling
	headerCell := func(text string) *tview.TableCell {
		return tview.NewTableCell(text).SetSelectable(false).
			SetBackgroundColor(AccentColor).SetTextColor(tcell.ColorWhite).SetAttributes(tcell.AttrBold)
	}

	table.SetCell(0, 0, headerCell("Host"))
	table.SetCell(0, 1, headerCell("HostName"))
	table.SetCell(0, 2, headerCell("User"))

	infoBox := tview.NewGrid()

	infoBox.SetBorder(true).SetTitle("Info")

	//TODO: Remove this later
	infoBox.AddItem(tview.NewTextView().SetText("<ESC>: Quit Application"), 0, 0, 1, 1, 1, 1, false)
	infoBox.AddItem(tview.NewTextView().SetText("<ENTER>: Connect to the selected entry"), 1, 0, 1, 1, 1, 1, false)
	infoBox.AddItem(tview.NewTextView().SetText("<CTRL+E>: Edit Entry (Not yet implemented)"), 2, 0, 1, 1, 1, 1, false)
	infoBox.AddItem(tview.NewTextView().SetText("<CTRL+N>: New Entry (Not yet implemented)"), 0, 1, 1, 1, 1, 1, false)
	infoBox.AddItem(tview.NewTextView().SetText("<CTRL+D>: Delete Entry (Not yet implemented)"), 1, 1, 1, 1, 1, 1, false)
	infoBox.AddItem(tview.NewTextView().SetText("<CTRL+U>: Duplicate Entry (Not yet implemented)"), 2, 1, 1, 1, 1, 1, false)

	// Add title and table to the flex container
	flex.AddItem(title, 1, 1, false).
		AddItem(infoBox, 5, 1, false).
		AddItem(table, 0, 8, true)

	// Set the root flex container
	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}

func loadForm(app *tview.Application, preload bool) {
	var name string
	var hostname string
	var user string
	var title string

	if preload {
		x, _ := table.GetSelection()

		name = table.GetCell(x, 0).Text
		hostname = table.GetCell(x, 1).Text
		user = table.GetCell(x, 2).Text
		title = fmt.Sprintf("Edit entry %s", name)
	} else {
		title = "New entry"
	}

	formFlex := tview.NewFlex()

	form := tview.NewForm().
		AddInputField("Name", name, 20, nil, nil).
		AddInputField("Hostname", hostname, 20, nil, nil).
		AddInputField("User", user, 20, nil, nil).
		AddTextArea("Notes", "", 40, 0, 0, nil).
		AddPasswordField("Password", "", 10, '*', nil).
		AddInputField("SSH-Key path", "", 20, nil, nil).
		AddButton("Save", nil).
		AddButton("Quit", func() {
			app.SetRoot(flex, true)
		})

	form.SetLabelColor(tcell.ColorWhite).
		SetFieldBackgroundColor(AccentColor).
		SetFieldTextColor(tcell.ColorWhite).
		SetButtonBackgroundColor(AccentColor).
		SetButtonTextColor(tcell.ColorWhite)

	form.SetTitle(title).SetBorder(true)

	formFlex.AddItem(form, 0, 1, true)

	app.SetRoot(formFlex, true)
}

func loadSSHConfig() {
	sshPath := path.Join(os.Getenv("HOME"), ".ssh", "config")
	config, err := os.ReadFile(sshPath)
	if err != nil {
		fmt.Printf("Failed to read SSH config file: %v\n", err)
		return
	}

	// Parse SSH config file and populate the table
	lines := strings.Split(string(config), "\n")
	rowIndex := 1 // Start after the headers

	var currentHost, currentName, currentUser string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue // Skip empty lines and comments
		}

		line = strings.ToLower(line)

		// Detect Host entry
		if strings.HasPrefix(line, "host ") {
			// Add previous host entry to table
			if currentHost != "" {
				addHostEntryToTable(rowIndex, currentHost, currentName, currentUser)
				rowIndex++
			}
			currentHost = strings.TrimSpace(strings.TrimPrefix(line, "host "))
		} else if strings.HasPrefix(line, "hostname ") {
			currentName = strings.TrimSpace(strings.TrimPrefix(line, "hostname "))
		} else if strings.HasPrefix(line, "user ") {
			currentUser = strings.TrimSpace(strings.TrimPrefix(line, "user "))
		}
	}

	// Add the last host entry to table
	if currentHost != "" {
		addHostEntryToTable(rowIndex, currentHost, currentName, currentUser)
	}
}

func addHostEntryToTable(row int, host, name, user string) {
	// Normal cell style
	tableCell := func(content string) *tview.TableCell {
		return tview.NewTableCell(content).
			SetSelectable(true).
			SetTextColor(tcell.ColorWhite).
			SetBackgroundColor(tcell.ColorBlack).
			SetExpansion(1) // Expand to fill available width
	}

	// Add the cell to the table
	table.SetCell(row, 0, tableCell(host))
	table.SetCell(row, 1, tableCell(name))
	table.SetCell(row, 2, tableCell(user))
}
