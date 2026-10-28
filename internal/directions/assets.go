// Package directions manages explicitly downloaded, local-only routing data.
package directions

import _ "embed"

//go:embed bridge/MapsRoutingBridge.class
var bridge []byte

//go:embed licenses/brouter.txt
var routerLicense []byte

type component struct {
	Name, URL, SHA256 string
	Bytes             uint64
}

var components = []component{
	{"brouter.zip", "https://github.com/abrensch/brouter/releases/download/v1.7.10/brouter-1.7.10.zip", "023fec3ba997758e8cd7ab9e1bae52e962af3f00b57683e3de86b84ffad01532", 6724983},
	{"runtime.zip", "https://github.com/adoptium/temurin21-binaries/releases/download/jdk-21.0.12%2B8/OpenJDK21U-jre_x64_windows_hotspot_21.0.12_8.zip", "b8aa18fef5edb69bee8618f99677d66d0873d22cb40d974c15ac9ffcdecf73ba", 48993215},
}

const segmentOrigin = "https://brouter.de/brouter/segments4/"
const Version = "BRouter 1.7.10 · Temurin 21.0.12+8"
