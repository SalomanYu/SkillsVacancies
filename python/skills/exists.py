from storages import mysqldb, mongodb
from filter_vacancies_skills import save_vacancies_skill_inMongo

print([i for i in mongodb.get_all(dbName="skills", collection="all_without_duplicates")][:50])
vacancies_skills = next(i for i in mongodb.get_all(dbName="skills", collection="all_without_duplicates"))
print(vacancies_skills)
sorted_by_bot_skills = (i.lower().strip() for i in mysqldb.get_valide_skills())

exists_skills = set(vacancies_skills) & set(sorted_by_bot_skills) # Выделяем общие навыки из вакансий и навыки, обработанные ботом
unknown_skills = set(vacancies_skills) | set(sorted_by_bot_skills) 
print(len(exists_skills))
print(len(unknown_skills))

save_vacancies_skill_inMongo(exists_skills, collectionName="exists_skills")
save_vacancies_skill_inMongo(unknown_skills, collectionName="unknown_skills")

