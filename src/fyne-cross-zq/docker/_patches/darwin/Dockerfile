ARG FYNE_CROSS_VERSION

FROM fyneio/fyne-cross:${FYNE_CROSS_VERSION}-darwin

RUN apt-get update \
    && apt-get install -y -q --no-install-recommends \
        libzmq3-dev \
    && apt-get -qy autoremove \
    && apt-get clean \
    && rm -r /var/lib/apt/lists/*;

COPY ./scripts/setup-zeromq.sh .
RUN chmod +x ./setup-zeromq.sh && ./setup-zeromq.sh

ENV CPATH $CPATH:/usr/include

RUN export CPATH=$CPATH:/usr/include

RUN LIBRARY_PATH=$LIBRARY_PATH:/usr/local/lib && export LIBRARY_PATH
