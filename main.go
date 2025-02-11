package main

import (
    "crypto/tls"
    "flag"
    "fmt"
    "log"
    "net"
    "net/http"
    "os"
    "strings"
    "sync"
    "time"
)

// Версия программы
const version = "0.0.1"

func sendRequestWithSourceIP(url string, srcIP string, src_port int, page string) error {
    // Создаем транспорт с указанием локального адреса
    transport := &http.Transport{
        DialTLS: func(network, addr string) (net.Conn, error) {
            // Генерируем случайный SNI
            // sni := fmt.Sprintf("nginx%d", src_port)

            // Настроим TLS
            config := &tls.Config{
                InsecureSkipVerify: true,  
                MinVersion: tls.VersionTLS12,
                MaxVersion: tls.VersionTLS12,
                // ServerName: sni, // Устанавливаем SNI
            }

            if network != "tcp" {
                panic("unsupported network " + network)
            }

            // Выполняем TLS подключение с полным адресом и портом
            conn, err := tls.Dial(network, url, config)
            if err != nil {
                return nil, err
            }
            return conn, nil
        },

    }    

    client := &http.Client{
        Transport: transport,
        Timeout:   5 * time.Second,
    }    
    
    // Создаем HTTP GET запрос
    url2 := fmt.Sprintf("https://%s/%s", url, page)
	req, err := http.NewRequest("GET", url2, nil)
	if err != nil {
		return fmt.Errorf("Ошибка создания запроса:", err)
	}
    host := fmt.Sprintf("nginx%d", src_port)
    req.Host = host
    req.Header.Set("Connection", "close")

    // Отправляем запрос
    resp, err := client.Do(req)
    if err != nil {
        return fmt.Errorf("failed to send request: %v", err)
    }
    defer resp.Body.Close()   
    log.Printf("Response from %s with source IP %s:%d %s\n", url, srcIP, src_port, resp.Status)
    return nil
}

func main() {
    // declaring vars
    var url_str string
    var ips string
    var page_size string
    var cpsInt int
    var workers int
    var logfie string
    var showVersion bool

    // parsing
    flag.StringVar(&url_str, "target", "nginx01:443", "Введите target")
    flag.StringVar(&ips, "ips", "192.168.0.100,192.168.0.101", "Список IP-адресов, разделённых запятой")
    flag.StringVar(&page_size, "uri", "1kb.html", "Введите URI")
    flag.IntVar(&cpsInt, "cps", 10000, "Введите CPS")
    flag.IntVar(&workers, "wrk", 20000, "Введите кол-во workers")
    flag.StringVar(&logfie, "log", "default.log", "Введите название log файла")
    flag.BoolVar(&showVersion, "version", false, "Показать версию программы")
    flag.Parse()


    // Если указан флаг -version, выводим версию и завершаем работу
    if showVersion {
        fmt.Printf("Версия программы: %s\n", version)
        return
    }

    logFile, err := os.OpenFile(logfie, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
		log.Fatalf("Ошибка открытия файла для логирования: %v", err)
	}
    defer logFile.Close()
	log.SetOutput(logFile) // Устанавливаем логгер, чтобы писать в файл

    // base calculations
    srcIPs := strings.Split(ips, ",")
    var cps float32 = float32(cpsInt)
    concurrentLimit := workers
    limiter := make(chan struct{}, concurrentLimit)
    var microsecond float32 = 1000000
    var cps_sec_ratio float32 =  (1 / cps)
    duration := time.Duration(cps_sec_ratio * (microsecond)) * time.Microsecond
    port_strt := 1
    port_end := 65535
    

    var wg sync.WaitGroup // Используем WaitGroup для ожидания завершения всех горутин
    
    for _, src_ip := range srcIPs  {
        for src_port := port_strt; src_port <= port_end; src_port++ {
            wg.Add(1)
            limiter <- struct{}{}
            go func(src_ip string) {
                defer wg.Done() // Уменьшаем счетчик при завершении горутины
                defer func() { <-limiter }() // Когда горутина завершится, освобождаем место в канале
                err := sendRequestWithSourceIP(url_str, src_ip, src_port, page_size)
                if err != nil {
                    log.Printf("Error sending request from %s: %v\n", src_ip, err)
                }
            }(src_ip) // Передаем значение src_ip в горутину
            time.Sleep(duration) 
        }
    }
    wg.Wait()
}