# Partial code of [pHash](https://github.com/aetilius/pHash)

# License

https://github.com/aetilius/pHash?tab=GPL-3.0-1-ov-file#readme

## Build on macOS (Xcode is required)

### Dependencies

```shell
brew install libpng
```

### Shared library

```shell
rm -f phash.o libphash.so
g++ -c -fPIC phash.cpp -o phash.o
g++ -lpng -shared -o libphash.so phash.o
```
