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
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
	md "github.com/JohannesKaufmann/html-to-markdown"
)

const useHighPerformanceRenderer = false

// Establishing styles
var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.BorderStyle(b)
	}()
)


func main() {
    m := initialModel()
	tm, _ := tea.NewProgram(&m).Run()
	mm := tm.(model)
    fmt.Printf(mm.markdownString)   
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


// Most of the following boilerplate (and some of Main()) 
// is copied word-for-word from the Charm Bracelet examples here:
// https://github.com/charmbracelet/bubbletea/blob/master/examples

type model struct {
	filepicker        filepicker.Model
	viewport          viewport.Model
	selectedFile      string
	markdownString    string
	quitting          bool
	err               error
	content           string
	ready             bool
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
	var (
		cmd  tea.Cmd
		cmds tea.Cmd
	)
	
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "s":
			m.markdownString = doConvert(m.selectedFile)
			m.viewport.SetContent(m.markdownString)
			m.ready = true 
		case "enter":
			if m.ready {
				return m, tea.Quit
			}
		}
		
	case clearErrorMsg:
		m.err = nil

	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if m.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.HighPerformanceRendering = useHighPerformanceRenderer

			// This is only necessary for high performance rendering, which in
			// most cases you won't need.
			//
			// Render the viewport one line below the header.
			m.viewport.YPosition = headerHeight + 1
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}

		if useHighPerformanceRenderer {
			// Render (or re-render) the whole viewport. Necessary both to
			// initialize the viewport and when the window is resized.
			//
			// This is needed for high-performance rendering only.
			cmds = tea.Batch(cmds, viewport.Sync(m.viewport))
		}
	}

	m.filepicker, cmd = m.filepicker.Update(msg)
	cmds = tea.Batch(cmds, cmd)
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = tea.Batch(cmds, cmd)

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
		cmds = tea.Batch(cmds, tea.Batch(cmd, clearErrorAfter(2*time.Second)))
	}

	return m, cmds
}

func (m model) View() string {
	if m.quitting {
		return ""
	}
	if !m.ready {
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
	} else {
		return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
	}
}

// --- End copied Text --


// Viewport functions

func (m model) headerView() string {
	title := titleStyle.Render("Markdown Preview")
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m model) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

