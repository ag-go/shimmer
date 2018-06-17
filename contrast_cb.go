// +build js,wasm

package shimmer

import (
	"encoding/base64"
	"fmt"
	"syscall/js"
	"time"

	"github.com/anthonynsimon/bild/adjust"
	"github.com/anthonynsimon/bild/imgio"
)

func (s *Shimmer) setupContrastCb() {
	s.contrastCb = js.NewEventCallback(js.PreventDefault, func(ev js.Value) {
		// quick return if no source image is yet uploaded
		if s.sourceImg == nil {
			return
		}
		delta := ev.Get("target").Get("value").Float()
		start := time.Now()
		res := adjust.Contrast(s.sourceImg, delta)

		s.buf.Reset()
		enc := imgio.JPEGEncoder(90)
		err := enc(&s.buf, res)
		if err != nil {
			s.log(err.Error())
			return
		}
		// Updating the image
		js.Global.Get("document").
			Call("getElementById", "targetImg").
			Set("src", jpegPrefix+base64.StdEncoding.EncodeToString(s.buf.Bytes()))
		fmt.Println("time taken for contrast:", time.Now().Sub(start))
	})
}
