# Thoughts on hooks

*not done*

Hooks are used to configure scripts that will be run on the root image, directly before chrooting into the image for booting.

Hooks can be used to prep a /etc/fstab, /etc/hostname, etc.

# The boot process

At a minimum, an image will contain the following files, which are required to boot.

```
└── /var/darch/staged/image/tag
    ├── image.json
    ├── initramfs-linux.img
    └── rootfs.squash
```

Optionally, any number of hooks can be associated with the image.

```
└── /var/darch/staged/image/tag
    ├── ...
    └── hooks
        ├── custom-hook
        |   └── hook
        ├── another-hook
        |   └── hook
        └── 00_custom-hook
        └── 01_another-hook
```

After the initramfs mounts the rootfs for booting, all the scripts under ```hook``` will be run. The prefx ```00``` is determines the execution order of the hook (defined later).

# hooks-config.json

Hooks can be include/excluded into images depending upon the ```/var/darch/hooks/hooks-config.json``` file.

```json
{
    "_default": {
        "execution-order": 0,
        "include-images": ["*"],
        "exclude-images": []
    }
}
```

The ```_default``` property isn't required, since it is assumed. The above ```hooks-config.json``` is the same as the following.

```json
{
}
```

# Hooks

Hooks are placed into ```/var/darch/hookes/$HOOK_NAME```, for example, ```/var/darch/hooks/fstab```

The ```fstab``` hook can be configured to run before other hooks by updating the ```hooks-config.json``` as followed.

```json
{
    "fstab": {
        "execution-order": -1,
    }
}
```

By default, ```fstab``` will be applied to all images/tags, but it can be configured to only run on certain images.

```json
{
    "fstab": {
        "include-images": ["only-these-images-*:*", "and-these-also:with-this-tag-prefix-*"]
    }
}
```

You can also configure ```fstab``` to **not** be run on specific images.

```json
{
    "fstab": {
        "exclude-images": ["ignore-this-image:latest"]
    }
}
```

Inside of the ```/var/darch/hooks/fstab``` should be a ```hook``` script. This script will define what happens when the hook is applied to an image, and what happens when the hook is executed during boot.

```bash
#!/bin/bash

help() {
    # Used to output help documentation about the hook
    echo "..."
}

install() {
    # Called when the hook needs to be installed to an image/tag.
    # Some environment variables are used to provide information to installers.
    # DARCH_HOOKS_DIR=/var/darch/hooks
    # DARCH_HOOK_NAME=fstab
    # DARCH_HOOK_SRC_DIR=/var/darch/hooks/fstab
    # DARCH_HOOK_DEST_DIR=/var/darch/staged/image/tag/hooks/fstab
    echo "..."
}

run() {
    # Called when the hook is run, during boot.
    # Some environment variables are used to provide information to the running hook.
    # DARCH_HOOK_DIR=/tmp/before/chroot/hooks/fstab/
    # DARCH_ROOT_FS=/tmp/rootfs
}
```

Here is what a hook would look like for ``fstab```.

```bash
#!/bin/bash

help() {
    echo "Places /var/darch/defaultfstab into /etc/fstab, before booting..."
}

install() {
    cp "/var/darch/defaultfstab" "$DARCH_HOOK_DEST_DIR/fstab"
}

run() {
    cp "$DARCH_HOOK_DIR/fstab" "$DARCH_ROOT_FS/etc/fstab"
}
```