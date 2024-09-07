package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"image/color"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type GUIOptions struct {
	DataFolder   string
	JSONFilePath string
	JSONUrl      string
}

type Settings struct {
	UserName string `json:"username"`
	Discord  string `json:"discord"`
	Twitter  string `json:"twitter"`
	Reddit   string `json:"reddit"`
}

//go:embed pinecone.ui
var uiXML string

var outputContainer *gtk.Box

const (
	guiHeaderWidth = 50
)

func addHeader(title string) {
	title = strings.TrimSpace(title)
	if len(title) > guiHeaderWidth-6 {
		title = title[:guiHeaderWidth-4] + "..."
	}
	formattedTitle := "== " + title + " =="
	padLen := (guiHeaderWidth - len(formattedTitle)) / 2
	addText(color.RGBA{255, 255, 255, 255}, strings.Repeat("=", padLen)+formattedTitle+strings.Repeat("=", guiHeaderWidth-padLen-len(formattedTitle)))
}

func addText(textColor color.Color, format string, args ...interface{}) {
	label := gtk.NewLabel(fmt.Sprintf(format, args...))
	outputContainer.Append(label)
}

func loadSettings() (*Settings, error) {
	settingsPath := filepath.Join(dataPath, "pineconeSettings.json")
	settingsFile, err := os.Open(settingsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &Settings{}, nil
		}
		return nil, err
	}
	defer settingsFile.Close()

	settings := &Settings{}
	err = json.NewDecoder(settingsFile).Decode(settings)
	if err != nil {
		return nil, err
	}

	return settings, nil
}

func saveSettings(settings *Settings) error {
	settingsPath := filepath.Join(dataPath, "pineconeSettings.json")
	settingsFile, err := os.Create(settingsPath)
	if err != nil {
		return err
	}
	defer settingsFile.Close()

	encoder := json.NewEncoder(settingsFile)
	encoder.SetIndent("", "    ")

	err = encoder.Encode(settings)
	if err != nil {
		return err
	}

	return nil
}

func showSettingsDialog(settings *Settings) {
}

func setDumpFolder(window *gtk.Window) {
	dialog := gtk.NewFileDialog()
	dialog.SelectFolder(context.Background(), window, func(result gio.AsyncResulter) {
		folder, err := dialog.SelectFolderFinish(result)
		if err != nil {
			fmt.Println("ERROR: ", err.Error())
		}
		if folder != nil {
			tmpDumpPath := folder.PeekPath()
			if _, err := os.Stat(path.Join(tmpDumpPath, "TDATA")); os.IsNotExist(err) {
				dumpLocation = tmpDumpPath
				addText(color.RGBA{255, 255, 255, 255}, "Path set to: "+tmpDumpPath)
			} else {
				addText(color.RGBA{255, 0, 0, 255}, "Incorrect pathing. Please select a dump with TDATA folder.\n")
			}
		}
	})
}

func saveOutput(settings *Settings) {
	// // Get current time
	// t := time.Now()
	// // Format time to be used in filename
	// timestamp := t.Format("2006-01-02-15-04-05")
	// // Define the path to the output file
	// outputPath := filepath.Join(dataPath, "output", "output-"+timestamp+".txt")
	// // Create the 'output' directory if it doesn't exist
	// outputDir := filepath.Dir(outputPath)
	// if _, err := os.Stat(outputDir); os.IsNotExist(err) {
	// 	err = os.MkdirAll(outputDir, 0o755)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }
	// fileText := ""
	// // Add user info to top of file
	// if settings.UserName != "" {
	// 	fileText += fmt.Sprintf("Username: %s\n", settings.UserName)
	// }
	// if settings.Discord != "" {
	// 	fileText += fmt.Sprintf("Discord Username: @%s\n", settings.Discord)
	// }
	// if settings.Twitter != "" {
	// 	fileText += fmt.Sprintf("Twitter Username: @%s\n", settings.Twitter)
	// }
	// if settings.Reddit != "" {
	// 	fileText += fmt.Sprintf("Reddit Username: u/%s\n", settings.Reddit)
	// }
	// // Write output to file
	// for _, obj := range outputContainer {
	// 	if textObj, ok := obj.(*canvas.Text); ok {
	// 		// Append the text value to the string
	// 		fileText += textObj.Text + "\n"
	// 	}
	// }
	// err := os.WriteFile(outputPath, []byte(fileText), 0o644)
	// if err != nil {
	// 	panic(err)
	// }
	// // Debug output, show the path we're scanning
	// output := widget.NewLabel("Output saved to: " + outputPath + "\n")
	// outputContainer.Add(output)
}

func guiShowDownloadConfirmation(window *gtk.Window, filePath string, url string) {
	// message := fmt.Sprintf("The required JSON data is not found.\nIt can be downloaded from:\n%s\nDo you want to download it now?", url)
	//
	// dialog := &gtk.AlertDialog{}

	// confirmation := dialog.NewConfirm("Confirmation", message, func(confirmed bool) {
	// 	if confirmed {
	// 		// Action to perform if confirmed
	// 		err := loadJSONData(filePath, "Xbox-Preservation-Project", "Pinecone", dataPath+"/id_database.json", &titles, true)
	// 		if err != nil {
	// 			text := fmt.Sprintf("error downloading data: %v", err)
	// 			output := canvas.NewText(text, theme.ErrorColor())
	// 			outputContainer.Add(output)
	// 			return
	// 		}
	// 		guiScanDump()
	// 	} else {
	// 		// Action to perform if canceled
	// 		output := canvas.NewText("Download aborted by user", theme.ErrorColor())
	// 		outputContainer.Add(output)
	// 	}
	// }, window)
	//
	// // Show the confirmation dialog
	// confirmation.Show()
	//
	guiScanDump()
}

func guiScanDump() {
	err := checkDumpFolder(dumpLocation)
	if nil != err {
		fmt.Println("ERROR: ", err.Error())
		addText(nil, err.Error())
	}

	err = checkParsingSettings()
	if nil != err {
		fmt.Println("ERROR: ", err.Error())
		addText(nil, err.Error())
	}
}

func guiStartScan(options GUIOptions, win *gtk.Window) {
	// outputContainer.
	if dumpLocation == "" {
		addText(color.RGBA{255, 255, 255, 255}, "Please set a path first.")
	} else {
		addText(color.RGBA{255, 255, 255, 255}, "Checking for Content...")
		err := checkDatabaseFile(options.JSONFilePath, options.JSONUrl, updateFlag, win)
		if err != nil {
			fmt.Println("ERROR: ", err.Error())
			addText(color.RGBA{255, 0, 0, 255}, err.Error())
		}
	}
}

func startGUI(options GUIOptions, app *gtk.Application) {
	builder := gtk.NewBuilderFromString(uiXML)

	window := builder.GetObject("main_window").Cast().(*gtk.Window)

	notebook := builder.GetObject("main_notebook").Cast().(*gtk.Notebook)
	mainPage := builder.GetObject("main_box").Cast().(*gtk.Box)
	mainLable := builder.GetObject("main_tab").Cast().(*gtk.Label)

	settingsPage := builder.GetObject("update_database_box").Cast().(*gtk.Box)
	settingsLable := builder.GetObject("settings_tab").Cast().(*gtk.Label)

	notebook.AppendPage(mainPage, mainLable)
	notebook.AppendPage(settingsPage, settingsLable)

	// outputContainer = builder.GetObject("output_box").Cast().(*gtk.Box)
	// scroll := gtk.NewScrolledWindow()
	// scroll.SetChild(outputContainer)
	//
	// startScanButton := builder.GetObject("start_scan_button").Cast().(*gtk.Button)
	// startScanButton.Connect("clicked", func() {
	// 	guiStartScan(options, window)
	// })
	//
	// saveOutputButton := builder.GetObject("save_output_button").Cast().(*gtk.Button)
	// saveOutputButton.Connect("clicked", func() {
	// 	settings, err := loadSettings()
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		settings = &Settings{}
	// 	}
	//
	// 	saveOutput(settings)
	// })
	//
	// updateDatabaseButton := builder.GetObject("update_database_button").Cast().(*gtk.Button)
	// updateDatabaseButton.Connect("clicked", func() {
	// 	updateJSON := true
	// 	err := checkDatabaseFile(options.JSONFilePath, options.JSONUrl, updateJSON, nil)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// })
	//
	// saveSettingsButton := builder.GetObject("save_settings_button").Cast().(*gtk.Button)
	// saveSettingsButton.Connect("clicked", func() {
	// })

	window.SetVisible(true)
}
