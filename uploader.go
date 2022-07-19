package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Println("usage error.")
		fmt.Println("mattermost_custom_emoji_uploader http://your-mattermost.example [personal access token] [path of dir or image file]")
		os.Exit(1)
	}
	baseURL := os.Args[1]
	token := os.Args[2]
	imagePath := os.Args[3]

	u, err := url.Parse(baseURL)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	apiUrl := *u
	apiUrl.Path = path.Join(apiUrl.Path, "api", "v4")

	c := model.Client4{
		Url:        u.String(),
		ApiUrl:     apiUrl.String(),
		HttpClient: &http.Client{},
		AuthToken:  token,
		AuthType:   model.HEADER_BEARER,
		HttpHeader: nil,
	}
	// me で自分自身のデータを取得することができる
	user, resp := c.GetUser("me", "")
	if resp.Error != nil {
		fmt.Println(resp.Error.Error())
		os.Exit(1)
	}

	info, err := os.Stat(imagePath)
	if err != nil {
		fmt.Println(resp.Error.Error())
		os.Exit(1)
	}
	var paths []string
	if info.IsDir() {
		pngFiles, err := filepath.Glob(filepath.Join(imagePath, "*.png"))
		if err != nil {
			fmt.Println(resp.Error.Error())
			os.Exit(1)
		}
		paths = append(paths, pngFiles...)

		jpegFiles, err := filepath.Glob(filepath.Join(imagePath, "*.jpg"))
		if err != nil {
			fmt.Println(resp.Error.Error())
			os.Exit(1)
		}
		paths = append(paths, jpegFiles...)

		jpegFiles, err = filepath.Glob(filepath.Join(imagePath, "*.jpeg"))
		if err != nil {
			fmt.Println(resp.Error.Error())
			os.Exit(1)
		}
		paths = append(paths, jpegFiles...)

		gifFiles, err := filepath.Glob(filepath.Join(imagePath, "*.gif"))
		if err != nil {
			fmt.Println(resp.Error.Error())
			os.Exit(1)
		}
		paths = append(paths, gifFiles...)
	} else {
		paths = append(paths, imagePath)
	}

	for _, imagePath := range paths {
		fmt.Printf("filepath: %v のファイルを処理中...\n", imagePath)
		err = createEmoji(c, imagePath, user.Id)
		if err != nil {
			fmt.Println("エラー発生。")
			fmt.Println(err.Error())
		} else {
			fmt.Println("処理完了")
		}
	}

}

func createEmoji(c model.Client4, imagePath string, creatorId string) error {
	f, err := os.Open(imagePath)
	defer closeFile(f)
	if err != nil {
		return err
	}

	decode, extension, err := image.Decode(f)
	if err != nil {
		return err
	}
	s, err := f.Stat()
	if err != nil {
		return err
	}
	name := strings.TrimSuffix(s.Name(), path.Ext(imagePath))
	emoji := &model.Emoji{Name: name, CreatorId: creatorId}

	buf := &bytes.Buffer{}

	if extension == "png" {
		err = png.Encode(buf, decode)
	} else if extension == "jpg" || extension == "jpeg" {
		err = jpeg.Encode(buf, decode, nil)
	} else if extension == "gif" {
		err = gif.Encode(buf, decode, nil)
	} else {
		return errors.New("読み込みに失敗しました。拡張子は jpg, jpeg, gif しか対応していません。")
	}
	if err != nil {
		return err
	}

	_, resp := c.CreateEmoji(emoji, buf.Bytes(), s.Name())
	if resp.Error != nil {
		err = errors.New(resp.Error.Error())
		return err
	} else {
		return nil
	}
}

func closeFile(f *os.File) {
	err := f.Close()
	// closeに失敗した場合は潔くコマンドを終了する
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
