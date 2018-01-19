FROM scratch

ADD root.x86_64 /

# # Add our cpio hooks to be used when building the image
# TODO: Make this an AUR package? or a pkg.tar.gz?
RUN mkdir -p /etc/mkinitcpio.d/
ADD /initcpio-hooks/darch.preset /etc/mkinitcpio.d/linux.preset
ADD /initcpio-hooks/darch-mkinitcpio.conf /etc/mkinitcpio.conf
RUN mkdir -p /usr/lib/initcpio/hooks
ADD /initcpio-hooks/darch.hook /usr/lib/initcpio/hooks/darch
RUN mkdir -p /usr/lib/initcpio/install
ADD /initcpio-hooks/darch.install usr/lib/initcpio/install/darch

ADD helpers/darch-longrunning /
ADD helpers/darch-prepare /
ADD helpers/darch-runimage /
ADD helpers/darch-teardown /
ADD helpers/darch-extract /
