package core

// GetMsg loads the message variable
func (d *DuOS) SetScreen(s string) {
	d.Wv.Dispatch(func() {
		//d.Render("status", d.GetDuOSstatus())
		//d.Wv.Eval(`document.getElementById("main").innerHTML = "` + lib.VUEscreens()[s] + `";`)
	})
}
