package packages

import (
	"encoding/json"
	"os/exec"
	"strings"
)

// PackageManager represents different package managers
type PackageManager string

const (
	NPM    PackageManager = "npm"
	Pip    PackageManager = "pip"
	Brew   PackageManager = "brew"
	Cargo  PackageManager = "cargo"
	Go     PackageManager = "go"
	Gem    PackageManager = "gem"
)

// Package represents a package that provides CLI tools
type Package struct {
	Name           string         `json:"name"`
	Version        string         `json:"version"`
	Manager        PackageManager `json:"manager"`
	Binaries       []string       `json:"binaries,omitempty"`
	Location       string         `json:"location,omitempty"`
	Global         bool           `json:"global"`
}

// Detector finds packages from various package managers
type Detector struct {
	enabledManagers []PackageManager
}

// NewDetector creates a new package detector
func NewDetector() *Detector {
	return &Detector{
		enabledManagers: []PackageManager{NPM, Pip, Brew, Cargo, Go, Gem},
	}
}

// DetectAll detects packages from all enabled package managers
func (d *Detector) DetectAll() ([]Package, error) {
	var packages []Package

	for _, manager := range d.enabledManagers {
		pkgs, err := d.detectByManager(manager)
		if err != nil {
			// Skip managers that fail (not installed, etc.)
			continue
		}
		packages = append(packages, pkgs...)
	}

	return packages, nil
}

// detectByManager detects packages for a specific manager
func (d *Detector) detectByManager(manager PackageManager) ([]Package, error) {
	switch manager {
	case NPM:
		return d.detectNPM()
	case Pip:
		return d.detectPip()
	case Brew:
		return d.detectBrew()
	case Cargo:
		return d.detectCargo()
	case Go:
		return d.detectGo()
	case Gem:
		return d.detectGem()
	default:
		return nil, nil
	}
}

// detectNPM detects globally installed npm packages
func (d *Detector) detectNPM() ([]Package, error) {
	cmd := exec.Command("npm", "list", "-g", "--json", "--depth=0")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var result struct {
		Dependencies map[string]struct {
			Version string `json:"version"`
		} `json:"dependencies"`
	}

	if err := json.Unmarshal(output, &result); err != nil {
		return nil, err
	}

	var packages []Package
	for name, info := range result.Dependencies {
		packages = append(packages, Package{
			Name:    name,
			Version: info.Version,
			Manager: NPM,
			Global:  true,
		})
	}

	return packages, nil
}

// detectPip detects installed pip packages
func (d *Detector) detectPip() ([]Package, error) {
	cmd := exec.Command("pip", "list", "--format=json")
	output, err := cmd.Output()
	if err != nil {
		// Try pip3
		cmd = exec.Command("pip3", "list", "--format=json")
		output, err = cmd.Output()
		if err != nil {
			return nil, err
		}
	}

	var result []struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}

	if err := json.Unmarshal(output, &result); err != nil {
		return nil, err
	}

	var packages []Package
	for _, item := range result {
		packages = append(packages, Package{
			Name:    item.Name,
			Version: item.Version,
			Manager: Pip,
			Global:  false,
		})
	}

	return packages, nil
}

// detectBrew detects installed homebrew packages
func (d *Detector) detectBrew() ([]Package, error) {
	cmd := exec.Command("brew", "list", "--versions")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(output), "\n")
	var packages []Package

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) >= 2 {
			packages = append(packages, Package{
				Name:    parts[0],
				Version: parts[1],
				Manager: Brew,
				Global:  true,
			})
		}
	}

	return packages, nil
}

// detectCargo detects installed cargo packages
func (d *Detector) detectCargo() ([]Package, error) {
	cmd := exec.Command("cargo", "install", "--list")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(output), "\n")
	var packages []Package

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, " ") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) >= 2 {
			name := parts[0]
			version := strings.Trim(parts[1], "v:")
			packages = append(packages, Package{
				Name:    name,
				Version: version,
				Manager: Cargo,
				Global:  true,
			})
		}
	}

	return packages, nil
}

// detectGo detects installed go packages
func (d *Detector) detectGo() ([]Package, error) {
	// Go doesn't have a built-in list command, so this is limited
	// We could scan $GOPATH/bin but that requires more work
	return nil, nil
}

// detectGem detects installed ruby gems
func (d *Detector) detectGem() ([]Package, error) {
	cmd := exec.Command("gem", "list", "--local")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(output), "\n")
	var packages []Package

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Format: "package_name (version1, version2)"
		parts := strings.Split(line, " (")
		if len(parts) == 2 {
			name := parts[0]
			version := strings.TrimSuffix(parts[1], ")")
			// Take first version if multiple
			version = strings.Split(version, ",")[0]
			packages = append(packages, Package{
				Name:    name,
				Version: strings.TrimSpace(version),
				Manager: Gem,
				Global:  false,
			})
		}
	}

	return packages, nil
}

// FindPackageByName finds a package by name across all managers
func FindPackageByName(packages []Package, name string) *Package {
	for _, pkg := range packages {
		if pkg.Name == name {
			return &pkg
		}
	}
	return nil
}

// GroupByManager groups packages by their package manager
func GroupByManager(packages []Package) map[PackageManager][]Package {
	grouped := make(map[PackageManager][]Package)
	for _, pkg := range packages {
		grouped[pkg.Manager] = append(grouped[pkg.Manager], pkg)
	}
	return grouped
}
