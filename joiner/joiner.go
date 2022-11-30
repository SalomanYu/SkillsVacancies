package joiner

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/SalomanYu/SkillsVacancies/mongoStorage"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	ProfessionsFolderPath 	 = 	"Data\\Professions"
	VacanciesFolderPath 	 = 	"Data/Vacancies"
	UpdatedVacanciesFolder	 =	"Data\\UpdatedVacancies"
	ProfessionsListSheetName = 	"Вариации названий"
	PairProfessionsJSON 	 = 	"PairProfessions"
)
var (
	// similarVacanciesAndProfessionsList []PairProfessions
	documents []*mongoStorage.Vacancy
	vacancyIdIndex = 0
)
type PairProfessions struct {
	IdVacancy		string	`json:"vacancy_id"`
	TitleVacancy	string	`json:"vacancy_title"`
	MatchProfession	string	`json:"edwica_profession"`
	PathVacancy		string	`json:"csvfile"`
}


func JoinAllEdwicaProfessionsInOneFile(){
	var rowNum = 2
	files, err := ioutil.ReadDir(ProfessionsFolderPath)
	CheckErr(err)
	xlsx := excelize.NewFile() 
	xlsx.SetCellValue("Sheet1", "A1", "id")	
	xlsx.SetCellValue("Sheet1", "B1", "name")	
	xlsx.SetCellValue("Sheet1", "C1", "file")	

	for _, file := range files{
		if strings.HasSuffix(file.Name(), ".xlsx"){
			professions := getAllEdwicaProfessions(filepath.Join(ProfessionsFolderPath, file.Name()))
			for _, prof := range professions{
				xlsx.SetCellValue("Sheet1", fmt.Sprintf("A%d", rowNum), rowNum-1)	
				xlsx.SetCellValue("Sheet1", fmt.Sprintf("B%d", rowNum), prof)	
				xlsx.SetCellValue("Sheet1", fmt.Sprintf("C%d", rowNum), file.Name())	
				rowNum++
			}
		}
	}
	if err := xlsx.SaveAs("Data/JoindedProfessions.xlsx"); err != nil{
		log.Fatal(err)
	}
}

func getAllEdwicaProfessions(file string) (professions []string){
	xlsx, err := excelize.OpenFile(file)
	CheckErr(err)
	columns, err := xlsx.GetCols(ProfessionsListSheetName)
	CheckErr(err)
	professions = removeEmpty(columns[6][1:]) // Не берем заголовок колонки
	return 

}

func removeEmpty(slice []string) (sliceWithoutEmptyStrings []string){
	for _, str := range slice{
		if str != ""{
			sliceWithoutEmptyStrings = append(sliceWithoutEmptyStrings, str)
		}
	}
	return 
}

func CombineTheFoundPairsWithVacancies() {
	// similarVacanciesAndProfessionsList = getSimilarVacanciesAndProfessionsFromJson(pairsJsonPath)
	vacancyFiles, err := ioutil.ReadDir(VacanciesFolderPath)
	CheckErr(err)
	for _, file := range vacancyFiles{
		if !strings.HasSuffix(file.Name(), ".csv") {
			continue
		}
		csvPath := filepath.Join(VacanciesFolderPath, file.Name())
		pathForCompare := fmt.Sprintf("%s\\%s", VacanciesFolderPath, file.Name())
		// findMatchesInCSV(filepath.Join(VacanciesFolderPath, file.Name()))
		documents, err = mongoStorage.GetVacanciesCsvFile(pathForCompare)
		if err == mongo.ErrNoDocuments{
			continue
		}
		CheckErr(err)
		findMatchesInCSV(csvPath)
	}
}

func getSimilarVacanciesAndProfessionsFromJson(jsonPath string) (payload []PairProfessions) {
	data, err := ioutil.ReadFile(jsonPath)
	CheckErr(err)
	err = json.Unmarshal(data, &payload)
	CheckErr(err)
	return
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
	CheckErr(err)
	defer file.Close()

	csvReader := csv.NewReader(file)
	csvReader.Comma = ';'
	records, err = csvReader.ReadAll()
	CheckErr(err)
	return
}

func getRowWithMatchColumn(row []string) (record []string){
	// for _, couple := range similarVacanciesAndProfessionsList{
	for _, vacancy := range documents{
		if vacancy.IdVacancy == row[vacancyIdIndex] {
			record := append(row, []string{vacancy.MatchTitle}...)
			return record
		}
	}
	return row
}

func saveUpdatedCsv(filename string, data [][]string) {
	new_filename := strings.Split(filename, "\\")[2]
	file, err := os.Create(filepath.Join(UpdatedVacanciesFolder, new_filename))
	CheckErr(err)
	defer file.Close()
	writter := csv.NewWriter(file)
	writter.WriteAll(data)
	CheckErr(writter.Error())
	fmt.Println("Updated file!")
}

func CheckErr(err error){
	if err != nil{
		log.Fatal(err)
	}
}