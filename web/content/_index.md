---
title: "Darch"
subtitle: "Boot multiple stateless/clean environments."
title: Requirements
---

# Darch

## What is it?

Think Dockerfiles, but for bootable, stateless, graphical (or not) environments.

After each reboot, all changes made to your operating system are wiped, except where explicitly configured not to (home directory, game data, etc). Easily switch between images with individual grub menu entries.

Your images can be layered/inherited (like Dockerfiles) using recipes.