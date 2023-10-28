package utils

const TEMPLATE_DIR = "newtemplate/"
const HTML_PATH_PREFIX = "/bca"

type LinksType struct {
	Budget     string
	Invoices   string
	Closure    string
	Actual     string
	Historic   string
	Supplier   string
	BudgetItem string
	Proyects   string
	EditUser   string
	Password   string
}

var Links = &LinksType{
	Budget:     HTML_PATH_PREFIX + "/budget",
	Invoices:   HTML_PATH_PREFIX + "/invoices",
	Closure:    HTML_PATH_PREFIX + "/closure",
	Actual:     HTML_PATH_PREFIX + "/actual",
	Historic:   HTML_PATH_PREFIX + "/historic",
	Supplier:   HTML_PATH_PREFIX + "/supplier",
	BudgetItem: HTML_PATH_PREFIX + "/budget-item",
	Proyects:   HTML_PATH_PREFIX + "/proyects",
	EditUser:   HTML_PATH_PREFIX + "/edit-user",
	Password:   HTML_PATH_PREFIX + "/password",
}
