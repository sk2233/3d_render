/*
@author: sk
@date: 2024/5/30
*/
package main

import "math"

type Vec struct {
	X, Y, Z, W float64
}

func (v Vec) Scale(scale float64) Vec {
	return Vec{X: v.X * scale, Y: v.Y * scale, Z: v.Z * scale, W: v.W * scale}
}

func (v Vec) Add(vec Vec) Vec {
	return Vec{
		X: v.X + vec.X,
		Y: v.Y + vec.Y,
		Z: v.Z + vec.Z,
		W: v.W + vec.W,
	}
}

func (v Vec) Mul(matrix Matrix) Vec {
	return Vec{
		X: v.X*matrix[0][0] + v.Y*matrix[0][1] + v.Z*matrix[0][2] + v.W*matrix[0][3],
		Y: v.X*matrix[1][0] + v.Y*matrix[1][1] + v.Z*matrix[1][2] + v.W*matrix[1][3],
		Z: v.X*matrix[2][0] + v.Y*matrix[2][1] + v.Z*matrix[2][2] + v.W*matrix[2][3],
		W: v.X*matrix[3][0] + v.Y*matrix[3][1] + v.Z*matrix[3][2] + v.W*matrix[3][3],
	}
}

func (v Vec) Normal() Vec {
	l := v.Len()
	return v.Scale(1 / l)
}

func (v Vec) Len() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

func (v Vec) Cross(vec Vec) Vec {
	return Vec{
		X: v.Y*vec.Z - v.Z*vec.Y,
		Y: v.X*vec.Z - v.Z*vec.X,
		Z: v.X*vec.Y - v.Y*vec.X,
	}
}

func (v Vec) Dot(vec Vec) float64 {
	return v.X*vec.X + v.Y*vec.Y + v.Z*vec.Z
}

func (v Vec) Reflect(vn Vec) Vec {
	// 计算 v 以 vn 为法线的反射光线 v vn 都是单位向量
	return v.Sub(vn.Scale(2 * v.Dot(vn)))
}

func (v Vec) Sub(vec Vec) Vec {
	return Vec{
		X: v.X - vec.X,
		Y: v.Y - vec.Y,
		Z: v.Z - vec.Z,
		W: v.W - vec.W,
	}
}

func NewVec3(x float64, y float64, z float64) Vec {
	return Vec{X: x, Y: y, Z: z, W: 1} // w 默认需要为 1 否则平移不会有作用
}

func NewVec2(x float64, y float64) Vec {
	return Vec{X: x, Y: y}
}
