## A utility to run commands (such as add/delete routes/dns records) in Keenetic routers via REST API

### Video

https://github.com/user-attachments/assets/404e89cc-4675-42c4-ae93-4a0955b06348

---

### Version with UI

There is a GUI `gokeenapi` version available [here](https://github.com/Noksa/gokeenapiui)

If you don't like or don't know how to use CLI programs, consider using the GUI version

---

### Important notes
* No additional configuration is required on a router - just specify the router address in `yaml` config file (for example, `http://192.168.1.1`)
* `gokeenapi` works with Keenetic routers over LAN or Internet using internal router IP address (like `192.168.1.1`) or domain from KeenDNS (like `my-router.keenetic.pro`)
---

### What `gokeenapi` can already do:
* Add AWG connection from conf files and start them (also wait until they are up and running)
* Apply ASC parameters to existing WG connections from AWG conf files
* Display a list of interfaces that have already been added to the router - for easy search of the interface ID for which you need to add/remove routes
* Delete static routes only for the specified interface. In the Web Configurator of the router, at the moment you can only delete all created static routes for all interfaces at once
* Add\update static routes for the specified interface from bat files from disk
* Add\update static routes for the specified interface from links that download bat file (for example [from here](https://iplist.opencck.org/?format=bat&data=cidr4&site=youtube.com))
* Add/delete static DNS records
---

### Configuration

`gokeenapi` should be configured using `yaml` config file ([example](https://github.com/Noksa/gokeenapi/blob/main/config_example.yaml))

Use `--config <path>` flag to pass config file to the utility

---

### Videos

Videos with examples (note that language is **Russian**:
* [Routes](https://www.youtube.com/watch?v=lKX74btFypY)

---

### Examples

The easiest way to start using `gokeenapi` is through docker containers or using the latest available release from [here](https://github.com/Noksa/gokeenapi/releases)

---

### Docker

It is recommended to use `noksa/gokeenapi:stable` image

* Check all existing commands
```shell
export GOKEENAPI_IMAGE="noksa/gokeenapi:stable"
docker pull "${GOKEENAPI_IMAGE}"
docker run --rm -ti "${GOKEENAPI_IMAGE}" --help
```

* Show interfaces on the router
```shell
export GOKEENAPI_IMAGE="noksa/gokeenapi:stable"
docker run --rm -ti -v "$(pwd)/config_example.yaml":"/gokeenapi/config.yaml" "${GOKEENAPI_IMAGE}" show-interfaces --config "/gokeenapi/config.yaml"
```

### Commands

#### Help

To see all available commands/subcommands/flags use `--help` flag

Examples:
```shell
./gokeenapi --help
./gokeenapi show-interfaces --help
./gokeenapi add-routes --help
```

#### Show interfaces

To see all interfaces available in a keenetic router, use `show-interfaces` command

It is possible to limit which interface types should be in output using `--type` flag

```shell
# all interfaces including internal ones
./gokeenapi show-interfaces --config my_config.yaml
# show only wireguard interfaces
./gokeenapi show-interfaces --config my_config.yaml --type Wireguard
```

#### Add routes
To add static routes use `add-routes|addroutes|ar` command

The `add-routes` command uses the yaml config file to determine which routes to which interface should be added (check [config_examle.yaml](https://github.com/Noksa/gokeenapi/blob/main/config_example.yaml))

Once the config file is ready, run the following command to add routes

```shell
./gokeenapi add-routes --config my_config.yaml
```

#### Delete routes

To delete added static routes use `delete-routes|deleteroutes|dr` command

By default, routes are deleted only for interfaces which are specified in `routes` field in yaml config

The `--interface-id` flag is optional and can be used to explicitly specify the interface id for which routes should be deleted instead of using ids from yaml config

**The routes deletion is only done for the specified interfaces - all other routes stay**.


```shell
# delete routes for all interface ids which are described in the yaml config
./gokeenapi delete-routes --config my_config.yaml
# delete routes only for specified interface via flag --interface-id
# interface ids from yaml config are not used and ignored
./gokeenapi delete-routes --config my_config.yaml --interface-id <your-interface-id>
```

#### Add DNS records

To add static DNS records use `add-dns-records|adddnsrecords|adr` command

The `add-dns-records` command uses the yaml config file to determine which DNS records should be added (check [config_examle.yaml](https://github.com/Noksa/gokeenapi/blob/main/config_example.yaml))

Once the config file is ready, run the following command to add DNS records

```shell
./gokeenapi add-dns-records --config my_config.yaml
```

#### Delete DNS records

To delete added static routes use `delete-routes|deleteroutes|dr` command

The `--interface-id` flag is required to pass to the command

**The deletion is only done for the specified interface**.

All other routes which relate to another interfaces will stay


```shell
./gokeenapi delete-dns-records --config my_config.yaml --interface-id <your-interface-id>
```

#### Add new AWG connection from conf file

To add new AWG connection from conf file, use `add-awg|addawg|aawg` command

The command uses AWG configuration file to add a new WG connection in a keenetic router, configure and run it.

The command works as follows:

* Adds new WG connection
* Checks if ASC parameters should be added/updated in the created connection from conf file and updates if needed (before `4.3.6` keenetic routes didn't add ASC parameters automatically)
* Moves the interface to `up` state
* Waits until the interface is `up` and `running` which means it is ready to use 

`--conf-file` is required flag, path to the conf file from which connection should be added

`--name` is optional flag, name for new WG connection

```shell
./gokeenapi add-awg --config my_config.yaml --conf-file <path-to-conf> --name MySuperInterface
```

#### Update existing AWG connection from conf file

To update existing AWG connection from conf file, use `update-awg|updateawg|uawg` command

The command uses AWG configuration file to update an existing WG connection in a keenetic router, reconfigure and run it.

The command works as follows:

* Finds required interface-id in the router
* Checks if ASC parameters should be added/updated in the connection from conf file and updates if needed (before `4.3.6` keenetic routes didn't add ASC parameters automatically)
* Moves the interface to `up` state
* Waits until the interface is `up` and `running` which means it is ready to use

`--conf-file` is required flag, path to the conf file from which connection should be added

`--interface-id` is required flag, id of an interface to update

```shell
./gokeenapi update-awg --config my_config.yaml --conf-file <path-to-conf> --interface-id <interface-id>
```

---
