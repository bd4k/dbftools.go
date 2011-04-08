package main
import(
    "os"
    "fmt"
    dbf "dbftools"
)

func main(){
    if len(os.Args)<2 { os.Exit(1)}
    
    f,e:=os.Open(os.Args[1],os.O_RDONLY,0666)
    if e!=nil{
        fmt.Println(e.String())
        os.Exit(1)
    }
    defer f.Close()
    reader,er:=dbf.NewReader(f,dbf.CP_866)
    if e!=nil{
        fmt.Println(er.String())
        os.Exit(1)
    }
    fs:= make(map[string]int,0)
    for x:=0;x<reader.FieldCount();x++{
        fs[reader.FieldName(x)]=x
    }
    
    var ok bool
    for ok,e=reader.Read();ok;ok,e=reader.Read(){
        for x:=0;x<reader.FieldCount();x++{
            fmt.Println(reader.FieldName(x),reader.String(x))
        }
        fmt.Println()
    }
    if e!=nil && e!=os.EOF { fmt.Println(e.String()) }
        
    
    
}    
