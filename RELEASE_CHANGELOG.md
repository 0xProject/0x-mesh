- [Docker image](https://hub.docker.com/r/0xorg/mesh/tags)
- [README](https://github.com/0xProject/0x-mesh/blob/v9.0.1/README.md)

## Summary

### Bug fixes 🐞

- Fix bug where we weren't enforcing that we never store more than `miniHeaderRetentionLimit` block headers in the DB. This caused [issue #667](https://github.com/0xProject/0x-mesh/issues/667) and also caused the Mesh node's DB storage to continuously grow over time. ([#716](https://github.com/0xProject/0x-mesh/pull/716))


