from typing import NamedTuple
from bson import ObjectId


class Test(NamedTuple):
    Id: ObjectId
    VacancyId: str
    VacancyTitle: str
    Edwica: str
    File: str

class SimilarCouple(NamedTuple):
    id: int
    demand_id : int
    demand_name: str
    dup_demand_id: int
    dup_demand_name: str
    similarity: int
    is_duplicate: bool
