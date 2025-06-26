# buf-plugin-pv-strict

A [Buf check plugin](https://buf.build/docs/cli/buf-plugins/) for strictly enforcing consistency
in Protovalidate annotations.

The rule `PROTOVALIDATE_STRICT` expects that Protovalidate annotations on existing messages
and fields are not mutated in any way: annotations cannot be added, removed, or adjusted.

New messages and fields are not subject to validation (only existing ones are checked for changes). This results in the following behavior:
- Messages are checked by name, so if the name of a message has changed, then it is treated
  as a new message
- Fields are checked by field number, so if field numbers have shifted, then this could affect
  check results.

To enforce message and field consistency, use [existing Buf breaking change detection rules](https://buf.build/docs/breaking/rules/).

## Usage

### Running the check plugin locally

Build and install the plugin binary to your `$PATH`:

```
$ go install github.com/bufbuild/buf-plugin-pv-strict@latest
```

Configure the plugin in your module's `buf.yaml`:

```
version: v2
breaking:
  use:
    - PROTOVALIDATE_STRICT
plugins:
  - plugin: "buf-plugin-pv-strict"
```

To use this plugin remotely through the BSR, follow [the instructions for publishing plugins to the BSR](https://buf.build/docs/cli/buf-plugins/publish/).

## Status: Dev

buf-plugin-pv-strict is currently under development and distributed for testing purposes.

## Legal

Offered under the [Apache 2 license](https://github.com/bufbuild/buf-plugin-pv-strict/blob/main/LICENSE).
