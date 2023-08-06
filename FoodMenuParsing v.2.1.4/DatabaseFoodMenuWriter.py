import sqlite3
from tools.logging_tools import log_class_methods


@log_class_methods
class DatabaseFoodMenuWriter:
    _instance = None
    __slots__ = ["_sqlite_connection", "_cursor", "_database"]

    def __new__(cls, *args, **kwargs):
        if not cls._instance:
            cls._instance = super(DatabaseFoodMenuWriter, cls).__new__(cls)
        return cls._instance

    def __init__(self, database_file: str = "../database.db"):
        self._sqlite_connection = None
        self._cursor = None
        self._database = database_file

    def __enter__(self):
        self._connect()
        return self

    def __exit__(self, exc_type, exc_val, exc_tb):
        self._commit_and_close_connection()
        return True

    def _connect(self):
        if not self._sqlite_connection:
            self._sqlite_connection = sqlite3.connect(self._database)
            self._cursor = self._sqlite_connection.cursor()

    def _commit_and_close_connection(self):
        self._sqlite_connection.commit()
        self._cursor.close()

    def _clear_food_menu(self):
        self._cursor.execute("DELETE FROM FoodMenu")

    def __get_next_id(self):
        self._cursor.execute('SELECT MAX(id) FROM FoodMenu')
        result = self._cursor.fetchone()
        return (result[0] + 1) if result[0] else 1

    def get_next_id(self) -> int:
        """
        :return: First unclaimed ID from FoodMenu table.
        """
        self._connect()
        requested_id = self.__get_next_id()
        self._commit_and_close_connection()
        return requested_id

    def _execute_data_writing(self, data):
        food_menu_insert_query = "INSERT INTO FoodMenu (id, date, breakfast, lunch, dinner) VALUES (?, ?, ?, ?, ?)"
        self._cursor.executemany(food_menu_insert_query, data)

    def execute_food_menu_data_writing(self, data):
        if not self._sqlite_connection:
            message = "Use context manager \"with\" while working with this class."
            raise type("NoDatabaseConnectionException", (Exception,), {})(message)
        self._clear_food_menu()
        self._execute_data_writing(data)
