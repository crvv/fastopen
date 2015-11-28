#include <sys/socket.h>
#include <stdio.h>
#include <netdb.h>
#include <strings.h>
#include <unistd.h>

int main() {
    int socketfd = socket(AF_INET, SOCK_STREAM, 0);
    struct hostent* server = gethostbyname("nuc");

    struct sockaddr_in server_addr;
    bzero((char*)&server_addr, sizeof(server_addr));
    server_addr.sin_family = AF_INET;
    bcopy((char*)server->h_addr_list[0], (char*)&server_addr.sin_addr.s_addr, server->h_length);
    server_addr.sin_port = htons(50007);
    server_addr.sin_len = sizeof(server_addr);

    sa_endpoints_t endpoints;
    bzero((char*)&endpoints, sizeof(endpoints));
    endpoints.sae_dstaddr = (struct sockaddr*)&server_addr;
    endpoints.sae_dstaddrlen = sizeof(server_addr);

    int rc = connectx(socketfd,
            &endpoints,
            SAE_ASSOCID_ANY,
            CONNECT_RESUME_ON_READ_WRITE | CONNECT_DATA_IDEMPOTENT,
            NULL, 0, NULL, NULL);

    write(socketfd, "Hello, world", 12);
    char buffer[256];
    bzero(buffer, 256);
    read(socketfd, buffer, 255);
    printf("%s\n", buffer);
    close(socketfd);
    return 0;
}

