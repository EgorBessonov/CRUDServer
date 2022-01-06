package repository

import (
	"CRUDServer/internal/model"
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
)

// PostgresRepository type replies for accessing to postgres database
type PostgresRepository struct {
	DBconn *pgxpool.Pool
}

// Create save Order object into postgresql database
func (rps PostgresRepository) Create(ctx context.Context, order model.Order) error {
	result, err := rps.DBconn.Exec(ctx, "insert into orders (orderName, orderCost, isDelivered) "+
		"values ($1, $2, $3)", order.OrderName, order.OrderCost, order.IsDelivered)
	if err != nil {
		postgresOperationError(err, "Create()")
		return err
	}
	postgresOperationSuccess(result, "Create()")
	return nil
}

// Read returns Order object from postgresql database
// with selection by OrderID
func (rps PostgresRepository) Read(ctx context.Context, orderID string) (model.Order, error) {
	var order model.Order
	err := rps.DBconn.QueryRow(ctx, "select orderName, orderCost, isDelivered from orders "+
		"where orderID=$1", orderID).Scan(&order.OrderName, &order.OrderCost, &order.IsDelivered)
	if err != nil {
		postgresOperationError(err, "Read()")
		return model.Order{}, err
	}
	postgresOperationSuccess(nil, "Read()")
	return order, nil
}

// Update update Order object from postgresql database
// with selection by OrderID
func (rps PostgresRepository) Update(ctx context.Context, order model.Order) error {
	result, err := rps.DBconn.Exec(ctx, "update orders "+
		"set orderName=$2, orderCost=$3, isDelivered=$4"+
		"where orderID=$1", order.OrderID, order.OrderName, order.OrderCost, order.IsDelivered)
	if err != nil {
		postgresOperationError(err, "Update()")
		return err
	}
	postgresOperationSuccess(result, "Update()")
	return nil
}

// Delete delete Order object from postgresql database
// with selection by OrderID
func (rps PostgresRepository) Delete(ctx context.Context, orderID string) error {
	result, err := rps.DBconn.Exec(ctx, "delete from orders where orderID=$1", orderID)
	if err != nil {
		postgresOperationError(err, "Delete()")
		return err
	}
	postgresOperationSuccess(result, "Delete()")
	return nil
}

// CreateAuthUser method saves authentication info about user into
// postgres database
func (rps PostgresRepository) CreateAuthUser(ctx context.Context, authUser *model.AuthUser) error {
	result, err := rps.DBconn.Exec(ctx, "insert into authusers (username, email, password)"+
		"values($1, $2, $3)", authUser.UserName, authUser.Email, authUser.Password)
	if err != nil {
		postgresOperationError(err, "CreateAuthUser()")
		return err
	}
	postgresOperationSuccess(result, "CreateAuthUser()")
	return nil
}

// GetAuthUser method returns authentication info about user from
// postgres database with selection by email
func (rps PostgresRepository) GetAuthUser(ctx context.Context, email string) (model.AuthUser, error) {
	var authUser model.AuthUser
	err := rps.DBconn.QueryRow(ctx, "select useruuid, username, email, password from authusers "+
		"where email=$1", email).Scan(&authUser.UserUUID, &authUser.UserName, &authUser.Email, &authUser.Password)
	if err != nil {
		postgresOperationError(err, "GetAuthUser()")
		return model.AuthUser{}, err
	}
	postgresOperationSuccess(nil, "GetAuthUser()")
	return authUser, nil
}

// GetAuthUserByID method returns authentication info about user from
// postgres database with selection by id
func (rps PostgresRepository) GetAuthUserByID(ctx context.Context, userUUID string) (model.AuthUser, error) {
	var authUser model.AuthUser
	err := rps.DBconn.QueryRow(ctx, "select useruuid, username, email, password, refreshtoken from authusers "+
		"where useruuid=$1", userUUID).Scan(&authUser.UserUUID, &authUser.UserName, &authUser.Email, &authUser.Password, &authUser.RefreshToken)
	if err != nil {
		postgresOperationError(err, "GetAuthUserByID()")
		return model.AuthUser{}, err
	}
	postgresOperationSuccess(nil, "GetAuthUserByID()")
	return authUser, nil
}

// UpdateAuthUser is method to set refresh token into authuser info
func (rps PostgresRepository) UpdateAuthUser(ctx context.Context, email, refreshToken string) error {
	result, err := rps.DBconn.Exec(ctx, "update authusers "+
		"set refreshtoken=$2"+
		"where email=$1", email, refreshToken)
	if err != nil {
		postgresOperationError(err, "UpdateAuthUser()")
		return err
	}
	postgresOperationSuccess(result, "UpdateAuthUser()")
	return nil
}

// CloseDBConnection is using to close current postgres database connection
func (rps PostgresRepository) CloseDBConnection() error {
	rps.DBconn.Close()
	postgresOperationSuccess(nil, "CloseDBConnection")
	return nil
}

func postgresOperationError(err error, method string) {
	log.WithFields(log.Fields{
		"method": method,
		"status": "Operation failed.",
		"error":  err,
	}).Info("Postgres repository info.")
}

func postgresOperationSuccess(result pgconn.CommandTag, method string) {
	log.WithFields(log.Fields{
		"method": method,
		"status": "Operation ended successfully",
		"result": result,
	}).Info("Postgres repository info.")
}
