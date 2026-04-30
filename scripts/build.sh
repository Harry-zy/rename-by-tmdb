#!/bin/bash
set -euo pipefail

# 获取脚本所在目录的绝对路径
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
# 获取项目根目录
PROJECT_ROOT="$( cd "$SCRIPT_DIR/.." && pwd )"
# 设置输出目录
OUTPUT_DIR="$PROJECT_ROOT/dist"
VERSION="v1.0.6"

echo "Building from $PROJECT_ROOT"
echo "Output directory: $OUTPUT_DIR"

# 创建输出目录
mkdir -p "$OUTPUT_DIR"

# 构建函数
build_binary() {
    local os=$1
    local arch=$2
    local ext=""
    local binary_name="rename-by-tmdb"

    # Windows二进制文件添加.exe后缀
    if [ "$os" = "windows" ]; then
        ext=".exe"
    fi

    echo "Building $binary_name for $os/$arch..."

    # 设置交叉编译环境变量
    export GOOS=$os
    export GOARCH=$arch

    # 构建主程序
    go build -o "$OUTPUT_DIR/${binary_name}-${os}-${arch}${ext}" "$PROJECT_ROOT"

    # 打包文件
    if [ "$os" = "windows" ]; then
        local archive="$OUTPUT_DIR/${binary_name}-${VERSION}-${os}-${arch}.zip"
        rm -f "$archive"
        zip -j "$archive" \
            "$OUTPUT_DIR/${binary_name}-${os}-${arch}${ext}" \
            "$PROJECT_ROOT/README.md" \
            "$PROJECT_ROOT/.env.example"
    else
        local archive="$OUTPUT_DIR/${binary_name}-${VERSION}-${os}-${arch}.tar.gz"
        rm -f "$archive"
        tar -czf "$archive" \
            -C "$OUTPUT_DIR" "${binary_name}-${os}-${arch}" \
            -C "$PROJECT_ROOT" "README.md" ".env.example"
    fi

    echo "Built and packaged $os/$arch"
}

# 构建所有平台的二进制文件
# macOS (Intel & Apple Silicon)
build_binary "darwin" "amd64"
build_binary "darwin" "arm64"

# Linux (x86_64 & ARM64)
build_binary "linux" "amd64"
build_binary "linux" "arm64"

# Windows (x86_64, x86 & ARM64)
build_binary "windows" "amd64"
build_binary "windows" "386"
build_binary "windows" "arm64"

echo -e "\nBuild complete! Files in $OUTPUT_DIR:"
ls -l "$OUTPUT_DIR"
