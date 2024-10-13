# gophash

Perceptual hash in golang with CGO.

Only implemented for PNG image.

## Usage

### Install dependencies

#### MacOS

- Xcode
- C libs
  ```shell
  brew install libpng
  ```

#### Ubuntu

```shell
sudo apt-get update && sudo apt-get install -y build-essential libpng-dev
```

### Build

```shell
git clone https://github.com/allape/gophash.git
cd gophash/cphash
g++ -c -fPIC phash.cpp -o phash.o
g++ -lpng -shared -o libphash.so phash.o
cp libphash.so /usr/local/lib/libphash.so
```

### FAQ

- `Library not loaded: libphash.so`.
    - Put libphash.so in the root directory of your project.
      ```shell
      cp /usr/local/lib/libphash.so /path/to/your/project
      ```

### License for C language source code

https://github.com/aetilius/pHash?tab=GPL-3.0-1-ov-file#readme
