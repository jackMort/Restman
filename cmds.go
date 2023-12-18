package main

import (
	"errors"
	"fmt"
	neturl "net/url"
	"os"
	"restman/components/collections"
	"restman/components/footer"
	"restman/components/results"
	"restman/components/tabs"
	"restman/components/url"

	tea "github.com/charmbracelet/bubbletea"
	zone "github.com/lrstanley/bubblezone"
	"github.com/spf13/cobra"

	boxer "github.com/treilik/bubbleboxer"
)

var (
	version = "dev"
)

var rootCmd = &cobra.Command{
	Use:   "restman [http://example.com/api/v1]",
	Short: "A CLI tool for RESTful API",
	Long: `
┏┓┏┓┏╋┏┳┓┏┓┏┓
┛ ┗ ┛┗┛┗┗┗┻┛┗

Restman is a CLI tool for RESTful API.`,
	Version: version,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("accepts only one optional URL")
		}

		if len(args) == 1 {
			u, err := neturl.Parse(args[0])
			if err != nil || u.Scheme == "" || u.Host == "" {
				return errors.New("invalid URL provided")
			}
		}
		return nil
	},

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) >= 1 {
			fmt.Printf("You provided an optional argument: %s\n", args[0])
		} else {
			fmt.Println("You did not provide an optional argument.")
		}

		zone.NewGlobal()

		// layout-tree defintion
		m := Model{tui: boxer.Boxer{}, focused: "url"}

		url := url.New()
		middle := results.New()
		footerBox := footer.New()
		colBox := collections.New()
		tabs := tabs.New()

		centerNode := boxer.CreateNoBorderNode()
		centerNode.VerticalStacked = true
		centerNode.SizeFunc = func(node boxer.Node, widthOrHeight int) []int {
			return []int{
				2,
				3,
				widthOrHeight - 5,
			}
		}
		centerNode.Children = []boxer.Node{
			stripErr(m.tui.CreateLeaf("tabs", tabs)),
			stripErr(m.tui.CreateLeaf("url", url)),
			stripErr(m.tui.CreateLeaf("middle", middle)),
		}

		// middle Node
		middleNode := boxer.CreateNoBorderNode()
		middleNode.SizeFunc = func(node boxer.Node, widthOrHeight int) []int {
			gap := 30
			if m.tui.ModelMap["collections"].(collections.Collections).IsMinified() {
				gap = 6
			}
			return []int{
				gap,
				widthOrHeight - gap,
			}
		}
		middleNode.Children = []boxer.Node{
			stripErr(m.tui.CreateLeaf("collections", colBox)),
			centerNode,
		}

		rootNode := boxer.CreateNoBorderNode()
		rootNode.VerticalStacked = true
		rootNode.SizeFunc = func(node boxer.Node, widthOrHeight int) []int {
			return []int{
				widthOrHeight - 1,
				1,
			}
		}
		rootNode.Children = []boxer.Node{
			middleNode,
			stripErr(m.tui.CreateLeaf("footer", footerBox)),
		}

		m.tui.LayoutTree = rootNode

		if f, err := tea.LogToFile("debug.log", "debug"); err != nil {
			fmt.Println("Couldn't open a file for logging:", err)
			os.Exit(1)
		} else {
			defer f.Close()
		}

		p := tea.NewProgram(
			m,
			tea.WithAltScreen(),       // use the full size of the terminal in its "alternate screen buffer"
			tea.WithMouseCellMotion(), // turn on mouse support so we can track the mouse wheel
		)

		if _, err := p.Run(); err != nil {
			fmt.Println("could not run program:", err)
			os.Exit(1)
		}
	},
}
