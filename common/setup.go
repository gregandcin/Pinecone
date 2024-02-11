package common

import (
	"fmt"
	"os"
	"pinecone/cli"
	"pinecone/gui"
	"runtime"
)

func CheckDataFolder(dataFolder string) error {
	// Ensure data folder exists
	if _, err := os.Stat(dataFolder); os.IsNotExist(err) {
		fmt.Println("Data folder not found. Creating...")
		if mkDirErr := os.Mkdir(dataFolder, 0755); mkDirErr != nil {
			return fmt.Errorf("error creating data folder: %v", mkDirErr)
		}
	}
	return nil
}

func CheckDatabaseFile(cliOpts *cli.CLIOptions, guiOpts *gui.GUIOptions) error {
	if nil != cliOpts {
	// Check if JSON file exists
	if _, err := os.Stat(cliOpts.JSONFilePath); os.IsNotExist(err) {
		// Prompt for download if JSON file doesn't exist
		if cli.PromptForDownload(cliOpts.JSONUrl) {
			err := loadJSONData(cliOpts.JSONFilePath, "Xbox-Preservation-Project", "Pinecone", "data/id_database.json", &cliOpts.TitleList, true)
			if err != nil {
				return fmt.Errorf("error downloading data: %v ", err)
			}
		} else {
			return fmt.Errorf("download aborted by user")
		}
	} else if cliOpts.UpdateFlag {
		// Handle manual update
		err := loadJSONData(cliOpts.JSONFilePath, "Xbox-Preservation-Project", "Pinecone", cliOpts.JSONFilePath, &cliOpts.TitleList, true)
		if err != nil {
			return fmt.Errorf("error updating data: %v", err)
		}
	} else {
		// Load existing JSON data
		err := loadJSONData(cliOpts.JSONFilePath, "Xbox-Preservation-Project", "Pinecone", cliOpts.JSONFilePath, &cliOpts.TitleList, false)
		if err != nil {
			return fmt.Errorf("error loading data: %v", err)
		}
	}	
	}

	
	return nil
}

func CheckDumpFolder(dumpLocation string) error {
	if dumpLocation != "dump" {
		if _, err := os.Stat(dumpLocation); os.IsNotExist(err) {
			return fmt.Errorf("directory does not exist, exiting")
		}
	} else {
		if _, err := os.Stat(dumpLocation); os.IsNotExist(err) {
			fmt.Println("Default dump folder not found. Creating...")
			if mkDirErr := os.Mkdir(dumpLocation, 0755); mkDirErr != nil {
				return fmt.Errorf("error creating dump folder: %v", mkDirErr)
			}
			return fmt.Errorf("please place TDATA folder in the \"dump\" folder")
		}
	}
	return nil
}

func CheckParsingSettings(cliOpts *cli.CLIOptions, guiOpts *gui.GUIOptions) error {
	if nil != cliOpts {
	if cliOpts.TitleID != "" {
		// if the titleID flag is set, print stats for that title
		cliOpts.PrintStats(cliOpts.TitleID, false)
	} else if cliOpts.Summarize {
		// if the summarize flag is set, print stats for all titles
		cliOpts.PrintStats("", true)
	} else if cliOpts.FatXplorer {
		if runtime.GOOS == "windows" {
			if _, err := os.Stat(`X:\`); os.IsNotExist(err) {
				return fmt.Errorf(`FatXplorer's X: drive not found`)
			} else {
				fmt.Println("Checking for Content...")
				fmt.Println("====================================================================================================")
				err := CheckForContent("X:\\TDATA")
				if err != nil {
					return err
				}
			}
		} else {
			return fmt.Errorf("FatXplorer mode is only available on Windows")
		}
	} else {
		// If no flag is set, proceed normally
		// Check if TDATA folder exists
			if _, err := os.Stat(cliOpts.DumpLocation + "/TDATA"); os.IsNotExist(err) {
				return fmt.Errorf("TDATA folder not found. Please place TDATA folder in the dump folder")
			}
			fmt.Println("Checking for Content...")
			fmt.Println("====================================================================================================")
			err := CheckForContent("dump/TDATA")
			if err != nil {
				return err
			}		
	}	
	}

	

	return nil
}
