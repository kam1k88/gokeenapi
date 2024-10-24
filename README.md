## A utility to add/delete static routes in Keenetic routers via REST API

###### Russian readme (русская инструкция): [README_ru.md](https://github.com/Noksa/gokeenapi/README_ru.md)

---

#### Prerequisites
* To start using it, REST API must be configured first. Refer to this page to find out how to do that
---

#### Features:
* Deletes static routes only with specified interface and not all like Web Panel does
* Adds/updates static routes from bat file

---

#### Examples

* The simplest way to use this tool on PC/Laptop is just run via docker
```shell
# Common way to run commands:
docker run --rm -ti <image_name> -- <command>
```

* Delete all static routes on 