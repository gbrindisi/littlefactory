// Package tui provides terminal user interface components for littlefactory.
package tui

import (
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/yourusername/littlefactory/internal/config"
	"github.com/yourusername/littlefactory/internal/driver"
	"github.com/yourusername/littlefactory/internal/tasks"
)

const (
	// leftPanelWidth is the fixed width of the task list panel
	leftPanelWidth = 30
)

// FileChangedMsg is emitted when the progress file changes on disk.
type FileChangedMsg struct{}

// Model holds the state for the TUI.
type Model struct {
	// config holds the application configuration
	config *config.Config

	// projectRoot is the project root directory path
	projectRoot string

	// progressFilePath is the computed path to progress.md file
	progressFilePath string

	// tasks is the full list of tasks with current status
	tasks []tasks.Task

	// activeTaskID is the ID of the currently executing task
	activeTaskID string

	// outputPanel is the output panel component
	outputPanel *OutputPanel

	// progressContent holds the content from progress.md file
	progressContent string

	// autoFollow controls whether the viewport scrolls to bottom on new output
	autoFollow bool

	// eventChan receives events from the driver
	eventChan <-chan interface{}

	// width and height track terminal dimensions
	width  int
	height int

	// iteration tracks the current iteration number
	iteration int

	// maxIterations is the maximum number of iterations configured
	maxIterations int

	// runComplete indicates if the run has finished
	runComplete bool

	// finalStatus holds the final run status when complete
	finalStatus driver.RunStatus
}

// New creates a new TUI model.
// The eventChan parameter receives events from the driver running in a goroutine.
// The cfg parameter provides configuration including the state directory.
// The projectRoot parameter is used to construct the full path to progress.md.
func New(eventChan <-chan interface{}, cfg *config.Config, projectRoot string) *Model {
	return &Model{
		config:           cfg,
		projectRoot:      projectRoot,
		progressFilePath: filepath.Join(projectRoot, cfg.StateDir, "progress.md"),
		tasks:            []tasks.Task{},
		outputPanel:      NewOutputPanel(0, 0),
		autoFollow:       true, // Start with auto-follow enabled
		eventChan:        eventChan,
	}
}

// Init implements tea.Model.
// It returns a command that waits for events from the driver and starts watching the progress file.
func (m *Model) Init() tea.Cmd {
	// Load initial progress file content if it exists
	_ = m.loadProgressFile() // Ignore error, file may not exist yet

	return tea.Batch(
		waitForEvent(m.eventChan),
		watchProgressFile(m.progressFilePath),
	)
}

// Update implements tea.Model.
// It handles all message types from bubbletea and the driver.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle keyboard input
		switch msg.String() {
		case "q", "ctrl+c":
			// Quit the TUI
			return m, tea.Quit

		case "f":
			// Toggle auto-follow mode
			m.autoFollow = !m.autoFollow
			if m.autoFollow {
				// Scroll to bottom when enabling auto-follow
				m.outputPanel.GotoBottom()
			}

		case "up":
			// Scroll viewport up
			m.outputPanel.LineUp(1)
			m.autoFollow = false // Disable auto-follow when manually scrolling

		case "down":
			// Scroll viewport down
			m.outputPanel.LineDown(1)

		case "pgup":
			// Page up in viewport
			m.outputPanel.ViewUp()
			m.autoFollow = false

		case "pgdown":
			// Page down in viewport
			m.outputPanel.ViewDown()
		}

	case tea.WindowSizeMsg:
		// Handle terminal resize
		m.width = msg.Width
		m.height = msg.Height
		m.recalculateLayout()

	case driver.RunStartedMsg:
		// Run has started
		m.maxIterations = msg.MaxIterations
		cmds = append(cmds, waitForEvent(m.eventChan))

	case driver.IterationStartedMsg:
		// New iteration has started
		m.iteration = msg.Iteration
		m.activeTaskID = msg.TaskID

		cmds = append(cmds, waitForEvent(m.eventChan))

	case driver.OutputMsg:
		// Received output from agent - now handled via file watching
		cmds = append(cmds, waitForEvent(m.eventChan))

	case driver.IterationCompleteMsg:
		// Iteration has completed
		m.activeTaskID = ""
		cmds = append(cmds, waitForEvent(m.eventChan))

	case driver.TasksRefreshedMsg:
		// Task list has been updated
		m.tasks = msg.Tasks
		cmds = append(cmds, waitForEvent(m.eventChan))

	case driver.RunCompleteMsg:
		// Run has finished
		m.runComplete = true
		m.finalStatus = msg.Status
		// Don't wait for more events

	case FileChangedMsg:
		// Progress file has changed on disk, reload it
		_ = m.loadProgressFile() // Ignore error, file may not exist yet
		// Continue watching for file changes
		cmds = append(cmds, watchProgressFile(m.progressFilePath))
	}

	return m, tea.Batch(cmds...)
}

// View implements tea.Model.
// It renders the two-panel layout with task list and output.
func (m *Model) View() string {
	if m.width == 0 || m.height == 0 {
		// Not initialized yet
		return "Initializing..."
	}

	// Build left panel (task list)
	contentHeight := m.height - 1 // Reserve 1 line for status bar
	leftPanel := renderTasksPanel(m.tasks, m.activeTaskID, leftPanelWidth, contentHeight)

	// Build right panel (output viewport)
	var rightPanel string
	if len(m.progressContent) == 0 {
		// Show placeholder when progress file doesn't exist or is empty
		placeholder := "Waiting for progress..."
		if m.runComplete {
			placeholder = "Run complete."
		}
		rightPanel = lipgloss.NewStyle().
			Faint(true).
			Width(m.width - leftPanelWidth).
			Height(contentHeight).
			Padding(1, 2).
			Render(placeholder)
	} else {
		rightPanel = m.outputPanel.View()
	}

	// Combine panels horizontally
	mainContent := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftPanel,
		rightPanel,
	)

	// Build status bar
	statusBar := renderStatusBar(
		m.tasks,
		m.iteration,
		m.maxIterations,
		m.autoFollow,
		m.runComplete,
		m.finalStatus,
		m.width,
	)

	// Combine vertically
	return lipgloss.JoinVertical(
		lipgloss.Left,
		mainContent,
		statusBar,
	)
}

// recalculateLayout adjusts panel dimensions based on terminal size.
func (m *Model) recalculateLayout() {
	// Reserve 1 line for status bar
	contentHeight := m.height - 1

	// Right panel gets remaining width after left panel
	rightPanelWidth := m.width - leftPanelWidth
	if rightPanelWidth < 0 {
		rightPanelWidth = 0
	}

	// Update output panel size
	m.outputPanel.SetSize(rightPanelWidth, contentHeight)
}

// loadProgressFile reads the progress.md file and updates the viewport content.
// Returns an error if the file cannot be read.
func (m *Model) loadProgressFile() error {
	// Check if file exists first using os.Stat
	_, err := os.Stat(m.progressFilePath)
	if os.IsNotExist(err) {
		// File doesn't exist - set empty string to trigger placeholder
		m.progressContent = ""
		m.outputPanel.SetContent("")
		return err
	}

	content, err := os.ReadFile(m.progressFilePath)
	if err != nil {
		// File can't be read - use empty string
		m.progressContent = ""
		m.outputPanel.SetContent("")
		return err
	}

	// Store file content
	m.progressContent = string(content)
	m.outputPanel.SetContent(m.progressContent)

	// Scroll to bottom if auto-follow is enabled
	if m.autoFollow {
		m.outputPanel.GotoBottom()
	}

	return nil
}

// waitForEvent creates a command that waits for the next event from the driver.
func waitForEvent(eventChan <-chan interface{}) tea.Cmd {
	return func() tea.Msg {
		msg, ok := <-eventChan
		if !ok {
			// Channel closed, return nil to stop waiting
			return nil
		}
		return msg
	}
}

// watchProgressFile creates a command that watches the progress file for changes.
// It blocks until the file is modified, created, or deleted and then returns FileChangedMsg.
// If the file doesn't exist, it watches the directory for file creation.
func watchProgressFile(path string) tea.Cmd {
	return func() tea.Msg {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			// Watcher creation failed, return nil to stop watching
			return nil
		}
		defer watcher.Close()

		// Try to watch the file first
		err = watcher.Add(path)
		if err != nil {
			// File doesn't exist - watch the directory instead for file creation
			dir := filepath.Dir(path)
			err = watcher.Add(dir)
			if err != nil {
				// Can't watch directory either, return nil
				return nil
			}
		}

		// Block until we receive a relevant event
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return nil
				}
				// Watch for Write (file modified), Create (file created), or Remove (file deleted)
				if event.Name == path {
					if event.Op&fsnotify.Write == fsnotify.Write ||
						event.Op&fsnotify.Create == fsnotify.Create ||
						event.Op&fsnotify.Remove == fsnotify.Remove {
						return FileChangedMsg{}
					}
				}
			case <-watcher.Errors:
				// Error occurred, stop watching
				return nil
			}
		}
	}
}
