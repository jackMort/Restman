package main

import (
	"errors"
	"fmt"
	neturl "net/url"
	"os"
	"restman/app"
	"restman/components/collections"
	"restman/components/config"
	"restman/components/footer"
	"restman/components/request"
	"restman/components/results"
	"restman/components/url"
	"restman/utils"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	zone "github.com/lrstanley/bubblezone"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	boxer "github.com/treilik/bubbleboxer"
)

var version = "dev"

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
		config.SetVersion(version)
		call := app.NewCall()
		if len(args) >= 1 {
			call.Url = args[0]
		}

		// TODO: what if both args and flags are provided?
		curl, _ := cmd.Flags().GetString("url")
		if curl != "" {
			call.Url = curl
		}

		method, _ := cmd.Flags().GetString("request")
		if method != "" {
			call.Method = method
		}

		data, _ := cmd.Flags().GetString("data")
		dataRaw, _ := cmd.Flags().GetString("data-raw")
		call.Data = data
		if call.Data == "" {
			call.Data = dataRaw
		}

		if call.Data != "" {
			call.DataType = "Text"
		}

		// make sure the method is POST if data is provided
		if call.Data != "" && call.Method == "GET" {
			call.Method = "POST"
		}

		headers, _ := cmd.Flags().GetStringArray("header")
		if headers != nil {
			// split headers into key-value pairs
			// to check authorization for bearer token
			processed_headers := []string{}
			for _, h := range headers {
				if h == "" {
					continue
				}
				pair := strings.Split(h, ":")
				if len(pair) == 2 {
					if strings.ToLower(pair[0]) == "authorization" && strings.Contains(pair[1], "Bearer") {
						call.Auth = &app.Auth{Type: "bearer_token", Token: strings.TrimSpace(strings.ReplaceAll(pair[1], "Bearer", ""))}
						continue
					}

					if strings.ToLower(pair[0]) == "content-type" && strings.Contains(pair[1], "application/json") {
						call.DataType = "JSON"
						call.Data = utils.FormatJSON(call.Data)
					}
				}
				processed_headers = append(processed_headers, h)
			}
			call.Headers = processed_headers
		}

		viper.SetConfigName("config")         // name of config file (without extension)
		viper.SetConfigType("json")           // REQUIRED if the config file does not have the extension in the name
		viper.AddConfigPath("/etc/restman/")  // path to look for the config file in
		viper.AddConfigPath("$HOME/.restman") // call multiple times to add many search paths
		err := viper.ReadInConfig()           // Find and read the config file
		if err != nil {                       // Handle errors reading the config file
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				// NOTE: ignore if config file is not found
			} else {
				panic(fmt.Errorf("fatal error config file: %w", err))
			}
		}

		var default_headers map[string]string = viper.GetStringMapString("default_headers")
		for k, v := range default_headers {
			call.Headers = append(call.Headers, fmt.Sprintf("%s: %s", k, v))
		}

		// ----
		zone.NewGlobal()

		// layout-tree defintion
		m := Model{tui: boxer.Boxer{}, focused: "url", initialCall: call}

		url := url.New()
		resultsBox := results.New()
		requestBox := request.New()
		footerBox := footer.New()
		colBox := collections.New()

		splitNode := boxer.CreateNoBorderNode()
		splitNode.SizeFunc = func(node boxer.Node, widthOrHeight int) []int {
			paramsSize := int(float64(widthOrHeight) * 0.4)
			return []int{
				paramsSize,
				widthOrHeight - paramsSize,
			}
		}
		splitNode.Children = []boxer.Node{
			stripErr(m.tui.CreateLeaf("request", requestBox)),
			stripErr(m.tui.CreateLeaf("results", resultsBox)),
		}

		centerNode := boxer.CreateNoBorderNode()
		centerNode.VerticalStacked = true
		centerNode.SizeFunc = func(node boxer.Node, widthOrHeight int) []int {
			return []int{
				3,
				widthOrHeight - 3,
			}
		}
		centerNode.Children = []boxer.Node{
			stripErr(m.tui.CreateLeaf("url", url)),
			splitNode,
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
