package main

import (
	"encoding/json" // JSONのエンコード/デコード用パッケージ
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/natefinch/lumberjack"               // ログローテーションライブラリ
	"github.com/playwright-community/playwright-go" // Playwright-Goをインポート
)

type Message struct {
	Message string `json:"Message"` // JSONのフィールド名を指定
}

func main() {
	logFile := &lumberjack.Logger{
		Filename:   "./logs/my_application.log", // ログファイルのパス
		MaxSize:    1,                           // MB単位。この例では1MBを超えるとローテーション
		MaxBackups: 3,                           // 保持するバックアップログファイルの最大数
		MaxAge:     28,                          // 日単位。ログファイルを保持する最大日数
		Compress:   true,                        // ローテーションされたファイルをgzipで圧縮するかどうか
	}

	// 標準のロガーの出力を設定
	// logFile に加えて、標準エラー出力 (os.Stderr) にもログを出力するように設定
	// io.MultiWriter を使うことで、複数の Writer に同時に書き込めます。
	mw := os.Stderr
	log.SetOutput(mw)      // 開発中は両方に出すことが多いですが、本番ではlogFileだけにするなど調整します。
	log.SetOutput(logFile) // ログ出力を lumberjack に設定
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	log.Println("サーバー起動中...")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // デフォルトのポート
		log.Printf("環境変数PORTが設定されていないため、デフォルトポート %s を使用します。", port)
	} else {
		log.Printf("環境変数PORTからポート %s を取得しました。", port)
	}
	// httpサーバーを起動
	// ここではバックグラウンドでHTTPサーバーを起動して、後で値を取得するためのエンドポイントを提供します

	// HTTPサーバーをバックグラウンドで起動
	http.HandleFunc("/GeneralCsv", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "POSTメソッドのみ許可されています", http.StatusMethodNotAllowed)
			return
		}
		// 受信したデータをログに出力
		txtID2 := r.FormValue("txtID2")
		txtID1 := r.FormValue("txtID1")
		txtPass := r.FormValue("txtPass")
		w.Header().Set("Content-Type", "application/json")
		if txtID2 == "" || txtID1 == "" || txtPass == "" { // いずれかの値が空の場合,responseにエラーメッセージを返す
			returnJson(w, Message{Message: "txtID2, txtID1, txtPassのいずれかが空です。"})
			return
		}
		go func() {
			// Playwrightを使ってウェブサイトをスクレイピング
			err := getPage(txtID2, txtID1, txtPass)
			if err != nil {
				log.Printf("スクレイピング中にエラーが発生しました: %v", err)
			}

		}()
		w.WriteHeader(http.StatusOK)
		returnJson(w, Message{Message: "スクレイピングを開始しました。"})

	})
	log.Printf("HTTPサーバーを :%s で起動します", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("HTTPサーバーの起動に失敗しました: %v", err)
	}

	// 例えば、http://localhost:8080/value にPOSTリクエストを送信して値を取得することができます
	// ここでは、Playwrightを使ってウェブサイトをスクレイピング

	// ...
}

// return json.NewEncoder(w).Encode(msg) // JSONエンコードしてレスポンスに書き込む
func returnJson(w http.ResponseWriter, msg Message) {
	if err := json.NewEncoder(w).Encode(msg); err != nil {
		log.Printf("JSONエンコードエラー: %v", err)
		return
	}
}

func getPage(txtID2 string, txtID1 string, txtPass string) error {

	if txtID2 == "" || txtID1 == "" || txtPass == "" {
		return errors.New("txtID2, txtID1, txtPassのいずれかが空です。")
	}
	// Playwrightのインストール（初回のみ、またはCI/CDなどで）
	err := playwright.Install()
	if err != nil {
		log.Fatalf("Playwright のインストールに失敗しました: %v", err)
	}
	log.Println("Playwright ブラウザがインストールされました！")

	err = os.MkdirAll("./file", 0755)
	if err != nil {
		log.Fatalf("fileディレクトリの作成に失敗しました: %v", err)
	} else {
		log.Println("fileディレクトリが作成されました。")
	}

	// 環境変数またはhttp経由で値を取得する例
	// ここでは http://localhost:8080/value から値を取得して変数に格納

	//fileディレクトリになにかfileが存在するか確認
	files, err := os.ReadDir("./file")
	if err != nil {
		log.Fatalf("fileディレクトリの読み取りに失敗しました: %v", err)
	} else {
		if len(files) > 0 {
			log.Println("fileディレクトリに既存のファイルがあります。")
			// 既存のファイルを削除
			for _, file := range files {
				err := os.Remove("./file/" + file.Name())
				if err != nil {
					log.Printf("ファイル '%s' の削除に失敗しました: %v", file.Name(), err)
				} else {
					log.Printf("ファイル '%s' を削除しました。\n", file.Name())
				}
			}
			log.Println("fileディレクトリの既存のファイルを削除しました。")
		} else {
			log.Println("fileディレクトリは空です。")
		}
	}

	// Playwrightの起動
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("Playwright の起動に失敗しました: %v", err)
	}
	defer pw.Stop() // プログラム終了時にPlaywrightを確実に停止

	// ブラウザの起動 (ヘッドレスモードがデフォルト)
	// GUIを表示したい場合は Launch(playwright.BrowserTypeLaunchOptions{Headless: playwright.Bool(false)}) を使う
	//gui を表示
	// ブラウザの起動
	log.Println("ブラウザを起動しています...")
	// ヘッドレスモードを無効にしてGUIを表示する場合は、Headless: playwright.Bool(false) を指定
	// browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: playwright.Bool(false)})
	browser, err := pw.Chromium.Launch()
	if err != nil {
		log.Fatalf("ブラウザの起動に失敗しました: %v", err)
	}
	defer browser.Close() // プログラム終了時にブラウザを確実に閉じる

	// 新しいページ (タブ) の作成
	page, err := browser.NewPage()
	if err != nil {
		log.Fatalf("ページの作成に失敗しました: %v", err)
	}

	// 目的のURLに移動
	targetURL := "http://theearth-np.com/F-OES1010[Login].aspx" // スクレイピングしたいウェブサイトのURLに変更してください
	log.Printf("URLにアクセス中: %s", targetURL)
	_, err = page.Goto(targetURL)
	if err != nil {
		log.Fatalf("URLへの移動に失敗しました: %v", err)
	}

	// ページのタイトルを取得
	title, err := page.Title()
	if err != nil {
		log.Printf("タイトル取得中にエラー: %v", err)
		title = "取得できませんでした"
		return err
	}
	log.Printf("ページのタイトル: %s\n", title)

	//#popup_1を探してクリック
	clickSelector(page, "#popup_1", 3000) // ポップアップを閉じるためのセレクターをクリック
	clickSelector(page, "#txtID2")
	inputSelector(page, "#txtID2", txtID2) // ユーザー名を入力
	clickSelector(page, "#txtID1")
	inputSelector(page, "#txtID1", txtID1) // ユーザー名を入力
	clickSelector(page, "#txtPass")
	inputSelector(page, "#txtPass", txtPass) // ユーザー名を入力

	takeScreenshot(page, "screenshot.png") // スクリーンショットを撮る
	clickSelector(page, "#imgLogin")       // ログインボタンをクリック
	// 3秒待機してからポップアップを閉じる
	//表示されなかった場合はそのまま次の処理に進む
	log.Println("ログインボタンをクリックしました。3秒待機します。")
	// ポップアップが表示される場合は、#popup_1 セレクターをクリックして閉じる
	// ポップアップが表示されるまで待機してからクリック
	// ポップアップが表示されるまで待機
	// #popup_1 が表示されるまで待機してからクリック

	takeScreenshot(page, "screenshot_01_afterLoginButton.png")
	page.Locator("#popup_1").WaitFor(playwright.LocatorWaitForOptions{
		State:   playwright.WaitForSelectorStateVisible,
		Timeout: playwright.Float(3000),
	})
	// popup_1 が表示されたか確認

	// ポップアップが表示されたらクリック
	if exists, _ := selectorExists(page, "#popup_1"); exists {
		log.Println("ポップアップが表示されました。クリックして閉じます。")
		//#popup_1のvalueを確認
		value, err := page.Locator("#popup_1").Evaluate("el => el.value", nil)
		if err != nil {
			log.Printf("ポップアップの値の取得に失敗しました: %v", err)
		} else {
			if value == "接続ユーザー確認" {
				//popup を新しいタブで開く
				// 新しいウィンドウ（ページ）が開くのを待つ
				popupPage, err := page.ExpectPopup(func() error {
					return clickSelector(page, "#popup_1", 3000) // ポップアップを閉じる
					// return clickSelector(page, "#Button1st_2") // 例: ボタンをクリックして新しいウィンドウを開く
				})
				if err != nil {
					log.Println("新しいウィンドウの取得に失敗しました")
				} else {
					log.Printf("新しいウィンドウを捕捉しました: %s\n", popupPage.URL())
					//10秒待機してポップアップの内容を確認
					time.Sleep(10 * time.Second)
					//tr td の内容を取得
					rows, err := popupPage.Locator("tr").All()
					if err != nil {
						log.Printf("ポップアップのテーブル行の取得に失敗しました: %v", err)
					} else {
						log.Printf("ポップアップのテーブル行の数: %d\n", len(rows))
						for i, row := range rows {
							// 各行の内容を取得
							cells, err := row.Locator("td").All()
							if err != nil {
								log.Printf("ポップアップのテーブル行 %d のセルの取得に失敗しました: %v", i, err)
								continue
							}
							// log.Printf("ポップアップのテーブル行 %d のセルの数: %d\n", i, len(cells))
							for j, cell := range cells {
								// 各セルの内容を取得
								cellText, err := cell.InnerText()
								if err != nil {
									log.Printf("ポップアップのテーブル行 %d のセル %d の内容の取得に失敗しました: %v", i, j, err)
									continue
								}
								if j == 2 && !inArray(cellText, []string{"auto2", "auto1", "auto3", "autoload"}) {
									log.Printf("ポップアップのテーブル行 %d のセル %d の内容: %s\n", i, j, cellText)
								}
							}
						}
						// ここで popupPage に対して操作ができます
					}
					log.Println("ポップアップを閉じました。")

				}

				log.Printf("ポップアップの値: %v\n", value)
			}
			clickSelector(page, "#popup_1", 3000) // ポップアップを閉じる
		}
	}
	//#Button1st_2が表示されるまで待機してからクリック
	err = page.Locator("#Button1st_2").WaitFor(playwright.LocatorWaitForOptions{
		State:   playwright.WaitForSelectorStateVisible,
		Timeout: playwright.Float(10000),
	})
	if err != nil {
		log.Printf("Button1st_2の表示待機中にエラーが発生しました: %v", err)
		return err
	}
	takeScreenshot(page, "screenshot_02_afterLogin.png") // スクリーンショットを撮る

	log.Println("ログインが完了しました。")

	//3秒待機
	log.Println("3秒待機します。")
	time.Sleep(3 * time.Second)

	Button1st_2, err := page.Locator("#Button1st_2").Count()
	if err != nil {
		log.Printf("Button1st_2のカウント取得中にエラーが発生しました: %v", err)
		return err
	}
	if Button1st_2 == 0 {
		log.Println("Button1st_2が見つかりませんでした。ログインに失敗した可能性があります。")
	}
	page.Locator("Button2nd_5").WaitFor(playwright.LocatorWaitForOptions{
		State:   playwright.WaitForSelectorStateVisible,
		Timeout: playwright.Float(10000),
	})
	// Button1st_2が存在する場合はクリック
	log.Println("Button1st_2が存在します。クリックします。")
	clickSelector(page, "#Button1st_2", 3000)   // Button1st_2をクリック
	waitForSelector(page, "#Button2nd_5", 3000) // Button2nd_5が表示されるまで待機
	clickSelector(page, "#Button2nd_5", 3000)   // Button1st_2をクリック
	waitForSelector(page, "#Button3rd_0", 3000) // Button2nd_5が表示されるまで待機
	clickSelector(page, "#Button3rd_0", 3000)   // Button1st_2をクリック

	waitForSelector(page, "#rdoSelect1", 10000) // 日付入力フィールドが表示されるまで待機
	//https://theearth-np.com/F-NOS3010[GeneralCsv].aspxに移動
	// targetURL = "https://theearth-np.com/F-NOS3010[GeneralCsv].aspx"
	// log.Printf("次のURLにアクセス中: %s", targetURL)
	// _, err = page.Goto(targetURL)
	// if err != nil {
	// 	log.Fatalf("次のURLへの移動に失敗しました: %v", err)
	// }
	// ページのタイトルを取得
	title, err = page.Title()
	if err != nil {
		log.Printf("次のページのタイトル取得中にエラー: %v", err)
		title = "取得できませんでした"
		return err
	}
	log.Printf("ページのタイトル: %s\n", title)
	err = clickSelector(page, "#rdoSelect1", 3000) // ポップアップを閉じる
	if err != nil {
		return err
	} // エラーが発生した場合は終了
	clickSelector(page, "#rdoDate1", 3000) // ポップアップを閉じる
	//日付をyesterday_yy, yesterday_mm, yesterday_ddに設定
	// 昨日の日付を取得
	yesterday := time.Now().AddDate(0, 0, -1)
	yesterdayYY := fmt.Sprintf("%02d", yesterday.Year()%100) // 年を2桁で取得
	yesterdayMM := fmt.Sprintf("%02d", int(yesterday.Month()))
	yesterdayDD := fmt.Sprintf("%02d", yesterday.Day())
	//今日の日付を取得
	today := time.Now()
	todayYY := fmt.Sprintf("%02d", today.Year()%100) // 年を2桁で取得
	todayMM := fmt.Sprintf("%02d", int(today.Month()))
	todayDD := fmt.Sprintf("%02d", today.Day())

	clickSelector(page, "#MainContent_ucStartDate_txtYear")
	inputSelector(page, "#MainContent_ucStartDate_txtYear", yesterdayYY)
	clickSelector(page, "#MainContent_ucStartDate_txtMonth")
	inputSelector(page, "#MainContent_ucStartDate_txtMonth", yesterdayMM)
	clickSelector(page, "#MainContent_ucStartDate_txtDay")
	inputSelector(page, "#MainContent_ucStartDate_txtDay", yesterdayDD)
	clickSelector(page, "#MainContent_ucEndDate_txtYear")
	inputSelector(page, "#MainContent_ucEndDate_txtYear", todayYY)
	clickSelector(page, "#MainContent_ucEndDate_txtMonth")
	inputSelector(page, "#MainContent_ucEndDate_txtMonth", todayMM)
	clickSelector(page, "#MainContent_ucEndDate_txtDay")
	inputSelector(page, "#MainContent_ucEndDate_txtDay", todayDD)

	clickSelector(page, "#btnCsv") // 検索ボタンをクリック
	//ダウンロードが完了するまで待機
	log.Println("CSVダウンロードを開始しました。ダウンロードが完了するまで待機します。")
	// ダウンロードが完了するまで待機
	download, err := page.ExpectDownload(func() error {
		return nil // 既にクリック済みなので何もしない
	}, playwright.PageExpectDownloadOptions{
		Timeout: playwright.Float(60000), // 60秒待機
	})
	if err != nil {
		log.Printf("ダウンロードの待機中にエラーが発生しました: %v", err)
	} else {
		log.Printf("ダウンロードが完了しました: %s", download.URL())
	}
	// ダウンロードしたファイルを保存
	downloadPath := "file/downloaded_file.zip" // 保存するファイル名
	err = download.SaveAs(downloadPath)
	if err != nil {
		log.Printf("ダウンロードファイルの保存に失敗しました: %v", err)
		return err
	} else {
		log.Printf("ダウンロードファイルを '%s' に保存しました。\n", downloadPath)
	}

	// スクリーンショットを撮って保存 (デバッグや証拠として便利)

	log.Println("スクレイピングが完了しました。")

	// ここからPlaywrightのコードを記述できます
	// 例: ブラウザを起動してGoogleにアクセス
	return nil // ここではエラーがないことを示すために nil を返します
}

// clickSelector は指定されたセレクターをクリックするヘルパー関数です
// エラーが発生した場合はログに出力し、エラーを返します
// 成功した場合はクリックしたセレクターをログに出力します
func clickSelector(page playwright.Page, selector string, timeout ...int32) error {
	var opts playwright.LocatorClickOptions
	if len(timeout) > 0 {
		opts.Timeout = playwright.Float(float64(timeout[0]))
	}
	err := page.Locator(selector).Click(opts)
	if err != nil {
		log.Printf("セレクター '%s' のクリックに失敗しました: %v", selector, err)
		return err
	}
	log.Printf("セレクター '%s' をクリックしました。\n", selector)
	return nil
}

// waitforSelector は指定されたセレクターが表示されるまで待機するヘルパー関数です
// エラーが発生した場合はログに出力し、エラーを返します
// 成功した場合は表示されたセレクターをログに出力します
func waitForSelector(page playwright.Page, selector string, timeout ...int32) error {
	var opts playwright.LocatorWaitForOptions
	if len(timeout) > 0 {
		opts.Timeout = playwright.Float(float64(timeout[0]))
	}
	err := page.Locator(selector).WaitFor(opts)
	if err != nil {
		log.Printf("セレクター '%s' の表示待機に失敗しました: %v", selector, err)
		return err
	}
	log.Printf("セレクター '%s' が表示されました。\n", selector)
	return nil
}

// inputSelector は指定されたセレクターにテキストを入力するヘルパー関数です
// エラーが発生した場合はログに出力し、エラーを返します
// 成功した場合は入力したセレクターをログに出力します
func inputSelector(page playwright.Page, selector string, text string) error {
	err := page.Locator(selector).Fill(text)
	if err != nil {
		log.Printf("セレクター '%s' への入力に失敗しました: %v", selector, err)
		return err
	}
	log.Printf("セレクター '%s' にテキスト '%s' を入力しました。\n", selector, text)
	return nil
}

// スクリーンショットを撮る関数
func takeScreenshot(page playwright.Page, path string) error {
	_, err := page.Screenshot(playwright.PageScreenshotOptions{Path: playwright.String("./file/" + path)})
	if err != nil {
		log.Printf("スクリーンショットの撮影に失敗しました: %v", err)
		return err
	}
	log.Printf("スクリーンショットを '%s' に保存しました。\n", path)
	return nil
}

// selectorExists は指定されたセレクターが存在するかどうかを確認するヘルパー関数です
// 存在する場合は true、存在しない場合は false を返します
func selectorExists(page playwright.Page, selector string) (bool, error) {
	count, err := page.Locator(selector).Count()
	if err != nil {
		log.Printf("セレクター '%s' の存在確認に失敗しました: %v", selector, err)
		return false, err
	}
	exists := count > 0
	log.Printf("セレクター '%s' の存在確認: %v\n", selector, exists)
	return exists, nil
}

// inArray は target が arr に含まれているかどうかを判定します
func inArray(target string, arr []string) bool {
	for _, v := range arr {
		if v == target {
			return true
		}
	}
	return false
}
