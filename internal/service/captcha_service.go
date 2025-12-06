package service

import (
	"fmt"
	"image"
	"log"
	"time"

	"github.com/golang/freetype/truetype"
	"github.com/google/uuid"
	"github.com/wenlng/go-captcha-assets/resources/fonts/fzshengsksjw"
	"github.com/wenlng/go-captcha-assets/resources/imagesv2"
	"github.com/wenlng/go-captcha/v2/base/option"
	"github.com/wenlng/go-captcha/v2/click"

	"github.com/qs-lzh/movie-reservation/internal/cache"
)

type CaptchaService interface {
	Generate() (mBase64, tBase64, cacheKey string, err error)
	Verify(clickData []Dot, dotAnswerData map[int]*click.Dot) bool
	VerifyWithKey(clickData []Dot, cacheKey string) (bool, error)
}

type captchaService struct {
	Cache *cache.RedisCache
}

func NewCaptchaService(cache *cache.RedisCache) *captchaService {
	return &captchaService{
		Cache: cache,
	}
}

var _ CaptchaService = (*captchaService)(nil)

var textCapt click.Captcha

func (s *captchaService) Generate() (mBase64, tBase64, cacheKey string, err error) {
	builder := click.NewBuilder(
		click.WithRangeLen(option.RangeVal{Min: 4, Max: 6}),
		click.WithRangeVerifyLen(option.RangeVal{Min: 2, Max: 4}),
	)

	// You can use preset material resources：https://github.com/wenlng/go-captcha-assets
	fontN, err := fzshengsksjw.GetFont()
	if err != nil {
		return "", "", "", fmt.Errorf("failed to load font: %w", err)
	}
	bgImage, err := imagesv2.GetImages()
	if err != nil {
		return "", "", "", fmt.Errorf("failed to load background images: %w", err)
	}

	builder.SetResources(
		click.WithChars([]string{"1A", "5E", "3d", "0p", "78", "DL", "CB", "9M"}),
		click.WithFonts([]*truetype.Font{
			fontN,
		}),
		click.WithBackgrounds([]image.Image{
			bgImage[0],
		}),
	)

	textCapt = builder.Make()

	captData, err := textCapt.Generate()
	if err != nil {
		return "", "", "", fmt.Errorf("failed to generate captcha: %w", err)
	}

	// Extract the answer data to store in cache (this should be serializable)
	dotAnswerData := captData.GetData()
	captchaID := uuid.New().String()

	err = s.Cache.Set(captchaID, dotAnswerData, 5*time.Minute)
	if err != nil {
		log.Printf("Warning: failed to save captcha answer to redis: %v", err)
	}

	mBase64, err = captData.GetMasterImage().ToBase64()
	if err != nil {
		fmt.Println(err)
	}
	tBase64, err = captData.GetThumbImage().ToBase64()
	if err != nil {
		fmt.Println(err)
	}
	return mBase64, tBase64, captchaID, nil
}

// for parse data from frontend
type Dot struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func (s *captchaService) Verify(clickData []Dot, dotAnswerData map[int]*click.Dot) bool {

	fmt.Println("start running captchaService.Verify......")
	fmt.Println("0")
	fmt.Printf("len of user: %d\n", len(clickData))
	fmt.Printf("len of answer: %d\n", len(dotAnswerData))
	if len(clickData) != len(dotAnswerData) {
		return false
	}
	fmt.Println("1")

	// the key of dotAnswerData begin with 0
	chkRet := false
	for idx, dot := range clickData {
		fmt.Printf("cycle: %d\n", idx)
		fmt.Printf("userDot: {%d, %d}\n", dot.X, dot.Y)
		answerDot := dotAnswerData[idx]
		fmt.Printf("answerDot: {%d, %d}\n", answerDot.X, answerDot.Y)
		fmt.Printf("answerRec: {%d, %d}\n", answerDot.Width, answerDot.Height)
		chkRet = click.Validate(dot.X, dot.Y, answerDot.X, answerDot.Y, answerDot.Width, answerDot.Height, 5)
		if !chkRet {
			fmt.Println(chkRet)
			return false
		}
	}
	return true
}

func (s *captchaService) VerifyWithKey(clickData []Dot, cacheKey string) (bool, error) {
	dotAnswerData := make(map[int]*click.Dot)
	if err := s.Cache.Get(cacheKey, &dotAnswerData); err != nil {
		return false, fmt.Errorf("failed to get captcha answer data from cache: %v", err)
	}

	valid := s.Verify(clickData, dotAnswerData)

	return valid, nil
}
