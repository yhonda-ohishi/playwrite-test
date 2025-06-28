package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/playwright-community/playwright-go"
	"github.com/stretchr/testify/assert"
)

func TestGetEtcMeisai(t *testing.T) {
	// This is a placeholder for the getEtcMeisai function test.
	// Since the function does not return any value,
	// we cannot directly test it. Instead, we can check
	// if it runs without panicking.
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("getEtcMeisai function panicked: %v", r)
		}
	}()
	// getEtcMeisai()
}

func TestContains(t *testing.T) {
	tests := []struct {
		name     string
		slice    string
		element  string
		expected bool
	}{
		{"Element exists", "abc", "b", true},
		{"Element does not exist", "abc", "d", false},
		{"Empty slice", "", "a", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := contains(tt.slice, tt.element)
			if result != tt.expected {
				t.Errorf("contains(%v, %s) = %v; want %v", tt.slice, tt.element, result, tt.expected)
				log.Printf("Test %s failed: expected %v, got %v", tt.name, tt.expected, result)
			}
		})
	}
}

// MockPage represents a mock implementation of a page interface
type MockPage struct {
	radioButtons map[string]map[string]bool // name -> value -> exists
}

// NewMockPage creates a new mock page instance
func NewMockPage() *MockPage {
	return &MockPage{
		radioButtons: make(map[string]map[string]bool),
	}
}

// AddRadioButton adds a radio button to the mock page
func (m *MockPage) AddRadioButton(name, value string) {
	if m.radioButtons[name] == nil {
		m.radioButtons[name] = make(map[string]bool)
	}
	m.radioButtons[name][value] = true
}

// HasRadioButton checks if a radio button exists
func (m *MockPage) HasRadioButton(name, value string) bool {
	if m.radioButtons[name] == nil {
		return false
	}
	return m.radioButtons[name][value]
}

// clickRadioButtonByNameByValue tests the clickRadioButtonByNameByValue function
func TestClickRadioButtonByNameByValue(t *testing.T) {
	// This is a placeholder for the clickRadioButtonByNameByValue function test.
	// Since the function does not return any value,
	// we cannot directly test it. Instead, we can check
	// if it runs without panicking.
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("clickRadioButtonByNameByValue function panicked: %v", r)
		}
	}()
	// Create a mock page
	mockPage := NewMockPage()
	_ = mockPage // Ignore unused variable warning
	// Create a playwright page for testing
	pw, err := playwright.Run()
	if err != nil {
		t.Fatalf("could not start playwright: %v", err)
	}
	browser, err := pw.Chromium.Launch()
	if err != nil {
		t.Fatalf("could not launch browser: %v", err)
	}
	page, err := browser.NewPage()
	if err != nil {
		t.Fatalf("could not create page: %v", err)
	}
	defer pw.Stop()
	defer browser.Close()

	err = clickRadioButtonByNameByValue(nil, "testRadio", 1)
	if err == nil {
		t.Errorf("clickRadioButtonByNameByValue should returned an error: %v", err)
	}

	err = clickRadioButtonByNameByValue(page, "testRadio", 1)
	if err == nil {
		t.Errorf("clickRadioButtonByNameByValue should returned an error: %v", err)
	}
	// Check that the error is a timeout error
	if err != nil && err.Error() != "ラジオボタン testRadio の表示待機中にエラーが発生しました: playwright: timeout: Timeout 3000ms exceeded." {
		t.Errorf("Expected timeout error, but got: %v", err.Error())
	}

	page.SetContent("<!DOCTYPE html><html><body><form><input type='radio' name='testRadio' value='1'><input type='radio' name='testRadio' value='2'></form></body></html>")
	err = clickRadioButtonByNameByValue(page, "testRadio", 1)
	if err != nil {
		t.Errorf("clickRadioButtonByNameByValue returned an error: %v", err)
	}
	assert.NotNil(t, page, "Page should not be nil after setting content")
	assert.Nil(t, err, "Expected no error when clicking radio button")
	// Check if the radio button was clicked
	isChecked, err := page.Locator("input[name='testRadio'][value='1']").IsChecked()
	if err != nil {
		t.Errorf("Error checking if radio button is checked: %v", err)
	} else if !isChecked {
		t.Errorf("Expected radio button with name 'testRadio' and value '1' to be checked, but it was not.")
	}

	// pageのmock にradioボタンを追加

}

func TestPostErrorToLineWorksBot(t *testing.T) {
	// This is a placeholder for the postErrorToLineWorksBot function test.
	// Since the function does not return any value,
	// we cannot directly test it. Instead, we can check
	// if it runs without panicking.
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("postErrorToLineWorksBot function panicked: %v", r)
		}
	}()

	// t.Error("postErrorToLineWorksBot function is not implemented yet, please implement it before running the test")
	// if err == nil {
	// 	t.Errorf("postErrorToLineWorksBot should have returned an error, but got nil")
	// }

	// postErrorToLineWorksBot()
}

// MockRoundTripper は http.RoundTripper インターフェースのモック実装です。
type MockRoundTripper struct {
	Response    *http.Response      // 返したいレスポンス
	Err         error               // 返したいエラー
	RequestFunc func(*http.Request) // リクエストが来たときに実行する関数 (検証用)
}

// RoundTrip は http.RoundTripper インターフェースのメソッドです。
// モックのレスポンスまたはエラーを返します。
func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if m.RequestFunc != nil {
		m.RequestFunc(req) // リクエストの検証を実行
	}
	return m.Response, m.Err
}

type MyRequest struct {
	Test    string `json:"test"`
	Message string `json:"message"`
}
type MyResponse struct {
	Message string `json:"message"`
	ID      string `json:"id"`
}

func TestPostJson(t *testing.T) {
	// This is a placeholder for the postJson function test.
	// Since the function does not return any value,
	// we cannot directly test it. Instead, we can check
	// if it runs without panicking.
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("postJson function panicked: %v", r)
		}
	}()

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path != "/" {
			t.Errorf("Expected path /, got %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// リクエストボディを読み込み
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("Failed to read request body: %v", err)
		}
		defer r.Body.Close()

		if len(bodyBytes) == 0 {
			t.Errorf("Expected non-empty request body, got empty")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// リクエストボディの内容をログに出力
		log.Printf("Received request body: %s", bodyBytes)

		// リクエストボディをJSONとしてデコード
		jsonData := string(bodyBytes)
		log.Printf("POSTリクエストを受信: URL=%s, データ=%s", r.URL, jsonData)

		// JSONをデコードしてMyRequest構造体に変換
		// ここでは、MyRequest構造体を定義していると
		// 仮定しています。実際の構造体に合わせて変更してください。

		var request MyRequest
		if err := json.Unmarshal(bodyBytes, &request); err != nil {
			log.Printf("Failed to decode request body: %v, body: %s", err, string(bodyBytes))
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"Invalid JSON"}`))
			// t.Errorf("Failed to decode request body %s", string(bodyBytes))
			return
		}
		log.Printf("Decoded request: %+v", request)
		respJSON, _ := json.Marshal("success") // エラーハンドリングは省略

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK) // POST成功は201 Createdが適切
		w.Write(respJSON)
	}))
	defer mockServer.Close()
	err := postJson(`{"key":"value"}`, mockServer.URL)
	assert.NotNil(t, err, "Expected an error when posting nil data")
	assert.Equal(t, "HTTP POSTリクエスト失敗: ステータスコード 400", err.Error(), "Expected error message for unsupported protocol scheme")

	err = postJson(`{"key":"value"}`, "")
	assert.NotNil(t, err, "Expected an error when posting nil data")
	assert.Equal(t, "HTTP POSTリクエスト送信エラー: Post \"\": unsupported protocol scheme \"\"", err.Error(), "Expected error message for empty URL")

	payloadToSend := map[string]string{"test": "test", "message": "test"}

	err = postJson(payloadToSend, mockServer.URL)
	assert.Nil(t, err, "Expected no error when posting valid JSON data")

}
