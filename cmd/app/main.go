package main

import "github.com/levchenki/tea-api/internal/app"

//	@title			Tea API
//	@version		1.0
//	@description	This is a Tea API for tea cafe.

//	@contact.name	Danila Levchenko

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization

// @query.collection.format	multi
func main() {
	app.Run()
}
