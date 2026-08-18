// SPDX-License-Identifier: GPL-2.0-only
//
// Copyright (C) 2026 Sovereignite contributors

package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net"
	"net/netip"
	"os"
	"os/signal"
	"syscall"

	"github.com/sovereignite/drydock"
	"google.golang.org/grpc"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()
	if err := run(ctx, os.Args[1:]); err != nil &&
		!errors.Is(err, context.Canceled) {
		log.Fatal(err)
	}
}

func run(ctx context.Context, args []string) error {
	if len(args) != 0 {
		return fmt.Errorf(
			"sovereignite-bootstrap accepts no command-line operations; got %d arguments",
			len(args),
		)
	}

	coordinator, err := newStubCoordinator()
	if err != nil {
		return fmt.Errorf("create bootstrap coordinator: %w", err)
	}

	server, err := drydock.NewServer(coordinator)
	if err != nil {
		return fmt.Errorf("create bootstrap server: %w", err)
	}

	grpcServer := grpc.NewServer()
	server.RegisterGrpcServer(grpcServer)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return fmt.Errorf("create gRPC listener: %w", err)
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- grpcServer.Serve(listener)
	}()

	log.Printf("bootstrap gRPC server listening on %s", listener.Addr().String())

	select {
	case <-ctx.Done():
		grpcServer.GracefulStop()
		return <-errCh
	case err := <-errCh:
		return err
	}
}

func newStubCoordinator() (*drydock.Coordinator, error) {
	config := stubConfiguration()
	store := &memoryJournalStore{}
	deps := drydock.Dependencies{
		CA:         &stubCASigner{},
		Kubernetes: &stubKubernetesInstaller{},
		Calico:     &stubCalicoInstaller{},
		Istio:      &stubIstioInstaller{},
		SPIRE:      &stubSPIREInstaller{},
		Dex:        &stubDexInstaller{},
		Prepared:   &stubPreparedPublisher{},
	}
	return drydock.NewCoordinator(config, store, deps)
}

var errStubDependencyUnavailable = errors.New(
	"stub dependency unavailable; production wiring required",
)

type memoryJournalStore struct {
	journal drydock.Journal
}

func (m *memoryJournalStore) Load() (drydock.Journal, error) {
	return m.journal, nil
}

func (m *memoryJournalStore) Save(
	_ uint64,
	journal drydock.Journal,
) error {
	m.journal = journal
	return nil
}

type stubCASigner struct{}

func (stubCASigner) EnsureSigning(
	context.Context,
	drydock.Operation,
	drydock.CARequest,
) (drydock.Observation, error) {
	return drydock.Observation{}, errStubDependencyUnavailable
}

func (stubCASigner) VerifySigning(
	context.Context,
	drydock.Operation,
	drydock.CARequest,
	drydock.Observation,
) error {
	return errStubDependencyUnavailable
}

type stubKubernetesInstaller struct{}

func (stubKubernetesInstaller) PrepareKubelet(
	context.Context,
	drydock.Operation,
	drydock.Artifact,
	drydock.Artifact,
) (drydock.Observation, error) {
	return drydock.Observation{}, errStubDependencyUnavailable
}

func (stubKubernetesInstaller) VerifyKubelet(
	context.Context,
	drydock.Operation,
	drydock.Artifact,
	drydock.Artifact,
	drydock.Observation,
) error {
	return errStubDependencyUnavailable
}

func (stubKubernetesInstaller) InitializeAPIServer(
	context.Context,
	drydock.Operation,
	drydock.Artifact,
	drydock.Artifact,
) (drydock.Observation, error) {
	return drydock.Observation{}, errStubDependencyUnavailable
}

func (stubKubernetesInstaller) VerifyAPIServer(
	context.Context,
	drydock.Operation,
	drydock.Artifact,
	drydock.Artifact,
	drydock.Observation,
) error {
	return errStubDependencyUnavailable
}

func (stubKubernetesInstaller) WaitControlPlane(
	context.Context,
	drydock.Operation,
	drydock.Artifact,
) (drydock.Observation, error) {
	return drydock.Observation{}, errStubDependencyUnavailable
}

func (stubKubernetesInstaller) VerifyControlPlane(
	context.Context,
	drydock.Operation,
	drydock.Artifact,
	drydock.Observation,
) error {
	return errStubDependencyUnavailable
}

func (stubKubernetesInstaller) Reconcile(
	context.Context,
	drydock.Operation,
	drydock.Artifact,
) (drydock.Observation, error) {
	return drydock.Observation{}, errStubDependencyUnavailable
}

func (stubKubernetesInstaller) VerifyReconciled(
	context.Context,
	drydock.Operation,
	drydock.Artifact,
	drydock.Observation,
) error {
	return errStubDependencyUnavailable
}

func (stubKubernetesInstaller) CheckReady(
	context.Context,
	drydock.Operation,
) (drydock.Observation, error) {
	return drydock.Observation{}, errStubDependencyUnavailable
}

func (stubKubernetesInstaller) VerifyReady(
	context.Context,
	drydock.Operation,
	drydock.Observation,
) error {
	return errStubDependencyUnavailable
}

type stubCalicoInstaller struct{}

func (stubCalicoInstaller) PrepareIPv6(
	context.Context,
	drydock.Operation,
	netip.Prefix,
	drydock.Artifact,
) (drydock.Observation, error) {
	return drydock.Observation{}, errStubDependencyUnavailable
}

func (stubCalicoInstaller) VerifyIPv6(
	context.Context,
	drydock.Operation,
	netip.Prefix,
	drydock.Artifact,
	drydock.Observation,
) error {
	return errStubDependencyUnavailable
}

func (stubCalicoInstaller) Reconcile(
	context.Context,
	drydock.Operation,
	drydock.Artifact,
) (drydock.Observation, error) {
	return drydock.Observation{}, errStubDependencyUnavailable
}

func (stubCalicoInstaller) VerifyReconciled(
	context.Context,
	drydock.Operation,
	drydock.Artifact,
	drydock.Observation,
) error {
	return errStubDependencyUnavailable
}

func (stubCalicoInstaller) CheckReady(
	context.Context,
	drydock.Operation,
) (drydock.Observation, error) {
	return drydock.Observation{}, errStubDependencyUnavailable
}

func (stubCalicoInstaller) VerifyReady(
	context.Context,
	drydock.Operation,
	drydock.Observation,
) error {
	return errStubDependencyUnavailable
}

type stubIstioInstaller struct{}

func (stubIstioInstaller) PrepareIngress(
	context.Context,
	drydock.Operation,
	drydock.Artifact,
	drydock.Artifact,
) (drydock.Observation, error) {
	return drydock.Observation{}, errStubDependencyUnavailable
}

func (stubIstioInstaller) VerifyIngress(
	context.Context,
	drydock.Operation,
	drydock.Artifact,
	drydock.Artifact,
	drydock.Observation,
) error {
	return errStubDependencyUnavailable
}

func (stubIstioInstaller) Reconcile(
	context.Context,
	drydock.Operation,
	drydock.Artifact,
) (drydock.Observation, error) {
	return drydock.Observation{}, errStubDependencyUnavailable
}

func (stubIstioInstaller) VerifyReconciled(
	context.Context,
	drydock.Operation,
	drydock.Artifact,
	drydock.Observation,
) error {
	return errStubDependencyUnavailable
}

func (stubIstioInstaller) CheckReady(
	context.Context,
	drydock.Operation,
) (drydock.Observation, error) {
	return drydock.Observation{}, errStubDependencyUnavailable
}

func (stubIstioInstaller) VerifyReady(
	context.Context,
	drydock.Operation,
	drydock.Observation,
) error {
	return errStubDependencyUnavailable
}

type stubSPIREInstaller struct{}

func (stubSPIREInstaller) PrepareTPM(
	context.Context,
	drydock.Operation,
	drydock.SPIRERequest,
) (drydock.Observation, error) {
	return drydock.Observation{}, errStubDependencyUnavailable
}

func (stubSPIREInstaller) VerifyTPM(
	context.Context,
	drydock.Operation,
	drydock.SPIRERequest,
	drydock.Observation,
) error {
	return errStubDependencyUnavailable
}

func (stubSPIREInstaller) Reconcile(
	context.Context,
	drydock.Operation,
	drydock.Artifact,
) (drydock.Observation, error) {
	return drydock.Observation{}, errStubDependencyUnavailable
}

func (stubSPIREInstaller) VerifyReconciled(
	context.Context,
	drydock.Operation,
	drydock.Artifact,
	drydock.Observation,
) error {
	return errStubDependencyUnavailable
}

func (stubSPIREInstaller) CheckReady(
	context.Context,
	drydock.Operation,
) (drydock.Observation, error) {
	return drydock.Observation{}, errStubDependencyUnavailable
}

func (stubSPIREInstaller) VerifyReady(
	context.Context,
	drydock.Operation,
	drydock.Observation,
) error {
	return errStubDependencyUnavailable
}

type stubDexInstaller struct{}

func (stubDexInstaller) Reconcile(
	context.Context,
	drydock.Operation,
	drydock.DexRequest,
) (drydock.Observation, error) {
	return drydock.Observation{}, errStubDependencyUnavailable
}

func (stubDexInstaller) VerifyReconciled(
	context.Context,
	drydock.Operation,
	drydock.DexRequest,
	drydock.Observation,
) error {
	return errStubDependencyUnavailable
}

func (stubDexInstaller) CheckReady(
	context.Context,
	drydock.Operation,
) (drydock.Observation, error) {
	return drydock.Observation{}, errStubDependencyUnavailable
}

func (stubDexInstaller) VerifyReady(
	context.Context,
	drydock.Operation,
	drydock.Observation,
) error {
	return errStubDependencyUnavailable
}

type stubPreparedPublisher struct{}

func (stubPreparedPublisher) Publish(context.Context) error {
	return errStubDependencyUnavailable
}

func (stubPreparedPublisher) Verify(context.Context) error {
	return errStubDependencyUnavailable
}

func stubArtifact(name, version, content string) drydock.Artifact {
	sum := sha256.Sum256([]byte(content))
	return drydock.Artifact{
		Name:    name,
		Version: version,
		SHA256:  hex.EncodeToString(sum[:]),
		Content: []byte(content),
	}
}

func stubConfiguration() drydock.Configuration {
	const tpmReference = "tpm-device-key:0x81010040"
	return drydock.Configuration{
		BootstrapVersion: "bootstrap-v5.0.0",
		Versions: drydock.PinnedVersions{
			Kubernetes: "v1.35.1",
			Calico:     "v3.31.2",
			Istio:      "v1.29.3",
			SPIRE:      "v1.14.1",
			Dex:        "v2.45.1",
		},
		ULA:                   netip.MustParsePrefix("fd18:4f1c:14d::/48"),
		TPMDeviceKeyReference: tpmReference,
		Artifacts: drydock.Artifacts{
			CARequest: stubArtifact(
				drydock.ArtifactCARequest,
				"bootstrap-v5.0.0",
				"public certificate request",
			),
			Kubelet: stubArtifact(
				drydock.ArtifactKubelet,
				"v1.35.1",
				"apiVersion: kubelet.config.k8s.io/v1beta1\n",
			),
			Calico: stubArtifact(
				drydock.ArtifactCalico,
				"v3.31.2",
				"calico IPv6 desired configuration\n",
			),
			Istio: stubArtifact(
				drydock.ArtifactIstio,
				"v1.29.3",
				"istio ingress desired configuration\n",
			),
			SPIRE: stubArtifact(
				drydock.ArtifactSPIRE,
				"v1.14.1",
				"plugin = tpm_devid\nreference = "+tpmReference+"\n",
			),
			Dex: stubArtifact(
				drydock.ArtifactDex,
				"v2.45.1",
				"dex desired configuration\n",
			),
			ControlPlane: stubArtifact(
				drydock.ArtifactControlPlane,
				"v1.35.1",
				"complete authorized control plane inputs\n",
			),
			Cluster: stubArtifact(
				drydock.ArtifactCluster,
				"v1.35.1",
				"kubernetes cluster desired configuration\n",
			),
		},
		Authority: drydock.AuthorityContracts{
			KubernetesTopology: stubArtifact(
				drydock.ArtifactTopologyContract,
				"decision-006-r1",
				"complete topology and component inventory\n",
			),
			TPMKeyInventory: stubArtifact(
				drydock.ArtifactKeyInventory,
				"decision-001-014-r1",
				"authorized TPM role inventory and policy\n",
			),
			IssuerHierarchy: stubArtifact(
				drydock.ArtifactIssuerHierarchy,
				"decision-004-r1",
				"authorized issuer hierarchy\n",
			),
			IngressTokenFlow: stubArtifact(
				drydock.ArtifactIngressTokenFlow,
				"decision-007-r1",
				"authorized ingress and token flow\n",
			),
			Compatibility: stubArtifact(
				drydock.ArtifactCompatibility,
				"compatibility-r1",
				"pinned upstream compatibility matrix\n",
			),
		},
	}
}
