package provisioning

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/pkg/errors"

	"k8s.io/apimachinery/pkg/runtime"
	clientgotesting "k8s.io/client-go/testing"

	"github.com/kyma-project/control-plane/components/provisioner/internal/apperrors"
	"github.com/kyma-project/control-plane/components/provisioner/internal/model"
	provisioning_mocks "github.com/kyma-project/control-plane/components/provisioner/internal/operations/stages/provisioning/mocks"
	"github.com/kyma-project/control-plane/components/provisioner/internal/util"
	"github.com/kyma-project/control-plane/components/provisioner/internal/util/k8s/mocks"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

const dynamicKubeconfig = "dynamic_kubeconfig"
const testAdministrator1 = "testadmin1"
const testAdministrator2 = "testadmin2"

func TestCreateBindingsForOperatorsStep_Run(t *testing.T) {

	cluster := model.Cluster{
		Kubeconfig: util.PtrTo("kubeconfig"),
		ClusterConfig: model.GardenerConfig{
			Name: "shoot",
		},
		Administrators: []string{testAdministrator1, testAdministrator2},
	}

	operatorBindingConfig := OperatorRoleBinding{
		CreatingForAdmin: true,
	}

	dynamicKubeconfigProvider := &provisioning_mocks.DynamicKubeconfigProvider{}
	dynamicKubeconfigProvider.On("FetchFromRequest", "shoot").Return([]byte("dynamic_kubeconfig"), nil)

	t.Run("should return next step when finished", func(t *testing.T) {
		// given
		k8sClient := fake.NewSimpleClientset()
		k8sClientProvider := &mocks.K8sClientProvider{}
		k8sClientProvider.On("CreateK8SClient", dynamicKubeconfig).Return(k8sClient, nil)

		step := NewCreateBindingsForOperatorsStep(k8sClientProvider, operatorBindingConfig, dynamicKubeconfigProvider, nextStageName, time.Minute)

		// when
		result, err := step.Run(cluster, model.Operation{}, &logrus.Entry{})

		// then
		require.NoError(t, err)
		assert.Equal(t, nextStageName, result.Stage)
		assert.Equal(t, time.Duration(0), result.Delay)
	})

	t.Run("should not fail if cluster role binding already exists", func(t *testing.T) {
		// given
		k8sClient := fake.NewSimpleClientset()
		for i, administrator := range cluster.Administrators {
			clusterRoleBinding := buildClusterRoleBinding(
				fmt.Sprintf("%s%d", administratorOperatorClusterRoleBindingName, i),
				administrator,
				ownerClusterRoleBindingRoleRefName,
				userKindSubject,
				map[string]string{"app": "kyma"},
			)
			_, err := k8sClient.RbacV1().ClusterRoleBindings().Create(context.Background(), &clusterRoleBinding, metav1.CreateOptions{})
			require.NoError(t, err)
		}

		k8sClientProvider := &mocks.K8sClientProvider{}
		k8sClientProvider.On("CreateK8SClient", dynamicKubeconfig).Return(k8sClient, nil)

		step := NewCreateBindingsForOperatorsStep(k8sClientProvider, operatorBindingConfig, dynamicKubeconfigProvider, nextStageName, time.Minute)

		// when
		result, err := step.Run(cluster, model.Operation{}, &logrus.Entry{})

		// then
		require.NoError(t, err)
		assert.Equal(t, nextStageName, result.Stage)
		assert.Equal(t, time.Duration(0), result.Delay)
	})

	t.Run("should not fail if namespace istio-system already exists", func(t *testing.T) {
		// given
		k8sClient := fake.NewSimpleClientset()

		ns := &core.Namespace{
			ObjectMeta: metav1.ObjectMeta{Name: "istio-system"},
		}

		_, err := k8sClient.CoreV1().Namespaces().Create(context.Background(), ns, metav1.CreateOptions{})

		require.NoError(t, err)

		k8sClientProvider := &mocks.K8sClientProvider{}
		k8sClientProvider.On("CreateK8SClient", dynamicKubeconfig).Return(k8sClient, nil)

		step := NewCreateBindingsForOperatorsStep(k8sClientProvider, operatorBindingConfig, dynamicKubeconfigProvider, nextStageName, time.Minute)

		// when
		result, err := step.Run(cluster, model.Operation{}, &logrus.Entry{})

		// then
		require.NoError(t, err)
		assert.Equal(t, nextStageName, result.Stage)
		assert.Equal(t, time.Duration(0), result.Delay)
	})

	t.Run("should attempt retry when failed to get dynamic kubeconfig", func(t *testing.T) {
		// given
		dynamicKubeconfigProvider := &provisioning_mocks.DynamicKubeconfigProvider{}
		dynamicKubeconfigProvider.On("FetchFromRequest", "shoot").Return(nil, errors.New("some error"))

		step := NewCreateBindingsForOperatorsStep(nil, operatorBindingConfig, dynamicKubeconfigProvider, nextStageName, time.Minute)

		// when
		result, err := step.Run(cluster, model.Operation{}, &logrus.Entry{})

		// then
		require.NoError(t, err)
		assert.Equal(t, model.CreatingBindingsForOperators, result.Stage)
		assert.Equal(t, 20*time.Second, result.Delay)
	})

	t.Run("should return error when failed to provide k8s client", func(t *testing.T) {
		// given
		k8sClientProvider := &mocks.K8sClientProvider{}
		k8sClientProvider.On("CreateK8SClient", dynamicKubeconfig).Return(nil, apperrors.Internal("error"))

		step := NewCreateBindingsForOperatorsStep(k8sClientProvider, operatorBindingConfig, dynamicKubeconfigProvider, nextStageName, time.Minute)

		// when
		_, err := step.Run(cluster, model.Operation{}, &logrus.Entry{})

		// then
		require.Error(t, err)
	})

	t.Run("should return error when failed to create cluster role binding", func(t *testing.T) {
		// given
		k8sClient := fake.NewSimpleClientset()
		k8sClient.PrependReactor(
			"*",
			"*",
			func(action clientgotesting.Action) (handled bool, ret runtime.Object, err error) {
				return true, nil, fmt.Errorf("error")
			})

		k8sClientProvider := &mocks.K8sClientProvider{}
		k8sClientProvider.On("CreateK8SClient", dynamicKubeconfig).Return(k8sClient, nil)

		step := NewCreateBindingsForOperatorsStep(k8sClientProvider, operatorBindingConfig, dynamicKubeconfigProvider, nextStageName, time.Minute)

		// when
		_, err := step.Run(cluster, model.Operation{}, &logrus.Entry{})

		// then
		require.Error(t, err)
	})
}
