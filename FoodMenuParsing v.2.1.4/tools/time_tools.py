from .logging_tools import log_execution
from datetime import datetime, timedelta
from typing import Tuple


@log_execution
def get_current_date() -> Tuple[str, str, str]:
    now = datetime.now()
    return now.strftime("%Y"), now.strftime("%m"), now.strftime("%d")


@log_execution
def get_next_week_date(date_tuple) -> Tuple[str, str, str]:
    current_date = datetime(int(date_tuple[0]), int(date_tuple[1]), int(date_tuple[2]))
    next_week_date = current_date + timedelta(weeks=1)
    return next_week_date.strftime("%Y"), next_week_date.strftime("%m"), next_week_date.strftime("%d")
