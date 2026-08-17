# gomatrix

`gomatrix` is a fast terminal digital rain (Matrix green code) effect program
built with Go and the `tcell` engine. It is optimized for Linux and Raspberry
Pi 5. It provides a smooth 60 FPS render with zero memory allocation.

## Features

- **High Performance**: Zero memory allocation in the core render loop. This
  prevents Garbage Collection (GC) pauses.
- **Multiple Themes**: Includes four high-quality themes: Classic (Green),
  AntiGravity (Blue/Purple), QuantumGold (Gold), and Cyberpunk (Pink/Green).
- **Auto Fallback**: Switches to ASCII mode if the terminal does not support
  UTF-8.
- **Zero External Dependencies**: Config parsing only uses the Go standard
  library `encoding/json`.

## Requirements

### Build Environment

- Go 1.24 or newer (`sudo apt install golang-go`)

### Recommended Fonts (Debian / Ubuntu / Raspberry Pi OS)

To display half-width Katakana perfectly, please install these open-source
fonts:

```bash
sudo apt install fonts-noto-cjk
```

## Install and Build

```bash
# 1. Get the project
git clone https://github.com/TheRetroMars/gomatrix.git
cd gomatrix

# 2. Build and run
go build -o gomatrix
./gomatrix
```

### Global Install and Uninstall

**Install:**
In the project folder, use Go to install the program globally:

```bash
go install
```

*(Setup Environment Variables)*
To run it from any folder, add `~/go/bin` to your `PATH`:

```bash
echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.bashrc
source ~/.bashrc
```

**Uninstall:**
In the project folder, run this command to remove the program:

```bash
go clean -i
```

## Keybindings
While running, use the following keys:
- `+` / `=` : Increase speed
- `-` / `_` : Decrease speed
- `c` : Cycle through color themes (Classic, AntiGravity, QuantumGold, Cyberpunk)
- `Space` : Trigger a manual matrix flash
- `g` : Toggle CRT bleed effect
- `s` : Toggle Star mode (background flashing)
- `[` / `]` : Decrease/Increase the maximum number of stars
- `m` : Cycle Gradient Mode (Classic, Smooth, VerySmooth)
- `t` : Toggle Render Engine (TrueColor Interpolation vs Dithering)
- `h` : Toggle Help Menu
- `r` : Reset to default speed and theme
- `q` / `Esc` / `Ctrl+C` : Quit

## Configuration
`gomatrix` automatically saves your settings to a config file. It looks for `config.json` in:
1. The path specified by the `-config` flag.
2. `~/.config/gomatrix/config.json`
3. `./gomatrix.json`

If no file is found, it uses built-in defaults. Config example:

```json
{
  "speed": 5,
  "color_theme": "Classic",
  "ascii_only": false,
  "star_mode": false,
  "star_count": 7,
  "crt_mode": false,
  "t1_percent": 15,
  "t2_percent": 50,
  "true_color": true,
  "gradient_mode": 1
}
```

## Raspberry Pi OS Fullscreen (lxterminal F11)

On **Raspberry Pi OS**, the default `lxterminal` does not natively support
`F11` fullscreen and displays an internal menu bar that breaks immersion.
You can run the following script in your terminal to automatically configure
`F11` fullscreen for both **Labwc (Wayland)** and **Openbox (X11)**, and
hide the terminal's internal menu and scrollbar.

```bash
#!/bin/bash
# Raspberry Pi OS Fullscreen Fix for lxterminal
echo "Configuring lxterminal F11 Fullscreen..."

LX_CONF="$HOME/.config/lxterminal/lxterminal.conf"
if [ -f "$LX_CONF" ]; then
    sed -i 's/^hidemenubar=.*/hidemenubar=true/' "$LX_CONF"
    sed -i 's/^hidescrollbar=.*/hidescrollbar=true/' "$LX_CONF"
fi

F11_XML="<keybind key=\"F11\">\n      <action name=\"ToggleFullscreen\" />\n    </keybind>"
LABWC_CONF="$HOME/.config/labwc/rc.xml"
OPENBOX_CONF="$HOME/.config/openbox/lxde-pi-rc.xml"

if [ -f "$LABWC_CONF" ] && ! grep -q "ToggleFullscreen" "$LABWC_CONF"; then
    sed -i "/<\/keyboard>/i \    $F11_XML" "$LABWC_CONF"
    labwc -r || echo "Restart Labwc."
elif [ -f "$OPENBOX_CONF" ] && ! grep -q "ToggleFullscreen" "$OPENBOX_CONF"; then
    sed -i "/<\/keyboard>/i \    $F11_XML" "$OPENBOX_CONF"
    openbox --reconfigure || echo "Restart Openbox."
fi
echo "Done! Press F11 in lxterminal to toggle fullscreen."
```

---

# gomatrix (中文版)

`gomatrix` 是一個使用 Go 語言與 `tcell` 引擎打造的高效能終端機數位雨
(Matrix 綠色代碼) 特效程式。專為 Linux 與 Raspberry Pi 5 最佳化，提供
零記憶體分配 (Zero-Allocation) 的 60 FPS 流暢渲染體驗。

## 特色

- **極致效能**：核心渲染迴圈 0 記憶體分配，避免 GC 卡頓。
- **多種主題**：內建 Classic (經典綠)、AntiGravity (極光藍紫)、
  QuantumGold (太陽金) 與 Cyberpunk (霓虹粉綠) 四種高品質主題。
- **自動降級**：當偵測到不支援 UTF-8 的終端機時，自動降級為純 ASCII 模式。
- **零外部依賴**：除了 `tcell` 終端繪圖庫外，設定檔解析完全依賴 Go 標準庫
  `encoding/json`。

## 環境需求

### 開發與編譯環境

- Go 1.24 或以上版本 (`sudo apt install golang-go`)

### 推薦字型 (Debian / Ubuntu / Raspberry Pi OS)

若要完美呈現半角假名特效，建議安裝以下開源字型套件：

```bash
sudo apt install fonts-noto-cjk
```

## 安裝與建置

```bash
# 1. 取得專案
git clone https://github.com/TheRetroMars/gomatrix.git
cd gomatrix

# 2. 測試編譯與執行
go build -o gomatrix
./gomatrix
```

### 全域安裝與反安裝

**安裝：**
在專案目錄下，您可以透過 Go 內建指令將程式安裝到全域二進位目錄：

```bash
go install
```

*(設定環境變數)*
為了讓您能在任何目錄直接執行，請執行以下指令將 `~/go/bin` 加入 `PATH`：

```bash
echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.bashrc
source ~/.bashrc
```

**反安裝：**
在專案目錄下執行以下指令，Go 會自動清除安裝的二進位檔：

```bash
go clean -i
```

## 快捷鍵
在程式執行時，支援以下快捷鍵操作：
- `+` / `=` : 加速下落
- `-` / `_` : 減速下落
- `c` : 循環切換色彩主題 (Classic, AntiGravity, QuantumGold, Cyberpunk)
- `Space` : 觸發隨機的手動閃爍
- `g` : 開關 CRT 掃描線溢色特效
- `s` : 開關星空模式 (背景自動閃爍)
- `[` / `]` : 減少/增加星星數量上限
- `m` : 切換漸層模式 (Classic, Smooth, VerySmooth)
- `t` : 切換渲染引擎 (TrueColor 內插 與 機率抖動)
- `h` : 開關快捷鍵說明選單
- `r` : 重置為預設速度與主題
- `q` / `Esc` / `Ctrl+C` : 安全退出

## 設定檔
`gomatrix` 會依序在以下路徑尋找 `config.json` 進行套用：
1. `-config` 參數所指定之路徑。
2. `~/.config/gomatrix/config.json`
3. `./gomatrix.json`

若找不到設定檔，程式會使用內建預設值。設定檔範例如下：

```json
{
  "speed": 5,
  "color_theme": "Classic",
  "ascii_only": false,
  "star_mode": false,
  "star_count": 7,
  "crt_mode": false,
  "t1_percent": 15,
  "t2_percent": 50,
  "true_color": true,
  "gradient_mode": 1
}
```

## 樹莓派全螢幕設定 (Raspberry Pi OS)

在 **Raspberry Pi OS** 中，預設的 `lxterminal` 不支援 `F11` 全螢幕，
且內建的選單列會破壞沉浸感。您可以直接複製以下腳本並在終端機執行，
它將自動識別 **Labwc (Wayland)** 或 **Openbox (X11)** 架構，
為您綁定 `F11` 快捷鍵，並徹底隱藏終端機的選單列與捲軸。

```bash
#!/bin/bash
# Raspberry Pi OS Fullscreen Fix for lxterminal
echo "Configuring lxterminal F11 Fullscreen..."

LX_CONF="$HOME/.config/lxterminal/lxterminal.conf"
if [ -f "$LX_CONF" ]; then
    sed -i 's/^hidemenubar=.*/hidemenubar=true/' "$LX_CONF"
    sed -i 's/^hidescrollbar=.*/hidescrollbar=true/' "$LX_CONF"
fi

F11_XML="<keybind key=\"F11\">\n      <action name=\"ToggleFullscreen\" />\n    </keybind>"
LABWC_CONF="$HOME/.config/labwc/rc.xml"
OPENBOX_CONF="$HOME/.config/openbox/lxde-pi-rc.xml"

if [ -f "$LABWC_CONF" ] && ! grep -q "ToggleFullscreen" "$LABWC_CONF"; then
    sed -i "/<\/keyboard>/i \    $F11_XML" "$LABWC_CONF"
    labwc -r || echo "Restart Labwc."
elif [ -f "$OPENBOX_CONF" ] && ! grep -q "ToggleFullscreen" "$OPENBOX_CONF"; then
    sed -i "/<\/keyboard>/i \    $F11_XML" "$OPENBOX_CONF"
    openbox --reconfigure || echo "Restart Openbox."
fi
echo "Done! Press F11 in lxterminal to toggle fullscreen."
```
