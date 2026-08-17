# gomatrix Design Specification

## 1. Project Basics

- Project Name: gomatrix
- Go Module: gomatrix
- Output Binary: gomatrix
- Target Platforms: Linux (x86_64, ARM64 / Raspberry Pi 5), macOS, Windows

## 2. Tech Stack and Dependencies

- Language: Go (version 1.24 or newer)
- Terminal UI Engine: github.com/gdamore/tcell/v2
- Config Parser: Go standard library `encoding/json` (zero external
  dependencies)
- Timing Control: Go standard library `time.Ticker`

## 3. Performance Specification

- Zero-Allocation Render Loop: Pre-allocate `Grid []Cell` and
  `Columns []Column` during engine initialization. Reuse existing memory slices
  to achieve 0 heap allocations during the render loop. This prevents Go Garbage
  Collection (GC) pauses and lag.
- Diff Rendering: Use `tcell` double buffering. Only output ANSI escape
  sequences for cells where the character or color has changed.
- Raspberry Pi 5 Performance Target: Total CPU usage stays under 1.5% at 60 FPS.

## 4. Display and Font Specification

- Character Sets:
  - Default Mode: Unicode half-width Katakana (U+FF66-U+FF9D), numbers (0-9),
    and symbols.
  - ASCII Fallback Mode: Auto-switch to pure ASCII (letters and numbers) when
    `LC_ALL`, `LANG`, or `TERM` indicates a non-UTF-8 environment, or if
    `TERM=dumb` or `TERM=linux` is detected.
- Recommended Open-Source Fonts (Debian/Ubuntu/PiOS):
  - fonts-noto-cjk
  - fonts-vlgothic
  - fonts-ipafont-gothic

## 5. Color Theme Specification

Themes include head highlight, trail gradient, and background color:

- Classic: Head (White #FFFFFF), Trail (Neon Green #00FF00 to Dark Green
  #003300), Background (Pure Black #000000).
- AntiGravity: Head (Aurora White #F0F8FF), Trail (Cyan #00F0FF to Neon Purple
  #BD00FF to Dark Blue Purple #120033), Background (Deep Void Blue #0A0B1E).
- QuantumGold: Head (Golden White #FFF8DC), Trail (Solar Gold #FFB800 to Plasma
  Orange #FF5500 to Dark Red #330A00), Background (Obsidian #0C0A05).
- Cyberpunk: Head (Pink White #FFE6F2), Trail (Neon Pink #FF007F to Acid Green
  #7FFF00 to Dark Purple #200020), Background (Midnight Black #000000).

## 6. Keyboard Controls

- `+` / `=`: Increase falling speed.
- `-` / `_`: Decrease falling speed.
- `s`: Toggle auto Star background effect.
- `]` / `[`: Increase or decrease automatic Star max limit.
- `Space`: Manually trigger a random column head flash.
- `c`: Cycle through color themes.
- `r`: Reset to default Matrix settings (Classic theme, default speed).
- `q` / `Esc` / `Ctrl+C`: Safely exit and restore terminal state using
  `screen.Fini()`.

## 7. Config File and CLI Specification

- Config Format: JSON (`config.json`).
- Fixed Checking Locations (XDG standard compatible):
  1. `~/.config/gomatrix/config.json`
  2. `./gomatrix.json`
  If neither exists, load built-in defaults.
- CLI Usage: `gomatrix [-config /path/to/config.json]`
- Config Example:

  ```json
  {
    "speed": 5,
    "color_theme": "Classic",
    "ascii_only": false,
    "star_mode": false,
    "star_count": 7,
    "crt_mode": false,
    "t1_percent": 30,
    "t2_percent": 50,
    "true_color": true,
    "gradient_mode": 1
  }
  ```

## 8. Environment Check (Debian/Ubuntu/PiOS)

- **Compiler**: Verify Go is installed via `go version`. Requires Go 1.24+.
  If missing, install via `sudo apt install golang-go`.
- **Fonts**: Verify CJK fonts for half-width Katakana. Check installation via
  `dpkg -l | grep fonts-noto-cjk`.
  If missing, install via `sudo apt install fonts-noto-cjk`.
- **CGO**: Not required. The engine relies on pure Go `tcell`.

## 9. File Structure and Responsibilities

- `config.go`: Configuration structures, JSON parsing, and bounds sanitization.
- `config_test.go`: Unit tests for configuration parsing.
- `engine.go`: Core matrix logic, state updates, random range math,
  and layout calculation.
- `engine_test.go`: Unit tests for matrix state updates.
- `main.go`: Entry point, lifecycle management, and keyboard event loop.
- `main_test.go`: Unit tests for main entry logic.
- `theme.go`: Color palettes and theme management.
- `theme_test.go`: Unit tests for theme configurations.
- `ui.go`: Rendering layer, responsible for gradient math (`getStyle`,
  `lerpColor`, `clampFloat`) and drawing the screen.
- `ui_test.go`: Unit tests for gradient boundary logic and UI constraints.
- `version.go`: Stores the application version string.

---

## gomatrix 完整開發規格書
### 零、 系統架構與檔案職責

- `config.go`：設定檔結構、JSON 解析模組與邊界防呆 (`sanitizeConfig`, `clampSpeed`)。
- `config_test.go`：設定檔解析模組之單元測試。
- `engine.go`：核心引擎，管理矩陣狀態更新、亂數邏輯 (`randomRange`) 與排版。
- `engine_test.go`：引擎狀態更新邏輯之單元測試。
- `main.go`：程式進入點，負責初始化、生命週期與按鍵迴圈 (`cycleNext`)。
- `main_test.go`：主程式邏輯之單元測試。
- `theme.go`：色彩與主題配套配置管理。
- `theme_test.go`：主題配置之單元測試。
- `ui.go`：介面層，專責漸層數學運算 (`getStyle`, `clampFloat`) 與畫面繪製。
- `ui_test.go`：介面單元測試，驗證漸層邊界邏輯與 Config 配置正確性。
- `version.go`：儲存應用程式之版本字串。

### 一、 專案基本資訊與目標平台

- 專案名稱：`gomatrix`
- Go 模組名稱：`gomatrix` (`go mod init gomatrix`)
- 產出可執行檔：`gomatrix`
- 目標執行平台：Linux (x86_64, ARM64 / Raspberry Pi 5), macOS, Windows

### 二、 技術選型與依賴庫

- 開發語言：Go 1.24+
- 終端機 UI 繪圖引擎：`github.com/gdamore/tcell/v2`
- 設定檔解析器：Go 標準庫內建 `encoding/json`（零第三方依賴）
- 時間與渲染控制：Go 標準庫 `time.Ticker`

### 三、 效能與架構規格

- Render Loop 零分配 (Zero-Allocation)：
  於引擎初始化時預先分配全螢幕 `Grid []Cell` 與 `Columns []Column` 切片
  記憶體。主渲染迴圈內部重用既存切片與值型別，確保 Render Loop 達成
  0 Heap Allocation，消除 Go GC 停頓引起的卡頓。
- 差異渲染 (Diff Rendering)：
  利用 `tcell` 雙重緩衝機制，僅對字元或前背景色彩有變更的單元格輸出 ANSI
  轉義序列。
- Raspberry Pi 5 效能指標：於 60 FPS 下，CPU 總佔用率保持在 1.5% 以下。

### 四、 顯示與字型渲染規格

- 字元集 (Character Sets)：
  - 預設模式：Unicode 半角假名 (`U+FF66`–`U+FF9D`) + 數字 + 符號。
  - ASCII 自動降級模式：系統自動檢查環境變數 `LC_ALL`、`LANG` 與 `TERM`。
    若語系非 UTF-8 或 `TERM=dumb`/`linux`，自動切換為純 ASCII 模式。
- 開源字型建議 (Debian / Ubuntu / Raspberry Pi OS)：
  - `fonts-noto-cjk`, `fonts-vlgothic`, `fonts-ipafont-gothic`

### 五、 主題配色配套規格 (Color Themes)

每個主題均包含「頭部白光」、「尾跡漸層」與對應的「背景底色配套」：

- `Classic` (預設)：頭部亮白 (`#FFFFFF`)，尾跡霓虹綠 (`#00FF00`) 降至暗綠
  (`#003300`)，背景純黑 (`#000000`)。
- `AntiGravity`：頭部極光白 (`#F0F8FF`)，尾跡電光青藍 (`#00F0FF`) 漸變至
  霓虹紫 (`#BD00FF`) 降至暗藍紫 (`#120033`)，背景深邃虛空藍 (`#0A0B1E`)。
- `QuantumGold`：頭部閃耀金白 (`#FFF8DC`)，尾跡太陽金黃 (`#FFB800`) 漸變至
  電漿橙 (`#FF5500`) 降至暗焦紅 (`#330A00`)，背景黑曜石底色 (`#0C0A05`)。
- `Cyberpunk`：頭部鮮粉白 (`#FFE6F2`)，尾跡電光螢光粉 (`#FF007F`) 漸變至酸性
  螢光綠 (`#7FFF00`) 降至深夜暗紫 (`#200020`)，背景午夜純黑 (`#000000`)。

### 六、 即時按鍵控制規格 (Keyboard Controls)

- `+` / `=`：增加下落速度
- `-` / `_`：降低下落速度
- `s`：開啟或關閉自動星星 (Star) 背景特效
- `]` / `[`：增加或減少星星的數量上限
- `Space`：手動觸發隨機一條瀑布的頭部閃爍
- `c`：循環切換主題色彩配套
- `m`：切換漸層過渡模式（Classic, Smooth, VerySmooth）
- `t`：切換漸層渲染引擎（TrueColor, Dithering）
- `h`：開啟或關閉多行說明選單
- `r`：重置為 Matrix 預設設定（Classic 主題 + 預設速度）
- `q` / `Esc` / `Ctrl+C`：安全退出程式（執行 `screen.Fini()` 恢復原終端機畫面）

### 七、 設定檔與 CLI 參數規格

- 設定檔格式：JSON (`config.json`)，使用 Go 標準庫 `encoding/json` 解析。
- 固定位置檢查邏輯：
  1. `~/.config/gomatrix/config.json`
  2. `./gomatrix.json`
  若兩者皆不存在，則載入內建預設值。
- `config.json` 欄位與預設值範例：

  ```json
  {
    "speed": 5,
    "color_theme": "Classic",
    "ascii_only": false,
    "star_mode": false,
    "star_count": 7,
    "crt_mode": false,
    "t1_percent": 30,
    "t2_percent": 50,
    "true_color": true,
    "gradient_mode": 1
  }
  ```

- CLI 參數：`gomatrix [-config /path/to/config.json]`

### 八、 環境檢查步驟 (Debian/Ubuntu/PiOS)

- **編譯器檢查**：使用 `go version` 確認 Go 環境 (需 1.24+)。
  若缺失，可透過 `sudo apt install golang-go` 安裝。
- **字型檢查**：使用 `dpkg -l | grep fonts-noto-cjk` 檢查是否具備半角假名
  顯示能力。
  若缺失，可透過 `sudo apt install fonts-noto-cjk` 安裝。
- **CGO 依賴**：無。`tcell` 於 Linux 平台為純 Go 實作，無需額外 C 編譯器。

### 九、 漸層渲染與邊界防護規格

- **色彩錨點 (Color Anchors)**：
  藉由 `t1_percent` 與 `t2_percent` 定義純色段的物理位置。
- **漸層模式 (Gradient Modes)**：
  - `Classic` (0)：錨點作為絕對邊界，瞬間切斷無過渡。
  - `Smooth` (1)：錨點前後建立 45% 緩衝區 (`BufferSmooth`) 進行平滑混色。
  - `VerySmooth` (2)：錨點作為 100% 濃度極值，區段間全範圍連續漸層。
- **渲染引擎 (Render Engines)**：
  - `TrueColor`：運用線性內插 (Lerp) 即時計算 RGB 中間色。
  - `Dithering`：依比例作為機率權重，隨機交錯相鄰純色，相容於舊終端機。
- **邊界防護 (Boundary Defenses)**：
  - **設定檔防呆 (`sanitizeConfig`)**：強制修正負數、超出 100%、或邏輯倒置的 T1/T2，
    並確保模式列舉不越界。
  - **數學箝制 (`clampFloat`)**：嚴格限制內插比例 `p` 在 `0.0~1.0` 之間，
    根絕 `tcell` 底層 RGB 位移損毀 (Bitmask Corruption) 與 Panic 的風險。
  - **除以零防禦**：確保長度計算時 `length <= 0` 不會觸發 Panic。
