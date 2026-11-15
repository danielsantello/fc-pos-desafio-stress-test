package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sort"
)

func main() {
	var target string
	var requests, concurrency int

	flag.StringVar(&target, "url", "", "URL do serviço a ser testado")
	flag.IntVar(&requests, "requests", 0, "Número total de requisições")
	flag.IntVar(&concurrency, "concurrency", 0, "Número máximo de chamadas simultâneas")
	flag.Parse()

	parsedURL, err := url.Parse(target)

	var problem string
	switch {
	case err != nil || parsedURL.Hostname() == "" ||
		(parsedURL.Scheme != "http" && parsedURL.Scheme != "https"):
		problem = "--url deve ser uma URL HTTP ou HTTPS válida"
	case requests < 1:
		problem = "--requests deve ser maior que zero"
	case concurrency < 1:
		problem = "--concurrency deve ser maior que zero"
	case flag.NArg() != 0:
		problem = "use apenas as flags --url, --requests e --concurrency"
	}

	if problem != "" {
		fmt.Fprintln(os.Stderr, problem)
		flag.Usage()
		os.Exit(2)
	}

	result := runLoadTest(target, requests, concurrency)
	printReport(result)
}

func printReport(result loadResult) {
	fmt.Printf("Tempo total: %s\n", result.duration)
	fmt.Printf("Total de requisições realizadas: %d\n", result.total)
	fmt.Printf("HTTP 200: %d\n", result.statusCounts[http.StatusOK])

	statuses := make([]int, 0, len(result.statusCounts))
	for status := range result.statusCounts {
		if status != http.StatusOK {
			statuses = append(statuses, status)
		}
	}
	sort.Ints(statuses)

	fmt.Println("Outros status HTTP:")
	if len(statuses) == 0 {
		fmt.Println("  Nenhum")
	}
	for _, status := range statuses {
		fmt.Printf("  HTTP %d: %d\n", status, result.statusCounts[status])
	}

	fmt.Printf("Falhas sem resposta HTTP: %d\n", result.failures)
}
