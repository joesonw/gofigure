package realworld_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/joesonw/gofigure"
	"github.com/joesonw/gofigure/feature"
)

var _ = Describe("All Realworld Tests", func() {
	It("should resolve after load (can reference  another file)", func() {
		loader := gofigure.New().WithFeatures(feature.All...)
		Expect(loader.Load("app.yaml", []byte(`env: dev
port: 8080
host: localhost
listen: !tpl |
  {{ config "app.host" }}:{{ config "app.port" }}
name: !ref app.product
db_host: !ref storage.db.host
database: !tpl |
  mysql://{{ config "storage.db.user" }}:{{ config "storage.db.password" }}@{{ config "storage.db.host" }}:{{ config "storage.db.port" }}
`))).To(BeNil())
		Expect(loader.Load("storage/db.yaml", []byte(`host: localhost
port: 3306
user: root`))).To(BeNil())
		Expect(loader.Load("app.yaml", []byte(`env: prod 
product: test
port: 80
`))).To(BeNil())
		Expect(loader.Load("storage/db.yaml", []byte(`env: prod 
host: remote-address
password: supersecret
`))).To(BeNil())
		var listen, name, database, dbHost string
		Expect(loader.Get(context.Background(), "app.listen", &listen)).To(BeNil())
		Expect(loader.Get(context.Background(), "app.name", &name)).To(BeNil())
		Expect(loader.Get(context.Background(), "app.database", &database)).To(BeNil())
		Expect(loader.Get(context.Background(), "app.db_host", &dbHost)).To(BeNil())
		Expect(listen).To(Equal("localhost:80"))
		Expect(name).To(Equal("test"))
		Expect(database).To(Equal("mysql://root:supersecret@remote-address:3306"))
		Expect(dbHost).To(Equal("remote-address"))
	})
})
