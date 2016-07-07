package main

//write by qibin(http://blog.csdn.net/qibin0506)
import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/russross/blackfriday"
	"github.com/toqueteos/webbrowser"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
	"text/template"
	"time"
)

const (
	create = "create" // create new file
	view   = "view"   // view as html
	build  = "build"  // convert markdown file to html file
	cate   = "cate"   // create category

	defaultContent = `### hello markdown`

	categoryJs = `
var defaultPageSize = 5;
var arr = eval('{{.}}');
function get(currentPage) {
	return getResult(currentPage, defaultPageSize);
}
function getResult(currentPage, pageSize) {
	currentPage = parseInt(currentPage);
	pageSize = parseInt(pageSize);
	var startIndex = (currentPage - 1) * pageSize;
	var endIndex = startIndex + pageSize;
	if (arr.length <= startIndex) { return null;}
	if (endIndex > arr.length) { endIndex = arr.length;}
	return arr.slice(startIndex, endIndex);
}
function pageCount() {
	return getPageCount(defaultPageSize);
}
function getPageCount(pageSize) {
	return Math.ceil(arr.length / pageSize);
}
function getQueryString(query) {
	var uri = window.location.search;
    var re = new RegExp("" +query+ "=([^&?]*)", "ig");
    return ((uri.match(re))?(uri.match(re)[0].substr(query.length+1)):null);
}
	`
)

// md -type create -file my
// md -type view -file my
// md -type build -file my (-tmpl template.html -author loader -datefmt 2016-06-28)
var tp *string = flag.String("type", create, "use create build or view")
var fileName *string = flag.String("file", "default", "the filename to create or build or view")
var tmpl *string = flag.String("tmpl", "", "the html template file you want to use")
var author *string = flag.String("author", "", "the author of this article")
var hlp *string = flag.String("help", "", "")

var mdDir string = "./raw/"
var htmlDir string = "./html/"
var dateFmt string = "2006-01-02 15:04:05"

func init() {
	flag.Parse()
}

func main() {
	if *hlp != "" {
		help(*hlp)
		return
	}

	checkType(*tp)
	checkDir()

	switch *tp {
	case view:
		viewContent(*fileName)
	case create:
		createContent(*fileName)
	case build:
		buildContent(*fileName)
	case cate:
		buildCategory()
	}
}

func help(arg string) {
	fmt.Print("Usage:\n")
	switch arg {
	case "type":
		fmt.Println(" the type you want to use, you can use 'create', 'build' or 'view' here")
	case "file":
		fmt.Println(" the file name you want to operate")
	case "author":
		fmt.Println(" the author of this article")
	case "tmpl":
		fmt.Println(" the html template file you want to use when build markdown file")
	case "create":
		fmt.Println(" create a new markdown file \n e.g. tt -type create -file fileName")
	case "build":
		fmt.Println(" convert a html file from a markdown file \n e.g. tt -type build -file fileName -author qibin -tmpl ./content.html")
	case "view":
		fmt.Println(" view the html file you builded \n e.g. tt -type view -file fileName")
	case "detail":
		fallthrough
	default:
		fmt.Println(" 1. use 'tt -type create -file fileName' to create a new file")
		fmt.Println(" 2. use 'tt -type build -file fileName' to convert a markdown file to html")
		fmt.Println(" 3. use 'tt -type view -file fileName' to view a html file in browser")
	}
}

// 生成目录js
func buildCategory() {
	jsFile := htmlDir + "category.auto.js"
	deleteOldFile(jsFile)

	dirs, err := ioutil.ReadDir(htmlDir)
	checkError(err)

	var cates CategorySlice
	for _, item := range dirs {
		title := item.Name()
		if strings.HasSuffix(title, ".html") {
			title = strings.TrimSuffix(title, ".html")
		} else if strings.HasSuffix(title, ".htm") {
			title = strings.TrimSuffix(title, ".htm")
		}

		cate := &Category{
			Title: title,
			date:  item.ModTime(),
			Date:  item.ModTime().Format(dateFmt),
			Desc:  getContentDesc(item)}
		cates = append(cates, cate)
	}
	sort.Sort(cates)

	jsonData, err := json.Marshal(cates)
	checkError(err)

	outFile, err := os.Create(jsFile)
	checkError(err)
	defer outFile.Close()
	t := template.Must(template.New("category.auto.js").Parse(categoryJs))
	t.Execute(outFile, string(jsonData))
	fmt.Println("category data '" + jsFile + "' build success! now you can use it in your html file")
}

// 生成html内容文件
func buildContent(name string) {
	mdFile := mdDir + name + ".md"
	if !fileExist(mdFile) {
		createContent(name)
		return
	}

	f, err := os.Open(mdFile)
	checkError(err)
	defer f.Close()
	fd, err := ioutil.ReadAll(f)
	checkError(err)

	html := blackfriday.MarkdownCommon(fd)
	art := &Article{
		Title:   name,
		Date:    time.Now().Format(dateFmt),
		Author:  *author,
		Content: string(html),
		Desc:    parseContentDesc(fd)}
	convertHtml(name, html, art)
}

// 查看html文件
func viewContent(name string) {
	file := htmlDir + name + ".html"
	if !fileExist(file) {
		if fileExist(mdDir + name + ".md") {
			buildContent(name)
		} else {
			createContent(name)
		}

		return
	}

	fmt.Println("wait a moment...")
	time.AfterFunc(time.Duration(1)*time.Second, func() {
		webbrowser.Open("http://127.0.0.1:8080/" + name + ".html")
	})
	HandleHttp()
}

// 创建markdown文件
func createContent(name string) {
	file := mdDir + name + ".md"
	deleteOldFile(file)

	defer func() {
		if err := recover(); err != nil {
			log.Fatal(err)
		}
	}()

	mdFile, err := os.Create(file)
	checkError(err)

	defer mdFile.Close()
	_, err = mdFile.WriteString(defaultContent)
	checkError(err)

	fmt.Println(file + " create successed! now you can edit it!!!")
}

func HandleHttp() {
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(htmlDir))))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

// 从html中提取description
func getContentDesc(fileInfo os.FileInfo) (desc string) {
	name := fileInfo.Name()
	info, err := ioutil.ReadFile(htmlDir + name)
	checkError(err)
	content := string(info)
	re := regexp.MustCompile(`(?i)<meta.*?name=["|\']description["|\'].*?content=["|\'](?P<desc>.*?)["|\']`)
	res := re.FindAllStringSubmatch(content, -1)
	if res != nil && len(res) > 0 {
		desc = res[0][len(res[0])-1]
	}
	return
}

// 从markdown中提取一行作为description
func parseContentDesc(content []byte) (desc string) {
	buf := bufio.NewScanner(strings.NewReader(string(content)))
	for buf.Scan() {
		line := buf.Text()
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}
		desc = line
		break
	}

	if desc != "" {
		html := blackfriday.MarkdownCommon([]byte(desc))
		re, err := regexp.Compile(`(?iU)<.*>|\s`)
		checkError(err)
		desc = re.ReplaceAllString(string(html), "")
	}

	return
}

// 将markdown内容转化为html，并写入模板中
func convertHtml(name string, html []byte, content interface{}) {
	htmlFile := htmlDir + name + ".html"
	deleteOldFile(htmlFile)

	if *tmpl == "" {
		// raw
		err := ioutil.WriteFile(htmlFile, html, 0777)
		checkError(err)
	} else {
		// template
		outFile, err := os.Create(htmlFile)
		checkError(err)
		defer outFile.Close()
		t := template.Must(template.ParseFiles(*tmpl))
		err = t.Execute(outFile, content)
		checkError(err)
	}

	var arg string
	scan(&arg, "build '"+htmlFile+"' success!!!\nview it now? (Y/N)")

	switch arg {
	case "y":
		fallthrough
	case "Y":
		viewContent(name)
	}
}

// 删除旧文件
func deleteOldFile(file string) {
	if !fileExist(file) {
		return
	}

	var arg string
	scan(&arg, "old file '"+file+"' exist!!! \n do you want to create a new one? (Y/N)")

	switch arg {
	case "n":
		fallthrough
	case "N":
		log.Fatal("exit by user")
	case "y":
		fallthrough
	case "Y":
		if err := os.Remove(file); err != nil {
			log.Fatal(err)
		}
	}
}

func scan(arg *string, notice string) {
	fmt.Print(notice)
	fmt.Scanln(arg)
}

func fileExist(file string) bool {
	_, err := os.Stat(file)
	if err == nil {
		return true
	}

	return os.IsExist(err)
}

func checkDir() {
	err := os.MkdirAll(htmlDir, 0777)
	checkError(err)
	err = os.MkdirAll(mdDir, 0777)
	checkError(err)
}

func checkType(tp string) {
	if tp != create && tp != view && tp != build && tp != cate {
		log.Fatal("type" + tp + " is unavailable")
	}
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type CategorySlice []*Category

func (this CategorySlice) Len() int {
	return len(this)
}

func (this CategorySlice) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

func (this CategorySlice) Less(i, j int) bool {
	return this[i].date.After(this[j].date)
}

// 文章
type Article struct {
	Title   string
	Date    string
	Author  string
	Content string
	Desc    string
}

// 目录
type Category struct {
	Title string `json:"title"`
	Date  string `json:"date"`
	Desc  string `json:"desc"`
	date  time.Time
}
