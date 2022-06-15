package main

/*
Written by Kevin Gillanders - 2022-06-15
*/


import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/k0kubun/pp/v3"
)

type SqlOutline struct{
	source, destination, key string
	columns [] string
}


func main() {
	colDumpFile := "colInput.txt"
	log.SetFlags(log.Lshortfile)
	
	log.Println("Generating SQL outline")
	SqlOutlineFile :=  GenerateSQLOutlineFile(colDumpFile)
	log.Println("Done Generating SQL outline")

	log.Println("Reading SQL outline from : ", SqlOutlineFile)
	sqlOutline := *ReadInSQLFile(SqlOutlineFile)
	log.Println("Done reading in outline")

	pp.Print(sqlOutline)

	OutputSQL(sqlOutline)	
}

func GenerateSQLOutlineFile(colDumpFile string) string {
	outPutFileName := "SqlOutline.txt"
	outPutFile, _  := os.Create(outPutFileName)
	defer outPutFile.Close()
	file, _ := os.Open(colDumpFile)

	scanner := bufio.NewScanner(file)
	var currentTable string
	var cols []string
	for scanner.Scan(){
		tableOutline := strings.Fields(scanner.Text())
		if tableOutline[2] != currentTable{
			
			if currentTable != ""{
				outPutFile.WriteString(fmt.Sprint(strings.Join(cols, "	"), "\n"))
			}
			cols = nil

			db := tableOutline[0]
			dbType := tableOutline[1]
			currentTable = tableOutline[2]
			ID := tableOutline[3]
			
			outPutFile.WriteString(fmt.Sprintf("%v [%v].[%v].%v %v\n", 
				currentTable, 
				db,
				dbType,
				currentTable, 
				ID))
		}
		cols = append(cols, tableOutline[3])

	}	

	return outPutFileName
}


func OutputSQL(sqlOutlines [] SqlOutline) {
	
	// update {tableName}(destination)
	// set {foo = b.foo
	// 		...} (cols)
	// from {[srcDB].[dbo].tableName}(source) as b
	// where {tableName}(destinaion).{Member_ID}(Key) = b.{Member_ID}(Key) 

	outPutFileName := "OutputSQL.sql"
	outPutFile, _  := os.Create(outPutFileName)

	outPutFile.WriteString("\nBEGIN TRAN\n\n")

	sqlToPrint := "\nUPDATE %v\nSET\n%v\nFROM\n\t%v AS b\nWHERE\n\t%v.%v = b.%v\n--========================\n\n"
	var colDef string

	for _, sqlOutline := range sqlOutlines{

		for i, col := range sqlOutline.columns{
			
			colDef = fmt.Sprint(colDef, "\t", col, " = b.", col)
			
			if i != len(sqlOutline.columns) - 1 {
				colDef = fmt.Sprint(colDef, ",\n")
			}	
		}

		outPutFile.WriteString(fmt.Sprintf(sqlToPrint, sqlOutline.destination, colDef, sqlOutline.source, sqlOutline.destination, sqlOutline.key, sqlOutline.key))
	}
	outPutFile.WriteString("\n\nROLLBACK\n--COMMIT")

}
	
func ReadInSQLFile(inputFile string) *[]SqlOutline{
	
	sqlOutlines := []SqlOutline{}
	file, _ := os.Open(inputFile)

	scanner := bufio.NewScanner(file)


	for scanner.Scan(){
		infoLine := strings.Split(scanner.Text(), " ")

		scanner.Scan()
		cols := strings.Fields(scanner.Text())

		sqlOutline := new(SqlOutline)

		sqlOutline.destination = infoLine[0]
		sqlOutline.source = infoLine[1]
		sqlOutline.key = infoLine[2]

		sqlOutline.columns = cols

		sqlOutlines = append(sqlOutlines, *sqlOutline)

	}

	return &sqlOutlines

}
