import pymysql
from pymysql.connections import Connection
import os
from pydantic import BaseSettings

from storages.models import SimilarCouple

class MYSQL(BaseSettings):
    HOST = os.getenv("EDWICA_DB_HOST")
    USER = os.getenv("EDWICA_DB_USER")
    PASS = os.getenv("EDWICA_DB_PASS")
    PORT = 3306
    DB = "edwica"
    TABLE = "demand_duplicate"


def connect() -> Connection:
    settings = MYSQL()
    try:
        connection = pymysql.connect(
            host=settings.HOST,
            port=settings.PORT,
            database=settings.DB,
            password=settings.PASS,
            user=settings.USER,
            cursorclass=pymysql.cursors.DictCursor
        )
        return connection
    except Exception as err:
        exit(f"Error by connection: {err}")


def get_all_skills() -> list[SimilarCouple]:
    connection = connect()
    with connection.cursor() as cursor:
        cursor.execute(f"SELECT * FROM {MYSQL().TABLE}")
        result = [SimilarCouple(*item.values()) for item in cursor.fetchall()]
        return result

def get_valide_skills() -> list[str]:
    valide_skills: list[str] = []
    skills = get_all_skills()
    
    for skill in skills:
        if skill.is_duplicate:
            valide_skills.append(skill.demand_name)
        else:
            valide_skills.extend((skill.demand_name, skill.dup_demand_name))
    return valide_skills

