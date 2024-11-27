## log-checker

Zetta Log Checker/Analyzer

### Synopsis

Application that performs the collection and verification of the Zetta Station Logs
This CLI application can also collect (snapshot) the station logs and run a HTTP simulator server for development and testing

### Options

```
      --config string           config file (default is $HOME/.log-checker.yaml)
  -t, --date string             Optional date to query Zetta Logs (YYYY-MM-DD)
  -h, --help                    help for log-checker
  -i, --host string             Zetta Server Host (default "localhost")
  -k, --key string              Zetta API Key
  -n, --port int                Zetta API Server Port (default 3139)
  -p, --zetta_password string   Zetta Password (default "admin")
  -u, --zetta_username string   Zetta Username (default "admin")
```

### SEE ALSO

-   [log-checker completion](docs/log-checker_completion.md) - Generate the autocompletion script for the specified shell
-   [log-checker run](docs/log-checker_run.md) - The run command executes the zetta log check analyzer
-   [log-checker server](docs/log-checker_server.md) - Simulated Simple Zetta API Server for development and testing
-   [log-checker snapshot](docs/log-checker_snapshot.md) - Fetch the station logs and save them to the disk to simulate for later

###### Auto generated by spf13/cobra on 27-Nov-2024
