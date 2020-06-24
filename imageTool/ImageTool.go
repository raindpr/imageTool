
package main

import (
	"bytes"
	"code.google.com/p/graphics-go/graphics"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
)


//帮助提示信息
var usage = func() {
	fmt.Println("输入错误，请按照下面的格式输入：")
	fmt.Println("使用: imagetool [OPTION]  source_image [output]")
	fmt.Println("  Options is flow:")
	fmt.Println("    -r         图片颜色翻转")
	fmt.Println("    -g         图片灰度")
	fmt.Println("    -c         缩放文本，该参数时，可以传入图片缩放的宽度 如：imagetool -c 1.jpg c1.jpg 100")
	fmt.Println("    -t         转成文本")
	fmt.Println("    -cut       剪切单个图片，该参数时，可以传入剪切的图片 如：imagetool -cut 1.jpg")
	fmt.Println("    -cut-dir   剪切文件夹所以图片，该参数时，可以传入剪切的图片和生成文件夹 如：imagetool -cut-dir filename")
}

func main() {
	args := os.Args //获取用户输入的所有参数
	if args == nil || len(args) < 3 || !(args[1] == "-r" || args[1] == "-g" || args[1] == "-t" || args[1] == "-c" || args[1] == "-cut" || args[1] == "-cut-dir") {
		usage()
		return
	}

	option := args[1]
	source := args[2]
	target := ""
	if len(args)>3 {
		target = args[3]
	}

	//option := "-cut-dir"
	//source := "Texture2D"
	//读取文件
	ff, _ := ioutil.ReadFile(source)
	bbb := bytes.NewBuffer(ff)
	m, _, _ := image.Decode(bbb)

	if option == "-r" {
		newRgba := fzImage(m)
		f, _ := os.Create(target)
		defer f.Close()
		encode(source, f, newRgba)
	} else if option == "-g" {
		newGray := hdImage(m)
		f, _ := os.Create(target)
		defer f.Close()
		encode(source, f, newGray)
	} else if option == "-c" {
		rectWidth := 200
		if len(args) > 4 {
			rectWidth, _ = strconv.Atoi(args[4])
		}
		newRgba := rectImage(m, rectWidth)
		f, _ := os.Create(target)
		defer f.Close()
		encode(source, f, newRgba)
	} else if option == "-cut" {
		clipping(m,source)
	}else if option == "-cut-dir" {
		GetAllFile(source)
	}else {
		ascllimage(m, target)
	}
	fmt.Println("转换完成...")
}

//图片编码
func encode(inputName string, file *os.File, rgba *image.RGBA) {
	if strings.HasSuffix(inputName, "jpg") || strings.HasSuffix(inputName, "jpeg") {
		jpeg.Encode(file, rgba, nil)
	} else if strings.HasSuffix(inputName, "png") {
		png.Encode(file, rgba)
	} else if strings.HasSuffix(inputName, "gif") {
		gif.Encode(file, rgba, nil)
	} else {
		fmt.Errorf("不支持的图片格式")
	}
}


//图片色彩反转
func fzImage(m image.Image) *image.RGBA {
	bounds := m.Bounds()
	dx := bounds.Dx()
	dy := bounds.Dy()
	newRgba := image.NewRGBA(bounds)
	for i := 0; i < dx; i++ {
		for j := 0; j < dy; j++ {
			colorRgb := m.At(i, j)
			r, g, b, a := colorRgb.RGBA()
			r_uint8 := uint8(r >> 8)
			g_uint8 := uint8(g >> 8)
			b_uint8 := uint8(b >> 8)
			a_uint8 := uint8(a >> 8)
			r_uint8 = 255 - r_uint8
			g_uint8 = 255 - g_uint8
			b_uint8 = 255 - b_uint8
			newRgba.SetRGBA(i, j, color.RGBA{r_uint8, g_uint8, b_uint8, a_uint8})
		}
	}
	return newRgba
}

//图片灰化处理
func hdImage(m image.Image) *image.RGBA {
	bounds := m.Bounds()
	dx := bounds.Dx()
	dy := bounds.Dy()
	newRgba := image.NewRGBA(bounds)
	for i := 0; i < dx; i++ {
		for j := 0; j < dy; j++ {
			colorRgb := m.At(i, j)
			_, g, _, a := colorRgb.RGBA()
			g_uint8 := uint8(g >> 8)
			a_uint8 := uint8(a >> 8)
			newRgba.SetRGBA(i, j, color.RGBA{g_uint8, g_uint8, g_uint8, a_uint8})
		}
	}
	return newRgba
}

//图片缩放, add at 2018-9-12
func rectImage(m image.Image, newdx int) *image.RGBA {
	bounds := m.Bounds()
	dx := bounds.Dx()
	dy := bounds.Dy()
	newRgba := image.NewRGBA(image.Rect(0, 0, newdx, newdx*dy/dx))
	graphics.Scale(newRgba, m)
	return newRgba
}
//图片转为字符画（简易版）
func ascllimage(m image.Image, target string) {
	if m.Bounds().Dx() > 300 {
		m = rectImage(m, 300)
	}
	bounds := m.Bounds()
	dx := bounds.Dx()
	dy := bounds.Dy()
	arr := []string{"M", "N", "H", "Q", "$", "O", "C", "?", "7", ">", "!", ":", "–", ";", "."}

	fileName := target
	dstFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer dstFile.Close()
	for i := 0; i < dy; i++ {
		for j := 0; j < dx; j++ {
			colorRgb := m.At(j, i)
			_, g, _, _ := colorRgb.RGBA()
			avg := uint8(g >> 8)
			num := avg / 18
			dstFile.WriteString(arr[num])
			if j == dx-1 {
				dstFile.WriteString("\n")
			}
		}
	}
}

var mapPic = make(map[int]*PicCell)
var picFlog = make(map[int]int)


type PicCell struct {
	point []PicPoint
	name  int
	minX  int
	minY  int
	maxX  int
	maxY  int
}

func NewPicCell(n int) *PicCell {
	return &PicCell{
		point: []PicPoint{},
		name:  n,
		minX:  65535,
		minY:  65535,
		maxX:  0,
		maxY:  0,
	}
}

func (c *PicCell)add(p PicPoint)  {
	c.point = append(c.point, p)
}

type PicPoint struct {
	pox int
	poy int
}

func isHavePix(r uint32,g uint32,b uint32, a uint32)bool{
	return !(r == 0 && g == 0 && b == 0 && a == 0)
}

func clipping(m image.Image,source string) {
	fmt.Println("扫描图片"+source)
	mapPic = make(map[int]*PicCell)
	picFlog = make(map[int]int)
	bounds := m.Bounds()
	dx := bounds.Dx()
	dy := bounds.Dy()
	for x := 0; x < dx; x++ {
		for y := 0; y < dy; y++ {
			colorRgb := m.At(x, y)
			r, g, b, a := colorRgb.RGBA()
			if isHavePix(r,g,b,a) {
				key := x * 10000 + y
				_, ok := picFlog [key]
				if !ok {
					mapPic[key] = NewPicCell(key)
					FindPoint(x, y, m, key)
				}
			}
		}
	}
	baseFile := path.Base(source)
	dirFile := path.Dir(source)
	filesplit := strings.Split(baseFile,".")
	outDir := path.Join(dirFile,filesplit[0])
	if !isExist(outDir){
		os.Mkdir(outDir, os.ModePerm)
	}
	for _,cell := range mapPic {
		newRgba := image.NewRGBA(image.Rect(0, 0, cell.maxX - cell.minX, cell.maxY - cell.minY))
		for i := cell.minX; i <= cell.maxX; i++ {
			for j := cell.minY; j <= cell.maxY; j++ {
				colorRgb := m.At(i, j)
				r, g, b, a := colorRgb.RGBA()
				r_uint8 := uint8(r >> 8)
				g_uint8 := uint8(g >> 8)
				b_uint8 := uint8(b >> 8)
				a_uint8 := uint8(a >> 8)
				newRgba.SetRGBA(i-cell.minX, j-cell.minY, color.RGBA{r_uint8, g_uint8, b_uint8, a_uint8})
			}
		}
		f, _ := os.Create(fmt.Sprintf("%s/%d.png", outDir, cell.name))
		encode(source, f, newRgba)
		f.Close()
		newRgba = nil
	}

}
func FindPoint(x int,y int,m image.Image,key int)  {
	if x < 0 || x > 4096 {
		return
	}
	if y < 0 || y > 4096 {
		return
	}
	colorRgb := m.At(x, y)
	r, g, b, a := colorRgb.RGBA()
	if isHavePix(r,g,b,a) {
		cell, picFind := mapPic[key]
		if !picFind {
			fmt.Println("NewPicCell")
			mapPic[key] = NewPicCell(key)
			cell = mapPic[key]
		}
		//cell.add(PicPoint{pox:x,poy:y})
		k := x * 10000  + y
		if x < cell.minX {
			cell.minX = x
		}
		if x > cell.maxX {
			cell.maxX = x
		}
		if y < cell.minY {
			cell.minY = y
		}
		if y > cell.maxY {
			cell.maxY = y
		}
		value, ok := picFlog[k]
		if !ok {
			picFlog[k] = key
			FindPoint(x-1,y+1,m,key)
			FindPoint(x,y+1,m,key)
			FindPoint(x+1,y+1,m,key)

			FindPoint(x-1,y,m,key)
			FindPoint(x+1,y,m,key)

			FindPoint(x-1,y-1,m,key)
			FindPoint(x,y-1,m,key)
			FindPoint(x+1,y-1,m,key)
		}else {
			if value == key{

			}else {

			}
		}
	}
}
func GetAllFile(pathname string) error {
	rd, err := ioutil.ReadDir(pathname)
	for _, fi := range rd {
		if fi.IsDir() {
			fmt.Printf("[%s]\n", pathname+"/"+fi.Name())
			GetAllFile(pathname + fi.Name() + "/")
		} else {
			ff, _ := ioutil.ReadFile(path.Join(pathname,fi.Name()))
			bbb := bytes.NewBuffer(ff)
			m, _, _ := image.Decode(bbb)
			if m != nil {
				clipping(m, path.Join(pathname,fi.Name()))
			}
		}
	}
	return err
}
//判断文件或文件夹是否存在
func isExist(path string)bool{
	_, err := os.Stat(path)
	if err != nil{
		if os.IsExist(err){
			return true
		}
		if os.IsNotExist(err){
			return false
		}
		fmt.Println(err)
		return false
	}
	return true
}