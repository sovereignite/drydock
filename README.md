# bootstrap

Kubernetes bootstrap coordinator for Sovereignite nodes. Extracted from
[`github.com/sovereignite/sovereignite`](https://github.com/sovereignite/sovereignite)
(`internal/bootstrap` + `cmd/bootstrap`).

The root package `bootstrap` drives the ordered, journaled steps that turn a
bare Sovereignite node into a cluster control plane: TPM-keymanager CA
signing, kubelet config, Calico IPv6 (injected ULA), Istio ingress, SPIRE TPM
plugin/device-key config, Kubernetes API-server initialization, control-plane
wait, and cluster config application. Each step is recorded in a versioned,
crash-recoverable journal. The `cmd/bootstrap` binary exposes the coordinator
as a gRPC service (Bootstrap.GetStatus / Bootstrap.StartBootstrap) over a
loopback listener.

The only non-stdlib dependency is the sibling
[`github.com/sovereignite/api`](https://github.com/sovereignite/api) module
(the shared protobuf/gRPC contract), resolved via a local `replace` directive
until it is published.

## Build & test

```sh
go build ./...
go test ./...
```

## Deploy

- `kubernetes/sovereignite.io/bootstrap/` — kustomize manifests (DaemonSet,
  namespace, service account).
- `.ko.yaml` — ko build entry (`main: ./cmd/bootstrap`, static, `CGO_ENABLED=0`).
- `os/systemd/sovereignite-bootstrap.service` — systemd unit.

## License

GPL-2.0-only. See [LICENSE.md](LICENSE.md).
