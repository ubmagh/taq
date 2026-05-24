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
taq               # launch interactive search
taq --help,    -h # show help
taq --version, -v # show version
taq --validate    # parse inventories, report host count, then exit
taq --debug,   -d # enable verbose output (combine with --validate or normal run)
```

**Keybindings:**
```
↑/↓       navigate the list
Enter     select host / confirm username
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

## Todos

- Multiple taq-inventory sources ?
- Nested groups on taq-inventory ?
- Remote inventory (URL, repository) ?
- Interactive SSH port forwarding helper (`-L`/`-R`) ?
