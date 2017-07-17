package gocsv

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sync"
)

type App struct {
	File  string
	Keys  []string
	Items []map[string]string
}

func (a *App) parseCSVData(parseItems chan []string) {
	csvfile, err := os.Open(a.File)

	if err != nil {
		fmt.Println(err)
		return
	}
	defer csvfile.Close()

	reader := csv.NewReader(csvfile)
	reader.FieldsPerRecord = -1
	rowCount := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if rowCount == 0 {
			a.Keys = record
		}
		rowCount++
		// Stop at EOF.
		go func(pi []string) {
			parseItems <- pi
		}(record)
	}
}

func (a *App) moldObject(parseItems <-chan []string, wg *sync.WaitGroup, lineItems chan<- map[string]string) {
	defer wg.Done()
	for pi := range parseItems {
		var l map[string]string
		for i, key := range a.Keys {
			fmt.Println(l[key])
			fmt.Println(pi[i])
			fmt.Println(i)
			l[key] = pi[i]
		}
		lineItems <- l
	}
}

func (a *App) Run() {
	parseItems := make(chan []string)
	lineItems := make(chan map[string]string)

	go a.parseCSVData(parseItems)

	wg := new(sync.WaitGroup)
	for i := 0; i <= 3; i++ {
		wg.Add(1)
		go a.moldObject(parseItems, wg, lineItems)
	}

	go func() {
		wg.Wait()
		close(lineItems)
	}()

	for li := range lineItems {
		a.Items = append(a.Items, li)
	}
}
