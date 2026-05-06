package gobs

/*
#include <obs/obs.h>
#include <obs/graphics/vec2.h>

extern bool scene_items_cb(obs_scene_t*, obs_sceneitem_t*, void*);
*/
import "C"

import (
	"context"
	"iter"
	"runtime/cgo"
	"unsafe"
)

func SceneCreate(name string) Scene {
	// #cgo noescape obs_scene_create
	// #cgo nocallback obs_scene_create
	scene := Scene{C.obs_scene_create(fromString(name).cptr())}
	if scene.c == nil {
		panic("failed to create scene")
	}
	return scene
}

func SceneFromSource(s Source) Scene {
	// #cgo noescape obs_scene_from_source
	// #cgo nocallback obs_scene_from_source
	return Scene{C.obs_scene_from_source(s.c)}
}

type Scene struct {
	c *C.obs_scene_t
}

func (s Scene) Release() {
	// #cgo noescape obs_scene_release
	// #cgo nocallback obs_scene_release
	C.obs_scene_release(s.c)
}

func (s Scene) Add(source Source) SceneItem {
	// #cgo noescape obs_scene_add
	// #cgo nocallback obs_scene_add
	return SceneItem{C.obs_scene_add(s.c, source.c)}
}

func (s Scene) Source() Source {
	// #cgo noescape obs_scene_get_source
	// #cgo nocallback obs_scene_get_source
	return Source{C.obs_scene_get_source(s.c)}
}

func (s Scene) FindSource(name string, recursive bool) SceneItem {
	if recursive {
		// #cgo noescape obs_scene_find_source_recursive
		// #cgo nocallback obs_scene_find_source_recursive
		return SceneItem{C.obs_scene_find_source_recursive(s.c, fromString(name).cptr())}
	}
	// #cgo noescape obs_scene_find_source
	// #cgo nocallback obs_scene_find_source
	return SceneItem{C.obs_scene_find_source(s.c, fromString(name).cptr())}
}

//export scene_items_cb
func scene_items_cb(_ *C.obs_scene_t, item *C.obs_sceneitem_t, p unsafe.Pointer) C.bool {
	h := (*cgo.Handle)(p)
	yield := h.Value().(func(SceneItem) bool)

	return C.bool(yield(SceneItem{item}))
}

func (s Scene) Items() iter.Seq[SceneItem] {
	return func(yield func(SceneItem) bool) {
		h := cgo.NewHandle(yield)
		defer h.Delete()

		// #cgo noescape obs_scene_enum_items
		C.obs_scene_enum_items(s.c, (*[0]byte)(C.scene_items_cb), unsafe.Pointer(&h))
	}
}

func (s Scene) ReorderItems(items []SceneItem) bool {
	// #cgo noescape obs_scene_reorder_items
	// #cgo nocallback obs_scene_reorder_items
	return bool(C.obs_scene_reorder_items(s.c, (**C.obs_sceneitem_t)(unsafe.Pointer(unsafe.SliceData(items))), C.size_t(len(items))))
}

type SceneDuplicateType uint32

const (
	SceneDupRefs        = SceneDuplicateType(C.OBS_SCENE_DUP_REFS)         // Duplicates the scene, but scene items are only duplicated with references
	SceneDupCopy        = SceneDuplicateType(C.OBS_SCENE_DUP_COPY)         // Duplicates the scene, and scene items are also fully duplicated when possible
	SceneDupPrivateRefs = SceneDuplicateType(C.OBS_SCENE_DUP_PRIVATE_REFS) // Duplicates with references, but the scene is a private source
	SceneDupPrivateCopy = SceneDuplicateType(C.OBS_SCENE_DUP_PRIVATE_COPY) // Fully duplicates scene items when possible, but the scene and duplicates sources are private sources
)

func (s Scene) Duplicate(name string, duplicateType SceneDuplicateType) Scene {
	return Scene{C.obs_scene_duplicate(s.c, fromString(name).cptr(), C.enum_obs_scene_duplicate_type(duplicateType))}
}

func (s Scene) Events(ctx context.Context) <-chan SignalScene {
	ch := make(chan SignalScene, 8)

	sh := GlobalSignalHandler()
	connections := []SignalConnection{
		sh.Connect("item_add", func(data any, cd CallData) {
			data.(chan SignalScene) <- SignalSceneItemAdd{
				Scene: Scene{(*C.obs_scene_t)(cd.Ptr("scene"))},
				Item:  SceneItem{(*C.obs_sceneitem_t)(cd.Ptr("item"))},
			}
		}, ch),
		sh.Connect("item_remove", func(data any, cd CallData) {
			data.(chan SignalScene) <- SignalSceneItemRemove{
				Scene: Scene{(*C.obs_scene_t)(cd.Ptr("scene"))},
				Item:  SceneItem{(*C.obs_sceneitem_t)(cd.Ptr("item"))},
			}
		}, ch),
		sh.Connect("reorder", func(data any, cd CallData) {
			data.(chan SignalScene) <- SignalSceneReorder{
				Scene: Scene{(*C.obs_scene_t)(cd.Ptr("scene"))},
			}
		}, ch),
		sh.Connect("refresh", func(data any, cd CallData) {
			data.(chan SignalScene) <- SignalSceneRefresh{
				Scene: Scene{(*C.obs_scene_t)(cd.Ptr("scene"))},
			}
		}, ch),
		sh.Connect("item_visible", func(data any, cd CallData) {
			data.(chan SignalScene) <- SignalSceneItemVisible{
				Scene:   Scene{(*C.obs_scene_t)(cd.Ptr("scene"))},
				Item:    SceneItem{(*C.obs_sceneitem_t)(cd.Ptr("item"))},
				Visible: cd.Bool("visible"),
			}
		}, ch),
		sh.Connect("item_locked", func(data any, cd CallData) {
			data.(chan SignalScene) <- SignalSceneItemLocked{
				Scene:  Scene{(*C.obs_scene_t)(cd.Ptr("scene"))},
				Item:   SceneItem{(*C.obs_sceneitem_t)(cd.Ptr("item"))},
				Locked: cd.Bool("locked"),
			}
		}, ch),
		sh.Connect("item_select", func(data any, cd CallData) {
			data.(chan SignalScene) <- SignalSceneItemSelect{
				Scene: Scene{(*C.obs_scene_t)(cd.Ptr("scene"))},
				Item:  SceneItem{(*C.obs_sceneitem_t)(cd.Ptr("item"))},
			}
		}, ch),
		sh.Connect("item_deselect", func(data any, cd CallData) {
			data.(chan SignalScene) <- SignalSceneItemDeselect{
				Scene: Scene{(*C.obs_scene_t)(cd.Ptr("scene"))},
				Item:  SceneItem{(*C.obs_sceneitem_t)(cd.Ptr("item"))},
			}
		}, ch),
		sh.Connect("item_transform", func(data any, cd CallData) {
			data.(chan SignalScene) <- SignalSceneItemTransform{
				Scene: Scene{(*C.obs_scene_t)(cd.Ptr("scene"))},
				Item:  SceneItem{(*C.obs_sceneitem_t)(cd.Ptr("item"))},
			}
		}, ch),
	}

	go func() {
		<-ctx.Done()
		for _, connection := range connections {
			connection.Disconnect()
		}
		close(ch)
	}()

	return ch
}

type SignalScene interface {
	Signal
	isSignalScene()
}

type SignalSceneItemAdd struct {
	Scene Scene
	Item  SceneItem
}

func (SignalSceneItemAdd) isSignal()      {}
func (SignalSceneItemAdd) isSignalScene() {}

type SignalSceneItemRemove struct {
	Scene Scene
	Item  SceneItem
}

func (SignalSceneItemRemove) isSignal()      {}
func (SignalSceneItemRemove) isSignalScene() {}

type SignalSceneReorder struct {
	Scene Scene
}

func (SignalSceneReorder) isSignal()      {}
func (SignalSceneReorder) isSignalScene() {}

type SignalSceneRefresh struct {
	Scene Scene
}

func (SignalSceneRefresh) isSignal()      {}
func (SignalSceneRefresh) isSignalScene() {}

type SignalSceneItemVisible struct {
	Scene   Scene
	Item    SceneItem
	Visible bool
}

func (SignalSceneItemVisible) isSignal()      {}
func (SignalSceneItemVisible) isSignalScene() {}

type SignalSceneItemLocked struct {
	Scene  Scene
	Item   SceneItem
	Locked bool
}

func (SignalSceneItemLocked) isSignal()      {}
func (SignalSceneItemLocked) isSignalScene() {}

type SignalSceneItemSelect struct {
	Scene Scene
	Item  SceneItem
}

func (SignalSceneItemSelect) isSignal()      {}
func (SignalSceneItemSelect) isSignalScene() {}

type SignalSceneItemDeselect struct {
	Scene Scene
	Item  SceneItem
}

func (SignalSceneItemDeselect) isSignal()      {}
func (SignalSceneItemDeselect) isSignalScene() {}

type SignalSceneItemTransform struct {
	Scene Scene
	Item  SceneItem
}

func (SignalSceneItemTransform) isSignal()      {}
func (SignalSceneItemTransform) isSignalScene() {}

type Vec2 struct {
	X, Y float32
}

type BoundsType uint32

const (
	BOUNDS_NONE            BoundsType = iota // No bounding box
	BOUNDS_STRETCH                           // Stretch to the bounding box without preserving aspect ratio
	BOUNDS_SCALE_INNER                       // Scales with aspect ratio to inner bounding box rectangle
	BOUNDS_SCALE_OUTER                       // Scales with aspect ratio to outer bounding box rectangle
	BOUNDS_SCALE_TO_WIDTH                    // Scales with aspect ratio to the bounding box width
	BOUNDS_SCALE_TO_HEIGHT                   // Scales with aspect ratio to the bounding box height
	BOUNDS_MAX_ONLY                          // Scales with aspect ratio, but only to the size of the source maximum
)

type SceneItem struct {
	c *C.obs_sceneitem_t
}

func (si SceneItem) Release() {
	// #cgo noescape obs_sceneitem_release
	// #cgo nocallback obs_sceneitem_release
	C.obs_sceneitem_release(si.c)
}

func (si SceneItem) Remove() {
	// #cgo noescape obs_sceneitem_remove
	// #cgo nocallback obs_sceneitem_remove
	C.obs_sceneitem_remove(si.c)
}

func (si SceneItem) BoundsType() BoundsType {
	// #cgo noescape obs_sceneitem_get_bounds_type
	// #cgo nocallback obs_sceneitem_get_bounds_type
	return BoundsType(C.obs_sceneitem_get_bounds_type(si.c))
}

func (si SceneItem) SetBoundsType(tipe BoundsType) {
	// #cgo noescape obs_sceneitem_set_bounds_type
	// #cgo nocallback obs_sceneitem_set_bounds_type
	C.obs_sceneitem_set_bounds_type(si.c, uint32(tipe))
}

func (si SceneItem) Bounds() (v Vec2) {
	// #cgo noescape obs_sceneitem_get_bounds
	// #cgo nocallback obs_sceneitem_get_bounds
	C.obs_sceneitem_get_bounds(si.c, (*C.struct_vec2)(unsafe.Pointer(&v)))
	return
}

func (si SceneItem) SetBounds(v Vec2) {
	// #cgo noescape obs_sceneitem_set_bounds
	// #cgo nocallback obs_sceneitem_set_bounds
	C.obs_sceneitem_set_bounds(si.c, (*C.struct_vec2)(unsafe.Pointer(&v)))
}
