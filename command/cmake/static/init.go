package static

import (
	"embed"
)

//go:embed cmake/cc_binary.cmake
//go:embed cmake/cc_library.cmake
//go:embed cmake/cc_test.cmake
//go:embed cmake/README.md
//go:embed main.cpp.txt
//go:embed CMakeLists.txt
var cmake_fs embed.FS

const CMakeListsFile = "CMakeLists.txt"

func ReadAllFiles() (map[string]string, error) {
	result := make(map[string]string, 0)

	readFile := func(name string) error {
		file, err := cmake_fs.ReadFile(name)
		if err != nil {
			return err
		}
		if name == "main.cpp.txt" { // cgo 编译会扫描文件下的.c/.cpp文件
			name = "main.cpp"
		}
		result[name] = string(file)
		return nil
	}

	if err := readFile("cmake/cc_binary.cmake"); err != nil {
		return nil, err
	}
	if err := readFile("cmake/cc_library.cmake"); err != nil {
		return nil, err
	}
	if err := readFile("cmake/cc_test.cmake"); err != nil {
		return nil, err
	}
	if err := readFile("cmake/README.md"); err != nil {
		return nil, err
	}
	if err := readFile("main.cpp.txt"); err != nil {
		return nil, err
	}
	if err := readFile(CMakeListsFile); err != nil {
		return nil, err
	}
	return result, nil
}
