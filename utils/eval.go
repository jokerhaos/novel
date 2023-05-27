package utils

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func Eval(packageName string, funcName string, args ...interface{}) error {
	// 动态导入包
	pkg, err := importPackage(packageName)
	if err != nil {
		return errors.New("Failed to import package:" + err.Error())
	}

	// 获取函数的反射值
	funcValue := reflect.ValueOf(pkg).MethodByName(funcName)

	// 检查函数是否存在
	if !funcValue.IsValid() {
		return errors.New("Failed to import package:" + funcName)
	}

	// 调用函数
	funcValue.Call(nil)
	return nil
}

// 动态导入包
func importPackage(packageName string) (interface{}, error) {
	pkg, err := importPackageByPath(packageName)
	if err != nil {
		// 尝试将包路径转换为小写再次导入
		pkg, err = importPackageByPath(strings.ToLower(packageName))
		if err != nil {
			return nil, err
		}
	}

	return pkg, nil
}

// 根据包路径动态导入包
func importPackageByPath(packagePath string) (interface{}, error) {
	pkg, err := importWithReflect(packagePath)
	if err != nil {
		// 如果导入失败，尝试使用全局包名导入
		pkg, err = importWithReflect(getPackageNameFromPath(packagePath))
		if err != nil {
			return nil, err
		}
	}

	return pkg, nil
}

// 使用 reflect.ValueOf 导入包
func importWithReflect(packagePath string) (interface{}, error) {
	pkgValue := reflect.New(reflect.TypeOf(nil)).Elem()
	pkgType := pkgValue.Type()
	for i := 0; i < pkgValue.NumField(); i++ {
		field := pkgType.Field(i)
		if field.PkgPath == packagePath {
			pkgValue = pkgValue.Field(i)
			break
		}
	}

	if pkgValue.IsValid() {
		return pkgValue.Interface(), nil
	}

	return nil, fmt.Errorf("Package not found: %s", packagePath)
}

// 从包路径中获取包名
func getPackageNameFromPath(packagePath string) string {
	parts := strings.Split(packagePath, "/")
	return parts[len(parts)-1]
}
