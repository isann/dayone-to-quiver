package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	Version  string
	Revision string
)

const (
	location = "Asia/Tokyo"
)

type QuiverEntryMeta struct {
	CreatedAt int64    `json:"created_at"`
	Tags      []string `json:"tags"`
	Title     string   `json:"title"`
	UpdatedAt int64    `json:"updated_at"`
	Uuid      string   `json:"uuid"`
}

func NewQuiverEntryMeta() *QuiverEntryMeta {
	return &QuiverEntryMeta{}
}

type QuiverEntryCell struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

func NewQuiverEntryCell() *QuiverEntryCell {
	return &QuiverEntryCell{}
}

type QuiverEntryContent struct {
	Title string            `json:"title"`
	Cells []QuiverEntryCell `json:"cells"`
}

func NewQuiverEntryContent() *QuiverEntryContent {
	return &QuiverEntryContent{}
}

func isExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

// 2014-01-21T09:47:15Z や 2018-10-31T10:18:54Z の文字列
func convertStringToTime(dateStr string) int64 {
	year, err := strconv.Atoi(dateStr[0:4])
	if err != nil {
		panic(err)
	}
	month, err := strconv.Atoi(dateStr[5:7])
	if err != nil {
		panic(err)
	}
	day, err := strconv.Atoi(dateStr[8:10])
	if err != nil {
		panic(err)
	}
	hour, err := strconv.Atoi(dateStr[11:13])
	if err != nil {
		panic(err)
	}
	min, err := strconv.Atoi(dateStr[14:16])
	if err != nil {
		panic(err)
	}
	sec, err := strconv.Atoi(dateStr[17:19])
	if err != nil {
		panic(err)
	}
	date := time.Date(year, time.Month(month), day, hour, min, sec, 0, time.Local)
	return date.Unix()
}

func exec(input, output string, _e interface{}) {
	entry := _e.(map[string]interface{})

	// テキストを抽出
	text := ""
	if _, ok := entry["text"]; ok {
		text = entry["text"].(string)
	} else if _, ok := entry["richText"]; ok {
		text = entry["richText"].(string)
	}

	// タグをリスト
	_t := entry["tags"]
	var tags []interface{}
	if _t != nil {
		tags = _t.([]interface{})
	}
	var tagArray []string
	var tagStr string
	for _, tag := range tags {
		tagStr = tag.(string)
		tagArray = append(tagArray, tagStr)
	}
	// タグが要素 0 だと正常に import されないので dummy 設定
	if len(tagArray) == 0 {
		tagArray = append(tagArray, "notag")
	}

	// Quiver ディレクトリつくる
	uuid := entry["uuid"].(string)
	quiverEntryDir := filepath.Join(output, fmt.Sprintf("%s.qvnote", uuid))
	err := os.MkdirAll(quiverEntryDir, 0777)
	if err != nil {
		log.Fatalln(err)
	}

	// resource あればディレクトリつくる
	var imgArray []string
	if _, ok := entry["photos"]; ok {
		resourcesDir := filepath.Join(quiverEntryDir, "resources")
		err := os.MkdirAll(resourcesDir, 0777)
		if err != nil {
			log.Fatalln(err)
		}
		photos := entry["photos"].([]interface{})
		dayonePhotoBasePath := filepath.Join(input, "photos")
		for _, _p := range photos {
			photo := _p.(map[string]interface{})
			id := photo["identifier"].(string)
			md5 := photo["md5"].(string)
			jpeg := fmt.Sprintf("%s.jpeg", md5)
			gif := fmt.Sprintf("%s.gif", md5)
			_1 := filepath.Join(dayonePhotoBasePath, jpeg)
			_2 := filepath.Join(dayonePhotoBasePath, gif)
			dayOnePhotoPath := _1
			filePath := jpeg
			if isExist(_2) {
				dayOnePhotoPath = _2
				filePath = gif
			}
			// resource をコピーし ID をつける
			copyFile(dayOnePhotoPath, filepath.Join(resourcesDir, filePath))
			dayoneImg := fmt.Sprintf(`![](dayone-moment://%s)`, id)
			// ![greenback.png](quiver-image-url/19FA8414EC71158F7C1C6A18E75AA0F1.png =480x480)
			quiverImg := fmt.Sprintf(`![](quiver-image-url/%s)`, filePath)
			imgArray = append(imgArray, dayoneImg)
			imgArray = append(imgArray, quiverImg)
			text = strings.Replace(text, dayoneImg, quiverImg, -1)
		}
	}

	dateString := entry["creationDate"].(string)
	creationTime := convertStringToTime(dateString)

	lines := strings.Split(text, "\n")
	title := lines[0]

	// create meta.json
	meta := NewQuiverEntryMeta()
	meta.CreatedAt = creationTime
	meta.Tags = tagArray
	meta.Title = title
	meta.UpdatedAt = creationTime
	meta.Uuid = uuid
	metaJsonPath := filepath.Join(quiverEntryDir, "meta.json")
	bytes, err := json.Marshal(meta)
	if err != nil {
		panic(err)
	}
	_f, err := os.Create(metaJsonPath)
	if err != nil {
		panic(err)
	}
	defer _f.Close()
	_f.Write(bytes)

	// create content.json
	cell := NewQuiverEntryCell()
	cell.Type = "markdown"
	cell.Data = text
	content := NewQuiverEntryContent()
	content.Title = title
	content.Cells = []QuiverEntryCell{*cell}
	contentJsonPath := filepath.Join(quiverEntryDir, "content.json")
	bytes, err = json.Marshal(content)
	if err != nil {
		panic(err)
	}
	_ff, err := os.Create(contentJsonPath)
	if err != nil {
		panic(err)
	}
	defer _ff.Close()
	_ff.Write(bytes)

}

func copyFile(srcPath, dstPath string) {
	//log.Println(srcPath, dstPath)
	src, err := os.Open(srcPath)
	if err != nil {
		panic(err)
	}
	defer src.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		panic(err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		panic(err)
	}
}

func main() {
	version := flag.Bool("v", false, "")
	inputFilePath := flag.String("i", "", "input file path")
	outputFilePath := flag.String("o", "", "output file path")
	flag.Parse()
	if *version {
		println("Application version : ", Version)
		println("Built revision : ", Revision)
		// println("Env name : ", env.NAME)
		os.Exit(0)
	}

	if inputFilePath == nil {
		println("Invalid parameter. Display help \"-h\"")
		os.Exit(1)
	}
	if !isExist(*inputFilePath) {
		println("Not found file.")
		os.Exit(1)
	}
	if outputFilePath == nil {
		println("Invalid parameter. Display help \"-h\"")
		os.Exit(1)
	}
	if isExist(*outputFilePath) {
		println("Already exist output dir.")
		os.Exit(1)
	}
	err := os.MkdirAll(*outputFilePath, 0777)
	if err != nil {
		panic(err)
	}

	// Logging 設定
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// time locale
	loc, err := time.LoadLocation(location)
	if err != nil {
		loc = time.FixedZone(location, 9*60*60)
	}
	time.Local = loc

	// load json
	file, err := os.Open(*inputFilePath)
	if err != nil {
		log.Println(err)
		return
	}
	bytes, err := ioutil.ReadAll(file)
	m := map[string]interface{}{}
	err = json.Unmarshal(bytes, &m)
	if err != nil {
		log.Println(err)
		return
	}
	entries := m["entries"].([]interface{})
	cnt := 0
	for _, _e := range entries {
		cnt += 1
		exec(filepath.Dir(*inputFilePath), *outputFilePath, _e)
	}
	log.Println(cnt)

	_f, err := os.Create(filepath.Join(*outputFilePath, "meta.json"))
	if err != nil {
		panic(err)
	}
	defer _f.Close()
	//uuid := strings.Replace(filepath.Base(*outputFilePath), filepath.Ext(*outputFilePath), "", -1)
	uuid := "DayOne"
	_f.Write([]byte(fmt.Sprintf(`{
  "name" : "%s",
  "uuid" : "%s"
}
`, uuid, uuid)))
}
