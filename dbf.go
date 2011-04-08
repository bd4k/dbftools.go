package dbftools
import(
    "os"
    "io"
    enc "encoding/binary"
)


// 1-st part of header
// common info 
// 32 bytes
type rawheader struct{
    TableType byte              
    LastUpdate [3] byte
    RecordCount uint32
    HeaderLength uint16
    RecordLength uint16
    Reserved [20] byte
}


// field info
// 32 bytes
type rawfield struct{
    Name [11]byte
    Type byte
    Offset uint32
    Length byte
    Dot byte
    Unused [14]byte
}


type Reader struct{
    reader io.Reader
    header *rawheader
    fields []rawfield
    fieldnames []string
    Options int
    tmp []byte
    row int
    enc int
}

func NewReader(source io.Reader, codepage int) (*Reader,os.Error){
    var e os.Error
    header:=&rawheader{}
    if e:=enc.Read( source, enc.LittleEndian, header);e!=nil{
        return nil,e
    }
        
    fcount:=(header.HeaderLength-33)/32
    fields:=make([]rawfield,fcount)
    if e=enc.Read(source, enc.LittleEndian, fields);e!=nil{ 
        return nil,e 
    }

    names:=make([]string,len(fields))
    for x,v:= range fields{
        names[x]=string(trimtrail(v.Name[:],0))
    }

    mark:=[1]byte{0}
    if _,e = source.Read(mark[:]); e!=nil{
        return nil,e
    }
    if mark[0]!=0xD {
        return nil,os.NewError("Bad header or some wrong")
    }

    return &Reader{ 
        reader:source, 
        header: header, 
        fields:fields,
        fieldnames:names, 
        tmp:make([]byte,header.RecordLength),
        enc:codepage }, nil
}

func (r *Reader) FieldCount() int{
    return len(r.fields)
}
func (r *Reader) RecordCount() uint32{
    return r.header.RecordCount
}
func (r *Reader) FieldName(ord int) (string){
    return r.fieldnames[ord]
}
func (r *Reader) FieldLen(ord int) (byte){
    return r.fields[ord].Length
}
func (r *Reader) RecordNo() (int){
    return r.row
}
func (r *Reader) Read() (bool, os.Error){
    r.row++
    _,e:=r.reader.Read(r.tmp)
    if e!=nil {return false, e}
    if r.tmp[0]==0x1A{return false,os.EOF}
    return true,nil
}

func (r *Reader) Bytes(field int) ([]byte){
    left:=r.fields[field].Offset
    right:=r.fields[field].Offset+uint32(r.fields[field].Length)
    f:=make([]byte,r.fields[field].Length)
    copy( f, r.tmp[left:right])
    return f
}

func (r *Reader) String(field int) (string){
    return Decode(trimtrail(r.Bytes(field),' '),r.enc)
}
func trimtrail(b []byte, v byte) []byte{
    for i:=len(b);i>0;{
        i--
        if b[i]!=v { return b[0:i+1]}
    }
    return b[0:0]
}
func trimlead(b []byte, v byte) []byte{
    for i:=0;i<len(b);{
        if b[i]!=v { return b[i:]}
    }
    return b[0:0]
}


