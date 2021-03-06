package fitasks

import (
	"fmt"
	"k8s.io/kops/upup/pkg/fi"
)

//go:generate fitask -type=Secret
type Secret struct {
	Name *string
}

var _ fi.HasCheckExisting = &Secret{}

// It's important always to check for the existing Secret, so we don't regenerate tokens e.g. on terraform
func (e *Secret) CheckExisting(c *fi.Context) bool {
	return true
}

func (e *Secret) Find(c *fi.Context) (*Secret, error) {
	secrets := c.SecretStore

	name := fi.StringValue(e.Name)
	if name == "" {
		return nil, nil
	}

	secret, err := secrets.FindSecret(name)
	if err != nil {
		return nil, err
	}
	if secret == nil {
		return nil, nil
	}

	actual := &Secret{
		Name: &name,
	}

	return actual, nil
}

func (e *Secret) Run(c *fi.Context) error {
	return fi.DefaultDeltaRunMethod(e, c)
}

func (s *Secret) CheckChanges(a, e, changes *Secret) error {
	if a != nil {
		if changes.Name != nil {
			return fi.CannotChangeField("Name")
		}
	}
	return nil
}

func (_ *Secret) Render(c *fi.Context, a, e, changes *Secret) error {
	name := fi.StringValue(e.Name)
	if name == "" {
		return fi.RequiredField("Name")
	}

	secrets := c.SecretStore

	_, _, err := secrets.GetOrCreateSecret(name)
	if err != nil {
		return fmt.Errorf("error creating secret %q: %v", name, err)
	}

	return nil
}
