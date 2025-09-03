# gokeenapi 🚀

<p align="center">
  <img src="https://github.com/user-attachments/assets/404e89cc-4675-42c4-ae93-4a0955b06348" alt="gokeenapi-video" width="100%">
</p>

**A powerful and easy-to-use command-line utility to manage your Keenetic router via REST API.**

Tired of clicking through web interfaces? `gokeenapi` lets you automate common networking tasks, saving you time and effort.

---

### ✨ Why use gokeenapi?

*   **💻 Automate Everything:** Manage routes, DNS records, and WireGuard connections with simple commands.
*   **⚙️ No Router Configuration:** No complex setup is needed on your router. Just provide the address in a `yaml` config file.
*   **🌐 LAN or Internet:** Works seamlessly whether you're on the same network as your router or accessing it from the internet via KeenDNS.
*   **🎯 Precise Control:** Unlike the web interface, `gokeenapi` allows you to delete static routes for a *specific* interface without affecting others.

---

###  GUI Version Available! 🎨

Not a fan of the command line? No problem! A user-friendly GUI version of `gokeenapi` is available [here](https://github.com/Noksa/gokeenapiui).

---

### 🚀 Getting Started

The easiest way to get started is by using Docker or by downloading the latest release.

#### Docker (Recommended)

Using Docker is the recommended way to run `gokeenapi`.

1.  **Pull the Docker image:**
    ```shell
    export GOKEENAPI_IMAGE="noksa/gokeenapi:stable"
    docker pull "${GOKEENAPI_IMAGE}"
    ```

2.  **Run a command:**
    ```shell
    docker run --rm -ti -v "$(pwd)/config_example.yaml":/gokeenapi/config.yaml \
      "${GOKEENAPI_IMAGE}" show-interfaces --config /gokeenapi/config.yaml
    ```

#### Latest Release

You can find the latest release [here](https://github.com/Noksa/gokeenapi/releases).

---

### 🔧 Configuration

`gokeenapi` is configured using a `yaml` file. You can find an example [here](https://github.com/Noksa/gokeenapi/blob/main/config_example.yaml).

To use your configuration file, pass the `--config <path>` flag with your command.

---

### 🎬 Video Demos

Check out these video demonstrations (in Russian) to see `gokeenapi` in action:

*   [Routes Management](https://www.youtube.com/watch?v=lKX74btFypY)

---

### 📚 Commands

Here are some of the things you can do with `gokeenapi`. For a full list of commands and options, use the `--help` flag.

```shell
./gokeenapi --help
```

#### `show-interfaces`

*Aliases: `showinterfaces`, `showifaces`, `si`*

Displays all available interfaces on your Keenetic router.

```shell
# Show all interfaces
./gokeenapi show-interfaces --config my_config.yaml

# Show only WireGuard interfaces
./gokeenapi show-interfaces --config my_config.yaml --type Wireguard
```

#### `add-routes`

*Aliases: `addroutes`, `ar`*

Adds static routes to your router.

```shell
./gokeenapi add-routes --config my_config.yaml
```

#### `delete-routes`

*Aliases: `deleteroutes`, `dr`*

Deletes static routes for a specific interface.

```shell
# Delete routes for all interfaces in the config file
./gokeenapi delete-routes --config my_config.yaml

# Delete routes for a specific interface
./gokeenapi delete-routes --config my_config.yaml --interface-id <your-interface-id>
```

#### `add-dns-records`

*Aliases: `adddnsrecords`, `adr`*

Adds static DNS records.

```shell
./gokeenapi add-dns-records --config my_config.yaml
```

#### `delete-dns-records`

*Aliases: `deletednsrecords`, `ddr`*

Deletes static DNS records based on your configuration file.

```shell
./gokeenapi delete-dns-records --config my_config.yaml
```

#### `add-awg`

*Aliases: `addawg`, `aawg`*

Adds a new WireGuard connection from a `.conf` file.

```shell
./gokeenapi add-awg --config my_config.yaml --conf-file <path-to-conf> --name MySuperInterface
```

#### `update-awg`

*Aliases: `updateawg`, `uawg`*

Updates an existing WireGuard connection from a `.conf` file.

```shell
./gokeenapi update-awg --config my_config.yaml --conf-file <path-to-conf> --interface-id <interface-id>
```

---

### 🤝 Contributing

Contributions are welcome! If you have any ideas, suggestions, or bug reports, please open an issue or create a pull request.

---

### 📄 License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
