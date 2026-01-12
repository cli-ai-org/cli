package cmd

import (
	"os"

	"github.com/cli-ai-org/cli/internal/display"
	"github.com/cli-ai-org/cli/internal/packages"
	"github.com/cli-ai-org/cli/internal/scanner"
	"github.com/spf13/cobra"
)

var (
	listAll  bool
	listJSON bool
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List package-managed CLI tools",
	Long: `Lists CLI tools installed through package managers (npm, pip, brew, cargo, gem, etc.).

By default, shows only tools from known packages to provide a clean list of intentionally
installed CLI tools. Use --all flag to show all executables in your PATH.

Use --json flag to output in JSON format for programmatic access or AI agent consumption.`,
	Example: `  # List package-managed CLI tools (default)
  cli list

  # List ALL executables in PATH
  cli list --all

  # List in JSON format for AI agents
  cli list --json`,
	Run: func(cmd *cobra.Command, args []string) {
		s := scanner.New()
		d := display.New(os.Stdout)

		// Scan for tools
		tools, err := s.ScanAllDetailed()
		if err != nil {
			cmd.PrintErrf("Error scanning for tools: %v\n", err)
			os.Exit(1)
		}

		// By default, show only tools from packages (unless --all is specified)
		if !listAll {
			detector := packages.NewDetector()
			pkgs, err := detector.DetectAll()
			if err != nil {
				cmd.PrintErrf("Error detecting packages: %v\n", err)
				os.Exit(1)
			}

			linker := packages.NewLinker(pkgs)
			linkedTools := linker.LinkTools(tools)

			// Count binaries per package to filter out library packages
			pkgBinaryCount := make(map[string]int)
			for _, tool := range linkedTools {
				if tool.PackageName != "" {
					pkgBinaryCount[tool.PackageName]++
				}
			}

			// Packages to exclude (development libraries, codecs, not user-facing CLIs)
			excludePackages := map[string]bool{
				"gcc": true, "netpbm": true, "gd": true, "gdal": true,
				"gettext": true, "libtiff": true, "libpng": true, "fontconfig": true,
				"glib": true, "hdf5": true, "graphviz": true, "gts": true,
				"mbedtls": true, "nss": true, "perl": true, "tesseract": true,
				"pcre": true, "pcre2": true, "python@3.11": true, "python@3.13": true,
				"xz": true, "ffmpeg": true, "libsndfile": true, "little-cms2": true,
				"jpeg-xl": true, "libfido2": true, "libgcrypt": true, "libheif": true,
				"c-ares": true, "libtasn1": true, "libavif": true, "libbluray": true,
				"cairo": true, "jpeg-turbo": true, "zeromq": true, "tcl-tk": true,
				"libdap": true, "libde265": true, "libgeotiff": true, "libidn2": true,
				"librist": true, "libvmaf": true, "lua": true, "autoconf": true,
				"brotli": true, "flac": true, "giflib": true, "lame": true,
				"leptonica": true, "libassuan": true, "libdeflate": true, "libevent": true,
				"libgpg-error": true, "libksba": true, "lz4": true, "m4": true,
				"miniupnpc": true, "mpg123": true, "nettle": true, "nghttp2": true,
				"oniguruma": true, "openexr": true, "openjpeg": true, "opus": true,
				"p11-kit": true, "pango": true, "pkgconf": true, "proj": true,
				"qhull": true, "rav1e": true, "rubberband": true, "sdl2": true,
				"speex": true, "srt": true, "unbound": true, "uriparser": true,
				"webp": true, "x264": true, "x265": true, "dav1d": true, "aom": true,
				"gnupg": true, "gnutls": true, "gpgme": true, "gobject-introspection": true,
				"grpc": true, "guile": true, "harfbuzz": true, "jasper": true,
				"jemalloc": true, "libtool": true, "nspr": true,
				"cfitsio": true, "gdbm": true, "netcdf": true, "freetype": true,
				"fribidi": true, "fmt": true, "gdk-pixbuf": true, "geos": true,
				"gflags": true, "fizz": true, "epsilon": true, "unixodbc": true,
				"openssl@3": true, "shared-mime-info": true,
				"apache-arrow": true, "protobuf": true, "protobuf@29": true,
			}

			// Get unique CLI names that are linked to packages
			// Exclude packages with too many binaries (>10) or known dev libraries
			seenTools := make(map[string]bool)
			var cliTools []string
			for _, tool := range linkedTools {
				pkgName := tool.PackageName
				if pkgName != "" && !seenTools[tool.Name] {
					// Skip if package has too many binaries or is in exclude list
					if pkgBinaryCount[pkgName] > 10 || excludePackages[pkgName] {
						continue
					}
					cliTools = append(cliTools, tool.Name)
					seenTools[tool.Name] = true
				}
			}
			d.ShowTools(cliTools)
			return
		}

		// With --all, show all executables
		if listJSON {
			if err := d.ShowToolsJSON(tools, true); err != nil {
				cmd.PrintErrf("Error encoding JSON: %v\n", err)
				os.Exit(1)
			}
		} else {
			// Simple name list
			var names []string
			for _, tool := range tools {
				names = append(names, tool.Name)
			}
			d.ShowTools(names)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVarP(&listAll, "all", "a", false, "show ALL executables in PATH (not just package-managed)")
	listCmd.Flags().BoolVarP(&listJSON, "json", "j", false, "output in JSON format for AI agents")
}
