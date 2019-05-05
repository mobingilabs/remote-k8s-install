package storage

import (
	"context"
	"mobingi/ocean/pkg/config"
	"mobingi/ocean/pkg/constants"
	"mobingi/ocean/pkg/kubernetes/bootstrap"
	"mobingi/ocean/pkg/kubernetes/service"
	"mobingi/ocean/pkg/log"
	"mobingi/ocean/pkg/tools/certs"
	"mobingi/ocean/pkg/tools/kubeconf"
	utilcert "mobingi/ocean/pkg/util/certs"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	pkiutil "mobingi/ocean/pkg/util/pki"

	"github.com/mongodb/mongo-go-driver/mongo"
)

type ClusterMongo struct{}

type certType struct {
	Cluster string `bson:"cluster"`
	Name    string `bson:"name"`
	Cert    []byte `bson:"cert"`
}

type kubeconfType struct {
	Cluster string `bson:"cluster"`
	Name    string `bson:"name"`
	Conf    []byte `bson:"conf"`
}

type etcdServersType struct {
	Cluster string `bson:"cluster"`
	Servers string `bson:"servers"`
}

type clusterType struct {
	Name string `bson:"name"`
}

const (
	clusterTableName     = "clusters"
	certsTableName       = "certs"
	kubeconfsTableName   = "kubeconfs"
	etcdServersTableName = "etcd_servers"
)

var db *mongo.Database

func NewMongoClient() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://admin:admin123456@localhost:27017"))
	if err != nil {
		log.Panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Panic(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Panic(err)
	}

	db = client.Database("k8s-cluster")
}

func (c *ClusterMongo) Init(cfg *config.Config) error {
	// TODO Alraedy exists validation
	err := c.CreateCerts(cfg)
	if err != nil {
		return err
	}
	caCert, err := c.GetCert(cfg.ClusterName, utilcert.PathForCert(constants.CACertAndKeyBaseName))
	if err != nil {
		return err
	}
	caKey, err := c.GetCert(cfg.ClusterName, utilcert.PathForKey(constants.CACertAndKeyBaseName))
	if err != nil {
		return err
	}
	err = c.CreateKubeconfs(cfg, caCert, caKey)
	if err != nil {
		return err
	}
	err = c.SetEtcdServers(cfg)
	if err != nil {
		return err
	}
	err = c.setCluster(cfg)
	if err != nil {
		return err
	}
	return nil
}

func (c *ClusterMongo) Exist(name string) (bool, error) {
	collection := db.Collection(clusterTableName)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	var result bson.M
	err := collection.FindOne(ctx, bson.M{"name": name}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return false, nil
	}
	if err != nil {
		return true, err
	}
	return true, nil
}

func (c *ClusterMongo) Drop(cfg *config.Config) error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	kubeconf := db.Collection(kubeconfsTableName)
	err := kubeconf.Drop(ctx)
	if err != nil {
		return err
	}
	certs := db.Collection(certsTableName)
	err = certs.Drop(ctx)
	if err != nil {
		return err
	}
	etcd := db.Collection(etcdServersTableName)
	err = etcd.Drop(ctx)
	if err != nil {
		return err
	}
	cluster := db.Collection(clusterTableName)
	err = cluster.Drop(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (c *ClusterMongo) setCluster(cfg *config.Config) error {
	// Store etcd servers in the database
	collection := db.Collection(clusterTableName)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	insertData := clusterType{
		Name: cfg.ClusterName,
	}
	_, err := collection.InsertOne(ctx, insertData)
	if err != nil {
		return err
	}
	return nil
}

func (c *ClusterMongo) SetEtcdServers(cfg *config.Config) error {
	privateIPs := cfg.GetMasterPrivateIPs()
	etcdServers := service.GetEtcdServers(privateIPs)
	// Store etcd servers in the database
	collection := db.Collection(etcdServersTableName)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	insertData := etcdServersType{
		Cluster: cfg.ClusterName,
		Servers: etcdServers,
	}
	_, err := collection.InsertOne(ctx, insertData)
	if err != nil {
		return err
	}
	return nil
}

func (c *ClusterMongo) GetEtcdServers(clusterName string) (string, error) {
	collection := db.Collection(etcdServersTableName)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	var result etcdServersType
	err := collection.FindOne(ctx, bson.M{"cluster": clusterName}).Decode(&result)
	if err != nil {
		return "", err
	}
	return result.Servers, nil
}

func (c *ClusterMongo) CreateCerts(cfg *config.Config) error {
	// Get cert
	sans := cfg.GetSANs()
	certs, err := certs.CreatePKIAssets(cfg.AdvertiseAddress, cfg.PublicIP, sans)
	if err != nil {
		return err
	}
	var insertData []interface{}
	for name, cert := range certs {
		insertData = append(insertData, certType{Cluster: cfg.ClusterName, Name: name, Cert: cert})
	}
	// Store certs in the database
	collection := db.Collection(certsTableName)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err = collection.InsertMany(ctx, insertData, &options.InsertManyOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (c *ClusterMongo) AllCerts(clusterName string) (map[string][]byte, error) {
	collection := db.Collection(certsTableName)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	cur, err := collection.Find(ctx, bson.M{"cluster": clusterName})
	if err != nil {
		return nil, err
	}
	certs := make(map[string][]byte)
	for cur.Next(ctx) {
		var result certType
		err := cur.Decode(&result)
		if err != nil {
			return nil, err
		}
		certs[result.Name] = result.Cert
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return certs, nil
}

func (c *ClusterMongo) GetCert(clusterName, name string) ([]byte, error) {
	collection := db.Collection(certsTableName)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	var result certType
	err := collection.FindOne(ctx, bson.M{"cluster": clusterName, "name": name}).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result.Cert, nil
}

func (c *ClusterMongo) CreateKubeconfs(cfg *config.Config, caCert []byte, caKey []byte) error {
	cert, err := pkiutil.ParseCertPEM(caCert)
	if err != nil {
		return err
	}
	key, err := pkiutil.ParsePrivateKeyPEM(caKey)
	if err != nil {
		return err
	}
	kubeconfs, err := kubeconf.CreateKubeconf(cfg, cert, key)
	if err != nil {
		return err
	}
	var insertData []interface{}
	for name, conf := range kubeconfs {
		insertData = append(insertData, kubeconfType{Cluster: cfg.ClusterName, Name: name, Conf: conf})
	}
	// Store kubeconfs in the database
	collection := db.Collection(kubeconfsTableName)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err = collection.InsertMany(ctx, insertData, &options.InsertManyOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (c *ClusterMongo) AllKubeconfs(clusterName string) (map[string][]byte, error) {
	collection := db.Collection(kubeconfsTableName)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	cur, err := collection.Find(ctx, bson.M{"cluster": clusterName})
	if err != nil {
		return nil, err
	}
	kubeconfs := make(map[string][]byte)
	for cur.Next(ctx) {
		var result kubeconfType
		err := cur.Decode(&result)
		if err != nil {
			return nil, err
		}
		kubeconfs[result.Name] = result.Conf
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return kubeconfs, nil
}

func (c *ClusterMongo) GetKubeconf(clusterName, name string) ([]byte, error) {
	collection := db.Collection(kubeconfsTableName)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	var result kubeconfType
	err := collection.FindOne(ctx, bson.M{"cluster": clusterName, "name": name}).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result.Conf, nil
}

func (c *ClusterMongo) SetBootstrapConf(clusterName string) error {
	adminConf, err := c.GetKubeconf(clusterName, "admin.conf")
	bootstrapConf, err := bootstrap.Bootstrap(adminConf)
	if err != nil {
		return err
	}
	insertData := kubeconfType{Cluster: clusterName, Name: "bootstrap.conf", Conf: bootstrapConf}
	// Store kubeconfs in the database
	collection := db.Collection(kubeconfsTableName)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err = collection.InsertOne(ctx, insertData)
	if err != nil {
		return err
	}
	return nil
}
