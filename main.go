package main

import (
	"fmt"
	"image/png"
	"path/filepath"
	"log"
	"os"
	"os/user"
	"net/http"
	"io/ioutil"
	"io"
	"encoding/json"
	"image/color"
	"image"
	_ "image/png"
	_ "image/jpeg"
	"math/rand"
	"math"
	"sort"
	"errors"
	"time"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os/exec"
)

type BingAPIResponse struct {
	Images   []struct {
		Startdate     string
		Fullstartdate string
		Enddate       string
		URL           string
		Urlbase       string
		Copyright     string
		Copyrightlink string
		Wp            bool
		Hsh           string
		Drk           int
		Top           int
		Bot           int
		Hs            []struct {
			Desc  string
			Link  string
			Query string
			Locx  int
			Locy  int
		}
		Msg           []interface{}
	}
	Tooltips struct {
				 Loading  string
				 Previous string
				 Next     string
				 Walle    string
				 Walls    string
			 }
}

const (
	threshold int = 50
	minBrightness int = 0
	maxBrightness int = 200
)

func changeDesktopBackground(path string) error {
	usr, err := user.Current()
	if err != nil {
		return err
	}

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("sqlite3", fmt.Sprintf("%s/Library/Application Support/Dock/desktoppicture.db", usr.HomeDir))
	if err != nil {
		return err
	}
	defer db.Close()

	file := fmt.Sprint(dir,"/",path)

	sqlSmt := fmt.Sprintf("update data set value = '%s';", file)

	_, err = db.Exec(sqlSmt)
	if err != nil {
		return err
	}

	cmd := exec.Command("killall", "Dock")
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func loadImage(filepath string) image.Image {
	infile, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer infile.Close()

	src, _, err := image.Decode(infile)
	if err != nil {
		log.Fatal(err)
	}
	return src
}

func abs(n int) int {
	if n >= 0 {
		return n
	}
	return -n
}

func randMinMax(min int, max int) int {
	if min == max {
		return min
	}
	return rand.Intn(max - min) + min
}

func randBool() bool {
	return rand.Intn(2) == 0
}

func colorDifference(col1 color.Color, col2 color.Color, threshold int) bool {
	c1 := col1.(color.NRGBA)
	c2 := col2.(color.NRGBA)

	rDiff := abs(int(c1.R) - int(c2.R))
	gDiff := abs(int(c1.G) - int(c2.G))
	bDiff := abs(int(c1.B) - int(c2.B))

	total := rDiff + gDiff + bDiff
	return total >= threshold
}

func getDistinctColors(colors []color.Color, threshold int, minBrightness, maxBrightness int) []color.Color {
	distinctColors := make([]color.Color, 0)
	for _, c := range colors {
		same := false
		if !colorDifference(c, color.NRGBAModel.Convert(color.Black), minBrightness * 3) {
			continue
		}
		if !colorDifference(c, color.NRGBAModel.Convert(color.White), (255 - maxBrightness) * 3) {
			continue
		}
		for _, k := range distinctColors {
			if !colorDifference(c, k, threshold) {
				same = true
				break
			}
		}
		if !same {
			distinctColors = append(distinctColors, c)
		}
	}
	return distinctColors
}

func colorsFromImage(filename string) ([]color.Color, error) {
	fuzzyness := 5
	img := loadImage(filename)
	w, h := img.Bounds().Max.X, img.Bounds().Max.Y
	colors := make([]color.Color, 0, w * h)
	for x := 0; x < w; x += fuzzyness {
		for y := 0; y < h; y += fuzzyness {
			col := color.NRGBAModel.Convert(img.At(x, y))
			colors = append(colors, col)
		}
	}
	distinctColors := getDistinctColors(colors, threshold, minBrightness, maxBrightness)

	count := 0
	for len(distinctColors) < 16 {
		count++
		distinctColors = append(distinctColors, getDistinctColors(colors, threshold - count, minBrightness, maxBrightness)...)
		if count == threshold {
			return nil, errors.New("Could not get colors from image with settings specified. Aborting.\n")
		}
	}

	if len(distinctColors) > 16 {
		distinctColors = distinctColors[:16]
	}
	return distinctColors, nil
}

type Circle struct {
	col  color.Color
	x, y int
	size int
}

type circleBySize []Circle

func (a circleBySize) Len() int { return len(a) }
func (a circleBySize) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a circleBySize) Less(i, j int) bool { return a[i].size < a[j].size }

func Circles(colors []color.Color, w int, h int, size int, sizevar int, overlap bool, large2small bool, filled bool, bordersize int) image.Image {
	img := image.NewNRGBA(image.Rect(0, 0, w, h))

	circles := make([]Circle, 0)
	bg := colors[0]

	for _, c := range colors {
		if c == bg {
			continue
		}
		circle := Circle{c, rand.Intn(w), rand.Intn(h), randMinMax(size - sizevar, size + sizevar)}
		circles = append(circles, circle)
	}

	if large2small {
		sort.Sort(sort.Reverse(circleBySize(circles)))
	}

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			img.Set(x, y, bg)
			for _, c := range circles {
				a := float64((x - c.x) * (x - c.x))
				b := float64((y - c.y) * (y - c.y))

				if filled {
					if int(math.Sqrt(a + b)) < c.size {
						img.Set(x, y, c.col)
					}
				} else {
					if int(math.Sqrt(a + b)) < c.size && int(math.Sqrt(a + b)) > (c.size - bordersize) {
						img.Set(x, y, c.col)
					}
				}
			}
		}
	}
	return img
}

type Ray struct {
	col   color.Color
	x, y  int // Middle point
	angle int // 0-180
	size  int
}

type rayBySize []Ray

func (a rayBySize) Len() int { return len(a) }
func (a rayBySize) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a rayBySize) Less(i, j int) bool { return a[i].size < a[j].size }

func Rays(colors []color.Color, w int, h int, size int, sizevar int, evendist bool, centered bool, large2small bool) image.Image {
	img := image.NewNRGBA(image.Rect(0, 0, w, h))

	rays := make([]Ray, 0)

	spacing := 360 / len(colors)
	current_angle := 0

	xpos := w / 2
	ypos := h / 2

	bg := colors[0]

	for _, c := range colors {
		if c == bg {
			continue
		}
		var ray Ray
		if !centered {
			xpos = rand.Intn(w)
			ypos = rand.Intn(h)
		}
		if !evendist {
			current_angle = rand.Intn(360)
		}
		ray = Ray{c, xpos, ypos, current_angle, randMinMax(size - sizevar, size + sizevar)}

		if evendist {
			current_angle += spacing + ray.size
		}
		rays = append(rays, ray)
	}

	if large2small {
		sort.Sort(sort.Reverse(rayBySize(rays)))
	}

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			img.Set(x, y, bg)
			for _, r := range rays {
				deltaX := float64(x - r.x)
				deltaY := float64(y - r.y)
				angle := math.Atan(deltaY / deltaX) * 180 / math.Pi
				if angle < 0 {
					angle += 360
				}
				if int(math.Abs(float64(int(angle) - r.angle))) < r.size {
					img.Set(x, y, r.col)
				}
			}
		}
	}
	return img
}

type Line struct {
	col      color.Color
	position int
	size     int
}

func Lines(colors []color.Color, w int, h int, size int, sizevar int, horizontal bool, equalspacing bool, spacingsize int, offset int) image.Image {
	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	var maxsize int
	if horizontal {
		maxsize = h
	} else {
		maxsize = w
	}

	currentposition := offset
	spacing := spacingsize

	lines := make([]Line, 0)
	bg := colors[0]

	for _, c := range colors {
		if c == bg {
			continue
		}
		line := Line{c, currentposition, randMinMax(size - sizevar, size + sizevar)}
		lines = append(lines, line)
		if !equalspacing {
			spacing = rand.Intn(maxsize / 16)
		}
		currentposition += line.size + spacing
	}

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			img.Set(x, y, bg)
			for _, l := range lines {
				var pixelpos int
				if horizontal {
					pixelpos = y
				} else {
					pixelpos = x
				}

				if pixelpos > l.position && pixelpos < l.position + l.size {
					img.Set(x, y, l.col)
				}
			}
		}
	}

	return img
}

func randomImage(colors []color.Color, w int, h int) image.Image {
	rand.Seed(time.Now().UnixNano())
	switch rand.Intn(3) {
	case 0:
		return Circles(colors, w, h, rand.Intn(w / 2), rand.Intn(w / 2), randBool(), randBool(), randBool(), rand.Intn(20))
	case 1:
		return Rays(colors, w, h, rand.Intn(h / 32) + 1, rand.Intn(h / 32), randBool(), true, randBool())
	case 2:
		return Lines(colors, w, h, rand.Intn(h / 32) + 1, rand.Intn(h / 32), randBool(), randBool(), rand.Intn(h / 32), rand.Intn(h / 2) + 1)
	}
	return nil
}

func getInputImage(body []byte) (string, error) {
	var b = new(BingAPIResponse)
	err := json.Unmarshal(body, &b)
	if (err != nil) {
		fmt.Println("NOPE:", err)
	}
	return "https://bing.com" + b.Images[0].URL, err
}

func main() {
	url := "http://www.bing.com/HPImageArchive.aspx?format=js&idx=0&n=1"
	res, err := http.Get(url)

	defer res.Body.Close()
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	imageUrl, err := getInputImage([]byte(body))

	response, err := http.Get(imageUrl)
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	usr, err := user.Current()
	err = os.Mkdir(fmt.Sprint(usr.HomeDir, "/.tapet"), 0666)
	if err != nil {
		println("dir NOPE")
		log.Fatal(err)
	}


	file, err := os.Create(fmt.Sprint(usr.HomeDir, "/.tapet/", "temp.jpg"))
	if err != nil {
		println("file NOPE")
		log.Fatal(err)
	}

	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	colors, err := colorsFromImage(fmt.Sprint(usr.HomeDir, "/.tapet/", "temp.jpg"))

	if len(colors) > 16 {
		colors = colors[:16]
	} else if len(colors) < 16 {
		log.Fatal("Less than 16 colors. Aborting.")
	}

	file, err = os.OpenFile(fmt.Sprint(usr.HomeDir, "/.tapet/", "background.jpg"), os.O_CREATE | os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	img := randomImage(colors, 1920, 1080)

	png.Encode(file, img)
	changeDesktopBackground(fmt.Sprint(usr.HomeDir, "/.tapet/", "temp.jpg"))
}
