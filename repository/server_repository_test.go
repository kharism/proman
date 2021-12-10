package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	db "github.com/kharism/proman/connection"
	"github.com/kharism/proman/model"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
)

func TestCreateServer(t *testing.T) {
	serverRepo := NewServer()
	serverModel := model.Server{}

	Convey("Clean up", t, func() {
		cli1, err := db.NewClient()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, err = cli1.Database(viper.GetString("db")).Collection(serverModel.TableName()).DeleteMany(ctx, bson.M{})
		So(err, ShouldBeNil)

		Convey("Try Save", func() {
			for i := 0; i < 10; i++ {
				serverModel.ID = fmt.Sprintf("%d", i)
				serverModel.Name = fmt.Sprintf("SSDSDD_%d", i)
				serverModel.Path = "/dummy/path/"
				serverModel.AccessorType = "Dummy"
				serverModel.Url = fmt.Sprintf("asdasdqweqwe_%d", i)
				serverModel, err := serverRepo.Save(serverModel)
				So(err, ShouldBeNil)
				So(serverModel.ID, ShouldNotBeEmpty)
			}

			count, err := serverRepo.Count(nil)

			So(err, ShouldBeNil)
			So(count, ShouldEqual, 10)

			baseParam := BaseParam{}
			baseParam.Limit = 10
			allData, err := serverRepo.FindAll(&baseParam)

			So(err, ShouldBeNil)
			So(len(allData), ShouldEqual, 10)

			//try filter
			filter := bson.M{}
			filter["Name"] = "SSDSDD_1"

			baseParam.Filter = filter
			oneData, err := serverRepo.FindAll(&baseParam)

			So(err, ShouldBeNil)
			So(len(oneData), ShouldEqual, 1)
			So(oneData[0].Name, ShouldEqual, "SSDSDD_1")
			deleteKey := oneData[0].ID

			//test paging
			baseParam.Filter = bson.M{}
			baseParam.Limit = 2
			baseParam.Skip = 2
			baseParam.Orders = []string{"Address"}

			nextData, err := serverRepo.FindAll(&baseParam)
			So(err, ShouldBeNil)
			So(len(nextData), ShouldEqual, 2)
			//t.Log(nextData)
			So(nextData[0].Name, ShouldEqual, "SSDSDD_2")
			So(nextData[1].Name, ShouldEqual, "SSDSDD_3")

			pp, err := nextData[0].ReadData()
			So(err, ShouldBeNil)
			So(pp, ShouldEqual, fmt.Sprintf("dummy:///%s", nextData[0].Path))

			filter = bson.M{}
			filter["_id"] = deleteKey

			//test count
			count, err = serverRepo.Count(&filter)
			So(err, ShouldBeNil)
			So(count, ShouldEqual, 1)

			//test delete
			err = serverRepo.DeleteByID(deleteKey)
			So(err, ShouldBeNil)

			baseParam.Filter = filter
			baseParam.Skip = 0
			baseParam.Limit = 0
			oneData, err = serverRepo.FindAll(&baseParam)
			So(err, ShouldBeNil)
			So(len(oneData), ShouldEqual, 0)

			//test delete all
			filter = bson.M{}
			err = serverRepo.Delete(&filter)
			So(err, ShouldBeNil)

			count, err = serverRepo.Count(nil)
			So(err, ShouldBeNil)
			So(count, ShouldEqual, 0)
		})

	})
}
