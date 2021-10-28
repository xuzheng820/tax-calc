package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/xuri/excelize"
)

type Salary struct {
	Index    int
	Tax      string
	AfterTax string
	CoverTax string
}

type SalaryModel struct {
	walk.TableModelBase
	items               []*Salary
	city_average_salary []float64
	taxFree             []float64
	insuranceRate       float64
	gjjRate             float64
	totalIncome         float64
	totalTax            float64
	totalIncomeCoverTax float64
	totalTaxCoverTax    float64
}

func (m *SalaryModel) AddItem(month int, money string) {
	if money == "" {
		return
	}
	salary, err := walk.ParseFloat(money)
	if err != nil {
		fmt.Println("add item failed", err)
		return
	}
	var insurace float64
	if salary > m.city_average_salary[month-1] {
		insurace = m.city_average_salary[month-1] * (m.insuranceRate + m.gjjRate) / 100
	} else {
		insurace = salary * (m.insuranceRate + m.gjjRate) / 100
	}
	income := salary - 5000 - insurace
	if len(m.taxFree) >= month {
		income -= m.taxFree[month-1]
	}
	if income > 0 {
		m.totalIncome += income
		totalTax := calcTotalTax(m.totalIncome)
		tax := totalTax - m.totalTax
		m.totalTax = totalTax

		//cover tax
		//m.totalIncomeCoverTax += income
		//income2 := calcTotalCoverTax(tax, m.totalIncomeCoverTax)
		//m.totalIncomeCoverTax += income2
		//new_total_tax := calcTotalTax(m.totalIncomeCoverTax)
		//new_tax := new_total_tax - m.totalTaxCoverTax
		//m.totalTaxCoverTax = new_total_tax

		//fmt.Println("tax", tax, "income2", income2, "delta", income2-new_tax+tax)

		m.items = append(m.items, &Salary{
			Index:    month,
			Tax:      fmt.Sprintf("%.2f", tax),
			AfterTax: fmt.Sprintf("%.2f", salary-insurace-tax),
			CoverTax: "0",
		})
		//fmt.Println("month", month)
		//fmt.Println("tax", fmt.Sprintf("%.2f", tax))
		//fmt.Println("income", fmt.Sprintf("%.2f", salary-insurace-tax))
	} else {
		m.items[month-1] = &Salary{
			Index:    month,
			Tax:      "0",
			AfterTax: "0",
			CoverTax: "0",
		}
		//fmt.Println("month", month)
	}
}

func (m *SalaryModel) Reset() {
	m.city_average_salary = nil
	m.taxFree = nil
	m.insuranceRate = 0
	m.gjjRate = 0
	m.totalTax = 0
	m.totalIncome = 0
	m.items = make([]*Salary, 0)
}

func (m *SalaryModel) RowCount() int {
	return len(m.items)
}

func (m *SalaryModel) Value(row, col int) interface{} {
	item := m.items[row]

	switch col {
	case 0:
		return item.Index

	case 1:
		return item.Tax

	case 2:
		return item.AfterTax

	case 3:
		return item.CoverTax
	}
	return 0
}

func (m *SalaryModel) initCityAverageSalary(pay1, pay2, pay1_from, pay1_to, pay2_from, pay2_to string) error {
	m.city_average_salary = make([]float64, 0, 12)
	f_pay1, err := walk.ParseFloat(pay1)
	if err != nil {
		return err
	}
	from, er := strconv.Atoi(pay1_from)
	if er != nil {
		return er
	}
	to, e := strconv.Atoi(pay1_to)
	if e != nil {
		return e
	}
	//fmt.Println("city", f_pay1, from, to)
	for index := from; index <= to; index++ {
		if index < 1 || index > 12 {
			return errors.New("city average salary: month counter over 12")
		}
		if len(m.city_average_salary) >= index {
			m.city_average_salary[index-1] = f_pay1
			//fmt.Println("Index more", index, f_pay1, m.city_average_salary)
		} else {
			for len(m.city_average_salary) < index {
				m.city_average_salary = append(m.city_average_salary, 0)
			}
			m.city_average_salary[index-1] = f_pay1
			//fmt.Println("Index less", index, f_pay1, m.city_average_salary)
		}
	}
	if pay2 != "" {
		f_pay2, err := walk.ParseFloat(pay2)
		if err != nil {
			return err
		}
		from2, er := strconv.Atoi(pay2_from)
		if er != nil {
			return er
		}
		to2, e := strconv.Atoi(pay2_to)
		if e != nil {
			return e
		}
		for index := from2; index <= to2; index++ {
			if index < 1 || index > 12 {
				return errors.New("city average salary: month counter over 12")
			}
			if len(m.city_average_salary) >= index {
				m.city_average_salary[index-1] = f_pay2
			} else {
				for len(m.city_average_salary) < index {
					m.city_average_salary = append(m.city_average_salary, 0)
				}
				m.city_average_salary[index-1] = f_pay2
			}
		}
		//fmt.Println("city 2", f_pay2, from2, to2)
	}
	//fmt.Println("city", m.city_average_salary)

	return nil
}

func (m *SalaryModel) initRate(insurace, gjj string) error {
	f_insurace, err := walk.ParseFloat(insurace)
	if err != nil {
		return err
	}
	f_gjj, er := walk.ParseFloat(gjj)
	if er != nil {
		return er
	}
	m.insuranceRate = f_insurace
	m.gjjRate = f_gjj
	//fmt.Println(m.insuranceRate, m.gjjRate)
	//fmt.Println(m.gjjRate)
	return nil
}

func (m *SalaryModel) initTaxFree(free_money, free_from, free_to string) error {
	if m.taxFree == nil {
		m.taxFree = make([]float64, 12, 12)
	}
	if free_money == "" {
		return nil
	}
	from, err := strconv.Atoi(free_from)
	if err != nil {
		return err
	}
	to, er := strconv.Atoi(free_to)
	if er != nil {
		return er
	}
	money, e := walk.ParseFloat(free_money)
	if e != nil {
		return e
	}
	for index := from; index <= to; index++ {
		if index < 0 || index > 12 {
			return errors.New("tax free over 12 month")
		}
		m.taxFree[index-1] += money
	}
	//fmt.Println("tax free", m.taxFree)
	return nil
}

func NewSalaryModel() *SalaryModel {
	m := new(SalaryModel)
	return m
}

func calcTotalTax(money float64) float64 {
	if money <= 36000 {
		return money * 0.03
	} else if money <= 144000 {
		return money*0.1 - 2520
	} else if money <= 300000 {
		return money*0.2 - 16920
	} else if money <= 420000 {
		return money*0.25 - 31920
	} else if money <= 660000 {
		return money*0.3 - 52920
	} else if money <= 96000 {
		return money*0.35 - 85920
	} else {
		return money*0.45 - 181920
	}
}

func calcTotalCoverTax(tax, money float64) float64 {
	var income float64
	delta := tax
	switch {
	case money <= 36000:
		temp := (36000 - money) * (1 - 0.03)
		if temp > delta {
			return income + delta/(1-0.03)
		} else {
			income = 36000 - money
			delta = delta - temp
		}
		fallthrough
	case money <= 144000:
		temp := (144000 - money) * (1 - 0.1)
		if temp > delta {
			return income + delta/(1-0.1)
		} else {
			income = 144000 - money
			delta = delta - temp
		}
		fallthrough
		//return money*0.1 - 2520
	case money <= 300000:
		temp := (300000 - money) * (1 - 0.2)
		if temp > delta {
			return income + delta/(1-0.2)
		} else {
			income = 300000 - money
			delta = delta - temp
		}
		//return money*0.2 - 16920
		fallthrough
	case money <= 420000:
		temp := (420000 - money) * (1 - 0.25)
		if temp > delta {
			return income + delta/(1-0.25)
		} else {
			income = 420000 - money
			delta = delta - temp
		}
		//return money*0.25 - 31920
		fallthrough
	case money <= 660000:
		temp := (660000 - money) * (1 - 0.3)
		if temp > delta {
			return income + delta/(1-0.3)
		} else {
			income = 660000 - money
			delta = delta - temp
		}
		//return money*0.3 - 52920
		fallthrough
	case money <= 960000:
		temp := (960000 - money) * (1 - 0.35)
		if temp > delta {
			return income + delta/(1-0.35)
		} else {
			income = 960000 - money
			delta = delta - temp
		}
		//return money*0.35 - 85920
		fallthrough
	default:
		return income + delta/(1-0.45)
	}

}

type FileName struct {
	Name string
}

func RunFileNameDialy(owner walk.Form, file *FileName) (int, error) {
	var dlg *walk.Dialog
	var db *walk.DataBinder
	var acceptPB, cancelPB *walk.PushButton

	return Dialog{
		AssignTo:      &dlg,
		DefaultButton: &acceptPB,
		CancelButton: &cancelPB,
		DataBinder: DataBinder{
			AssignTo:       &db,
			Name:           "file",
			DataSource:     file,
			ErrorPresenter: ToolTipErrorPresenter{},
		},
		MinSize: Size{200, 100},
		Layout:  VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					Label{
						Text: "Name:",
					},
					LineEdit{
						Text: Bind("Name"),
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{
						AssignTo: &acceptPB,
						Text:     "OK",
						OnClicked: func() {
							if err := db.Submit(); err != nil {
								return
							}

							dlg.Accept()
						},
					},
					PushButton{
						AssignTo:  &cancelPB,
						Text:      "Cancel",
						OnClicked: func() { dlg.Cancel() },
					},
				},
			},
		},
	}.Run(owner)
}

func main() {
	var perMonthPay *walk.LineEdit
	var month1, month2, month3, month4, month5, month6, month7, month8, month9, month10, month11, month12 *walk.LineEdit
	var city_average_pay1, city_average_pay1_3, city_average_pay1_from, city_average_pay1_to *walk.LineEdit
	var city_average_pay2, city_average_pay2_3, city_average_pay2_from, city_average_pay2_to *walk.LineEdit
	var tax_discount_1, tax_discount_1_from, tax_discount_1_to *walk.LineEdit
	var tax_discount_2, tax_discount_2_from, tax_discount_2_to *walk.LineEdit
	var insurace_percent *walk.LineEdit
	var gongjijin_percent *walk.LineEdit

	model := NewSalaryModel()

	var mw *walk.MainWindow
	var result *walk.TableView

	MainWindow{
		AssignTo: &mw,
		Title:  "tax-calc",
		Size:   Size{1200, 800},
		Layout: HBox{},
		Children: []Widget{
			Composite{
				MaxSize: Size{600, 800},
				Layout:  VBox{},
				Children: []Widget{
					GroupBox{
						Layout: HBox{},
						Children: []Widget{
							Label{Text: "Salary Per Month"},
							LineEdit{AssignTo: &perMonthPay},
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
							HSpacer{
								Size: 300,
							},
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
							Label{Text: "city average salary"},
							LineEdit{AssignTo: &city_average_pay1},
							PushButton{
								Text: "*3",
								OnClicked: func() {
									salary, err := strconv.ParseFloat(city_average_pay1.Text(), 32)
									if err != nil {
										city_average_pay1.SetText("")
									}
									city_average_pay1_3.SetText(strconv.FormatFloat(salary*3, 'f', -1, 32))
								},
							},
							LineEdit{AssignTo: &city_average_pay1_3},
							Label{Text: "From"},
							LineEdit{AssignTo: &city_average_pay1_from, Text: "1"},
							Label{Text: "To"},
							LineEdit{AssignTo: &city_average_pay1_to, Text: "12"},
						},
					}, // city average salary 1
					GroupBox{
						Layout: HBox{},
						Children: []Widget{
							Label{Text: "city average salary 2"},
							LineEdit{AssignTo: &city_average_pay2},
							PushButton{
								Text: "*3",
								OnClicked: func() {
									salary, err := strconv.ParseFloat(city_average_pay2.Text(), 32)
									if err != nil {
										city_average_pay2.SetText("")
									}
									city_average_pay2_3.SetText(strconv.FormatFloat(salary*3, 'f', -1, 32))
								},
							},
							LineEdit{AssignTo: &city_average_pay2_3},
							Label{Text: "From"},
							LineEdit{AssignTo: &city_average_pay2_from, Text: "1"},
							Label{Text: "To"},
							LineEdit{AssignTo: &city_average_pay2_to, Text: "12"},
						},
					}, // city average salary 2
					GroupBox{
						Layout: HBox{},
						Children: []Widget{
							Label{Text: "tax free 1"},
							LineEdit{AssignTo: &tax_discount_1},
							Label{Text: "from"},
							LineEdit{AssignTo: &tax_discount_1_from, Text: "1"},
							Label{Text: "to"},
							LineEdit{AssignTo: &tax_discount_1_to, Text: "12"},
						},
					}, // tax discount 1
					GroupBox{
						Layout: HBox{},
						Children: []Widget{
							Label{Text: "tax free 2"},
							LineEdit{AssignTo: &tax_discount_2},
							Label{Text: "from"},
							LineEdit{AssignTo: &tax_discount_2_from, Text: "1"},
							Label{Text: "to"},
							LineEdit{AssignTo: &tax_discount_2_to, Text: "12"},
						},
					}, // tax discount 2
					GroupBox{
						Layout: HBox{},
						Children: []Widget{
							Label{Text: "insurance rate"},
							LineEdit{AssignTo: &insurace_percent},
							Label{Text: "%"},
							HSpacer{},
							Label{Text: "gongjijin rate"},
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
									fmt.Println("Calc Clicked")
									model.Reset()
									if err := model.initCityAverageSalary(city_average_pay1_3.Text(), city_average_pay2_3.Text(), city_average_pay1_from.Text(),
										city_average_pay1_to.Text(), city_average_pay2_from.Text(), city_average_pay2_to.Text()); err != nil {
										fmt.Println("init city failed", err)
										return
									}
									if err := model.initRate(insurace_percent.Text(), gongjijin_percent.Text()); err != nil {
										fmt.Println("init rate failed", err)
										return
									}
									fmt.Println("rate", insurace_percent.Text())
									if err := model.initTaxFree(tax_discount_1.Text(), tax_discount_1_from.Text(), tax_discount_1_to.Text()); err != nil {
										fmt.Println("init tax free failed", err)
										return
									}
									if err := model.initTaxFree(tax_discount_2.Text(), tax_discount_2_from.Text(), tax_discount_2_to.Text()); err != nil {
										fmt.Println("init tax free 2 failed", err)
										return
									}
									model.AddItem(1, month1.Text())
									model.AddItem(2, month2.Text())
									model.AddItem(3, month3.Text())
									model.AddItem(4, month4.Text())
									model.AddItem(5, month5.Text())
									model.AddItem(6, month6.Text())
									model.AddItem(7, month7.Text())
									model.AddItem(8, month8.Text())
									model.AddItem(9, month9.Text())
									model.AddItem(10, month10.Text())
									model.AddItem(11, month11.Text())
									model.AddItem(12, month12.Text())
									model.PublishRowsReset()

								},
							},
							PushButton{
								Text: "Cover Tax",
								OnClicked: func() {

								},
							},
						},
					}, // calc button
				},
			},
			Composite{
				Layout: VBox{},
				Children: []Widget{
					TableView{
						AssignTo: &result,
						Columns: []TableViewColumn{
							{Title: "month"},
							{Title: "tax"},
							{Title: "income"},
							{Title: "cover tax"},
						},
						Model: model,
					},
					PushButton{
						Text: "Save",
						OnClicked: func() {
							file := new(FileName)
							if cmd, err := RunFileNameDialy(mw, file); err != nil {
								fmt.Println(err)
								return
							} else if cmd == walk.DlgCmdOK {
								if file.Name == "" {
									file.Name = "salary"
								}
								f := excelize.NewFile()
								f.SetSheetName("Sheet1", "salary")
								index := f.GetSheetIndex("salary")
								f.SetActiveSheet(index)
								for line, item := range model.items {
									f.SetCellValue("salary", "A"+strconv.Itoa(line+1), item.Index)
									f.SetCellValue("salary", "B"+strconv.Itoa(line+1), item.Tax)
									f.SetCellValue("salary", "C"+strconv.Itoa(line+1), item.AfterTax)
									f.SetCellValue("salary", "D"+strconv.Itoa(line+1), item.CoverTax)
								}
								if err := f.SaveAs(file.Name +".xlsx"); err != nil {
									fmt.Println("save failed", err)
								}
							}
						},
					},
				},
			},
		},
	}.Run()
}
