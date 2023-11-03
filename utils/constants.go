package utils

const TEMPLATE_DIR = "templates/"
const HTML_PATH_PREFIX = "/bca"

type LinksType struct {
	Budget     string
	Invoices   string
	Closure    string
	Actual     string
	Historic   string
	Supplier   string
	BudgetItem string
	Projects   string
	EditUser   string
	Password   string
}

var Links = &LinksType{
	Budget:     HTML_PATH_PREFIX + "/transacciones/presupuesto",
	Invoices:   HTML_PATH_PREFIX + "/transacciones/factura",
	Closure:    HTML_PATH_PREFIX + "/transacciones/cierre",
	Actual:     HTML_PATH_PREFIX + "/reportes/actual",
	Historic:   HTML_PATH_PREFIX + "/reportes/historico",
	Supplier:   HTML_PATH_PREFIX + "/parametros/proveedor/",
	BudgetItem: HTML_PATH_PREFIX + "/parametros/partidas",
	Projects:   HTML_PATH_PREFIX + "/parametros/proyectos",
	EditUser:   HTML_PATH_PREFIX + "/usuarios/",
	Password:   HTML_PATH_PREFIX + "/usuarios/contrasena",
}

var (
	BaseTemplate  string = TEMPLATE_DIR + "base.html"
	NavTemplate   string = TEMPLATE_DIR + "nav.html"
	TitleTemplate string = TEMPLATE_DIR + "title.html"

	RequiredFiles []string = []string{
		BaseTemplate,
		NavTemplate,
		TitleTemplate,
	}
)
