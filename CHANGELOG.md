# Changelog

## [1.1.1](https://github.com/mjmorales/doan/compare/v1.1.0...v1.1.1) (2023-01-08)


### Bug Fixes

* :bug: added missing agent config to init call ([3de7298](https://github.com/mjmorales/doan/commit/3de729890a986916bfe4dece015acdf81d93ad4a))

## [1.1.0](https://github.com/mjmorales/doan/compare/v1.0.1...v1.1.0) (2023-01-08)


### Features

* :sparkles: --init now deploys ansible repo and starts a single agent run ([8f5c53a](https://github.com/mjmorales/doan/commit/8f5c53a5408b201f59331b24d68b950c7b7eb2e3))

## [1.0.1](https://github.com/mjmorales/doan/compare/v1.0.0...v1.0.1) (2023-01-08)


### Bug Fixes

* :bug: VERSION now pulls from git tag instead of release json ([e575e8e](https://github.com/mjmorales/doan/commit/e575e8e1b76cdc3a07d9e67bba9df3ff6f5a5a41))

## 1.0.0 (2023-01-08)


### Features

* :tada: migrated project from previous repository due to release-please bug ([6308ee7](https://github.com/mjmorales/doan/commit/6308ee71efedd4226737c8bca4faee9e3d4ba5a6))


### Bug Fixes

* **ghactions:** :bug: update release-type to go instead of node ([babf0ea](https://github.com/mjmorales/doan/commit/babf0eac9078c038e1d67cfc3557a77e085d8b17))
* **ghactions:** :bug: update triggering branch to master instead of main ([3c5e6c8](https://github.com/mjmorales/doan/commit/3c5e6c820a84921b3b78decda5d186c9709ed771))

## 1.0.0 (2023-01-08)


### Features

* :beers: started ansible agent code ([d14954f](https://github.com/mjmorales/doan/commit/d14954fbf498c25519d79814add3141699e414f5))
* :sparkles: added check to DeployRepo for matching md5 sums of bundles ([ec6d2df](https://github.com/mjmorales/doan/commit/ec6d2dfbd8dcd5d31aa92880c70e60f2cac274c1))
* :sparkles: added Daemon mode support for agent ([ddb4cb3](https://github.com/mjmorales/doan/commit/ddb4cb302978d87f5c8042f1f340b0633433a183))
* :sparkles: added support for reading from yaml config ([c999bb0](https://github.com/mjmorales/doan/commit/c999bb05f9706d7e94da116e7f95e78fe5fb6280))
* :sparkles: added symlinking for active directory ([21f4754](https://github.com/mjmorales/doan/commit/21f4754d84bea90782b220a35a3054a95516bdac))
* :sparkles: fixed tar handling ([d292171](https://github.com/mjmorales/doan/commit/d2921711a0a3d9aaf1e2dd236c2be8e9dbb964b9))
* :tada: initial commit ([e0bd91e](https://github.com/mjmorales/doan/commit/e0bd91e6d05543a5d00df8f15ca7f57597360efc))
