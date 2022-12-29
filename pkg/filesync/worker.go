package filesync

import (
	"fmt"
	"github.com/jacenr/filediff/diff"
	"log"
	"net"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// WorkerStartupSync worker在启动时进行文件同步
func WorkerStartupSync(host, port, auth, dir string) (err error) {
	serverAddr := fmt.Sprintf("%s:%s", host, port)
	// 1 连接到server
	var conn net.Conn
	conn, err = net.Dial("tcp", serverAddr)
	if err != nil {
		return
	}
	defer conn.Close()
	gbc := initGobConn(conn)
	// 2 发送SYNC请求
	msgSync := Message{MgType: MsgSync, MgAuthKey: auth}
	err = gbc.gobConnWt(msgSync)
	if err != nil {
		return
	}
	// 3 服务器返回信息
	var hostMessage Message
	err = gbc.Dec.Decode(&hostMessage)
	if err != nil {
		return
	}
	if hostMessage.MgType != MsgMd5List {
		err = fmt.Errorf("无法获取同步文件信息:%v", hostMessage.MgString)
		return
	}
	// 4 获取服务器所有文件及md5值,并预处理本地的路径和文件
	var transFiles []string
	transFiles, err = doFileMd5List(&hostMessage, dir)
	if err == nil {
		if len(transFiles) > 0 {
			log.Printf("需要同步的文件数量: %v\n", len(transFiles))
			// 5 同步文件
			for i, file := range transFiles {
				err = doTranFile(file, auth, gbc, dir)
				if err != nil {
					log.Printf("%v %v 同步失败\n", i+1, file)
				} else {
					log.Printf("%v %v 同步成功\n", i+1, file)
				}
			}
			log.Println("同步完成")
		}
	}
	// 6 结束同步
	endMsg := Message{MgType: MsgEnd, MgAuthKey: auth}
	err = gbc.gobConnWt(endMsg)
	return
}

// doFileMd5List 读取worker本地文件列表及md5值，并与服务端进行对比，确定需要同步的文件列表
func doFileMd5List(mg *Message, dir string) (transFiles []string, err error) {
	var slinkNeedCreat = make(map[string]string)
	var slinkNeedChange = make(map[string]string)
	var needDelete = make([]string, 0)
	var needCreDir = make([]string, 0)
	var srcPath string
	srcPath, err = filepath.Abs(dir)
	if err != nil {
		return
	}
	// 遍历本地目标路径失败
	var localFilesMd5 []string
	localFilesMd5, err = Traverse(srcPath)
	if err != nil {
		return
	}
	sort.Strings(localFilesMd5)
	var diffrm []string
	var diffadd []string
	if len(localFilesMd5) != 0 {
		diffrm, diffadd = diff.DiffOnly(mg.MgStrings, localFilesMd5)
	} else {
		diffrm, diffadd = mg.MgStrings, localFilesMd5
	}
	if len(diffrm) == 0 && len(diffadd) == 0 {
		return
	}
	// 重组成map
	diffrmM := make(map[string]string)
	diffaddM := make(map[string]string)
	for _, v := range diffrm {
		s := strings.Split(v, ",,")
		if len(s) != 1 {
			diffrmM[s[0]] = s[1]
		}
	}
	for _, v := range diffadd {
		s := strings.Split(v, ",,")
		if len(s) != 1 {
			diffaddM[s[0]] = s[1]
		}
	}
	// 整理
	for k := range diffaddM {
		v2, ok := diffrmM[k]
		if ok {
			if !mg.Overwrite {
				delete(diffrmM, k)
			}
			if mg.Overwrite {
				if strings.HasPrefix(v2, "symbolLink&&") {
					slinkNeedChange[k] = strings.TrimPrefix(v2, "symbolLink&&")
					delete(diffrmM, k)
				}
				needDelete = append(needDelete, k)
			}
		}
		if !ok && mg.Del {
			needDelete = append(needDelete, k)
		}

	}
	for k, v := range diffrmM {
		if strings.HasPrefix(v, "symbolLink&&") {
			slinkNeedCreat[k] = strings.TrimPrefix(v, "symbolLink&&")
			delete(diffrmM, k)
			continue
		}
		if v == "Directory" {
			needCreDir = append(needCreDir, k)
			delete(diffrmM, k)
		}
	}
	// 接收新文件的本地操作
	err = os.Chdir(srcPath)
	if err != nil {
		return
	}
	defer os.Chdir(cwd)

	err = localOP(slinkNeedCreat, slinkNeedChange, needDelete, needCreDir)
	if err != nil {
		return
	}
	// do request needTrans files
	for k := range diffrmM {
		transFiles = append(transFiles, k)
	}
	sort.Strings(transFiles)
	return
}

// doTranFile worker向server请求同步一个文件
func doTranFile(filePathName, authKey string, gbc *GobConn, dir string) (err error) {
	var srcPath string
	srcPath, err = filepath.Abs(dir)
	if err != nil {
		return
	}
	mg := Message{
		MgAuthKey: authKey,
		MgType:    MsgTran,
		MgString:  filePathName,
	}
	err = gbc.gobConnWt(mg)
	if err != nil {
		return
	}
	var hostMessage Message
	err = gbc.Dec.Decode(&hostMessage)
	if err != nil {
		return
	}
	if hostMessage.MgType == MsgTranData {
		dstFilePathName := filepath.Join(srcPath, filePathName)
		err = os.WriteFile(dstFilePathName, hostMessage.MgByte, hostMessage.MgFileMode)
		if err != nil {
			return
		}
		return
	}
	err = fmt.Errorf("msgTranData error")
	return
}
