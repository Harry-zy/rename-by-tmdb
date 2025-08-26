package main

import (
	"fmt"
	"io"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/harry/rename-by-tmdb/internal/models"
	"github.com/harry/rename-by-tmdb/internal/services"
	"github.com/harry/rename-by-tmdb/internal/utils"
)

// romanToArabic 将罗马数字转换为阿拉伯数字
func romanToArabic(roman string) int {
	romanMap := map[byte]int{
		'i': 1, 'v': 5, 'x': 10, 'l': 50,
		'c': 100, 'd': 500, 'm': 1000,
	}

	roman = strings.ToLower(roman)
	total := 0
	prevValue := 0

	for i := len(roman) - 1; i >= 0; i-- {
		value := romanMap[roman[i]]
		if value < prevValue {
			total -= value
		} else {
			total += value
		}
		prevValue = value
	}

	return total
}

// extractPartInfo 从文件标题中提取part信息
func extractPartInfo(fileTitle string) string {
	// 将输入转为小写进行匹配
	lowerTitle := strings.ToLower(fileTitle)

	// 先检查数字格式的part
	digitPatterns := []string{
		`part\.?(\d+)`,    // part1, part.1
		`part\.?\s+(\d+)`, // part 1
	}

	for _, pattern := range digitPatterns {
		re := regexp.MustCompile(pattern)
		if matches := re.FindStringSubmatch(lowerTitle); len(matches) > 1 {
			return "part" + matches[1]
		}
	}

	// 检查罗马数字格式的part
	romanPatterns := []string{
		`part\.?([ivxlcdm]+)`,    // part.i, part.ii, part.iii等
		`part\.?\s+([ivxlcdm]+)`, // part i, part ii等
	}

	for _, pattern := range romanPatterns {
		re := regexp.MustCompile(pattern)
		if matches := re.FindStringSubmatch(lowerTitle); len(matches) > 1 {
			// 将罗马数字转换为阿拉伯数字
			arabicNum := romanToArabic(matches[1])
			return fmt.Sprintf("part%d", arabicNum)
		}
	}

	return ""
}

// findExistingWordGroup 查找已存在的词组
func findExistingWordGroup(wordGroupService *services.WordGroupService, namingFormat string) (*models.WordGroup, error) {
	list, err := wordGroupService.GetWordGroupList()
	if err != nil {
		return nil, fmt.Errorf("获取词组列表失败: %v", err)
	}

	for _, group := range list.List {
		if group.Title == namingFormat {
			return &group, nil
		}
	}

	return nil, nil
}

// 处理电影重命名
func handleMovie(tmdbService *services.TMDBService) error {
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

		// 查找是否存在相同的命名格式
		existingGroup, err := findExistingWordGroup(wordGroupService, namingFormat)
		if err != nil {
			return err
		}

		if existingGroup != nil {
			// 使用已存在的词组
			wordGroup = existingGroup
			fmt.Printf("使用已存在的词组，ID: %d\n", wordGroup.ID)
		} else {
			// 创建新词组
			wordGroup, err = wordGroupService.CreateWordGroup(namingFormat)
			if err != nil {
				return fmt.Errorf("创建词组失败: %v", err)
			}
			fmt.Printf("词组创建成功，ID: %d\n", wordGroup.ID)
		}

		// 获取用户当前文件名中的标题部分
		fileTitle, err := utils.GetUserInput("请输入当前文件名中的标题部分（例如：The.Matrix）: ")
		if err != nil {
			return fmt.Errorf("错误: %v", err)
		}

		// 检测并提取part信息
		partInfo := extractPartInfo(fileTitle)

		// 构建电影的替换规则
		beReplaced := fmt.Sprintf("%s.*", regexp.QuoteMeta(fileTitle))
		var replace string
		if partInfo != "" {
			replace = fmt.Sprintf("%s.%s.%s.{[tmdbid=%s;type=movie]}", movieName, year, partInfo, movieID)
		} else {
			replace = fmt.Sprintf("%s.%s.{[tmdbid=%s;type=movie]}", movieName, year, movieID)
		}

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

	// 询问是否以日期判断集数
	fmt.Print("是否以日期判断集数？(y/N，直接回车默认为N): ")
	useDateForEpisode, err := utils.GetUserInput("")
	if err != nil {
		return fmt.Errorf("错误: %v", err)
	}
	useDateForEpisode = strings.ToLower(strings.TrimSpace(useDateForEpisode))
	isDateMode := useDateForEpisode == "y" || useDateForEpisode == "yes"

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

		// 查找是否存在相同的命名格式
		existingGroup, err := findExistingWordGroup(wordGroupService, namingFormat)
		if err != nil {
			return err
		}

		if existingGroup != nil {
			// 使用已存在的词组
			wordGroup = existingGroup
			fmt.Printf("使用已存在的词组，ID: %d\n", wordGroup.ID)
		} else {
			// 创建新词组
			wordGroup, err = wordGroupService.CreateWordGroup(namingFormat)
			if err != nil {
				return fmt.Errorf("创建词组失败: %v", err)
			}
			fmt.Printf("词组创建成功，ID: %d\n", wordGroup.ID)
		}
	}

	// 获取用户当前文件名中的标题部分
	fileTitle, err := utils.GetUserInput("请输入当前文件名中的标题部分（例如：One.Piece）: ")
	if err != nil {
		return fmt.Errorf("错误: %v", err)
	}

	// 获取是否有part剧集
	var hasPartEpisodes bool
	var partEpisodeInfo map[int][]int
	var specificSeasons []int
	var generateAllSeasons bool
	var includeSpecialSeason bool

	if !isDateMode {
		hasPartEpisodes, err = utils.GetPartEpisodeChoice()
		if err != nil {
			return fmt.Errorf("错误: %v", err)
		}

		if hasPartEpisodes {
			// 先询问要生成的季数
			fmt.Println("\n由于选择了part模式，需要先确定要生成的季数")
			specificSeasons, generateAllSeasons, err = utils.GetSpecificSeasons()
			if err != nil {
				return fmt.Errorf("错误: %v", err)
			}

			// 如果选择生成所有季，询问是否包含第0季
			if generateAllSeasons {
				includeSpecialSeason, err = utils.GetIncludeSpecialSeason()
				if err != nil {
					return fmt.Errorf("错误: %v", err)
				}
			}

			// 然后询问part剧集信息
			partEpisodeInfo, err = utils.GetPartEpisodeInfo()
			if err != nil {
				return fmt.Errorf("错误: %v", err)
			}

			// 显示用户输入的part剧集信息
			fmt.Printf("\n=== Part剧集信息 ===\n")
			for episodeNum, parts := range partEpisodeInfo {
				fmt.Printf("第%d集: part%d", episodeNum, parts[0])
				for i := 1; i < len(parts); i++ {
					fmt.Printf(", part%d", parts[i])
				}
				fmt.Println()
			}
		}
	}

	// 获取是否包含季数的选择
	var hasSeason bool
	if hasPartEpisodes {
		// 如果有part剧集，自动设置为不使用原文件名季数
		hasSeason = false
		fmt.Println("\n注意：由于选择了part剧集，自动设置为不使用原文件名季数")
	} else {
		hasSeason, err = utils.GetHasSeasonChoice()
		if err != nil {
			return fmt.Errorf("错误: %v", err)
		}
	}

	// 如果没有part剧集且不使用原文件名季数，则询问用户要生成哪些季
	if !hasPartEpisodes && !hasSeason {
		// 如果不使用原文件名季数，则询问用户要生成哪些季
		specificSeasons, generateAllSeasons, err = utils.GetSpecificSeasons()
		if err != nil {
			return fmt.Errorf("错误: %v", err)
		}

		// 如果选择生成所有季，询问是否包含第0季
		if generateAllSeasons {
			includeSpecialSeason, err = utils.GetIncludeSpecialSeason()
			if err != nil {
				return fmt.Errorf("错误: %v", err)
			}
		}
	}

	// 获取集数偏移量（如果有part剧集则不询问）
	var episodeOffset int
	if !isDateMode && !hasPartEpisodes {
		episodeOffset, err = utils.GetEpisodeOffset()
		if err != nil {
			return fmt.Errorf("错误: %v", err)
		}
	}

	// 获取是否需要补0站位
	var padZero, episodeContinuous bool
	if !isDateMode {
		padZero, err = utils.GetPadZeroChoice()
		if err != nil {
			return fmt.Errorf("错误: %v", err)
		}

		// 如果需要补0，询问集数是否连续
		if padZero {
			episodeContinuous, err = utils.GetEpisodeContinuousChoice()
			if err != nil {
				return fmt.Errorf("错误: %v", err)
			}
		}
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

		// 如果是第0季（特别篇），且是生成所有季的情况，检查是否需要包含第0季
		if season.SeasonNumber == 0 && generateAllSeasons && !includeSpecialSeason {
			continue
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

		if isDateMode {
			// 日期模式：为每一集生成替换规则
			fmt.Printf("\n=== 第 %d 季 - 日期模式 ===\n", season.SeasonNumber)

			for _, episode := range seasonDetails.Episodes {
				// 只处理有播出日期的集数
				if episode.AirDate == "" {
					fmt.Printf("第%d集：未获取到播出日期，跳过\n", episode.EpisodeNumber)
					continue
				}

				// 格式化播出日期为YYYYMMDD
				airDate := strings.ReplaceAll(episode.AirDate, "-", "")

				// 构建被替换词：标题+播出日期+后面所有字符
				beReplaced := fmt.Sprintf("%s.*%s.*", regexp.QuoteMeta(fileTitle), airDate)

				// 构建替换词：剧集名称.S季数.E集数.年份.{[tmdbid=ID;type=tv]}
				replace := fmt.Sprintf("%s.S%02dE%02d.%s.{[tmdbid=%s;type=%s]}",
					showName, season.SeasonNumber, episode.EpisodeNumber, year, seriesID, showType)

				fmt.Printf("\n第%d集 (播出日期: %s):\n", episode.EpisodeNumber, episode.AirDate)
				fmt.Printf("被替换词：\n%s\n", beReplaced)
				fmt.Printf("替换词：\n%s\n", replace)

				// 上传替换规则
				if utils.IsUploadEnabled() {
					err = wordGroupService.AddWordUnit(wordGroup.ID, beReplaced, replace, "", "", 0)
					if err != nil {
						return fmt.Errorf("上传第 %d 季第 %d 集替换规则失败: %v", season.SeasonNumber, episode.EpisodeNumber, err)
					}
					fmt.Printf("第 %d 季第 %d 集替换规则上传成功\n", season.SeasonNumber, episode.EpisodeNumber)
				}
			}
			continue // 跳过原有的集数范围处理逻辑
		}

		// Part模式处理
		if hasPartEpisodes {
			fmt.Printf("\n=== 第 %d 季 - Part模式 ===\n", season.SeasonNumber)

			// 计算需要的位数
			var digits int
			if !padZero {
				digits = 1
			} else {
				if episodeContinuous {
					digits = len(strconv.Itoa(maxEpisodeNumber))
				} else {
					digits = len(strconv.Itoa(endEp))
				}
				if digits < 2 {
					digits = 2
				}
			}

			// 为每个part生成替换规则
			// 先按集数排序，确保按正确顺序处理
			var sortedEpisodes []int
			for episodeNum := range partEpisodeInfo {
				sortedEpisodes = append(sortedEpisodes, episodeNum)
			}
			sort.Ints(sortedEpisodes)

			for _, episodeNum := range sortedEpisodes {
				parts := partEpisodeInfo[episodeNum]
				// 检查这一季是否包含这个集数
				if episodeNum < startEp || episodeNum > endEp {
					continue
				}

				for _, partNum := range parts {
					// 重新设计的偏移量计算逻辑
					offset := 0

					// 计算前面所有part2的总数（不包括part1）
					offset = 0
					for checkEpisodeNum, checkParts := range partEpisodeInfo {
						if checkEpisodeNum < episodeNum {
							// 只计算part2的数量（每个集数的part数-1）
							part2Count := len(checkParts) - 1
							if part2Count > 0 {
								offset += part2Count
							}
						}
					}

					// 如果是part2，需要额外递增+1
					if partNum > 1 {
						offset += 1
					}

					// 调试信息
					fmt.Printf("调试 - 第%d集part%d: 前面part总数=%d, 最终偏移量=%d\n",
						episodeNum, partNum, offset, offset)

					// 计算实际集数（原集数 + 偏移量）
					actualEpisodeNum := episodeNum + offset

					// 构建被替换词：包含part信息，季数可有可无，part兼容大小写，集数补0
					// 计算需要的位数（基于总集数）
					digits := len(strconv.Itoa(seasonDetails.Episodes[len(seasonDetails.Episodes)-1].EpisodeNumber))
					if digits < 2 {
						digits = 2 // 确保至少使用2位数
					}

					beReplaced := fmt.Sprintf("%s.*?(?:S%02d)?(?:E|Ep|EP|[Ee]pisode|[Ee]p)?%0*d.*?(?:[Pp]art|PART|Part)%d.*",
						regexp.QuoteMeta(fileTitle), season.SeasonNumber, digits, episodeNum, partNum)

					// 构建替换词：使用原集数，而不是偏移后的集数，集数补0
					replace := fmt.Sprintf("%s.S%02dE%0*d.%s.{[tmdbid=%s;type=%s]}",
						showName, season.SeasonNumber, digits, episodeNum, year, seriesID, showType)

					fmt.Printf("\n第%d集 part%d (偏移量:+%d, 实际集数:%d):\n",
						episodeNum, partNum, offset, actualEpisodeNum)
					fmt.Printf("被替换词：\n%s\n", beReplaced)
					fmt.Printf("替换词：\n%s\n", replace)

					// 构建前定位词和后定位词（如果有偏移量）
					var prefix, suffix string
					if offset > 0 {
						prefix = fmt.Sprintf("%s.S%02dE", showName, season.SeasonNumber)
						suffix = backPositionWord
					}

					// 上传替换规则
					if utils.IsUploadEnabled() {
						err = wordGroupService.AddWordUnit(wordGroup.ID, beReplaced, replace, prefix, suffix, offset)
						if err != nil {
							return fmt.Errorf("上传第 %d 季第 %d 集part%d替换规则失败: %v",
								season.SeasonNumber, episodeNum, partNum, err)
						}
						fmt.Printf("第 %d 季第 %d 集part%d替换规则上传成功\n",
							season.SeasonNumber, episodeNum, partNum)
					}
				}
			}

			// 为非part集数生成区间替换规则
			// 使用已经排序的episodes

			// 为每个区间生成替换规则
			for i := 0; i <= len(sortedEpisodes); i++ {
				var startEp, endEp int
				var offset int

				if i == 0 {
					// 第一个区间：从第1集到第一个part集数之前
					startEp = 1
					endEp = sortedEpisodes[0] - 1
					offset = 0
				} else if i == len(sortedEpisodes) {
					// 最后一个区间：从最后一个part集数之后到季末
					startEp = sortedEpisodes[len(sortedEpisodes)-1] + 1

					// 计算最后一个part集数的最大偏移量
					lastEpisodeNum := sortedEpisodes[len(sortedEpisodes)-1]
					lastParts := partEpisodeInfo[lastEpisodeNum]
					maxOffset := 0

					// 计算前面所有part2的总数
					for checkEpisodeNum, checkParts := range partEpisodeInfo {
						if checkEpisodeNum < lastEpisodeNum {
							// 只计算part2的数量（每个集数的part数-1）
							part2Count := len(checkParts) - 1
							if part2Count > 0 {
								maxOffset += part2Count
							}
						}
					}

					// 最后一个part2的偏移量
					if len(lastParts) > 1 {
						maxOffset += 1
					}

					// 计算最后一个part集数偏移后的最大集数
					maxActualEpisode := lastEpisodeNum + maxOffset

					// 如果偏移后超过了季末，则不生成这个区间
					if maxActualEpisode >= seasonDetails.Episodes[len(seasonDetails.Episodes)-1].EpisodeNumber {
						continue
					}

					endEp = seasonDetails.Episodes[len(seasonDetails.Episodes)-1].EpisodeNumber
					offset = maxOffset
				} else {
					// 中间区间：两个part集数之间
					startEp = sortedEpisodes[i-1] + 1
					endEp = sortedEpisodes[i] - 1
					// 计算这个区间的偏移量（继承前面最近的part2的偏移量）
					// 找到前面最近的part2的偏移量
					offset = 0
					for checkEpisodeNum := sortedEpisodes[i-1]; checkEpisodeNum >= 1; checkEpisodeNum-- {
						if checkParts, exists := partEpisodeInfo[checkEpisodeNum]; exists {
							// 如果前面有part2，则计算其part2的偏移量
							if len(checkParts) > 1 {
								// 计算前面所有part2的总数
								for checkCheckEpisodeNum, checkCheckParts := range partEpisodeInfo {
									if checkCheckEpisodeNum < checkEpisodeNum {
										// 只计算part2的数量（每个集数的part数-1）
										part2Count := len(checkCheckParts) - 1
										if part2Count > 0 {
											offset += part2Count
										}
									}
								}
								// part2的偏移量
								offset += 1
								break
							}
						}
					}
				}

				// 如果区间有效，生成替换规则
				if startEp <= endEp && startEp >= startEp && endEp <= seasonDetails.Episodes[len(seasonDetails.Episodes)-1].EpisodeNumber {
					// 构建被替换词：匹配区间内的集数
					beReplaced := fmt.Sprintf("%s.*?(?:S%02d)?(?:E|Ep|EP|[Ee]pisode|[Ee]p)?(%s)(?!.*(?:[Pp]art|PART|Part))",
						regexp.QuoteMeta(fileTitle), season.SeasonNumber,
						utils.GenerateRangePattern(startEp, endEp, digits))

					// 构建替换词：使用捕获组和偏移量
					replace := fmt.Sprintf("%s.S%02dE\\1.%s.{[tmdbid=%s;type=%s]}",
						showName, season.SeasonNumber, year, seriesID, showType)

					fmt.Printf("\n区间 %d-%d 非part集数规则 (偏移量:+%d):\n", startEp, endEp, offset)
					fmt.Printf("被替换词：\n%s\n", beReplaced)
					fmt.Printf("替换词：\n%s\n", replace)
					fmt.Printf("说明：区间内集数的实际集数 = 原集数 + %d\n", offset)

					// 构建前定位词和后定位词（如果有偏移量）
					var prefix, suffix string
					if offset > 0 {
						prefix = fmt.Sprintf("%s.S%02dE", showName, season.SeasonNumber)
						suffix = backPositionWord
					}

					// 上传替换规则
					if utils.IsUploadEnabled() {
						err = wordGroupService.AddWordUnit(wordGroup.ID, beReplaced, replace, prefix, suffix, offset)
						if err != nil {
							return fmt.Errorf("上传第 %d 季区间 %d-%d 非part集数替换规则失败: %v",
								season.SeasonNumber, startEp, endEp, err)
						}
						fmt.Printf("第 %d 季区间 %d-%d 非part集数替换规则上传成功\n",
							season.SeasonNumber, startEp, endEp)
					}
				}
			}

			continue // 跳过原有的集数范围处理逻辑
		}

		// 计算原文件中的集数范围（减去偏移量，因为原文件需要减去这个值）
		sourceStartEp := startEp - episodeOffset
		sourceEndEp := endEp - episodeOffset

		// 确保源文件的集数不会变成负数
		if sourceStartEp <= 0 || sourceEndEp <= 0 {
			fmt.Printf("警告：偏移量 %d 会导致源文件集数小于等于0，跳过此季\n", episodeOffset)
			continue
		}

		// 计算需要的位数
		var digits int
		if !padZero {
			digits = 1 // 如果不需要补0，则使用1位数
		} else {
			if episodeContinuous {
				// 连续集数：使用全剧最大集数来确定位数
				digits = len(strconv.Itoa(maxEpisodeNumber))
			} else {
				// 不连续集数：使用当前季最大集数来确定位数
				digits = len(strconv.Itoa(endEp))
			}
			if digits < 2 {
				digits = 2 // 确保至少使用2位数
			}
		}

		// 显示集数范围和对应关系
		if padZero {
			if episodeContinuous {
				fmt.Printf("集数范围：%d-%d（连续，使用%d位数）\n", sourceStartEp, sourceEndEp, digits)
			} else {
				fmt.Printf("集数范围：%d-%d（不连续，使用%d位数）\n", sourceStartEp, sourceEndEp, digits)
			}
		} else {
			fmt.Printf("集数范围：%d-%d（不补0）\n", sourceStartEp, sourceEndEp, digits)
		}
		if episodeOffset != 0 {
			fmt.Printf("集数偏移量：%+d\n", episodeOffset)
			fmt.Printf("原始集数示例：%d → 实际集数：%d\n",
				sourceStartEp, startEp)
		}

		// 构建匹配范围的正则表达式
		var beReplaced string
		if hasSeason {
			beReplaced = fmt.Sprintf("%s.*S%02d(?:E|Ep|EP|[Ee]pisode|[Ee]p)?(%s)",
				regexp.QuoteMeta(fileTitle),
				season.SeasonNumber,
				utils.GenerateRangePattern(sourceStartEp, sourceEndEp, digits))
		} else {
			beReplaced = fmt.Sprintf("%s.*?(?:S\\d{2})?(?:E|Ep|EP|[Ee]pisode|[Ee]p)?(%s)",
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
	if isDateMode {
		fmt.Println("5. 日期模式：使用播出日期匹配文件名，每集生成独立的替换规则")
		fmt.Printf("6. 播出日期格式：YYYYMMDD（如：%s）\n", strings.ReplaceAll(show.FirstAirDate, "-", ""))
		fmt.Println("7. 只处理有播出日期的集数，未获取到播出日期的集数将被跳过")
	} else {
		fmt.Printf("5. 所有集数都使用相同的位数（由最大集数决定），不足位数补0\n")
		fmt.Printf("   例如：如果最大集数是500（3位），则第1集应该写作001\n")
	}
	if !hasSeason {
		fmt.Printf("8. 原文件名不包含季数，仅匹配集数部分\n")
	}
	if episodeOffset != 0 {
		fmt.Printf("9. 被替换词中的集数范围已经过调整，可以直接匹配原文件名中的集数\n")
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
