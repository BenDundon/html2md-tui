package main

import (
	"fmt"
    "log"
    "os"
    "io"

    "strings"
	"time"
	"errors"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
	md "github.com/JohannesKaufmann/html-to-markdown"
)

func main() {
    m := initialModel()
	tm, _ := tea.NewProgram(&m).Run()
	mm := tm.(model)
    fileName := mm.selectedFile
    markdown := doConvert(fileName)
    fmt.Printf(markdown)   
}

func initialModel() model {

	fp := filepicker.New()
	fp.AllowedTypes = []string{".html"}
	fp.CurrentDirectory, _ = os.UserHomeDir()
	return model{
		filepicker: fp,
	}
}

func doConvert(fileName string) string {
	
	converter := md.NewConverter("", true, nil)
    html := readFile(fileName) 
    
	markdown, err := converter.ConvertString(html)

	if err != nil {
  		log.Fatal(err)
	}
	
	fmt.Println("md ->", markdown)
	return markdown
}

func readFile(fileName string) string {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}

    // Write the file text somewhere else so we can close it
    fileText, err := io.ReadAll(file)

    if err != nil {
    	log.Fatal(err)
    }

    // File *Must* be manually closed
    file.Close()
    
	return string(fileText)
}


// The following boilerplate (and some of Main()) 
// is copied word-for-word from the Charm Bracelet example here:
// https://github.com/charmbracelet/bubbletea/blob/master/examples/file-picker/main.go

type model struct {
	filepicker   filepicker.Model
	selectedFile string
	quitting     bool
	err          error
}

type clearErrorMsg struct{}

func clearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return clearErrorMsg{}
	})
}

func (m model) Init() tea.Cmd {
	return m.filepicker.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "s":
			return m, nil		
		}
		
	case clearErrorMsg:
		m.err = nil
	}

	var cmd tea.Cmd
	m.filepicker, cmd = m.filepicker.Update(msg)

	// Did the user select a file?
	if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
		// Get the path of the selected file.
		m.selectedFile = path
	}

	// Did the user select a disabled file?
	// This is only necessary to display an error to the user.
	if didSelect, path := m.filepicker.DidSelectDisabledFile(msg); didSelect {
		// Let's clear the selectedFile and display an error.
		m.err = errors.New(path + " is not valid.")
		m.selectedFile = ""
		return m, tea.Batch(cmd, clearErrorAfter(2*time.Second))
	}

	return m, cmd
}

func (m model) View() string {
	if m.quitting {
		return ""
	}
	var s strings.Builder
	s.WriteString("\n  ")
	if m.err != nil {
		s.WriteString(m.filepicker.Styles.DisabledFile.Render(m.err.Error()))
	} else if m.selectedFile == "" {
		s.WriteString("Pick a file:")
	} else {
    	s.WriteString("Selected file: " + m.filepicker.Styles.Selected.Render(m.selectedFile))
	}
	s.WriteString("\n\n" + m.filepicker.View() + "\n")
	return s.String()
}

// --- End copied Text --
