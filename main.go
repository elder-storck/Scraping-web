/*
  Descrição: Programa em Go para realizar scraping em uma página web de compras usando Colly.
  Autor: Elder Ribeiro Storck
  Email: elder.storck@gmail.com
  Data de Criação: 07 de Dezembro de 2023
  Última Atualização: 07 de Dezembro de 2023
*/

package main

import (
	"encoding/csv"
	"github.com/gocolly/colly"
	"log"
	"os"
)

// definindo a estrtutra de dado para armazenar os dados scrapeados
type T_itensProd struct {
	url, image, name, price string
}

// retorna true se a string estiver presente no slice
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func main() {

	var produtos []T_itensProd // inicializa o slice de structs que conterá os dados scrapeados

	var paginasToScrape []string // Iniicializa o slice de páginas a serem "raspadas" com um slice vazio

	paginaToScrape := "https://www.extrabom.com.br/c/bebidas/28/?limit=60" // A primeira URL a ser raspada

	paginasDescobertas := []string{paginaToScrape} // Inicializando o slice de páginas descobertas com a primeira página a ser raspada

	i := 1
	limit := 5

	c := colly.NewCollector() // Iniicializando a intancia Colly

	/*Definindo um User-Agent valido*/
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"

	c.OnHTML("a.page-numbers", func(e *colly.HTMLElement) {
		linkNovaPag := e.Attr("href") //descobrindo nova pagina
		if !contains(paginasToScrape, linkNovaPag) {
			if !contains(paginasDescobertas, linkNovaPag) {
				paginasToScrape = append(paginasToScrape, linkNovaPag)
			}
			paginasDescobertas = append(paginasDescobertas, linkNovaPag)
		}
	})

	// Scraping os dados do produto
	c.OnHTML(".carousel__item", func(e *colly.HTMLElement) {
		produto := T_itensProd{}
		produto.url = e.ChildAttr("a", "href")
		produto.image = e.ChildAttr("img", "src")
		produto.name = e.ChildText(".name-produto")
		produto.price = e.ChildText(".item-por")

		produtos = append(produtos, produto)
	})

	c.OnScraped(func(response *colly.Response) {
		if len(paginasToScrape) != 0 && i < limit {
			paginaToScrape = paginasToScrape[0]
			paginasToScrape = paginasToScrape[1:]
			i++
			c.Visit(paginaToScrape)
		}
	})

	c.Visit(paginaToScrape) // Visitando a primeira pag

	file, err := os.Create("products.csv") //Abrindo arquivo
	if err != nil {
		log.Fatalln("Erro ao criar arquivo CSV", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file) //iniciando escrita do arquivo

	headers := []string{
		"url",
		"imagem",
		"nome",
		"preço",
	}
	writer.Write(headers) //escrevendo Headers

	/*Escrevendo os dados dos produtos*/
	for _, produto := range produtos {

		record := []string{
			produto.url, produto.image, produto.name, produto.price, // Convertendo o produto para array de strings
		}

		writer.Write(record)
	}
	for _, produto := range produtos {
		record := []string{
			produto.url,
			produto.image,
			produto.name,
			produto.price,
		}
		println(record[2])
		println(record[3])
	}

	defer writer.Flush()
}
