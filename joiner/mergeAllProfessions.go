package joiner

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/xuri/excelize/v2"
)

const (
	FolderProfessions 	 	 = 	"Data/Professions"
	FolderVacancies 	 	 = 	"Data/Vacancies"
	FileJoinedProfessions 	 = 	"Data/JoindedProfessions.xlsx"
	SheetProfessions		 = 	"Вариации названий"
)



// Метод собирает все эксель файлы с профессиями объединяет их в один общий файл для того,
// Чтобы на питоне проще было искать совпадения названий вакансий с названиями профессий Эдвики.
// Результат выполнения метода не используется дальше языком Go.
// Почему метод написан на Golang, а не на Python? 1) Golang выполнит эту задачу быстрее. 2) Способ попрактироваться
func JoinAllEdwicaProfessionsInOneFile(){
	var rowNum = 2
	files, err := ioutil.ReadDir(FolderProfessions)
	checkErr(err)
	xlsx := excelize.NewFile() 
	xlsx.SetCellValue("Sheet1", "A1", "id")	
	xlsx.SetCellValue("Sheet1", "B1", "name")	
	xlsx.SetCellValue("Sheet1", "C1", "file")	

	for _, file := range files{
		if strings.HasSuffix(file.Name(), ".xlsx"){
			professions := getAllEdwicaProfessions(filepath.Join(FolderProfessions, file.Name()))
			for _, prof := range professions{
				xlsx.SetCellValue("Sheet1", fmt.Sprintf("A%d", rowNum), rowNum-1)	
				xlsx.SetCellValue("Sheet1", fmt.Sprintf("B%d", rowNum), prof)	
				xlsx.SetCellValue("Sheet1", fmt.Sprintf("C%d", rowNum), file.Name())	
				rowNum++
			}
		}
	}
	if err := xlsx.SaveAs(FileJoinedProfessions); err != nil{
		log.Fatal(err)
	}
}

func getAllEdwicaProfessions(file string) (professions []string){
	xlsx, err := excelize.OpenFile(file)
	checkErr(err)
	columns, err := xlsx.GetCols(SheetProfessions)
	checkErr(err)
	professions = removeEmpty(columns[6][1:]) // Не берем заголовок колонки
	return 

}

// Метод excelize.GetCols() возвращает все доступные ячейки, поэтому нужен removeEmpty, который будет вырезать пустые ячейки
func removeEmpty(slice []string) (sliceWithoutEmptyStrings []string){
	for _, str := range slice{
		if str != ""{
			sliceWithoutEmptyStrings = append(sliceWithoutEmptyStrings, str)
		}
	}
	return 
}

func checkErr(err error){
	if err != nil{
		log.Fatal(err)
	}
}