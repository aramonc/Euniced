# RMQ Point Haired Boss

The generic RabbitMQ worker manager.

## Usage

```
$> sudo phbd /etc/phb.conf
```

The only argument to the command is the location of the configuration file which defaults to `/etc/phb.conf`. If the
configuration file is not found from the argument or the default location then it exits with a code of `2`. The process
needs to run with root permissions to check number of cores, load averages, and other system metrics.

## The configuration

The configuration is written in YAML format.

### Sample

```
cmd: /var/run/worker option1 option2
min_workers: 1
max_workers: 5
max_load_override: 4
max_mem_override: 512M
idle_time: 30
rmq_host: http://192.168.100.25
rmq_port: 4545
rmq_user: guest
rmq_password: guest
rmq_exchange: data_exchange
```
