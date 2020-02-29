# Resync

Cross platform directory synchronization.

![CI Linux Mac](https://github.com/enr/resync/workflows/CI%20Linux%20Mac/badge.svg)
![CI Windows](https://github.com/enr/resync/workflows/CI%20Windows/badge.svg) https://enr.github.io/resync/

Resync synchronize target directory to source and write a report file (`.resync`) in both directories.

Under the hood, it uses `rsync` on Linux/Mac or `robocopy` or `xcopy` on Windows.

Run with "`--noop`" to see the actual command, but not executing it.

*Configuration*

Resync uses a configuration file `${HOME}/.resyncrc`:

```yaml
folders:
    pictures:
        local_path: ~/Pictures
        external_path: Pictures
    documents:
        local_path: ~/Documents
        external_path: Documents
```

*Usage*

To synchronize a single directory:

```
resync [--noop] /path/to/source/ /path/to/target/
```

To synchronize all directories registered in configuration file `${HOME}/.resyncrc`:

```
resync to [--noop] /path/to/target
```

Using the above configuration file, this command will synchronize `~/Pictures` to
`/path/to/target/Pictures` and `~/Documents` to `/path/to/target/Documents`.

## License

**Apache 2.0**

```
Copyright 2017 resync contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
