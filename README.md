# v2dat

A cli tool that can unpack v2ray data packages, also known as `geoip.dat` and `geosite.dat`.

## Usage

```shell
v2dat unpack geoip [-d output_dir] [-f tag]... geoip_file
v2dat unpack geosite [-d output_dir] [-f tag[@attr]...]... geosite_file
```

- If `-d` was omitted, the current working dir `.` will be used.
- If no filter `-f` was given. All tag will be unpacked.
- If multiple `@attr` were given. Entries that don't contain any of given attrs will be ignored.
- Unpacked text files will be named as `<geo_file>_<suffix>.txt`.