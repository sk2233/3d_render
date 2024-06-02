/*
@author: sk
@date: 2024/5/30
*/
package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type FaceIndex struct {
	VIndex, VtIndex, VnIndex int // 顶点，纹理坐标，法向量 索引
}

type Face struct {
	Indexes [3]FaceIndex
}

type Triangle struct {
	Vs  [3]Vec // 顶点
	Vts [3]Vec // 纹理坐标
	Vns [3]Vec // 法向量
}

type Model struct {
	Vns []Vec // 法向量
	Vts []Vec // 纹理坐标
	Vs  []Vec // 顶点
	Fs  []Face
}

func (m *Model) ToTriangle(face Face) Triangle {
	res := Triangle{}
	for i := 0; i < 3; i++ {
		index := face.Indexes[i]
		res.Vs[i] = m.Vs[index.VIndex]
		res.Vts[i] = m.Vts[index.VtIndex]
		res.Vns[i] = m.Vns[index.VnIndex]
	}
	return res
}

func LoadModel(path string) *Model {
	open, err := os.Open(path)
	HandleErr(err)
	scanner := bufio.NewScanner(open)
	vns := make([]Vec, 0)
	vts := make([]Vec, 0)
	vs := make([]Vec, 0)
	fs := make([]Face, 0)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "vn") {
			vns = append(vns, parseVn(line))
		} else if strings.HasPrefix(line, "vt") {
			vts = append(vts, parseVt(line))
		} else if strings.HasPrefix(line, "v") {
			vs = append(vs, parseV(line))
		} else if strings.HasPrefix(line, "f") {
			fs = append(fs, parseF(line))
		}
	}
	return &Model{
		Vns: vns,
		Vts: vts,
		Vs:  vs,
		Fs:  fs,
	}
}

func parseF(line string) Face {
	items := strings.Split(line, " ")
	return Face{
		Indexes: [3]FaceIndex{parseFIndex(items[1]), parseFIndex(items[2]), parseFIndex(items[3])},
	}
}

func parseFIndex(line string) FaceIndex {
	items := strings.Split(line, "/")
	return FaceIndex{
		VIndex:  parseInt(items[0]) - 1, // 索引是从 1 开始的
		VtIndex: parseInt(items[1]) - 1,
		VnIndex: parseInt(items[2]) - 1,
	}
}

func parseInt(val string) int {
	res, err := strconv.ParseInt(val, 10, 64)
	HandleErr(err)
	return int(res)
}

func parseV(line string) Vec {
	items := strings.Split(line, " ")
	return NewVec3(parseFloat(items[1]), parseFloat(items[2]), parseFloat(items[3]))
}

func parseVt(line string) Vec {
	items := strings.Split(line, " ")
	return NewVec2(parseFloat(items[1]), parseFloat(items[2]))
}

func parseVn(line string) Vec {
	items := strings.Split(line, " ")
	return NewVec3(parseFloat(items[1]), parseFloat(items[2]), parseFloat(items[3]))
}

func parseFloat(val string) float64 {
	res, err := strconv.ParseFloat(val, 64)
	HandleErr(err)
	return res
}
