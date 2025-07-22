#!/bin/bash

# 获取脚本所在目录的绝对路径
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
# 获取项目根目录
PROJECT_ROOT="$( cd "$SCRIPT_DIR/.." && pwd )"
# 设置输出目录
OUTPUT_DIR="$PROJECT_ROOT/dist"

echo "Building from $PROJECT_ROOT"
echo "Output directory: $OUTPUT_DIR"

# 创建输出目录
mkdir -p "$OUTPUT_DIR"

# 构建函数
build_binary() {
    local cmd=$1
    local os=$2
    local arch=$3
    local ext=""
    local binary_name="rename-by-tmdb"
    local list_binary_name="list"

    # Windows二进制文件添加.exe后缀
    if [ "$os" = "windows" ]; then
        ext=".exe"
    fi

    echo "Building $cmd for $os/$arch..."
    
    # 设置交叉编译环境变量
    export GOOS=$os
    export GOARCH=$arch

    if [ "$cmd" = "rename" ]; then
        # 构建主程序
        go build -o "$OUTPUT_DIR/${binary_name}-${os}-${arch}${ext}" "$PROJECT_ROOT"
        
        # 打包文件
        if [ "$os" = "windows" ]; then
            (cd "$OUTPUT_DIR" && zip "${binary_name}-v1.0.3-${os}-${arch}.zip" "${binary_name}-${os}-${arch}${ext}" "$PROJECT_ROOT/README.md" "$PROJECT_ROOT/.env.example")
        else
            tar -czf "$OUTPUT_DIR/${binary_name}-v1.0.3-${os}-${arch}.tar.gz" -C "$OUTPUT_DIR" "${binary_name}-${os}-${arch}" -C "$PROJECT_ROOT" "README.md" ".env.example"
        fi
    elif [ "$cmd" = "list" ]; then
        # 构建list命令
        go build -o "$OUTPUT_DIR/${list_binary_name}-${os}-${arch}${ext}" "$PROJECT_ROOT/cmd/list"
        
        # 打包文件
        if [ "$os" = "windows" ]; then
            (cd "$OUTPUT_DIR" && zip "${list_binary_name}-v1.0.3-${os}-${arch}.zip" "${list_binary_name}-${os}-${arch}${ext}")
        else
            tar -czf "$OUTPUT_DIR/${list_binary_name}-v1.0.3-${os}-${arch}.tar.gz" -C "$OUTPUT_DIR" "${list_binary_name}-${os}-${arch}"
        fi
    fi

    echo "Built and packaged $os/$arch"
}

# 构建所有平台的二进制文件
for cmd in "rename" "list"; do
    # macOS (Intel & Apple Silicon)
    build_binary $cmd "darwin" "amd64"
    build_binary $cmd "darwin" "arm64"

    # Linux (x86_64 & ARM64)
    build_binary $cmd "linux" "amd64"
    build_binary $cmd "linux" "arm64"

    # Windows (x86_64, x86 & ARM64)
    build_binary $cmd "windows" "amd64"
    build_binary $cmd "windows" "386"
    build_binary $cmd "windows" "arm64"
done

echo -e "\nBuild complete! Files in $OUTPUT_DIR:"
ls -l "$OUTPUT_DIR" 