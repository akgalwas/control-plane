package installation

import (
	"io/ioutil"

	"github.com/kyma-project/control-plane/components/provisioner/internal/model"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

const (
	defaultKymaNamespace = "kyma-system"
)

type Component struct {
	Name      string
	Namespace string
}
type List struct {
	DefaultNamespace string
	Prerequisites    []Component
	Components       []Component
}

func CreateFile(filename string, components []model.KymaComponentConfig) error {
	var cl List

	for _, component := range components {
		if component.Prerequisite {
			cl.Prerequisites = append(cl.Prerequisites, Component{
				Name:      string(component.Component),
				Namespace: component.Namespace,
			})
			continue
		}
		cl.Components = append(cl.Components, Component{
			Name:      string(component.Component),
			Namespace: component.Namespace,
		})
	}
	cl.DefaultNamespace = defaultKymaNamespace

	content, err := yaml.Marshal(&cl)
	if err != nil {
		return errors.Wrap(err, "while marshaling List")
	}

	err = ioutil.WriteFile(filename, content, 0644)
	if err != nil {
		return errors.Wrap(err, "while writing content to file")
	}

	return nil
}
