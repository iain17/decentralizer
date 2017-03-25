package decentralizer

import(
	"testing"
	"gopkg.in/stretchr/testify.v1/assert"
)

func TestAddService(t *testing.T) {
	d, err := New()
	assert.NoError(t, err)
	err = d.AddService("iain", 0)
	assert.NoError(t, err)

	service := d.GetService("iain")
	assert.NotNil(t, service)
	service.SetDetail("dedicated", "yes")
	select{}
}
