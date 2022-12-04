import csv
import os
import time
from storages import mongodb
from storages.models import SimilarCouple
from fuzzywuzzy import process, fuzz


FOLDER_VACANCIES = "C:\Projects\Go\src\Vacancies"
SKILLS_COLUMN = 9

def get_vacancies_skills():
    skills = []
    for file in os.listdir(FOLDER_VACANCIES):
        file_skills = get_skills_from_CSV(os.path.join(FOLDER_VACANCIES, file))
        skills += file_skills
    return skills


def get_skills_from_CSV(csvPath: str):
    print(csvPath)
    file_skills = []
    with open(csvPath, newline="", encoding="utf-8") as csvfile:
        reader = csv.reader(csvfile, delimiter=";")
        for index, row in enumerate(reader):
            if index == 0: continue # Пропускаем название колонки
            try:
                file_skills += row[SKILLS_COLUMN].split("|")
            except IndexError:
                continue
    return file_skills

def save_vacancies_skill_inMongo(skills: set[str], collectionName: str):
    db = mongodb.connect_to_db("skills")
    collection = db[collectionName]
    for skill in skills:
        if skill: collection.insert_one({"skill": skill}) 
    

def clear_vacancies_skill_from_duplicates(skills: set[str] = None):
    # skills = mongodb.get_all(dbName="skills", collection="vacancies")
    skills = get_vacancies_skills()
    skills_without_duplicates = simple_detection_duplicates(skills)
    save_vacancies_skill_inMongo(skills_without_duplicates, "all_without_duplicates")


def detect_duplicates_in_skills(skills: set[str]):
    PairIndex = 0
    SimilarityIndex = 1
    LimitOfDifferenceInLength = 5
    AcceptablePercentOfSimilariry = 90

    comparison_list = skills
    for skill in skills[:1000]:
        if not skill: continue
        duplicates = []
        suspect_skills = process.extract(skill, comparison_list)
        if len(suspect_skills) < 2: continue
        for index, suspect in enumerate(suspect_skills):
            if index == 0: continue # первый в списке подозреваемых сам навык, поэтому пропускаем его
            if suspect[SimilarityIndex] >= AcceptablePercentOfSimilariry and (abs(len(suspect[PairIndex]) - len(skill)) < LimitOfDifferenceInLength): # Проверяем, что разница в количестве символов не превышает 5
                duplicates.append(suspect[PairIndex])
        
        comparison_list = remove_duplicates(comparison_list, duplicates)
    return comparison_list

def simple_detection_duplicates(skills):
    duplicates = []
    for index1 in range(len(skills)):
        skill1 = skills[index1]
        for index2 in range(index1+1, len(skills)):
            skill2 = skills[index2]
            if skill1.lower().strip() == skill2.lower().strip():
                duplicates.append(skill2)
    return remove_duplicates(skills, duplicates)


def remove_duplicates(skills: list[str], duplicates: list[str]):
    print("Количество повторений ", len(duplicates))
    for item in duplicates:
        try:skills.remove(item)
        except ValueError: continue
    return skills

if __name__ == "__main__":
    start = time.time()
    # skills = get_vacancies_skills()
    # save_vacancies_skill_inMongo(set(skills))
    clear_vacancies_skill_from_duplicates()
    print("time:", time.time()-start)