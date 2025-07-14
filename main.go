package main

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/harry/rename-by-tmdb/internal/models"
	"github.com/harry/rename-by-tmdb/internal/services"
	"github.com/harry/rename-by-tmdb/internal/utils"
)

// 处理电影重命名
func handleMovie(tmdbService *services.TMDBService) error {
	// 获取用户当前文件名中的标题部分
	fileTitle, err := utils.GetUserInput("请输入当前文件名中的标题部分（例如：The.Matrix）: ")
	if err != nil {
		return fmt.Errorf("错误: %v", err)
	}

	// 获取电影ID
	movieID, err := utils.GetUserInput("请输入电影ID: ")
	if err != nil {
		return fmt.Errorf("错误: %v", err)
	}

	// 获取电影信息
	movie, err := tmdbService.FetchMovieInfo(movieID)
	if err != nil {
		return fmt.Errorf("获取电影信息失败: %v", err)
	}

	// 从发布日期中提取年份
	year := ""
	if len(movie.ReleaseDate) >= 4 {
		year = movie.ReleaseDate[:4]
	}

	// 将空格替换为点号
	movieName := strings.ReplaceAll(movie.Title, " ", ".")

	// 创建命名格式
	namingFormat := fmt.Sprintf("%s.%s.{[tmdbid=%s;type=movie]}",
		movieName, year, movieID)
	fmt.Printf("命名格式：\n%s\n", namingFormat)

	var wordGroup *models.WordGroup
	if utils.IsUploadEnabled() {
		// 创建词组服务
		wordGroupService, err := services.NewWordGroupService()
		if err != nil {
			return fmt.Errorf("创建词组服务失败: %v", err)
		}

		// 创建词组
		wordGroup, err = wordGroupService.CreateWordGroup(namingFormat)
		if err != nil {
			return fmt.Errorf("创建词组失败: %v", err)
		}
		fmt.Printf("词组创建成功，ID: %d\n", wordGroup.ID)

		// 构建电影的替换规则
		beReplaced := fmt.Sprintf("%s.*", regexp.QuoteMeta(fileTitle))
		replace := fmt.Sprintf("%s.%s.{[tmdbid=%s;type=movie]}", movieName, year, movieID)

		fmt.Printf("\n被替换词：\n%s\n", beReplaced)
		fmt.Printf("替换词：\n%s\n", replace)

		// 上传替换规则
		err = wordGroupService.AddWordUnit(wordGroup.ID, beReplaced, replace, "", "", 0)
		if err != nil {
			return fmt.Errorf("上传电影替换规则失败: %v", err)
		}
		fmt.Printf("电影替换规则上传成功\n")

		fmt.Println("\n注意：")
		fmt.Println("1. 正则表达式中的点号（.）已经被转义")
		fmt.Println("2. 替换词中的'\\1'表示保留原始集数")
		fmt.Println("3. [^.]* 匹配除点号外的任意字符，用于处理标题和集数之间可能存在的额外字符")
		fmt.Println("4. 替换后的文件名使用TMDB中的官方电影名称")
	}

	return nil
}

// 处理剧集重命名
func handleTVShow(tmdbService *services.TMDBService) error {
	// 获取剧集ID
	seriesID, err := utils.GetUserInput("请输入剧集ID: ")
	if err != nil {
		return fmt.Errorf("错误: %v", err)
	}

	// 获取剧集信息
	show, err := tmdbService.FetchShowInfo(seriesID)
	if err != nil {
		return fmt.Errorf("获取剧集信息失败: %v", err)
	}

	// 从首播日期中提取年份
	year := ""
	if len(show.FirstAirDate) >= 4 {
		year = show.FirstAirDate[:4]
	}

	// 显示剧集命名格式
	showType := "tv"
	if show.Type == "movie" {
		showType = "movie"
	}

	// 将空格替换为点号
	showName := strings.ReplaceAll(show.Name, " ", ".")

	// 获取最后一季的最大集数
	var maxEpisodeNumber int
	if len(show.Seasons) > 0 {
		// 找到最后一个非第0季
		var lastSeason *models.TMDBSeason
		for i := len(show.Seasons) - 1; i >= 0; i-- {
			if show.Seasons[i].SeasonNumber > 0 {
				lastSeasonDetails, err := tmdbService.FetchSeasonDetails(seriesID, show.Seasons[i].SeasonNumber)
				if err == nil && len(lastSeasonDetails.Episodes) > 0 {
					lastSeason = lastSeasonDetails
					break
				}
			}
		}
		if lastSeason != nil && len(lastSeason.Episodes) > 0 {
			maxEpisodeNumber = lastSeason.Episodes[len(lastSeason.Episodes)-1].EpisodeNumber
		}
	}

	// 确保至少使用2位数
	digits := len(strconv.Itoa(maxEpisodeNumber))
	if digits < 2 {
		digits = 2
	}

	// 创建命名格式
	namingFormat := fmt.Sprintf("%s.%s.{[tmdbid=%s;type=%s]}",
		showName, year, seriesID, showType)
	fmt.Printf("命名格式：\n%s\n", namingFormat)

	var wordGroup *models.WordGroup
	var wordGroupService *services.WordGroupService
	if utils.IsUploadEnabled() {
		// 创建词组服务
		var err error
		wordGroupService, err = services.NewWordGroupService()
		if err != nil {
			return fmt.Errorf("创建词组服务失败: %v", err)
		}

		// 创建词组
		wordGroup, err = wordGroupService.CreateWordGroup(namingFormat)
		if err != nil {
			return fmt.Errorf("创建词组失败: %v", err)
		}
		fmt.Printf("词组创建成功，ID: %d\n", wordGroup.ID)
	}

	// 获取用户当前文件名中的标题部分
	fileTitle, err := utils.GetUserInput("请输入当前文件名中的标题部分（例如：One.Piece）: ")
	if err != nil {
		return fmt.Errorf("错误: %v", err)
	}

	// 获取是否包含季数的选择
	hasSeason, err := utils.GetHasSeasonChoice()
	if err != nil {
		return fmt.Errorf("错误: %v", err)
	}

	var specificSeasons []int
	var generateAllSeasons bool

	if !hasSeason {
		// 如果不使用原文件名季数，则询问用户要生成哪些季
		specificSeasons, generateAllSeasons, err = utils.GetSpecificSeasons()
		if err != nil {
			return fmt.Errorf("错误: %v", err)
		}
	}

	// 获取集数偏移量
	episodeOffset, err := utils.GetEpisodeOffset()
	if err != nil {
		return fmt.Errorf("错误: %v", err)
	}

	// 设置后定位词为".年份."
	backPositionWord := fmt.Sprintf(".%s.", year)

	fmt.Printf("\n=== %s 各季重命名正则表达式 ===\n", show.Name)

	// 为每一季生成替换规则
	for _, season := range show.Seasons {
		// 如果不使用原文件名季数且不是生成所有季，则检查是否是用户指定的季
		if !hasSeason && !generateAllSeasons {
			isSpecificSeason := false
			for _, s := range specificSeasons {
				if s == season.SeasonNumber {
					isSpecificSeason = true
					break
				}
			}
			if !isSpecificSeason {
				continue
			}
		}

		// 显示季数信息（为第0季添加特别说明）
		if season.SeasonNumber == 0 {
			fmt.Printf("\n--- 特别篇 ---\n")
		} else {
			fmt.Printf("\n--- 第 %d 季 ---\n", season.SeasonNumber)
		}

		// 获取该季的详细信息
		seasonDetails, err := tmdbService.FetchSeasonDetails(seriesID, season.SeasonNumber)
		if err != nil {
			fmt.Printf("获取第 %d 季信息失败: %v\n", season.SeasonNumber, err)
			continue
		}

		if len(seasonDetails.Episodes) == 0 {
			fmt.Printf("第 %d 季没有找到任何剧集\n", season.SeasonNumber)
			continue
		}

		// 计算这一季的起始和结束集数
		startEp := seasonDetails.Episodes[0].EpisodeNumber
		endEp := seasonDetails.Episodes[len(seasonDetails.Episodes)-1].EpisodeNumber

		// 计算原文件中的集数范围（减去偏移量，因为原文件需要减去这个值）
		sourceStartEp := startEp - episodeOffset
		sourceEndEp := endEp - episodeOffset

		// 确保源文件的集数不会变成负数
		if sourceStartEp <= 0 || sourceEndEp <= 0 {
			fmt.Printf("警告：偏移量 %d 会导致源文件集数小于等于0，跳过此季\n", episodeOffset)
			continue
		}

		// 构建匹配范围的正则表达式
		var beReplaced string
		if hasSeason {
			beReplaced = fmt.Sprintf("%s.*S%02dE(%s)",
				regexp.QuoteMeta(fileTitle),
				season.SeasonNumber,
				utils.GenerateRangePattern(sourceStartEp, sourceEndEp, digits))
		} else {
			beReplaced = fmt.Sprintf("%s.*(?:S\\d{2})?E(%s)",
				regexp.QuoteMeta(fileTitle),
				utils.GenerateRangePattern(sourceStartEp, sourceEndEp, digits))
		}

		// 生成替换词和前后定位词
		replace := fmt.Sprintf("%s.S%02dE\\1.%s.{[tmdbid=%s;type=%s]}", showName, season.SeasonNumber, year, seriesID, showType)

		// 只在有偏移量时设置前后定位词
		var prefix, suffix string
		if episodeOffset != 0 {
			prefix = fmt.Sprintf("%s.S%02dE", showName, season.SeasonNumber)
			suffix = backPositionWord
		}

		// 显示集数范围和对应关系
		fmt.Printf("集数范围：%d-%d（使用%d位数）\n", sourceStartEp, sourceEndEp, digits)
		if episodeOffset != 0 {
			fmt.Printf("集数偏移量：%+d\n", episodeOffset)
			fmt.Printf("原始集数示例：%0*d → 实际集数：%0*d\n",
				digits, sourceStartEp, digits, startEp)
		}

		fmt.Printf("\n被替换词：\n%s\n", beReplaced)
		fmt.Printf("替换词：\n%s\n", replace)

		// 只在有偏移量时显示前后定位词
		if episodeOffset != 0 {
			fmt.Printf("\n前定位词：\n%s\n", prefix)
			fmt.Printf("后定位词：\n%s\n", suffix)
		}

		// 只在启用上传时上传替换规则
		if utils.IsUploadEnabled() {
			err = wordGroupService.AddWordUnit(wordGroup.ID, beReplaced, replace, prefix, suffix, episodeOffset)
			if err != nil {
				return fmt.Errorf("上传第 %d 季替换规则失败: %v", season.SeasonNumber, err)
			}
			fmt.Printf("第 %d 季替换规则上传成功\n", season.SeasonNumber)
		}
	}

	fmt.Println("\n注意：")
	fmt.Println("1. 正则表达式中的点号（.）已经被转义")
	fmt.Println("2. 替换词中的'\\1'表示保留原始集数")
	fmt.Println("3. [^.]* 匹配除点号外的任意字符，用于处理标题和集数之间可能存在的额外字符")
	fmt.Println("4. 替换后的文件名使用TMDB中的官方剧集名称")
	fmt.Printf("5. 所有集数都使用相同的位数（由最大集数决定），不足位数补0\n")
	fmt.Printf("   例如：如果最大集数是500（3位），则第1集应该写作001\n")
	if !hasSeason {
		fmt.Println("6. 原文件名不包含季数，仅匹配集数部分")
	}
	if episodeOffset != 0 {
		fmt.Printf("7. 被替换词中的集数范围已经过调整，可以直接匹配原文件名中的集数\n")
	}

	return nil
}

func main() {
	// 加载环境变量
	if err := utils.LoadEnv(); err != nil {
		fmt.Printf("错误: %v\n", err)
		fmt.Println("\n请确保 .env 文件存在并包含必要的环境变量")
		return
	}

	// 检查环境变量
	if err := utils.CheckRequiredEnvVars(); err != nil {
		fmt.Printf("错误: %v\n", err)
		fmt.Println("\n请在 .env 文件中设置以下环境变量：")
		fmt.Println("TMDB_API_KEY='your_tmdb_api_key'")
		return
	}

	// 创建TMDB服务
	tmdbService, err := services.NewTMDBService()
	if err != nil {
		fmt.Printf("创建TMDB服务失败: %v\n", err)
		return
	}

	// 获取媒体类型选择
	fmt.Println("请选择媒体类型：")
	fmt.Println("1. 电影")
	fmt.Println("2. 剧集")
	mediaType, err := utils.GetUserInput("请输入选项（1或2）: ")
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}

	var handleErr error
	switch mediaType {
	case "1":
		handleErr = handleMovie(tmdbService)
	case "2":
		handleErr = handleTVShow(tmdbService)
	default:
		fmt.Println("无效的选项，请输入1或2")
		return
	}

	if handleErr != nil {
		fmt.Printf("%v\n", handleErr)
		return
	}

	// 等待用户按回车键退出
	fmt.Print("\n按回车键退出...")
	if _, err := fmt.Scanln(); err != nil && err != io.EOF {
		fmt.Printf("读取输入失败: %v\n", err)
	}
}
