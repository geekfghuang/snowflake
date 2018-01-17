A distributed, thread safe and self increasing ID generator based on the theory
of Twitter-Snowflake 64 bit self increasing ID algorithm

Usage：

    http:
        ./uid-http ../conf/worker.xml

    application：
        go get github.com/geekfghuang/snowflake

        func main() {
            worker, err := NewWorker(0)
        	if err != nil {
                // ...
        	}
        	id, err := iw.NextId()
        	// ...
        }