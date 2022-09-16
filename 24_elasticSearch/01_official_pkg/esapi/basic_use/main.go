package main

import "log"

func main() {
	log.SetFlags(0)

	err := initClient()
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	// 1. Get cluster info
	getClusterInfo()

	// 2. Index documents concurrently
	insertIndex()

	// 3. Search for the indexed documents
	searchIndex()
	// Use the helper methods of the client.

}

// 结果

//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// [200 OK] updated; version=1
// [200 OK] updated; version=1
// -------------------------------------
// [200 OK] 2 hits; took: 7ms
//  * ID=1, map[title:Test One]
//  * ID=2, map[title:Test Two]
// =====================================
