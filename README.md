# taq

Easily find that instance you wanna ssh into, instead of browsing commands history for so long. 

## How it works 

- Parses inventory files, either;
  - A dedicated file `TAQ_INVENTORY_PATH=~/.config/taq.inventory.yaml` following examples on `example-inventories/`.
  - Or Ansible projects yaml inventories, specified on the variable `TAQ_ANSIBLE_INVS`
- Runs a fuzzy search on the parsed items based on input.
- Returns an interactive list of found hosts.
- SSH-es to the selected one.


```
	taq - fast SSH search and connect CLI
	
	Usage:
	taq               # launch interactive search
	taq --help,-h     # show this help message
	taq --version,-v  # show version

	Environment Variables:
	TAQ_DEFAULT_USER   : Specifies default SSH username [$USER]
	TAQ_ANSIBLE_INVS   : List of ansible projects inventories, (;) separated.  []
	TAQ_INVENTORY_PATH : Inventory file path ["~/.config/taq/inventory.yaml"]
```


## Todos 

- default user + can override upon selection
- add ansible inventories parsing
- multiple files/sources ? 
- better design ? 
- nested groups ? 
- support remote inventory (url, repo)
- customized list styling ? 