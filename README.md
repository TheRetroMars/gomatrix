# gomatrix

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
git clone https://github.com/your-username/gomatrix.git
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
grep -q '$HOME/go/bin' ~/.bashrc || echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.bashrc
source ~/.bashrc
```

**反安裝：**
在專案目錄下執行以下指令，Go 會自動清除安裝的二進位檔：

```bash
go clean -i
```

## 操作說明

程式執行期間，您可以使用以下快捷鍵即時控制：

- `+` / `=`：增加下落速度
- `-` / `_`：降低下落速度
- `c`：循環切換主題色彩配套
- `r`：重置為預設設定（Classic 主題 + 預設速度）
- `q` / `Esc` / `Ctrl+C`：安全退出程式

## 設定檔

您可透過 `-config` 參數指定設定檔路徑：

```bash
gomatrix -config /path/to/config.json
```

若未指定，`gomatrix` 會依序在以下位置尋找
`config.json`：

1. `~/.config/gomatrix/config.json`
2. `./gomatrix.json`

若設定檔不存在，程式將載入內建預設值。設定檔範例：

```json
{
  "speed": 5,
  "color_theme": "Classic",
  "ascii_only": false
}
```
