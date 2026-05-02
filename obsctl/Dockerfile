# Build context: repository root
# docker build -f example/Dockerfile.gamescope -t gobs-gamescope .
#
# Run with GPU passthrough (recommended):
#   docker run --device /dev/dri --cap-add SYS_NICE gobs-gamescope
#
# Override gamescope backend:
#   docker run --device /dev/dri --cap-add SYS_NICE gobs-gamescope \
#     gamescope --backend sdl -W 1920 -H 1080 -w 1920 -h 1080 -- /app/app

# ============================================================
# Stage 1: Build OBS Studio from source
# ============================================================
FROM fedora:43 AS obs-builder

RUN dnf install -y \
    https://download1.rpmfusion.org/free/fedora/rpmfusion-free-release-$(rpm -E %fedora).noarch.rpm

RUN dnf install -y \
    wget xz cmake ninja-build gcc gcc-c++ gcc-objc git make \
    alsa-lib-devel asio-devel extra-cmake-modules \
    ffmpeg-free-devel fontconfig-devel freetype-devel \
    glib2-devel jansson-devel json-devel \
    libavcodec-free-devel libavdevice-free-devel \
    librist-devel libcurl-devel libdatachannel-devel \
    libdrm-devel libglvnd-devel libqrcodegencpp-devel \
    libuuid-devel libva-devel libv4l-devel libvpl-devel \
    libX11-devel libXcomposite-devel libXdamage \
    libXinerama-devel libxkbcommon-devel \
    luajit-devel mbedtls-devel nv-codec-headers \
    pciutils-devel pipewire-devel pulseaudio-libs-devel \
    python3-devel qt6-qtbase-devel qt6-qtbase-private-devel \
    qt6-qtsvg-devel qt6-qtwayland-devel \
    rnnoise-devel simde-devel speexdsp-devel srt-devel \
    swig systemd-devel uthash-devel vlc-devel \
    wayland-devel websocketpp-devel x264-devel

RUN wget https://cdn-fastly.obsproject.com/downloads/cef_binary_6533_linux_x86_64_v6.tar.xz && \
    tar -xf cef_binary_6533_linux_x86_64_v6.tar.xz -C /opt && \
    rm cef_binary_6533_linux_x86_64_v6.tar.xz

RUN git clone --depth 1 --branch 32.0.4 --recurse -j8 \
    https://github.com/obsproject/obs-studio.git /obs-studio

WORKDIR /obs-studio

RUN cmake -S . -B build -G Ninja \
    -DCMAKE_BUILD_TYPE=Release \
    -DCMAKE_INSTALL_PREFIX=/opt/obs \
    -DCMAKE_INSTALL_LIBDIR=lib \
    -DCMAKE_INSTALL_DATAROOTDIR=data \
    -DCMAKE_POSITION_INDEPENDENT_CODE=ON \
    -DENABLE_AJA=OFF \
    -DENABLE_VLC=ON \
    -DENABLE_WAYLAND=ON \
    -DENABLE_WEBRTC=ON \
    -DENABLE_BROWSER=ON \
    -DCEF_ROOT_DIR=/opt/cef_binary_6533_linux_x86_64 \
    -DENABLE_UI=OFF \
    -DENABLE_FRONTEND=OFF \
    -DENABLE_RELOCATABLE=ON \
    -DENABLE_PORTABLE_CONFIG=ON

RUN cmake --build build
RUN DESTDIR=/tmp/obs-install cmake --install build

# ============================================================
# Stage 2: Build gamescope from source
# ============================================================
FROM fedora:43 AS gamescope-builder

RUN dnf install -y \
    meson ninja-build gcc-c++ cmake git pkg-config \
    wayland-devel wayland-protocols-devel \
    libdrm-devel mesa-libEGL-devel mesa-libgbm-devel \
    vulkan-headers vulkan-loader-devel \
    libXcomposite-devel libXdamage-devel libXrender-devel \
    libXres-devel libXtst-devel libXmu-devel \
    libXi-devel libXcursor-devel libXext-devel \
    pixman-devel libcap-devel \
    xorg-x11-server-Xwayland-devel \
    glm-devel stb-devel hwdata \
    libdisplay-info-devel libliftoff-devel \
    libinput-devel libseat-devel \
    SDL2-devel \
    pipewire-devel avahi-devel \
    libei-devel libdecor-devel \
    xcb-util-errors-devel libXrandr-devel \
    libXxf86vm-devel libxkbcommon-devel \
    xcb-util-wm-devel glslang luajit-devel

RUN git clone --recurse-submodules \
    https://github.com/ValveSoftware/gamescope.git /gamescope

WORKDIR /gamescope

RUN meson setup build \
    --prefix=/usr \
    --buildtype=release \
    -Dpipewire=enabled \
    -Dsdl2_backend=enabled

RUN ninja -C build
RUN DESTDIR=/tmp/gamescope-install ninja -C build install

# ============================================================
# Stage 3: Build Go application
# ============================================================
FROM fedora:43 AS go-builder

ARG GO_VERSION=1.26.0

RUN dnf install -y \
    https://download1.rpmfusion.org/free/fedora/rpmfusion-free-release-$(rpm -E %fedora).noarch.rpm

RUN dnf install -y wget pkg-config gcc gcc-c++ simde-devel \
    pulseaudio-libs-devel ffmpeg-free-devel \
    libX11-devel libxcb-devel wayland-devel libxkbcommon-devel

RUN wget -q https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz && \
    rm go${GO_VERSION}.linux-amd64.tar.gz

ENV PATH="/usr/local/go/bin:$PATH"

COPY --from=obs-builder /tmp/obs-install /

COPY . /src/gobs

WORKDIR /src/gobs/obsctl

RUN go mod edit -replace github.com/dusk125/gobs=/src/gobs
RUN go mod tidy

ENV GOWORK=off
RUN PKG_CONFIG_PATH=/opt/obs/lib/pkgconfig \
    CGO_LDFLAGS="-L/opt/obs/lib" \
    CGO_CFLAGS="-I/opt/obs/include" \
    go build -v -o /app

# ============================================================
# Stage 4: Runtime
# ============================================================
FROM fedora:43

RUN dnf install -y \
    https://download1.rpmfusion.org/free/fedora/rpmfusion-free-release-$(rpm -E %fedora).noarch.rpm

RUN dnf install -y --setopt=install_weak_deps=False \
    mesa-vulkan-drivers mesa-dri-drivers \
    vulkan-loader libglvnd-glx libglvnd-egl \
    xorg-x11-server-Xwayland \
    ffmpeg-free vlc-core vlc-plugins-base \
    alsa-lib pipewire pipewire-alsa pulseaudio-libs \
    x264-libs libva libdrm \
    libXcomposite libXres libXmu libxkbcommon \
    wayland-devel luajit \
    jansson mbedtls-devel srt-libs \
    librist speexdsp rnnoise \
    libinput libseat SDL2 libdecor pixman pciutils \
    libdisplay-info libliftoff \
    libei avahi-libs hwdata \
    qt6-qtbase qt6-qtbase-gui nss nspr \
    libdatachannel \
    && dnf clean all

COPY --from=obs-builder /tmp/obs-install/opt/obs /opt/obs
COPY --from=gamescope-builder /tmp/gamescope-install/usr /usr

RUN mkdir -p /run/user/0

WORKDIR /app

COPY --from=go-builder /app ./app

ENV XDG_RUNTIME_DIR=/run/user/0
ENV LD_LIBRARY_PATH=/opt/obs/lib

ENTRYPOINT ["gamescope", "--backend", "headless", "-W", "1920", "-H", "1080", "-w", "1920", "-h", "1080", "-r", "60", "--", "/app/app"]
