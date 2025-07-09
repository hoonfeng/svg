package font

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// TextAnchor 定义文本锚点类型
type TextAnchor string

const (
	TextAnchorStart  TextAnchor = "start"
	TextAnchorMiddle TextAnchor = "middle"
	TextAnchorEnd    TextAnchor = "end"
)

// AlignmentBaseline 定义基线对齐类型
type AlignmentBaseline string

const (
	AlignmentBaselineAlphabetic AlignmentBaseline = "alphabetic"
	AlignmentBaselineMiddle     AlignmentBaseline = "middle"
	AlignmentBaselineHanging    AlignmentBaseline = "hanging"
	AlignmentBaselineTop        AlignmentBaseline = "top"
	AlignmentBaselineBottom     AlignmentBaseline = "bottom"
)

// FontMetrics 字体度量信息
type FontMetrics struct {
	Ascent  float64 // 上升高度
	Descent float64 // 下降高度
	Height  float64 // 总高度
	Advance float64 // 字符前进宽度
}

// FontWeight 定义字体粗细类型 / Font weight type definition
type FontWeight string

const (
	FontWeightNormal   FontWeight = "normal"   // 400
	FontWeightBold     FontWeight = "bold"     // 700
	FontWeightBolder   FontWeight = "bolder"   // 相对于父元素更粗 / Bolder than parent
	FontWeightLighter  FontWeight = "lighter"  // 相对于父元素更细 / Lighter than parent
	FontWeight100      FontWeight = "100"      // 最细 / Thinnest
	FontWeight200      FontWeight = "200"      // 很细 / Extra light
	FontWeight300      FontWeight = "300"      // 细 / Light
	FontWeight400      FontWeight = "400"      // 正常 / Normal
	FontWeight500      FontWeight = "500"      // 中等 / Medium
	FontWeight600      FontWeight = "600"      // 半粗 / Semi bold
	FontWeight700      FontWeight = "700"      // 粗 / Bold
	FontWeight800      FontWeight = "800"      // 很粗 / Extra bold
	FontWeight900      FontWeight = "900"      // 最粗 / Black
	FontWeightLight    FontWeight = "light"    // 轻量 / Light
	FontWeightMedium   FontWeight = "medium"   // 中等 / Medium
	FontWeightSemibold FontWeight = "semibold" // 半粗 / Semi bold
	FontWeightBlack    FontWeight = "black"    // 超粗 / Black
)

// FontStyle 定义字体样式类型 / Font style type definition
type FontStyle string

const (
	FontStyleNormal  FontStyle = "normal"  // 正常 / Normal
	FontStyleItalic  FontStyle = "italic"  // 斜体 / Italic
	FontStyleOblique FontStyle = "oblique" // 倾斜 / Oblique
)

// TextStyle 文本样式 / Text style definition
type TextStyle struct {
	FontFamily        string            // 字体族 / Font family
	FontSize          float64           // 字体大小 / Font size
	FontWeight        FontWeight        // 字体粗细 / Font weight (supports numeric and keyword values)
	FontStyle         FontStyle         // 字体样式 / Font style
	TextAnchor        TextAnchor        // 文本锚点 / Text anchor
	AlignmentBaseline AlignmentBaseline // 基线对齐 / Alignment baseline
	Fill              image.Image       // 填充颜色 / Fill color
	Stroke            image.Image       // 描边颜色 / Stroke color
	StrokeWidth       float64           // 描边宽度 / Stroke width
	LetterSpacing     float64           // 字符间距 / Letter spacing
	WordSpacing       float64           // 单词间距 / Word spacing
	TextDecoration    string            // 文本装饰 / Text decoration (underline, overline, line-through)
}

// TextRenderer 是文本渲染器接口
type TextRenderer interface {
	// RenderText 在图像上渲染文本
	RenderText(img draw.Image, text string, x, y float64, style *TextStyle) error

	// MeasureText 测量文本尺寸
	MeasureText(text string, style *TextStyle) (*FontMetrics, error)

	// GetFontMetrics 获取字体度量信息
	GetFontMetrics(style *TextStyle) (*FontMetrics, error)
}

// SVGTextRenderer 是符合SVG标准的文本渲染器实现
type SVGTextRenderer struct {
	fontCache map[string]font.Face // 字体缓存
	fontPaths []string             // 字体搜索路径
}

// NewSVGTextRenderer 创建新的SVG文本渲染器 / Create a new SVG text renderer
func NewSVGTextRenderer() *SVGTextRenderer {
	return &SVGTextRenderer{
		fontCache: make(map[string]font.Face),
		fontPaths: getSystemFontPaths(),
	}
}

// NewSVGTextRendererWithFonts 创建带自定义字体路径的SVG文本渲染器 / Create SVG text renderer with custom font paths
func NewSVGTextRendererWithFonts(customFontPaths []string) *SVGTextRenderer {
	allPaths := append(customFontPaths, getSystemFontPaths()...)
	return &SVGTextRenderer{
		fontCache: make(map[string]font.Face),
		fontPaths: allPaths,
	}
}

// getSystemFontPaths 获取系统字体路径
func getSystemFontPaths() []string {
	return []string{
		"C:\\Windows\\Fonts",     // Windows系统字体目录
		"/System/Library/Fonts",  // macOS系统字体目录
		"/usr/share/fonts",       // Linux系统字体目录
		"/usr/local/share/fonts", // Linux用户字体目录
		"./fonts",                // 项目本地字体目录
	}
}

// normalizeFontWeight 标准化字体权重 / Normalize font weight
func normalizeFontWeight(weight FontWeight) FontWeight {
	switch weight {
	case FontWeight100, FontWeight200, FontWeight300, FontWeightLighter:
		return FontWeightLight
	case FontWeight400, FontWeightNormal:
		return FontWeightNormal
	case FontWeight500:
		return FontWeightMedium
	case FontWeight600:
		return FontWeightSemibold
	case FontWeight700, FontWeightBold:
		return FontWeightBold
	case FontWeight800, FontWeight900, FontWeightBolder:
		return FontWeightBlack
	case FontWeightLight:
		return FontWeightLight
	case FontWeightMedium:
		return FontWeightMedium
	case FontWeightSemibold:
		return FontWeightSemibold
	case FontWeightBlack:
		return FontWeightBlack
	default:
		return FontWeightNormal
	}
}

// getFontWeightIntensity 获取字体权重强度 / Get font weight intensity
func getFontWeightIntensity(weight FontWeight) float64 {
	switch weight {
	case FontWeight100:
		return 0.1
	case FontWeight200:
		return 0.2
	case FontWeight300, FontWeightLight:
		return 0.3
	case FontWeight400, FontWeightNormal:
		return 0.4
	case FontWeight500, FontWeightMedium:
		return 0.5
	case FontWeight600, FontWeightSemibold:
		return 0.6
	case FontWeight700, FontWeightBold:
		return 0.7
	case FontWeight800:
		return 0.8
	case FontWeight900, FontWeightBlack:
		return 0.9
	case FontWeightLighter:
		return 0.2
	case FontWeightBolder:
		return 0.8
	default:
		return 0.4
	}
}



// loadFont 加载字体文件 / Load font file
func (r *SVGTextRenderer) loadFont(fontFamily string, fontSize float64, fontWeight FontWeight, fontStyle FontStyle) (font.Face, error) {
	// 标准化字体权重 / Normalize font weight
	normalizedWeight := normalizeFontWeight(fontWeight)
	
	// 生成字体缓存键 / Generate font cache key
	cacheKey := fmt.Sprintf("%s-%.1f-%s-%s", fontFamily, fontSize, normalizedWeight, fontStyle)

	// 检查缓存 / Check cache
	if face, exists := r.fontCache[cacheKey]; exists {
		return face, nil
	}

	// 对于斜体样式，始终加载正常样式的字体，通过软件变换实现斜体效果 / For italic styles, always load normal style font and implement italic effect through software transformation
	loadFontStyle := fontStyle
	if fontStyle == FontStyleItalic || fontStyle == FontStyleOblique {
		loadFontStyle = FontStyleNormal
	}

	// 尝试加载字体文件 / Try to load font file
	fontFile := r.findFontFile(fontFamily, string(normalizedWeight), string(loadFontStyle))
	if fontFile == "" {
		// 如果找不到字体文件，尝试加载普通样式 / If font file not found, try normal style
		fontFile = r.findFontFile(fontFamily, string(FontWeightNormal), string(FontStyleNormal))
		if fontFile == "" {
			// 最终回退到基础字体 / Final fallback to basic font
			return basicfont.Face7x13, nil
		}
	}

	// 读取字体文件 / Read font file
	fontBytes, err := ioutil.ReadFile(fontFile)
	if err != nil {
		return basicfont.Face7x13, nil // 回退到基础字体 / Fallback to basic font
	}

	// 解析TrueType字体 / Parse TrueType font
	tt, err := truetype.Parse(fontBytes)
	if err != nil {
		return basicfont.Face7x13, nil // 回退到基础字体 / Fallback to basic font
	}

	// 创建字体选项 / Create font options
	options := &truetype.Options{
		Size:    fontSize,
		DPI:     72,
		Hinting: font.HintingFull,
	}

	// 根据字体权重调整字体大小来模拟不同粗细 / Adjust font size based on weight to simulate different thickness
	weightIntensity := getFontWeightIntensity(fontWeight)
	if weightIntensity > 0.4 {
		// 检查是否真的加载了对应权重的字体文件 / Check if corresponding weight font file was actually loaded
		originalFile := r.findFontFile(fontFamily, string(FontWeightNormal), string(FontStyleNormal))
		if fontFile == originalFile {
			// 如果加载的是普通字体，通过调整字体大小来模拟权重效果 / Simulate weight by adjusting font size
			weightFactor := 1.0 + (weightIntensity-0.4)/0.6 * 0.2 // 渐进式权重调整 / Progressive weight adjustment
			options.Size = fontSize * weightFactor
		}
	} else if weightIntensity < 0.4 {
		// 对于较轻的字体权重，稍微减小字体大小 / For lighter weights, slightly reduce font size
		originalFile := r.findFontFile(fontFamily, string(FontWeightNormal), string(FontStyleNormal))
		if fontFile == originalFile {
			weightFactor := 0.9 + weightIntensity * 0.25 // 轻量级权重调整 / Light weight adjustment
			options.Size = fontSize * weightFactor
		}
	}

	// 创建字体面 / Create font face
	face := truetype.NewFace(tt, options)

	// 缓存字体面 / Cache font face
	r.fontCache[cacheKey] = face
	return face, nil
}

// findFontFile 查找字体文件 / Find font file matching the given criteria
func (r *SVGTextRenderer) findFontFile(fontFamily, fontWeight, fontStyle string) string {
	// 字体文件名映射，支持不同样式的字体变体 / Font file name mapping with style variants
	fontMappings := map[string]map[string][]string{
		"sans-serif": {
			"normal": {
				// 中文字体 / Chinese fonts
				"msyh.ttc", "msyh.ttf", "Microsoft YaHei.ttf", "微软雅黑.ttf",
				"simhei.ttf", "SimHei.ttf", "黑体.ttf",
				"simsun.ttc", "simsun.ttf", "SimSun.ttf", "宋体.ttf",
				"NotoSansCJK-Regular.ttc", "SourceHanSansCN-Regular.otf",
				// 英文字体 / English fonts
				"arial.ttf", "Arial.ttf", "DejaVuSans.ttf", "LiberationSans-Regular.ttf",
			},
			"light": {
				// 轻量字体 / Light fonts
				"msyhl.ttf", "Microsoft YaHei Light.ttf", "微软雅黑细体.ttf",
				"NotoSansCJK-Light.ttc", "SourceHanSansCN-Light.otf",
				"ariall.ttf", "Arial-Light.ttf", "DejaVuSans-Light.ttf", "LiberationSans-Light.ttf",
				// 回退到普通字体 / Fallback to normal fonts
				"msyh.ttc", "msyh.ttf", "arial.ttf", "Arial.ttf",
			},
			"medium": {
				// 中等字体 / Medium fonts
				"msyhm.ttf", "Microsoft YaHei Medium.ttf", "微软雅黑中等.ttf",
				"NotoSansCJK-Medium.ttc", "SourceHanSansCN-Medium.otf",
				"arialm.ttf", "Arial-Medium.ttf", "DejaVuSans-Medium.ttf", "LiberationSans-Medium.ttf",
				// 回退到普通字体 / Fallback to normal fonts
				"msyh.ttc", "msyh.ttf", "arial.ttf", "Arial.ttf",
			},
			"semibold": {
				// 半粗字体 / Semi-bold fonts
				"msyhsb.ttf", "Microsoft YaHei Semibold.ttf", "微软雅黑半粗.ttf",
				"NotoSansCJK-DemiBold.ttc", "SourceHanSansCN-SemiBold.otf",
				"arialsb.ttf", "Arial-SemiBold.ttf", "DejaVuSans-SemiBold.ttf", "LiberationSans-SemiBold.ttf",
				// 回退到粗体字体 / Fallback to bold fonts
				"msyhbd.ttc", "msyhbd.ttf", "arialbd.ttf", "Arial-Bold.ttf",
				// 最终回退到普通字体 / Final fallback to normal fonts
				"msyh.ttc", "msyh.ttf", "arial.ttf", "Arial.ttf",
			},
			"bold": {
				// 中文粗体字体 / Chinese bold fonts
				"msyhbd.ttc", "msyhbd.ttf", "Microsoft YaHei Bold.ttf", "微软雅黑粗体.ttf",
				"simhei.ttf", "SimHei.ttf", "黑体.ttf", // 黑体本身较粗
				"NotoSansCJK-Bold.ttc", "SourceHanSansCN-Bold.otf",
				// 英文粗体字体 / English bold fonts
				"arialbd.ttf", "Arial-Bold.ttf", "DejaVuSans-Bold.ttf", "LiberationSans-Bold.ttf",
				// 回退到普通字体 / Fallback to normal fonts
				"msyh.ttc", "msyh.ttf", "arial.ttf", "Arial.ttf",
			},
			"black": {
				// 超粗字体 / Black/Extra bold fonts
				"msyhblk.ttf", "Microsoft YaHei Black.ttf", "微软雅黑超粗.ttf",
				"simhei.ttf", "SimHei.ttf", "黑体.ttf", // 黑体作为超粗字体
				"NotoSansCJK-Black.ttc", "SourceHanSansCN-Heavy.otf",
				"arialblk.ttf", "Arial-Black.ttf", "DejaVuSans-ExtraBold.ttf", "LiberationSans-ExtraBold.ttf",
				// 回退到粗体字体 / Fallback to bold fonts
				"msyhbd.ttc", "msyhbd.ttf", "arialbd.ttf", "Arial-Bold.ttf",
				// 最终回退到普通字体 / Final fallback to normal fonts
				"msyh.ttc", "msyh.ttf", "arial.ttf", "Arial.ttf",
			},
			"italic": {
				// 中文斜体字体（大多数中文字体不支持斜体，使用楷体作为替代）/ Chinese italic fonts
				"simkai.ttf", "SimKai.ttf", "楷体.ttf", // 楷体作为中文斜体替代
				"NotoSansCJK-Regular.ttc", "SourceHanSansCN-Regular.otf",
				// 英文斜体字体 / English italic fonts
				"ariali.ttf", "Arial-Italic.ttf", "DejaVuSans-Oblique.ttf", "LiberationSans-Italic.ttf",
				// 回退到普通字体 / Fallback to normal fonts
				"msyh.ttc", "msyh.ttf", "arial.ttf", "Arial.ttf",
			},
			"bold-italic": {
				// 中文粗斜体字体 / Chinese bold italic fonts
				"simkai.ttf", "SimKai.ttf", "楷体.ttf", // 楷体作为中文粗斜体替代
				"NotoSansCJK-Bold.ttc", "SourceHanSansCN-Bold.otf",
				// 英文粗斜体字体 / English bold italic fonts
				"arialbi.ttf", "Arial-BoldItalic.ttf", "DejaVuSans-BoldOblique.ttf", "LiberationSans-BoldItalic.ttf",
				// 回退到粗体或普通字体 / Fallback to bold or normal fonts
				"msyhbd.ttc", "msyhbd.ttf", "msyh.ttc", "msyh.ttf",
			},
		},
		"serif": {
			"normal": {
				"simsun.ttc", "simsun.ttf", "SimSun.ttf", "宋体.ttf",
				"simkai.ttf", "SimKai.ttf", "楷体.ttf",
				"NotoSerifCJK-Regular.ttc", "SourceHanSerifCN-Regular.otf",
				"times.ttf", "Times.ttf", "DejaVuSerif.ttf", "LiberationSerif-Regular.ttf",
			},
			"bold": {
				"NotoSerifCJK-Bold.ttc", "SourceHanSerifCN-Bold.otf",
				"timesbd.ttf", "Times-Bold.ttf", "DejaVuSerif-Bold.ttf", "LiberationSerif-Bold.ttf",
				"simsun.ttc", "simsun.ttf", "times.ttf", "Times.ttf",
			},
			"italic": {
				"simkai.ttf", "SimKai.ttf", "楷体.ttf",
				"timesi.ttf", "Times-Italic.ttf", "DejaVuSerif-Italic.ttf", "LiberationSerif-Italic.ttf",
				"simsun.ttc", "simsun.ttf", "times.ttf", "Times.ttf",
			},
			"bold-italic": {
				"simkai.ttf", "SimKai.ttf", "楷体.ttf",
				"timesbi.ttf", "Times-BoldItalic.ttf", "DejaVuSerif-BoldItalic.ttf", "LiberationSerif-BoldItalic.ttf",
				"simsun.ttc", "simsun.ttf", "times.ttf", "Times.ttf",
			},
		},
		"monospace": {
			"normal": {
				"msyh.ttc", "msyh.ttf", "Microsoft YaHei.ttf",
				"simhei.ttf", "SimHei.ttf", "simsun.ttc", "simsun.ttf", "SimSun.ttf",
				"consola.ttf", "Consolas.ttf", "SourceCodePro-Regular.ttf", "DejaVuSansMono.ttf",
				"cour.ttf", "Courier.ttf", "LiberationMono-Regular.ttf",
			},
			"bold": {
				"msyhbd.ttc", "msyhbd.ttf", "simhei.ttf", "SimHei.ttf",
				"consolab.ttf", "Consolas-Bold.ttf", "SourceCodePro-Bold.ttf", "DejaVuSansMono-Bold.ttf",
				"courbd.ttf", "Courier-Bold.ttf", "LiberationMono-Bold.ttf",
				"msyh.ttc", "msyh.ttf", "consola.ttf", "Consolas.ttf",
			},
			"italic": {
				"simkai.ttf", "SimKai.ttf", "楷体.ttf",
				"consolai.ttf", "Consolas-Italic.ttf", "SourceCodePro-Italic.ttf", "DejaVuSansMono-Oblique.ttf",
				"couri.ttf", "Courier-Italic.ttf", "LiberationMono-Italic.ttf",
				"msyh.ttc", "msyh.ttf", "consola.ttf", "Consolas.ttf",
			},
			"bold-italic": {
				"simkai.ttf", "SimKai.ttf", "楷体.ttf",
				"consolaz.ttf", "Consolas-BoldItalic.ttf", "SourceCodePro-BoldItalic.ttf", "DejaVuSansMono-BoldOblique.ttf",
				"courbi.ttf", "Courier-BoldItalic.ttf", "LiberationMono-BoldItalic.ttf",
				"msyhbd.ttc", "msyhbd.ttf", "msyh.ttc", "msyh.ttf",
			},
		},
		// 专门的中文字体族 / Dedicated Chinese font families
		"microsoft-yahei": {
			"normal": {
				"msyh.ttc", "msyh.ttf", "Microsoft YaHei.ttf", "微软雅黑.ttf",
			},
			"bold": {
				"msyhbd.ttc", "msyhbd.ttf", "Microsoft YaHei Bold.ttf", "微软雅黑粗体.ttf",
				"msyh.ttc", "msyh.ttf", // 回退到普通微软雅黑
			},
			"italic": {
				"simkai.ttf", "SimKai.ttf", "楷体.ttf", // 使用楷体作为斜体替代
				"msyh.ttc", "msyh.ttf", // 回退到普通微软雅黑
			},
			"bold-italic": {
				"simkai.ttf", "SimKai.ttf", "楷体.ttf", // 使用楷体作为粗斜体替代
				"msyhbd.ttc", "msyhbd.ttf", // 回退到粗体微软雅黑
				"msyh.ttc", "msyh.ttf", // 最终回退到普通微软雅黑
			},
		},
		"simhei": {
			"normal": {
				"simhei.ttf", "SimHei.ttf", "黑体.ttf",
			},
			"bold": {
				"simhei.ttf", "SimHei.ttf", "黑体.ttf", // 黑体本身较粗
			},
			"italic": {
				"simkai.ttf", "SimKai.ttf", "楷体.ttf", // 使用楷体作为斜体替代
				"simhei.ttf", "SimHei.ttf", // 回退到黑体
			},
			"bold-italic": {
				"simkai.ttf", "SimKai.ttf", "楷体.ttf", // 使用楷体作为粗斜体替代
				"simhei.ttf", "SimHei.ttf", // 回退到黑体
			},
		},
		"simsun": {
			"normal": {
				"simsun.ttc", "simsun.ttf", "SimSun.ttf", "宋体.ttf",
			},
			"bold": {
				"simsun.ttc", "simsun.ttf", "SimSun.ttf", "宋体.ttf", // 宋体没有专门的粗体
			},
			"italic": {
				"simkai.ttf", "SimKai.ttf", "楷体.ttf", // 使用楷体作为斜体替代
				"simsun.ttc", "simsun.ttf", // 回退到宋体
			},
			"bold-italic": {
				"simkai.ttf", "SimKai.ttf", "楷体.ttf", // 使用楷体作为粗斜体替代
				"simsun.ttc", "simsun.ttf", // 回退到宋体
			},
		},
		"simkai": {
			"normal": {
				"simkai.ttf", "SimKai.ttf", "楷体.ttf",
			},
			"bold": {
				"simkai.ttf", "SimKai.ttf", "楷体.ttf", // 楷体没有专门的粗体
			},
			"italic": {
				"simkai.ttf", "SimKai.ttf", "楷体.ttf", // 楷体本身就是斜体风格
			},
			"bold-italic": {
				"simkai.ttf", "SimKai.ttf", "楷体.ttf", // 楷体本身就是斜体风格
			},
		},
	}

	// 构建样式键 / Build style key
	styleKey := "normal"
	if fontWeight == "bold" && fontStyle == "italic" {
		styleKey = "bold-italic"
	} else if fontWeight == "bold" {
		styleKey = "bold"
	} else if fontStyle == "italic" {
		styleKey = "italic"
	}

	// 获取字体文件候选列表 / Get font file candidates
	familyKey := strings.ToLower(fontFamily)
	var candidates []string

	if familyMappings, exists := fontMappings[familyKey]; exists {
		if styleCandidates, styleExists := familyMappings[styleKey]; styleExists {
			candidates = styleCandidates
		} else {
			// 如果没有找到特定样式，回退到普通样式 / Fallback to normal style
			candidates = familyMappings["normal"]
		}
	} else {
		// 如果不是标准字体族，尝试直接使用字体族名 / Try using font family name directly
		candidates = []string{
			fontFamily + ".ttf", fontFamily + ".TTF",
			fontFamily + ".ttc", fontFamily + ".TTC",
			fontFamily + ".otf", fontFamily + ".OTF",
		}
	}

	// 在字体路径中搜索 / Search in font paths
	for _, fontPath := range r.fontPaths {
		for _, candidate := range candidates {
			fullPath := filepath.Join(fontPath, candidate)
			if _, err := os.Stat(fullPath); err == nil {
				return fullPath
			}
		}
	}

	return "" // 未找到字体文件 / Font file not found
}

// AddFontPath 添加自定义字体路径 / Add custom font path
func (r *SVGTextRenderer) AddFontPath(fontPath string) {
	// 检查路径是否已存在 / Check if path already exists
	for _, existingPath := range r.fontPaths {
		if existingPath == fontPath {
			return
		}
	}
	r.fontPaths = append(r.fontPaths, fontPath)
}

// LoadFontFromFile 直接从文件加载字体 / Load font directly from file
func (r *SVGTextRenderer) LoadFontFromFile(fontPath, fontFamily string, fontSize float64) error {
	// 检查文件是否存在 / Check if file exists
	if _, err := os.Stat(fontPath); err != nil {
		return fmt.Errorf("字体文件不存在: %s / Font file not found: %s", fontPath, err)
	}

	// 读取字体文件 / Read font file
	fontBytes, err := ioutil.ReadFile(fontPath)
	if err != nil {
		return fmt.Errorf("读取字体文件失败 / Failed to read font file: %v", err)
	}

	// 解析字体 / Parse font
	tt, err := truetype.Parse(fontBytes)
	if err != nil {
		return fmt.Errorf("解析字体文件失败 / Failed to parse font file: %v", err)
	}

	// 创建字体面 / Create font face
	face := truetype.NewFace(tt, &truetype.Options{
		Size:    fontSize,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	// 生成缓存键并存储 / Generate cache key and store
	cacheKey := fmt.Sprintf("%s-%.1f-normal-normal", fontFamily, fontSize)
	r.fontCache[cacheKey] = face

	return nil
}

// RegisterFontFamily 注册字体族映射 / Register font family mapping
func (r *SVGTextRenderer) RegisterFontFamily(familyName string, fontFiles []string) {
	// 这个功能可以让用户自定义字体族映射 / This allows users to define custom font family mappings
	// 由于当前的findFontFile方法使用硬编码映射，这里提供一个扩展点
	// Since current findFontFile uses hardcoded mappings, this provides an extension point
	// 实际实现需要重构findFontFile方法来支持动态映射
	// Actual implementation would require refactoring findFontFile to support dynamic mappings
}

// ClearFontCache 清空字体缓存 / Clear font cache
func (r *SVGTextRenderer) ClearFontCache() {
	r.fontCache = make(map[string]font.Face)
}

// GetLoadedFonts 获取已加载的字体列表 / Get list of loaded fonts
func (r *SVGTextRenderer) GetLoadedFonts() []string {
	fonts := make([]string, 0, len(r.fontCache))
	for key := range r.fontCache {
		fonts = append(fonts, key)
	}
	return fonts
}

// RenderText 在图像上渲染文本 / Render text on image
func (r *SVGTextRenderer) RenderText(img draw.Image, text string, x, y float64, style *TextStyle) error {
	// 加载字体 / Load font
	face, err := r.loadFont(style.FontFamily, style.FontSize, style.FontWeight, style.FontStyle)
	if err != nil {
		return err
	}

	// 测量文本尺寸用于锚点计算 / Measure text for anchor calculation
	metrics, _ := r.MeasureText(text, style)

	// 根据文本锚点调整X坐标 / Adjust X coordinate based on text anchor
	switch style.TextAnchor {
	case TextAnchorMiddle:
		x -= metrics.Advance / 2
	case TextAnchorEnd:
		x -= metrics.Advance
	}

	// 根据基线对齐调整Y坐标 / Adjust Y coordinate based on alignment baseline
	switch style.AlignmentBaseline {
	case AlignmentBaselineMiddle:
		y += metrics.Height / 2
	case AlignmentBaselineHanging:
		y += metrics.Ascent
	case AlignmentBaselineTop:
		y += metrics.Ascent
	case AlignmentBaselineBottom:
		y -= metrics.Descent
	}

	// 检查是否需要软件字体效果 / Check if software font effects are needed
	needsBoldEffect := false
	needsItalicEffect := false

	// 检查是否需要软件粗体效果 / Check if software bold effect is needed
	if style.FontWeight != FontWeightNormal && style.FontWeight != FontWeight100 && style.FontWeight != FontWeight200 && style.FontWeight != FontWeight300 {
		originalFile := r.findFontFile(style.FontFamily, string(FontWeightNormal), string(FontStyleNormal))
		boldFile := r.findFontFile(style.FontFamily, string(style.FontWeight), string(FontStyleNormal))
		needsBoldEffect = (boldFile == originalFile || boldFile == "")
	}

	// 检查是否需要软件斜体效果 / Check if software italic effect is needed
	// 斜体和倾斜都统一使用软件模拟的15度倾斜效果 / Both italic and oblique use software-simulated 15-degree skew effect
	if style.FontStyle == FontStyleItalic || style.FontStyle == FontStyleOblique {
		needsItalicEffect = true
	}

	// 使用标准字体绘制器 / Use standard font drawer
	d := &font.Drawer{
		Dst:  img,
		Src:  style.Fill,
		Face: face,
	}

	// 应用字体效果 / Apply font effects
	if needsBoldEffect && needsItalicEffect {
		// 粗斜体：先应用粗体效果，再应用斜体变换 / Bold italic: apply bold effect first, then italic transformation
		r.renderBoldItalicText(d, text, x, y, style.FontStyle)
	} else if needsBoldEffect {
		// 粗体：多次绘制实现粗体效果 / Bold: multiple draws for bold effect
		r.renderBoldText(d, text, x, y)
	} else if needsItalicEffect {
		// 斜体：使用变换矩阵实现斜体效果 / Italic: use transformation matrix for italic effect
		r.renderItalicText(d, text, x, y, style.FontStyle)
	} else {
		// 普通绘制 / Normal drawing
		d.Dot = fixed.Point26_6{
			X: fixed.Int26_6(x * 64),
			Y: fixed.Int26_6(y * 64),
		}
		d.DrawString(text)
	}

	return nil
}

// renderBoldText 渲染粗体文本 / Render bold text
func (r *SVGTextRenderer) renderBoldText(d *font.Drawer, text string, x, y float64) {
	// 根据字体大小动态调整粗体效果强度 / Dynamically adjust bold effect intensity based on font size
	fontSize := d.Face.Metrics().Height >> 6
	boldStrength := math.Max(0.2, math.Min(1.0, float64(fontSize)/24.0)) // 根据字体大小调整粗体强度 / Adjust bold strength based on font size
	
	// 使用改进的多次绘制实现粗体效果 / Use improved multiple draws for bold effect
	offsets := []struct{ dx, dy float64 }{
		{0, 0},                    // 原始位置 / Original position
		{boldStrength, 0},         // 右偏移 / Right offset
		{0, boldStrength},         // 下偏移 / Down offset
		{boldStrength, boldStrength}, // 右下偏移 / Right-down offset
		{-boldStrength*0.3, 0},    // 轻微左偏移增加厚度 / Slight left offset for thickness
		{0, -boldStrength*0.3},    // 轻微上偏移增加厚度 / Slight up offset for thickness
	}

	for _, offset := range offsets {
		d.Dot = fixed.Point26_6{
			X: fixed.Int26_6((x + offset.dx) * 64),
			Y: fixed.Int26_6((y + offset.dy) * 64),
		}
		d.DrawString(text)
	}
}

// renderItalicText 渲染斜体文本 / Render italic text
func (r *SVGTextRenderer) renderItalicText(d *font.Drawer, text string, x, y float64, fontStyle FontStyle) {
	// 获取文本度量信息 / Get text metrics
	advance := font.MeasureString(d.Face, text)
	metrics := d.Face.Metrics()
	
	// 计算文本边界框 / Calculate text bounding box
	textWidth := int(advance >> 6)
	textHeight := int(metrics.Height >> 6)
	ascent := int(metrics.Ascent >> 6)
	
	// 创建适当大小的临时图像用于斜体变换 / Create appropriately sized temporary image for italic transformation
	padding := textHeight / 2 // 为斜体变换预留空间 / Reserve space for italic transformation
	tempBounds := image.Rect(0, 0, textWidth+padding*2, textHeight+padding)
	tempImg := image.NewRGBA(tempBounds)

	// 在临时图像上绘制文本 / Draw text on temporary image
	tempDrawer := &font.Drawer{
		Dst:  tempImg,
		Src:  d.Src,
		Face: d.Face,
	}
	tempDrawer.Dot = fixed.Point26_6{
		X: fixed.Int26_6(padding * 64),
		Y: fixed.Int26_6(ascent * 64),
	}
	tempDrawer.DrawString(text)

	// 斜体统一使用15度倾斜 / Italic uses 15 degrees skew uniformly
	skewAngle := 15.0

	// 应用改进的斜体变换 / Apply improved italic transformation
	r.applyAdvancedItalicTransform(d.Dst, tempImg, x-float64(padding), y-float64(ascent), skewAngle)
}

// renderBoldItalicText 渲染粗斜体文本 / Render bold italic text
func (r *SVGTextRenderer) renderBoldItalicText(d *font.Drawer, text string, x, y float64, fontStyle FontStyle) {
	// 获取文本度量信息 / Get text metrics
	advance := font.MeasureString(d.Face, text)
	metrics := d.Face.Metrics()
	
	// 计算文本边界框 / Calculate text bounding box
	textWidth := int(advance >> 6)
	textHeight := int(metrics.Height >> 6)
	ascent := int(metrics.Ascent >> 6)
	
	// 为粗斜体效果创建更大的临时图像 / Create larger temporary image for bold italic effect
	padding := textHeight // 为粗斜体变换预留更多空间 / Reserve more space for bold italic transformation
	tempBounds := image.Rect(0, 0, textWidth+padding*2, textHeight+padding)
	tempImg := image.NewRGBA(tempBounds)
	tempDrawer := &font.Drawer{
		Dst:  tempImg,
		Src:  d.Src,
		Face: d.Face,
	}

	// 使用改进的粗体渲染到临时图像 / Render improved bold to temporary image
	r.renderBoldText(tempDrawer, text, float64(padding), float64(ascent))

	// 粗斜体统一使用15度倾斜 / Bold italic uses 15 degrees skew uniformly
	skewAngle := 15.0

	// 对粗体结果应用改进的斜体变换 / Apply improved italic transformation to bold result
	r.applyAdvancedItalicTransform(d.Dst, tempImg, x-float64(padding), y-float64(ascent), skewAngle)
}

// MeasureText 测量文本尺寸
func (r *SVGTextRenderer) MeasureText(text string, style *TextStyle) (*FontMetrics, error) {
	// 加载字体
	face, err := r.loadFont(style.FontFamily, style.FontSize, style.FontWeight, style.FontStyle)
	if err != nil {
		return nil, err
	}

	// 获取字体度量
	fontMetrics := face.Metrics()

	// 测量文本宽度
	advance := font.MeasureString(face, text)

	return &FontMetrics{
		Ascent:  float64(fontMetrics.Ascent) / 64.0,
		Descent: float64(fontMetrics.Descent) / 64.0,
		Height:  float64(fontMetrics.Height) / 64.0,
		Advance: float64(advance) / 64.0,
	}, nil
}

// applyShearTransform 应用斜切变换实现斜体效果 / Apply shear transformation for italic effect
func (r *SVGTextRenderer) applyShearTransform(dst draw.Image, src *image.RGBA, skewAngle float64) {
	bounds := src.Bounds()
	dstBounds := dst.Bounds()
	if dstImg, ok := dst.(*image.RGBA); ok {
		dstBounds = dstImg.Bounds()
	}

	// 计算斜切系数 / Calculate skew factor
	skewFactor := math.Tan(skewAngle * math.Pi / 180.0)

	// 遍历源图像的每个像素 / Iterate through each pixel of source image
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			srcPixel := src.RGBAAt(x, y)
			if srcPixel.A == 0 {
				continue
			}

			// 斜切变换：x' = x + skewFactor * (height - y), y' = y
			// Shear transformation: x' = x + skewFactor * (height - y), y' = y
			newX := float64(x) + skewFactor*float64(bounds.Max.Y-y)
			newY := float64(y)

			// 检查目标坐标是否在边界内 / Check if target coordinates are within bounds
			if newX >= float64(dstBounds.Min.X) && newX < float64(dstBounds.Max.X) &&
				newY >= float64(dstBounds.Min.Y) && newY < float64(dstBounds.Max.Y) {
				// 使用双线性插值进行抗锯齿 / Use bilinear interpolation for anti-aliasing
				r.setPixelWithBlending(dst, newX, newY, srcPixel)
			}
		}
	}
}

// applyAdvancedItalicTransform 应用改进的斜体变换 / Apply advanced italic transformation
func (r *SVGTextRenderer) applyAdvancedItalicTransform(dst draw.Image, src *image.RGBA, offsetX, offsetY, skewAngle float64) {
	skewFactor := math.Tan(skewAngle * math.Pi / 180.0)
	bounds := src.Bounds()
	height := float64(bounds.Dy())
	dstBounds := dst.Bounds()

	// 使用高精度子像素采样提高平滑度 / Use high-precision sub-pixel sampling for better smoothness
	subSamples := 4 // 每个像素分解为4x4子像素 / Decompose each pixel into 4x4 sub-pixels
	subPixelSize := 1.0 / float64(subSamples)
	
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		// 计算当前行的相对位置 / Calculate relative position for current row
		relY := float64(y - bounds.Min.Y)
		
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := src.RGBAAt(x, y)
			if c.A > 0 { // 只处理非透明像素 / Only process non-transparent pixels
				relX := float64(x - bounds.Min.X)
				
				// 对每个像素进行子像素采样 / Sub-pixel sampling for each pixel
				for sy := 0; sy < subSamples; sy++ {
					for sx := 0; sx < subSamples; sx++ {
						// 计算子像素位置 / Calculate sub-pixel position
						subY := relY + float64(sy)*subPixelSize
						subX := relX + float64(sx)*subPixelSize
						
						// 计算子像素的偏移量 / Calculate sub-pixel offset
						subRowOffset := skewFactor * (height - subY)
						targetX := subX + subRowOffset + offsetX
						targetY := subY + offsetY
						
						// 边界检查 / Boundary check
						if targetX >= float64(dstBounds.Min.X) && targetX < float64(dstBounds.Max.X) &&
							targetY >= float64(dstBounds.Min.Y) && targetY < float64(dstBounds.Max.Y) {
							// 计算子像素的alpha值 / Calculate sub-pixel alpha
							subAlpha := float64(c.A) / float64(subSamples*subSamples) / 255.0
							subColor := color.RGBA{
								R: c.R,
								G: c.G,
								B: c.B,
								A: uint8(subAlpha * 255),
							}
							// 使用高精度混合设置子像素 / Set sub-pixel with high-precision blending
							r.setPixelWithHighPrecisionBlending(dst, targetX, targetY, subColor)
						}
					}
				}
			}
		}
	}
}

// setPixelWithHighPrecisionBlending 使用高精度混合算法设置像素 / Set pixel with high-precision blending algorithm
func (r *SVGTextRenderer) setPixelWithHighPrecisionBlending(dst draw.Image, x, y float64, srcColor color.RGBA) {
	// 获取整数坐标 / Get integer coordinates
	x1 := int(math.Floor(x))
	y1 := int(math.Floor(y))
	
	// 计算小数部分 / Calculate fractional parts
	dx := x - float64(x1)
	dy := y - float64(y1)
	
	// 使用高精度双线性插值 / Use high-precision bilinear interpolation
	if dstRGBA, ok := dst.(*image.RGBA); ok {
		bounds := dstRGBA.Bounds()
		
		// 计算四个邻近像素的权重 / Calculate weights for four neighboring pixels
		weights := [4]float64{
			(1.0 - dx) * (1.0 - dy), // 左上 / Top-left
			dx * (1.0 - dy),         // 右上 / Top-right
			(1.0 - dx) * dy,         // 左下 / Bottom-left
			dx * dy,                 // 右下 / Bottom-right
		}
		
		// 四个邻近像素的坐标 / Coordinates of four neighboring pixels
		positions := [4][2]int{
			{x1, y1},         // 左上 / Top-left
			{x1 + 1, y1},     // 右上 / Top-right
			{x1, y1 + 1},     // 左下 / Bottom-left
			{x1 + 1, y1 + 1}, // 右下 / Bottom-right
		}
		
		// 对每个邻近像素应用加权混合 / Apply weighted blending to each neighboring pixel
		for i, pos := range positions {
			px, py := pos[0], pos[1]
			weight := weights[i]
			
			// 跳过权重过小的像素 / Skip pixels with very small weights
			if weight < 0.001 {
				continue
			}
			
			// 边界检查 / Boundary check
			if px >= bounds.Min.X && px < bounds.Max.X && py >= bounds.Min.Y && py < bounds.Max.Y {
				existingPixel := dstRGBA.RGBAAt(px, py)
				alpha := float64(srcColor.A) / 255.0 * weight
				r.blendPixelPrecise(dstRGBA, px, py, srcColor, alpha, existingPixel)
			}
		}
	} else {
		// 回退到简单设置 / Fallback to simple set
		dst.Set(x1, y1, srcColor)
	}
}

// blendPixelPrecise 高精度像素混合 / High-precision pixel blending
func (r *SVGTextRenderer) blendPixelPrecise(dst *image.RGBA, x, y int, srcColor color.RGBA, alpha float64, existingPixel color.RGBA) {
	if alpha < 0.0001 {
		return // 忽略极小的alpha值 / Ignore very small alpha values
	}
	
	// 使用更精确的alpha混合算法 / Use more precise alpha blending algorithm
	srcAlpha := alpha
	dstAlpha := float64(existingPixel.A) / 255.0
	outAlpha := srcAlpha + dstAlpha*(1.0-srcAlpha)
	
	if outAlpha > 0 {
		// 预乘alpha混合 / Premultiplied alpha blending
		newPixel := color.RGBA{
			R: uint8(math.Min(255, (srcAlpha*float64(srcColor.R)+dstAlpha*(1.0-srcAlpha)*float64(existingPixel.R))/outAlpha)),
			G: uint8(math.Min(255, (srcAlpha*float64(srcColor.G)+dstAlpha*(1.0-srcAlpha)*float64(existingPixel.G))/outAlpha)),
			B: uint8(math.Min(255, (srcAlpha*float64(srcColor.B)+dstAlpha*(1.0-srcAlpha)*float64(existingPixel.B))/outAlpha)),
			A: uint8(math.Min(255, outAlpha*255)),
		}
		dst.SetRGBA(x, y, newPixel)
	}
}

// setPixelWithSimpleBlending 使用简化的混合算法设置像素，避免条纹问题 / Set pixel with simplified blending algorithm to avoid stripe artifacts
func (r *SVGTextRenderer) setPixelWithSimpleBlending(dst draw.Image, x, y float64, srcColor color.RGBA) {
	// 获取整数坐标 / Get integer coordinates
	x1 := int(math.Floor(x))
	y1 := int(math.Floor(y))
	
	// 计算小数部分 / Calculate fractional parts
	dx := x - float64(x1)
	dy := y - float64(y1)
	
	// 使用更简单的双线性插值 / Use simpler bilinear interpolation
	if dstRGBA, ok := dst.(*image.RGBA); ok {
		bounds := dstRGBA.Bounds()
		
		// 主像素 / Main pixel
		if x1 >= bounds.Min.X && x1 < bounds.Max.X && y1 >= bounds.Min.Y && y1 < bounds.Max.Y {
			existingPixel := dstRGBA.RGBAAt(x1, y1)
			alpha := float64(srcColor.A) / 255.0 * (1.0 - dx) * (1.0 - dy)
			r.blendPixel(dstRGBA, x1, y1, srcColor, alpha, existingPixel)
		}
		
		// 右邻像素 / Right neighbor
		if dx > 0.1 && x1+1 >= bounds.Min.X && x1+1 < bounds.Max.X && y1 >= bounds.Min.Y && y1 < bounds.Max.Y {
			existingPixel := dstRGBA.RGBAAt(x1+1, y1)
			alpha := float64(srcColor.A) / 255.0 * dx * (1.0 - dy)
			r.blendPixel(dstRGBA, x1+1, y1, srcColor, alpha, existingPixel)
		}
		
		// 下邻像素 / Bottom neighbor
		if dy > 0.1 && x1 >= bounds.Min.X && x1 < bounds.Max.X && y1+1 >= bounds.Min.Y && y1+1 < bounds.Max.Y {
			existingPixel := dstRGBA.RGBAAt(x1, y1+1)
			alpha := float64(srcColor.A) / 255.0 * (1.0 - dx) * dy
			r.blendPixel(dstRGBA, x1, y1+1, srcColor, alpha, existingPixel)
		}
		
		// 右下邻像素 / Bottom-right neighbor
		if dx > 0.1 && dy > 0.1 && x1+1 >= bounds.Min.X && x1+1 < bounds.Max.X && y1+1 >= bounds.Min.Y && y1+1 < bounds.Max.Y {
			existingPixel := dstRGBA.RGBAAt(x1+1, y1+1)
			alpha := float64(srcColor.A) / 255.0 * dx * dy
			r.blendPixel(dstRGBA, x1+1, y1+1, srcColor, alpha, existingPixel)
		}
	} else {
		// 回退到简单设置 / Fallback to simple set
		dst.Set(x1, y1, srcColor)
	}
}

// blendPixel 混合单个像素 / Blend single pixel
func (r *SVGTextRenderer) blendPixel(dst *image.RGBA, x, y int, srcColor color.RGBA, alpha float64, existingPixel color.RGBA) {
	if alpha < 0.01 {
		return // 忽略过小的alpha值 / Ignore very small alpha values
	}
	
	// 简化的alpha混合 / Simplified alpha blending
	invAlpha := 1.0 - alpha
	newPixel := color.RGBA{
		R: uint8(alpha*float64(srcColor.R) + invAlpha*float64(existingPixel.R)),
		G: uint8(alpha*float64(srcColor.G) + invAlpha*float64(existingPixel.G)),
		B: uint8(alpha*float64(srcColor.B) + invAlpha*float64(existingPixel.B)),
		A: uint8(math.Min(255, float64(existingPixel.A)+alpha*255)),
	}
	
	dst.SetRGBA(x, y, newPixel)
}

// setPixelWithAdvancedBlending 使用改进的混合算法设置像素 / Set pixel with advanced blending algorithm
func (r *SVGTextRenderer) setPixelWithAdvancedBlending(dst draw.Image, x, y float64, srcColor color.RGBA) {
	// 获取整数坐标 / Get integer coordinates
	x1 := int(math.Floor(x))
	y1 := int(math.Floor(y))
	x2 := x1 + 1
	y2 := y1 + 1

	// 计算插值权重 / Calculate interpolation weights
	dx := x - float64(x1)
	dy := y - float64(y1)

	// 四个邻近像素的权重 / Weights for four neighboring pixels
	weights := []struct {
		x, y   int
		weight float64
	}{
		{x1, y1, (1 - dx) * (1 - dy)},
		{x2, y1, dx * (1 - dy)},
		{x1, y2, (1 - dx) * dy},
		{x2, y2, dx * dy},
	}

	// 对每个邻近像素进行高质量Alpha混合 / High-quality alpha blend with each neighboring pixel
	for _, w := range weights {
		if w.weight > 0.005 { // 使用更低的阈值以获得更好的质量 / Use lower threshold for better quality
			// 计算加权后的颜色 / Calculate weighted color
			weightedColor := color.RGBA{
				R: srcColor.R,
				G: srcColor.G,
				B: srcColor.B,
				A: uint8(float64(srcColor.A) * w.weight),
			}

			// 设置像素 / Set pixel
			if dstRGBA, ok := dst.(*image.RGBA); ok {
				bounds := dstRGBA.Bounds()
				if w.x >= bounds.Min.X && w.x < bounds.Max.X && w.y >= bounds.Min.Y && w.y < bounds.Max.Y {
					existingPixel := dstRGBA.RGBAAt(w.x, w.y)

					// 改进的Alpha混合算法 / Improved alpha blending algorithm
					alpha := float64(weightedColor.A) / 255.0
					existingAlpha := float64(existingPixel.A) / 255.0
					newAlpha := alpha + existingAlpha*(1-alpha)
					
					var newPixel color.RGBA
					if newAlpha > 0 {
						// 预乘Alpha混合 / Premultiplied alpha blending
						newPixel = color.RGBA{
							R: uint8((alpha*float64(weightedColor.R) + existingAlpha*(1-alpha)*float64(existingPixel.R)) / newAlpha),
							G: uint8((alpha*float64(weightedColor.G) + existingAlpha*(1-alpha)*float64(existingPixel.G)) / newAlpha),
							B: uint8((alpha*float64(weightedColor.B) + existingAlpha*(1-alpha)*float64(existingPixel.B)) / newAlpha),
							A: uint8(math.Min(255, newAlpha*255)),
						}
					} else {
						newPixel = existingPixel
					}

					dstRGBA.SetRGBA(w.x, w.y, newPixel)
				}
			} else {
				dst.Set(w.x, w.y, weightedColor)
			}
		}
	}
}

// setPixelWithBlending 使用双线性插值设置像素并进行抗锯齿
// setPixelWithBlending sets pixel with bilinear interpolation and anti-aliasing
func (r *SVGTextRenderer) setPixelWithBlending(dst draw.Image, x, y float64, srcColor color.RGBA) {
	// 获取整数坐标 / Get integer coordinates
	x1 := int(math.Floor(x))
	y1 := int(math.Floor(y))
	x2 := x1 + 1
	y2 := y1 + 1

	// 计算插值权重 / Calculate interpolation weights
	dx := x - float64(x1)
	dy := y - float64(y1)

	// 四个邻近像素的权重 / Weights for four neighboring pixels
	weights := []struct {
		x, y   int
		weight float64
	}{
		{x1, y1, (1 - dx) * (1 - dy)},
		{x2, y1, dx * (1 - dy)},
		{x1, y2, (1 - dx) * dy},
		{x2, y2, dx * dy},
	}

	// 对每个邻近像素进行Alpha混合 / Alpha blend with each neighboring pixel
	for _, w := range weights {
		if w.weight > 0.01 { // 忽略权重过小的像素 / Ignore pixels with very small weights
			// 计算加权后的颜色 / Calculate weighted color
			weightedColor := color.RGBA{
				R: srcColor.R,
				G: srcColor.G,
				B: srcColor.B,
				A: uint8(float64(srcColor.A) * w.weight),
			}

			// 设置像素 / Set pixel
			if dstRGBA, ok := dst.(*image.RGBA); ok {
				bounds := dstRGBA.Bounds()
				if w.x >= bounds.Min.X && w.x < bounds.Max.X && w.y >= bounds.Min.Y && w.y < bounds.Max.Y {
					existingPixel := dstRGBA.RGBAAt(w.x, w.y)

					// Alpha混合 / Alpha blending
					alpha := float64(weightedColor.A) / 255.0
					invAlpha := 1.0 - alpha

					newPixel := color.RGBA{
						R: uint8(alpha*float64(weightedColor.R) + invAlpha*float64(existingPixel.R)),
						G: uint8(alpha*float64(weightedColor.G) + invAlpha*float64(existingPixel.G)),
						B: uint8(alpha*float64(weightedColor.B) + invAlpha*float64(existingPixel.B)),
						A: uint8(math.Min(255, float64(existingPixel.A)+alpha*255)),
					}

					dstRGBA.SetRGBA(w.x, w.y, newPixel)
				}
			} else {
				dst.Set(w.x, w.y, weightedColor)
			}
		}
	}
}

// bicubicInterpolation 双三次插值函数，提供更高质量的图像重采样 / Bicubic interpolation function for higher quality image resampling
func (r *SVGTextRenderer) bicubicInterpolation(src *image.RGBA, x, y float64) color.RGBA {
	bounds := src.Bounds()

	// 获取整数坐标和小数部分 / Get integer coordinates and fractional parts
	x0 := int(math.Floor(x))
	y0 := int(math.Floor(y))
	dx := x - float64(x0)
	dy := y - float64(y0)

	// 双三次插值需要4x4像素网格 / Bicubic interpolation requires 4x4 pixel grid
	var red, green, blue, alpha float64

	// 对每个颜色通道进行双三次插值 / Perform bicubic interpolation for each color channel
	for i := -1; i <= 2; i++ {
		for j := -1; j <= 2; j++ {
			px := x0 + j
			py := y0 + i

			// 获取像素值，边界外使用透明像素 / Get pixel value, use transparent for out-of-bounds
			var pixel color.RGBA
			if px >= bounds.Min.X && px < bounds.Max.X && py >= bounds.Min.Y && py < bounds.Max.Y {
				pixel = src.RGBAAt(px, py)
			}

			// 计算双三次权重 / Calculate bicubic weights
			weightX := r.cubicWeight(float64(j) - dx)
			weightY := r.cubicWeight(float64(i) - dy)
			weight := weightX * weightY

			// 累加加权像素值 / Accumulate weighted pixel values
			red += weight * float64(pixel.R)
			green += weight * float64(pixel.G)
			blue += weight * float64(pixel.B)
			alpha += weight * float64(pixel.A)
		}
	}

	// 确保值在有效范围内 / Ensure values are within valid range
	red = math.Max(0, math.Min(255, red))
	green = math.Max(0, math.Min(255, green))
	blue = math.Max(0, math.Min(255, blue))
	alpha = math.Max(0, math.Min(255, alpha))

	return color.RGBA{
		R: uint8(red),
		G: uint8(green),
		B: uint8(blue),
		A: uint8(alpha),
	}
}

// cubicWeight 计算双三次插值的权重函数（使用Catmull-Rom样条） / Calculate bicubic interpolation weight function (using Catmull-Rom spline)
func (r *SVGTextRenderer) cubicWeight(t float64) float64 {
	t = math.Abs(t)
	if t <= 1.0 {
		// 核心区域：1.5*t^3 - 2.5*t^2 + 1 / Core region: 1.5*t^3 - 2.5*t^2 + 1
		return 1.5*t*t*t - 2.5*t*t + 1.0
	} else if t <= 2.0 {
		// 边缘区域：-0.5*t^3 + 2.5*t^2 - 4*t + 2 / Edge region: -0.5*t^3 + 2.5*t^2 - 4*t + 2
		return -0.5*t*t*t + 2.5*t*t - 4.0*t + 2.0
	}
	// 超出范围返回0 / Return 0 for out of range
	return 0.0
}

// bilinearInterpolation 双线性插值函数（保留作为备用） / Bilinear interpolation function (kept as backup)
func (r *SVGTextRenderer) bilinearInterpolation(src *image.RGBA, x, y float64) color.RGBA {
	// 获取四个邻近像素的坐标 / Get coordinates of four neighboring pixels
	x1 := int(math.Floor(x))
	y1 := int(math.Floor(y))
	x2 := x1 + 1
	y2 := y1 + 1

	// 计算插值权重 / Calculate interpolation weights
	dx := x - float64(x1)
	dy := y - float64(y1)

	bounds := src.Bounds()

	// 获取四个邻近像素，边界外使用透明像素 / Get four neighboring pixels, use transparent for out-of-bounds
	var c1, c2, c3, c4 color.RGBA

	if x1 >= bounds.Min.X && x1 < bounds.Max.X && y1 >= bounds.Min.Y && y1 < bounds.Max.Y {
		c1 = src.RGBAAt(x1, y1)
	}
	if x2 >= bounds.Min.X && x2 < bounds.Max.X && y1 >= bounds.Min.Y && y1 < bounds.Max.Y {
		c2 = src.RGBAAt(x2, y1)
	}
	if x1 >= bounds.Min.X && x1 < bounds.Max.X && y2 >= bounds.Min.Y && y2 < bounds.Max.Y {
		c3 = src.RGBAAt(x1, y2)
	}
	if x2 >= bounds.Min.X && x2 < bounds.Max.X && y2 >= bounds.Min.Y && y2 < bounds.Max.Y {
		c4 = src.RGBAAt(x2, y2)
	}

	// 双线性插值计算 / Bilinear interpolation calculation
	red := (1-dx)*(1-dy)*float64(c1.R) + dx*(1-dy)*float64(c2.R) + (1-dx)*dy*float64(c3.R) + dx*dy*float64(c4.R)
	green := (1-dx)*(1-dy)*float64(c1.G) + dx*(1-dy)*float64(c2.G) + (1-dx)*dy*float64(c3.G) + dx*dy*float64(c4.G)
	blue := (1-dx)*(1-dy)*float64(c1.B) + dx*(1-dy)*float64(c2.B) + (1-dx)*dy*float64(c3.B) + dx*dy*float64(c4.B)
	alpha := (1-dx)*(1-dy)*float64(c1.A) + dx*(1-dy)*float64(c2.A) + (1-dx)*dy*float64(c3.A) + dx*dy*float64(c4.A)

	// 确保值在有效范围内 / Ensure values are within valid range
	red = math.Max(0, math.Min(255, red))
	green = math.Max(0, math.Min(255, green))
	blue = math.Max(0, math.Min(255, blue))
	alpha = math.Max(0, math.Min(255, alpha))

	return color.RGBA{
		R: uint8(red),
		G: uint8(green),
		B: uint8(blue),
		A: uint8(alpha),
	}
}

// GetFontMetrics 获取字体度量信息
func (r *SVGTextRenderer) GetFontMetrics(style *TextStyle) (*FontMetrics, error) {
	// 加载字体
	face, err := r.loadFont(style.FontFamily, style.FontSize, style.FontWeight, style.FontStyle)
	if err != nil {
		return nil, err
	}

	// 获取字体度量
	fontMetrics := face.Metrics()

	return &FontMetrics{
		Ascent:  float64(fontMetrics.Ascent) / 64.0,
		Descent: float64(fontMetrics.Descent) / 64.0,
		Height:  float64(fontMetrics.Height) / 64.0,
		Advance: 0, // 对于字体度量，Advance通常为0
	}, nil
}

// 辅助函数：创建纯色图像
func CreateSolidColor(c color.Color) *image.Uniform {
	return image.NewUniform(c)
}

// DefaultTextRenderer 是默认的文本渲染器
var DefaultTextRenderer TextRenderer = NewSVGTextRenderer()

// NewTextStyle 创建新的文本样式
func NewTextStyle() *TextStyle {
	return &TextStyle{
		FontFamily:        "sans-serif",
		FontSize:          16,
		FontWeight:        "normal",
		FontStyle:         "normal",
		TextAnchor:        TextAnchorStart,
		AlignmentBaseline: AlignmentBaselineAlphabetic,
		Fill:              CreateSolidColor(color.RGBA{0, 0, 0, 255}), // 黑色
		Stroke:            nil,
		StrokeWidth:       0,
	}
}

// DefaultTextStyle 创建默认文本样式
func DefaultTextStyle() *TextStyle {
	return &TextStyle{
		FontFamily:        "sans-serif",
		FontSize:          12,
		FontWeight:        "normal",
		FontStyle:         "normal",
		TextAnchor:        TextAnchorStart,
		AlignmentBaseline: AlignmentBaselineAlphabetic,
		Fill:              CreateSolidColor(color.RGBA{0, 0, 0, 255}), // 黑色
		Stroke:            nil,
		StrokeWidth:       0,
	}
}
