# Automated video maker

Project `auvima` - automated video maker app

## Usage:

To begin development:

1. Alter `settings.yaml` file configuration
2. Then just run:

```bash
docker-compose up --build
```

## Production build:

To build lightweight production image under `10 megabytes` just run:

```
docker build -f docker/prod.Dockerfile -t auvima .
```
