from typing import List
import json

from tools.dataclasses_for_parsing import ParsingRecord


class FormatChanger:
    def __init__(self, parsed_content: List[ParsingRecord]):
        self.__content_to_change = parsed_content

    @staticmethod
    def _dish_list(dishes):
        if dishes is None:
            return []
        return [
            {
                "name": dish.name,
                "ingredients": dish.ingredients,
                "nutrition_value": dish.nutrition_value,
                "group": dish.group
            }
            for dish in dishes
        ]

    def get_output(self):
        result = list()
        id_to_use = 1

        for parsing_record in self.__content_to_change:
            r = list()

            r.append(id_to_use)
            id_to_use += 1

            r.append(str(parsing_record.date))
            r.extend(json.dumps(self._dish_list(dish_group)) for dish_group in (parsing_record.breakfast,
                                                                                parsing_record.lunch,
                                                                                parsing_record.dinner))
            result.append(tuple(r))

        return tuple(result)
