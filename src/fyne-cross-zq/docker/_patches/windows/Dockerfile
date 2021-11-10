ARG FYNE_CROSS_VERSION

FROM fyneio/fyne-cross:${FYNE_CROSS_VERSION}-windows

RUN apt-get update \
    && apt-get install -y -q --no-install-recommends \
        libtool pkg-config build-essential autoconf automake uuid-dev libzmq3-dev \
    && apt-get -qy autoremove \
    && apt-get clean \
    && rm -r /var/lib/apt/lists/*;

#ENV LD_DEBUG "all make"

#ENV CPATH $CPATH:/usr/share/mingw-w64/include/:/usr/include



# COPY ./win_build/lib/* ./win_build/lib
# COPY ./win_build/include/* ./win_build/include

# COPY ./win_build/lib/* /usr/local/lib
# COPY ./win_build/include/* /usr/local/include

# COPY ./scripts/setup-zeromq-windows.sh .
# RUN chmod +x ./setup-zeromq-windows.sh && ./setup-zeromq-windows.sh

# ENV CGO_CFLAGS -I/usr/local/include
# ENV CGO_LDFLAGS -L/usr/local/lib

# CXX_FOR_TARGET="i686-w64-mingw32-g++" CC_FOR_TARGET="i686-w64-mingw32-gcc" CGO_ENABLED=1 CGO_CFLAGS="-I/usr/local/include" CGO_LDFLAGS="-L/usr/local/lib" GOOS=windows GOARCH=amd64 go build -ldflags '-H windowsgui -w -s' -o ./aaaaaaaaaaa.exe -v ./cmd/app
# CGO_CFLAGS=-I/usr/local/include CGO_LDFLAGS=-L/usr/local/lib GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build -ldflags '-H windowsgui -w -s -extld=$CC' ./cmd/app
#ENV CGO_CPPFLAGS "-I/usr/include"
#ENV CGO_LDFLAGS "-L/usr/lib -L/usr/lib/x86_64-linux-gnu -lzmq -lpthread -lsodium -lrt -lstdc++ -lm -lc -lgcc"
# go build -v --ldflags '-extldflags "-static"' -a subscriber.go

# ENV CGO_CFLAGS '-ID:/usr/local/include'
# ENV CGO_LDFLAGS '-L/usr/lib -L/usr/lib/x86_64-linux-gnu -LD:/usr/local/lib -lzmq'


# ENV LIBRARY_PATH $LIBRARY_PATH:/usr/local/lib:/usr/lib/x86_64-linux-gnu
# ENV CPATH $CPATH:/usr/share/mingw-w64/include/:/usr/include

# ENV CGO_LDFLAGS "-L/usr/local/lib -lzmq"
# ENV CGO_CPPFLAGS "-I/usr/include"
# ENV CGO_LDFLAGS "-L/usr/lib -L/usr/lib/x86_64-linux-gnu -lzmq -lpthread -lsodium -lrt -lstdc++ -lm -lc -lgcc"

# RUN export CPATH=$CPATH:/usr/include
# RUN LIBRARY_PATH=$LIBRARY_PATH:/usr/local/lib && export LIBRARY_PATH
