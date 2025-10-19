# taq

Easily find that instance you wanna ssh into, instead of browsing commands history for so long. 

## How it works 

Construct the default inventory file `~/.config/taq.inventory.yaml` following examples on `example-inventories/`.
Then execute the binary `taq` to find quickly the target instance by comma/space separated keywords.
The default inventory file `~/.config/taq.inventory.yaml` can be altered by setting the env variable `TAQ_INVENTORY_PATH`.


## Todos 

- nested groups ? 
- support remote inventory (url, repo)
- multiple files ? 
- customized list styling ? 
