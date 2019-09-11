package archive

import (
	"archive/tar"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"time"
)

func WriteTarFile(sourceDir, dest string, uid, gid int) (string, error) {
	hasher := sha256.New()
	f, err := os.Create(dest)
	if err != nil {
		return "", err
	}
	defer f.Close()
	w := io.MultiWriter(hasher, f)

	if WriteTarArchive(w, sourceDir, uid, gid) != nil {
		return "", err
	}
	sha := hex.EncodeToString(hasher.Sum(make([]byte, 0, hasher.Size())))
	return "sha256:" + sha, nil
}

func WriteTarArchive(w io.Writer, srcDir string, uid, gid int) error {
	tw := tar.NewWriter(w)
	defer tw.Close()

	err := addParentDirs(srcDir, tw, uid, gid)
	if err != nil {
		return err
	}

	if runtime.GOOS == "windows" {

	}

	return filepath.Walk(srcDir, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fi.Mode()&os.ModeSocket != 0 {
			return nil
		}
		var header *tar.Header
		var target string
		if fi.Mode()&os.ModeSymlink != 0 {
			target, err = os.Readlink(file)
			if err != nil {
				return err
			}
		}
		header, err = tar.FileInfoHeader(fi, target)
		if err != nil {
			return err
		}
		if runtime.GOOS == "windows" {
			header.Name = path.Join("Files", filepath.ToSlash(file))
		} else {
			header.Name = file
		}
		header.ModTime = time.Date(1980, time.January, 1, 0, 0, 1, 0, time.UTC)
		header.Uid = uid
		header.Gid = gid
		header.Uname = ""
		header.Gname = ""

		if err := tw.WriteHeader(header); err != nil {
			return err
		}
		if fi.Mode().IsRegular() {
			f, err := os.Open(file)
			if err != nil {
				return err
			}
			defer f.Close()
			if _, err := io.Copy(tw, f); err != nil {
				return err
			}
		}
		return nil
	})
}

func addParentDirs(tarDir string, tw *tar.Writer, uid, gid int) error {
	parent := filepath.Dir(tarDir)
	if parent == "." || parent == "c:\\" || parent == "\\" || parent == "/" {
		if err := tw.WriteHeader(&tar.Header{
			Typeflag: tar.TypeDir,
			Name:     "Files",
			Mode:     0755,
			ModTime:  time.Date(1980, time.January, 1, 0, 0, 1, 0, time.UTC),
		}); err != nil {
			return err
		}
		return tw.WriteHeader(&tar.Header{
			Typeflag: tar.TypeDir,
			Name:     "Hives",
			Mode:     0755,
			ModTime:  time.Date(1980, time.January, 1, 0, 0, 1, 0, time.UTC),
		})
	}

	if err := addParentDirs(parent, tw, uid, gid); err != nil {
		return err
	}

	info, err := os.Stat(parent)
	if err != nil {
		return err
	}

	header, err := tar.FileInfoHeader(info, parent)
	if err != nil {
		return err
	}
	if runtime.GOOS == "windows" {
		header.Name = path.Join("Files", filepath.ToSlash(parent))
	} else {
		header.Name = parent
	}
	header.ModTime = time.Date(1980, time.January, 1, 0, 0, 1, 0, time.UTC)

	return tw.WriteHeader(header)
}

func Untar(r io.Reader, dest string) error {
	tr := tar.NewReader(r)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		path := filepath.Join(dest, hdr.Name)

		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(path, hdr.FileInfo().Mode()); err != nil {
				return err
			}
		case tar.TypeReg, tar.TypeRegA:
			_, err := os.Stat(filepath.Dir(path))
			if os.IsNotExist(err) {
				if err := os.MkdirAll(filepath.Dir(path), 0777); err != nil {
					return err
				}
			}
			if err := writeFile(tr, path, hdr.FileInfo().Mode()); err != nil {
				return err
			}
		case tar.TypeSymlink:
			if err := os.Symlink(hdr.Linkname, path); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown file type in tar %d", hdr.Typeflag)
		}
	}
}

func writeFile(in io.Reader, path string, mode os.FileMode) error {
	fh, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer fh.Close()
	_, err = io.Copy(fh, in)
	return err
}
