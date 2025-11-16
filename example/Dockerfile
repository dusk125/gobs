FROM ubuntu:24.04 AS base

RUN apt-get update \
   && apt-get install -y software-properties-common python3-launchpadlib \
   && add-apt-repository ppa:obsproject/obs-studio \
   && apt-get update \
   && apt-get install -y \
   obs-studio \
   ffmpeg \
   wget \
   xvfb \
   dbus-x11 \
   x11-utils \
   vlc \
   && apt-get clean \
   && rm -rf /var/lib/apt/lists/*

FROM base AS go

RUN wget https://golang.google.cn/dl/go1.25.4.linux-amd64.tar.gz \
   &&  tar -C /usr/local -xzf go1.25.4.linux-amd64.tar.gz

WORKDIR /app 

COPY . .

RUN /usr/local/go/bin/go build -v -o app

FROM base

RUN mkdir -p /run/user/1000 \
   && chown 1000:1000 /run/user/1000

USER 1000
WORKDIR /app 

COPY --from=go /app/app .

ENV NEED_XVFB=1

ENV WINDOWID=0
ENV DRI_PRIME=pci-0000_03_00_0
ENV XMODIFIERS=@im=ibus
ENV IM_CONFIG_CHECK_ENV=1
ENV IM_CONFIG_PHASE=1
ENV WAYLAND_DISPLAY=wayland-1
ENV CLUTTER_IM_MODULE=ibus
ENV SHLVL=1
ENV MOZ_ENABLE_WAYLAND=1
ENV GSM_SKIP_SSH_AGENT_WORKAROUND=true

ENV GDK_BACKEND=wayland,x11

ENV XDG_SEAT=seat0
ENV XDG_SESSION_DESKTOP=COSMIC
ENV XDG_SESSION_TYPE=wayland
ENV XDG_CURRENT_DESKTOP=COSMIC
ENV XDG_VTNR=1
ENV XDG_SESSION_ID=3
ENV XDG_RUNTIME_DIR=/run/user/1000

ENV GTK_IM_MODULE=ibus
ENV GTK_MODULES=gail:atk-bridge

ENV QT_ACCESSIBILITY=1
ENV QT_QPA_PLATFORM=wayland;xcb
ENV QT_IM_MODULE=ibus
ENV QT_AUTO_SCREEN_SCALE_FACTOR=1
ENV QT_ENABLE_HIGHDPI_SCALING=1

CMD [ "/app/app" ]
