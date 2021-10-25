package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func main() {
	var perMonthPay *walk.LineEdit
	var month1, month2, month3, month4, month5, month6, month7, month8, month9, month10, month11, month12 *walk.LineEdit
	var city_average_pay, city_average_pay_3 *walk.LineEdit
	var tax_discount_1, tax_discount_1_from, tax_discount_1_to *walk.LineEdit
	var tax_discount_2, tax_discount_2_from, tax_discount_2_to *walk.LineEdit
	var insurace_percent *walk.LineEdit
	var gongjijin_percent *walk.LineEdit
	result := TableView{
		Column: 4,
		Row:    13,
	}

	MainWindow{
		Title:   "tax-calc",
		MinSize: Size{400, 800},
		Layout:  VBox{},
		Children: []Widget{
			GroupBox{
				Layout: HBox{},
				Children: []Widget{
					Label{Text: "Per Month"},
					LineEdit{AssignTo: &perMonthPay, MaxLength: 30},
					PushButton{
						Text: "Sync",
						OnClicked: func() {
							month1.SetText(perMonthPay.Text())
							month2.SetText(perMonthPay.Text())
							month3.SetText(perMonthPay.Text())
							month4.SetText(perMonthPay.Text())
							month5.SetText(perMonthPay.Text())
							month6.SetText(perMonthPay.Text())
							month7.SetText(perMonthPay.Text())
							month8.SetText(perMonthPay.Text())
							month9.SetText(perMonthPay.Text())
							month10.SetText(perMonthPay.Text())
							month11.SetText(perMonthPay.Text())
							month12.SetText(perMonthPay.Text())

						},
					},
					HSpacer{},
				},
			},
			GroupBox{
				Layout: HBox{},
				Children: []Widget{
					Label{Text: "Jan"},
					LineEdit{AssignTo: &month1},
					Label{Text: "Feb"},
					LineEdit{AssignTo: &month2},
					Label{Text: "Mar"},
					LineEdit{AssignTo: &month3},
					Label{Text: "Apr"},
					LineEdit{AssignTo: &month4},
					Label{Text: "May"},
					LineEdit{AssignTo: &month5},
					Label{Text: "Jun"},
					LineEdit{AssignTo: &month6},
				},
			}, // 12month
			GroupBox{
				Layout: HBox{},
				Children: []Widget{
					Label{Text: "Jul"},
					LineEdit{AssignTo: &month7},
					Label{Text: "Aug"},
					LineEdit{AssignTo: &month8},
					Label{Text: "Sep"},
					LineEdit{AssignTo: &month9},
					Label{Text: "Oct"},
					LineEdit{AssignTo: &month10},
					Label{Text: "Nov"},
					LineEdit{AssignTo: &month11},
					Label{Text: "Dec"},
					LineEdit{AssignTo: &month12},
				},
			}, // 12month
			GroupBox{
				Layout: HBox{},
				Children: []Widget{
					Label{Text: "average pay"},
					LineEdit{AssignTo: &city_average_pay},
					Label{Text: "*3"},
					LineEdit{AssignTo: &city_average_pay_3},
				},
			}, // city average payment
			GroupBox{
				Layout: HBox{},
				Children: []Widget{
					Label{Text: "tax discount 1"},
					LineEdit{AssignTo: &tax_discount_1},
					Label{Text: "from"},
					LineEdit{AssignTo: &tax_discount_1_from},
					Label{Text: "to"},
					LineEdit{AssignTo: &tax_discount_1_to},
				},
			}, // tax discount 1
			GroupBox{
				Layout: HBox{},
				Children: []Widget{
					Label{Text: "tax discount 2"},
					LineEdit{AssignTo: &tax_discount_2},
					Label{Text: "from"},
					LineEdit{AssignTo: &tax_discount_2_from},
					Label{Text: "to"},
					LineEdit{AssignTo: &tax_discount_2_to},
				},
			}, // tax discount 2
			GroupBox{
				Layout: HBox{},
				Children: []Widget{
					Label{Text: "insurance per"},
					LineEdit{AssignTo: &insurace_percent},
					Label{Text: "%"},
					Label{Text: "gongjijin per"},
					LineEdit{AssignTo: &gongjijin_percent},
					Label{Text: "%"},
				},
			}, // insurance percent & gongjijin percent
			GroupBox{
				Layout: HBox{},
				Children: []Widget{
					PushButton{
						Text: "Calc",
						OnClicked: func() {

						},
					},
					PushButton{
						Text: "Cover Tax",
						OnClicked: func() {

						},
					},
				},
			}, // calc button
			result,
		},
	}.Run()
}
