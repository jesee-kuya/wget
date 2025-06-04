package downloader

// Options holds configuration flags passed to the downloader
type Options struct {
	OutputName  string   // -O: custom filename
	OutputDir   string   // -P: directory to save file in
	InputFile   string   // -i: input file with URLs
	RateLimit   float64  // --rate-limit: in bytes per second
	RunInBg     bool     // -B: download in background
	LogFilePath string   // if -B is set, logs are redirected here
	Reject      []string // -R: file suffixes to skip (e.g. []string{"jpg","gif"})
	Exclude     []string // -X: directory paths to skip (e.g. []string{"/js","/assets"})
}
