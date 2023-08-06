import functools
import inspect
import logging
import time
import traceback


def log_execution(func):
    @functools.wraps(func)
    def wrapper(*args, **kwargs):
        start_time = time.time()
        class_name = None
        try:
            if args and inspect.ismethod(func):
                class_name = args[0].__class__.__name__
            class_name = str(class_name + '.') if class_name else ''

            result = func(*args, **kwargs)
            elapsed_time = time.time() - start_time
            logging.info(f"[{time.asctime()}] {class_name}{func.__name__} "
                         f"was executed successfully in {elapsed_time:.6f}s.")
            return result

        except Exception:
            exc_traceback = traceback.format_exc()
            max_line_length = max(len(line) for line in exc_traceback.splitlines())
            separator = '-' * max_line_length
            logging.error(f"[{time.asctime()}] {class_name}{func.__name__} "
                          f"threw an exception:\n{separator}\n{traceback.format_exc()}{separator}")

    return wrapper


def log_class_methods(cls):
    for attr_name, attr_value in cls.__dict__.items():
        if callable(attr_value):
            setattr(cls, attr_name, log_execution(attr_value))
    return cls


@log_execution
def setup_logging(filename: str = "log"):
    if not filename.endswith(".log"):
        filename += ".log"
    logging.basicConfig(filename=filename, level=logging.INFO, format='%(message)s')
    logging.info("START".center(130, "="))
