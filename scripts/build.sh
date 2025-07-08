#!/bin/bash

# 获取脚本所在目录的绝对路径
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
# 获取项目根目录
PROJECT_ROOT="$( cd "$SCRIPT_DIR/.." && pwd )"

# 版本号
VERSION="v1.0.0"

# 程序名称
APP_NAME="rename-by-tmdb"

# 构建目录
BUILD_DIR="${PROJECT_ROOT}/dist"

# 清理构建目录
rm -rf "$BUILD_DIR"
mkdir -p "$BUILD_DIR"

# 编译函数
build() {
    os=$1
    arch=$2
    extension=$3
    
    output="${BUILD_DIR}/${APP_NAME}-${os}-${arch}${extension}"
    
    echo "Building for ${os}/${arch}..."
    cd "$PROJECT_ROOT"
    GOOS=$os GOARCH=$arch go build -o "$output" .
    
    # 创建发布包
    cd "$BUILD_DIR"
    if [ "$os" = "windows" ]; then
        zip "${APP_NAME}-${VERSION}-${os}-${arch}.zip" \
            "${APP_NAME}-${os}-${arch}${extension}" \
            "${PROJECT_ROOT}/README.md" \
            "${PROJECT_ROOT}/.env.example"
    else
        tar -czf "${APP_NAME}-${VERSION}-${os}-${arch}.tar.gz" \
            "${APP_NAME}-${os}-${arch}${extension}" \
            -C "${PROJECT_ROOT}" README.md .env.example
    fi
    
    echo "Built and packaged ${os}/${arch}"
}

echo "Building from $PROJECT_ROOT"
echo "Output directory: $BUILD_DIR"

# 为不同平台构建
# macOS
build "darwin" "amd64" ""      # Intel Mac
build "darwin" "arm64" ""      # Apple Silicon Mac

# Linux
build "linux" "amd64" ""       # 64位 Linux
build "linux" "arm64" ""       # ARM64 Linux

# Windows
build "windows" "amd64" ".exe" # 64位 Windows (x64)
build "windows" "386" ".exe"   # 32位 Windows (x86)
build "windows" "arm64" ".exe" # ARM64 Windows

echo -e "\nBuild complete! Files in ${BUILD_DIR}:"
ls -l "$BUILD_DIR" 