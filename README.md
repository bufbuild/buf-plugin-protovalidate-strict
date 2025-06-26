# buf-plugin-protovalidate-strict

An experimental [Buf check plugin](https://buf.build/docs/cli/buf-plugins) for [Protovalidate](https://buf.build/docs/protovalidate) that enforces:

- Protovalidate annotations are not added to existing fields without Protovalidate annotations.
- Existing Protovalidate annotations are never modified or removed.

This plugin only allows Protovalidate annotations to be added to new fields.

This guarantees safety in Protovalidate annotation evolution via the strictest means possible. Generally, this level of strictness is not requires and is not even desirable.

The single rule `PROTOVALIDATE_STRICT` must be added to `buf.yaml` to enable this check:

```
version: v2
breaking:
  use:
    - STANDARD # Or whatever your existing rules are
    - PROTOVALIDATE_STRICT
plugins:
  - plugin: buf-plugin-protovalidate-strict
```

Then, install this plugin locally:

```bash
go install github.com/bufbuild/buf-plugin-protovalidate-strict@latest
```

You can also [publish this plugin to your enterprise BSR instance]((https://buf.build/docs/cli/buf-plugins/publish) and use it without installing it locally.

## Status: Alpha

buf-plugin-protovalidate-strict is an experimental plugin and is in active development. It may be changed or archived without notice.

## Legal

Offered under the [Apache 2 license](https://github.com/bufbuild/buf-plugin-protovalidate-strict/blob/main/LICENSE).
