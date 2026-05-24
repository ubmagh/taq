# taq

Quickly find and SSH into any host in your inventory — no more digging through command history.

## How it works

1. Parses your inventory file (or Ansible inventories)
2. Launches an interactive fuzzy search over all hosts
3. Prompts for the SSH username (defaults to the host's configured user)
4. Opens the SSH session

## Installation

**From a release binary** — download the latest binary for your platform from the [Releases](../../releases) page.

**With Go:**
```sh
go install github.com/ubmagh/taq@latest
```

**From source:**
```sh
git clone https://github.com/ubmagh/taq
cd taq
make install
```

## Usage

```
taq                          # launch interactive SSH search
taq -l, --local-forward      # launch in local port-forward mode (-L)
taq -r, --remote-forward     # launch in remote/reverse port-forward mode (-R)
taq --list [query]           # list all hosts (or filter by query), then exit
taq --list [query] -o fmt    # same with output format: table (default), json, yaml, plain
taq --validate               # parse inventories, report host count, then exit
taq --debug,   -d            # enable verbose output (combine with any flag)
taq --version, -v            # show version
taq --help,    -h            # show help
```

**Keybindings:**
```
↑/↓       navigate the list
Enter     select host / confirm
Tab       toggle compact / detailed view
Esc       back / exit
Ctrl+C    exit
```

## Configuration

| Variable | Default | Description |
|---|---|---|
| `TAQ_INVENTORY_PATH` | `$HOME/.config/taq/inventory.yaml` | Path to inventory file |
| `TAQ_DEFAULT_USER` | `$USER` | Default SSH username |
| `TAQ_DEFAULT_SSH_KEY_PATH` | _(none)_ | Default SSH key path |
| `TAQ_ANSIBLE_INVS` | _(none)_ | Semicolon-separated list of Ansible project inventory **directories** |
| `TAQ_DISPLAY_MODE` | `detailed` | List display mode: `detailed` or `compact` |
| `TAQ_SSH_TIMEOUT` | _(none)_ | SSH connect timeout in seconds (e.g. `5`) |
| `TAQ_DEBUG` | _(none)_ | Set to any value to enable verbose/debug output |

Paths support `$HOME` and other environment variable expansion.

`TAQ_ANSIBLE_INVS` expects the **inventory root** directory (e.g. `inventories/`), not the whole project root. taq recursively walks each directory for `.yaml`/`.yml`/`hosts` yaml files, skipping known non-inventory dirs (`group_vars`, `host_vars`, `roles`, etc.). Multiple paths are separated by `;`:

```sh
export TAQ_ANSIBLE_INVS="~/projects/infra/inventory;~/projects/app/inventory"
```

## Inventory file

```yaml
hosts:
  - name: my-server
    address: 192.168.1.10
    user: ubuntu
    port: 22
    key_path: ~/.ssh/id_rsa

groups:
  production:
    labels:
      env: prod
    hosts:
      - name: web-01
        address: 10.0.0.1
        user: deploy
```

See `example-inventories/` for more examples.

## Non-interactive mode

`--list` prints hosts without launching the TUI — useful for scripting, auditing, or piping into other tools.

```sh
taq --list                        # all hosts, table format
taq --list web                    # fuzzy-filter "web", table format
taq --list -o json                # all hosts as JSON
taq --list -o json prod           # filter "prod", output JSON
taq --list -o yaml                # all hosts as YAML
taq --list -o plain               # name + address, one per line
```

**Output formats:**

| Format | Description |
|--------|-------------|
| `table` | Aligned columns: NAME, ADDRESS, USER, PORT, GROUPS *(default)* |
| `json`  | JSON array of host objects |
| `yaml`  | YAML sequence of host objects |
| `plain` | `name address` one per line — pipe-friendly |

**Examples:**

```sh
# Count all hosts
taq --list -o plain | wc -l

# Get all hosts in a group, pipe to fzf
taq --list -o plain prod | fzf

# Extract addresses with jq
taq --list -o json | jq -r '.[].address'

# Find all hosts on a specific subnet
taq --list -o json | jq '.[] | select(.address | startswith("10.0.1."))'
```

## Port Forwarding

Use `-l` / `--local-forward` or `-r` / `--remote-forward` to launch taq in port-forwarding mode. The host list and search work exactly the same — the flag just changes what happens after you pick a host.

### Flow

1. Launch with `taq -l` or `taq -r`
2. Search and navigate to a host with `↑/↓`, press `Enter`
3. Confirm the SSH username
4. Add one or more forwarding rules — press `Enter` on an empty line when done
5. The terminal blocks while the tunnel is active — `Ctrl+C` to stop

### Rule format

Use the shorthand `localPort->remotePort` — taq assumes `localhost` on the remote side:

| You type | SSH arg produced |
|---|---|
| `8080->3000` | `-L 8080:localhost:3000` or `-R 8080:localhost:3000` |
| `5432->5432` | `-L 5432:localhost:5432` |

You can also type the full SSH spec (`8080:somehost:3000`) if you need a non-localhost target.

### Local vs Remote

| Key | Flag | Direction | Typical use |
|-----|------|-----------|-------------|
| `l` | `-L` | local → remote | Reach a service on the server from your machine |
| `r` | `-R` | remote → local | Expose a local service to the remote server |

## Todos

- Multiple taq-inventory sources ?
- Nested groups on taq-inventory ?
- Remote inventory (URL, repository) ?
