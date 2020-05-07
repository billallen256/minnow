# minnow

A minimalist publish/subscribe file processing framework with no dependencies.  Processors can be implemented in any language since they are simply called as shell commands.  Each time a processor is executed, it is given an input directory to read from and an output directory to write to.  All data files are accompanied by a metadata `.properties` file with the same base name.

To start minnow, just run:

`minnow /path/to/config.properties`

Minnow is meant to be run on a single machine, though there is no reason a processor can't be written to transfer files to another machine's minnow ingest directory.

## Configuration
Minnow needs very little configuration to get started.  Here is an example `config.properties` file:

```
ingest_dir = /data/minnow/ingest
ingest_min_age = 600
work_dir = /data/minnow/work
processor_definitions_dir = /data/minnow/processors
```

### `ingest_dir`
This tells minnow where to look for incoming files.  All data files must be accompanied by a metadata `.properties` file with the same base name (eg. `blueprints.dwg` and `blueprints.dwg.properties`).

### `ingest_min_age`
This tells minnow how long to wait before actually ingesting a file, in seconds.  This gives time for copies/writes to complete before ingest.

### `work_dir`
This is scratch space that the processors will use during execution.  You can also look in here during debugging, as failed processor execution leave behind a dated output directory containing a file with output from stdout and stderr.

### `processor_definitions_dir`
Minnow will look for processors in this directory.  Each processor has its own definition directory.  Minnow will scan for new processors every five minutes.

Processors can be hot-swapped while minnow is running by changing the contents of the processor's definition directory.  Note that it's probably best to make changes in a separate directory, then drop the changed files in with an atomic `mv` (move) command.

## Defining a Processor
A processor is simply a directory containing a `config.properties` file, a start script, and a hook file.  Minnow parses the processor's `config.properties` file to find the name of the start script and hook file, for example:

```
start_script = start.sh
hook_file = hook.properties
```

### Start Script
Minnow will pass only two parameters to the start script: an input directory to read files from, and an output directory to write files to.  The start script can then use them, or pass them on to a more complex script.  For example, `start.sh` could contain the following:

```sh
#!/bin/sh

# pass all params on to the python script
# param 1 is the input directory to read from
# param 2 is the output directory to write to
python3 amazing.py $@
```

Calling containers is also straight forward, just bind the input and output directories so they are accessible within the container.  Here's how to do it with Singularity (note that filesystem overlay must be enabled):

```sh
singularity exec --no-home --bind $1:/input:ro,$2:/output:rw /path/to/awesome.sif
```

And for Docker:

```sh
docker run --rm -v $1:/input -v $2:/output awesome:latest
```

### Hook File
Processors subscribe to the data/metadata they want by defining a hook `.properties` file.  If the metadata `.properties` file matches all the key=value pairs in the processor's hook, it will get a copy of that data/metadata to process.  For example, `hook.properties` could contain the following:

```
type = blueprints
orientation = above
```

This processor will get a copy of all data/metadata files that contain these properties.  For example, this would match:

```
type = blueprints
size = huge
orientation = above
style = line
```

But this would not match:

```
type = blueprints
size = huge
orientation = front
style = line
```

### Pool Size
The `config.properties` file can also take an optional `pool_size` parameter to indicate the maximum number of instances of the processor should run simultaneously.  This helps prevent resource-intensive processors from taking over the whole system.  By default, `pool_size` is equal to the number of logical CPU's on the system.

### Output
To output data back into minnow, it must go in the output directory that was provided to the start script.  Data/metadata files must come in pairs, or they will be ignored.  For example, if your data file is named `blueprints.dwg`, then there must also be a metadata file called `blueprints.dwg.properties`.  If your script doesn't generate any data, just touch the file so it exists.

To output data out of minnow, just write it somewhere else in the file system per your project.  In fact, **you'll have to create a processor to subscribe and output the data somewhere else in the file system**, as minnow deletes data that does not have a subscriber from its `work_dir`.  The output processor can be as simple as creating a hook for the relevant data and putting the following in the start script:

```sh
#!/bin/sh
cp -R $1/* /path/to/project/output/
```

## Property Files
Minnow uses a dead-simple file format for configuration and metadata `.properties` files.  The only reserved characters are newline and equals (`=`).  Leading and trailing whitespace is stripped from all keys and values.

This format was chosen so it would be as easy as possible for anyone to implement a parser in their language of choice.  For example, a Python parser can be written in just a few lines:

```python
from pathlib import Path  # pathlib is awesome, btw
input_path = Path(sys.argv[1])
output_path = Path(sys.argv[2])
metadata_path = list(input_path.glob('*.properties'))[0]

# properties parser starts here
metadata = metadata_path.read_text().split('\n')
metadata = [tuple(x.split('=')) for x in metadata if len(x.strip()) > 0]
metadata = [x for x in metadata if len(x) == 2]
metadata = {x[0].strip():x[1].strip() for x in metadata}
```

Writing `.properties` files with Python is also simple:

```python
metadata = {'type':'blueprints',
            'orientation':'above',
            'size':'huge',
            'style':'line'}

# this turns the dictionary into a string, ready to be written
metadata_str = '\n'.join(['{} = {}'.format(k, v) for k, v in metadata.items()])

output_metadata_path = output_path.joinpath('blueprints.properties')
output_metadata_path.write_text(metadata_str)
output_data_path = output_path.joinpath('blueprints')  # make sure the data file exists
output_data_path.touch()                               # even if there's no data
```

The original plan was to use JSON, because JSON is awesome.  But minnow is implemented in Go, and Go does not do well with JSON's dynamic structure and types, even with external dependencies.

## Security Notes
It's worth mentioning that minnow does not enforce any security constraints itself, and instead relies on the underlying operating system.  Be aware that all processors run as the user that started minnow.  For example, if user A starts minnow and sets `processor_definition_dir` to a path that's writable by users B and C, then users B and C can put *any script* in there, and it will execute as user A, with access to all of user A's files.  User A better trust users B and C.

One way to help mitigate this concern is to have processors point at trusted containers that have already been reviewed, and to limit the containers' accesses with options like `--net --network=none` for Singularity (see above).
