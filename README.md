# Table Populator

# Introduction

An enterprise-y placeName to planeName,city,state translator,
which effectively reads placeName data from a csv file, extracts the city and state using google maps api,
and creates an output csv file with the placeName, city, state associated together.

- Uses Ports and Adaptors architecture
  - **domain**: Contains domain structures and interfaces (ports). Has no dependencies.
  - **application**: Performs the application logic over the domain. Only depends on the domain.
  - **logger**, **locator**, **dataio**: These implement the interfaces to interact with the outside world,
    i.e. are adapters to the ports found in the domain. Depends on **domain**, **config**.
  - **config**: Reads the configuration information found in `.env` configuration file.
  - **main**: Ties everything together.

# How to Run

- Configure environment
- Build and Run

```
go build -o build/table-populator
./build/table-populator
```

-

## Example `.env` File

```
# API KEYS
MAPS_API_KEY=

# FILE PATHS
CSV_DATA_FILE_PATH=
OUTPUT_FILE_PATH=

# DATAIO
DATAIO_KIND="csv"

# LOCATOR
LOCATOR_KIND="google-maps"

# LOGGER
LOGGER_KIND="multi"
LOG_TO_STDOUT=1
LOG_FILE_PATH=
```

## Example Input Format

```
ignore1,ignore2,my-place-1,ignore3
ignore1,ignore2,my-place-2,ignore3
ignore1,ignore2,my-place-3,ignore3
```

## Example Output Format

```
my-place-1,city-1,state-1
my-place-2,city-2,state-2
my-place-3,city-3,state-3
```

# Todos

- [ ] Add unit tests
- [ ] Integration tests
- [ ] Acceptance tests
