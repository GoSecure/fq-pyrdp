# fq-pyrdp

[`fq`](https://github.com/wader/fq) format for parsing PyRDP replays.

_In progress._

## How to use the format

### Get the codebase

You will need Go >= 1.17 installed. If present, clone the fq repository:

```bash
$ git clone git@github.com:wader/fq.git
$ cd fq
```


### Add the format

Then add the `pyrdp` format to the import list at `format/all/all.go`:

```go
_ "github.com/wader/fq/format/pyrdp"
```

And also to the constants at `format/format.go`:

```go
        PYRDP               = "pyrdp"
```

Now that the format was added to the list of available formats, the actual code needs to be copied or linked in the `format/` directory. Here's how you can do the later (as it can be updated independently):

```bash
$ cd format/
$ ln -s /the/path/to/fq-pyrdp/pyrdp pyrdp
```

### Build

To build fq with the new format added just go to the root of the fq repository and use `make fq`. That should create an fq binary in the same directory with the new PyRDP replay format added.

## fq usage

To parse the replay files using fq you will need to specify the format using `-d` and a query, just as you will do with jq:

```bash
$ ./fq -d pyrdp '.events[]' /the/path/to/replay.pyrdp
```

More complex information can be extracted depending on the PDUs that the pyrdp format can parse. For example, we can get the password used by the user that connected to the RDP service:

```json
$ ./fq -d pyrdp '.events[1].client_info|{password:.password,username:.username}' /the/path/to/replay.pyrdp
{
  "password": "admin",
  "username": "administrator"
}
```
