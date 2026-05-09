package main

import (
	"flag"
	"log"
	"os"
	"fmt"
)


func main() {
	var err error

	var flagOutputPath string
	var flagUpdateBView bool
	var flagWhois bool

	var records TRecords

	flag.StringVar(&flagOutputPath, "o", "", "Path to write the output to a CSV file")
	flag.BoolVar(&flagUpdateBView, "update-latest-bview", false, "Force download the latest BGP view")
	flag.BoolVar(&flagWhois, "whois", false, "Force download the latest BGP view")
	flag.Parse()

	if _, err = os.Stat(_CONST_LATEST_BVIEW_FILEPATH); os.IsNotExist(err) || flagUpdateBView {
		err = downloadAndDecompress()
		if err != nil {
			log.Fatalf("Failed to download or decompress file: %v", err)
		}
	}

	records = Parser(_CONST_LATEST_BVIEW_FILEPATH, flagWhois)

	if flagOutputPath != "" {
		writeCSV(records, flagOutputPath)
		log.Printf("Saved to %s", flagOutputPath)
	} else {
		fmt.Printf("%s", "\nResults:\n")
		fmt.Printf("%s", "AS, AS Name, Loops, Repeats, Consecutive Repeats (Prependings), Non-Consecutive Repeats\n")
		for as, loop := range records {
			fmt.Printf(
				"%s, %s, %d, %d, %d, %d\n",

				"AS"+fmt.Sprintf("%d", as),
				loop.Name,
				loop.Loops,
				loop.Repeats,
				loop.ConsecutiveRepeats, loop.NonConsecutiveRepeats,
			)
		}
	}
}
