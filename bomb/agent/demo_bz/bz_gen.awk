#! /usr/bin/awk -f
# 获取到所有的数据

BEGIN {
    LoadTypeAlias()
    ApiCount = 1
    PayLoadCount = 1
    LoadAPI()
    LoadPayLoad()
    # OutPutAPI()
    # OutPutPayLoad()
    OutPutBz()
}

# 类型别名.
function LoadTypeAlias() {
    TypeMap["uint32"] = "uint32"
    TypeMap["int32"] = "int32"

    TypeMap["uint16"] = "uint16"
    TypeMap["int16"] = "int16"

    TypeMap["string"] = "string"
    TypeMap["int"] = "int32"
    TypeMap["integer"] = "int32"
}


function FindType(Type) {
    if (TypeMap[Type] != "") {
        return TypeMap[Type]
    }else {
        return Type
    }
}



# 加载api.txt
function LoadAPI() {
    while( getline line < "api.txt" ) { #由指定的文件中读取测验数据
        if (line ~ /^#.*/ || line ~ /^\s*$/ ) {
            if (ApiList[ApiCount,"packet_type"] != "") {
                ApiCount+=1
            }
            continue
        }
        split(line,a,":")
        ApiList[ApiCount,a[1]] = a[2]
    }
}


# 加载payload.txt
function LoadPayLoad() {
    while( getline line < "proto.txt" ) { #由指定的文件中读取测验数据
        if (line ~ /^#.*/ || line ~ /^\s*$/ || line == "===" ) {
            continue
        }
        if(match(line,/^[^=].+=/) > 0 ) {
            name=substr(line,0,length(line)-1)
            PayLoadList[name,"count"] = 0
            PayLoadNames[PayLoadCount] = name
            PayLoadCount+=1
        } else {
            split(line,a," ")
            if (a[2] == "array") {
                fc = PayLoadList[name,"count"]+1
                PayLoadList[name,fc,"name"]=a[1]
                PayLoadList[name,fc,"type"]="array"
                PayLoadList[name,fc,"addtion"]=FindType(a[3])
                PayLoadList[name,"count"]= fc
            } else {
                fc = PayLoadList[name,"count"]+1
                PayLoadList[name,fc,"name"]=a[1]
                PayLoadList[name,fc,"type"]=FindType(a[2])
                PayLoadList[name,"count"]= fc
            }
        }
    }
}

# 数据PayLoad
function OutPutPayLoad() {
    for (i = 1; i< PayLoadCount; i++) {
        name = PayLoadNames[i]
        print i,name
    }
}

# 输出go的struct结构.
function OutPutPayLoadStruct(Name) {
    count=PayLoadList[Name,"count"]
    printf("type Pkt%s struct {\n", Name) > "demo.go"
    for(StructI=1; StructI<= count; StructI++) {
        name = PayLoadList[Name,StructI,"name"]
        type = PayLoadList[Name,StructI,"type"]
        addtion = PayLoadList[Name,StructI,"addtion"]
        if(type == "array") {
            printf("\t%s []%s\n", name,addtion) > "demo.go"
        } else {
            printf("\t%s %s\n", name,type) > "demo.go"
        }
    }
    print "}\n" > "demo.go"
}

# 输出所有的PayLoadStruct.
function OutPutAllPayLoadStruct() {
    for(PayLoadi =1; PayLoadi< PayLoadCount; PayLoadi ++) {
        name = PayLoadNames[PayLoadi]
        OutPutPayLoadStruct(name)
    }
}

# 输出所有的API常量定义
function OutPutAllApiConst() {
    print "const (" > "demo.go"
    for (APIi = 1; APIi < ApiCount; APIi++) {
        printf("\tBZ_%s = %s\n", toupper(ApiList[APIi, "name"]),
               ApiList[APIi, "packet_type"]) > "demo.go"
    }
    print ")\n" > "demo.go"
}


# 输出所有的MapHandler
function OutPutMapHandler(Name) {
    printf("func MakeBz%sHandler() BzHandlerMAP {\n",Name)  > "demo.go"
    print "\tProtocalHandler := BzHandlerMAP{" > "demo.go"
    for(APIi=1; APIi <=ApiCount; APIi++) {
        if (ApiList[APIi, "name"] ~ /Req$/) {
            printf("\t\tBZ_%s: Bz%s,\n", toupper(ApiList[APIi, "name"]),
                   ApiList[APIi, "name"]) > "demo.go"
        }
    }
    print "\t}\n\treturn ProtocalHandler\n}\n" > "demo.go"
}


function OutPutErrCheck(prefix) {
    print prefix, "\tif err != nil {"  > "demo.go"
    print prefix, "\t\treturn"  > "demo.go"
    print prefix, "\t}"  > "demo.go"
}

# 输出序列化:read
function OutPutSerializeRead(Name) {
    count=PayLoadList[Name,"count"]
    printf("func BzReadPkt%s(datai []byte) (data []byte, ret *Pkt%s, err error) {\n",
           Name,Name) > "demo.go"
    printf("\tdata = datai\n") > "demo.go"
    printf("\tret = &Pkt%s{}\n",Name) > "demo.go"
    for(Filedi = 1 ;Filedi<= count; Filedi++) {
        type = PayLoadList[Name,Filedi,"type"]
        addtion = PayLoadList[Name,Filedi,"addtion"]
        name = PayLoadList[Name,Filedi,"name"]

        if (type == "array") {
            printf("\tvar %s_v %s\n", name, addtion) > "demo.go"
            printf("\tdata, %s_size, err := BzReaduint16(data)\n", name) > "demo.go"
            printf("\tfor i := 0; i < int(%s_size); i++ {\n", name)  > "demo.go"
            printf("\t\tdata, %s_v, err = BzRead%s(data)\n",name,addtion) > "demo.go"
            OutPutErrCheck("\t")
            printf("\t\tret.%s = append(ret.%s, %s_v)\n",name,name,name) > "demo.go"
            printf("\t}\n") > "demo.go"

        } else {
            printf("\tdata, ret.%s, err = BzRead%s(data)\n",name, type ) > "demo.go"
        }
        OutPutErrCheck()
    }
    print "\treturn" > "demo.go"
    printf("}\n") > "demo.go"
}


# 输出序列化:write
function OutPutSerializeWrite(Name) {
    count=PayLoadList[Name,"count"]
    printf("func BzWritePkt%s(datai []byte, ret *Pkt%s) (data []byte, err error) {\n",
           Name,Name) > "demo.go"
    printf("\tdata = datai\n") > "demo.go"
    for(Filedi = 1 ;Filedi<= count; Filedi++) {
        type = PayLoadList[Name,Filedi,"type"]
        addtion = PayLoadList[Name,Filedi,"addtion"]
        name = PayLoadList[Name,Filedi,"name"]
        if (type == "array") {
            printf("\tdata, err = BzWriteuint16(data, uint16(len(ret.%s)))\n",
                   name) > "demo.go"
            printf("\tfor _, %s_v := range ret.%s {\n", name, name)  > "demo.go"
            printf("\t\tdata, err = BzWrite%s(data, %s_v)\n",addtion,name) > "demo.go"
            printf("\t}\n") > "demo.go"

        } else {
            printf("\tdata, err = BzWrite%s(data, ret.%s)\n",type, name) > "demo.go"
        }
    }
    print "\treturn" > "demo.go"
    printf("}\n") > "demo.go"
}


# 输出序列化函数.
function OutPutSerialize(Name) {
    OutPutSerializeRead(Name)
    OutPutSerializeWrite(Name)
}

# 输出所有的序列化函数.
function OutPutAllSerialize() {
    for (i = 1; i< PayLoadCount; i++) {
        Name = PayLoadNames[i]
        OutPutSerialize(Name)
    }
}

# 输出BzGo文件.
function OutPutBz() {
    print "package agent\n" > "demo.go"
    OutPutAllApiConst()
    OutPutAllPayLoadStruct()
    OutPutMapHandler("Gs")
    OutPutAllSerialize()
}
