# Darch

A lightweight app that allows you to build and boot Docker images, bare-metal.

# Getting started

## Install

*NOTE: Only 64-bit Linux with grub is supported.*

**Arch user repositories**

```
pacaur -S darch
```

**Generic Linux**

```
curl -s https://raw.githubusercontent.com/pauldotknopf/darch/master/scripts/install | sudo bash /dev/stdin
```

## Booting

```bash
# Pull the Arch image from docker
docker pull pauldotknopf/darch-base-common:latest
# Extract/prepare the Arch image to a format suitable for
# bare-metal booting.
darch extract pauldotknopf/darch-base-common:latest
# Move the Arch image to the "/boot" directory and update
# the "/boot/grub/grub.cfg" file with the new entries.
sudo darch stage pauldotknopf/darch-base-common:latest
```

Reboot and select your Arch image during boot!

# TODO

More documentation