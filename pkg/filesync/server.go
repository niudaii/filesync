package filesync

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"github.com/projectdiscovery/gologger"
	"net"
	"os"
	"path/filepath"
)

type GobConn struct {
	cnRd *bufio.Reader
	cnWt *bufio.Writer
	Dec  *gob.Decoder
	enc  *gob.Encoder
}

// initGobConn 初始化连接
func initGobConn(conn net.Conn) (gbc *GobConn) {
	gbc = &GobConn{}
	gbc.cnRd = bufio.NewReader(conn)
	gbc.cnWt = bufio.NewWriter(conn)
	gbc.Dec = gob.NewDecoder(gbc.cnRd)
	gbc.enc = gob.NewEncoder(gbc.cnWt)
	return
}

// gobConnWt 发送消息与数据
func (gbc *GobConn) gobConnWt(mg interface{}) (err error) {
	err = gbc.enc.Encode(mg)
	if err != nil {
		return
	}
	err = gbc.cnWt.Flush()
	return
}

// StartFileSyncServer 启动文件同步服务监听
func StartFileSyncServer(host, port, auth, dir string, blackList []string) {
	syncFileBlackList = blackList
	serverAddr := fmt.Sprintf("%s:%s", host, port)
	srv, err := net.Listen("tcp", serverAddr)
	if err != nil {
		gologger.Error().Msgf("net.Listen() err", err)
		return
	}
	var conn net.Conn
	for {
		conn, err = srv.Accept()
		if err != nil {
			gologger.Error().Msgf("srv.Accept() err", err)
			continue
		}
		go handleSync(conn, auth, dir)
	}
}

func handleSync(conn net.Conn, auth, dir string) {
	defer conn.Close()
	gbc := initGobConn(conn)
	for {
		mg := Message{}
		err := gbc.Dec.Decode(&mg)
		if err != nil {
			writeErrorMg("gbc Decode() error!", gbc)
			return
		}
		// 检查authKey，如果不通过直接返回
		if !checkSyncAuthKey(auth, mg.MgAuthKey) {
			writeErrorMg("authKey error!", gbc)
			return
		}
		switch mg.MgType {
		// 请求同步
		case MsgSync:
			if err = hdSync(gbc, dir); err != nil {
				gologger.Error().Msgf("hdSync() err", err)
			}
		// 请求文件传输
		case MsgTran:
			if err = hdTranFile(&mg, gbc, dir); err != nil {
				gologger.Error().Msgf("hdTranFile() err", err)
			}
		// 结束
		case MsgEnd:
			return
		// 未知消息
		default:
			writeErrorMg("error, not a recognizable message.", gbc)
			return
		}
	}
}

// hdSync 处理worker的全部文件同步请求
func hdSync(gbc *GobConn, dir string) (err error) {
	var srcPath string
	srcPath, err = filepath.Abs(dir)
	if err != nil {
		writeErrorMg("filepath.Abs() error", gbc)
		return
	}
	var fileMd5List []string
	fileMd5List, err = Traverse(srcPath)
	if err != nil {
		writeErrorMg("Traverse() error", gbc)
		return
	}
	if len(fileMd5List) == 0 {
		writeErrorMg("emtry file list", gbc)
		return
	}
	cr := Message{
		MgStrings: fileMd5List,
		MgType:    MsgMd5List,
		Overwrite: true,
	}
	err = gbc.gobConnWt(cr)
	return
}

// hdTranFile 向worker同步一个文件
func hdTranFile(mg *Message, gbc *GobConn, dir string) (err error) {
	if len(mg.MgString) <= 0 {
		writeErrorMg("no file to transfer", gbc)
		return
	}
	var srcPath string
	srcPath, err = filepath.Abs(dir)
	if err != nil {
		writeErrorMg(fmt.Sprintf("read sync file:%s error", srcPath), gbc)
		return
	}
	srcPathFileName := filepath.Join(srcPath, mg.MgString)
	var cr Message
	st, err := os.Stat(srcPathFileName)
	if err != nil {
		writeErrorMg(fmt.Sprintf("read sync file:%s error", err), gbc)
		return
	}
	cr.MgFileMode = st.Mode()
	cr.MgByte, err = os.ReadFile(srcPathFileName)
	if err != nil {
		writeErrorMg(fmt.Sprintf("read sync file:%s error", srcPathFileName), gbc)
		return
	}
	cr.MgType = MsgTranData
	err = gbc.gobConnWt(cr)
	return
}

// writeErrorMg 返回错误信息的消息
func writeErrorMg(message string, gbc *GobConn) {
	var errMsg Message
	errMsg.MgType = MsgError
	errMsg.MgString = message
	if err := gbc.gobConnWt(errMsg); err != nil {
		gologger.Error().Msgf("gbc.gobConnWt() err", err)
	}
}

// checkSyncAuthKey 同步的认证检查
func checkSyncAuthKey(authKey, workerAuthKey string) (success bool) {
	if workerAuthKey == authKey {
		return true
	}
	return false
}
