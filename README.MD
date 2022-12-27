# Getpack
**Tool for fetching server jars and modpacks.**

### TODO
- Docker building and testing
- Webhook/Discord error handling
- Store currently installed version of pack
- GUI?

### Features
- bedrock fetching
- paper fetching
- purpur fetching
- spigot compiling
- forge installing
- fabric installing
- curseforge installing
- technicpack installing
- ftb installing

## Config
config.yml layout:
```yaml
target: /path/to/folder
enabled:
  curse: false
  technic: false
  ftb: false
  paper: false
  purpur: false
  forge: false
  fabric: false
  bedrock: false
  spigot: false
curse:
  apikey: insert your curseforge api key here
  modpacks:
    - insert modpack id here
technic:
  modpacks:
    - insert modpack name here
ftb:
  modpacks:
    - insert modpack id here
paper:
  version: insert minecraft version here
purpur:
  version: insert minecraft version here
fabric:
  versions:
    - insert minecraft version here
forge:
  version: insert minecraft version here
spigot:
  version: insert minecraft version here
```

packs.yml layout:
```yaml
packs:
  insert pack id here:
    - insert pack instruction here
    - insert another pack instruction here
```

### Instructions
See practical examples in packs.yml
#### delete
`delete:filename` delete specified file.

`delete:foldername/` delete specified folder.

#### deletereg
`deletereg:filenamereg` parse generic regex format and delete all files that match.

#### rename
`rename:oldname:newname` rename file from **oldname** to **newname**.

#### renamreg
`rename:oldnamereg:newname` parse generic regex format and rename matched file to **newname**.

#### cdreg
`cdreg:foldernamereg` parse generic regex format and cd into found folder.

#### javarun
`javarun:filenamereg:arg` parse generic regex format and run found jar with argument.

#### forgegrep
`forgegrep:filename` read specified file, search for forge version and install it.

#### bashrun
`bashrun:arg` execute specified bash command.

#### shforgegrep
`shforgegrep:filename` read specified file, search for forge >=1.17 startup format and create generic run.sh.

#### forge
`forge:version` install latest forge version for specified **mcversion** in current folder.