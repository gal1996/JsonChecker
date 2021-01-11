package lib

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/thoas/go-funk"
	"io"
	"os"
	"strings"
)

func JsonCheck(target, correctFilePath string) error {
	f, err := os.Open(correctFilePath)
	if err != nil {
		return err
	}

	// ファイルを読み込みインデントをつける
	indentedCorrectJson, err := addIndentIntoRowJson(f)
	if err != nil {
		return err
	}

	// ファイルから期待するjsonのkeyの抽出
	expectedKeys := extractJsonKeys(indentedCorrectJson)
	fmt.Printf("[Info] expectedKeys(len: %d): %#v\n", len(expectedKeys), expectedKeys)

	// 入力されたjsonを読み込みインデントをつける
	indentedTargetJson, err := addIndentIntoRowJson(strings.NewReader(target))
	if err != nil {
		return err
	}
	// リクエストからきたjsonのkeyの抽出
	actualKeys := extractJsonKeys(indentedTargetJson)
	fmt.Printf("[Info] acctualKeys(len: %d): %#v\n", len(actualKeys), actualKeys)

	// そもそも二つの長さが違うのはアウト
	if len(expectedKeys) != len(actualKeys) {
		return fmt.Errorf("length is not equal")
	}

	invalidKey := make([]string, 0, len(actualKeys))
	// 同じkeyで複数データが存在する場合、長さが同じでも存在しないkeyがすり抜けてしまうので、actualとexpected両方でチェックする必要がある
	for _, actualKey := range actualKeys {
		if !funk.ContainsString(expectedKeys, actualKey) {
			actualKey = strings.TrimSpace(actualKey)
			data := fmt.Sprintf("%s(is invalid)", actualKey)
			invalidKey = append(invalidKey, data)
		}
	}

	for _, expectedKey := range expectedKeys {
		if !funk.ContainsString(actualKeys, expectedKey) {
			expectedKey = strings.TrimSpace(expectedKey)
			data := fmt.Sprintf("%s(is not in your correctJsonBytes)", expectedKey)
			invalidKey = append(invalidKey, data)
		}
	}

	if len(invalidKey) != 0 {
		resStr := fmt.Sprintf("your correctJsonBytes is not perfect invalid keys: %s", invalidKey)
		return fmt.Errorf(resStr)
	}

	fmt.Printf("your json is correct!!")

	return nil
}

func addIndentIntoRowJson(f io.Reader) (io.Reader, error) {
	buf := make([]byte, 2048)
	var JsonBytes []byte
	for {
		n, err := f.Read(buf)
		// バイト数が0は読み取り終了
		if n == 0 {
			break
		}
		if err != nil {
			break
		}

		JsonBytes = buf[:n]
	}

	var indentedJson bytes.Buffer
	err := json.Indent(&indentedJson, JsonBytes, "", "  ")
	if err != nil {
		fmt.Printf("error has occurred at indenting json: %s\n", err.Error())
		return nil, fmt.Errorf("%s\n", err.Error())
	}
	return &indentedJson, nil
}

func extractJsonKeys(json io.Reader) []string {
	jsonKeys := make([]string, 0, 30)
	scanner := bufio.NewScanner(json)
	for scanner.Scan() {
		text := scanner.Text()
		if strings.Contains(text, ":") {
			// 基本的にjsonは 'key: value' の形なので、':'の前がキーになる
			jsonKey := strings.TrimSpace(strings.Replace((strings.Split(text, ":"))[0], "\"", "", -1))
			jsonKeys = append(jsonKeys, jsonKey)
		}
	}

	return jsonKeys
}
