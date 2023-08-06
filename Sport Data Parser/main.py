from datetime import datetime
from typing import List, Tuple, Any

from Parser import Parser
import traceback
import logging
import sqlite3
import argparse
import time
import re

import ConstantExpressions

database_filename = "../database.db"


def time_remaining(task):
    """
    :param task: Current task from scheduled ones
    :return: Time until task will be executed
    """
    next_run = task.next_run
    now = datetime.now()
    time_diff = next_run - now
    return time_diff


def is_time24_format(time_str) -> str:
    """
    :param time_str: Possible string to check if it satisfies 24h time format.
    :return: True if time_str satisfies format HH:mm else False
    """
    if not re.match(ConstantExpressions.HOURS24_FORMAT_CHECK, time_str):
        raise argparse.ArgumentTypeError(f"Invalid time format: '{time_str}'. Expected format: HH:mm (24-hour).")
    return time_str


def parse_data() -> Tuple[List[Any], List[Any]] or None:
    logging.info("Starting parsing...")

    # Creating parser object - required for parsing any information.
    parser = Parser()

    try:
        events, sport_infos = parser.parse_all_info()
        logging.info("Successfully parsed all information.")
        return events, sport_infos
    except Exception as e:
        logging.critical(f"Parsing failed. Exception: {e}")
    return None


def write_to_database(events: List[Any], sport_infos: List[Any]) -> bool:
    if events is None or sport_infos is None:
        return False
    else:
        try:
            # Connecting to the database
            sqlite_connection = sqlite3.connect(database_filename)
            cursor = sqlite_connection.cursor()
            logging.info("Successfully connected to database")

            # Clearing contents of SportsGames before writing new data
            logging.info("Started clearing the contents of SportsGames table")
            cursor.execute(ConstantExpressions.SPORTS_GAMES_SQL_DELETE_QUERY)
            logging.info("Cleared the contents of SportsGames table")

            # Inserting parsed data to SportsGames
            #
            logging.info("Inserting information to SportsGames")
            # Getting records from SportGame record (EventRecord)
            cursor.executemany(ConstantExpressions.SPORTS_GAMES_SQL_INSERT_QUERY,
                               tuple(el.get_tuple() for el in events))
            sqlite_connection.commit()

            # Clearing contents of SportsInfo before writing new data
            logging.info("Started clearing the contents of SportsInfo table")
            cursor.execute(ConstantExpressions.SPORTS_INFO_SQL_DELETE_QUERY)
            logging.info("Cleared the contents of SportsInfo table")

            # Inserting parsed data to SportsInfo
            #
            logging.info("Inserting information to SportsInfo")
            # Getting records from SportInfo record (SportInfoRecord)
            cursor.executemany(ConstantExpressions.SPORTS_INFO_SQL_INSERT_QUERY,
                               tuple(el.get_tuple() for el in sport_infos))
            sqlite_connection.commit()

            # Closing connection with database
            cursor.close()
            logging.info("Successfully uploaded information to database")
            return True

        # Logging information about working with database
        except sqlite3.Error as error:
            logging.error(f"Writing to database failed with sqlite3.Error: {error}")

        except Exception as ex:
            tb = traceback.format_exc()
            logging.error(f"Writing to database failed with unhandled exception: {ex}\n{tb}")
    return False


def setup() -> None:
    logging.basicConfig(filename=f'launch_{datetime.now()}', level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')


def update() -> None:
    sport_events, sport_information = parse_data()
    commit_factor = False
    if not(sport_events is None) and not(sport_information is None):
        commit_factor = write_to_database(sport_events, sport_information)
    if commit_factor:
        logging.info("Update finished.")
        return
    else:
        time.sleep(360)
        update()


if __name__ == "__main__":
    setup()     # Setup of argument parser, logging
    update()    # Regular parsing based on a schedule
