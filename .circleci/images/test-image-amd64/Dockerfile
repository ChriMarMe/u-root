# Copyright 2018-2019 the u-root Authors. All rights reserved
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

FROM circleci/golang:1.12

# Install dependencies
RUN sudo apt-get update &&                          \
    sudo apt-get install -y --no-install-recommends \
        `# Linux dependencies`                      \
        git                                         \
        bc                                          \
        bison                                       \
        flex                                        \
        gcc                                         \
        make                                        \
        `# QEMU dependencies`                       \
        libglib2.0-dev                              \
        libfdt-dev                                  \
        libpixman-1-dev                             \
        zlib1g-dev                                  \
        libcap-dev                                  \
        libattr1-dev                                \
        `# Linux kernel build deps`                 \
        libelf-dev                                  \
        `# Multiboot kernel build deps`             \
        gcc-multilib                                \
        gzip &&                                     \
        `# libtpms & swtpm dependencies`            \
        dh-exec	                            		\
        gawk     			                        \
	    libfuse-dev 		                        \
	    libglib2.0-dev 		                        \
	    libgmp-dev 		                            \
	    expect 			                            \
	    libtasn1-dev 	                        	\
	    socat 			                            \
	    tpm-tools 		                            \
        gnutls-dev                          		\
	    gnutls-bin 		                            \	
	    net-tools 		                            \
	    libseccomp-dev		                        \
        trousers                                    \
        libtool 		                            \
	    libssl-dev		                            \                                
    sudo rm -rf /var/lib/apt/lists/*

# Create working directory
WORKDIR /home/circleci
COPY config_linux4.17_x86_64.txt .config

# Build linux
RUN set -eux;                                                             \
    git clone --depth=1 --branch=v4.17 https://github.com/torvalds/linux; \
    sudo chmod 0444 .config;                                              \
    mv .config linux/;                                                    \
    cd linux;                                                             \
    make olddefconfig;                                                    \
    make -j$(($(nproc) * 2 + 1));                                         \
    cd ~;                                                                 \
    cp linux/arch/x86_64/boot/bzImage bzImage;                            \
    rm -rf linux/

# Build QEMU
RUN set -eux;                                                          \
    git clone --depth=1 --branch=v2.12.0 https://github.com/qemu/qemu; \
    cd qemu;                                                           \
    mkdir build;                                                       \
    cd build;                                                          \
    ../configure                                                       \
        --target-list=x86_64-softmmu                                   \
        --enable-virtfs                                                \
        --disable-docs                                                 \
        --disable-sdl                                                  \
        --disable-kvm;                                                 \
    make -j$(($(nproc) * 2 + 1));                                      \
    cd ~;                                                              \
    cp -rL qemu/build/pc-bios/ ~/pc-bios;                              \
    cp qemu/build/x86_64-softmmu/qemu-system-x86_64 .;                 \
    rm -rf qemu/

# Build Multiboot kernel
RUN set -eux;                                                          \
    git clone --depth=1 --branch=v1.01                                 \
           https://github.com/u-root/multiboot-test-kernel;            \
    cd multiboot-test-kernel;                                          \
    make;                                                              \
    cd ~;                                                              \
    cp multiboot-test-kernel/kernel ./;                                \
    gzip -k kernel;                                                    \
    rm -rf multiboot-test-kernel/

# Build libtpms & swtpm
RUN set -eux;                                                          \
    git clone --depth=1 --branch=v0.7.0                                \
        https://github.com/stefanberger/libtpms;                       \
    cd libtpms;                                                        \
    ./autogen.sh --with-tpm2 --with-openssl;                           \
    make                                                               \
    make install;                                                      \
    cd ~;                                                              \             
    rm -rf libtpms/;                                                   \
    git clone --depth=1 --branch=v0.2.0                                \
        https://github.com/stefanberger/swtpm;                         \
    cd swtpm;                                                          \
    make -j$(($(nproc) * 2 + 1));                                      \
    make install;                                                      \
    mkdir /tmp/tpm1
    mkdir /tmp/tpm2                                                    \
    cd ~;                                                              \
    rm -rf swtpm/

RUN ldconfig;                                                          \
     swtpm_setup --tpmstate /tmp/tpm1/ --create-ek-cert;               \               
     swtpm_setup --tpm2 --tpmstate /tmp/tpm2/ --create-ek-cert

RUN swtpm socket --tpmstate dir=/tmp/tpm1 --ctrl type=unixio,path=/tmp/tpm1/swtpm-sock; \
    swtpm socket --tpmstate dir=/tmp/tpm2 --tpm2 --ctrl type=unixio,path=/tmp/tpm2/swtpm-sock

# Export paths to binaries.
ENV UROOT_KERNEL /home/circleci/bzImage
ENV UROOT_QEMU "/home/circleci/qemu-system-x86_64 -L /home/circleci/pc-bios -m 1G"
ENV UROOT_TESTARCH amd64
#ENV UROOT_SWTPM "not sure about the path usually its in /usr/bin/swtpm"