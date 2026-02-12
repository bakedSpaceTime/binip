package libip

import (
	"fmt"
	"os"
	"runtime"

	"github.com/bakedSpaceTime/binip/libip/app"
	"github.com/bakedSpaceTime/binip/libip/config"
	"github.com/bakedSpaceTime/binip/libip/db"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/davecgh/go-spew/spew"
)

func Info(c *config.Config) error {
	fmt.Println("System info")
	fmt.Println("\tNum Logical CPU:", runtime.NumCPU())
	fmt.Println("\tOperating System:", runtime.GOOS)
	fmt.Println("\tArchitecture:", runtime.GOARCH)
	fmt.Println("\tDebug:", c.Debug)
	fmt.Println("\tDebug File:", c.DebugFile)

	for i := 0; i < 256; i++ {
		style := lipgloss.NewStyle().
			Background(lipgloss.Color(fmt.Sprintf("%d", i))).
			Width(4)
		fmt.Print(style.Render(fmt.Sprintf("%3d", i)))
		if (i+1)%8 == 0 && i < 16 || (i+1-16)%6 == 0 && i >= 16 {
			fmt.Println()
		}
	}
	fmt.Println()
	fmt.Println(db.New(c).String())
	return nil
}

func App(c *config.Config) error {
	var dump *os.File
	var err error

	if c.Debug {
		var err error
		dump, err = os.OpenFile(c.DebugFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
		if err != nil {
			spew.Dump(err)
			os.Exit(1)
		}
		defer dump.Close()
	}
	c.DebugWriter = dump

	p := tea.NewProgram(app.New(c), tea.WithAltScreen())
	if _, err = p.Run(); err != nil {
		fmt.Println("could not start program:", err)
	}
	return err
}

func Test(c *config.Config) error {
	return nil
}

func Reset(c *config.Config) error {
	return db.New(c).Reset()
}
