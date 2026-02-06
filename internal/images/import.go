package images

import (
	"archive/tar"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"example.com/containeredu/internal/paths"
)

type ManifestEntry struct {
	Config   string   `json:"Config"`
	RepoTags []string `json:"RepoTags"`
	Layers   []string `json:"Layers"`
}

// ImportDockerSaveTar imports a docker save tarball into local image store.
// It extracts layer tarballs into images/<name>/layers/<n>/ and writes metadata.json.
func ImportDockerSaveTar(tarPath, name string) error {
	if err := paths.EnsureDirs(); err != nil {
		return err
	}
	f, err := os.Open(tarPath)
	if err != nil {
		return err
	}
	defer f.Close()
	tr := tar.NewReader(f)

	var manifest []ManifestEntry
	tempDir := filepath.Join(os.TempDir(), "cede-import")
	_ = os.RemoveAll(tempDir)
	if err := os.MkdirAll(tempDir, 0o755); err != nil {
		return err
	}
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		target := filepath.Join(tempDir, hdr.Name)
		switch hdr.Typeflag {
		case tar.TypeReg, tar.TypeRegA:
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				return err
			}
			out, err := os.Create(target)
			if err != nil {
				return err
			}
			if _, err := io.Copy(out, tr); err != nil {
				out.Close()
				return err
			}
			out.Close()
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0o755); err != nil {
				return err
			}
		default:
			// ignore other types
		}
	}
	manifestBytes, err := os.ReadFile(filepath.Join(tempDir, "manifest.json"))
	if err != nil {
		return fmt.Errorf("manifest.json missing: %w", err)
	}
	if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
		return err
	}
	if len(manifest) == 0 {
		return errors.New("empty manifest")
	}
	entry := manifest[0]
	imgRoot := filepath.Join(paths.ImagesRoot(), name)
	layersRoot := filepath.Join(imgRoot, "layers")
	if err := os.MkdirAll(layersRoot, 0o755); err != nil {
		return err
	}
	for i, l := range entry.Layers {
		src := filepath.Join(tempDir, l)
		dst := filepath.Join(layersRoot, fmt.Sprintf("%02d", i))
		if err := os.MkdirAll(dst, 0o755); err != nil {
			return err
		}
		if err := extractTar(src, dst); err != nil {
			return fmt.Errorf("extract layer %s: %w", l, err)
		}
	}
	meta := struct {
		Name   string   `json:"name"`
		Layers []string `json:"layers"`
	}{
		Name:   name,
		Layers: entry.Layers,
	}
	metaBytes, _ := json.MarshalIndent(meta, "", "  ")
	if err := os.WriteFile(filepath.Join(imgRoot, "metadata.json"), metaBytes, 0o644); err != nil {
		return err
	}
	_ = os.RemoveAll(tempDir)
	return nil
}

func extractTar(srcTar, dstDir string) error {
	f, err := os.Open(srcTar)
	if err != nil {
		return err
	}
	defer f.Close()
	tr := tar.NewReader(f)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		target := filepath.Join(dstDir, hdr.Name)
		switch hdr.Typeflag {
		case tar.TypeReg, tar.TypeRegA:
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				return err
			}
			out, err := os.Create(target)
			if err != nil {
				return err
			}
			if _, err := io.Copy(out, tr); err != nil {
				out.Close()
				return err
			}
			if err := os.Chmod(target, os.FileMode(hdr.Mode)); err != nil {
				out.Close()
				return err
			}
			out.Close()
		case tar.TypeDir:
			if err := os.MkdirAll(target, os.FileMode(hdr.Mode)); err != nil {
				return err
			}
		case tar.TypeSymlink:
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				return err
			}
			if err := os.Symlink(hdr.Linkname, target); err != nil && !os.IsExist(err) {
				return err
			}
		default:
			// ignore
		}
	}
	return nil
}
