package storage

import (
	"context"
	"mobingi/ocean/pkg/config"
	"mobingi/ocean/pkg/constants"
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

type ClusterMongo struct {
	cfg *config.Config
}

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

const (
	certsTableName     = "certs"
	kubeconfsTableName = "kubeconfs"
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
	c.cfg = cfg

	// TODO Alraedy exists validation
	err := c.CreateCerts()
	if err != nil {
		return err
	}

	caCert, _ := c.GetCert(utilcert.PathForCert(constants.CACertAndKeyBaseName))
	caKey, _ := c.GetCert(utilcert.PathForKey(constants.CACertAndKeyBaseName))

	err = c.CreateKubeconfs(caCert, caKey)
	if err != nil {
		return err
	}
	return nil
}

func (c *ClusterMongo) CreateCerts() error {
	// Get cert
	sans := c.cfg.GetSANs()
	certs, err := certs.CreatePKIAssets(c.cfg.AdvertiseAddress, c.cfg.PublicIP, sans)
	if err != nil {
		return err
	}
	var insertData []interface{}
	for name, cert := range certs {
		insertData = append(insertData, certType{Cluster: c.cfg.ClusterName, Name: name, Cert: cert})
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

func (c *ClusterMongo) AllCerts() (map[string][]byte, error) {
	collection := db.Collection(certsTableName)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	cur, err := collection.Find(ctx, bson.M{"cluster": c.cfg.ClusterName})
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

func (c *ClusterMongo) GetCert(name string) ([]byte, error) {
	collection := db.Collection(certsTableName)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	var result certType
	err := collection.FindOne(ctx, bson.M{"cluster": c.cfg.ClusterName, "name": name}).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result.Cert, nil
}

func (c *ClusterMongo) CreateKubeconfs(caCert []byte, caKey []byte) error {
	cert, err := pkiutil.ParseCertPEM(caCert)
	if err != nil {
		return err
	}
	key, err := pkiutil.ParsePrivateKeyPEM(caKey)
	if err != nil {
		return err
	}
	kubeconfs, err := kubeconf.CreateKubeconf(c.cfg, cert, key)
	if err != nil {
		return err
	}
	var insertData []interface{}
	for name, conf := range kubeconfs {
		insertData = append(insertData, kubeconfType{Cluster: c.cfg.ClusterName, Name: name, Conf: conf})
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

func (c *ClusterMongo) AllKubeconfs() (map[string][]byte, error) {
	collection := db.Collection(kubeconfsTableName)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	cur, err := collection.Find(ctx, bson.M{"cluster": c.cfg.ClusterName})
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

func (c *ClusterMongo) GetKubeconf(name string) ([]byte, error) {
	collection := db.Collection(kubeconfsTableName)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	var result kubeconfType
	err := collection.FindOne(ctx, bson.M{"cluster": c.cfg.ClusterName, "name": name}).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result.Conf, nil
}
