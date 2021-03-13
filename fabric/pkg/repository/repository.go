package repository

type FabricI interface {
	ProduceProduct(code, power int) error
	CompareQuantity (code, required int) (int, error)
}

type Repository struct {
	Fabric
}

