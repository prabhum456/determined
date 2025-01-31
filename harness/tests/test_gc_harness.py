import os
from pathlib import Path
from typing import Any, List

import pytest

from determined.common import storage
from determined.exec.gc_checkpoints import delete_checkpoints
from tests.storage import util as storage_util


@pytest.fixture()
def manager(tmp_path: Path) -> storage.StorageManager:
    return storage.SharedFSStorageManager(str(tmp_path))


@pytest.fixture(params=[0, 1, 5])
def to_delete(request: Any, manager: storage.StorageManager) -> List[str]:
    storage_ids = []
    for _ in range(request.param):
        with manager.store_path() as (storage_id, path):
            storage_util.create_checkpoint(path)
            storage_ids.append(storage_id)

    assert len(os.listdir(manager._base_path)) == request.param
    return storage_ids


def test_delete_checkpoints(manager: storage.StorageManager, to_delete: List[str]) -> None:
    delete_checkpoints(manager, to_delete, dry_run=False)
    assert len(os.listdir(manager._base_path)) == 0


def test_dry_run(manager: storage.StorageManager, to_delete: List[str]) -> None:
    delete_checkpoints(manager, to_delete, dry_run=True)
    assert len(os.listdir(manager._base_path)) == len(to_delete)
