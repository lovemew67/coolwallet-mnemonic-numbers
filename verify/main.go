package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/gocarina/gocsv"
)

func main() {
	log.Println("start")
	defer log.Println("end")

	file, err := os.Open("../number-seed-conversion.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	type mappingRecord struct {
		Number string `csv:"NUMBER"`
		Word   string `csv:"WORD"`
	}
	var records = make([]*mappingRecord, 0)
	err = gocsv.UnmarshalFile(file, &records)
	if err != nil {
		panic(err)
	}
	log.Printf("🧐 number of records:\t%d\n", len(records))

	recordMap := make(map[string]string)
	numberMap := make(map[string]string)
	for _, record := range records {
		recordMap[record.Word] = record.Number
		numberMap[record.Number] = record.Word
	}
	log.Printf("🧐 size of record map:\t%d\n", len(recordMap))
	log.Printf("🧐 size of number map:\t%d\n", len(numberMap))

	file, err = os.Open("../english.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	wordMap := make(map[string]struct{})
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		wordMap[sc.Text()] = struct{}{}
	}
	log.Printf("🧐 number of words:\t\t%d\n", len(wordMap))

	err = sc.Err()
	if err != nil {
		panic(err)
	}

	for word := range wordMap {
		if _, ok := recordMap[word]; !ok {
			panic(fmt.Errorf("🤕 missing record: %s", word))
		}
	}
	log.Println("🥳 all words are in mapping table")

	file, err = os.Create("../seed-number-conversion-full.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	type revertedRecord struct {
		Word   string `csv:"WORD"`
		Number string `csv:"NUMBER"`
	}

	var fullRevertedRecords = make([]*revertedRecord, 0)
	for _, record := range records {
		fullRevertedRecords = append(fullRevertedRecords, &revertedRecord{Word: record.Word, Number: record.Number})
	}
	err = gocsv.MarshalFile(&fullRevertedRecords, file)
	if err != nil {
		panic(err)
	}

	file, err = os.Create("../seed-number-conversion-trimmed.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var trimmedRevertedRecords = make([]*revertedRecord, 0)
	for _, record := range records {
		word := record.Word
		if len(word) > 3 {
			word = word[:4]
		}

		trimmedRevertedRecords = append(trimmedRevertedRecords, &revertedRecord{
			Word:   word,
			Number: record.Number,
		})
	}
	err = gocsv.MarshalFile(&trimmedRevertedRecords, file)
	if err != nil {
		panic(err)
	}
}
