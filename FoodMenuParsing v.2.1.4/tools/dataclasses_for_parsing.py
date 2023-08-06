from dataclasses import dataclass
from typing import Optional, List, Dict


@dataclass(init=True, repr=True, eq=True, order=False, unsafe_hash=True, frozen=False)
class Dish:
    __doc__ = "Complete information about particular dish."
    __slots__ = ["name", "ingredients", "nutrition_value", "group"]
    name: Optional[str]
    ingredients: Optional[str]
    nutrition_value: Optional[Dict[str, float]]
    group: Optional[str]


@dataclass(init=True, repr=True, eq=False, order=False, unsafe_hash=False, frozen=False)
class ParsingRecord:
    __doc__ = "Contains all data that will be stored in the database except ID. " \
              "This class is being compared (__eq__) by date attribute."
    __slots__ = ["date", "breakfast", "lunch", "dinner"]

    date: Optional[str]
    breakfast: Optional[List[Dish]]
    lunch: Optional[List[Dish]]
    dinner: Optional[List[Dish]]

    def __eq__(self, other):
        if isinstance(other, ParsingRecord):
            return self.date == other.date
        return NotImplemented

    def __hash__(self):
        return hash(self.date)


