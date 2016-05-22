package carloscloud

import (
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/cloud/storage"
	"io"
	"net/http"
)

func init() {
	http.HandleFunc("/", handler)
}

const gcsBucket = "meme-1299.appspot.com"

type demo struct {
	ctx    context.Context
	res    http.ResponseWriter
	bucket *storage.BucketHandle
	client *storage.Client
}

func handler(res http.ResponseWriter, req *http.Request) {

	if req.URL.Path != "/" {
		http.NotFound(res, req)
		return
	}

	ctx := appengine.NewContext(req)

	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Errorf(ctx, "ERROR handler NewClient: ", err)
		return
	}
	defer client.Close()

	d := &demo{
		ctx:    ctx,
		res:    res,
		client: client,
		bucket: client.Bucket(gcsBucket),
	}

	d.delFiles()
	d.createFiles()
	d.listFiles()
	io.WriteString(d.res, "\nFile Structure using /, aka windows/unix ( Delimeter: / )\n")
	d.listDelim()

}

func (d *demo) listDelim() {

	query := &storage.Query{
		Delimiter: "/",
	}

	objs, err := d.bucket.List(d.ctx, query)
	if err != nil {
		log.Errorf(d.ctx, "listBucketDirMode: unable to list bucket %q: %v", gcsBucket, err)
		return
	}

	for _, obj := range objs.Results {
		fmt.Fprintf(d.res, "%v\n", obj.Name)
	}

	fmt.Fprintf(d.res, "\nPREFIXES ( storage.ObjectList Prefixes )\n%v", objs.Prefixes)
}

func (d *demo) listFiles() {
	io.WriteString(d.res, "every file individually\n")

	objs, err := d.bucket.List(d.ctx, nil)
	if err != nil {
		log.Errorf(d.ctx, "%v", err)
		return
	}

	for _, obj := range objs.Results {
		io.WriteString(d.res, obj.Name+"\n")
	}
}

func (d *demo) createFiles() {
	for _, n := range []string{"firstobject", "folder1/folder2/secondobject", "whatwhat/ok/yup", "allkindsofthings", "does-this-make-pseudo-folders"} {
		d.createFile(n)
	}
}

func (d *demo) createFile(fileName string) {

	wc := d.bucket.Object(fileName).NewWriter(d.ctx)
	wc.ContentType = "text/plain"

	if _, err := wc.Write([]byte("BLAHBLAHBLAHBLAH\n")); err != nil {
		log.Errorf(d.ctx, "File Creation is messed!! here's why: %q, file %q: %v", gcsBucket, fileName, err)
		return
	}
	if err := wc.Close(); err != nil {
		log.Errorf(d.ctx, "createFile: BUCKET WILL NOT CLOSE, this stuff: %q, file %q: %v", gcsBucket, fileName, err)
		return
	}
}

func (d *demo) delFiles() {
	objs, err := d.bucket.List(d.ctx, nil)
	if err != nil {
		log.Errorf(d.ctx, "%v", err)
		return
	}

	for _, obj := range objs.Results {
		if err := d.bucket.Object(obj.Name).Delete(d.ctx); err != nil {
			log.Errorf(d.ctx, "deleteFiles: Can not Delete Files %q, file %q: %v", d.bucket, obj.Name, err)
			return
		}
	}
}
