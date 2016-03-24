package main
 
import (
	"flag"

	"fmt"
	"net/http"
	"io/ioutil"
   
	//"net/url"
	//"log"
    //"github.com/opesun/goquery"
  
)

func main() {
/*
    x, err := goquery.ParseUrl(url_link)
    if err != nil {
        panic(err)
    }
*/
	var url_link = "https://vk.com/album19839792_000"
	
	resp, err := http.Get(url_link)
	if err != nil {
		//fmt.Printf(string(err.))
		//results, err = ioutil.ReadAll(resp.Body)
		//fmt.Printf("%s\r\n", string(results))
		// resp.Body.Close()
		
	}else{
		fmt.Printf("%s\r\n", err)
	}
	defer resp.Body.Close()
	//body, err := ioutil.ReadAll(resp.Body)

	results, err = ioutil.ReadAll(reader)
    if err != nil {
        panic(err)
    }
    reader.Close()
    println(string(results))
	 
	
	
/*
	u, err := url.Parse(url_link)
	if err != nil {
		log.Fatal(err)
	}
	u.Scheme = "https"
	u.Host = "vk.com"
	
	q := u.Query()
	q.Set("q", "golang")
	u.RawQuery = q.Encode()
	fmt.Println(u)
	*/
	

	
	// resp, err := http.Get("http://www.google.co.jp")
	// defer resp.Body.Close()
	// body, err := ioutil.ReadAll(resp.Body)
	
}
 


/*
-flag
-flag=x
-flag x  // non-boolean flags only

-h // all flag names list
*/

// var ip = flag.Int("flagname", 1234, "help message for flagname")
 
var (
    WORKERS       int             = 2            //кол-во "потоков"
    REPORT_PERIOD int             = 10           //частота отчетов (сек)
    DUP_TO_STOP   int             = 500          //максимум повторов до останова
    HASH_FILE     string          = "hash.bin"   //файл с хешами
    QUOTES_FILE   string          = "quotes.txt" //файл с цитатами
	
	// Int flags accept +- 1234, 0664, 0x1234 
	// bool_flag bool = false // 1, 0, t, f, T, F, true, false, TRUE, FALSE, True, False
)
func init() {
    //Задаем правила разбора:
    flag.IntVar(&WORKERS, "w", WORKERS, "количество потоков")
    flag.IntVar(&REPORT_PERIOD, "r", REPORT_PERIOD, "частота отчетов (сек)")
    flag.IntVar(&DUP_TO_STOP, "d", DUP_TO_STOP, "кол-во дубликатов для остановки")
    flag.StringVar(&HASH_FILE, "hf", HASH_FILE, "файл хешей")
    flag.StringVar(&QUOTES_FILE, "qf", QUOTES_FILE, "файл записей")
    
	
	// flag.BoolVar(&bool_flag, "t_f", bool_flag, "bool test")
	
}

/*
func main(){
	fmt.Printf("hello, world\n")
	
	//запускаем разбор аргументов
	flag.Parse() 
	
	
	// fmt.Println("ip: has value ", ip) // get pointer addres value
	// fmt.Println("ip: has value po ", *ip) // get value
	// fmt.Println("ip: has value ", *ip)
	// fmt.Println("Bool_flag: has value ", bool_flag)
	
	
	// fmt.Println("workers: has value ", WORKERS)
		
	
}
*/