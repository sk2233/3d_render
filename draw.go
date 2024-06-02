/*
@author: sk
@date: 2024/5/30
*/
package main

import (
	"image"
	"image/color"
	"math"
)

func DrawLine(screen *image.RGBA, v1, v2 Vec, clr color.Color) {
	offX, offY := float64(0), float64(0)
	count := 0
	if math.Abs(v2.X-v1.X) > math.Abs(v2.Y-v1.Y) {
		count = int(math.Abs(v2.X - v1.X))
		offX = Sign(v2.X - v1.X)
		offY = (v2.Y - v1.Y) / (v2.X - v1.X)
	} else {
		count = int(math.Abs(v1.Y - v2.Y))
		offX = (v2.X - v1.X) / (v2.Y - v1.Y)
		offY = Sign(v2.Y - v1.Y)
	}
	x, y := v1.X, v1.Y
	for i := 0; i < count; i++ {
		screen.Set(int(x), int(y), clr)
		x += offX
		y += offY
	}
}

func DrawTrianglePro(img *image.RGBA, w, h int, shader IShader, as ...FArgs) {
	// 计算绘制的区域
	minX, minY, maxX, maxY := float64(w), float64(h), float64(0), float64(0)
	for _, arg := range as {
		minX = math.Min(minX, arg.NewV.X)
		maxX = math.Max(maxX, arg.NewV.X)
		minY = math.Min(minY, arg.NewV.Y)
		maxY = math.Max(maxY, arg.NewV.Y)
	}
	minX = math.Max(minX, 0)
	maxX = math.Min(maxX, float64(w))
	minY = math.Max(minY, 0)
	maxY = math.Min(maxY, float64(h))
	// 循环绘制
	minXI, minYI, maxXI, maxYI := int(minX), int(minY), int(maxX), int(maxY)
	for x := minXI; x <= maxXI; x++ {
		for y := minYI; y <= maxYI; y++ {
			weight := GetWight(float64(x), float64(y), as) // 获取占3个顶点的权重
			if weight.X >= 0 && weight.X <= 1 && weight.Y >= 0 &&
				weight.Y <= 1 && weight.Z >= 0 && weight.Z <= 1 { // 在三角形内才绘制
				z := weight.X*as[0].NewV.Z + weight.Y*as[1].NewV.Z + weight.Z*as[2].NewV.Z
				w0 := weight.X*as[0].NewV.W + weight.Y*as[1].NewV.W + weight.Z*as[2].NewV.W
				if res, ok := shader.Fragment(FArgs{
					NewV:  Vec{X: float64(x), Y: float64(y), Z: z, W: w0},
					OldV:  NewVecByWeight(weight, as[0].OldV, as[1].OldV, as[2].OldV), // 可能会有误差 不再使用
					NewVn: NewVecByWeight(weight, as[0].NewVn, as[1].NewVn, as[2].NewVn),
					OldVt: NewVecByWeight(weight, as[0].OldVt, as[1].OldVt, as[2].OldVt),
				}); ok {
					img.Set(x, y, res.Clr)
				}
			}
		}
	}
}

func NewVecByWeight(weight Vec, v1, v2, v3 Vec) Vec {
	return NewVec3(v1.X*weight.X+v2.X*weight.Y+v3.X*weight.Z, v1.Y*weight.X+v2.Y*weight.Y+v3.Y*weight.Z,
		v1.Z*weight.X+v2.Z*weight.Y+v3.Z*weight.Z)
}

func GetWight(x float64, y float64, as []FArgs) Vec {
	//b := ((x-as[2].NewV.X)*(as[0].NewV.Y-as[2].NewV.Y) + (y-as[2].NewV.Y)*(as[0].NewV.X-as[2].NewV.X)) /
	//	((as[1].NewV.X-as[2].NewV.X)*(as[0].NewV.Y-as[2].NewV.Y) + (as[1].NewV.Y-as[2].NewV.Y)*(as[0].NewV.X-as[2].NewV.X))
	//a := (x - as[2].NewV.X - b*(as[1].NewV.X-as[2].NewV.X)) / (as[0].NewV.X - as[2].NewV.X)
	a := ((x-as[1].NewV.X)*(as[2].NewV.Y-as[1].NewV.Y) - (y-as[1].NewV.Y)*(as[2].NewV.X-as[1].NewV.X)) /
		((as[0].NewV.X-as[1].NewV.X)*(as[2].NewV.Y-as[1].NewV.Y) - (as[0].NewV.Y-as[1].NewV.Y)*(as[2].NewV.X-as[1].NewV.X))
	b := ((x-as[2].NewV.X)*(as[0].NewV.Y-as[2].NewV.Y) - (y-as[2].NewV.Y)*(as[0].NewV.X-as[2].NewV.X)) /
		((as[1].NewV.X-as[2].NewV.X)*(as[0].NewV.Y-as[2].NewV.Y) - (as[1].NewV.Y-as[2].NewV.Y)*(as[0].NewV.X-as[2].NewV.X))
	c := 1 - a - b
	return NewVec3(a, b, c)
}

func DrawTriangle(screen *image.RGBA, r1, r2, r3 FArgs, shader IShader) {
	// 确保按y v1 v2 v3 排序
	if r1.NewV.Y > r2.NewV.Y {
		r1, r2 = r2, r1
	}
	if r1.NewV.Y > r3.NewV.Y {
		r1, r3 = r3, r1
	}
	if r2.NewV.Y > r3.NewV.Y {
		r2, r3 = r3, r2
	}
	//// 先用纹理的中间点作为纹理差值中间点
	//mid := FArgs{
	//	OldVt: r1.OldVt.Add(r2.OldVt).Add(r3.OldVt).Scale(1.0 / 3),
	//}
	// 先绘制上部分 v1.Y ~ v2.Y
	vn := r1.NewVn.Add(r2.NewVn).Add(r3.NewVn).Scale(1.0 / 3) // 法向量一个三角形面是一样的，先计算出来
	if r1.NewV.Y < r2.NewV.Y {                                // 存在上部分
		y1, y2 := int(r1.NewV.Y), int(r2.NewV.Y)
		x1f, z1f, x2f, z2f := r1.NewV.X, r1.NewV.Z, r1.NewV.X, r1.NewV.Z
		// 纹理坐标计算
		tx1, ty1, tx2, ty2 := r1.OldVt.X, r1.OldVt.Z, r1.OldVt.X, r1.OldVt.Z
		xOff1, zOff1, xOff2, zOff2 := (r2.NewV.X-r1.NewV.X)/(r2.NewV.Y-r1.NewV.Y), (r2.NewV.Z-r1.NewV.Z)/(r2.NewV.Y-r1.NewV.Y),
			(r3.NewV.X-r1.NewV.X)/(r3.NewV.Y-r1.NewV.Y), (r3.NewV.Z-r1.NewV.Z)/(r3.NewV.Y-r1.NewV.Y)
		// 与绘制点类似也是使用偏移
		txOff1, tyOff1, txOff2, tyOff2 := (r2.OldVt.X-r1.OldVt.X)/(r2.NewV.Y-r1.NewV.Y), (r2.OldVt.Y-r1.OldVt.Y)/(r2.NewV.Y-r1.NewV.Y),
			(r3.OldVt.X-r1.OldVt.X)/(r3.NewV.Y-r1.NewV.Y), (r3.OldVt.Y-r1.OldVt.Y)/(r3.NewV.Y-r1.NewV.Y)
		if xOff1 > xOff2 {
			xOff1, xOff2 = xOff2, xOff1
			zOff1, zOff2 = zOff2, zOff1
			txOff1, txOff2 = txOff2, txOff1
			tyOff1, tyOff2 = tyOff2, tyOff1
		}
		for y := y1; y <= y2; y++ { // 上包含中间的那条线，下不包
			x1, x2 := int(x1f), int(x2f)
			z := z1f
			zOff := (z2f - z1f) / (x2f - x1f)
			tx, ty := tx1, ty1
			txOff, tyOff := (tx2-tx1)/(x2f-x1f), (ty2-ty1)/(x2f-x1f)
			for x := x1; x <= x2; x++ {
				args := FArgs{
					NewV:  NewVec3(float64(x), float64(y), z),
					OldVt: NewVec2(tx, ty),
					NewVn: vn,
				}
				if res, ok := shader.Fragment(args); ok {
					screen.Set(x, y, res.Clr)
				}
				z += zOff
				tx += txOff
				ty += tyOff
			}
			// 坐标偏移
			x1f += xOff1
			x2f += xOff2
			z1f += zOff1
			z2f += zOff2
			// 纹理偏移
			tx1 += txOff1
			ty1 += tyOff1
			tx2 += txOff2
			ty2 += tyOff2
		}
	}
	// 再绘制下部分 v2.Y ~ v3.Y
	if r2.NewV.Y < r3.NewV.Y { // 存在下部分 与上面类似不过是从底部绘制的
		y1, y2 := int(r3.NewV.Y), int(r2.NewV.Y)
		x1f, z1f, x2f, z2f := r3.NewV.X, r3.NewV.Z, r3.NewV.X, r3.NewV.Z
		// 纹理坐标计算
		tx1, ty1, tx2, ty2 := r3.OldVt.X, r3.OldVt.Z, r3.OldVt.X, r3.OldVt.Z
		xOff1, zOff1, xOff2, zOff2 := (r1.NewV.X-r3.NewV.X)/(r3.NewV.Y-r1.NewV.Y), (r1.NewV.Z-r3.NewV.Z)/(r3.NewV.Y-r1.NewV.Y),
			(r2.NewV.X-r3.NewV.X)/(r3.NewV.Y-r2.NewV.Y), (r2.NewV.Z-r3.NewV.Z)/(r3.NewV.Y-r2.NewV.Y)
		// 与绘制点类似也是使用偏移
		txOff1, tyOff1, txOff2, tyOff2 := (r1.OldVt.X-r2.OldVt.X)/(r3.NewV.Y-r1.NewV.Y), (r1.OldVt.Y-r3.OldVt.Y)/(r3.NewV.Y-r1.NewV.Y),
			(r2.OldVt.X-r3.OldVt.X)/(r3.NewV.Y-r2.NewV.Y), (r2.OldVt.Y-r3.OldVt.Y)/(r3.NewV.Y-r2.NewV.Y)
		if xOff1 > xOff2 {
			xOff1, xOff2 = xOff2, xOff1
			zOff1, zOff2 = zOff2, zOff1
			txOff1, txOff2 = txOff2, txOff1
			tyOff1, tyOff2 = tyOff2, tyOff1
		}
		for y := y1; y > y2; y-- {
			x1, x2 := int(x1f), int(x2f)
			z := z1f
			zOff := (z2f - z1f) / (x2f - x1f)
			tx, ty := tx1, ty1
			txOff, tyOff := (tx2-tx1)/(x2f-x1f), (ty2-ty1)/(x2f-x1f)
			for x := x1; x <= x2; x++ {
				args := FArgs{
					NewV:  NewVec3(float64(x), float64(y), z),
					OldVt: NewVec2(tx, ty),
					NewVn: vn,
				}
				if res, ok := shader.Fragment(args); ok {
					screen.Set(x, y, res.Clr)
				}
				z += zOff
				tx += txOff
				ty += tyOff
			}
			x1f += xOff1
			x2f += xOff2
			z1f += zOff1
			z2f += zOff2
			// 纹理偏移
			tx1 += txOff1
			ty1 += tyOff1
			tx2 += txOff2
			ty2 += tyOff2
		}
	}
}

func Fill(screen *image.RGBA, clr color.Color) {
	c1, c2, c3, c4 := clr.RGBA()
	r, g, b, a := uint8(c1>>8), uint8(c2>>8), uint8(c3>>8), uint8(c4>>8)
	for i := 0; i < len(screen.Pix); i += 4 {
		screen.Pix[i], screen.Pix[i+1], screen.Pix[i+2], screen.Pix[i+3] = r, g, b, a
	}
}
