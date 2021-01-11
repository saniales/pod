package gui

import (
	l "gioui.org/layout"
	"gioui.org/text"
	
	"github.com/p9c/pod/pkg/gui"
)

func (wg *WalletGUI) ReceivePage() l.Widget {
	return func(gtx l.Context) l.Dimensions {
		if wg.State != nil {
			// Debug(wg.State.isAddress)
			// Debug(wg.State.isAddress.Load())
			if wg.State.isAddress.Load() {
				ad := wg.State.currentReceivingAddress.Load()
				wg.currentReceiveAddress = ad.EncodeAddress()
				// var err error
				// // Debug(ad.ScriptAddress())
				// var conv []byte
				// if conv, err = bech32.ConvertBits(ad.ScriptAddress(), 8, 5, true); Check(err) {
				// }
				// if bech, err = bech32.Encode("pc", conv); Check(err) {
				// }
			}
		}
		return wg.Responsive(*wg.Size, gui.Widgets{
			{
				Size: 1280,
				Widget:
				wg.Fill("PanelBg", l.W, wg.TextSize.V, l.Center,
					wg.VFlex().AlignMiddle().
						Flexed(1, gui.EmptyMaxWidth()).
						Rigid(
							wg.Flex().AlignMiddle().
								Rigid(
									wg.VFlex().AlignMiddle().
										Rigid(
											wg.Inset(0.25,
												wg.Body2("Scan to send or click to copy").Alignment(text.Middle).Fn,
											).Fn,
										).
										Rigid(
											wg.currentReceiveQR,
										).
										Rigid(
											wg.Inset(0.25,
												wg.Caption(wg.currentReceiveAddress).Font("go regular").Fn,
											).Fn,
										).Fn,
								).
								
								Rigid(
									wg.VFlex().AlignMiddle().
										Rigid(
											wg.Inset(0.25,
												func(gtx l.Context) l.
												Dimensions {
													gtx.Constraints.Max.X = int(wg.TextSize.V * 29)
													return wg.inputs["receiveAmount"].Fn(gtx)
												},
											).Fn,
										).
										Rigid(
											wg.Inset(0.25,
												func(gtx l.Context) l.Dimensions {
													gtx.Constraints.Max.X = int(wg.TextSize.V * 29)
													return wg.inputs["receiveMessage"].Fn(gtx)
												},
											).Fn,
										).
										Rigid(
											wg.Inset(0.25,
												func(gtx l.Context) l.Dimensions {
													gtx.Constraints.Max.X = int(wg.TextSize.V * 29)
													return wg.ButtonLayout(wg.currentReceiveRegenClickable.SetClick(func() {
														Debug("clicked regenerate button")
														wg.currentReceiveGetNew.Store(true)
													})).Background("Primary").
														Embed(
															wg.Inset(0.5,
																wg.H6("regenerate").Color("Light").Fn,
															).Fn,
														).
														Fn(gtx)
												}).
												Fn,
										).Fn,
								).
								Fn,
						).
						Flexed(1, gui.EmptyMaxWidth()).
						Fn,
					// l.W, wg.TextSize.V).Fn,
				).
					Fn,
			},
			{
				Size: 0,
				Widget:
				wg.Fill("PanelBg", l.W, wg.TextSize.V, l.Center,
					wg.VFlex().AlignMiddle().
						Flexed(1, gui.EmptyMaxWidth()).
						Rigid(
							wg.VFlex().AlignMiddle().
								Rigid(
									wg.Inset(0.25,
										wg.Body2("Scan to send or click to copy").Alignment(text.Middle).Fn,
									).Fn,
								).
								Rigid(
									wg.currentReceiveQR,
								).
								Rigid(
									wg.Inset(0.25,
										wg.Caption(wg.currentReceiveAddress).Font("go regular").Fn,
									).Fn,
								).
								Rigid(
									wg.Inset(0.25,
										func(gtx l.Context) l.
										Dimensions {
											gtx.Constraints.Max.X = int(wg.TextSize.V * 29)
											return wg.inputs["receiveAmount"].Fn(gtx)
										},
									).Fn,
								).
								Rigid(
									wg.Inset(0.25,
										func(gtx l.Context) l.Dimensions {
											gtx.Constraints.Max.X = int(wg.TextSize.V * 29)
											return wg.inputs["receiveMessage"].Fn(gtx)
										},
									).Fn,
								).
								Rigid(
									wg.Inset(0.25,
										func(gtx l.Context) l.Dimensions {
											gtx.Constraints.Max.X = int(wg.TextSize.V * 29)
											return wg.ButtonLayout(wg.currentReceiveRegenClickable.SetClick(func() {
												Debug("clicked regenerate button")
												wg.currentReceiveGetNew.Store(true)
											})).Background("Primary").
												Embed(
													wg.Inset(0.5,
														wg.H6("regenerate").Color("Light").Fn,
													).Fn,
												).
												Fn(gtx)
										}).
										Fn,
								).
								Fn,
						).
						Flexed(1, gui.EmptyMaxWidth()).
						Fn,
					// l.W, wg.TextSize.V).Fn,
				).
					Fn,
			},
		}).
			Fn(gtx)
	}
}
