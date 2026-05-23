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
taq --help,  -h   # show help
taq --version,-v  # show version
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

Paths support `$HOME` and other environment variable expansion.

`TAQ_ANSIBLE_INVS` expects directories, not files. taq walks each directory for `.yaml`/`.yml` files (skipping `group_vars`, `host_vars`, and other subdirectories). Multiple directories are separated by `;`:

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

- Multiple inventory sources ?
- Nested groups ?
- Remote inventory (URL, repository) ?
