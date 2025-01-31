# Ancient Checkpoints

The checkpoint loading part of the system is different than most parts of the
system because it specifically needs to be built to handle the outputs of older
versions of the system.

We have some old checkpoints laying around here to ensure that we keep working
with old checkpoints in the future.

- `0.12.3-keras`: a stripped-down checkpoint, only useful for warmstarting
- `0.13.7-keras`: a stripped-down checkpoint, only useful for warmstarting
- `0.13.8-keras`: a stripped-down checkpoint, only useful for warmstarting
- `0.17.6-estimator`: fetched by Checkpoint.download() to populate metadata.json
- `0.17.6-keras`: fetched by Checkpoint.download() to populate metadata.json
- `0.17.6-pytorch`: fetched by Checkpoint.download() to populate metadata.json
- `0.17.7-estimator`: fetch by direct access to checkpoint files
- `0.17.7-keras`: fetch by direct access to checkpoint files
- `0.17.7-pytorch`: fetch by direct access to checkpoint files
