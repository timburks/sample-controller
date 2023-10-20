package main

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
	cid "github.com/ipfs/go-cid"
	mc "github.com/multiformats/go-multicodec"
	mh "github.com/multiformats/go-multihash"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

var bucketName = "nerdvana.timbx.me"
var serviceAccountFile = ".config/wdyk/control.json"

func main() {
	router := gin.Default()
	router.GET("/v1/blobs", GetBlobs)
	router.GET("/v1/blobs/:id", GetBlobById)
	router.POST("/v1/blobs", CreateBlob)
	router.Run("0.0.0.0:8080")
}

// GetBlobs responds with the list of all blobs as JSON.
func GetBlobs(c *gin.Context) {
	ctx := context.Background()
	client, err := GetClient(ctx, storage.ScopeFullControl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	srv, err := storage.NewClient(ctx, option.WithHTTPClient(client))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	bkt := srv.Bucket(bucketName)
	query := &storage.Query{Prefix: ""}
	var names []string
	it := bkt.Objects(ctx, query)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		names = append(names, attrs.Name)
	}
	c.IndentedJSON(http.StatusOK, names)
}

// CreateBlob adds a blob from bytes received in the request body.
func CreateBlob(c *gin.Context) {
	ctx := context.Background()
	client, err := GetClient(ctx, storage.ScopeFullControl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	srv, err := storage.NewClient(ctx, option.WithHTTPClient(client))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	bytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	bkt := srv.Bucket(bucketName)

	// Create a cid manually by specifying the 'prefix' parameters
	pref := cid.Prefix{
		Version:  1,
		Codec:    uint64(mc.Raw),
		MhType:   mh.SHA1, // MD5, //SHA2_256,
		MhLength: -1,      // default length
	}

	// And then feed it some data
	chunk, err := pref.Sum(bytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	name := chunk.String()
	obj := bkt.Object(name)
	objAttrs, err := obj.Attrs(ctx)

	// object exists
	if err == nil {
		if objAttrs.Size != int64(len(bytes)) {
			err := errors.New("size mismatch, hash collision?")
			c.JSON(http.StatusInternalServerError, err)
			return
		}
	} else if err != storage.ErrObjectNotExist {
		c.JSON(http.StatusInternalServerError, err)
		return
	} else {
		// Write bytes to obj.
		log.Printf("writing %d bytes", len(bytes))
		w := obj.NewWriter(ctx)
		defer w.Close()
		_, err = w.Write(bytes)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
	}

	c.IndentedJSON(http.StatusCreated, name)
}

// GetBlobById locates the blob whose ID value matches the id
// parameter sent by the client, then returns that blob as a response.
func GetBlobById(c *gin.Context) {
	ctx := context.Background()
	client, err := GetClient(ctx, storage.ScopeFullControl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	srv, err := storage.NewClient(ctx, option.WithHTTPClient(client))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	id := c.Param("id")
	bkt := srv.Bucket(bucketName)
	obj := bkt.Object(id)
	r, err := obj.NewReader(ctx)
	if err != nil {
		if err == storage.ErrObjectNotExist {
			c.JSON(http.StatusNotFound, err)
			return
		}
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	defer r.Close()
	bytes, err := io.ReadAll(r)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.Data(http.StatusOK, "application/octet-stream", bytes)
}

func GetClient(ctx context.Context, scope string) (*http.Client, error) {
	var client *http.Client
	if serviceAccountFile != "" {
		dirname, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		data, err := os.ReadFile(path.Join(dirname, serviceAccountFile))
		if err != nil {
			return nil, err
		}
		creds, err := google.CredentialsFromJSON(ctx, data, scope)
		if err != nil {
			return nil, err
		}
		client = oauth2.NewClient(ctx, creds.TokenSource)
	}
	return client, nil
}
