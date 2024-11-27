## log-checker completion zsh

Generate the autocompletion script for zsh

### Synopsis

Generate the autocompletion script for the zsh shell.

If shell completion is not already enabled in your environment you will need
to enable it.  You can execute the following once:

	echo "autoload -U compinit; compinit" >> ~/.zshrc

To load completions in your current shell session:

	source <(log-checker completion zsh)

To load completions for every new session, execute once:

#### Linux:

	log-checker completion zsh > "${fpath[1]}/_log-checker"

#### macOS:

	log-checker completion zsh > $(brew --prefix)/share/zsh/site-functions/_log-checker

You will need to start a new shell for this setup to take effect.


```
log-checker completion zsh [flags]
```

### Options

```
  -h, --help              help for zsh
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
      --config string           config file (default is $HOME/.log-checker.yaml)
  -t, --date string             Optional date to query Zetta Logs (YYYY-MM-DD)
  -i, --host string             Zetta Server Host (default "localhost")
  -k, --key string              Zetta API Key
  -n, --port int                Zetta API Server Port (default 3139)
  -p, --zetta_password string   Zetta Password (default "admin")
  -u, --zetta_username string   Zetta Username (default "admin")
```

### SEE ALSO

* [log-checker completion](log-checker_completion.md)	 - Generate the autocompletion script for the specified shell

###### Auto generated by spf13/cobra on 27-Nov-2024