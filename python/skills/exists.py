from storages import mysqldb, mongodb
from filter_vacancies_skills import clear_vacancies_skill_from_duplicates

vacancies_skills = (i.lower() for i in mongodb.get_all(dbName="skills", collection="without_duplicates"))
sorted_by_bot_skills = (i.lower() for i in mysqldb.get_valide_skills())

skills = set(vacancies_skills) & set(sorted_by_bot_skills)
print(len(skills))

