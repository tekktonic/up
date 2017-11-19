package main;

import (
	"fmt"
	"net/http"
	"strings"
	"log"
	"io/ioutil"
	//	"github.com/antchfx/xpath"
	"github.com/antchfx/xquery/xml"
	"github.com/husobee/vestigo"
	"github.com/kr/pretty"
)

func getTemplate(tree *xmlquery.Node) string {
	query := "//XRD/Link[@type=\"application/xrd+xml\"]/@template"
	//		pretty.Println((string)(body))
	//		fmt.Println(tree)

	result := xmlquery.FindOne(tree, query)
	if (result != nil) {}
	fmt.Println("RESULT ISSSSSS")
	pretty.Println(result)
	template := getAttr(result, "", "template")
	return template
}


type FingerInfo struct {
	id string;
	canonical string;
	atomfeed string;
//	magic string;
//	salmon string;
	feed string;
}

func getAtomFeed(feedurl string) string {
	// Yes, Get follows redirects.
	resp, _ := http.Get(feedurl)

	
	//		body, err := ioutil.ReadAll(resp.Body)
	tree, _ := xmlquery.Parse(resp.Body)

	result := xmlquery.FindOne(tree, "//service/workspace/collection")
	return getAttr(result, "", "href")
}


func getCanonical(tree *xmlquery.Node) string {
	// Yes, Get follows redirects.
	result := xmlquery.Find(tree, "//XRD/Alias")[0]
	return result.InnerText()
}
func getFingerInfo(name string) FingerInfo {

	fmt.Println("Getting finger info for " + name)
	id := strings.Split(name, "@")
	site := id[1]


	host_meta := "http://" + site + "/.well-known/host-meta"

	// Yes, Get follows redirects.
	resp, err := http.Get(host_meta)
	fmt.Println("Fetching " + host_meta)
	
	//		body, err := ioutil.ReadAll(resp.Body)
	tree, err := xmlquery.Parse(resp.Body)
	
	if (err != nil) {
		log.Fatal(err)
	}		

	template := getTemplate(tree);
	fmt.Println("URL Template is " + template)
	// Find our atom feed
	url := strings.Replace(template, "{uri}", name, 1)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/xrd+xml");
	resp, err = client.Do(req);
	tree, err = xmlquery.Parse(resp.Body)
	fmt.Println("Grabbed xml for webfinger profile")
	feed := getFeedURL(tree)
	atomfeed := getAtomFeed(feed)
	canonical := getCanonical(tree);

	fmt.Println("atom feed:" + atomfeed + "\nid: " + name)
	return FingerInfo{feed: feed, atomfeed: atomfeed, canonical: canonical, id: name}
//	feedurl, err := http.Get(atomfeed)
//	feed, err := ioutil.ReadAll(feedurl.Body)
}


// Given a webfinger document, extract the atom feed.
func getFeedURL(tree *xmlquery.Node) string {
	result := xmlquery.FindOne(tree, "//XRD/Link[@type=\"application/atomsvc+xml\"]")
	feed := getAttr(result, "", "href")

	return feed
}

func webfingerCB(w http.ResponseWriter, r *http.Request) {
	remote := vestigo.Param(r, "remote")
	fmt.Println("REMOTE IS " + remote)
	id := strings.Split(remote,"@")
	
	str := auth(r.Header.Get("X-Up-Auth"))

	

	if (str == "") {
		if (len(id) != 2) {
			w.WriteHeader(400);
			w.Write(([]byte)("Invalid webfinger ID"))
		}
		fmt.Println(remote + " " + id[1]);
		resp, err := http.Get(id[1] + "/.well-known/host-meta");
		if (err != nil) {
			w.WriteHeader(400);
			w.Write(([]byte)("Remote site failed to respond to host-meta"));
		}

		

		body, _ := ioutil.ReadAll(resp.Body)
		w.Write(([]byte)(body))

		
		return
		
	}

	w.WriteHeader(401)
	w.Write(([]byte)(str))
}




func HostMetaCB(w http.ResponseWriter, r *http.Request) {
	hostmetafile := `<?xml version="1.0" encoding="UTF-8"?>
     <XRD xmlns="http://docs.oasis-open.org/ns/xri/xrd-1.0">
        <Link rel="lrdd" type="application/xrd+xml" template="https://` + config.Server + `/.well-known/webfinger?resource={uri}"/>
     </XRD>
`
	w.Write(([]byte)(hostmetafile))
}

func WebfingerCB(w http.ResponseWriter, r *http.Request) {
	fingerfile := `<?xml version="1.0" encoding="UTF-8"?>
<XRD xmlns="http://docs.oasis-open.org/ns/xri/xrd-1.0">
  <Subject>acct:` + config.Owner + "@" + config.Server + `</Subject>
  <Link rel="http://schemas.google.com/g/2010#updates-from" type="application/atom+xml" href="https://` + config.Server + `/atom/feed.atom"/>
</XRD>
`
	w.Write(([]byte)(fingerfile))
}
