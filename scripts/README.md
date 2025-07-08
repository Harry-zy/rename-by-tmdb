# rename-by-tmdb v1.0.0

基于 TMDB API 的影视重命名工具，支持电影和剧集的标准命名格式生成和重命名规则。

## 新特性

- 支持电影和剧集的重命名规则生成
- 从 TMDB API 获取影视信息
- 生成标准的命名格式，包含年份和 TMDB ID
- 生成正则表达式替换规则
- 支持剧集集数偏移（适用于合集重命名）
- 支持剧集季数匹配（可选保留原文件名季数）
- 可选上传规则到远程服务器

## 下载

请根据您的系统选择对应的版本：

### Windows
- [rename-by-tmdb-v1.0.0-windows-amd64.zip](rename-by-tmdb-v1.0.0-windows-amd64.zip) - x64系统（推荐）
- [rename-by-tmdb-v1.0.0-windows-arm64.zip](rename-by-tmdb-v1.0.0-windows-arm64.zip) - ARM64系统（Surface Pro X等搭载ARM处理器的设备）
- [rename-by-tmdb-v1.0.0-windows-386.zip](rename-by-tmdb-v1.0.0-windows-386.zip) - 32位系统

### macOS
- [rename-by-tmdb-v1.0.0-darwin-amd64.tar.gz](rename-by-tmdb-v1.0.0-darwin-amd64.tar.gz) - Intel芯片
- [rename-by-tmdb-v1.0.0-darwin-arm64.tar.gz](rename-by-tmdb-v1.0.0-darwin-arm64.tar.gz) - Apple Silicon/M系列芯片

### Linux
- [rename-by-tmdb-v1.0.0-linux-amd64.tar.gz](rename-by-tmdb-v1.0.0-linux-amd64.tar.gz) - x86_64系统
- [rename-by-tmdb-v1.0.0-linux-arm64.tar.gz](rename-by-tmdb-v1.0.0-linux-arm64.tar.gz) - ARM64系统

## 使用说明

详细的使用说明请参考 [README.md](https://github.com/Harry-zy/rename-by-tmdb/blob/v1.0.0/README.md)

## 系统要求

- Windows 7 及以上 / macOS 10.13 及以上 / 主流 Linux 发行版
- 如需上传功能，需要有效的远程服务器认证令牌

## 注意事项

1. Windows用户如何选择正确的版本：
    - 右键"此电脑" -> 属性
    - 查看"系统类型"
    - 如果是 Surface Pro X 等 ARM 设备，选择 arm64 版本
    - 如果显示"64位操作系统"，选择 amd64 版本
    - 如果显示"32位操作系统"，选择 386 版本
2. 大多数2010年后的电脑都应该使用64位(amd64)版本
3. 使用前请确保已正确配置 TMDB API 密钥
4. 建议在文件名中包含年份等信息，以避免与其他影视作品重名

## 更新日志

- 首次发布
- 支持电影和剧集的重命名规则生成
- 支持多平台（Windows、macOS、Linux）
- 提供完整的命名规则生成和上传功能 