---
subtitle: "Boot multiple stateless/clean environments."
---

# Darch

## What is it?

Think Dockerfiles, but for bootable, stateless, graphical (or not) environments for your everyday usage.

After each reboot, all changes made to your operating system are wiped, except where explicitly configured not to (home directory, game data, etc). Easily switch between images with individual grub menu entries.

Your images can be layered/inherited (like Dockerfiles) using recipes. For example:

* ```base``` -  This can be Arch Linux, Debian, Gentoo, etc.
  * ```common``` - User setup, common tools.
    * ```steam``` - An image that is tweaked for Stream gaming.
    * ```development``` - Your dev tools (make/gcc/etc).
      * ```plasma``` - A KDE desktop environment.
      * ```i3``` - Maybe you want to test your development life in an i3 desktop environment?

It is up to you how you configure your layers. Or maybe want just one recipe for everyday usage. It is up to you how granular you get.