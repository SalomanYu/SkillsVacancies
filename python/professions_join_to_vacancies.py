import os
import csv
import xlrd
import pymongo

from fuzzywuzzy import process
from multiprocessing import Pool

VACANCIES_FOLDER = "Data/Vacancies"
JOINED_PROFESSIONS = "Data/JoindedProfessions.xlsx"


def get_all_edwica_professions() -> set[str]:
    """Возвращает генератор"""

    book = xlrd.open_workbook(JOINED_PROFESSIONS)
    sheet = book.sheet_by_index(0)
    professions = set(prof for prof in sheet.col_values(1)[1:])
    return professions


def find_similared_professions_if_file(file: str):
    print(file)
    csvfile = file + "" # копируем строку
    edwica_professions = get_all_edwica_professions()
    with open(file, newline="", encoding="utf-8") as file:
        rows = tuple(csv.reader(file, delimiter=";"))[1:] # Убираем заголовки
        for row in rows:
            prof = row[2]
            search_professions = edwica_professions.copy() # Копируем, потому что если бы мы дальше передавали edwica_professions, То они бы очищались
            similared = process.extract(prof, search_professions, limit=1)
            if similared:
                if similared[0][1] > 90:
                    client = pymongo.MongoClient("mongodb://localhost:27017/")
                    db = client["hh_ru"]
                    collection = db["vacancies"]
                    try:collection.insert_one({
                        "vacancy_id": row[0],
                        "vacancy_title": prof,
                        "edwica_profession": similared[0][0],
                        "csvfile": csvfile
                    })
                    except:continue

def main():
    files = (os.path.join(VACANCIES_FOLDER, file) for file in os.listdir(VACANCIES_FOLDER) if file.endswith(".csv"))
    with Pool(5) as process:
        process.map_async(
            func=find_similared_professions_if_file,
            iterable=files, 
            error_callback=lambda x: exit(f"Ошибка многопоточности:x"))
        process.close()
        process.join()

if __name__ == "__main__":    
    main()
