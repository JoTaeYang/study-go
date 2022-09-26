package main

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/oauth2/google"
	"gopkg.in/Iwark/spreadsheet.v2"
)

type SheetName struct {
	Names []string `json:"names"`
}

var sheetId string = "1T3La6Nk9-iWhmydex7zqc3ynIV5IAaQIyU4dEuz1DMQ"

func main() {
	data, err := ioutil.ReadFile("google-csv-read.json")
	if err != nil {
		panic(err)
	}

	readSheet, err := ioutil.ReadFile("sheet-name.json")
	if err != nil {
		panic(err)
	}
	_ = readSheet

	conf, err := google.JWTConfigFromJSON(data, spreadsheet.Scope)
	if err != nil {
		panic(err)
	}

	client := conf.Client(context.Background())
	service := spreadsheet.NewServiceWithClient(client)

	log.Println(sheetId)
	spreadsheet, _ := service.FetchSpreadsheet("1T3La6Nk9-iWhmydex7zqc3ynIV5IAaQIyU4dEuz1DMQ")

	{
		sheet, err := spreadsheet.SheetByTitle("TableData_Account")
		if err != nil {
			panic(err)
		}

		datas := make([]([]string), 0)
		for _, row := range sheet.Rows {
			rowCell := make([]string, 0)
			for _, cell := range row {
				value := fmt.Sprintf("\"%s\"", cell.Value)
				rowCell = append(rowCell, value)
			}
			datas = append(datas, rowCell)
		}

		file, err := os.Create("TableData_Account.csv")
		defer file.Close()

		if err != nil {
			panic(err)
		}

		// csv writer 생성
		wr := bufio.NewWriter(file)
		// utf8 bom mark
		wr.WriteByte(0xEF)
		wr.WriteByte(0xBB)
		wr.WriteByte(0xBF)

		// csv 내용 쓰기
		for _, row := range datas {
			for _, v := range row {

				wr.Write([]byte(v))
				wr.Write([]byte(","))
			}
			wr.Write([]byte("\r\n"))
		}
		wr.Flush()
	}

}

//https://github.com/Iwark/spreadsheet
