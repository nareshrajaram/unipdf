/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package optimize

import (
	"github.com/unidoc/unidoc/pdf/core"
)

// CompressStreams compresses uncompressed streams.
// It implements interface model.Optimizer.
type CompressStreams struct {
}

// Optimize optimizes PDF objects to decrease PDF size.
func (c *CompressStreams) Optimize(objects []core.PdfObject) (optimizedObjects []core.PdfObject, err error) {
	optimizedObjects = make([]core.PdfObject, len(objects))
	copy(optimizedObjects, objects)
	for _, obj := range objects {
		stream, isStreamObj := core.GetStream(obj)
		if !isStreamObj {
			continue
		}
		if _, found := core.GetName(stream.PdfObjectDictionary.Get("Filter")); found {
			continue
		}
		encoder := core.NewLZWEncoder()
		encoder.EarlyChange = 0
		var data []byte
		data, err = encoder.EncodeBytes(stream.Stream)
		if err != nil {
			return optimizedObjects, err
		}
		dict := encoder.MakeStreamDict()
		// compare compressed and uncompressed sizes
		if len(data)+len(dict.DefaultWriteString()) < len(stream.Stream) {
			stream.Stream = data
			stream.PdfObjectDictionary.Merge(dict)
			stream.PdfObjectDictionary.Set("Length", core.MakeInteger(int64(len(stream.Stream))))
		}
	}
	return optimizedObjects, nil
}
