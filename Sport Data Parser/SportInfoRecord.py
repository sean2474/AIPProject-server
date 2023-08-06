import random


class SportsInfoRecord:
    def __init__(self, sport_label: str, category: str, season: str, coach_data: dict, roster: list):
        self.id = self.__generate_id()
        self.sport_label = sport_label
        self.category = self.__define_category(category)
        self.season = self.__define_season(season)
        self.coach_names = str(list(coach_data.keys())) if len(coach_data.keys()) > 0 else "NULL"
        self.coach_contacts = str(list(coach_data.values())) if len(coach_data.values()) > 0 else "NULL"
        self.roster = roster

    @staticmethod
    def __generate_id() -> int:
        return random.randint(0, 2147483647)

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
    def __define_season(season: str) -> int:
        season_enum = {
            "Fall": 0,
            "Winter": 1,
            "Spring": 2,
            "NA": 3
        }
        for k in season_enum.keys():
            if k in season:
                return season_enum[k]
        return season_enum["NA"]

    def __str__(self) -> str:
        return (f"SportsInfoRecord[sport_label={self.sport_label}, "
                f"category={self.category}, "
                f"season={self.season}, "
                f"coach_name={self.coach_names}, "
                f"coach_number={self.coach_contacts}, "
                f"roster={str(self.roster)}")

    def get_tuple(self) -> tuple:
        return self.id, self.sport_label, self.category, self.season, self.coach_names, self.coach_contacts, \
            str(self.roster)

