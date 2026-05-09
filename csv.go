package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

func writeCSV(loops TRecords, filepath string) {
	file, err := os.Create(filepath)
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"AS", "AS Name", "Loops", "Repeats", "Consecutive Repeats (Prependings)", "Non-Consecutive Repeats"}
	if err := writer.Write(header); err != nil {
		log.Fatalf("Failed to write header: %v", err)
	}

	for as, loop := range loops {
		err := writer.Write(
			[]string{
				"AS"+fmt.Sprintf("%d",as),
				fmt.Sprintf("%s", loop.Name),
				fmt.Sprintf("%d", loop.Loops),
				fmt.Sprintf("%d", loop.Repeats),
				fmt.Sprintf("%d", loop.ConsecutiveRepeats), fmt.Sprintf("%d", loop.NonConsecutiveRepeats),
			},
		)
		if err != nil {
			log.Fatalf("Failed to write record: %v", err)
		}
	}
}
