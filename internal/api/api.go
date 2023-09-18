package api

import (
	"github.com/go-chi/chi/v5"
)

func Api() *chi.Mux {
	mux := chi.NewRouter()
	// 一般接口
	mux.Route("/", func(r chi.Router) {
		r.Get("", nil)
	})
	mux.Route("/login/{id}", func(r chi.Router) {
		r.Get("", nil)
		r.Post("", nil)
	})
	mux.Route("/menu", func(r chi.Router) {
		r.Post("", nil)
	})
	mux.Route("/partner/{id}", func(r chi.Router) {
		r.Get("", nil)
		r.Post("", nil)
		r.Put("", nil)
		r.Delete("", nil)
	})

	// 地址管理
	mux.Route("/address", func(r chi.Router) {
		r.Route("/country/{id}", func(sr chi.Router) {
			sr.Get("", nil)
			sr.Post("", nil)
			sr.Put("", nil)
		})
		r.Route("/province/{id}", func(sr chi.Router) {
			sr.Get("", nil)
			sr.Post("", nil)
			sr.Put("", nil)
		})
		r.Route("/city/{id}", func(sr chi.Router) {
			sr.Get("", nil)
			sr.Post("", nil)
			sr.Put("", nil)
			sr.Delete("", nil)
		})
		r.Route("/district/{id}", func(sr chi.Router) {
			sr.Get("", nil)
			sr.Put("", nil)
		})
	})
	// 设置管理
	mux.Route("/setting", func(r chi.Router) {
		r.Route("/user/{id}", func(sr chi.Router) {
			sr.Get("", nil)
		})
		r.Route("/group/{id}", func(sr chi.Router) {
			sr.Get("", nil)
		})
	})
	// 产品管理
	mux.Route("/product", func(r chi.Router) {
		r.Route("/attribute/{id}", func(sr chi.Router) {
			sr.Get("", nil)
			sr.Post("", nil)
			sr.Put("", nil)
		})
		r.Route("/attribute/line/{id}", func(sr chi.Router) {
			sr.Get("", nil)
			sr.Post("", nil)
			sr.Put("", nil)
		})
		r.Route("/attributevalue/{id}", func(sr chi.Router) {
			sr.Get("", nil)
			sr.Post("", nil)
			sr.Put("", nil)
			sr.Delete("", nil)
		})
		r.Route("/template/{id}", func(sr chi.Router) {
			sr.Get("", nil)
			sr.Post("", nil)
			sr.Put("", nil)
		})
		r.Route("/product/{id}", func(sr chi.Router) {
			sr.Get("", nil)
			sr.Post("", nil)
			sr.Put("", nil)
		})
		r.Route("/uom/{id}", func(sr chi.Router) {
			sr.Get("", nil)
			sr.Post("", nil)
			sr.Put("", nil)
		})
		r.Route("/uomcateg/{id}", func(sr chi.Router) {
			sr.Get("", nil)
			sr.Post("", nil)
			sr.Put("", nil)
		})
		r.Route("/category/{id}", func(sr chi.Router) {
			sr.Get("", nil)
			sr.Post("", nil)
			sr.Put("", nil)
		})
	})
	// 销售管理
	mux.Route("/sale", func(r chi.Router) {
		r.Route("/order/{id}", func(sr chi.Router) {
			sr.Get("", nil)
			sr.Post("", nil)
			sr.Put("", nil)
			sr.Delete("", nil)
		})
		r.Route("/order/line/{id}", func(sr chi.Router) {
			sr.Get("", nil)
			sr.Post("", nil)
			sr.Put("", nil)
			sr.Delete("", nil)
		})
	})
	return mux
}
