package table

import "context"

type Storage interface {
	Create(ctx context.Context, dto RecordDTO) (string, error)
	CreateCollection(ctx context.Context, dto CollectionDTO) error
	FindById(ctx context.Context, dto RecordDTO) (Record, error)
	FindByUserId(ctx context.Context, dto RecordDTO) (Record, error)
	FindAll(ctx context.Context, dto RecordDTO) ([]Record, error)
	FindCollectionNames(ctx context.Context) ([]string, error)
	Update(ctx context.Context, dto RecordDTO) error
	Delete(ctx context.Context, dto RecordDTO) error
	DeleteCollectionByName(ctx context.Context, dto CollectionDTO) error
}
