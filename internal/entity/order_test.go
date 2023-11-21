package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIfItgetAnErrorIfIDsIsBlank(t *testing.T) {
	order := Order{}

	assert.Error(t, order.Validate(), "ID is required")
}

func TestIfPriceIsGreaterThan0(t *testing.T) {
	order := Order{
		ID: "1",
	}
	assert.Error(t, order.Validate(), "Price should be > than 0")
}

func TestIfTaxIsGreaterThan0(t *testing.T) {
	order := Order{
		ID:    "1",
		Price: 10,
	}
	assert.Error(t, order.Validate(), "Tax should be > than 0")
}

func TestFinalPrice(t *testing.T) {
	order := Order{
		ID:    "1",
		Price: 10,
		Tax:   1,
	}

	assert.NoError(t, order.Validate())
	assert.Equal(t, "1", order.ID)
	assert.Equal(t, 10.0, order.Price)
	assert.Equal(t, 1.0, order.Tax)
	order.CalculateFinalprice()
	assert.Equal(t, 11.0, order.FinalPrice)
}
