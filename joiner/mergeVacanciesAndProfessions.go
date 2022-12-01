package joiner

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/SalomanYu/SkillsVacancies/storage"

	"go.mongodb.org/mongo-driver/mongo"
)

const (
	FolderUpdatedVacancies	 =	"Data/UpdatedVacancies"
	ColumnVacancyCode		 = 	0
)

var documents []*storage.Vacancy

func CombineTheFoundPairsWithVacancies() {
	vacancyFiles, err := ioutil.ReadDir(FolderVacancies)
	checkErr(err)
	for _, file := range vacancyFiles{
		if !strings.HasSuffix(file.Name(), ".csv") {
			continue
		}
		csvPath := filepath.Join(FolderVacancies, file.Name())
		pathForCompare := fmt.Sprintf("%s\\%s", FolderVacancies, file.Name())
		
		documents, err = storage.GetVacanciesCsvFile(pathForCompare)
		fmt.Println(csvPath)
		findMatchesInCSV(pathForCompare)
		if err == mongo.ErrNoDocuments{
			log.Printf("В этом файле нет ни одного совпадения: %s\n", csvPath)
			continue
		}
		checkErr(err)
	}
}

func findMatchesInCSV(csvPath string) {
	rows := readCsvFile(csvPath)
	rowsUpdated := [][]string{}
	tableHeader := append(rows[0], []string{"Подходящая_профессия"}...)
	rowsUpdated = append(rowsUpdated, tableHeader)
	for _, row := range rows[1:]{
		record := getRowWithMatchColumn(row)
		rowsUpdated = append(rowsUpdated, record)
	}
	saveUpdatedCsv(csvPath, rowsUpdated)
	
}

func readCsvFile(filePath string) (records [][]string) {
	file, err := os.Open(filePath)
	checkErr(err)
	defer file.Close()

	csvReader := csv.NewReader(file)
	csvReader.Comma = ';'
	csvReader.LazyQuotes = true
	csvReader.FieldsPerRecord = -1
	records, err = csvReader.ReadAll()
	checkErr(err)
	return
}

func getRowWithMatchColumn(row []string) (record []string){
	for _, vacancy := range documents{
		if vacancy.IdVacancy == row[ColumnVacancyCode] {
			record := append(row, []string{vacancy.MatchTitle}...)
			return record
		}
	}
	return row
}

func saveUpdatedCsv(filename string, data [][]string) {
	new_filename := strings.Split(filename, "\\")[1]
	file, err := os.Create(filepath.Join(FolderUpdatedVacancies, new_filename))
	checkErr(err)
	defer file.Close()
	writter := csv.NewWriter(file)
	writter.WriteAll(data)
	checkErr(writter.Error())
	fmt.Println("Updated file!")
}

