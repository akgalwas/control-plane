package installation

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/kyma-incubator/hydroform/install/installation"
	"github.com/kyma-incubator/hydroform/parallel-install/pkg/config"
	"github.com/kyma-incubator/hydroform/parallel-install/pkg/deployment"
	"github.com/kyma-incubator/hydroform/parallel-install/pkg/metadata"
	"github.com/kyma-incubator/hydroform/parallel-install/pkg/preinstaller"

	"github.com/kyma-project/control-plane/components/provisioner/internal/model"

	"github.com/avast/retry-go"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type parallelInstallationService struct {
}

func NewParallelInstallationService() Service {
	return &parallelInstallationService{}
}

func (p parallelInstallationService) CheckInstallationState(kubeconfig *rest.Config) (installation.InstallationState, error) {
	var (
		metadataName      = "kyma"
		metadataNamespace = "kyma-system"
		description       = ""
	)
	state := installation.InstallationState{State: installation.NoInstallationState}

	kubeClient, err := kubernetes.NewForConfig(kubeconfig)
	if err != nil {
		return state, errors.Wrap(err, "while creating kubernetes client")
	}

	cm, err := kubeClient.CoreV1().ConfigMaps(metadataNamespace).Get(context.TODO(), metadataName, metav1.GetOptions{})
	if err != nil {
		logrus.Errorf("cannot get %s/%s ConfigMap: %s", metadataName, metadataNamespace, err)
		return state, nil
	}

	stateRaw, ok := cm.Data["status"]
	if !ok {
		return state, errors.New("field status does not exist in CM with Kyma installation status")
	}

	if reason, ok := cm.Data["reason"]; ok {
		description = reason
	}

	// TODO method should return directly metadata.StatusEnum not installation.InstallationState
	switch metadata.StatusEnum(stateRaw) {
	case metadata.Deployed:
		state.State = "Installed"
		state.Description = description
		return state, nil
	case metadata.DeploymentInProgress:
		// TODO: handle in progress state better (e.g. from hydroform state)
		state.State = "InProgress"
		state.Description = description
		return state, nil
	case metadata.DeploymentError:
		return state, errors.Errorf("kyma deployment failed: %s", description)
	}

	return state, errors.Errorf("not supported Kyma installation status: %s", stateRaw)
}

func (p parallelInstallationService) TriggerInstallation(kubeconfigRaw *rest.Config, kymaProfile *model.KymaProfile, release model.Release, globalConfig model.Configuration, componentsConfig []model.KymaComponentConfig) error {
	kubeClient, err := kubernetes.NewForConfig(kubeconfigRaw)
	if err != nil {
		return errors.Wrap(err, "while creating kubernetes client")
	}

	dynamicClient, err := dynamic.NewForConfig(kubeconfigRaw)
	if err != nil {
		return errors.Wrap(err, "while creating dynamic client")
	}

	// prepare installation
	builder := &deployment.OverridesBuilder{}
	err = SetOverrides(builder, componentsConfig, globalConfig)
	if err != nil {
		return errors.Wrap(err, "while set overrides to the OverridesBuilder")
	}

	componentFilepath := fmt.Sprintf("/app/files/component-%s.yaml", uuid.New().String())
	err = CreateFile(componentFilepath, componentsConfig)
	if err != nil {
		return errors.Wrap(err, "while creating component file path")
	}

	installationCfg := &config.Config{
		WorkersCount:                  4,
		CancelTimeout:                 20 * time.Minute,
		QuitTimeout:                   25 * time.Minute,
		HelmTimeoutSeconds:            60 * 8,
		BackoffInitialIntervalSeconds: 3,
		BackoffMaxElapsedTimeSeconds:  60 * 5,
		Log:                           logrus.New(),
		HelmMaxRevisionHistory:        10,
		Profile:                       string(*kymaProfile),
		ComponentsListFile:            componentFilepath,
		ResourcePath:                  "/app/kyma-master/resources",
		InstallationResourcePath:      "/app/kyma-master/installation/resources",
		Version:                       release.Version,
	}

	cfg := preinstaller.Config{
		InstallationResourcePath: installationCfg.InstallationResourcePath,
		Log:                      installationCfg.Log,
	}

	resourceParser := &preinstaller.GenericResourceParser{}
	resourceManager := preinstaller.NewDefaultResourceManager(dynamicClient, cfg.Log, []retry.Option{})
	resourceApplier := preinstaller.NewGenericResourceApplier(installationCfg.Log, resourceManager)
	preInstaller := preinstaller.NewPreInstaller(resourceApplier, resourceParser, cfg, dynamicClient, []retry.Option{})

	// Install CRDs and create namespace
	result, err := preInstaller.InstallCRDs()
	if err != nil || len(result.NotInstalled) > 0 {
		return errors.Wrap(err, "while installing CRDs")
	}

	result, err = preInstaller.CreateNamespaces()
	if err != nil || len(result.NotInstalled) > 0 {
		return errors.Wrap(err, "while creating namespace")
	}

	// Install Kyma
	var progressCh chan deployment.ProcessUpdate
	deployer, err := deployment.NewDeployment(installationCfg, builder, kubeClient, progressCh)
	if err != nil {
		return errors.Wrap(err, "while creating deployer")
	}

	err = deployer.StartKymaDeployment()
	if err != nil {
		return errors.Wrap(err, "while starting Kyma installation")
	}
	logrus.Infof("Kyma installation started")

	// remove component file
	err = os.Remove(componentFilepath)
	if err != nil {
		return errors.Wrap(err, "while removing component file path")
	}

	return nil
}

func (p parallelInstallationService) TriggerUpgrade(_ *rest.Config, _ *model.KymaProfile, _ model.Release, _ model.Configuration, _ []model.KymaComponentConfig) error {
	panic("TriggerUpgrade is not implemented")
}

func (p parallelInstallationService) TriggerUninstall(_ *rest.Config) error {
	panic("TriggerUninstall is not implemented")
}

func (p parallelInstallationService) PerformCleanup(_ *rest.Config) error {
	panic("PerformCleanup is not implemented ")
}
