package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/c-rainbow/simplechatbot/install"
)

var languageCodeFlag = flag.String("language", "", "Language to run installer.")

var allLanguageCodes = []string{"en", "ko"}

// Installer entry point main code
func main() {
	languageCode := getLanguageCode()
	if languageCode != "" {
		installer := getInstaller(languageCode)
		installer.Install()
	}
}

func getLanguageCode() string {
	languageCode := sanitizeLanguageCode(*languageCodeFlag)

	scanner := bufio.NewScanner(os.Stdin)
	for !isValidLanguageCode(languageCode) {
		fmt.Println("For English, type 'en'. 한국어 설치는 'ko'를 입력하세요.")

		for scanner.Scan() {
			languageCode = sanitizeLanguageCode(scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("Quitting installer because of error. 에러가 발생하여 인스톨을 종료합니다.")
			return ""
		}
	}
	return languageCode
}

// This function assumes that the language code is valid
func getInstaller(languageCode string) *install.Installer {
	switch languageCode {
	case "en":
		return install.NewInstallerEng()
	case "ko":
		return install.NewInstallerKor()
	}
	return nil
}

func isValidLanguageCode(languageCode string) bool {
	languageCode = strings.ToLower(strings.TrimSpace(languageCode))
	for _, code := range allLanguageCodes {
		if languageCode == code {
			return true
		}
	}
	return false
}

func sanitizeLanguageCode(code string) string {
	return strings.ToLower(strings.TrimSpace(code))
}
