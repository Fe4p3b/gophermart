package order

import "go.uber.org/zap"

var _ Order = (*order)(nil)

type Order interface {
	List() ([]Order, error)
	addBonus(id string) error
}

type order struct {
	l *zap.SugaredLogger
}

func New(l *zap.SugaredLogger) *order {
	return &order{l: l}
}

func (o *order) List() ([]Order, error) {
	return []Order{}, nil
}

func (o *order) addBonus(id string) error {
	return nil
}
