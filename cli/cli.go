package cli

import (
	"fmt"
	"log"
	"strings"

	"pinecone/common"

	"github.com/fatih/color"
)

const (
	headerWidth = 100
	separator   = ""
)

type CLIOptions struct {
	DataFolder   string
	JSONFilePath string
	JSONUrl      string
	DumpLocation string
	TitleList 	common.TitleList
	UpdateFlag bool
	Summarize bool
	TitleID string
	FatXplorer bool
}

func printHeader(title string) {
	title = strings.TrimSpace(title)
	if len(title) > headerWidth-6 { // -6 to account for spaces and equals signs
		title = title[:headerWidth-9] + "..."
	}
	formattedTitle := "== " + title + " =="
	padLen := (headerWidth - len(formattedTitle)) / 2
	color.New(color.FgCyan).Println(strings.Repeat("=", padLen) + formattedTitle + strings.Repeat("=", headerWidth-padLen-len(formattedTitle)))
}

func printInfo(colorCode color.Attribute, format string, args ...interface{}) {
	color.New(colorCode).Printf("    "+format, args...)
}

// Prints statistics for a specific title or for all titles if batch is true.
func (options *CLIOptions)PrintStats(titleID string, batch bool) {
	if batch {
		options.printTotalStats()
	} else {
		data, ok := options.TitleList.Titles[titleID]
		if !ok {
			fmt.Printf("No data found for title ID %s\n", titleID)
			return
		}
		fmt.Printf("Statistics for title ID %s:\n", titleID)
		options.printTitleStats(&data)
	}
}

// Prints statistics for TitleData.
func (options *CLIOptions)printTitleStats(data *common.TitleData) {
	fmt.Println("Title:", data.TitleName)
	fmt.Println("Total number of Content IDs:", len(data.ContentIDs))
	fmt.Println("Total number of Title Updates:", len(data.TitleUpdates))
	fmt.Println("Total number of Known Title Updates:", len(data.TitleUpdatesKnown))
	fmt.Println("Total number of Archived items:", len(data.Archived))
	fmt.Println()
}

func (options *CLIOptions)printTotalStats() {
	totalTitles := len(options.TitleList.Titles)
	totalContentIDs := 0
	totalTitleUpdates := 0
	totalKnownTitleUpdates := 0
	totalArchivedItems := 0

	// Set to store unique hashes of known title updates and archived items
	knownTitleUpdateHashes := make(map[string]struct{})
	archivedItemHashes := make(map[string]struct{})

	for _, data := range options.TitleList.Titles {
		totalContentIDs += len(data.ContentIDs)
		totalTitleUpdates += len(data.TitleUpdates)

		// Count unique known title updates
		for _, knownUpdate := range data.TitleUpdatesKnown {
			for hash := range knownUpdate {
				knownTitleUpdateHashes[hash] = struct{}{}
			}
		}

		// Count unique archived items
		for _, archivedItem := range data.Archived {
			for hash := range archivedItem {
				archivedItemHashes[hash] = struct{}{}
			}
		}
	}

	totalKnownTitleUpdates = len(knownTitleUpdateHashes)
	totalArchivedItems = len(archivedItemHashes)

	fmt.Println("Total Titles:", totalTitles)
	fmt.Println("Total Content IDs:", totalContentIDs)
	fmt.Println("Total Title Updates:", totalTitleUpdates)
	fmt.Println("Total Known Title Updates:", totalKnownTitleUpdates)
	fmt.Println("Total Archived Items:", totalArchivedItems)
}

func PromptForDownload(url string) bool {
	var response string
	fmt.Printf("The required JSON data is not found. It can be downloaded from %s\n", url)
	fmt.Print("Do you want to download it now? (yes/no): ")
	fmt.Scanln(&response)

	return strings.ToLower(response) == "yes"
}

func (options *CLIOptions) StartCLI() {
	err := common.CheckDataFolder(options.DataFolder)
	if err != nil {
		log.Fatalln(err)
	}

	err = common.CheckDatabaseFile(options, nil)
	if err != nil {
		log.Fatalln(err)
	}

	err = common.CheckDumpFolder(options.DumpLocation)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Pinecone v0.4.2b")
	fmt.Println("Please share output of this program with the Pinecone team if you find anything interesting!")

	err = common.CheckParsingSettings()
	if err != nil {
		log.Fatalln(err)
	}
}
