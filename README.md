y360c is a command line interface (cli) for Yandex 360, which uses the standard API to access the corresponding entities/actions in Yandex360.

## Install and PATH

From this directory:

```
go install .
```

The command writes the `y360c` binary to `$(go env GOPATH)/bin`. That directory is `$HOME/go/bin` when `GOPATH` is `$HOME/go`.

Add it to `PATH` in `~/.zshrc`:

```
export PATH="$HOME/go/bin:$PATH"
```

Open a new terminal. `y360c --help` then runs from any directory.

## Credentials

Each run of the `y360c` binary reads `Y360_TOKEN` and `Y360_ORG_ID` from the process environment. Neither value is read from `.env` or stored in `y360c.json`.

```
export Y360_TOKEN='...'
export Y360_ORG_ID='123456'
```

`user add` generates a random initial password. Pass `--password` when the password should be chosen in advance. The created employee must change it at the first sign-in. The password is not read from the environment.

```
y360c user add --nickname person --department-id 42 --first-name Person --last-name Example --position Engineer --output creds.csv
y360c user set --email person@example.com --position Engineer --department-id 42
```

`--output` appends one CSV row (`nickname`, `email`, `user_id`, `password`). Repeating the command with the same file keeps several accounts in one table. The same row is printed in the terminal.

Copy `.env.example` to `.env` only when you need to point the process at a local stand-in:

```
Y360_API_BASE=
```

Leave `Y360_API_BASE` empty. `api360.yandex.net` is the legacy host: it is still fully supported, and employee creation plus a full profile update remain there. `cloud-api.yandex.net` is the current host for new methods, including the department list and user alias add/remove. It does not replace the legacy create and update methods. `--org-id` overrides `Y360_ORG_ID` for one command. `--token` overrides `Y360_TOKEN` for one command.

Create an OAuth application at https://oauth.yandex.ru/ and obtain a token with the scopes required by the command you run.

Examples:

**List of organizations**:

`y360c org ls` (with `Y360_TOKEN` exported)

**List of departments**:

`y360c dept ls` (with `Y360_TOKEN` and `Y360_ORG_ID` exported)

or

`y360c dept ls --org-id <organization id>`

**List of employees**:

`y360c user ls` (with `Y360_TOKEN` and `Y360_ORG_ID` exported)

or

`y360c user ls --org-id <organization id>`

or (export to csv-file)

`y360c user ls --org-id <organization id> --csv`

or (search for a specific employee)

`y360c user ls --org-id <organization id> --id <employee id>`

or (search by email)

`y360c user ls --org-id <organization id> --email aaa@vvv.cc`

or (search by first/middle/last name)

`y360c user ls --org-id <organization id> --name vladimir`
