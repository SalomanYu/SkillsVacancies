from pymongo import MongoClient, database
from typing import  Literal
from storages.models import SimilarCouple




def connect_to_db(dbName: str) -> database.Database:
    client = MongoClient("localhost", 27017)
    db = client[dbName]
    return db


def get_all(dbName: str, collection: str):
    return filterDocuments(dbName=dbName, collectionName=collection, filter=None)


def filterDocuments(dbName: str, collectionName: str, filter: dict[Literal["csvfile", "_id", "vacancy_id"]] | None = None) -> list[SimilarCouple]:

    documents: list[SimilarCouple] = []
    db = connect_to_db(dbName)
    collection = db[collectionName]
    if filter is None:
        for doc in collection.find():
            documents.append(doc["skill"])
        return documents
    else:
        for doc in collection.find(filter):
            documents.append(doc["skill"])
        return documents

