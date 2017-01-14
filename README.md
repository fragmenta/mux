# fragmenta mux 

Fragmenta mux is now the default router for fragmenta (though any router may be used). It provides sophisticated param parsing and routes prioritised in order of addition, without sacrificing speed. It also offers caching of frequent request->handler mappings for very fast lookups on most routes.

For a few benchmarks, see https://github.com/kennygrant/routebench