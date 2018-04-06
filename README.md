![BPool](https://i.imgur.com/NZzf9il.png)

### Ride sharing made easy for UCLA

## Application Wiki
https://github.com/ucladevx/BPool/wiki


## Dependencies

1. Golang 1.9 or higher
2. Postgres
3. Go dep
4. GPG

## How to run

1. Run dep ensure
```bash
dep ensure
```

2. Unencrypt the config file
```bash
make config
```

3. Run the server
```bash
make run
```