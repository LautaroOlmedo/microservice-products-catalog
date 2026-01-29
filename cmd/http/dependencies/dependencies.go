package dependencies

import (
	"fmt"
	"microservice-products-catalog/cmd/http/config"
	"microservice-products-catalog/cmd/http/handlers/reader"
	"microservice-products-catalog/cmd/http/handlers/writer"
	my_sql "microservice-products-catalog/internal/infraestructure/my-sql"
	"microservice-products-catalog/internal/infraestructure/security/jwt"
	"microservice-products-catalog/internal/service/order"
	"microservice-products-catalog/internal/service/product"
	"time"
)

type Dependencies struct {
	TokenGenerator reader.TokenGenerator
	WriterHandler  writer.WriteHandler
	ReaderHandler  reader.ReaderHandler
}

func InitDependencies(cfg config.Config) Dependencies {
	// repository layer
	mySQLRepo, err := my_sql.NewRepository(cfg)
	if err != nil {
		panic(fmt.Sprintf("failed to connect mysql: %s", err.Error()))
	}
	txManager := my_sql.NewTxManager(mySQLRepo.DB())

	tokenGenerator := jwt.NewTokenGenerator(cfg.JWT.Secret, 15*time.Minute)

	// service layer
	productsService := product.NewService(mySQLRepo, txManager)
	ordersService := order.NewService(mySQLRepo, txManager, productsService)

	// handler layer
	writerHandler := writer.NewWriteHandler(productsService, ordersService)
	readerHandler := reader.NewReaderHandler(productsService, ordersService, tokenGenerator)

	return Dependencies{
		WriterHandler: *writerHandler,
		ReaderHandler: *readerHandler,
	}

}
