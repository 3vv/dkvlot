# RCS

RCS - Redis server/Cluster Synchronizer

Cmds:
  - -keys=false Retrieve [keys] from source host.
  - -sync=false [Sync]hronize keys specified by -l/-f from source host to destination host.

Opts:
  - -l=""  The [l]ist of keys specified by a file lines
  - -f="*" The [f]ilter of keys specified by a pattern
  - -s=""  The [s]ource list of host:port server/cluster (Ex. 127.0.0.1:6379)
  - -d=""  The [d]estination list of host:port server/cluster (Ex. 127.0.0.1:7000,127.0.0.1:7001,127.0.0.1:7002)

### Download and Install

#### Use Binary Distributions

1.Install docker on your operating system

2.Execute command
  - `docker pull dlot/dkvlot_rcs:latest`

3.Locate your redis source `<such as $HOST1:$PORT1,$HOST2:$PORT2,$HOST3:$PORT3,$HOST4:$PORT4>` & destination `<ex. 127.0.0.1:7000,127.0.0.1:7001,127.0.0.1:7002>` (if you need) server/cluster

4.Organize your commands & options to replace the following "$@"

5.Execute command
  - `docker run --rm dlot/dkvlot_rcs:latest /app/dkvlot/rcs "$@"`

6.(Optional) Save the above script as an executable file `rcs`
  - `./rcs $CMDs $OPTs`

#### Install From Source

1.Select your working directory (as environment variable D)

2.Execute command
  - `git clone $GIT_HOST/$USR_NAME/dkvlot.rcs.git "$D/dkvlot.rcs"`

5.Execute command
  - `cd "$D/dkvlot.rcs" && export GOPATH=$GOPATH:.`

6.Follow your Go knowledge & your wants...

### Advanced usage

#### Times task

Organize your task

#### Timing task

Organize your task

#### One-way server/cluster synchronization

Organize your synchronization

#### Peer-to-peer server/cluster synchronization

Organize your synchronization

### Contributing

To contribute, please pull requests of your contributing code to above.

## Let your hot heart go!