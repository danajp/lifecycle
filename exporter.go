package lifecycle

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/buildpack/imgutil"
	"github.com/pkg/errors"

	"github.com/buildpack/lifecycle/archive"
	"github.com/buildpack/lifecycle/cmd"
	"github.com/buildpack/lifecycle/metadata"
)

type Exporter struct {
	Buildpacks   []Buildpack
	ArtifactsDir string
	In           []byte
	Logger       Logger
	UID, GID     int
}

type LauncherConfig struct {
	Path     string
	Metadata metadata.LauncherMetadata
}

func (e *Exporter) Export(
	layersDir,
	appDir string,
	workingImage imgutil.Image,
	runImageRef string,
	origMetadata metadata.LayersMetadata,
	additionalNames []string,
	launcherConfig LauncherConfig,
	stack metadata.StackMetadata,
) error {
	var err error

	meta := metadata.LayersMetadata{}
	meta.App = make(map[string]string)

	meta.RunImage.TopLayer, err = workingImage.TopLayer()
	if err != nil {
		return errors.Wrap(err, "get run image top layer SHA")
	}

	meta.RunImage.Reference = runImageRef
	meta.Stack = stack

	var buildMetadata BuildMetadata
	metadataTomlPath := filepath.Join(layersDir, "config", "metadata.toml")
	_, err = toml.DecodeFile(metadataTomlPath, &buildMetadata)
	if err != nil && !os.IsNotExist(err) {
		return errors.Wrap(err, "failed to read metadata.toml")
	}

	for _, slice := range buildMetadata.Slices {
		var allGlobMatches []string

		for _, path := range slice.Paths {
			// make all paths absolute to the app dir
			path = e.toAbs(appDir, path)
			globMatches, err := filepath.Glob(path)
			if err != nil {
				return errors.Wrap(err, "bad pattern for glob path")
			}
			allGlobMatches = append(allGlobMatches, globMatches...)
		}

		if len(allGlobMatches) > 0 {
			sliceLayerID := fmt.Sprintf("slice-%x", sha256.Sum256([]byte(strings.Join(slice.Paths, ""))))
			sliceSHA, err := e.addSliceLayer(workingImage, sliceLayerID, origMetadata.App[sliceLayerID], allGlobMatches)
			if err != nil {
				return errors.Wrap(err, "exporting slice layer")
			}
			e.Logger.Infof("Slice sha = %s", sliceSHA)
			meta.App[sliceLayerID] = sliceSHA
		}
	}

	appSHA, err := e.addLayer(workingImage, &layer{path: appDir, identifier: "app"}, origMetadata.App["app"])
	if err != nil {
		return errors.Wrap(err, "exporting app layer")
	}
	meta.App["app"] = appSHA

	meta.Config.SHA, err = e.addLayer(workingImage, &layer{path: filepath.Join(layersDir, "config"), identifier: "config"}, origMetadata.Config.SHA)
	if err != nil {
		return errors.Wrap(err, "exporting config layer")
	}

	meta.Launcher.SHA, err = e.addLayer(workingImage, &layer{path: launcherConfig.Path, identifier: "launcher"}, origMetadata.Launcher.SHA)
	if err != nil {
		return errors.Wrap(err, "exporting launcher layer")
	}

	for _, bp := range e.Buildpacks {
		bpDir, err := readBuildpackLayersDir(layersDir, bp)
		if err != nil {
			return errors.Wrapf(err, "reading layers for buildpack '%s'", bp.ID)
		}
		bpMD := metadata.BuildpackLayersMetadata{ID: bp.ID, Version: bp.Version, Layers: map[string]metadata.BuildpackLayerMetadata{}}

		layers := bpDir.findLayers(launch)
		for i, layer := range layers {
			lmd, err := layer.read()
			if err != nil {
				return errors.Wrapf(err, "reading '%s' metadata", layer.Identifier())
			}

			if layer.hasLocalContents() {
				origLayerMetadata := origMetadata.MetadataForBuildpack(bp.ID).Layers[layer.name()]
				lmd.SHA, err = e.addLayer(workingImage, &layers[i], origLayerMetadata.SHA)
				if err != nil {
					return err
				}
			} else {
				if lmd.Cache {
					return fmt.Errorf("layer '%s' is cache=true but has no contents", layer.Identifier())
				}
				origLayerMetadata, ok := origMetadata.MetadataForBuildpack(bp.ID).Layers[layer.name()]
				if !ok {
					return fmt.Errorf("cannot reuse '%s', previous image has no metadata for layer '%s'", layer.Identifier(), layer.Identifier())
				}

				e.Logger.Infof("Reusing layer '%s'\n", layer.Identifier())
				e.Logger.Debugf("Layer '%s' SHA: %s\n", layer.Identifier(), origLayerMetadata.SHA)
				if err := workingImage.ReuseLayer(origLayerMetadata.SHA); err != nil {
					return errors.Wrapf(err, "reusing layer: '%s'", layer.Identifier())
				}
				lmd.SHA = origLayerMetadata.SHA
			}
			bpMD.Layers[layer.name()] = lmd
		}

		if malformedLayers := bpDir.findLayers(malformed); len(malformedLayers) > 0 {
			ids := make([]string, 0, len(malformedLayers))
			for _, ml := range malformedLayers {
				ids = append(ids, ml.Identifier())
			}
			return fmt.Errorf("failed to parse metadata for layers '%s'", ids)
		}

		meta.Buildpacks = append(meta.Buildpacks, bpMD)
	}

	data, err := json.Marshal(meta)
	if err != nil {
		return errors.Wrap(err, "marshall metadata")
	}

	if err = workingImage.SetLabel(metadata.LayerMetadataLabel, string(data)); err != nil {
		return errors.Wrap(err, "set app image metadata label")
	}

	buildMD := &BuildMetadata{}
	if _, err := toml.DecodeFile(metadata.FilePath(layersDir), buildMD); err != nil {
		return errors.Wrap(err, "read build metadata")
	}

	if err := e.addBuildMetadataLabel(workingImage, buildMD.BOM, launcherConfig.Metadata); err != nil {
		return errors.Wrapf(err, "add build metadata label")
	}

	if err = workingImage.SetEnv(cmd.EnvLayersDir, layersDir); err != nil {
		return errors.Wrapf(err, "set app image env %s", cmd.EnvLayersDir)
	}

	if err = workingImage.SetEnv(cmd.EnvAppDir, appDir); err != nil {
		return errors.Wrapf(err, "set app image env %s", cmd.EnvAppDir)
	}

	if err = workingImage.SetEntrypoint(launcherConfig.Path); err != nil {
		return errors.Wrap(err, "setting entrypoint")
	}

	if err = workingImage.SetCmd(); err != nil { // Note: Command intentionally empty
		return errors.Wrap(err, "setting cmd")
	}

	return saveImage(workingImage, additionalNames, e.Logger)
}

func (e *Exporter) addLayer(image imgutil.Image, layer identifiableLayer, previousSHA string) (string, error) {
	tarPath := filepath.Join(e.ArtifactsDir, escapeID(layer.Identifier())+".tar")
	sha, err := archive.WriteTarFile(layer.Path(), tarPath, e.UID, e.GID)
	if err != nil {
		return "", errors.Wrapf(err, "exporting layer '%s'", layer.Identifier())
	}
	if sha == previousSHA {
		e.Logger.Infof("Reusing layer '%s'\n", layer.Identifier())
		e.Logger.Debugf("Layer '%s' SHA: %s\n", layer.Identifier(), sha)
		return sha, image.ReuseLayer(previousSHA)
	}
	e.Logger.Infof("Adding layer '%s'\n", layer.Identifier())
	e.Logger.Debugf("Layer '%s' SHA: %s\n", layer.Identifier(), sha)
	return sha, image.AddLayer(tarPath)
}

func (e *Exporter) addBuildMetadataLabel(image imgutil.Image, plan []BOMEntry, launcherMD metadata.LauncherMetadata) error {
	var bps []metadata.BuildpackMetadata
	for _, bp := range e.Buildpacks {
		bps = append(bps, metadata.BuildpackMetadata{
			ID:      bp.ID,
			Version: bp.Version,
		})
	}

	buildJSON, err := json.Marshal(metadata.BuildMetadata{
		BOM:        plan,
		Buildpacks: bps,
		Launcher:   launcherMD,
	})
	if err != nil {
		return errors.Wrap(err, "parse build metadata")
	}

	if err := image.SetLabel(metadata.BuildMetadataLabel, string(buildJSON)); err != nil {
		return errors.Wrap(err, "set build image metadata label")
	}

	return nil
}

func (e *Exporter) addSliceLayer(image imgutil.Image, layerID string, previousSHA string, files []string) (string, error) {
	tarPath := filepath.Join(e.ArtifactsDir, escapeID(layerID)+".tar")
	sha, fileSet, err := archive.WriteFilesToTar(tarPath, e.UID, e.GID, files...)
	if err != nil {
		return "", errors.Wrapf(err, "exporting slice layer '%s'", layerID)
	}

	for file, _ := range fileSet {
		stat, _ := os.Stat(file)
		if !stat.IsDir() {
			err = os.Remove(file)
			if err != nil {
				e.Logger.Errorf("failed to delete %v", err)
			}
		}
	}

	if sha == previousSHA {
		e.Logger.Infof("Reusing layer '%s'\n", layerID)
		e.Logger.Debugf("Layer '%s' SHA: %s\n", layerID, sha)
		return sha, image.ReuseLayer(previousSHA)
	}
	e.Logger.Infof("Adding layer '%s'\n", layerID)
	e.Logger.Debugf("Layer '%s' SHA: %s\n", layerID, sha)

	return sha, image.AddLayer(tarPath)
}

func (e *Exporter) toAbs(baseDir, path string) string {
	path = filepath.Clean(path)

	// force relative path to be absolute from the base dir
	if !filepath.IsAbs(path) {
		path = filepath.Join(baseDir, path)
	}
	// force an absolute path to be absolute from base dir
	if len(path) > len(baseDir) && path[:len(baseDir)] != baseDir {
		path = filepath.Join(baseDir, path)
		e.Logger.Warnf("found absolute path %s outside of %s", path, baseDir)
	}

	return path
}
