from fake_useragent import UserAgent
from typing import List
import requests
import json

from tools.logging_tools import log_class_methods
from tools.time_tools import get_current_date
from tools.dataclasses_for_parsing import Dish, ParsingRecord


@log_class_methods
class FoodMenuParser:
    __doc__ = "Parser for avonoldfarms.api.flikisdining.com. Gathers information through API."
    __slots__ = ["__food_intakes", "data_parsed"]

    def __init__(self):
        self.__food_intakes = ["breakfast", "lunch", "dinner"]
        self.data_parsed = list()

    @staticmethod
    def __fetch(link):
        return requests.get(link, headers={"User-Agent": UserAgent().random}).text

    @staticmethod
    def __get_api_request_link(food_intake, from_year, from_month, from_day):
        return f"https://avonoldfarms.api.flikisdining.com/menu/api/weeks/school/avon-old-farms/menu-type/" \
               f"{food_intake}/{from_year}/{from_month}/{from_day}/?format=json"

    @staticmethod
    def __extract_group_from_json(j, dict_with_food_types=None):
        if dict_with_food_types and j['menu_id'] in dict_with_food_types.keys():
            return dict_with_food_types[j['menu_id']]
        try:
            if j['text'] == '':
                return "N/A"
            else:
                return j['text']
        except (AttributeError, TypeError):
            return "N/A"

    @staticmethod
    def __extract_ingredients_from_json(j):
        try:
            ingredients = j['food']['ingredients']
        except (AttributeError, TypeError):
            ingredients = "N/A"
        return ingredients

    @staticmethod
    def __extract_nutrition_value(j):
        try:
            nutrition_info = j['food']['rounded_nutrition_info']
        except (AttributeError, TypeError):
            nutrition_info = "N/A"
        return nutrition_info

    @staticmethod
    def __does_name_appear_in_json(j):
        try:
            _ = j['food']['name']
            return True
        except (AttributeError, TypeError):
            return False

    def _get_json_response_from_api(self, food_intake: str, from_date: tuple):
        if from_date:
            link = self.__get_api_request_link(food_intake, from_date[0], from_date[1], from_date[2])
        else:
            from_year, from_month, from_day = get_current_date()
            link = self.__get_api_request_link(food_intake, from_year, from_month, from_day)
        return self.__fetch(link)

    def __determine_active_record(self, date):
        active_record: ParsingRecord
        for record in self.data_parsed:
            if record.date == date:
                active_record = record
                break
        else:
            active_record = ParsingRecord(
                date=date,
                breakfast=None,
                lunch=None,
                dinner=None
            )
            self.data_parsed.append(active_record)

        return active_record

    def __extract_dishes_from_menu_items(self, menu_items):
        dishes = []
        dict_with_food_type = {}
        # for dish_information in menu_items:
        #     if (g := self.__extract_group_from_json(dish_information)) != "N/A":
        #         dict_with_food_type[dish_information['menu_id']] = g

        saved_group = 'N/A'
        for dish_information in menu_items:

            active_group = self.__extract_group_from_json(dish_information)
            if active_group == 'N/A':
                active_group = saved_group
            else:
                saved_group = active_group

            if not self.__does_name_appear_in_json(dish_information):
                continue

            ingredients = self.__extract_ingredients_from_json(dish_information)
            nutrition_value = self.__extract_nutrition_value(dish_information)

            dishes.append(Dish(group=active_group,
                               ingredients=ingredients,
                               nutrition_value=nutrition_value,
                               name=dish_information['food']['name']))

        return dishes

    def __process_days_of_food_intake(self, food_intake_json, food_intake):
        reselected_food_intake_json = food_intake_json.get("days")
        for day in range(0, 6):
            date = reselected_food_intake_json[day].get("date")
            menu_items = reselected_food_intake_json[day].get("menu_items")
            active_record = self.__determine_active_record(date)
            setattr(active_record, food_intake, self.__extract_dishes_from_menu_items(menu_items))

    def parse_data(self, from_date: tuple = None) -> None:
        for food_intake in self.__food_intakes:
            raw_response = self._get_json_response_from_api(food_intake, from_date=from_date)
            food_intake_json = json.loads(raw_response)
            self.__process_days_of_food_intake(food_intake_json, food_intake)

    def get_data(self) -> List[ParsingRecord]:
        return self.data_parsed
