/*
@author: sk
@date: 2024/5/30
*/
package main

import (
	"math"
)

// https://github.com/ssloy/tinyrenderer/wiki/Lesson-0:-getting-started
// https://www.bilibili.com/video/av6731067/?p=4&spm_id_from=333.788.0.0&vd_source=923fe31a18e5b835f6cc1eeb2a08340b
// https://space.bilibili.com/3493092766255793/video

func main() {
	//DrawLineTest()
	//DrawTriangleTest()
	//DrawOrthoTest()
	//DrawPerspTest()
	DrawTestPro()
}

func DrawTestPro() {
	// 默认模型文件数据在描述 模型空间
	// 使用 模型矩阵 把物体移动缩放到 世界空间
	// 使用 观察矩阵 把物体移动到观察者为原点的坐标系内(从相机看到的画面) 观察空间
	// 使用 投影矩阵 把所有物体压缩到 x,y,z 的 -1～1 齐次空间
	// 最后使用 viewport 把 ~1到1的内容铺满屏幕
	lightDir := NewVec3(0, 0, -1)
	eyePos := NewVec3(1, 1, 1)
	model := NewTrans(-800, -800, 0).Mul(NewScale(480, 480, 480))     // 模型矩阵
	lightView := LookAt(NewVec3(0, 0, 1), lightDir, NewVec3(0, 1, 0)) // 光源 观察矩阵
	eyeView := LookAt(eyePos, NewVec3(-1, -1, -1), NewVec3(0, 1, 0))  // 相机 观察矩阵
	project := ProjectOrtho(-1200, 1200, -1200, 1200, -1200, 1200)    // 投影矩阵
	viewport := Viewport(1200, 1200)                                  // 视口变换
	// 光源矩阵 & 相机矩阵
	depthM := viewport.Mul(project).Mul(lightView).Mul(model)
	eyeM := viewport.Mul(project).Mul(eyeView).Mul(model)
	// 初始化光源处的深度值为最小值
	deps := make([][]float64, 1200)
	for i := 0; i < 1200; i++ {
		deps[i] = make([]float64, 1200)
		for j := 0; j < 1200; j++ {
			deps[i][j] = math.MinInt64
		}
	}
	// 绘制深度值
	DrawShaderTestPro(NewDepthShader(deps, depthM))
	// 绘制结果 阴影矩阵需要先还原再应用新矩阵
	DrawShaderTestPro(NewLightShader(eyeM, depthM, 1200, 1200, OpenImage("res/reimu.png"),
		deps, lightDir.Mul(eyeM).Normal(), eyePos.Mul(eyeM)))
	//fmt.Println(eyeM)
	//fmt.Println(Min, Max)
	//min, max := float64(math.MaxInt64), float64(math.MinInt64)
	//for i := 0; i < len(deps); i++ {
	//	for j := 0; j < len(deps[i]); j++ {
	//		min = math.Min(min, deps[i][j])
	//		max = math.Max(max, deps[i][j])
	//	}
	//}
	//fmt.Println(min, max)
}

func DrawShaderTestPro(shader IShader) {
	img := NewImg(1200, 1200)
	Fill(img, NewRGB(64, 64, 64))
	model := LoadModel("res/reimu.obj")
	for _, face := range model.Fs {
		triangle := model.ToTriangle(face) // * 600/5
		r1 := shader.Vertex(VArgs{V: triangle.Vs[0], Vt: triangle.Vts[0], Vn: triangle.Vns[0]})
		r2 := shader.Vertex(VArgs{V: triangle.Vs[1], Vt: triangle.Vts[1], Vn: triangle.Vns[1]})
		r3 := shader.Vertex(VArgs{V: triangle.Vs[2], Vt: triangle.Vts[2], Vn: triangle.Vns[2]})
		DrawTrianglePro(img, 1200, 1200, shader, FArgs{NewV: r1.NewV, OldV: r1.OldV, NewVn: r1.NewVn, OldVt: r1.OldVt},
			FArgs{NewV: r2.NewV, OldV: r2.OldV, NewVn: r2.NewVn, OldVt: r2.OldVt}, FArgs{NewV: r3.NewV, OldV: r3.OldV, NewVn: r3.NewVn, OldVt: r3.OldVt})
	}
	SaveImg(img)
}

func DrawPerspTest() {
	// 默认模型文件数据在描述 模型空间
	// 使用 模型矩阵 把物体移动缩放到 世界空间
	// 使用 观察矩阵 把物体移动到观察者为原点的坐标系内(从相机看到的画面) 观察空间
	// 使用 投影矩阵 把所有物体压缩到 x,y,z 的 -1～1 齐次空间
	// 最后使用 viewport 把 ~1到1的内容铺满屏幕
	model := NewTrans(43.2, -45.5, 0).Mul(NewScale(2, 2, 2))                // 模型矩阵
	view := LookAt(NewVec3(0, 0, 106), NewVec3(0, 0, -1), NewVec3(0, 1, 0)) // 观察矩阵
	project := ProjectPersp(math.Pi/4, 1, 1, 10)                            // 投影矩阵 透视投影
	viewport := Viewport(1200, 1200)                                        // 视口变换
	DrawShaderTest(NewImageShader(OpenImage("res/reimu.png"), 1200, 1200,
		viewport.Mul(project).Mul(view).Mul(model), NewVec3(0, 0, 106)))
}

func DrawOrthoTest() {
	// 默认模型文件数据在描述 模型空间
	// 使用 模型矩阵 把物体移动缩放到 世界空间
	// 使用 观察矩阵 把物体移动到观察者为原点的坐标系内(从相机看到的画面) 观察空间
	// 使用 投影矩阵 把所有物体压缩到 x,y,z 的 -1～1 齐次空间
	// 最后使用 viewport 把 ~1到1的内容铺满屏幕
	model := NewTrans(-800, -800, 0).Mul(NewScale(480, 480, 480))           // 模型矩阵
	view := LookAt(NewVec3(1, 1, 1), NewVec3(-1, -1, -1), NewVec3(0, 1, 0)) // 观察矩阵
	project := ProjectOrtho(-1200, 1200, -1200, 1200, -1200, 1200)          // 投影矩阵 正交投影
	viewport := Viewport(1200, 1200)                                        // 视口变换
	DrawShaderTest(NewImageShader(OpenImage("res/reimu.png"), 1200, 1200,
		viewport.Mul(project).Mul(view).Mul(model), NewVec3(1, 1, 1)))
}

func DrawShaderTest(shader IShader) {
	img := NewImg(1200, 1200)
	Fill(img, NewRGB(64, 64, 64))
	model := LoadModel("res/reimu.obj")
	for _, face := range model.Fs {
		triangle := model.ToTriangle(face) // * 600/5
		// 这里已经转换为 2d 坐标了
		r1 := shader.Vertex(VArgs{V: triangle.Vs[0], Vt: triangle.Vts[0], Vn: triangle.Vns[0]})
		r2 := shader.Vertex(VArgs{V: triangle.Vs[1], Vt: triangle.Vts[1], Vn: triangle.Vns[1]})
		r3 := shader.Vertex(VArgs{V: triangle.Vs[2], Vt: triangle.Vts[2], Vn: triangle.Vns[2]})
		DrawTriangle(img, FArgs{NewV: r1.NewV, OldVt: r1.OldVt, NewVn: r1.NewVn}, FArgs{NewV: r2.NewV, OldVt: r2.OldVt, NewVn: r2.NewVn},
			FArgs{NewV: r3.NewV, OldVt: r3.OldVt, NewVn: r3.NewVn}, shader)
	}
	SaveImg(img)
}

func DrawTriangleTest() {
	//img := NewImg(1200, 1200)
	//matrix := LoadModel("res/reimu.obj")
	//shader := NewSimpleShader(NewRGB(255, 0, 0))
	//for _, face := range matrix.Fs {
	//	triangle := matrix.ToTriangle(face) // * 600/5
	//	v0 := triangle.Vs[0].Scale(240)
	//	v1 := triangle.Vs[1].Scale(240)
	//	v2 := triangle.Vs[2].Scale(240)
	//	DrawTriangle(img, v0, v1, v2, shader)
	//}
	//SaveImg(img)
}

func DrawLineTest() {
	img := NewImg(1200, 1200)
	model := LoadModel("res/reimu.obj")
	for _, face := range model.Fs {
		triangle := model.ToTriangle(face) // * 600/5
		v0 := triangle.Vs[0].Scale(240)
		v1 := triangle.Vs[1].Scale(240)
		v2 := triangle.Vs[2].Scale(240)
		DrawLine(img, v0, v1, NewRGB(255, 0, 0))
		DrawLine(img, v1, v2, NewRGB(255, 0, 0))
		DrawLine(img, v2, v0, NewRGB(255, 0, 0))
	}
	SaveImg(img)
}
