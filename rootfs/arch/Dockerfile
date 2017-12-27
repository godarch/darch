FROM base/archlinux

# Prepare the build container
RUN pacman -Syu --noconfirm
RUN pacman -S --noconfirm arch-install-scripts squashfs-tools

# Add our rootfs
ADD rootfs.tar.gz /

# Add custom arch-chroot that handles custom mounts
ADD arch-chroot-custom /usr/sbin/
ADD arch-chroot /usr/sbin

# Add our cpio hooks to be used when building the image
ADD initcpio-hooks/ /initcpio-hooks
RUN mkdir -p /root.x86_64/etc/mkinitcpio.d/
RUN cp /initcpio-hooks/darch.preset /root.x86_64/etc/mkinitcpio.d/linux.preset
RUN cp /initcpio-hooks/darch-mkinitcpio.conf /root.x86_64/etc/mkinitcpio.conf
RUN mkdir -p /root.x86_64/usr/lib/initcpio/hooks && cp /initcpio-hooks/darch.hook /root.x86_64/usr/lib/initcpio/hooks/darch
RUN mkdir -p /root.x86_64/usr/lib/initcpio/install && cp /initcpio-hooks/darch.install /root.x86_64/usr/lib/initcpio/install/darch

ADD helpers/darch-longrunning /
ADD helpers/darch-prepare /
ADD helpers/darch-runimage /
ADD helpers/darch-teardown /

CMD [ "/darch-longrunning" ]
