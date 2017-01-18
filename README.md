A utility for slurping Chef definition files and showing/exporting info from those

```
Usage:
  chef-slurper [command]

Available Commands:
  i2c         Exports an i2csshrc content
  inventory   Exports an Ansible inventory file
  nodes       A command to list (and filter) the nodes

Flags:
      --config string                   config file (default is $HOME/.chef-slurper.yaml)
  -n, --nodes string                    Location of your Chef nodes directory (default "./nodes")
      --strip-role-prefix stringArray   Prefixes to remove from role identifiers

Use "chef-slurper [command] --help" for more information about a command.
```

## Installation

### Mac OS

```
curl -sSL -o chef-slurper https://github.com/itzg/chef-slurper/releases/download/v1.0/chef-slurper_darwin_amd64
chmod +x chef-slurper
```