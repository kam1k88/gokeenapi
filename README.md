## A utility to run commands (such as add/delete routes/dns records) in Keenetic routers via REST API

#### Video

https://github.com/user-attachments/assets/404e89cc-4675-42c4-ae93-4a0955b06348

---

#### Version with UI

There is a GUI `gokeenapi` version available [here](https://github.com/Noksa/gokeenapiui)

If you don't like or don't know how to use CLI programs, consider using the GUI version

---

#### Important notes
* No additional configuration is required on a router - just specify the router address in `yaml` config file (for example, `http://192.168.1.1`)
* `gokeenapi` works with Keenetic routers over LAN or Internet using internal router IP address (like `192.168.1.1`) or domain from KeenDNS (like `my-router.keenetic.pro`)
---

#### What `gokeenapi` can already do:
* Add AWG connection from conf files and start them (also wait until they are up and running)
* Apply ASC parameters to existing WG connections from AWG conf files
* Display a list of interfaces that have already been added to the router - for easy search of the interface ID for which you need to add/remove routes
* Delete static routes only for the specified interface. In the Web Configurator of the router, at the moment you can only delete all created static routes for all interfaces at once
* Add\update static routes for the specified interface from bat files from disk
* Add\update static routes for the specified interface from links that download bat file (for example [from here](https://iplist.opencck.org/?format=bat&data=cidr4&site=youtube.com))
* Add/delete static DNS records
---

#### Configuration

`gokeenapi` should be configured using `yaml` config file ([example](https://github.com/Noksa/gokeenapi/blob/main/config_example.yaml))

---

#### Videos

Videos with examples (note that language is **Russian**:
* [Routes](https://www.youtube.com/watch?v=lKX74btFypY)

---

#### Examples

The easiest way to start using `gokeenapi` is through docker containers or using the latest available release from [here](https://github.com/Noksa/gokeenapi/releases)

---

#### Docker

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

#### Binary

```shell
# config can be specified via GOKEENAPI_CONFIG environment variable instead of flag
./gokeenapi --config my_config.yaml show-interfaces
./gokeenapi --config my_config.yaml add-dns-records
./gokeenapi --config my_config.yaml delete-dns-records
./gokeenapi --config my_config.yaml add-routes
./gokeenapi --config my_config.yaml delete-routes --interface-id <iface-id-to-delete-routes-on>
./gokeenapi --config my_config.yaml add-awg --conf-file <path-to-conf-file> --name MySuperInterface
./gokeenapi --config my_config.yaml configure-awg --conf-file <path-to-conf-file> --interface-id <iface-id-to-configure>

```

---
