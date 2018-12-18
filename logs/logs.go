package logs

import (
	"flag"
)

func Init() {
	// glog need do this, let it don't write log file
	flag.Parse()
	flag.Set("logtostderr", "true")

}

// // GlogWriter serves as a bridge between the standard log package and the glog package.
// type GlogWriter struct{}

// // Write implements the io.Writer interface, Depth 4 is the caller current frame.
// func (writer GlogWriter) Write(data []byte) (n int, err error) {
// 	glog.InfoDepth(4, string(data))
// 	return len(data), nil
// }
