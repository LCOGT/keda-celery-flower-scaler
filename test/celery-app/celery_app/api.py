from typing import TypeAlias
from uuid import uuid4

from fastapi import FastAPI
from pydantic import conint

from .tasks import work

app = FastAPI(
    debug=True,
    docs_url="/",
)


PostitiveInt: TypeAlias = conint(ge=1)


@app.post("/tasks")
def create_tasks(
    trace_id: str | None = None,
    amount: PostitiveInt = 1,
    workSeconds: PostitiveInt = 30,
):
    if trace_id is None:
        trace_id = str(uuid4())

    task_ids: list[str] = []
    for i in range(amount):
        t = work.delay(trace_id, i, workSeconds)
        task_ids.append(t.id)

    return {"task_ids": task_ids}
