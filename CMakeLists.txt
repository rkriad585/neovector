cmake_minimum_required(VERSION 3.20)
project(neovector
    VERSION 0.1.0
    DESCRIPTION "Image-to-vector conversion CLI"
    LANGUAGES NONE
)

# ── Options ────────────────────────────────────────────
option(BUILD_RELEASE "Build with release ldflags" ON)

# ── Version from .version file ─────────────────────────
if(EXISTS "${CMAKE_SOURCE_DIR}/.version")
    file(READ "${CMAKE_SOURCE_DIR}/.version" VERSION_RAW)
    string(STRIP "${VERSION_RAW}" PROJECT_VERSION)
endif()

# ── Git commit ─────────────────────────────────────────
execute_process(
    COMMAND git rev-parse --short=8 HEAD
    WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}
    OUTPUT_VARIABLE GIT_COMMIT
    OUTPUT_STRIP_TRAILING_WHITESPACE
    ERROR_QUIET
)

if(NOT GIT_COMMIT)
    set(GIT_COMMIT "unknown")
endif()

# ── Custom target: build ───────────────────────────────
set(LDFLAGS "-s -w")
set(LDFLAGS "${LDFLAGS} -X main.Version=${PROJECT_VERSION}")
set(LDFLAGS "${LDFLAGS} -X main.Commit=${GIT_COMMIT}")
set(LDFLAGS "${LDFLAGS} -X main.PublisherName=rkriad585")
set(LDFLAGS "${LDFLAGS} -X main.PublisherEmail=rkriad585@gmail.com")

add_custom_target(build
    COMMAND go build -trimpath -ldflags "${LDFLAGS}" -o neovector .
    WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}
    COMMENT "Building neovector ${PROJECT_VERSION} (${GIT_COMMIT})"
)

add_custom_target(release
    COMMAND ${CMAKE_COMMAND} -E echo "Building release binaries..."
    COMMAND GOOS=linux   GOARCH=amd64 go build -trimpath -ldflags "${LDFLAGS}" -o bin/neovector-linux-amd64     .
    COMMAND GOOS=linux   GOARCH=arm64 go build -trimpath -ldflags "${LDFLAGS}" -o bin/neovector-linux-arm64     .
    COMMAND GOOS=darwin  GOARCH=amd64 go build -trimpath -ldflags "${LDFLAGS}" -o bin/neovector-darwin-amd64    .
    COMMAND GOOS=darwin  GOARCH=arm64 go build -trimpath -ldflags "${LDFLAGS}" -o bin/neovector-darwin-arm64    .
    COMMAND GOOS=windows GOARCH=amd64 go build -trimpath -ldflags "${LDFLAGS}" -o bin/neovector-windows-amd64.exe .
    COMMAND GOOS=windows GOARCH=arm64 go build -trimpath -ldflags "${LDFLAGS}" -o bin/neovector-windows-arm64.exe .
    WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}
    COMMENT "Cross-compiling all release targets"
)

add_custom_target(clean
    COMMAND rm -f neovector neovector.exe
    COMMAND rm -rf bin/
    WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}
)

add_custom_target(test
    COMMAND go test ./...
    WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}
)

add_custom_target(lint
    COMMAND go vet ./...
    WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}
)
