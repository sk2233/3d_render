/*
@author: sk
@date: 2024/5/31
*/
package main

import (
	"image"
	"image/color"
	"math"
)

type VArgs struct {
	V  Vec // 3d空间坐标
	Vt Vec // 纹理坐标
	Vn Vec // 法向量
}

type VRes struct {
	// 转换后的
	NewV  Vec
	NewVn Vec // 法线也需要进行转换
	// 转换前的数据
	OldV  Vec // 顶点
	OldVt Vec // 纹理坐标
}

type FArgs struct {
	NewV  Vec // 2d 画布位置
	OldV  Vec
	NewVn Vec
	OldVt Vec
}

type FRes struct {
	Clr color.Color
}

type IShader interface {
	Vertex(args VArgs) VRes
	Fragment(args FArgs) (FRes, bool)
}

//===================ImageShader===================

type ImageShader struct {
	img      image.Image
	w, h     float64
	deps     [][]float64 // w h
	matrix   Matrix
	lightDir Vec // 这里固定为平行光
	eye      Vec // 观察位置
}

func (i *ImageShader) Vertex(args VArgs) VRes {
	pos := args.V.Mul(i.matrix)
	return VRes{
		NewV:  pos,
		OldVt: args.Vt, // 纹理坐标透传
		NewVn: args.Vn,
	}
}

// TODO 阴影需要渲染两次，第一次从光源方向渲染得到灰度图，第二次正常渲染但是点使用第一次渲染的矩阵转换到第一次的位置查看深度是否有问题（要有容差）

func (i *ImageShader) Fragment(args FArgs) (FRes, bool) {
	x, y := int(args.NewV.X), int(args.NewV.Y)
	if x < 0 || x >= len(i.deps) || y < 0 || y >= len(i.deps[x]) { // 屏幕裁剪
		return FRes{}, false
	}
	if i.deps[x][y] > args.NewV.Z { // 深度测试
		return FRes{}, false
	}
	i.deps[x][y] = args.NewV.Z
	//c := uint8(args.NewV.Z) // 深度图
	//return FRes{
	//	Clr: NewRGB(c, c, c),
	//}, true
	clr := i.img.At(int(args.OldVt.X*i.w), int(args.OldVt.Y*i.h)) // 采样
	r, g, b, _ := clr.RGBA()
	ambient := Vec{X: float64(r) * 0.5, Y: float64(g) * 0.5, Z: float64(b) * 0.5} // 环境光
	temp := math.Max(i.lightDir.Scale(-1).Dot(args.NewVn)*0.1*math.MaxUint16, 0)
	diffuse := Vec{X: temp, Y: temp, Z: temp}    // 漫反射
	reflectDir := i.lightDir.Reflect(args.NewVn) // TODO 法线也需要经历 坐标变换，不能直接使用，且若物体进过了不等比缩放法线也需要对其纠正
	eyeDir := i.eye.Sub(args.NewV).Normal()      // TODO 当前点在 3d 空间的位置没有传递进来 暂时使用的是视口空间下的位置
	temp = math.Max(math.Pow(reflectDir.Dot(eyeDir), 4)*0.6*math.MaxUint16, 0)
	reflect := Vec{X: temp, Y: temp, Z: temp} // 镜面反射
	all := ambient.Add(diffuse).Add(reflect)
	return FRes{
		Clr: color.RGBA{
			R: uint8(all.X / 0x100),
			G: uint8(all.Y / 0x100),
			B: uint8(all.Z / 0x100),
			A: 0xFF,
		},
	}, true
}

func NewImageShader(img image.Image, w, h int, matrix Matrix, eye Vec) *ImageShader {
	deps := make([][]float64, w)
	for i := 0; i < w; i++ {
		deps[i] = make([]float64, h)
		for j := 0; j < h; j++ {
			deps[i][j] = math.MinInt64
		}
	}
	bound := img.Bounds()
	return &ImageShader{img: img, w: float64(bound.Dx()), h: float64(bound.Dy()),
		deps: deps, matrix: matrix, eye: eye, lightDir: NewVec3(0, 0, -1)} // 光源方向暂时固定
}

//===================SimpleShader====================

type SimpleShader struct {
	clr color.Color
}

func NewSimpleShader(clr color.Color) *SimpleShader {
	return &SimpleShader{clr: clr}
}

func (s *SimpleShader) Vertex(args VArgs) VRes {
	pos := args.V.Scale(240)
	return VRes{
		NewV: pos,
	}
}

func (s *SimpleShader) Fragment(args FArgs) FRes {
	return FRes{
		Clr: s.clr,
	}
}

//=======================DepShader 只绘制深度的shader=============================

type DepthShader struct {
	deps   [][]float64 // 把深度值写进来
	matrix Matrix
}

func NewDepthShader(deps [][]float64, matrix Matrix) *DepthShader {
	return &DepthShader{deps: deps, matrix: matrix}
}

func (d *DepthShader) Vertex(args VArgs) VRes {
	return VRes{
		NewV: args.V.Mul(d.matrix),
	}
}

//var (
//	Min = float64(math.MaxInt64)
//	Max = float64(math.MinInt64)
//)

func (d *DepthShader) Fragment(args FArgs) (FRes, bool) {
	x, y := int(args.NewV.X), int(args.NewV.Y)
	if x < 0 || x >= len(d.deps) || y < 0 || y >= len(d.deps[x]) {
		return FRes{}, false
	}
	if d.deps[x][y] > args.NewV.Z {
		return FRes{}, false
	}
	d.deps[x][y] = args.NewV.Z
	//Min = math.Min(Min, d.deps[x][y])
	//Max = math.Max(Max, d.deps[x][y])
	// 随便返回一个颜色即可，主要是记录深度值
	c := math.Min((args.NewV.Z+0.28)*0xFF/0.56, 0xFF)
	c = math.Max(c, 0)
	return FRes{
		Clr: NewRGB(uint8(c), uint8(c), uint8(c)),
	}, true
}

//==========================LightShader 光照 shader需要与DepthShader结合使用============================

type LightShader struct {
	matrix, matrixInv Matrix
	shadowDeps        [][]float64
	shadowM           Matrix
	lightDir          Vec
	eyePos            Vec
	deps              [][]float64
	img               image.Image
	w, h              float64
}

func (l *LightShader) Vertex(args VArgs) VRes {
	return VRes{
		NewV:  args.V.Mul(l.matrix),
		OldV:  args.V,
		NewVn: args.Vn.Mul(l.matrix), // 法线也需要做转换 且若模型进过非等比缩放还要矫正法线
		OldVt: args.Vt,
	}
}

func (l *LightShader) Fragment(args FArgs) (FRes, bool) {
	// 处理深度与越界问题
	x, y := int(args.NewV.X), int(args.NewV.Y)
	if x < 0 || x >= len(l.deps) || y < 0 || y >= len(l.deps[x]) {
		return FRes{}, false
	}
	if l.deps[x][y] > args.NewV.Z {
		return FRes{}, false
	}
	l.deps[x][y] = args.NewV.Z
	// 颜色采样
	clr := l.GetClr(args.OldVt)
	// 处理阴影
	shadowV := args.NewV.Mul(l.matrixInv).Mul(l.shadowM) // shadowM 包含了相机的逆矩阵与光源矩阵
	x, y = int(shadowV.X), int(shadowV.Y)
	if x >= 0 && x < len(l.shadowDeps) && y >= 0 && y < len(l.shadowDeps[x]) {
		z := l.shadowDeps[x][y]
		if math.Abs(shadowV.Z-z) > 0.001 { // 有 0.001 的容差
			return FRes{
				Clr: NewRGBByVec(clr.Scale(0.5)), // 阴影减淡
			}, true
		}
	}
	// 环境光
	ambient := clr.Scale(0.5) // 取总光照的 0.5为权重
	// 漫反射
	vn := args.NewVn.Normal()
	temp := math.Abs(l.lightDir.Dot(vn)) * 0.1 * math.MaxUint16 // 占 0.1 的权重
	diffuse := Vec{X: temp, Y: temp, Z: temp}
	// 镜面反射
	reflectDir := l.lightDir.Reflect(vn)
	eyeDir := l.eyePos.Sub(args.NewV).Normal()
	temp = math.Pow(math.Abs(reflectDir.Dot(eyeDir)), 4) * 0.4 * math.MaxUint16 // 占 0.4 的权重
	reflect := Vec{X: temp, Y: temp, Z: temp}
	// 合并颜色
	clr = ambient.Add(diffuse).Add(reflect)
	return FRes{
		Clr: NewRGBByVec(clr),
	}, true
}

func (l *LightShader) GetClr(vt Vec) Vec {
	temp := l.img.At(int(vt.X*l.w), int(vt.Y*l.h))
	r, g, b, _ := temp.RGBA() // 暂时不关心透明
	return NewVec3(float64(r), float64(g), float64(b))
}

func NewLightShader(matrix, shadowM Matrix, w, h int, img image.Image, shadowDeps [][]float64, lightDir Vec, eyePos Vec) *LightShader {
	deps := make([][]float64, w)
	for i := 0; i < w; i++ {
		deps[i] = make([]float64, h)
		for j := 0; j < h; j++ {
			deps[i][j] = math.MinInt64
		}
	}
	bound := img.Bounds()
	return &LightShader{matrix: matrix, matrixInv: NewInverse(matrix), shadowM: shadowM, shadowDeps: shadowDeps, deps: deps, img: img,
		w: float64(bound.Dx()), h: float64(bound.Dy()), lightDir: lightDir, eyePos: eyePos}
}
