import logging
from FormatChanger import FormatChanger
from FoodMenuParser import FoodMenuParser
from DatabaseFoodMenuWriter import DatabaseFoodMenuWriter

from tools.logging_tools import setup_logging, log_execution
from tools.time_tools import get_current_date, get_next_week_date


@log_execution
def update():
    food_menu_parser = FoodMenuParser()

    date = get_current_date()
    food_menu_parser.parse_data(date)
    for _ in range(3):
        date = get_next_week_date(date)
        food_menu_parser.parse_data(date)

    gathered_data = food_menu_parser.get_data()
    
    format_changer = FormatChanger(gathered_data)
    content_to_write_to_db = format_changer.get_output()

    with DatabaseFoodMenuWriter() as database_food_menu_writer:
        database_food_menu_writer.execute_food_menu_data_writing(content_to_write_to_db)


@log_execution
def setup():
    setup_logging(filename="food_menu_parsing.log")


if __name__ == "__main__":
    setup()
    update()
    logging.info("END".center(130, "="))
