#include <iostream>
#include <sys/socket.h>
#include <sys/types.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <netdb.h>
#include <unistd.h>

//以下头文件是为了使样例程序正常运行
#include <string.h>
#include <stdio.h>
#include <stdlib.h>


using namespace std;

#define PORT 8080
#define IP "127.0.0.1"


int main(int argc, char *argv[])
{

    struct sockaddr_in pin;
    struct hostent *nlp_host;
    int sd;
    char host_name[256];

    //初始化主机名和端口。主机名可以是IP，也可以是可被解析的名称
    strcpy(host_name,IP);

    //解析域名，如果是IP则不用解析，如果出错，显示错误信息
    while ((nlp_host=gethostbyname(host_name))==0){
        printf("Resolve Error!\n");
    }


    //设置pin变量，包括协议、地址、端口等，此段可直接复制到自己的程序中
    bzero(&pin,sizeof(pin));
    pin.sin_family=AF_INET;                 //AF_INET表示使用IPv4
    pin.sin_addr.s_addr=htonl(INADDR_ANY);
    pin.sin_addr.s_addr=((struct in_addr *)(nlp_host->h_addr))->s_addr;
    pin.sin_port=htons(PORT);

    //建立socket
    sd=socket(AF_INET,SOCK_STREAM,0);

    //建立连接
    while (connect(sd,(struct sockaddr*)&pin,sizeof(pin))==-1) {
        printf("Connect Error!\n");
    }


    return 0;
}
