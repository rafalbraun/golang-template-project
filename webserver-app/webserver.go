package main

/*
https://golang.cafe/blog/golang-read-file-example.html
https://www.golangprograms.com/how-to-extract-text-from-between-html-tag-using-regular-expressions-in-golang.html
https://zetcode.com/golang/mysql/
https://tour.golang.org/moretypes/15
https://www.calhoun.io/concatenating-and-building-strings-in-go/
https://yourbasic.org/golang/regexp-cheat-sheet/
https://qvault.io/golang/replace-strings-golang/
https://zetcode.com/golang/mysql/


deploy:
https://hackersandslackers.com/deploy-golang-app-nginx/
https://firehydrant.io/blog/develop-a-go-app-with-docker-compose/



ibmcloud-first-app







*/

import (
	"net/http"
	"github.com/gorilla/mux"
	"io/ioutil"
	"fmt"
	"regexp"
	"strings"

    "database/sql"
    "log"

    _ "github.com/go-sql-driver/mysql"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

type QueryRow struct {
    RowString string
}

func executeSql (updateString string) {
    db, err := sql.Open("mysql", "gorm:gorm@tcp(godockerDB:3306)/gorm")
    defer db.Close()

    if err != nil {
        log.Fatal(err)
    }

    //sql := "INSERT INTO cities(name, population) VALUES ('Moscow', 12506000)"
    res, err := db.Exec(updateString)

    if err != nil {
        panic(err.Error())
    }

    lastId, err := res.LastInsertId()

    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("The last inserted row id: %d\n", lastId)

}

func runSql (queryString string) string {
    //db, err := sql.Open("mysql", "user7:s$cret@tcp(127.0.0.1:3306)/testdb")
    //db, err := sql.Open("mysql", "gorm:gorm@tcp(127.0.0.1:9910)/gorm")
    db, err := sql.Open("mysql", "gorm:gorm@tcp(godockerDB:3306)/gorm")

    defer db.Close()
	
    if err != nil {
        log.Fatal(err)
    }

    res, err := db.Query(queryString)

    defer res.Close()

    if err != nil {
        log.Fatal(err)
    }

	var sb strings.Builder

    for res.Next() {

        var row QueryRow
        err := res.Scan(&row.RowString)

        if err != nil {
            log.Fatal(err)
        }

        //fmt.Printf("%s\n", row.RowString)
		sb.WriteString (row.RowString)
		sb.WriteString ("\n")
    }
	//fmt.Println (sb.String());

	return sb.String()

/*
    var version string

    //err2 := db.QueryRow("SELECT VERSION()").Scan(&version)
	err2 := db.QueryRow("select CONCAT('<tr><td>',post_id,'</td>','<td>',post_content,'</td><td>',publish_date,'</td></tr>') from posts;").Scan(&version)

    if err2 != nil {
        log.Fatal(err2)
    }

    //fmt.Println(version)

*/
}

func getSql (htmlBody string) string {
    //fmt.Print(htmlBody)

	re := regexp.MustCompile(`<sql.*?>(.*)</sql>`)
/*
	submatchall := re.FindAllStringSubmatch(htmlBody, -1)
	for _, element := range submatchall {
		fmt.Println(element[1])
	}
*/
	queryOutput := re.FindString(htmlBody)

	return queryOutput[5:len(queryOutput)-6]
}

func getPost (w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    //id, _ := strconv.Atoi(vars["id"])
	id := vars["id"]
	fmt.Print (id)
}

func createPost (w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
	    dat, err := ioutil.ReadFile("create_post.html")
	    check(err)
	    w.Write([]byte(string(dat)))

    case http.MethodPost:
        contents := r.FormValue("post_contents")
		updateSql := fmt.Sprintf("insert into posts (post_content) values ('%s');", contents)
		fmt.Println (updateSql)
		executeSql (updateSql)
		http.Redirect(w, r, "/posts", 301)

    default:
        // Give an error message.
		fmt.Print ("unknown error \n")
    }
}

func status (w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}

func main() {
    r := mux.NewRouter()

    // IMPORTANT: you must specify an OPTIONS method matcher for the middleware to set CORS headers
    r.HandleFunc ("/posts", fooHandler).Methods(http.MethodGet, http.MethodPut, http.MethodPatch, http.MethodOptions)
    r.HandleFunc ("/post/{id:[0-9]+}", getPost).Methods(http.MethodGet)
    r.HandleFunc ("/post_create", createPost).Methods(http.MethodGet, http.MethodPost)
    r.HandleFunc ("/status", status).Methods(http.MethodGet)

    r.Use(mux.CORSMethodMiddleware(r))
    
    http.ListenAndServe(":8080", r)
}

func fooHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    if r.Method == http.MethodOptions {
        return
    }

    dat, err := ioutil.ReadFile("post.html")
    check(err)
	
	htmlBody := string(dat)
	queryString := getSql (htmlBody);

	queryOutput := runSql (queryString)

	x := strings.Replace (htmlBody, queryString, queryOutput, -1);
	//fmt.Println (x)

	fmt.Printf ("queryString %s \n", queryString)
	//fmt.Printf ("queryOutput %s \n", queryOutput)

    w.Write([]byte(x))
}

