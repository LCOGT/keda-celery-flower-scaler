from time import sleep

from celery import Celery


app = Celery("tasks")


@app.task
def work(trace_id: str, i: int, seconds: int):
    sleep(seconds)
    return trace_id, i
