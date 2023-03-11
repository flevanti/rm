package refresh

import "database/sql"

/*
This logic can be used to import the csv file from https://www.doogal.co.uk/files/postcodes.zip
Csv file needs to be extracted and present in the current working directory.
File must be names postcodes.csv
*/

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

var sqlParamsBuffer []string   // parameters usend in query
var sqlRecsInBuffer int        //number of records currently in buffer
var sqlBufferMaxRecords = 3000 //max number of records in buffer, this will trigger a buffer flush (db insert)
var sqlBufferFlushes int       //number of times the buffer has been flushed
var totRecsBuffered int
var totLinesProcessed int
var totLinesInFile int

var db *sql.DB

func Postcodes(dbParam *sql.DB) {
	db = dbParam
	ctx, ctxCancelFunc := context.WithCancel(context.Background())
	defer ctxCancelFunc()

	truncateTargetTable()

	file, err := os.Open("./postcodes.csv")
	if err != nil {
		log.Fatal(err)
	}

	totLinesInFile, err = lineCounter(file)
	if err != nil {
		log.Fatal(err)
	}
	//reset file pointer to the beginning of the file...
	file.Seek(0, io.SeekStart)

	parser := csv.NewReader(file)
	firstLoop := true
	go printProgress(ctx)

	for {
		record, err := parser.Read()
		totLinesProcessed++
		if err == io.EOF {
			break
		}
		if firstLoop {
			//skip headers...
			firstLoop = false
			continue
		}

		if err != nil {
			log.Fatal(err)
		}

		sqlRecsInBuffer++
		totRecsBuffered++
		addRecToBuffer([]string{
			strings.TrimSpace(record[0]),
			strings.ReplaceAll(strings.TrimSpace(record[0]), " ", ""),
			strings.TrimSpace(record[1]),
			strings.TrimSpace(record[12]),
			strings.TrimSpace(record[8]),
		})

	} //end for each row of file
	if sqlRecsInBuffer > 0 {
		flushBuffer()
	}
	fmt.Print("\n\nProcess completed!\n")
}

func lineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

func truncateTargetTable() {
	sql := "truncate table postcodes_source;"
	db.Exec(sql)
}

func printProgress(ctx context.Context) {
	ticker := time.NewTicker(250 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			fmt.Printf("Rows processed %d (%.2f%%), rows in buffer %d, flushes count %d          \r",
				totRecsBuffered,
				float64(totLinesProcessed)/float64(totLinesInFile)*100,
				sqlRecsInBuffer,
				sqlBufferFlushes,
			)
		case <-ctx.Done():
			fmt.Println()
			return
		}
	}
}

func addRecToBuffer(params []string) {
	sqlParamsBuffer = append(sqlParamsBuffer, params...)
	if sqlRecsInBuffer >= sqlBufferMaxRecords {
		sqlBufferFlushes++
		flushBuffer()
	}
}
func flushBuffer() {
	var sqlParamsBufferLocal []interface{}
	sqlFields := 5 // this is the number of fields we are saving
	sql := "insert into postcodes_source (postcode, postcode_spaces, in_use, country, county)" +
		"values " + createSqlPlaceholdersString(sqlFields, sqlRecsInBuffer) + ";"

	for _, v := range sqlParamsBuffer {
		sqlParamsBufferLocal = append(sqlParamsBufferLocal, v)
	}

	_, err := db.Exec(sql, sqlParamsBufferLocal...)
	if err != nil {
		panic(err)
	}
	resetBuffer()
}

func resetBuffer() {
	sqlRecsInBuffer = 0
	sqlParamsBuffer = nil
}

func repeatStringWithSeparator(str string, separator string, c int) string {
	return strings.TrimRight(strings.Repeat(str+separator, c), separator)
}

func createSqlPlaceholdersString(fieldsCount int, recordsCount int) string {
	return repeatStringWithSeparator("("+repeatStringWithSeparator("?", ",", fieldsCount)+")", ", ", recordsCount)
}
