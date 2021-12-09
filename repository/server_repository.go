package repository

import (
	"context"
	"time"

	db "github.com/kharism/proman/connection"
	"github.com/kharism/proman/model"
	"github.com/kharism/proman/pkg/module/confaccessor"
	"github.com/kharism/proman/util"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IServer interface {
	Count(filters ...*bson.M) (int64, error)
	Delete(filters ...*bson.M) error
	DeleteByID(id string) error
	FindAll(param *BaseParam) ([]model.Server, error)
	//FindAtasan(id string) ([]model.User, error)
	//FindByID(id string) (model.User, error)
	//FindByUsername(username string) (model.User, error)
	Save(data model.Server) (model.Server, error)
}

type serverRepo struct {
	client func() (*mongo.Client, error)
}

func NewServer() IServer {
	return &serverRepo{
		client: db.NewClient,
	}
}
func (r *serverRepo) Count(filters ...*bson.M) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	mongoCli, err := r.client()
	defer cancel()
	if err != nil {
		return -1, err
	}
	defer mongoCli.Disconnect(ctx)
	serverModel := model.Server{}
	db := mongoCli.Database(viper.GetString("db"))
	coll := db.Collection(serverModel.TableName())
	if filters == nil {
		return coll.CountDocuments(ctx, bson.D{})
	}
	return coll.CountDocuments(ctx, filters)
}
func (r *serverRepo) DeleteByID(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	mongoCli, err := r.client()
	defer cancel()
	if err != nil {
		return err
	}
	defer mongoCli.Disconnect(ctx)
	serverModel := model.Server{}
	db := mongoCli.Database(viper.GetString("db"))
	coll := db.Collection(serverModel.TableName())
	_, err = coll.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
func (r *serverRepo) Delete(filters ...*bson.M) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	mongoCli, err := r.client()
	defer cancel()
	if err != nil {
		return err
	}
	defer mongoCli.Disconnect(ctx)
	serverModel := model.Server{}
	db := mongoCli.Database(viper.GetString("db"))
	coll := db.Collection(serverModel.TableName())
	if filters == nil {
		_, err := coll.DeleteMany(ctx, bson.D{})
		if err != nil {
			return err
		}
	}
	_, err = coll.DeleteMany(ctx, filters)
	if err != nil {
		return err
	}
	return nil
}
func (r *serverRepo) FindAll(param *BaseParam) ([]model.Server, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	mongoCli, err := r.client()
	defer cancel()
	if err != nil {
		return []model.Server{}, err
	}
	result := []model.Server{}
	data := model.Server{}
	defer mongoCli.Disconnect(ctx)
	db := mongoCli.Database(viper.GetString("db"))
	options := options.Find()
	if len(param.Orders) > 0 {
		orderE := []bson.E{}
		for _, val := range param.Orders {
			if val[0] == '-' {
				field := val[1:]
				orderE = append(orderE, bson.E{Key: field, Value: -1})
			} else {
				field := val
				orderE = append(orderE, bson.E{Key: field, Value: 1})
			}
		}
		options.SetSort(orderE)
	}
	if param.Limit > 0 {
		p := int64(param.Limit)
		options.Limit = &p
	} else {
		options.SetLimit(10)
	}
	if param.Skip > 0 {
		p := int64(param.Skip)
		options.Skip = &p
	}
	//options.SetSkip(0)

	currsor, err := db.Collection(data.TableName()).Find(ctx, param.Filter, options)
	if err != nil {
		return []model.Server{}, err
	}
	err = currsor.All(ctx, &result)
	if err != nil {
		return []model.Server{}, err
	}
	currsor.Close(ctx)
	return PopulateAccessor(result), nil
}
func PopulateAccessor(a []model.Server) []model.Server {
	for idx := range a {
		a[idx].Accessor = confaccessor.Registry[a[idx].AccessorType]
	}
	return a
}
func (r *serverRepo) Save(data model.Server) (model.Server, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	mongoCli, err := r.client()
	defer cancel()
	if err != nil {
		return model.Server{}, err
	}
	defer mongoCli.Disconnect(ctx)
	db := mongoCli.Database(viper.GetString("db"))
	if data.ID == "" {
		data.ID = util.RandString(23)
		insertOneRes, err := db.Collection(data.TableName()).InsertOne(ctx, data)
		if err != nil {
			return model.Server{}, err
		}
		data.SetID([]interface{}{insertOneRes.InsertedID})
	} else {
		upsert := true
		option := options.ReplaceOptions{Upsert: &upsert}
		_, err = db.Collection(data.TableName()).ReplaceOne(ctx, bson.M{"_id": data.ID}, data, &option)
		if err != nil {
			return model.Server{}, err
		}
	}

	return data, nil
}
