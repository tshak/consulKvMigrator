# ConsulKvMigrator

A simple tool for [Consul](http://consul.io) KV data migrations. This is useful for situations where you want source control to be
the source of truth for certain KV items. This tool can be run interactively or as part of a deployment.

## Usage

`$ consulKvMigrator [options] migration_file.json`

### Options
```
  -address string
    	address for consul (default "localhost:8500")
  -dry-run
    	don't submit any changes to consul
  -prompt
    	prompt before submitting changes to consul (default true)
```

### Example output

```
Changes found: 1
+------+-----------+-----------+
| Key  | Old Value | New Value |
+------+-----------+-----------+
| foo  | 43        | 42        |
+------+-----------+-----------+

Are you sure you want to apply these changes? [Y/n]
```

## Migration file format

The migration file is a simple JSON document with each node in the object graph representing a key path. For example,
for the KV pair: `foo/bar` with the value `quux`, the migration file would look like:

```
{
  "foo": {
      "bar": { "value": "quux" }
  }
}
```

This format allows for more human readable organization of keys in their key "folders".


## FAQ

### Why not YAML?

YAML would be far more readable. Unfortunately, the only YAML library capable of easily parsing arbitrary document
structures is [Licensed under the LGPLv3](https://github.com/go-yaml/yaml/blob/v2/LICENSE), which is problematic for
many organizations. Other YAML libraries are just support basic marshalling, or are just wrappers around this one.

## Why did you code it like X?
This is my first attempt at golang, so I wouldn't expect this project to represent idiomatic golang :). I'd love feedback from experienced Gophers though!
