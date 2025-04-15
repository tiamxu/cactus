package utils

import (
	"bytes"
	"fmt"
	"math/rand"
	"time"

	svg "github.com/ajstarks/svgo"
)

func GenerateCaptcha(width, height int) ([]byte, string) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 生成4位随机数字或字母（1-9和A-Z，排除容易混淆的字符）
	charSet := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	code := make([]byte, 4)
	for i := range code {
		code[i] = charSet[rng.Intn(len(charSet))]
	}
	captchaText := string(code)

	var svgContent bytes.Buffer
	canvas := svg.New(&svgContent)

	// 1. 背景设置（浅灰色）
	canvas.Start(width, height)
	canvas.Rect(0, 0, width, height, "fill:#f5f5f5")

	// 2. 干扰元素 - 粉色斜线
	canvas.Line(
		10, height/3,
		width-10, height/3*2,
		"stroke:#ff99cc;stroke-width:2",
	)

	// 3. 数字设置 - 每个字符不同颜色
	colors := []string{"#9370DB", "#9ACD32", "#FFA500", "#A52A2A"} // 蓝紫色、黄绿色、橙色、棕色
	charWidth := 20                                                // 单个字符宽度估算
	startX := width/2 - (4*charWidth)/2 + charWidth/2              // 居中计算

	for i, c := range captchaText {
		canvas.Text(
			startX+i*charWidth,
			height/2+8, // 垂直居中微调
			string(c),
			fmt.Sprintf("text-anchor:middle;font-size:24px;font-weight:bold;fill:%s;font-family:Arial, sans-serif", colors[i%len(colors)]),
		)
	}

	canvas.End()

	return svgContent.Bytes(), captchaText
}
