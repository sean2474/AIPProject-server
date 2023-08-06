from datetime import datetime
import random
import re

import ConstantExpressions


class EventRecord:
    def __init__(self, data: list, sport: str, category: str):
        self.raw_data = data

        self.event_date = self.__search_event_date()
        self.date_time = self.__search_date_time()
        self.unified_data = self.event_date + " " + self.date_time

        self.opponent_school = self.__search_opponent()

        self.category = self.__define_category(category)
        self.sport = sport

        self.is_away = self.__search_if_away()
        '''
        0 -> Away
        1 -> Home
        2 -> NA
        '''
        self.is_win = self.__search_win_data()
        self.game_score = self.__search_score_data()
        self.is_cancelled = self.__was_game_cancelled()

        self.game_location = self.__decide_game_location()

    @staticmethod
    def __alphabetical_date_format_to_numeric(date_string):
        if "TBD" in date_string:
            return "TBD"
        current_date = datetime.now()
        dt = datetime.strptime(date_string, "%a %b %d %I:%M %p")
        dt = dt.replace(year=current_date.year)

        if (dt.month > current_date.month) or\
                (dt.month == current_date.month and
                 dt.day > current_date.day) or \
                (dt.month == current_date.month and
                 dt.day == current_date.day and
                 dt.time() > current_date.time()):
            dt = dt.replace(year=dt.year - 1)

        return dt.strftime("%Y-%m-%d %I:%M %p")

    @staticmethod
    def __convert_date_format(date_string):
        dt = datetime.strptime(date_string, "%m/%d/%y")
        return dt.strftime("%a %b %d")

    @staticmethod
    def __define_sport(sport: str) -> int:
        sport_list = {
            "Football": 0,
            "Soccer": 1,
            "Cross Country": 2,
            "Hockey": 3,
            "Basketball": 4,
            "Squash": 5,
            "Swimming": 6,
            "Wrestling": 7,
            "Baseball": 8,
            "Lacrosse": 9,
            "Golf": 10,
            "Track and Field": 11,
            "Tennis": 12,
            "NA": 13
        }
        if sport in sport_list.keys():
            return sport_list[sport]
        elif "Track" in sport and "Field" in sport:
            return sport_list["Track & Field"]
        elif "Cross" in sport and "Country" in sport:
            return sport_list["Cross Country"]
        return 13

    @staticmethod
    def __define_category(category: str) -> int:
        sport_groups = {
            "Varsity": 0,
            "JV": 1,
            "Varsity B": 2,
            "Thirds": 3,
            "Thirds Blue": 4,
            "Thirds Red": 5,
            "Fourths": 6,
            "Fifths": 7,
            "NA": 8
        }
        if category in sport_groups.keys():
            return sport_groups[category]
        else:
            return sport_groups["NA"]

    @staticmethod
    def __generate_id() -> int:
        return random.randint(0, 2147483647)

    @staticmethod
    def __remove_quotes_and_strip(el: str) -> str:
        return el.replace("\"", "").replace("'", "").strip()

    def __search_date_time(self) -> str or None:
        for r in self.raw_data:
            if re.match(ConstantExpressions.REGEX_DAYTIME, r):
                return r
        return "TBD"

    def __search_event_date(self) -> str or None:
        if re.match(ConstantExpressions.REGEX_DATE, self.raw_data[0]):
            return self.raw_data[0]
        else:
            if self.raw_data[0].count("/") == 2 and re.match(ConstantExpressions.REGEX_NUMERIC_DATE,
                                                             self.raw_data[0]):
                return self.__convert_date_format(self.raw_data[0])
        return None

    def __search_opponent(self) -> str:
        try:
            vs_index = self.raw_data.index('vs.')
            if vs_index > 1:
                return self.__remove_quotes_and_strip(self.raw_data[vs_index + 1])
        except ValueError:
            return self.__remove_quotes_and_strip(self.raw_data[3])

    def __was_game_cancelled(self) -> bool:
        if "CANCELLED" in self.raw_data:
            return True
        return False

    def __search_win_data(self) -> bool or None:
        if "Win" in self.raw_data:
            return True
        elif "Loss" in self.raw_data:
            return False
        else:
            return None

    def __search_if_away(self) -> str or None:
        if "Away" in self.raw_data or "away" in self.raw_data:
            return 0
        elif "Home" in self.raw_data or "home" in self.raw_data:
            return 1
        else:
            return 2

    def __search_score_data(self) -> str or None:
        for el in self.raw_data:
            if re.match(ConstantExpressions.REGEX_GAME_SCORE, el) and self.raw_data.index(el) > 4:
                return el
        return None

    def __decide_game_location(self) -> str or None:
        if self.is_away == 0 and self.opponent_school is not None:
            return self.opponent_school
        elif self.is_away == 1:
            return "Avon Old Farms"
        else:
            return "N/A"

    def __make_none_data_null(self):
        if self.sport is None:
            self.sport = "NULL"
        if self.category is None:
            self.category = "NULL"
        if self.game_location is None:
            self.game_location = "NULL"
        if self.opponent_school is None:
            self.opponent_school = "NULL"
        if self.is_away is None:
            self.is_away = "NULL"
        if self.game_score is None:
            self.game_score = "NULL"
        if self.unified_data is None:
            self.unified_data = "NULL"

    def get_tuple(self) -> tuple:
        self.__make_none_data_null()
        return str(self.__generate_id()), \
            self.sport, \
            str(self.category), \
            str(self.game_location), \
            self.opponent_school, \
            str(self.is_away), \
            str(self.game_score), \
            "NULL", \
            str(self.__alphabetical_date_format_to_numeric(self.unified_data))

    def __str__(self) -> str:
        self.__make_none_data_null()
        return "GameEvent[id=" + str(self.__generate_id()) + \
            ", sport_name=" + str(self.sport) + \
            ", sport_category=" + str(self.category) + \
            ", game_location=" + self.game_location + \
            ", opponent_school=" + self.opponent_school + \
            ", is_away=" + self.is_away + \
            ", match_result=" + self.game_score + \
            ", coach_comment=NULL" + \
            ", game_schedule=" + self.__alphabetical_date_format_to_numeric(self.unified_data) + "]"
