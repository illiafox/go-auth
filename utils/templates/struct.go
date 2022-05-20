package templates

type Templates struct {
	Main    Template `tmpl:"main"`
	Message Template `tmpl:"msg"`
	Mail    Template `tmpl:"mail"`
}
