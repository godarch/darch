#!/usr/bin/env bash
set -e

# Packages needed: debootstrap arch-install-scripts virtualbox

# Download the debs to be used to install debian
if [ ! -e "debs.tar.gz" ]; then
    debootstrap --verbose \
        --make-tarball=debs.tar.gz \
	    --include=linux-image-amd64,grub2 \
	    stable rootfs https://deb.debian.org/debian
fi

# Create our hard disk
rm -rf boot.img
truncate -s 20G boot.img
parted -s boot.img \
        mklabel msdos \
        mkpart primary 0% 2GiB \
        mkpart primary 2GiB 2.5GiB \
        mkpart primary 2.5GiB 3GiB \
        mkpart primary 3GiB 100%

# Partition layout:
# 1. The base recovery os
# 2. Darch configuration (/etc/darch)
# 3. Home directory
# 3. Darch stage/images

# Mount the newly created drive
loop_device=`losetup --partscan --show --find boot.img`

# Format the partitions
mkfs.ext4 ${loop_device}p1
mkfs.ext4 ${loop_device}p2
mkfs.ext4 ${loop_device}p3
mkfs.ext4 ${loop_device}p4

# Mount the new partitions
rm -rf rootfs && mkdir rootfs
mount ${loop_device}p1 rootfs
mkdir -p rootfs/etc/darch
mount ${loop_device}p2 rootfs/etc/darch
mkdir rootfs/home
mount ${loop_device}p3 rootfs/home
mkdir -p rootfs/var/lib/darch
mount ${loop_device}p4 rootfs/var/lib/darch

# Generate the rootfs
debootstrap --verbose \
    --unpack-tarball=$(pwd)/debs.tar.gz \
    --include=linux-image-amd64,grub2 \
    stable rootfs https://deb.debian.org/debian

# Generate fstab (removing comments and whitespace)
genfstab -U -p rootfs | sed -e 's/#.*$//' -e '/^$/d' > rootfs/etc/fstab

# Set the computer name
echo "darch-demo" > rootfs/etc/hostname

# Update all the packages
arch-chroot rootfs apt-get update

# Install network manager for networking
arch-chroot rootfs apt-get -y install network-manager

# Install GRUB
arch-chroot rootfs grub-install ${loop_device}
arch-chroot rootfs grub-mkconfig -o /boot/grub/grub.cfg

# Create the default users
arch-chroot rootfs apt-get -y install sudo
arch-chroot rootfs /usr/bin/bash -c 'echo -en "root\nroot" | passwd'
arch-chroot rootfs useradd -m -G users,sudo -s /usr/bin/bash darch
arch-chroot rootfs /usr/bin/bash -c 'echo -en "darch\ndarch" | passwd darch'

# Install Darch
arch-chroot rootfs apt-get -y install curl gnupg software-properties-common
arch-chroot rootfs /bin/bash -c "curl -L https://raw.githubusercontent.com/godarch/debian-repo/master/key.pub | apt-key add -"
arch-chroot rootfs add-apt-repository 'deb https://raw.githubusercontent.com/godarch/debian-repo/master/darch testing main'
arch-chroot rootfs apt-get update
arch-chroot rootfs apt-get -y install darch
arch-chroot rootfs mkdir -p /etc/containerd
echo "root = \"/var/lib/darch/containerd\"" > rootfs/etc/containerd/config.toml
arch-chroot rootfs systemctl enable containerd

# Setup the fstab hooks for Darch
cat rootfs/etc/fstab | tail -n +2 > rootfs/etc/darch/hooks/default_fstab
echo "*=default_fstab" > rootfs/etc/darch/hooks/fstab.config

# Run grub-mkconfig again to ensure it loads the Darch grub config file
arch-chroot rootfs grub-mkconfig -o /boot/grub/grub.cfg

# Clone our examples repo
arch-chroot rootfs apt-get -y install git
arch-chroot rootfs git clone https://github.com/godarch/example-recipes.git /home/darch/example-recipes
arch-chroot rootfs mkdir /home/darch/Desktop
arch-chroot rootfs ln -s /home/darch/example-recipes /home/darch/Desktop/Recipes
arch-chroot rootfs chown -R darch:darch /home/darch/

# Clean up
umount rootfs/etc/darch
umount rootfs/var/lib/darch
umount rootfs/home
umount rootfs
losetup -d ${loop_device}

# Generate the vdi for VirtualBox
VBoxManage convertdd boot.img boot.vdi --format VDI
chmod 777 boot.vdi