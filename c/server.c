#include <stdio.h>
#include <unistd.h>
#include <sys/socket.h>
#include <netinet/tcp.h>
#include <netdb.h>
#include <netinet/in.h>
#include <strings.h>

int main( int argc, char *argv[] ) {
    int portno = 50007;

    int sockfd = socket(AF_INET, SOCK_STREAM, 0);

    struct sockaddr_in serv_addr;
    bzero((char *) &serv_addr, sizeof(serv_addr));

    serv_addr.sin_family = AF_INET;
    serv_addr.sin_addr.s_addr = INADDR_ANY;
    serv_addr.sin_port = htons(portno);

    bind(sockfd, (struct sockaddr *) &serv_addr, sizeof(serv_addr));
    listen(sockfd, 1);

    int fastopen_enable = 1;
    setsockopt(sockfd, IPPROTO_TCP, TCP_FASTOPEN, &fastopen_enable, sizeof(fastopen_enable)); 

    struct sockaddr_in cli_addr;
    unsigned int clilen = sizeof(cli_addr);
    char buffer[256];
    while(1) {
        int newsockfd = accept(sockfd, (struct sockaddr *)&cli_addr, &clilen);
        write(newsockfd,"abcdefg", 7);

        bzero(buffer,256);
        read( newsockfd,buffer,255 );
        printf("Received message: %s\n", buffer);
        close(newsockfd);
    }

    return 0;
}
