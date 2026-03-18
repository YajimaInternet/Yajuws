package main

import (
	"bufio"
	"fmt"
	"github.com/mattn/go-runewidth"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/load"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/process"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	appVersion = "v4.5"
	maxHistory = 50
	notesFile  = "yajuws_notes.txt"
)

type AppState struct {
	start     time.Time
	fastBoot  bool
	useColor  bool
	theme     string
	history   []string
	notes     []string
	favorites map[string]bool
	lastQuote string
}

func main() {
	rand.Seed(time.Now().UnixNano())
	reader := bufio.NewReader(os.Stdin)
	state := newAppState()

	showWarning()
	bootSequence(state)

	for {
		clearScreen()
		printBanner(state)
		fmt.Println("[1] ファイルマネージャ (やりますねぇ！)")
		fmt.Println("[2] システム情報 (24歳、学生です)")
		fmt.Println("[3] 語録ジェネレーター (ファッ！？ 75+5伝説版)")
		fmt.Println("[4] エラー診断 (ファボられてますねぇ)")
		fmt.Println("[5] 便利ツール (時計/タイマー/ミニゲーム)")
		fmt.Println("[6] 語録履歴・お気に入り")
		fmt.Println("[7] 設定")
		fmt.Println("[8] タスクマネージャー (top風)")
		fmt.Println("[9] メモ帳")
		fmt.Println("[0] やめる (終了)")
		fmt.Println()
		fmt.Printf("稼働時間: %s\n", formatDuration(time.Since(state.start)))
		fmt.Printf("保存メモ数: %d\n", len(state.notes))
		fmt.Print("選択肢を入力 (0-9): ")

		choice, ok := readLine(reader)
		if !ok {
			fmt.Println()
			exitMessage(state)
			return
		}

		switch strings.TrimSpace(choice) {
		case "1":
			fileManager(reader)
		case "2":
			systemInfo(reader, state)
		case "3":
			quotes(reader, state)
		case "4":
			diagnosis(reader, state)
		case "5":
			toolsMenu(reader, state)
		case "6":
			historyMenu(reader, state)
		case "7":
			settingsMenu(reader, state)
		case "8":
			taskManager(reader, state)
		case "9":
			notesMenu(reader, state)
		case "0":
			exitMessage(state)
			return
		default:
			// stay on main
		}
	}
}

func newAppState() *AppState {
	return &AppState{
		start:     time.Now(),
		fastBoot:  false,
		useColor:  true,
		theme:     "amber",
		history:   []string{},
		notes:     loadNotes(),
		favorites: map[string]bool{},
	}
}

func showWarning() {
	clearScreen()
	fmt.Println("========================================================")
	fmt.Println("  注意: これはジョークソフトです")
	fmt.Println("  PCに一切影響を与えません。終了で元通り。")
	fmt.Println("  作成者: YajimaNetWorks")
	fmt.Println("========================================================")
	fmt.Println()
}

func bootSequence(state *AppState) {
	if state.fastBoot {
		return
	}

	clearScreen()
	fmt.Printf("Yajuws OS %s 起動シーケンス\n", appVersion)
	fmt.Println("高速起動は設定から切り替えできます")
	fmt.Println()

	steps := []string{
		"BIOSチェック中...OK",
		"野獣プロセッサ起動...OK",
		"語録バッファ展開...OK",
		"wawawa伝説ロード...OK",
		"UI準備中...OK",
	}

	for _, step := range steps {
		fmt.Println(step)
		sleepMillis(120)
	}

	fmt.Println()
	fmt.Println("起動完了。やりますねぇ！")
	sleepMillis(150)
}

func printBanner(state *AppState) {
	banner := []string{
		"   ██████╗██╗  ██╗███████╗    ██╗    ██╗██╗███╗   ██╗███████╗",
		"   ██╔══██╗██║  ██║██╔════╝    ██║    ██║██║████╗  ██║██╔════╝",
		"   ██████╔╝███████║███████╗    ██║ █╗ ██║██║██╔██╗ ██║█████╗",
		"   ██╔══██╗██╔══██║╚════██║    ██║███╗██║██║██║╚██╗██║██╔══╝",
		"   ██║  ██║██║  ██║███████║    ╚███╔███╔╝██║██║ ╚████║███████╗",
		"   ╚═╝  ╚═╝╚═╝  ╚═╝╚══════╝     ╚══╝╚══╝ ╚═╝╚═╝  ╚═══╝╚══════╝",
	}

	for _, line := range banner {
		fmt.Println(applyTheme(state, line))
	}
	fmt.Println()
	fmt.Println(applyTheme(state, fmt.Sprintf("                       やりますねぇ！Yajuws OS %s", appVersion)))
	fmt.Println(applyTheme(state, "                 王道を征く野獣+wawawa伝説システム！"))
	fmt.Println()
}

func fileManager(reader *bufio.Reader) {
	clearScreen()
	fmt.Println()
	entries := [][2]string{
		{"C:Yajuwsyaju.exe", "666MB"},
		{"C:Yajuwsyajusenpai.iso", "24GB"},
		{"C:Yajuwsgoro.txt", "359語録+5伝説 114514KB"},
		{"C:Yajuwswawawa_legend.mp3", "∞MB"},
	}

	leftWidth := 0
	rightWidth := 0
	for _, entry := range entries {
		if w := displayWidth(entry[0]); w > leftWidth {
			leftWidth = w
		}
		if w := displayWidth(entry[1]); w > rightWidth {
			rightWidth = w
		}
	}

	lines := []string{""}
	for _, entry := range entries {
		line := padRightDisplay(entry[0], leftWidth) + strings.Repeat(" ", 3) + padLeftDisplay(entry[1], rightWidth)
		lines = append(lines, line)
	}
	lines = append(lines, "")

	box := drawBox("やりますねぇ！ファイルマネージャ", lines)
	printBoxLines(box, nil)
	fmt.Println()
	fmt.Println("注: 実際のファイルは存在しません。")
	pause(reader)
}

func systemInfo(reader *bufio.Reader, state *AppState) {
	clearScreen()
	fmt.Println()
	printBoxLines(drawBox("システム情報", nil), nil)
	fmt.Println()
	fmt.Printf("OS名: Yajuws OS %s (野獣先輩+wawawa伝説エディション)\n", appVersion)
	fmt.Println("バージョン: 王道を征く！(語録75+5伝説)")
	fmt.Println("CPU: 野獣プロセッサ (24歳学生コア x 114514 + wawawaコア)")
	fmt.Println("RAM: ファッ！？ 伝説語録無限大")
	fmt.Println("ストレージ: 淫夢容量 + wawawa伝説容量")
	fmt.Println()
	fmt.Printf("実ホストOS: %s (%s)\n", runtime.GOOS, runtime.Version())
	if hostOS := os.Getenv("OS"); hostOS != "" {
		fmt.Printf("環境変数OS: %s\n", hostOS)
	}
	fmt.Printf("稼働時間: %s\n", formatDuration(time.Since(state.start)))
	fmt.Println()
	pause(reader)
}

func quotes(reader *bufio.Reader, state *AppState) {
	for {
		clearScreen()
		fmt.Println()
		printBoxLines(drawBox("語録ジェネレーター v4.1", []string{"野獣先輩75語録 + wawawa伝説5語 (5%確率)"}), nil)
		fmt.Println()

		var quote string
		if rand.Intn(20) == 0 {
			quote = randomWawawa()
		} else {
			quote = randomQuote()
		}

		state.lastQuote = quote
		addHistory(state, quote)
		fmt.Println(quote)

		fmt.Println()
		fmt.Println("[Enterで次へ / qでメイン / fでお気に入り登録]")
		fmt.Print(":")
		input, ok := readLine(reader)
		if !ok || strings.EqualFold(strings.TrimSpace(input), "q") {
			return
		}
		if strings.EqualFold(strings.TrimSpace(input), "f") {
			toggleFavorite(state, quote)
		}
	}
}

func diagnosis(reader *bufio.Reader, state *AppState) {
	clearScreen()
	fmt.Println()
	printBoxLines(drawBox("エラー診断ツール", nil), nil)
	fmt.Println()
	fmt.Println("診断中... (359語録+5伝説スキャン)")
	sleepSeconds(2)

	switch rand.Intn(3) {
	case 0:
		fmt.Println("✓ システム正常。野獣先輩+wawawaが守ってます！")
	case 1:
		fmt.Println("⚠️ 軽微なエラー: 語録不足。ジェネレーター(3)で補充を！")
	default:
		fmt.Println("❌ 深刻なエラー: ファッ！？ たまには大人しく死ね(伝説)")
	}
	fmt.Println()
	fmt.Printf("稼働時間: %s\n", formatDuration(time.Since(state.start)))
	fmt.Println()
	pause(reader)
}

func toolsMenu(reader *bufio.Reader, state *AppState) {
	for {
		clearScreen()
		fmt.Println()
		printBoxLines(drawBox("便利ツール", nil), nil)
		fmt.Println()
		fmt.Println("[1] リアルタイム時計")
		fmt.Println("[2] カウントダウンタイマー")
		fmt.Println("[3] 稼働時間チェッカー")
		fmt.Println("[4] じゃんけんミニゲーム")
		fmt.Println("[5] 電卓")
		fmt.Println("[6] カレンダー")
		fmt.Println("[7] 戻る")
		fmt.Println()
		fmt.Print("選択肢を入力 (1-7): ")

		choice, ok := readLine(reader)
		if !ok {
			return
		}
		switch strings.TrimSpace(choice) {
		case "1":
			showClock(reader)
		case "2":
			countdownTimer(reader)
		case "3":
			showUptime(reader, state)
		case "4":
			rockPaperScissors(reader)
		case "5":
			calculator(reader)
		case "6":
			showCalendar(reader)
		case "7":
			return
		default:
		}
	}
}

func taskManager(reader *bufio.Reader, state *AppState) {
	clearScreen()
	fmt.Println("タスクマネージャー (Enterで戻る)")
	stop := make(chan struct{})
	go func() {
		_, _ = readLine(reader)
		close(stop)
	}()

	renderTaskManager(state)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			renderTaskManager(state)
		}
	}
}

func renderTaskManager(state *AppState) {
	cpuUsage := "N/A"
	if percents, err := cpu.Percent(0, false); err == nil && len(percents) > 0 {
		cpuUsage = fmt.Sprintf("%.1f%%", percents[0])
	}

	processCount := "N/A"
	if pids, err := process.Pids(); err == nil {
		processCount = strconv.Itoa(len(pids))
	}

	memSummary := "N/A"
	if info, err := mem.VirtualMemory(); err == nil {
		memSummary = fmt.Sprintf("%s / %s (%.1f%%)", formatBytes(info.Used), formatBytes(info.Total), info.UsedPercent)
	}

	loadSummary := "N/A"
	if avg, err := load.Avg(); err == nil {
		loadSummary = fmt.Sprintf("1m %.2f / 5m %.2f / 15m %.2f", avg.Load1, avg.Load5, avg.Load15)
	}

	clearScreen()
	printBoxLines(drawBox("タスクマネージャー", nil), func(line string) string {
		return applyTheme(state, line)
	})
	fmt.Println()
	fmt.Printf("時刻: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Printf("CPU使用率: %s\n", cpuUsage)
	fmt.Printf("プロセス数: %s\n", processCount)
	fmt.Printf("メモリ使用量: %s\n", memSummary)
	fmt.Printf("ロード平均: %s\n", loadSummary)
	fmt.Printf("稼働時間: %s\n", formatDuration(time.Since(state.start)))
	fmt.Println()
	fmt.Println("Enterで戻る")
}

func historyMenu(reader *bufio.Reader, state *AppState) {
	for {
		clearScreen()
		fmt.Println()
		printBoxLines(drawBox("語録履歴・お気に入り", nil), nil)
		fmt.Println()

		if len(state.history) == 0 {
			fmt.Println("まだ履歴がありません。語録ジェネレーターを回してね。")
		} else {
			fmt.Println("最近の語録 (最新10件):")
			start := len(state.history) - 10
			if start < 0 {
				start = 0
			}
			idx := 1
			for i := len(state.history) - 1; i >= start; i-- {
				quote := state.history[i]
				mark := ""
				if state.favorites[quote] {
					mark = " ★"
				}
				fmt.Printf("%d) %s%s\n", idx, quote, mark)
				idx++
			}
		}

		fmt.Println()
		fmt.Println("[番号]=お気に入り切替 / f=お気に入り一覧 / c=履歴クリア / q=戻る")
		fmt.Print(":")
		input, ok := readLine(reader)
		if !ok {
			return
		}

		cmd := strings.TrimSpace(input)
		if cmd == "" {
			continue
		}
		switch strings.ToLower(cmd) {
		case "q":
			return
		case "f":
			showFavorites(reader, state)
		case "c":
			state.history = []string{}
			fmt.Println("履歴を消去しました。")
			sleepSeconds(1)
		default:
			index, err := strconv.Atoi(cmd)
			if err != nil || index <= 0 {
				continue
			}
			if len(state.history) == 0 {
				continue
			}
			start := len(state.history) - 10
			if start < 0 {
				start = 0
			}
			pos := len(state.history) - index
			if pos < start || pos >= len(state.history) {
				continue
			}
			toggleFavorite(state, state.history[pos])
		}
	}
}

func showFavorites(reader *bufio.Reader, state *AppState) {
	clearScreen()
	fmt.Println()
	printBoxLines(drawBox("お気に入り", nil), nil)
	fmt.Println()

	if len(state.favorites) == 0 {
		fmt.Println("お気に入りはまだ空です。")
	} else {
		i := 1
		for quote := range state.favorites {
			fmt.Printf("%d) %s\n", i, quote)
			i++
		}
	}
	fmt.Println()
	pause(reader)
}

func settingsMenu(reader *bufio.Reader, state *AppState) {
	for {
		clearScreen()
		fmt.Println()
		printBoxLines(drawBox("設定", nil), nil)
		fmt.Println()
		fmt.Printf("高速起動: %s\n", onOff(state.fastBoot))
		fmt.Printf("カラー表示: %s\n", onOff(state.useColor))
		fmt.Printf("テーマ: %s\n", state.theme)
		fmt.Println()
		fmt.Println("[1] 高速起動を切替")
		fmt.Println("[2] カラー表示を切替")
		fmt.Println("[3] テーマ変更 (amber/green/cyan)")
		fmt.Println("[4] 戻る")
		fmt.Println()
		fmt.Print("選択肢を入力 (1-4): ")

		choice, ok := readLine(reader)
		if !ok {
			return
		}
		switch strings.TrimSpace(choice) {
		case "1":
			state.fastBoot = !state.fastBoot
		case "2":
			state.useColor = !state.useColor
		case "3":
			setTheme(reader, state)
		case "4":
			return
		default:
		}
	}
}

func setTheme(reader *bufio.Reader, state *AppState) {
	fmt.Println()
	fmt.Print("テーマ名を入力 (amber/green/cyan): ")
	input, ok := readLine(reader)
	if !ok {
		return
	}
	theme := strings.ToLower(strings.TrimSpace(input))
	switch theme {
	case "amber", "green", "cyan":
		state.theme = theme
	default:
		fmt.Println("不明なテーマです。")
		sleepSeconds(1)
	}
}

func showClock(reader *bufio.Reader) {
	clearScreen()
	fmt.Println("リアルタイム時計: Enterで停止")
	stop := make(chan struct{})
	go func() {
		_, _ = readLine(reader)
		close(stop)
	}()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-stop:
			return
		case t := <-ticker.C:
			clearScreen()
			fmt.Println("リアルタイム時計: Enterで停止")
			fmt.Println()
			fmt.Println(t.Format("2006-01-02 15:04:05"))
		}
	}
}

func countdownTimer(reader *bufio.Reader) {
	clearScreen()
	fmt.Println("カウントダウンタイマー")
	fmt.Print("秒数を入力: ")
	input, ok := readLine(reader)
	if !ok {
		return
	}
	sec, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil || sec <= 0 {
		fmt.Println("正しい秒数を入力してください。")
		sleepSeconds(1)
		return
	}
	for i := sec; i >= 0; i-- {
		clearScreen()
		fmt.Println("カウントダウン中...")
		fmt.Printf("残り: %d 秒\n", i)
		sleepSeconds(1)
	}
	fmt.Println()
	fmt.Println("時間だああああ！")
	pause(reader)
}

func showUptime(reader *bufio.Reader, state *AppState) {
	clearScreen()
	fmt.Println()
	fmt.Println("稼働時間チェッカー")
	fmt.Printf("起動から: %s\n", formatDuration(time.Since(state.start)))
	fmt.Println()
	pause(reader)
}

func rockPaperScissors(reader *bufio.Reader) {
	clearScreen()
	fmt.Println("じゃんけんミニゲーム")
	fmt.Println("[1] グー  [2] チョキ  [3] パー  [4] 戻る")
	fmt.Print("選択: ")
	input, ok := readLine(reader)
	if !ok {
		return
	}
	choice := strings.TrimSpace(input)
	if choice == "4" {
		return
	}
	user, ok := parseHand(choice)
	if !ok {
		fmt.Println("無効な入力です。")
		sleepSeconds(1)
		return
	}
	cpu := rand.Intn(3)
	fmt.Printf("あなた: %s / CPU: %s\n", handName(user), handName(cpu))
	switch judgeRPS(user, cpu) {
	case 0:
		fmt.Println("引き分け！")
	case 1:
		fmt.Println("勝ち！やりますねぇ！")
	case 2:
		fmt.Println("負け…ファッ！？")
	}
	pause(reader)
}

func calculator(reader *bufio.Reader) {
	clearScreen()
	fmt.Println("電卓")
	fmt.Print("1つ目の数値: ")
	leftInput, ok := readLine(reader)
	if !ok {
		return
	}

	left, err := strconv.ParseFloat(strings.TrimSpace(leftInput), 64)
	if err != nil {
		fmt.Println("数値が不正です。")
		sleepSeconds(1)
		return
	}

	fmt.Print("演算子 (+ - * /): ")
	op, ok := readLine(reader)
	if !ok {
		return
	}

	fmt.Print("2つ目の数値: ")
	rightInput, ok := readLine(reader)
	if !ok {
		return
	}

	right, err := strconv.ParseFloat(strings.TrimSpace(rightInput), 64)
	if err != nil {
		fmt.Println("数値が不正です。")
		sleepSeconds(1)
		return
	}

	operator := strings.TrimSpace(op)
	var result float64
	switch operator {
	case "+":
		result = left + right
	case "-":
		result = left - right
	case "*":
		result = left * right
	case "/":
		if right == 0 {
			fmt.Println("0では割れません。")
			sleepSeconds(1)
			return
		}
		result = left / right
	default:
		fmt.Println("演算子が不正です。")
		sleepSeconds(1)
		return
	}

	fmt.Println()
	fmt.Printf("%.4f %s %.4f = %.4f\n", left, operator, right, result)
	fmt.Println()
	pause(reader)
}

func showCalendar(reader *bufio.Reader) {
	clearScreen()
	now := time.Now()
	firstDay := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	lastDay := firstDay.AddDate(0, 1, -1)

	fmt.Println()
	printBoxLines(drawBox("カレンダー", []string{
		now.Format("2006年01月"),
	}), nil)
	fmt.Println()
	fmt.Println("日 月 火 水 木 金 土")

	offset := int(firstDay.Weekday())
	cells := make([]string, 0, offset+lastDay.Day())
	for i := 0; i < offset; i++ {
		cells = append(cells, "  ")
	}
	for day := 1; day <= lastDay.Day(); day++ {
		cells = append(cells, fmt.Sprintf("%2d", day))
	}

	for i := 0; i < len(cells); i += 7 {
		end := i + 7
		if end > len(cells) {
			end = len(cells)
		}
		fmt.Println(strings.Join(cells[i:end], " "))
	}

	fmt.Println()
	fmt.Printf("今日は %d日 (%s)\n", now.Day(), weekdayName(now.Weekday()))
	fmt.Println()
	pause(reader)
}

func notesMenu(reader *bufio.Reader, state *AppState) {
	for {
		clearScreen()
		fmt.Println()
		printBoxLines(drawBox("メモ帳", []string{
			fmt.Sprintf("保存先: %s", notesFile),
			fmt.Sprintf("メモ数: %d", len(state.notes)),
		}), nil)
		fmt.Println()

		if len(state.notes) == 0 {
			fmt.Println("まだメモがありません。")
		} else {
			fmt.Println("最近のメモ (最新5件):")
			start := len(state.notes) - 5
			if start < 0 {
				start = 0
			}
			for i := len(state.notes) - 1; i >= start; i-- {
				fmt.Printf("%d) %s\n", i+1, state.notes[i])
			}
		}

		fmt.Println()
		fmt.Println("[1] メモ追加")
		fmt.Println("[2] 全件表示")
		fmt.Println("[3] メモ削除")
		fmt.Println("[4] 全消去")
		fmt.Println("[5] 戻る")
		fmt.Println()
		fmt.Print("選択肢を入力 (1-5): ")

		choice, ok := readLine(reader)
		if !ok {
			return
		}
		switch strings.TrimSpace(choice) {
		case "1":
			addNote(reader, state)
		case "2":
			showAllNotes(reader, state)
		case "3":
			deleteNote(reader, state)
		case "4":
			clearNotes(reader, state)
		case "5":
			return
		default:
		}
	}
}

func addNote(reader *bufio.Reader, state *AppState) {
	clearScreen()
	fmt.Println("メモ追加")
	fmt.Print("内容を入力: ")
	input, ok := readLine(reader)
	if !ok {
		return
	}

	note := strings.TrimSpace(input)
	if note == "" {
		fmt.Println("空メモは保存できません。")
		sleepSeconds(1)
		return
	}

	state.notes = append(state.notes, note)
	if err := saveNotes(state.notes); err != nil {
		fmt.Printf("保存失敗: %v\n", err)
		sleepSeconds(2)
		return
	}

	fmt.Println("保存しました。")
	sleepSeconds(1)
}

func showAllNotes(reader *bufio.Reader, state *AppState) {
	clearScreen()
	fmt.Println()
	printBoxLines(drawBox("保存済みメモ", nil), nil)
	fmt.Println()

	if len(state.notes) == 0 {
		fmt.Println("保存済みメモはありません。")
	} else {
		for i, note := range state.notes {
			fmt.Printf("%d) %s\n", i+1, note)
		}
	}

	fmt.Println()
	pause(reader)
}

func deleteNote(reader *bufio.Reader, state *AppState) {
	clearScreen()
	fmt.Println("メモ削除")
	if len(state.notes) == 0 {
		fmt.Println("削除できるメモがありません。")
		sleepSeconds(1)
		return
	}

	for i, note := range state.notes {
		fmt.Printf("%d) %s\n", i+1, note)
	}
	fmt.Println()
	fmt.Print("削除する番号を入力: ")
	input, ok := readLine(reader)
	if !ok {
		return
	}

	index, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil || index <= 0 || index > len(state.notes) {
		fmt.Println("番号が不正です。")
		sleepSeconds(1)
		return
	}

	state.notes = append(state.notes[:index-1], state.notes[index:]...)
	if err := saveNotes(state.notes); err != nil {
		fmt.Printf("保存失敗: %v\n", err)
		sleepSeconds(2)
		return
	}

	fmt.Println("削除しました。")
	sleepSeconds(1)
}

func clearNotes(reader *bufio.Reader, state *AppState) {
	clearScreen()
	fmt.Println("メモ全消去")
	if len(state.notes) == 0 {
		fmt.Println("消すメモがありません。")
		sleepSeconds(1)
		return
	}

	fmt.Print("本当に全部消すなら YES と入力: ")
	input, ok := readLine(reader)
	if !ok {
		return
	}
	if strings.TrimSpace(input) != "YES" {
		fmt.Println("キャンセルしました。")
		sleepSeconds(1)
		return
	}

	state.notes = []string{}
	if err := saveNotes(state.notes); err != nil {
		fmt.Printf("保存失敗: %v\n", err)
		sleepSeconds(2)
		return
	}

	fmt.Println("全消去しました。")
	sleepSeconds(1)
}

func parseHand(choice string) (int, bool) {
	switch choice {
	case "1":
		return 0, true
	case "2":
		return 1, true
	case "3":
		return 2, true
	default:
		return 0, false
	}
}

func handName(hand int) string {
	switch hand {
	case 0:
		return "グー"
	case 1:
		return "チョキ"
	case 2:
		return "パー"
	default:
		return "?"
	}
}

func judgeRPS(user int, cpu int) int {
	if user == cpu {
		return 0
	}
	if (user == 0 && cpu == 1) || (user == 1 && cpu == 2) || (user == 2 && cpu == 0) {
		return 1
	}
	return 2
}

func exitMessage(state *AppState) {
	clearScreen()
	fmt.Println()
	fmt.Printf("Yajuws OS %sを終了します。やりますねぇ！またどうぞ！\n", appVersion)
	fmt.Println("(野獣先輩75語録 + wawawa伝説5語 ありがとうございました)")
	fmt.Printf("稼働時間: %s\n", formatDuration(time.Since(state.start)))
	sleepSeconds(2)
}

func readLine(reader *bufio.Reader) (string, bool) {
	line, err := reader.ReadString('\n')
	if err != nil && line == "" {
		return "", false
	}
	return strings.TrimRight(line, "\r\n"), true
}

func pause(reader *bufio.Reader) {
	fmt.Println("Enterで戻る")
	_, _ = readLine(reader)
}

func sleepSeconds(sec int) {
	time.Sleep(time.Duration(sec) * time.Second)
}

func sleepMillis(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

func clearScreen() {
	fmt.Print("\x1b[2J\x1b[H")
}

func formatDuration(d time.Duration) string {
	total := int(d.Seconds())
	h := total / 3600
	m := (total % 3600) / 60
	s := total % 60
	return fmt.Sprintf("%02dh%02dm%02ds", h, m, s)
}

func formatBytes(value uint64) string {
	units := []string{"B", "KB", "MB", "GB", "TB"}
	size := float64(value)
	unit := 0
	for size >= 1024 && unit < len(units)-1 {
		size /= 1024
		unit++
	}
	return fmt.Sprintf("%.1f %s", size, units[unit])
}

func onOff(value bool) string {
	if value {
		return "ON"
	}
	return "OFF"
}

func applyTheme(state *AppState, text string) string {
	if !state.useColor {
		return text
	}
	switch state.theme {
	case "green":
		return "\x1b[32m" + text + "\x1b[0m"
	case "cyan":
		return "\x1b[36m" + text + "\x1b[0m"
	default:
		return "\x1b[33m" + text + "\x1b[0m"
	}
}

func drawBox(title string, lines []string) []string {
	width := displayWidth(title)
	for _, line := range lines {
		if w := displayWidth(line); w > width {
			width = w
		}
	}
	if width < 1 {
		width = 1
	}
	top := "╔" + strings.Repeat("═", width+2) + "╗"
	bottom := "╚" + strings.Repeat("═", width+2) + "╝"
	out := []string{top, "║ " + padCenterDisplay(title, width) + " ║"}
	for _, line := range lines {
		out = append(out, "║ "+padRightDisplay(line, width)+" ║")
	}
	out = append(out, bottom)
	return out
}

func printBoxLines(lines []string, decorate func(string) string) {
	for _, line := range lines {
		if decorate != nil {
			fmt.Println(decorate(line))
		} else {
			fmt.Println(line)
		}
	}
}

func displayWidth(text string) int {
	return runewidth.StringWidth(text)
}

func padRightDisplay(text string, width int) string {
	padding := width - displayWidth(text)
	if padding <= 0 {
		return text
	}
	return text + strings.Repeat(" ", padding)
}

func padCenterDisplay(text string, width int) string {
	padding := width - displayWidth(text)
	if padding <= 0 {
		return text
	}
	left := padding / 2
	right := padding - left
	return strings.Repeat(" ", left) + text + strings.Repeat(" ", right)
}

func padLeftDisplay(text string, width int) string {
	padding := width - displayWidth(text)
	if padding <= 0 {
		return text
	}
	return strings.Repeat(" ", padding) + text
}

func addHistory(state *AppState, quote string) {
	if quote == "" {
		return
	}
	state.history = append(state.history, quote)
	if len(state.history) > maxHistory {
		state.history = state.history[len(state.history)-maxHistory:]
	}
}

func toggleFavorite(state *AppState, quote string) {
	if quote == "" {
		return
	}
	if state.favorites[quote] {
		delete(state.favorites, quote)
	} else {
		state.favorites[quote] = true
	}
}

func loadNotes() []string {
	data, err := os.ReadFile(notesFile)
	if err != nil {
		return []string{}
	}

	text := strings.ReplaceAll(string(data), "\r\n", "\n")
	lines := strings.Split(text, "\n")
	notes := make([]string, 0, len(lines))
	for _, line := range lines {
		note := strings.TrimSpace(line)
		if note != "" {
			notes = append(notes, note)
		}
	}
	return notes
}

func saveNotes(notes []string) error {
	content := ""
	if len(notes) > 0 {
		content = strings.Join(notes, "\n") + "\n"
	}
	return os.WriteFile(notesFile, []byte(content), 0644)
}

func weekdayName(day time.Weekday) string {
	switch day {
	case time.Sunday:
		return "日"
	case time.Monday:
		return "月"
	case time.Tuesday:
		return "火"
	case time.Wednesday:
		return "水"
	case time.Thursday:
		return "木"
	case time.Friday:
		return "金"
	case time.Saturday:
		return "土"
	default:
		return "?"
	}
}

func randomWawawa() string {
	wawawa := []string{
		"🔥 wawaさん「たまには大人しく死ね」[伝説]",
		"🔥 wawaさん「もうボロボロだよ…」[伝説]",
		"🔥 wawaさん「それは重症ですね…」[伝説]",
		"🔥 wawaさん「やめとけ、頭おかしい」[伝説]",
		"🔥 wawaさん「ファッ！？意味わかんねぇ」[伝説]",
	}
	return wawawa[rand.Intn(len(wawawa))]
}

func randomQuote() string {
	quotes := []string{
		"やりますねぇ！",
		"24歳、学生です。",
		"王道を征く！",
		"ファッ！？",
		"気持ちぃ～",
		"こ↑こ↓入って、どうぞ。",
		"流行らせコラ！",
		"喉渇か・・・喉渇かない？",
		"ブッチッパ！",
		"ぬわああああああん疲れたもおおお！",
		"冷えてるか～？",
		"ハイ、ヨロシクゥ！",
		"何やってんだあいつら・・・",
		"おう、考えてやるよ（返すとは言っていない）",
		"おまんこぉ＾～（気さくな挨拶）",
		"気持ちいいですね（建前）気持ちよくはない！（本音）",
		"冗談はよしてくれ（タメ口）",
		"しょうがねぇなぁ（悟空）",
		"なんだこれは、たまげたなぁ。",
		"まずうちさぁ、屋上あんだけど・・・",
		"夜中腹減んないすか？",
		"イキスギィ！",
		"ほら行くどー",
		"先輩、何してるんですか",
		"あーもうめちゃくちゃだよ（呆れ）",
		"お前のことが好きだったんだよ",
		"警察だ！（インパルス板倉）",
		"出そうと思えば（王者の風格）",
		"あっ、おい待てい（江戸っ子）",
		"ケツの穴舐めろ（鬼畜）",
		"お前ノンケかよぉ！（驚愕）",
		"はぁ～～～（クソデカため息）",
		"OK？OK牧場？（激寒）",
		"布団の上で枕を…⁉（最重要事項）",
		"穴は一つしかないから（至言）",
		"やりませんか？行きましょうよ",
		"すっげえキツかったゾ～",
		"田所さん⁉",
		"あっすいません（素）",
		"やべぇ、撃っちゃった（他人事）",
		"よかったら、お話でもしましょうぞ（武士）",
		"多分変態だと思うんですけど（名推理）",
		"入んねぇのか…（落胆）",
		"やばいですね（冷静）",
		"頭にきますよ（憤怒）",
		"痛いんだよおおおおおお！（マジギレ）",
		"うん、おいしい！（味覚障害）",
		"YO！（日顕）",
		"〆鯖ァ！（石田彰）",
		"悔い改めて（戒め）",
		"は？（威圧）",
		"はい、よろしくお願いします！",
		"にゃーん",
		"これはひどい",
		"ぎゃああああ",
		"やるんすか？",
		"眠い",
		"それな",
		"いえーい",
		"ぽんこつ",
		"まじかよ",
		"オラつくな",
		"しょーがねぇな",
		"なんでやねん",
		"ひえー",
		"ふざけんな！",
		"まじ卍",
		"いけるやん",
		"はいはい",
		"まあまあ",
		"どんまい",
		"いいぞ",
		"ファッ！？",
		"どけどけー",
		"すげえよ",
	}
	return quotes[rand.Intn(len(quotes))]
}
