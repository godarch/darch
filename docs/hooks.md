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

# Hook installation/configuring

Each hook has a script that will be used for configuring an image.
