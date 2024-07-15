package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/skryvvara/gossht/internal/ssh"
)

var (
	table *tview.Table
)

func main() {
	StartTUI()
}

func StartTUI() {
	// Create a new application
	app := tview.NewApplication()

	// Create a flex container
	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	// Create a box for the title
	title := tview.NewTextView().
		SetText("SSH Config").
		SetTextAlign(tview.AlignCenter).
		SetTextColor(tview.Styles.PrimaryTextColor).
		SetDynamicColors(true) // Optional: enable dynamic colors

	// Create a new table
	table = tview.NewTable().
		SetBorders(true).
		SetSelectable(true, false).
		SetFixed(1, 1).
		SetEvaluateAllRows(true) // Evaluate all rows for width calculation

	loadSSHConfig()

	// Set selection handler for the table
	table.SetSelectable(true, false).
		SetSelectedFunc(func(row, column int) {
			// Example: Print selected host details
			//host := table.GetCell(row, 0).Text
			name := table.GetCell(row, 1).Text
			user := table.GetCell(row, 2).Text

			app.Stop()

			ssh.SSHConnect(name+":22", user)

			StartTUI()
		})

	// Add headers with styling
	headerCell := func(text string) *tview.TableCell {
		return tview.NewTableCell(text).SetAlign(tview.AlignCenter).SetSelectable(false).
			SetBackgroundColor(tcell.ColorDarkCyan).SetTextColor(tcell.ColorWhite).SetAttributes(tcell.AttrBold)
	}

	table.SetCell(0, 0, headerCell("Host"))
	table.SetCell(0, 1, headerCell("HostName"))
	table.SetCell(0, 2, headerCell("User"))

	// Add title and table to the flex container
	flex.AddItem(title, 1, 1, false).
		AddItem(table, 0, 1, true)

	// Set the root flex container
	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
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
	cellStyle := tview.NewTableCell(host).
		SetAlign(tview.AlignCenter).
		SetSelectable(true).
		SetTextColor(tcell.ColorWhite).
		SetBackgroundColor(tcell.ColorBlack).
		SetExpansion(1) // Expand to fill available width

	// Add the cell to the table
	table.SetCell(row, 0, cellStyle)
	table.SetCell(row, 1, tview.NewTableCell(name).SetAlign(tview.AlignCenter).SetSelectable(true).SetTextColor(tcell.ColorWhite).SetBackgroundColor(tcell.ColorBlack).SetExpansion(1))
	table.SetCell(row, 2, tview.NewTableCell(user).SetAlign(tview.AlignCenter).SetSelectable(true).SetTextColor(tcell.ColorWhite).SetBackgroundColor(tcell.ColorBlack).SetExpansion(1))
}
