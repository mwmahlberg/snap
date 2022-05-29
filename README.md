snap
====

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=mwmahlberg_snap&metric=alert_status)](https://sonarcloud.io/dashboard?id=mwmahlberg_snap) [![Coverage](https://sonarcloud.io/api/project_badges/measure?project=mwmahlberg_snap&metric=coverage)](https://sonarcloud.io/dashboard?id=mwmahlberg_snap) [![Technical Debt](https://sonarcloud.io/api/project_badges/measure?project=mwmahlberg_snap&metric=sqale_index)](https://sonarcloud.io/dashboard?id=mwmahlberg_snap)

[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=mwmahlberg_snap&metric=sqale_rating)](https://sonarcloud.io/dashboard?id=mwmahlberg_snap) [![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=mwmahlberg_snap&metric=reliability_rating)](https://sonarcloud.io/dashboard?id=mwmahlberg_snap) [![Code Smells](https://sonarcloud.io/api/project_badges/measure?project=mwmahlberg_snap&metric=code_smells)](https://sonarcloud.io/dashboard?id=mwmahlberg_snap) [![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=mwmahlberg_snap&metric=vulnerabilities)](https://sonarcloud.io/dashboard?id=mwmahlberg_snap)

snap is a command line tool to compress and decompress files using the [snappy][snappy]
compression algorithm.

Usage
-----

```plaintext
Usage: ./snap [<in-file>]

(de-)compress files using snappy algorithm

Arguments:
  [<in-file>]    file to (de)compress

Flags:
  -h, --help            Show context-sensitive help.
  -d, --unsnap          uncompress file instead of compressing it
  -k, --keep            keep original file
  -c, --stdout          write to stdout
  -S, --suffix=".sz"    set the suffix
      --version
```

[snappy]: https://google.github.io/snappy/ "Snappy project page on GitHub"