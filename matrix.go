/*
@author: sk
@date: 2024/6/1
*/
package main

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

type Matrix [4][4]float64

func (m Matrix) Mul(matrix Matrix) Matrix {
	res := [4][4]float64{}
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			for k := 0; k < 4; k++ {
				res[i][j] += m[i][k] * matrix[k][j]
			}
		}
	}
	return res
}

func NewTrans(x, y, z float64) Matrix {
	return [4][4]float64{
		{1, 0, 0, x},
		{0, 1, 0, y},
		{0, 0, 1, z},
		{0, 0, 0, 1},
	}
}

func NewIdent() Matrix {
	return [4][4]float64{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
}

func NewScale(x, y, z float64) Matrix {
	return [4][4]float64{
		{x, 0, 0, 0},
		{0, y, 0, 0},
		{0, 0, z, 0},
		{0, 0, 0, 1},
	}
}

func NewRotateX(angle float64) Matrix {
	cos := math.Cos(angle)
	sin := math.Sin(angle)
	return [4][4]float64{
		{1, 0, 0, 0},
		{0, cos, -sin, 0},
		{0, sin, cos, 0},
		{0, 0, 0, 1},
	}
}

func NewRotateY(angle float64) Matrix {
	cos := math.Cos(angle)
	sin := math.Sin(angle)
	return [4][4]float64{
		{cos, 0, -sin, 0},
		{0, 1, 0, 0},
		{sin, 0, cos, 0},
		{0, 0, 0, 1},
	}
}

func NewRotateZ(angle float64) Matrix {
	cos := math.Cos(angle)
	sin := math.Sin(angle)
	return [4][4]float64{
		{cos, -sin, 0, 0},
		{sin, cos, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
}

func LookAt(eye, dir, up Vec) Matrix {
	z := dir.Normal().Scale(-1) // eye 朝向的是-z
	y := up.Normal()            // 相机上方是 y 的方向
	x := z.Cross(y)
	rotate := Matrix{ // 以 x y z 为基向量的逆矩阵
		{x.X, x.Y, x.Z, 0}, // 这个矩阵比较特殊，其转置矩阵就是其逆矩阵
		{y.X, y.Y, y.Z, 0},
		{z.X, z.Y, z.Z, 0},
		{0, 0, 0, 1},
	}
	trans := NewTrans(-eye.X, -eye.Y, -eye.Z) // 移动到 eye 为原点的的矩阵
	return rotate.Mul(trans)
}

func ProjectOrtho(l, r, t, b, n, f float64) Matrix {
	return [4][4]float64{ // 正交投影 根据给定的 6 个平面位置，把压缩平移到 x,y,z -1～1 的空间 只需要压缩与平移操作
		{2 / (r - l), 0, 0, -(r + l) / (r - l)},
		{0, 2 / (b - t), 0, -(b + t) / (b - t)},
		{0, 0, 2 / (f - n), -(f + n) / (f - n)},
		{0, 0, 0, 1},
	}
}

func ProjectPerspUE(hFov, w, h, near, far float64) Matrix {
	tan := math.Tan(hFov)
	return [4][4]float64{
		{1 / tan, 0, 0, 0},
		{0, w / h / tan, 0, 0},
		{0, 0, far / (far - near), 1},
		{0, 0, -near * far / (far - near), 0},
	}
}

func ProjectPersp(fov, aspect, near, far float64) Matrix { // 透视矩阵
	// 台柱变立方体
	temp := [4][4]float64{
		{near, 0, 0, 0},
		{0, near, 0, 0},
		{0, 0, near + far, -near * far},
		{0, 0, -1, 1},
	}
	// 直接复用正交矩阵 立方体压缩
	hHeight := math.Tan(fov/2) * near
	hWidth := hHeight * aspect
	ortho := ProjectOrtho(-hWidth, hWidth, -hHeight, hHeight, near, far)
	return ortho.Mul(temp)
}

func Viewport(w, h float64) Matrix {
	return [4][4]float64{ // 平移缩放以适配屏幕
		{w / 2, 0, 0, w / 2},
		{0, h / 2, 0, h / 2},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
}

func NewInverse(matrix Matrix) Matrix {
	mat := ToMat4(matrix)
	return ToMatrix(mat.Inv()) // 暂时借助其实现求逆
	//det := matrix[0][0]*(matrix[1][1]*matrix[2][2]-matrix[1][2]*matrix[2][1]) -
	//	matrix[0][1]*(matrix[1][0]*matrix[2][2]-matrix[1][2]*matrix[2][0]) +
	//	matrix[0][2]*(matrix[1][0]*matrix[2][1]-matrix[1][1]*matrix[2][0])
	//inv := 1.0 / det
	//
	//return [4][4]float64{
	//	{inv * (matrix[1][1]*matrix[2][2] - matrix[1][2]*matrix[2][1]), -inv * (matrix[0][1]*matrix[2][2] - matrix[0][2]*matrix[2][1]), inv * (matrix[0][1]*matrix[1][2] - matrix[0][2]*matrix[1][1]), 0},
	//	{-inv * (matrix[1][0]*matrix[2][2] - matrix[1][2]*matrix[2][0]), inv * (matrix[0][0]*matrix[2][2] - matrix[0][2]*matrix[2][0]), -inv * (matrix[0][0]*matrix[1][2] - matrix[0][2]*matrix[1][0]), 0},
	//	{inv * (matrix[1][0]*matrix[2][1] - matrix[1][1]*matrix[2][0]), -inv * (matrix[0][0]*matrix[2][1] - matrix[0][1]*matrix[2][0]), inv * (matrix[0][0]*matrix[1][1] - matrix[0][1]*matrix[1][0]), 0},
	//	{-(matrix[3][0]*matrix[0][0] + matrix[3][1]*matrix[1][0] + matrix[3][2]*matrix[2][0]), -(matrix[3][0]*matrix[0][1] + matrix[3][1]*matrix[1][1] + matrix[3][2]*matrix[2][1]), -(matrix[3][0]*matrix[0][2] + matrix[3][1]*matrix[1][2] + matrix[3][2]*matrix[2][2]), 1},
	//}
}

func ToMatrix(mat mgl32.Mat4) Matrix {
	res := Matrix{}
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			res[i][j] = float64(mat[i+j*4])
		}
	}
	return res
}

func ToMat4(matrix Matrix) mgl32.Mat4 {
	res := mgl32.Mat4{}
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			res[i+j*4] = float32(matrix[i][j])
		}
	}
	return res
}
