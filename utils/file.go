package utils

import (
	"bufio"
	"encoding/json"
	"os"
	"strings"
)

//覆盖写入
func WriteFileContents(filename string, content string) error {

	/* os.WriteFile takes in file path, a []byte of the file content,
	   and permission bits in case file doesn't exist */

	err := os.WriteFile(filename, []byte(content), 0666)
	return err
}

func WriteSpecialFileContents(filename string, content string) error {

	/* os.WriteFile takes in file path, a []byte of the file content,
	   and permission bits in case file doesn't exist */

	err := os.WriteFile(filename, StringToBytes(content), 0666)
	return err
}

//追加写如字符串
func WriteString(filename string, p string) (n int, err error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return 0, err
	}
	//及时关闭file句柄
	defer file.Close()
	//写入文件时，使用带缓存的 *Writer
	write := bufio.NewWriter(file)
	write.WriteString(p)
	write.Flush()
	return len(p), err
}

//替换文件字符串，不存在指定字符串，则追加
func ReplaceStrInFile(filename string, str string, replace_str string) error {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	context := BytesToString(bytes)
	if strings.Contains(context, str) {
		strings.ReplaceAll(context, str, replace_str)
		return WriteSpecialFileContents(filename, context)
	}
	if !strings.Contains(context, replace_str) {
		_, err := WriteBytes(filename, StringToBytes("\n"+replace_str))
		return err

	}
	return nil
}

//追加写如字符数组
func WriteBytes(filename string, p []byte) (n int, err error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return 0, err
	}
	//及时关闭file句柄
	defer file.Close()
	//写入文件时，使用带缓存的 *Writer
	write := bufio.NewWriter(file)
	write.Write(p)
	write.Flush()
	return len(p), err
}
func ReadFile2String(filename string) string {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return ""
	}
	return string(bytes[:])
}

func ReadFileContents(filename string) map[string]string {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil
	}
	data := make(map[string]string)

	fileText := string(bytes[:]) // fileText is "Hello World!"
	for _, value := range strings.Split(fileText, "\n") {
		if !strings.HasPrefix(value, "#") {
			config := strings.Split(value, "=")
			data[config[0]] = config[1]
		}
	}
	return data
}

//覆盖写入
func WriteFileByte(filename string, p interface{}) error {
	b, err := json.Marshal(&p)
	if err != nil {
		return err
	}
	err = os.WriteFile(filename, b, 0666)
	return err
}

func ReadJsonFile(filename string) ([]byte, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return bytes, err
}
func ReadByteFormJsonFile(filename string) ([]byte, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return bytes, err
	}
	return bytes, err
}
func ChangeStatus(filename string, key map[string]interface{}) error {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	data := make(map[string]interface{})
	err = json.Unmarshal(bytes, &data)
	for a, b := range key {
		data[a] = b
	}
	WriteFileByte(filename, data)
	return err
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	//当为空文件或文件夹存在
	if err == nil {
		return true
	}
	//os.IsNotExist(err)为true，文件或文件夹不存在
	if os.IsNotExist(err) {
		return false
	}
	//其它类型，不确定是否存在
	return false
}
