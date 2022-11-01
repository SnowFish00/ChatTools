/*
 * @Author: git config SnowFish && git config 3200401354@qq.com
 * @Date: 2022-10-28 10:25:55
 * @LastEditors: git config SnowFish && git 3200401354@qq.com
 * @LastEditTime: 2022-10-28 10:33:31
 * @FilePath: \IM_V2\theme\fonts\Theme.go
 * @Description:
 *
 * Copyright (c) 2022 by snow-fish 3200401354@qq.com, All Rights Reserved.
 */

package Theme

import (
	_ "embed"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

var (
	//go:embed HanChengShiHeYuanFangShouShu-2.ttf
	NotoSansSC []byte
)

type MyTheme struct{}

var _ fyne.Theme = (*MyTheme)(nil)

//	HanChengShiHeYuanFangShouShu-2.ttf 为 fonts 目录下的 ttf 类型的字体文件名
func (m MyTheme) Font(fyne.TextStyle) fyne.Resource {
	return &fyne.StaticResource{
		StaticName:    "HanChengShiHeYuanFangShouShu-2.ttf",
		StaticContent: NotoSansSC,
	}
}

func (*MyTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(n, v)
}

func (*MyTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}

func (*MyTheme) Size(n fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(n)
}
