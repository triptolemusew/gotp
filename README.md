## Gotp


## Requirements
- [wl-clipboard]

<!-- Add the key -->
gotp add --name <name> --secret <secret>

<!-- Remove the key -->
gotp remove --name <name>

<!-- Sync with the path -->
gotp sync --path [OPTIONAL]

<!-- List all of the available keys -->
gotp list

<!-- Initialize the project, creating the directory, and the sqlite tables -->
gotp init

## TODO
[-] Remove all commands, make it more simple.
[-] Remove db dependencies, make all generated passcode on the fly
[-] Make the filters more intuitive
[] Default folder will be `.gotp` however, user can supply --path to override
[] Package the project

