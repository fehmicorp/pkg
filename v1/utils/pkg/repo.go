package pkg

var supportedOS = map[string]bool{
	"windows": true,
	"linux":   true,
	"darwin":  true,
	"ios":     true,
	"android": true,
}

var supportedArch = map[string]bool{
	"amd64":   true,
	"x86_64":  true, // Alias for amd64
	"386":     true,
	"x86":     true, // Alias for 386
	"arm64":   true,
	"aarch64": true, // Alias for arm64
	"arm":     true, // 32-bit ARM (Raspberry Pi, older mobile devices)
	"armv7":   true, // Alias for 32-bit ARM v7
	"riscv64": true,
	"ppc64le": true,
	"s390x":   true,
}

type Repo struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	URL      string `json:"url"`
	LocalDir string `json:"dir"`
	Enabled  bool   `json:"enabled"`
}
