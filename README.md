[![Actions Status](https://github.com/pavbis/repositories-api/actions/workflows/verify-pull-request.yml/badge.svg)](https://github.com/pavbis/repositories-api/actions)

## Requirements

* Go >= 1.14

## Installation

```bash
make init
```
This will bring up all services, install all dependencies, compile the go binary and seed the database.

After this check if the route http://localhost:7000/api/health is available.